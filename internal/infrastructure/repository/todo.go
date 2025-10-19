package repository

import (
	"context"
	"database/sql"
	"fmt"
	"todo-service/internal/domain/entity"
	iface "todo-service/internal/domain/interface"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) iface.TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *entity.TodoItem) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO todos (id, description, due_date, file_id)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, todo.ID, todo.Description, todo.DueDate, todo.FileID)
	if err != nil {
		return fmt.Errorf("execute error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}
