package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"
	"todo-service/internal/domain/entity"
	"todo-service/internal/infrastructure/helper"
	iface "todo-service/internal/port"

	"github.com/google/uuid"
)

const MaxFileSize = 5

type FileUseCase struct {
	fileRepo iface.FileRepository
	s3Repo   iface.S3Repository
	bucket   string
}

func NewFileUseCase(fileRepo iface.FileRepository, s3Repo iface.S3Repository, bucket string) *FileUseCase {
	return &FileUseCase{
		fileRepo: fileRepo,
		s3Repo:   s3Repo,
		bucket:   bucket,
	}
}

func (uc *FileUseCase) UploadFile(ctx context.Context, originalName string, fileContent []byte, contentType string) (*entity.FileUploadResponse, error) {
	if int64(len(fileContent)) > MaxFileSize {
		return nil, fmt.Errorf("file too large: max size is %d bytes", MaxFileSize)
	}

	if !isAllowedContentType(contentType) {
		return nil, errors.New("invalid file type: only image/* and text/* are allowed")
	}

	fileHash := helper.ComputeFileHash(fileContent)

	fileID := uuid.New().String()
	ext := filepath.Ext(originalName)
	storageKey := fileID + ext

	_, err := uc.s3Repo.UploadFile(ctx, uc.bucket, storageKey, bytes.NewReader(fileContent), int64(len(fileContent)))
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to S3: %w", err)
	}

	now := time.Now()
	fileMetadata := &entity.File{
		ID:           fileID,
		FileName:     storageKey,
		OriginalName: originalName,
		ContentType:  contentType,
		FileSize:     int64(len(fileContent)),
		FileHash:     fileHash,
		StoragePath:  storageKey,
		CreatedAt:    now,
	}

	if err := uc.fileRepo.CreateFile(ctx, fileMetadata); err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return &entity.FileUploadResponse{
		ID:           fileID,
		FileName:     storageKey,
		OriginalName: originalName,
		ContentType:  contentType,
		FileSize:     int64(len(fileContent)),
		FileHash:     fileHash,
		UploadedAt:   now,
	}, nil
}

func (uc *FileUseCase) ValidateFileID(ctx context.Context, fileID string) (bool, error) {
	file, err := uc.fileRepo.GetFileByID(ctx, fileID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve file: %w", err)
	}
	return file != nil, nil
}

func isAllowedContentType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"text/plain":      true,
		"text/csv":        true,
		"application/pdf": true,
	}
	return allowedTypes[contentType]
}
