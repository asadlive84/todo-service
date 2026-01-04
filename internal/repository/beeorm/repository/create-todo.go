package repository

import (
	"context"
	"fmt"
	domain "todo-service/internal/domain/entity"
	"todo-service/internal/repository/beeorm/mapper"
)

func (r *OrmEngine) Create(ctx context.Context, todo *domain.TodoItem) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("beeorm flush panic: %v", r)
		}
	}()

	ormEntity := mapper.ToOrmEntity(todo)

	r.orm.Flush(ormEntity)

	todo.ID = int(ormEntity.ID)

	return nil
}
