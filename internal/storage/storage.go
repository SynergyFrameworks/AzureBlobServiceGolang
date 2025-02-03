package storage

import "context"

// StorageAdapter defines an interface for storage operations.
type StorageAdapter interface {
	UploadFile(ctx context.Context, filePath string, data []byte) error
	WriteFile(ctx context.Context, path string, content []byte, overwrite bool) error
	ReadFile(ctx context.Context, filePath string) ([]byte, error)
	DeleteFile(ctx context.Context, filePath string) error
	ListFiles(ctx context.Context, dirPath string) ([]string, error)
}
