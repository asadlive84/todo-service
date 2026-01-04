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

	return nil
}
