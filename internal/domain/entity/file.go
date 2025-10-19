package entity

import (
	"time"
)

type File struct {
	ID           string    `json:"id"`
	FileName     string    `json:"fileName"`
	OriginalName string    `json:"originalName"`
	ContentType  string    `json:"contentType"`
	FileSize     int64     `json:"fileSize"`
	FileHash     string    `json:"fileHash,omitempty"`
	StoragePath  string    `json:"storagePath,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type FileUploadRequest struct {
	OriginalName string
	ContentType  string
	FileSize     int64
	UploadedBy   string
}

type FileUploadResponse struct {
	ID           string    `json:"id"`
	FileName     string    `json:"fileName"`
	OriginalName string    `json:"originalName"`
	ContentType  string    `json:"contentType"`
	FileSize     int64     `json:"fileSize"`
	FileHash     string    `json:"fileHash,omitempty"`
	UploadedAt   time.Time `json:"uploadedAt"`
}
