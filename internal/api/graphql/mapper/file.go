// internal/adapter/api/graphql/mapper/file_mapper.go
package mapper

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strconv"
	"time"
	"todo-service/internal/domain/dto"
	"todo-service/internal/domain/entity"

	"github.com/99designs/gqlgen/graphql"
)

// UploadToFileEntity converts GraphQL Upload to File domain entity
func UploadToFileEntity(upload graphql.Upload) (*entity.File, error) {
	// Generate unique filename (you can customize this logic)
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), upload.Filename)

	// Calculate file hash
	hash := sha256.New()
	if _, err := io.Copy(hash, upload.File); err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}
	fileHash := fmt.Sprintf("%x", hash.Sum(nil))

	// Reset file reader for later use
	if _, err := upload.File.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to reset file reader: %w", err)
	}

	return &entity.File{
		// ID:           "", // Will be set after DB insert
		FileName:     uniqueFilename,
		OriginalName: upload.Filename,
		ContentType:  upload.ContentType,
		FileSize:     upload.Size,
		FileHash:     fileHash,
		StoragePath:  "", // Will be set after S3 upload
		CreatedAt:    time.Now().UTC(),
	}, nil
}

// FileEntityToResponse converts File entity to upload response
func FileEntityToResponse(file *entity.File) *dto.FileUploadResponse {

	strconvs := strconv.FormatInt(file.ID, 10)

	return &dto.FileUploadResponse{
		ID:           strconvs,
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		ContentType:  file.ContentType,
		// FileSize:     file.FileSize,
		// FileHash:     file.FileHash,
		// UploadedAt:   file.CreatedAt,
	}
}
