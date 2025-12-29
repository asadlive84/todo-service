package port

import (
	"context"
	"time"
	domain "todo-service/internal/domain/entity"
)

type TodoUseCasePort interface {
	Create(ctx context.Context, todo *domain.TodoItem) error
	GetByID(ctx context.Context, id int) (*domain.TodoItem, error)
	Search(ctx context.Context, query string, offset int32, limit int32) ([]*domain.TodoItem, int64, error)
}

type TodoRepoPort interface {
	Create(ctx context.Context, todo *domain.TodoItem) error
	GetByID(ctx context.Context, id int) (*domain.TodoItem, error)
}

type SearchRepo interface {
	CreateTodoIndex(ctx context.Context) error
	IndexTodo(ctx context.Context, todoID uint64, description, fileID string, dueDate, createdAt time.Time) error
	SearchTodos(ctx context.Context, query string, offset, limit int) ([]map[string]interface{}, int64, error)
	// ParseSearchResults(result interface{}) ([]map[string]interface{}, int64, error)
}
