package httphandler

import (
	"context"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/mkeOrt/tasks-go/internal/dto"
)

// TaskService defines the business logic interface for tasks.
type TaskService interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
}

// TaskHandler handles HTTP requests for tasks.
type TaskHandler struct {
	svc TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) RegisterRoutes() *http.ServeMux {
	g := http.NewServeMux()
	g.HandleFunc("GET /", h.GetAll)
	return g
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll(r.Context())
	if err != nil {
		statusCode := MapErrorToStatusCode(err)
		var errorMsg string

		if statusCode == http.StatusInternalServerError {
			errorMsg = "Ocurrió un error interno al procesar la solicitud"
		} else {
			errorMsg = err.Error()
		}
		RespondWithErrorJson(w, statusCode, errorMsg)
		return
	}
	dtos := dto.MapTasksToDTO(tasks)
	RespondWithJson(w, http.StatusOK, dto.TasksResponse{Tasks: dtos})
}
