package app

import (
	"log/slog"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/config"
	"github.com/mkeOrt/tasks-go/internal/repository"
	"github.com/mkeOrt/tasks-go/internal/service"
	"github.com/mkeOrt/tasks-go/internal/transport/httphandler"
	"github.com/mkeOrt/tasks-go/internal/transport/middleware"
	"github.com/mkeOrt/tasks-go/pkg/database"
)

// Container centraliza las dependencias de la aplicación.
type Container struct {
	Handler http.Handler
	Cleanup func()
}

// NewContainer inicializa todas las dependencias y retorna el handler raíz y cleanup.
func NewContainer(cfg *config.Config, logger *slog.Logger) (*Container, error) {
	db, err := database.NewSqliteDB(cfg.DB.ConnectionString)
	if err != nil {
		return nil, err
	}

	repo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(logger.With("package", "task"), repo)
	taskHandler := httphandler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.Handle("/api/tasks", taskHandler.RegisterRoutes())

	handler := middleware.Logger(logger)(mux)
	handler = middleware.Cors(&cfg.Cors)(handler)

	cleanup := func() {
		db.Close()
	}

	return &Container{
		Handler: handler,
		Cleanup: cleanup,
	}, nil
}
