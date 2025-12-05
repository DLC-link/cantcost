package exporters

import (
	"context"

	"github.com/DLC-link/cantcost/internal/parser"
)

type Exporter interface {
	Export(ctx context.Context, line *parser.Line) error
}
