package entity

import (
	"errors"
	"time"
)

type TodoItem struct {
	ID          int
	Description string
	DueDate     time.Time
	FileID      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (todo *TodoItem) Validate() error {
	if todo.Description == "" {
		return errors.New("description cannot be empty")
	}
	if todo.DueDate.IsZero() {
		return errors.New("dueDate is required")
	}

	return nil
}
