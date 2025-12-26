package entity

import (
	"time"

	"git.ice.global/packages/beeorm/v4"
)

// type TodoItem struct {
// 	ID          string    `json:"id"`
// 	Description string    `json:"description"`
// 	DueDate     time.Time `json:"dueDate"`
// 	FileID      string    `json:"fileId,omitempty"`
// }

type TodoItemEntity struct {
	beeorm.ORM  `orm:"table=todos;redisCache"`
	ID          uint64 `orm:"pk"`
	Description string
	DueDate     time.Time `orm:"time"`
	CreatedAt   time.Time `orm:"time"`
	FileID      string
}

type CreateTodoInput struct {
	Description string    `json:"description"`
	DueDate     time.Time `json:"dueDate"`
	FileID      string    `json:"fileId,omitempty"`
}
