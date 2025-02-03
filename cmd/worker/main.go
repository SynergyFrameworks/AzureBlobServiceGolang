package main

import (
	"context"
	"log"

	"project-root/config"
	"project-root/internal/events"
	"project-root/internal/kafka"
	"project-root/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize storage adapter (Azure or Local)
	var storageAdapter storage.StorageAdapter
	storageAdapter, err = storage.NewAzureStorage(cfg.Azure.AccountName, cfg.Azure.AccountKey, cfg.Azure.ContainerName)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize Kafka client
	kafkaClient, err := kafka.NewKafkaClient(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer kafkaClient.Close()

	// Register event handlers
	kafkaClient.RegisterHandler(events.FileUploaded, func(event *events.StorageEvent) {
		log.Printf("📥 Processing uploaded file: %s", event.Path)

		content := []byte("Placeholder content")
		if err := storageAdapter.WriteFile(context.Background(), event.Path, content, true); err != nil {
			log.Printf("❌ Failed to process uploaded file: %v", err)
		} else {
			log.Printf("✅ Successfully processed uploaded file: %s", event.Path)
		}
	})

	kafkaClient.RegisterHandler(events.FileDeleted, func(event *events.StorageEvent) {
		log.Printf("🗑️ Processing deleted file: %s", event.Path)
		if err := storageAdapter.DeleteFile(context.Background(), event.Path); err != nil {
			log.Printf("❌ Failed to process deleted file: %v", err)
		} else {
			log.Printf("✅ Successfully processed deleted file: %s", event.Path)
		}
	})

	ctx := context.Background()
	if err := kafkaClient.StartConsumers(ctx, []string{cfg.Kafka.Topics.StorageEvents}); err != nil {
		log.Fatalf("Failed to start Kafka consumers: %v", err)
	}

	log.Println("Worker is now listening for Kafka events...")
	select {}
}
