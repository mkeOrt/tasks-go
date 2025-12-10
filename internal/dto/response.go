package dto

import "github.com/mkeOrt/tasks-go/internal/domain"

type TasksResponse struct {
	Tasks []domain.Task `json:"tasks"`
}
