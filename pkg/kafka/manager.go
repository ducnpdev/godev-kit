package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"
)

// Manager -.
type Manager struct {
	producer  *Producer
	consumers map[string]*Consumer
	logger    zerolog.Logger
	mu        sync.RWMutex
	brokers   []string
}

// NewManager -.
func NewManager(brokers []string, logger zerolog.Logger) *Manager {
	return &Manager{
		producer:  NewProducer(brokers, logger),
		consumers: make(map[string]*Consumer),
		logger:    logger,
		brokers:   brokers,
	}
}

// SendMessage -.
func (m *Manager) SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error {
	return m.producer.SendMessage(ctx, topic, key, value)
}

// AddConsumer -.
func (m *Manager) AddConsumer(topic, groupID string, handler MessageHandler) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.consumers[topic]; exists {
		return fmt.Errorf("consumer for topic %s already exists", topic)
	}

	consumer := NewConsumer(
		m.brokers,
		topic,
		groupID,
		handler,
		m.logger,
	)

	m.consumers[topic] = consumer
	return nil
}

// StartConsumer -.
func (m *Manager) StartConsumer(ctx context.Context, topic string) error {
	m.mu.RLock()
	consumer, exists := m.consumers[topic]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("consumer for topic %s not found", topic)
	}

	return consumer.Start(ctx)
}

// StartAllConsumers -.
func (m *Manager) StartAllConsumers(ctx context.Context) {
	for topic, consumer := range m.consumers {
		go func(t string, c *Consumer) {
			if err := c.Start(ctx); err != nil {
				m.logger.Error().Err(err).Str("topic", t).Msg("consumer stopped with error")
			}
		}(topic, consumer)
	}
}

// Close -.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Close all consumers
	for topic, consumer := range m.consumers {
		if err := consumer.Close(); err != nil {
			m.logger.Error().Err(err).Str("topic", topic).Msg("failed to close consumer")
		}
	}

	// Close producer
	return m.producer.Close()
}
