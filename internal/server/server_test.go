package server

import (
	"log/slog"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/mkeOrt/tasks-go/internal/config"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Addr:         ":9090",
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	}
	handler := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv := NewServer(cfg, handler, logger)

	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
	if srv.httpServer.Addr != ":9090" {
		t.Errorf("expected Addr :9090, got %s", srv.httpServer.Addr)
	}
	if srv.httpServer.ReadTimeout != 1*time.Second {
		t.Errorf("expected ReadTimeout 1s, got %v", srv.httpServer.ReadTimeout)
	}
}

func TestServer_Run_GracefulShutdown(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Addr:         ":0",
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
		},
	}
	handler := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	srv := NewServer(cfg, handler, logger)

	errChan := make(chan error, 1)

	go func() {
		errChan <- srv.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find current process: %v", err)
	}
	if err := p.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send SIGTERM: %v", err)
	}

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("Run() returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not shut down within timeout")
	}
}
