package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config struct defines all the configurations needed
type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Azure struct {
		AccountName   string `yaml:"accountName"`
		AccountKey    string `yaml:"accountKey"`
		ContainerName string `yaml:"containerName"`
	} `yaml:"azure"`

	Kafka struct {
		Brokers       []string `yaml:"brokers"`
		ConsumerGroup string   `yaml:"consumerGroup"`
		Topics        struct {
			StorageEvents string `yaml:"storageEvents"`
		} `yaml:"topics"`
		Producer struct {
			RequiredAcks int `yaml:"requiredAcks"`
			Compression  int `yaml:"compression"`
			Retries      int `yaml:"retries"`
		} `yaml:"producer"`
	} `yaml:"kafka"`

	Logging struct {
		ElasticsearchURL string `yaml:"elasticsearchURL"`
	} `yaml:"logging"`
}

// LoadConfig reads the configuration from file
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}
