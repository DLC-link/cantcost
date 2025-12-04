package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/DLC-link/cantcost/internal/catcher"
	"github.com/DLC-link/cantcost/internal/env"
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
	ctx := context.Background()
	err := catcher.Stream(ctx, func(ctx context.Context, line string) error {
		if strings.Contains(strings.ToLower(line), "eventcost") {
			parsedLine, err := parser.ProcessLine(line)
			if err != nil {
				slog.ErrorContext(ctx, "Failed to parse log line", slog.Any("error", err))
				return err
			}
			fmt.Printf("%s TraceID: %s, EventCostDetails: TotalCost=%d, GroupToMembersSize=%v\n",
				parsedLine.Timestamp.Format(time.TimeOnly),
				parsedLine.TraceID,
				parsedLine.CostDetails.EventCost,
				parsedLine.CostDetails.GroupToMembersSize)
			fmt.Println("----------------------------")
		}
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "Log streaming failed", slog.Any("error", err))
	}
}
