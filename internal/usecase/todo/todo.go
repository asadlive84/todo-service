package todo

import (
	port "todo-service/internal/domain/interface"
)

const MaxFileSize = 5 << 20 // 5MB
const AllowedImageTypes = "image/"
const AllowedTextTypes = "text/"

type TodoUseCase struct {
	repo   port.TodoRepoPort
	s3Repo port.S3Repository
	search port.SearchRepo
	bucket string
}

func NewTodoUseCase(repo port.TodoRepoPort, s3Repo port.S3Repository, search port.SearchRepo, bucket string) *TodoUseCase {
	return &TodoUseCase{
		repo:   repo,
		s3Repo: s3Repo,
		search: search,
		bucket: bucket,
	}
}
