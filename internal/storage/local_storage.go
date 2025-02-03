package storage

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// LocalStorage is a local file system storage adapter.
type LocalStorage struct {
	BasePath string
}

// Ensure LocalStorage satisfies StorageAdapter.
var _ StorageAdapter = (*LocalStorage)(nil)

// NewLocalStorage initializes local storage.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

// UploadFile writes a file to local storage (alias for WriteFile).
func (s *LocalStorage) UploadFile(ctx context.Context, filePath string, data []byte) error {
	return s.WriteFile(ctx, filePath, data, false)
}

// WriteFile writes data to a file with an overwrite option.
func (s *LocalStorage) WriteFile(ctx context.Context, path string, content []byte, overwrite bool) error {
	fullPath := filepath.Join(s.BasePath, path)

	// Ensure the directory exists.
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	// Prevent overwrite if not allowed.
	if !overwrite {
		if _, err := os.Stat(fullPath); err == nil {
			return fmt.Errorf("file already exists and overwrite is disabled: %s", path)
		}
	}

	// Write the file.
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	return nil
}

// ReadFile retrieves the content of a file.
func (s *LocalStorage) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	fullPath := filepath.Join(s.BasePath, filePath)

	// Check if the file exists.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	// Read the file content.
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return data, nil
}

// DeleteFile removes a file from local storage.
func (s *LocalStorage) DeleteFile(ctx context.Context, filePath string) error {
	fullPath := filepath.Join(s.BasePath, filePath)

	// Check if the file exists.
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Delete the file.
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// ListFiles lists all files in a directory.
func (s *LocalStorage) ListFiles(ctx context.Context, dirPath string) ([]string, error) {
	fullPath := filepath.Join(s.BasePath, dirPath)
	files := []string{}

	err := filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(s.BasePath, path)
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	return files, nil
}
