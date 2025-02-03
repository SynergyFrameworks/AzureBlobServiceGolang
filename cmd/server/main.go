package main

import (
	"fmt"
	"log"

	"project-root/config"
	"project-root/internal/api"
	"project-root/internal/kafka"
	"project-root/internal/storage"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	var storageAdapter storage.StorageAdapter
	if cfg.Azure.AccountName != "" && cfg.Azure.AccountKey != "" {
		storageAdapter, err = storage.NewAzureStorage(cfg.Azure.AccountName, cfg.Azure.AccountKey, cfg.Azure.ContainerName)
		if err != nil {
			log.Fatalf("Failed to initialize Azure Storage: %v", err)
		}
		log.Println("Using Azure Storage")
	} else {
		storageAdapter = storage.NewLocalStorage("./local_data")
		log.Println("Azure credentials missing, using Local Storage")
	}

	kafkaClient, err := kafka.NewKafkaClient(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer kafkaClient.Close()

	apiInstance := &api.API{
		Storage: storageAdapter,
		Kafka:   kafkaClient,
	}

	r := api.SetupRoutes(apiInstance)
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	r.Run(serverAddr)
}
