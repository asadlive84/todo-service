package bench

import (
	"context"
	"testing"
	"time"
	"todo-service/internal/domain/entity"
	"todo-service/internal/usecase"
	"todo-service/test/mock"

	"github.com/golang/mock/gomock"
)

func BenchmarkCreateTodo_DatabaseInsert(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Benchmark todo",
		DueDate:     time.Now().Add(24 * time.Hour),
		FileID:      "file-123",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := todoUC.CreateTodo(context.Background(), input)
		if err != nil {
			b.Fatalf("CreateTodo failed: %v", err)
		}
	}
}

func BenchmarkUploadFile_FileProcessing(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), "test-bucket", gomock.Any(), gomock.Any(), gomock.Any()).
		Return("file-key", nil).
		Times(b.N)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")

	fileContent := []byte("This is a test file for benchmarking upload operations")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := fileUC.UploadFile(
			context.Background(),
			"test.txt",
			fileContent,
			"text/plain",
		)
		if err != nil {
			b.Fatalf("UploadFile failed: %v", err)
		}
	}
}

func BenchmarkPublishTodo_RedisStream(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	input := &entity.TodoItem{
		Description: "Benchmark redis publish",
		DueDate:     time.Now().Add(48 * time.Hour),
		FileID:      "file-456",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := todoUC.CreateTodo(context.Background(), input)
		if err != nil {
			b.Fatalf("CreateTodo failed: %v", err)
		}
	}
}

func BenchmarkCreateTodo_WithMultipleTodos(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockTodoRepo := mock.NewMockTodoRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)
	mockRedisRepo := mock.NewMockRedisStreamRepository(ctrl)

	mockTodoRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	mockRedisRepo.EXPECT().
		PublishTodo(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	todoUC := usecase.NewTodoUseCase(mockTodoRepo, mockS3Repo, mockRedisRepo, "test-bucket")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		input := &entity.TodoItem{
			Description: "Batch todo",
			DueDate:     time.Now().Add(24 * time.Hour),
			FileID:      "file-batch",
		}
		_, err := todoUC.CreateTodo(context.Background(), input)
		if err != nil {
			b.Fatalf("CreateTodo failed: %v", err)
		}
	}
}

func BenchmarkUploadFile_LargeFile(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), "test-bucket", gomock.Any(), gomock.Any(), gomock.Any()).
		Return("file-key", nil).
		Times(b.N)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")

	// Create 1MB file
	largeFileContent := make([]byte, 1024*1024)
	for i := range largeFileContent {
		largeFileContent[i] = byte(i % 256)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := fileUC.UploadFile(
			context.Background(),
			"large.txt",
			largeFileContent,
			"text/plain",
		)
		if err != nil {
			b.Fatalf("UploadFile failed: %v", err)
		}
	}
}

func BenchmarkUploadFile_ParallelUploads(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), "test-bucket", gomock.Any(), gomock.Any(), gomock.Any()).
		Return("file-key", nil).
		Times(b.N)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(b.N)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")
	fileContent := []byte("Test file for parallel upload")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := fileUC.UploadFile(
				context.Background(),
				"test.txt",
				fileContent,
				"text/plain",
			)
			if err != nil {
				b.Fatalf("UploadFile failed: %v", err)
			}
		}
	})
}
