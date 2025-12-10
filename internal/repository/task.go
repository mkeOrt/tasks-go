package repository

import (
	"context"
	"database/sql"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

// TaskRepository provides access to task storage.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// GetAll retrieves all tasks from the database.
func (r *TaskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	q := "SELECT id, title, done, created_at, updated_at FROM tasks"
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if tasks == nil {
		return []domain.Task{}, nil
	}

	return tasks, nil
}
