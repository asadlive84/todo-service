package repository

import (
	"context"
	"fmt"
	"todo-service/internal/domain/entity"
	iface "todo-service/internal/domain/interface"

	"git.ice.global/packages/beeorm/v4"
)

type TodoRepository struct {
	orm *beeorm.Engine
}

func NewTodoRepository(orm *beeorm.Engine) iface.TodoRepository {
	return &TodoRepository{orm: orm}
}

func (r *TodoRepository) Create(ctx context.Context, todo *entity.TodoItemEntity) (err error) {
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
	fmt.Printf(":::todo::::%+v", todo)
	r.orm.Flush(todo)
	return nil
}

func (r *TodoRepository) GetTodoByID(ctx context.Context, id int) (en *entity.TodoItemEntity, err error) {
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

	en = &entity.TodoItemEntity{}

	r.orm.LoadByID(uint64(id), en)

	return en, nil
}

// func (r *TodoRepository) Create(ctx context.Context, todo *entity.TodoItem) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return fmt.Errorf("begin transaction error: %w", err)
// 	}
// 	defer tx.Rollback()

// 	query := `
// 		INSERT INTO todos (id, description, due_date, file_id)
// 		VALUES (?, ?, ?, ?)
// 	`

// 	stmt, err := tx.PrepareContext(ctx, query)
// 	if err != nil {
// 		return fmt.Errorf("prepare statement error: %w", err)
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.ExecContext(ctx, todo.ID, todo.Description, todo.DueDate, todo.FileID)
// 	if err != nil {
// 		return fmt.Errorf("execute error: %w", err)
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return fmt.Errorf("commit transaction error: %w", err)
// 	}

// 	return nil
// }
