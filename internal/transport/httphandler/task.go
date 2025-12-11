package httphandler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/mkeOrt/tasks-go/internal/transport/dto"
	"github.com/mkeOrt/tasks-go/internal/transport/response"
)

// TaskService defines the business logic interface for tasks.
type TaskService interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
}

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	logger *slog.Logger
	svc    TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(logger *slog.Logger, svc TaskService) *TaskHandler {
	return &TaskHandler{
		logger: logger,
		svc:    svc,
	}
}

func (h *TaskHandler) RegisterRoutes() *http.ServeMux {
	g := http.NewServeMux()
	g.HandleFunc("GET /", h.GetAll)
	return g
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to get all tasks", slog.String("error", err.Error()))
		response.RespondWithError(w, err)
		return
	}
	dtos := dto.MapTasksToDTO(tasks)
	response.RespondWithJson(w, http.StatusOK, dto.TasksResponse{Tasks: dtos})
}
