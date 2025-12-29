package port

import (
	"context"
	"io"
	"todo-service/internal/domain/entity"

)

type FileRepository interface {
	CreateFile(ctx context.Context, file *entity.File) error
	GetFileByID(ctx context.Context, id string) (*entity.File, error)
}

type S3Repository interface {
	UploadFile(ctx context.Context, bucket, key string, reader io.Reader, size int64) (string, error)
}

type RedisStreamRepository interface {
	PublishTodo(ctx context.Context, todo *entity.TodoItem) error
}
