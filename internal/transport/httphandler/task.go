package httphandler

import (
	"context"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/mkeOrt/tasks-go/internal/dto"
)

type TaskService interface {
	GetAll(ctx context.Context) ([]domain.Task, error)
}

type TaskHandler struct {
	svc TaskService
}

func NewTaskHandler(svc TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) RegisterRoutes() *http.ServeMux {
	g := http.NewServeMux()
	g.HandleFunc("/api/tasks", h.GetAll)
	return g
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAll(r.Context())
	if err != nil {
		RespondWithErrorJson(w, MapErrorToStatusCode(err), err.Error())
		return
	}
	RespondWithJson(w, http.StatusOK, dto.TasksResponse{Tasks: tasks})
}
