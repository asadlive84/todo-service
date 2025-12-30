package usecase

import (
	iface "todo-service/internal/port"
)

const MaxFileSize = 50000

type FileUseCase struct {
	fileRepo iface.FileRepoPort
	s3Repo   iface.S3Repository
	bucket   string
}

func NewFileUseCase(fileRepo iface.FileRepoPort, s3Repo iface.S3Repository, bucket string) *FileUseCase {
	return &FileUseCase{
		fileRepo: fileRepo,
		s3Repo:   s3Repo,
		bucket:   bucket,
	}
}
