package entity

import (
	"time"

	domain "todo-service/internal/domain/entity"

	"git.ice.global/packages/beeorm/v4"
)

type TodoEntity struct {
	beeorm.ORM  `orm:"table=todos;redisCache"`
	ID          uint64 `orm:"pk"`
	Description string
	DueDate     time.Time `orm:"time"`
	CreatedAt   time.Time `orm:"time"`
	FileID      string
}

// Mapper: domain → BeeORM
func ToOrmEntity(todo *domain.TodoItem) *TodoEntity {

	return &TodoEntity{
		// ID:          uint64(todo.ID),
		Description: todo.Description,
		DueDate:     todo.DueDate,
		FileID:      todo.FileID,
		CreatedAt:   todo.CreatedAt,
		
	}
}

// Mapper: BeeORM → domain
func ToDomainEntity(ormTodo *TodoEntity) *domain.TodoItem {
	return &domain.TodoItem{
		ID:          int(ormTodo.ID),
		Description: ormTodo.Description,
		DueDate:     ormTodo.DueDate,
		FileID:      ormTodo.FileID,
		CreatedAt:   ormTodo.CreatedAt,
	}
}
