package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	// Capture logs in a buffer
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	// Create a dummy handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))
	})

	// Wrap with middleware
	middleware := Logger(logger)
	handler := middleware(nextHandler)

	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check response
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	// Check logs
	logOutput := buf.String()
	if !strings.Contains(logOutput, "request completed") {
		t.Error("expected log to contain 'request completed'")
	}
	if !strings.Contains(logOutput, "method=POST") {
		t.Error("expected log to contain 'method=POST'")
	}
	if !strings.Contains(logOutput, "path=/test") {
		t.Error("expected log to contain 'path=/test'")
	}
	if !strings.Contains(logOutput, "status=201") {
		t.Error("expected log to contain 'status=201'")
	}
}
