package storage_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"project-root/internal/storage"

)

// 🔹 Test Azure Storage Write File
func TestAzureStorageWriteFile(t *testing.T) {
	azureStorage, err := storage.NewAzureStorage("testAccount", "testKey", "testContainer")
	if err != nil {
		t.Fatalf("❌ Failed to create AzureStorage: %v", err)
	}

	err = azureStorage.WriteFile(context.Background(), "test-blob", []byte("test data"), true)
	if err != nil {
		t.Errorf("❌ Failed to write file to Azure: %v", err)
	}
}

// 🔹 Test Local Storage Write File
func TestLocalStorageWriteFile(t *testing.T) {
	basePath := os.TempDir()
	localStorage := storage.NewLocalStorage(basePath)

	testFilePath := "testfile.txt"
	testData := []byte("test data")

	// Write file
	err := localStorage.WriteFile(context.Background(), testFilePath, testData, true)
	if err != nil {
		t.Errorf("❌ Failed to write file: %v", err)
	}

	// Verify file exists
	fullPath := filepath.Join(basePath, testFilePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("❌ File was not written correctly: %s", fullPath)
	}

	// Clean up
	defer os.Remove(fullPath)
}

// 🔹 Test Local Storage Read File
func TestLocalStorageReadFile(t *testing.T) {
	basePath := os.TempDir()
	localStorage := storage.NewLocalStorage(basePath)

	testFilePath := "testfile_read.txt"
	testData := []byte("read test")

	// Write file
	err := localStorage.WriteFile(context.Background(), testFilePath, testData, true)
	if err != nil {
		t.Fatalf("❌ Failed to write file for read test: %v", err)
	}

	// Read file
	data, err := localStorage.ReadFile(context.Background(), testFilePath)
	if err != nil {
		t.Errorf("❌ Failed to read file: %v", err)
	} else if string(data) != string(testData) {
		t.Errorf("❌ Read data mismatch. Expected: %s, Got: %s", testData, data)
	}

	// Clean up
	defer os.Remove(filepath.Join(basePath, testFilePath))
}

// 🔹 Test Local Storage Delete File
func TestLocalStorageDeleteFile(t *testing.T) {
	basePath := os.TempDir()
	localStorage := storage.NewLocalStorage(basePath)

	testFilePath := "testfile_delete.txt"
	testData := []byte("delete test")

	// Write file
	err := localStorage.WriteFile(context.Background(), testFilePath, testData, true)
	if err != nil {
		t.Fatalf("❌ Failed to write file for delete test: %v", err)
	}

	// Delete file
	err = localStorage.DeleteFile(context.Background(), testFilePath)
	if err != nil {
		t.Errorf("❌ Failed to delete file: %v", err)
	}

	// Verify deletion
	fullPath := filepath.Join(basePath, testFilePath)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Errorf("❌ File was not deleted correctly: %s", fullPath)
	}
}

// 🔹 Test Local Storage Handle Missing File
func TestLocalStorageReadMissingFile(t *testing.T) {
	basePath := os.TempDir()
	localStorage := storage.NewLocalStorage(basePath)

	_, err := localStorage.ReadFile(context.Background(), "nonexistent.txt")
	if err == nil {
		t.Errorf("❌ Expected error when reading missing file, got nil")
	}
}
