package env

import (
	"log/slog"
	"os"
)

const (
	targetDeployment = "TARGET_DEPLOYMENT"
	targetContainer  = "TARGET_CONTAINER"
	targetNamespace  = "TARGET_NAMESPACE"

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

func Print() {
	slog.Info("Environment Variables")
	slog.Info("TARGET_DEPLOYMENT", slog.String("value", GetTargetDeployment()))
	slog.Info("TARGET_CONTAINER", slog.String("value", GetTargetContainer()))
	slog.Info("TARGET_NAMESPACE", slog.String("value", GetTargetNamespace()))
}
