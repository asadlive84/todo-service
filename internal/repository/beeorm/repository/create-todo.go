package repository

import (
	"context"
	"fmt"
	domain "todo-service/internal/domain/entity"
	"todo-service/internal/repository/beeorm/mapper"

	"git.ice.global/packages/hitrix/pkg/helper"
)

func (r *OrmEngine) Create(ctx context.Context, todo *domain.TodoItem) error {
	// Validate
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context error: %w", err)
	}
	if todo == nil {
		return fmt.Errorf("todo cannot be nil")
	}

	return helper.DBTransaction(r.orm, func() error {
		// Save Todo
		ormEntity := mapper.ToOrmTodoEntity(todo)
		if err := r.orm.FlushWithCheck(ormEntity); err != nil {
			return fmt.Errorf("failed to save todo: %w", err)
		}

		// Validate ID
		if ormEntity.ID == 0 {
			return fmt.Errorf("todo ID not generated")
		}

		todo.ID = int(ormEntity.ID)

		// Check context
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled: %w", err)
		}

		// Save Outbox
		todoOutbox := mapper.ToOrmOutboxEntity(todo)
		if err := r.orm.FlushWithCheck(todoOutbox); err != nil {
			return fmt.Errorf("failed to save outbox: %w", err)
		}

		return nil
	})
}
