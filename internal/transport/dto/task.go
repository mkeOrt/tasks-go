package dto

import "github.com/mkeOrt/tasks-go/internal/domain"

// TaskDTO is a data transfer object for Task.
type TaskDTO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// TasksResponse is the response for a list of tasks.
type TasksResponse struct {
	Tasks []TaskDTO `json:"tasks"`
}

// MapTasksToDTO maps domain tasks to DTOs.
func MapTasksToDTO(tasks []domain.Task) []TaskDTO {
	dtos := make([]TaskDTO, len(tasks))
	for i, t := range tasks {
		dtos[i] = TaskDTO{
			ID:        t.ID,
			Title:     t.Title,
			Done:      t.Done,
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return dtos
}
