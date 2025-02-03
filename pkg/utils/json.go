package utils

import (
	"encoding/json"
	"fmt"
)

// 🔹 Convert an object to JSON
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// 🔹 Convert an object to a pretty-printed JSON string
func ToPrettyJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to serialize JSON: %v", err)
	}
	return string(data), nil
}

// 🔹 Deserialize JSON into an object
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
