package repository

import (
	"context"
	"fmt"
	domain "todo-service/internal/domain/entity"
	"todo-service/internal/repository/beeorm/entity"
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

	ormEntity := entity.ToOrmEntity(todo)
	r.orm.Flush(ormEntity)

	return nil
}
