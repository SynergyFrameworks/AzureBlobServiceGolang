package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type AzureStorage struct {
	AccountName   string
	AccountKey    string
	ContainerName string
	client        *azblob.Client
}

var _ StorageAdapter = (*AzureStorage)(nil)

// NewAzureStorage initializes an Azure Storage client.
func NewAzureStorage(accountName, accountKey, containerName string) (*AzureStorage, error) {
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)

	// ✅ Fix: Use SharedKeyCredential instead of a raw key string
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create shared key credential: %v", err)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Storage client: %v", err)
	}

	return &AzureStorage{
		AccountName:   accountName,
		AccountKey:    accountKey,
		ContainerName: containerName,
		client:        client,
	}, nil
}

// UploadFile
func (s *AzureStorage) UploadFile(ctx context.Context, filePath string, data []byte) error {
	return s.WriteFile(ctx, filePath, data, false)
}

// WriteFile
func (s *AzureStorage) WriteFile(ctx context.Context, path string, content []byte, overwrite bool) error {
	blobClient := s.client.ServiceClient().NewContainerClient(s.ContainerName).NewBlockBlobClient(path)

	_, err := blobClient.UploadStream(ctx, bytes.NewReader(content), &blockblob.UploadStreamOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file to Azure Storage: %v", err)
	}
	return nil
}

// ReadFile
func (s *AzureStorage) ReadFile(ctx context.Context, filePath string) ([]byte, error) {
	blobClient := s.client.ServiceClient().NewContainerClient(s.ContainerName).NewBlobClient(filePath)

	response, err := blobClient.DownloadStream(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from Azure Storage: %v", err)
	}

	// Read content
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return data, nil
}

// DeleteFile
func (s *AzureStorage) DeleteFile(ctx context.Context, filePath string) error {
	blobClient := s.client.ServiceClient().NewContainerClient(s.ContainerName).NewBlobClient(filePath)
	_, err := blobClient.Delete(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete file from Azure Storage: %v", err)
	}
	return nil
}

// ListFiles
func (s *AzureStorage) ListFiles(ctx context.Context, dirPath string) ([]string, error) {
	containerClient := s.client.ServiceClient().NewContainerClient(s.ContainerName)

	pager := containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{})

	files := []string{}
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %v", err)
		}

		for _, blob := range resp.Segment.BlobItems {
			files = append(files, *blob.Name)
		}
	}
	return files, nil
}
