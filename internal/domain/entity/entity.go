package entity

import (
	"time"
)

type TodoItem struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}

type CreateTodoInput struct {
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}
