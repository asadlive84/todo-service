package usecase

import (
	"context"
	"errors"
	"fmt"
	"todo-service/internal/domain/entity"
	iface "todo-service/internal/domain/interface"

	"github.com/google/uuid"
)

const MaxFileSize = 5 << 20 // 5MB
const AllowedImageTypes = "image/"
const AllowedTextTypes = "text/"

type TodoUseCase struct {
	todoRepo  iface.TodoRepository
	s3Repo    iface.S3Repository
	redisRepo iface.RedisStreamRepository
	bucket    string
}

func NewTodoUseCase(todoRepo iface.TodoRepository, s3Repo iface.S3Repository, redisRepo iface.RedisStreamRepository, bucket string) *TodoUseCase {
	return &TodoUseCase{
		todoRepo:  todoRepo,
		s3Repo:    s3Repo,
		redisRepo: redisRepo,
		bucket:    bucket,
	}
}

func (uc *TodoUseCase) CreateTodo(ctx context.Context, todo *entity.TodoItem) (*entity.TodoItem, error) {
	if todo.Description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if todo.DueDate.IsZero() {
		return nil, errors.New("dueDate is required")
	}

	todo.ID = uuid.New().String()

	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)

	}
	if err := uc.redisRepo.PublishTodo(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to publish todo: %w", err)
	}

	return todo, nil
}
