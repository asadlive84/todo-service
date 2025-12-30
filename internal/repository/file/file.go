package file

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"todo-service/internal/domain/entity"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) CreateFile(ctx context.Context, file *entity.File) error {

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO files (id, file_name, original_name, content_type, file_size, file_hash, storage_path, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare statement error: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, file.ID, file.FileName, file.OriginalName, file.ContentType, file.FileSize, file.FileHash, file.StoragePath, file.CreatedAt)
	if err != nil {
		return fmt.Errorf("execute error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}

func (r *FileRepository) GetFileByID(ctx context.Context, id string) (*entity.File, error) {
	query := `
		SELECT id, file_name, original_name, content_type, file_size, file_hash, storage_path, created_at
		FROM files
		WHERE id = ?
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare statement error: %w", err)
	}
	defer stmt.Close()

	file := &entity.File{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&file.ID,
		&file.FileName,
		&file.OriginalName,
		&file.ContentType,
		&file.FileSize,
		&file.FileHash,
		&file.StoragePath,
		&file.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("file not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return file, nil
}
