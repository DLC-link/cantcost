package exporters

import (
	"context"

	"github.com/DLC-link/cantcost/internal/parser"
)

type Exporter interface {
	Export(ctx context.Context, line *parser.Line) error
}

type Exporters struct {
	exporters []Exporter
}

func New(exporters ...Exporter) *Exporters {
	return &Exporters{
		exporters: exporters,
	}
}

func (e *Exporters) Export(ctx context.Context, line *parser.Line) error {
	for _, exporter := range e.exporters {
		if err := exporter.Export(ctx, line); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporters) AddExporter(exporter Exporter) {
	e.exporters = append(e.exporters, exporter)
}
