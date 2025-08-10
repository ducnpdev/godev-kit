package persistent

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/ducnpdev/godev-kit/pkg/kafka"
	"github.com/rs/zerolog"
)

// KafkaRepo -.
type KafkaRepo struct {
	manager *kafka.Manager
}

// NewKafkaRepo -.
func NewKafkaRepo(brokers []string, logger zerolog.Logger) repo.KafkaRepo {
	return &KafkaRepo{
		manager: kafka.NewManager(brokers, logger),
	}
}

// NewKafkaRepoWithConfig creates a new Kafka repository with configuration
func NewKafkaRepoWithConfig(brokers []string, logger zerolog.Logger, producerEnabled, consumerEnabled bool) repo.KafkaRepo {
	return &KafkaRepo{
		manager: kafka.NewManagerWithConfig(brokers, logger, producerEnabled, consumerEnabled),
	}
}

// SendMessage -.
func (k *KafkaRepo) SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error {
	return k.manager.SendMessage(ctx, topic, key, value)
}

// AddConsumer -.
func (k *KafkaRepo) AddConsumer(topic, groupID string, handler func(ctx context.Context, key, value []byte) error) error {
	return k.manager.AddConsumer(topic, groupID, handler)
}

// StartConsumer -.
func (k *KafkaRepo) StartConsumer(ctx context.Context, topic string) error {
	return k.manager.StartConsumer(ctx, topic)
}

// StartAllConsumers -.
func (k *KafkaRepo) StartAllConsumers(ctx context.Context) {
	k.manager.StartAllConsumers(ctx)
}

// Close -.
func (k *KafkaRepo) Close() error {
	return k.manager.Close()
}

// Control methods for Kafka producer and consumer

// EnableProducer enables the Kafka producer
func (k *KafkaRepo) EnableProducer() {
	k.manager.EnableProducer()
}

// DisableProducer disables the Kafka producer
func (k *KafkaRepo) DisableProducer() {
	k.manager.DisableProducer()
}

// IsProducerEnabled returns the current producer status
func (k *KafkaRepo) IsProducerEnabled() bool {
	return k.manager.IsProducerEnabled()
}

// EnableConsumer enables the Kafka consumer
func (k *KafkaRepo) EnableConsumer() {
	k.manager.EnableConsumer()
}

// DisableConsumer disables the Kafka consumer
func (k *KafkaRepo) DisableConsumer() {
	k.manager.DisableConsumer()
}

// IsConsumerEnabled returns the current consumer status
func (k *KafkaRepo) IsConsumerEnabled() bool {
	return k.manager.IsConsumerEnabled()
}

// GetStatus returns the current status of both producer and consumer
func (k *KafkaRepo) GetStatus() map[string]interface{} {
	return k.manager.GetStatus()
}
