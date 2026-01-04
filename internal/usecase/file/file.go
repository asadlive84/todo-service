package usecase

import (
	iface "todo-service/internal/domain/interface"
)

const MaxFileSize = 50000

type FileUseCase struct {
	fileRepo iface.FileRepoPort
	todoRepo iface.TodoRepoPort
	s3Repo   iface.S3Repository
	bucket   string
}

func NewFileUseCase(fileRepo iface.FileRepoPort, todoRepo iface.TodoRepoPort, s3Repo iface.S3Repository, bucket string) *FileUseCase {
	return &FileUseCase{
		fileRepo: fileRepo,
		todoRepo: todoRepo,
		s3Repo:   s3Repo,
		bucket:   bucket,
	}
}
