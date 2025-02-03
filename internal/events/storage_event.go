package events

import (
	"encoding/json"
	"time"
)

// EventType represents different types of storage events
type EventType string

const (
	FileUploaded     EventType = "FileUploaded"
	FileDeleted      EventType = "FileDeleted"
	FileAppended     EventType = "FileAppended"
	DirectoryCreated EventType = "DirectoryCreated"
	DirectoryDeleted EventType = "DirectoryDeleted"
)

// StorageEvent
type StorageEvent struct {
	ID        string            `json:"id"`
	Type      EventType         `json:"type"`
	Path      string            `json:"path"`
	Size      int64             `json:"size,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	UserID    string            `json:"userId,omitempty"`
	MetaData  map[string]string `json:"metadata,omitempty"`
}

// 🔹 Convert
func (e *StorageEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// 🔹 Parse JSON
func FromJSON(data []byte) (*StorageEvent, error) {
	var event StorageEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	// metadata
	if event.MetaData == nil {
		event.MetaData = make(map[string]string)
	}
	return &event, nil
}

// EventListener
type EventListener interface {
	Subscribe(topics []string) error
	Close()
}

// EventPublisher
type EventPublisher interface {
	Publish(topic string, event *StorageEvent) error
	Close()
}
