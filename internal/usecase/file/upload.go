package usecase

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"time"
	"todo-service/internal/domain/entity"
	"todo-service/internal/helper"

	"github.com/google/uuid"
)

// upload s3 and create file info into MYSQL
func (uc *FileUseCase) UploadFile(ctx context.Context, file *entity.File) error {

	file.MaxFileSize = MaxFileSize

	if err := file.Validate(); err != nil {
		return err
	}
	fileHash := helper.ComputeFileHash(file.FileContent)

	file.FileHash = fileHash

	fileID := uuid.New().String()
	ext := filepath.Ext(file.OriginalName)
	storageKey := fileID + ext

	_, err := uc.s3Repo.UploadFile(ctx, uc.bucket, storageKey, bytes.NewReader(file.FileContent), int64(len(file.FileContent)))
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	now := time.Now()

	file.CreatedAt = now
	// fileMetadata := &entity.File{
	// 	ID:           fileID,
	// 	FileName:     storageKey,
	// 	OriginalName: originalName,
	// 	ContentType:  contentType,
	// 	FileSize:     int64(len(fileContent)),
	// 	FileHash:     fileHash,
	// 	StoragePath:  storageKey,
	// 	CreatedAt:    now,
	// }

	if err := uc.todoRepo.CreateFile(ctx, file); err != nil {
		return fmt.Errorf("failed to save file metadata: %w", err)
	}

	return nil
}

func (uc *FileUseCase) ValidateFileID(ctx context.Context, fileID string) (bool, error) {
	file, err := uc.fileRepo.GetFileByID(ctx, fileID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve file: %w", err)
	}
	return file != nil, nil
}
