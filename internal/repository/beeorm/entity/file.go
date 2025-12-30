package entity

import (
	"time"
	domain "todo-service/internal/domain/entity"

	"git.ice.global/packages/beeorm/v4"
)

type FileEnity struct {
	beeorm.ORM   `orm:"table=files;redisCache"`
	ID           uint64 `orm:"pk"`
	FileName     string
	OriginalName string
	ContentType  string
	FileSize     int64
	FileHash     string
	StoragePath  string
	CreatedAt    time.Time `orm:"time"`
}

// Mapper: domain → BeeORM
func ToFileEnityOrmEntity(file *domain.File) *FileEnity {

	now := time.Now().UTC()
	return &FileEnity{
		FileName:     file.FileName,
		OriginalName: file.OriginalName,
		ContentType:  file.ContentType,
		FileSize:     file.FileSize,
		FileHash:     file.FileHash,
		StoragePath:  file.StoragePath,
		CreatedAt:    now,
	}
}
