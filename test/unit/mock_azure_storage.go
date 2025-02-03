package storage_test

import (
	"context"
	"testing"

	"project-root/internal/storage"

)

// 🔹 Test Mock Azure Storage Upload
func TestMockAzureStorageUpload(t *testing.T) {
	mockStorage := storage.NewMockAzureStorage()

	err := mockStorage.UploadFile(context.Background(), "test-blob", []byte("mock data"))
	if err != nil {
		t.Fatalf("❌ Failed to upload file: %v", err)
	}

	// Read back to verify
	data, err := mockStorage.ReadFile(context.Background(),  "test-blob")
	if err != nil {
		t.Fatalf("❌ Failed to read uploaded file: %v", err)
	}
	if string(data) != "mock data" {
		t.Errorf("❌ Data mismatch. Expected 'mock data', got '%s'", string(data))
	}
}

// 🔹 Test Mock Azure Storage Read Non-Existent File
func TestMockAzureStorageReadMissingFile(t *testing.T) {
	mockStorage := storage.NewMockAzureStorage()

	_, err := mockStorage.ReadFile(context.Background(),  "missing-blob")
	if err == nil {
		t.Errorf("❌ Expected error for missing file, got nil")
	}
}

// 🔹 Test Mock Azure Storage Delete File
func TestMockAzureStorageDeleteFile(t *testing.T) {
	mockStorage := storage.NewMockAzureStorage()

	// Upload first
	mockStorage.UploadFile(context.Background(),  "test-blob", []byte("mock data"))

	// Delete file
	err := mockStorage.DeleteFile(context.Background(),  "test-blob")
	if err != nil {
		t.Fatalf("❌ Failed to delete file: %v", err)
	}

	// Verify deletion
	_, err = mockStorage.ReadFile(context.Background(),  "test-blob")
	if err == nil {
		t.Errorf("❌ Expected error when reading deleted file, but got none")
	}
}
