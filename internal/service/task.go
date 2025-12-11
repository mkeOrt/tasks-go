package service

import (
	"context"
	"fmt"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

// TaskService provides business logic for tasks.
type TaskService struct {
	repo domain.TaskRepository
}

// NewTaskService creates a new TaskService.
func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

// GetAll returns all tasks from the repository.
func (s *TaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("TaskService.GetAll: %w: %w", domain.ErrTaskRetrievalFailed, err)
	}
	return tasks, nil
}
