package exporters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/DLC-link/cantcost/internal/parser"
)

var _ Exporter = (*HTTP)(nil)

type HTTP struct {
	URL                 string `json:"url"`
	AuthorizationHeader string `json:"authorization_header"`
	BatchSize           int    `json:"batch_size,omitempty"`

	batch []*parser.Line
	mutex *sync.Mutex
}

type HTTPRequest struct {
	Count int                   `json:"count"`
	Lines []*parser.MessageLine `json:"lines"`
}

func NewHTTPExporter(url string, authHeader string, batchSize int) *HTTP {
	if batchSize <= 0 {
		batchSize = 10
	}
	return &HTTP{
		URL:                 url,
		AuthorizationHeader: authHeader,
		BatchSize:           batchSize,
		batch:               make([]*parser.Line, 0),
		mutex:               &sync.Mutex{},
	}
}

func (h *HTTP) Export(ctx context.Context, line *parser.Line) error {
	var request = HTTPRequest{
		Count: 1,
		Lines: []*parser.MessageLine{line.ToMessageLine()},
	}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", h.AuthorizationHeader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnprocessableEntity {
			slog.ErrorContext(ctx, "HTTP exporter received 422 Unprocessable Entity. Check if the log line format matches the expected schema.", slog.String("data", string(data)))
		}
		return errors.New("failed to export log line, status code: " + resp.Status)
	}

	return nil
}
