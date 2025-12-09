package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/DLC-link/cantcost/internal/catcher"
	"github.com/DLC-link/cantcost/internal/env"
	"github.com/DLC-link/cantcost/internal/exporters"
	"github.com/DLC-link/cantcost/internal/parser"
	slogcontext "github.com/PumpkinSeed/slog-context"
)

func main() {
	slog.SetDefault(
		slog.New(
			slogcontext.NewHandler(
				slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: env.GetLogLevel(),
				}),
			),
		),
	)
	env.Print()

	var exporter = exporters.New()
	if env.GetExporterType() == "http" {
		httpExporter := exporters.NewHTTPExporter(
			env.GetHTTPExporterURL(),
			env.GetHTTPExporterAuthHeader(),
			env.GetHTTPExporterBatchSize(),
		)
		exporter.AddExporter(httpExporter)
		slog.Info("HTTP exporter configured",
			slog.String("url", env.GetHTTPExporterURL()),
		)
	}

	ctx := context.Background()
	err := catcher.Stream(ctx, func(ctx context.Context, line string) error {
		if strings.Contains(strings.ToLower(line), "eventcost") {
			parsedLine, err := parser.ProcessLine(line)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to parse log line", slog.Any("error", err))
				return err
			}
			if err := exporter.Export(ctx, &parsedLine); err != nil {
				slog.ErrorContext(ctx, "Failed to export parsed line", slog.Any("error", err))
				return err
			}
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "Log streaming failed", slog.Any("error", err))
	}
}
