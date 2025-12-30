package todo

import (
	"context"
	"fmt"
	"todo-service/internal/domain/entity"
)

func (uc *TodoUseCase) Create(ctx context.Context, todo *entity.TodoItem) error {

	if err := todo.Validate(); err != nil {
		return err
	}

	if err := uc.repo.Create(ctx, todo); err != nil {
		return fmt.Errorf("failed to create todo: %w", err)

	}

	if err := uc.redisRepo.PublishTodo(ctx, todo); err != nil {
		return fmt.Errorf("failed to publish todo: %w", err)
	}

	err := uc.search.IndexTodo(ctx, uint64(todo.ID), todo.Description, todo.FileID, todo.DueDate, todo.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to index redis search todo: %w", err)
	}

	return nil
}
