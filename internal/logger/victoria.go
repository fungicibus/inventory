package logger

import (
	"bytes"
	"fmt"
	"net/http"
)

type victoriaLogsWriter struct {
	endpoint string
}

func NewVictoriaLogsWriter(endpoint string) *victoriaLogsWriter {
	return &victoriaLogsWriter{endpoint: endpoint}
}

func (w victoriaLogsWriter) Write(p []byte) (n int, err error) {
	resp, err := http.Post(w.endpoint, "application/json", bytes.NewBuffer(p))
	if err != nil {
		return 0, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	return len(p), nil
}
