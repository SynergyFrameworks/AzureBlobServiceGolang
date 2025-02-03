package logger

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"project-root/config"
)

var elasticsearchURL string

// Initialize Elasticsearch URL from config
func init() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to load config, using default Elasticsearch URL")
		elasticsearchURL = "http://localhost:9200/logs/_doc/"
		return
	}

	if cfg.Logging.ElasticsearchURL == "" {
		log.Printf("⚠️ Warning: Elasticsearch URL not set in config, using default")
		elasticsearchURL = "http://localhost:9200/logs/_doc/"
	} else {
		elasticsearchURL = cfg.Logging.ElasticsearchURL
	}
}

// 🔹 Log Message Structure
type LogMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// 🔹 Send logs to Elasticsearch
func sendToElasticSearch(logData LogMessage) {
	jsonData, _ := json.Marshal(logData)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(elasticsearchURL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Printf("❌ Failed to send log to Elasticsearch: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("❌ Elasticsearch returned error: %s", resp.Status)
	}
}

// 🔹 Log an informational message
func LogInfo(message string) {
	logData := LogMessage{Level: "INFO", Message: message}
	log.Println("ℹ️", message)
	sendToElasticSearch(logData)
}

// 🔹 Log an error message
func LogError(message string, err error) {
	logData := LogMessage{Level: "ERROR", Message: message, Error: err.Error()}
	log.Printf("❌ %s: %v", message, err)
	sendToElasticSearch(logData)
}
