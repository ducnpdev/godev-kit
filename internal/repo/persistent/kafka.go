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
