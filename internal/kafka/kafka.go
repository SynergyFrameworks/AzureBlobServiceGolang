package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"

	"project-root/internal/events"
)

// Kafka producer and consumer.
type KafkaClient struct {
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
	handlers      map[events.EventType]func(event *events.StorageEvent)
	handlersMutex sync.RWMutex
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// NewKafkaClient
func NewKafkaClient(brokers []string, groupID string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	// Create producer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	// Create
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("failed to create Kafka consumer group: %v", err)
	}

	return &KafkaClient{
		producer:      producer,
		consumerGroup: consumerGroup,
		handlers:      make(map[events.EventType]func(event *events.StorageEvent)),
	}, nil
}

// Publish
func (k *KafkaClient) Publish(topic string, event *events.StorageEvent) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		log.Printf("❌ Failed to publish message to Kafka: %v", err)
		return err
	}

	log.Printf("✅ Message published to Kafka [Topic: %s, Partition: %d, Offset: %d]", topic, partition, offset)
	return nil
}

// RegisterHandler
func (k *KafkaClient) RegisterHandler(eventType events.EventType, handler func(event *events.StorageEvent)) {
	k.handlersMutex.Lock()
	defer k.handlersMutex.Unlock()
	k.handlers[eventType] = handler
	log.Printf("✅ Registered handler for event type: %s", eventType)
}

// StartConsumers begins consuming messages from Kafka topics.
func (k *KafkaClient) StartConsumers(ctx context.Context, topics []string) error {
	ctx, cancel := context.WithCancel(ctx)
	k.cancel = cancel

	// Run the consumer loop
	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			if err := k.consumerGroup.Consume(ctx, topics, k); err != nil {
				log.Printf("❌ Error during Kafka consumption: %v", err)
			}

			// Check if context is done
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

// Close shuts down the Kafka producer, consumer, and any active goroutines.
func (k *KafkaClient) Close() {
	if k.cancel != nil {
		k.cancel()
	}
	k.wg.Wait()
	k.producer.Close()
	k.consumerGroup.Close()
	log.Println("✅ Kafka client closed")
}

// ConsumeClaim processes Kafka messages from a specific topic/partition.
func (k *KafkaClient) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event events.StorageEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("❌ Failed to parse Kafka message: %v", err)
			continue
		}

		log.Printf("📩 Received event from Kafka: %+v", event)

		// Find and execute the handler for the event type
		k.handlersMutex.RLock()
		handler, exists := k.handlers[event.Type]
		k.handlersMutex.RUnlock()
		if exists {
			handler(&event)
		} else {
			log.Printf("⚠️ No handler registered for event type: %s", event.Type)
		}

		// Mark the message as processed
		sess.MarkMessage(message, "")
	}

	return nil
}

// Setup is called once when the consumer group session begins.
func (k *KafkaClient) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called once when the consumer group session ends.
func (k *KafkaClient) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
