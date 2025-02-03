package storage

import (
	"context"
	"fmt"
	"sync"
)

type MockAzureStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMockAzureStorage() *MockAzureStorage {
	return &MockAzureStorage{
		data: make(map[string][]byte),
	}
}

func (s *MockAzureStorage) UploadFile(ctx context.Context, filePath string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[filePath] = data
	return nil
}

func (s *MockAzureStorage) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, exists := s.data[filePath]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	return data, nil
}

func (s *MockAzureStorage) DeleteFile(ctx context.Context, filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[filePath]; !exists {
		return fmt.Errorf("file not found: %s", filePath)
	}
	delete(s.data, filePath)
	return nil
}

func (s *MockAzureStorage) ListFiles(ctx context.Context, dirPath string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var files []string
	for key := range s.data {
		files = append(files, key)
	}
	return files, nil
}
