package usecase

import (
	"context"
	"errors"
	"fmt"
	"todo-service/internal/domain/entity"
	iface "todo-service/internal/domain/interface"
	// "github.com/google/uuid"
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
		todoRepo: todoRepo,
		s3Repo:   s3Repo,
		// redisRepo: redisRepo,
		bucket: bucket,
	}
}

func (uc *TodoUseCase) CreateTodo(ctx context.Context, todo *entity.TodoItemEntity) (*entity.TodoItemEntity, error) {
	if todo.Description == "" {
		return nil, errors.New("description cannot be empty")
	}
	if todo.DueDate.IsZero() {
		return nil, errors.New("dueDate is required")
	}

	if err := uc.todoRepo.Create(ctx, todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)

	}
	// if err := uc.redisRepo.PublishTodo(ctx, todo); err != nil {
	// 	return nil, fmt.Errorf("failed to publish todo: %w", err)
	// }

	return todo, nil
}

func (uc *TodoUseCase) GetTodoByID(ctx context.Context, id int) (*entity.TodoItemEntity, error) {

	if id == 0 {
		return nil, fmt.Errorf("error id is %+v", id)
	}

	entity, err := uc.todoRepo.GetTodoByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)

	}
	// if err := uc.redisRepo.PublishTodo(ctx, todo); err != nil {
	// 	return nil, fmt.Errorf("failed to publish todo: %w", err)
	// }

	return entity, nil
}
