package todo

import (
	"context"
	"fmt"
	"todo-service/internal/domain/entity"
)

func (uc *TodoUseCase) GetByID(ctx context.Context, id int) (*entity.TodoItem, error) {

	if id == 0 {
		return nil, fmt.Errorf("error id is %+v", id)
	}

	entity, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)

	}
	// if err := uc.redisRepo.PublishTodo(ctx, todo); err != nil {
	// 	return nil, fmt.Errorf("failed to publish todo: %w", err)
	// }

	return entity, nil
}
