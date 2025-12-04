package catcher

import (
	"bufio"
	"context"
	"io"
	"log/slog"

	"github.com/DLC-link/cantcost/internal/env"
	slogcontext "github.com/PumpkinSeed/slog-context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Stream(ctx context.Context, lineHandler func(context.Context, string) error) error {
	clientSet, err := getKubernetesClient(ctx)
	if err != nil {
		return err
	}

	ctx = slogcontext.WithValue(ctx, "target_deployment", env.GetTargetDeployment())
	ctx = slogcontext.WithValue(ctx, "target_namespace", env.GetTargetNamespace())

	pod, err := getPod(ctx, clientSet)
	if err != nil {
		return err
	}

	stream, err := getLogStream(ctx, clientSet, pod)
	if err != nil {
		return err
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		if err := lineHandler(ctx, line); err != nil {
			slog.ErrorContext(ctx, "Error handling log line", slog.Any("error", err))
			continue
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		slog.ErrorContext(ctx, "Error reading log stream", slog.Any("error", err))
		return err
	}

	slog.InfoContext(ctx, "Log stream ended")
	return nil
}

func getKubernetesClient(ctx context.Context) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create in-cluster config", slog.Any("error", err))
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create clientSet", slog.Any("error", err))
		return nil, err
	}

	return clientSet, nil
}

// Pick the first Running pod if possible, otherwise just the first pod.
func selectFirstRunningPod(list *corev1.PodList) corev1.Pod {
	for _, p := range list.Items {
		if p.Status.Phase == corev1.PodRunning {
			return p
		}
	}
	// fallback
	return list.Items[0]
}

func getPod(ctx context.Context, clientSet *kubernetes.Clientset) (corev1.Pod, error) {
	// 1) Get the Deployment
	deploy, err := clientSet.AppsV1().
		Deployments(env.GetTargetNamespace()).
		Get(ctx, env.GetTargetDeployment(), metav1.GetOptions{})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get deployment", slog.Any("error", err))
		return corev1.Pod{}, err
	}

	// 2) Build label selector from Deployment spec
	podSelector := deploy.Spec.Selector
	if podSelector == nil {
		slog.ErrorContext(ctx, "Deployment has no selector")
		return corev1.Pod{}, ErrNoSelector
	}
	labelSelector := metav1.FormatLabelSelector(podSelector)

	// 3) List pods matching the selector
	podList, err := clientSet.CoreV1().
		Pods(env.GetTargetNamespace()).
		List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to list pods", slog.Any("error", err))
		return corev1.Pod{}, err
	}

	if len(podList.Items) == 0 {
		slog.ErrorContext(ctx, "No pods found for deployment", slog.String("selector", labelSelector))
		return corev1.Pod{}, ErrNoPods
	}

	// For simplicity, pick the first pod
	pod := selectFirstRunningPod(podList)
	slog.InfoContext(ctx, "Selected pod for log streaming", slog.String("pod_name", pod.Name))

	return pod, nil
}

func getLogStream(ctx context.Context, clientSet *kubernetes.Clientset, pod corev1.Pod) (io.ReadCloser, error) {
	// Stream logs from the chosen pod
	logOptions := &corev1.PodLogOptions{
		Follow:     true,
		Timestamps: true,
	}
	if env.GetTargetContainer() != "" {
		logOptions.Container = env.GetTargetContainer()
	}

	req := clientSet.CoreV1().
		Pods(env.GetTargetNamespace()).
		GetLogs(pod.Name, logOptions)

	stream, err := req.Stream(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error opening log stream for pod",
			slog.String("pod_name", pod.Name), slog.Any("error", err))
		return nil, err
	}

	return stream, nil
}
