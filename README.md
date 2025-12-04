# CantCost

The tool exports the event costs from the Canton participant nodes and process them.

### Work in progress

This tool only logs the costs to the standard output. I'm working on other exporters.

### Important

The tool heavily relies on the Kubernetes API because it gets the logs from the pod directly running the Canton participant node.

### Installation

```shell
docker build -t {LOCATION_WHERE_THE_CLUSTER_CAN_PULL_FROM}/cantcost:latest .
docker push {LOCATION_WHERE_THE_CLUSTER_CAN_PULL_FROM}/cantcost:latest
```

After the image push you should change the values in the zarf/deployment/devnet/manifest.yaml file.

- Change the image location in the `spec.containers.image` field.
- Change the namespace everywhere for your desired namespace.
- Change the TARGET_DEPLOYMENT environment variable to the deployment name of your Canton participant node.

```shell
kubectl apply -f zarf/deployment/devnet/manifest.yaml
```

### Project layout

- internal/catcher: Setups a pod log streamer and call the callback to process a log line one by one.
- internal/parser: Parses the log lines and extract the cost events. This is the tricky part, because the log lines are Scala object serialized and wrapped into structured JSON logging.
- bin/main.go: The main entry point of the application. Everything glues together here. You can change the export logic here in the callback function.
