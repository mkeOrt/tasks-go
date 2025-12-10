package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/mkeOrt/tasks-go/internal/config"
	"github.com/mkeOrt/tasks-go/internal/repository"
	"github.com/mkeOrt/tasks-go/internal/server"
	"github.com/mkeOrt/tasks-go/internal/service"
	transport "github.com/mkeOrt/tasks-go/internal/transport/httphandler"
	"github.com/mkeOrt/tasks-go/internal/transport/middleware"
	"github.com/mkeOrt/tasks-go/pkg/database"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg := config.NewConfig(logger)

	db, err := database.NewSqliteDB(cfg.DB.ConnectionString)
	if err != nil {
		logger.Error("failed to create database", "error", err)
		os.Exit(1)
	}
	logger.Info("database created")

	defer db.Close()

	mux := http.NewServeMux()

	repo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(logger.With("package", "task"), repo)
	taskHandler := transport.NewTaskHandler(taskService)

	mux.Handle("/api/tasks", taskHandler.RegisterRoutes())

	// Wrap the mux with the logging middleware
	handler := middleware.Logger(logger)(mux)
	handler = middleware.Cors(&cfg.Cors)(handler)

	srv := server.NewServer(cfg, handler, logger)
	if err := srv.Run(); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
