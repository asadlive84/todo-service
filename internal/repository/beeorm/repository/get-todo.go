package repository

import (
	"context"
	"fmt"

	domain "todo-service/internal/domain/entity"
	"todo-service/internal/repository/beeorm/entity"
	"todo-service/internal/repository/beeorm/mapper"
)

func (r *OrmEngine) GetByID(ctx context.Context, id int) (todo *domain.TodoItem, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("beeorm flush panic: %v", r)
		}
	}()

	en := &entity.TodoEntity{}

	r.orm.LoadByID(uint64(id), en)

	return mapper.ToDomainEntity(en), nil
}
