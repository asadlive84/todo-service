package usecase

import (
	"context"
	"errors"
	"testing"
	"todo-service/internal/usecase"
	"todo-service/test/mock"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestUploadFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	fileContent := []byte("test file content")
	fileName := "test.txt"
	contentType := "text/plain"

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), "test-bucket", gomock.Any(), gomock.Any(), int64(len(fileContent))).
		Return("file-key", nil).
		Times(1)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")

	response, err := fileUC.UploadFile(context.Background(), fileName, fileContent, contentType)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, fileName, response.OriginalName)
	assert.Equal(t, contentType, response.ContentType)
}

func TestUploadFile_S3Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	fileContent := []byte("test")

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), "test-bucket", gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", errors.New("s3 service unavailable")).
		Times(1)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Times(0)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")

	response, err := fileUC.UploadFile(context.Background(), "test.txt", fileContent, "text/plain")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to upload file")
}

func TestUploadFile_FileTooLarge(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFileRepo := mock.NewMockFileRepository(ctrl)
	mockS3Repo := mock.NewMockS3Repository(ctrl)

	largeFileContent := make([]byte, 5*1024*1024+1)

	mockS3Repo.EXPECT().
		UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	mockFileRepo.EXPECT().
		CreateFile(gomock.Any(), gomock.Any()).
		Times(0)

	fileUC := usecase.NewFileUseCase(mockFileRepo, mockS3Repo, "test-bucket")

	response, err := fileUC.UploadFile(context.Background(), "large.txt", largeFileContent, "text/plain")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "file too large")
}
