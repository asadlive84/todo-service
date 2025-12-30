package entity

import (
	"errors"
	"fmt"
	"time"
)

type File struct {
	ID           int64
	FileName     string
	OriginalName string
	ContentType  string
	FileContent  []byte
	FileSize     int64
	FileHash     string
	StoragePath  string
	CreatedAt    time.Time
	MaxFileSize  int64
}

func (f *File) Validate() error {
	if int64(len(f.FileContent)) > f.MaxFileSize {
		return fmt.Errorf("file too large: max size is %d bytes", f.MaxFileSize)
	}

	if !f.isAllowedContentType(f.ContentType) {
		return errors.New("invalid file type: only image/* and text/* are allowed")
	}

	return nil

}

func (f *File) isAllowedContentType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
		"text/plain":      true,
		"text/csv":        true,
		"application/pdf": true,
	}
	return allowedTypes[contentType]
}
