package repository

import (
	"context"
	"fmt"
	fileEntity "todo-service/internal/domain/entity"
	"todo-service/internal/repository/beeorm/entity"
)

func (r *OrmEngine) CreateFile(ctx context.Context, file *fileEntity.File) (err error) {

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

	ormFileEntity := entity.ToFileEnityOrmEntity(file)

	r.orm.Flush(ormFileEntity)

	file.ID = int64(ormFileEntity.ID)

	return nil
}
