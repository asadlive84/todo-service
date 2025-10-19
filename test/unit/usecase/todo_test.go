package usecase

import (
	"context"
	"errors"
	"testing"
	"time"
	"todo-service/internal/domain/entity"
	"todo-service/internal/usecase"
	"todo-service/test/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func TestCreateTodo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Learn GoMock",
		DueDate:     time.Now().Add(48 * time.Hour),
		FileID:      "file-789",
	}

	result, err := todoUC.CreateTodo(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.FileID, result.FileID)
	assert.NotEmpty(t, result.ID)
}

func TestCreateTodo_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errors.New("database connection error")).
		Times(1)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Times(0)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Test error handling",
		DueDate:     time.Now().Add(24 * time.Hour),
		FileID:      "file-789",
	}

	result, err := todoUC.CreateTodo(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create todo")
}

func TestCreateTodo_RedisError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(errors.New("redis connection error")).
		Times(1)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Test Redis error",
		DueDate:     time.Now().Add(24 * time.Hour),
		FileID:      "file-789",
	}

	result, err := todoUC.CreateTodo(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to publish todo")
}

func TestCreateTodo_WithFileID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Todo with attachment",
		DueDate:     time.Now().Add(72 * time.Hour),
		FileID:      "file-12345",
	}

	result, err := todoUC.CreateTodo(context.Background(), input)

	require.NoError(t, err)
	assert.Equal(t, "file-12345", result.FileID)
}

func TestCreateTodo_EmptyFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Todo without file",
		DueDate:     time.Now().Add(24 * time.Hour),
		FileID:      "",
	}

	result, err := todoUC.CreateTodo(context.Background(), input)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Description, result.Description)
	assert.Equal(t, input.DueDate, result.DueDate)
}
