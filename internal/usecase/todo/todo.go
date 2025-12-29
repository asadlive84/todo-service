package todo

import (
	"todo-service/internal/port"
)

const MaxFileSize = 5 << 20 // 5MB
const AllowedImageTypes = "image/"
const AllowedTextTypes = "text/"

type TodoUseCase struct {
	repo      port.TodoRepoPort
	s3Repo    port.S3Repository
	redisRepo port.RedisStreamRepository
	search    port.SearchRepo
	bucket    string
}

func NewTodoUseCase(repo port.TodoRepoPort, s3Repo port.S3Repository, redisRepo port.RedisStreamRepository, search port.SearchRepo, bucket string) *TodoUseCase {
	return &TodoUseCase{
		repo:   repo,
		s3Repo: s3Repo,
		// redisRepo: redisRepo,
		search: search,
		bucket: bucket,
	}
}
