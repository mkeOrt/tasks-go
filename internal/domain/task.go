package domain

import (
	"context"
	"time"
)

type Task struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TaskRepository interface {
	GetAll(ctx context.Context) ([]Task, error)
}
