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

	// Control flags
	producerEnabled bool
	consumerEnabled bool
	controlMu       sync.RWMutex
}

// NewManager -.
func NewManager(brokers []string, logger zerolog.Logger) *Manager {
	return &Manager{
		producer:        NewProducer(brokers, logger),
		consumers:       make(map[string]*Consumer),
		logger:          logger,
		brokers:         brokers,
		producerEnabled: true, // Producer enabled by default
		consumerEnabled: true, // Consumer enabled by default
	}
}

// NewManagerWithConfig creates a new manager with configuration
func NewManagerWithConfig(brokers []string, logger zerolog.Logger, producerEnabled, consumerEnabled bool) *Manager {
	return &Manager{
		producer:        NewProducer(brokers, logger),
		consumers:       make(map[string]*Consumer),
		logger:          logger,
		brokers:         brokers,
		producerEnabled: producerEnabled,
		consumerEnabled: consumerEnabled,
	}
}

// SendMessage -.
func (m *Manager) SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error {
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()

	if !m.producerEnabled {
		return fmt.Errorf("kafka producer is disabled")
	}

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
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()

	if !m.consumerEnabled {
		return fmt.Errorf("kafka consumer is disabled")
	}

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
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()

	if !m.consumerEnabled {
		m.logger.Warn().Msg("kafka consumer is disabled, skipping start all consumers")
		return
	}

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

// EnableProducer enables the Kafka producer
func (m *Manager) EnableProducer() {
	m.controlMu.Lock()
	defer m.controlMu.Unlock()
	m.producerEnabled = true
	m.logger.Info().Msg("kafka producer enabled")
}

// DisableProducer disables the Kafka producer
func (m *Manager) DisableProducer() {
	m.controlMu.Lock()
	defer m.controlMu.Unlock()
	m.producerEnabled = false
	m.logger.Info().Msg("kafka producer disabled")
}

// IsProducerEnabled returns the current producer status
func (m *Manager) IsProducerEnabled() bool {
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()
	return m.producerEnabled
}

// EnableConsumer enables the Kafka consumer
func (m *Manager) EnableConsumer() {
	m.controlMu.Lock()
	defer m.controlMu.Unlock()
	m.consumerEnabled = true
	m.logger.Info().Msg("kafka consumer enabled")
}

// DisableConsumer disables the Kafka consumer
func (m *Manager) DisableConsumer() {
	m.controlMu.Lock()
	defer m.controlMu.Unlock()
	m.consumerEnabled = false
	m.logger.Info().Msg("kafka consumer disabled")
}

// IsConsumerEnabled returns the current consumer status
func (m *Manager) IsConsumerEnabled() bool {
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()
	return m.consumerEnabled
}

// GetStatus returns the current status of both producer and consumer
func (m *Manager) GetStatus() map[string]interface{} {
	m.controlMu.RLock()
	defer m.controlMu.RUnlock()

	m.mu.RLock()
	consumerCount := len(m.consumers)
	m.mu.RUnlock()

	return map[string]interface{}{
		"producer_enabled": m.producerEnabled,
		"consumer_enabled": m.consumerEnabled,
		"consumer_count":   consumerCount,
		"brokers":          m.brokers,
	}
}
