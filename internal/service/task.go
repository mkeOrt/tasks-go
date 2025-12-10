package service

import (
	"context"
	"log/slog"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

type TaskService struct {
	logger *slog.Logger
	repo   domain.TaskRepository
}

func NewTaskService(logger *slog.Logger, repo domain.TaskRepository) *TaskService {
	return &TaskService{
		logger: logger,
		repo:   repo,
	}
}

func (s *TaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all tasks", "error", err)
		return nil, domain.ErrTasksRetrieveError
	}
	return products, nil
}
