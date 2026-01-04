package port

import (
	"context"
	"io"
	"todo-service/internal/domain/entity"
	domain "todo-service/internal/domain/entity"
)

type TodoUseCasePort interface {
	Create(ctx context.Context, todo *domain.TodoItem) error
	GetByID(ctx context.Context, id int) (*domain.TodoItem, error)
	Search(ctx context.Context, query string, offset int32, limit int32) ([]*domain.TodoItem, int64, error)
}

type FileUseCasePort interface {
	// CreateFile(ctx context.Context, file *entity.File) error
	// GetFileByID(ctx context.Context, id string) (*entity.File, error)
	ValidateFileID(ctx context.Context, fileID string) (bool, error)
	UploadFile(ctx context.Context, file *entity.File) error
}

//==========

type TodoRepoPort interface {
	Create(ctx context.Context, todo *domain.TodoItem) error
	GetByID(ctx context.Context, id int) (*domain.TodoItem, error)
	CreateFile(ctx context.Context, file *domain.File) error
}

type FileRepoPort interface {
	GetFileByID(ctx context.Context, id string) (*entity.File, error)
}

type SearchRepo interface {
	// CreateTodoIndex(ctx context.Context) error
	// IndexTodo(ctx context.Context, todoID uint64, description, fileID string, dueDate, createdAt time.Time) error
	SearchTodos(ctx context.Context, query string, offset, limit int) ([]*entity.TodoItem, error)
	// ParseSearchResults(result interface{}) ([]map[string]interface{}, int64, error)
}

type S3Repository interface {
	UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64) (string, error)
}

type RedisStreamRepository interface {
	PublishTodo(ctx context.Context, todo *entity.TodoItem) error
}
