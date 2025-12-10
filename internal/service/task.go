package service

import (
	"context"
	"log/slog"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

// TaskService provides business logic for tasks.
type TaskService struct {
	logger *slog.Logger
	repo   domain.TaskRepository
}

// NewTaskService creates a new TaskService.
func NewTaskService(logger *slog.Logger, repo domain.TaskRepository) *TaskService {
	return &TaskService{
		logger: logger,
		repo:   repo,
	}
}

// GetAll returns all tasks from the repository.
func (s *TaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all tasks", "error", err)
		return nil, domain.ErrTasksRetrieveError
	}
	return tasks, nil
}
