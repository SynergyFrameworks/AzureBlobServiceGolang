package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"project-root/internal/events"
	"project-root/internal/kafka"
	"project-root/internal/storage"
)

type API struct {
	Storage storage.StorageAdapter // Exported (uppercase S)
	Kafka   *kafka.KafkaClient     // Exported (uppercase K)
}

// 🔹 Upload File Handler
func (api *API) uploadFile(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file: " + err.Error()})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	defer fileContent.Close()

	content, err := ioutil.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	overwrite := c.DefaultQuery("overwrite", "false") == "true"
	err = api.Storage.WriteFile(c.Request.Context(), path, content, overwrite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Storage write failed: " + err.Error()})
		return
	}

	api.publishEvent(events.FileUploaded, path, int64(len(content)), map[string]string{
		"filename":    file.Filename,
		"contentType": file.Header.Get("Content-Type"),
		"overwrite":   fmt.Sprintf("%v", overwrite),
	})

	c.Status(http.StatusCreated)
}

// 🔹 Delete File Handler
func (api *API) deleteFile(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	err := api.Storage.DeleteFile(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or cannot be deleted: " + err.Error()})
		return
	}

	api.publishEvent(events.FileDeleted, path, 0, nil)
	c.Status(http.StatusNoContent)
}

// 🔹 Create Directory Handler
func (api *API) createDirectory(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	// Create an empty directory (depends on the storage adapter)
	err := api.Storage.WriteFile(c.Request.Context(), path+"/.keep", []byte{}, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory: " + err.Error()})
		return
	}

	api.publishEvent(events.DirectoryCreated, path, 0, nil)
	c.Status(http.StatusCreated)
}

// 🔹 Delete Directory Handler
func (api *API) deleteDirectory(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	err := api.Storage.DeleteFile(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Directory not found or cannot be deleted: " + err.Error()})
		return
	}

	api.publishEvent(events.DirectoryDeleted, path, 0, nil)
	c.Status(http.StatusNoContent)
}

// 🔹 Read File Handler
func (api *API) readFile(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path parameter is required"})
		return
	}

	content, err := api.Storage.ReadFile(c.Request.Context(), path)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found: " + err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", content)
}

// 🔹 List Files Handler
func (api *API) listFiles(c *gin.Context) {
	dirPath := c.Param("path")
	if dirPath == "" {
		dirPath = "." // Root directory by default
	}

	files, err := api.Storage.ListFiles(c.Request.Context(), dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// 🔹 Publish Event to Kafka
func (api *API) publishEvent(eventType events.EventType, path string, size int64, metadata map[string]string) {
	event := &events.StorageEvent{
		Type:     eventType,
		Path:     path,
		Size:     size,
		MetaData: metadata,
	}

	if err := api.Kafka.Publish("storage-events", event); err != nil {
		fmt.Printf("Failed to publish event %s: %v\n", eventType, err)
	}
}
