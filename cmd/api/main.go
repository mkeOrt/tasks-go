package main

import (
	"log/slog"
	"os"

	"github.com/mkeOrt/tasks-go/internal/app"
	"github.com/mkeOrt/tasks-go/internal/config"
	"github.com/mkeOrt/tasks-go/internal/server"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg := config.NewConfig(logger)

	container, err := app.NewContainer(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}
	defer container.Cleanup()

	srv := server.NewServer(cfg, container.Handler, logger)
	if err := srv.Run(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
