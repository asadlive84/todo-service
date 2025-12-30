package dto

type FileUploadRequest struct {
	OriginalName string
	ContentType  string
	FileSize     int64
	UploadedBy   string
}
