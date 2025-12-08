package env

import (
	"log/slog"
	"os"
	"strconv"
)

const (
	targetDeployment       = "TARGET_DEPLOYMENT"
	targetContainer        = "TARGET_CONTAINER"
	targetNamespace        = "TARGET_NAMESPACE"
	exporterType           = "EXPORTER_TYPE"
	httpExporterURL        = "HTTP_EXPORTER_URL"
	httpExporterAuthHeader = "HTTP_EXPORTER_AUTH_HEADER"
	httpExporterBatchSize  = "HTTP_EXPORTER_BATCH_SIZE"

	logLevel = "LOG_LEVEL"
)

func GetLogLevel() slog.Level {
	if v := os.Getenv(logLevel); v != "" {
		switch v {
		case "debug":
			return slog.LevelDebug
		case "info":
			return slog.LevelInfo
		case "warn":
			return slog.LevelWarn
		case "error":
			return slog.LevelError
		default:
			return slog.LevelInfo
		}
	}
	return slog.LevelInfo
}

func GetTargetDeployment() string {
	if v := os.Getenv(targetDeployment); v != "" {
		return v
	}
	return "" // TODO
	//return DefaultReceiverPartyID
}

func GetTargetContainer() string {
	if v := os.Getenv(targetContainer); v != "" {
		return v
	}
	return ""
}

func GetTargetNamespace() string {
	if v := os.Getenv(targetNamespace); v != "" {
		return v
	}
	return "default"
}

func GetExporterType() string {
	if v := os.Getenv(exporterType); v != "" {
		return v
	}
	return "http"
}

func GetHTTPExporterURL() string {
	if v := os.Getenv(httpExporterURL); v != "" {
		return v
	}
	return ""
}

func GetHTTPExporterAuthHeader() string {
	if v := os.Getenv(httpExporterAuthHeader); v != "" {
		return v
	}
	return ""
}

func GetHTTPExporterBatchSize() int {
	if v := os.Getenv(httpExporterBatchSize); v != "" {
		strconvV, err := strconv.Atoi(v)
		if err == nil {
			return strconvV
		}
	}
	return 10
}

func Print() {
	slog.Info("Environment Variables")
	slog.Info("TARGET_DEPLOYMENT", slog.String("value", GetTargetDeployment()))
	slog.Info("TARGET_CONTAINER", slog.String("value", GetTargetContainer()))
	slog.Info("TARGET_NAMESPACE", slog.String("value", GetTargetNamespace()))
}
