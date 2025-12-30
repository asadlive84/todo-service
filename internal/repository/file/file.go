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
