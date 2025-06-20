// Package kafka implements Kafka producer and consumer functionality.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// Producer -.
type Producer struct {
	writer *kafka.Writer
	logger zerolog.Logger
}

// NewProducer -.
func NewProducer(brokers []string, logger zerolog.Logger) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		Logger:       kafka.LoggerFunc(logger.Printf),
	}

	return &Producer{
		writer: writer,
		logger: logger,
	}
}

// SendMessage -.
func (p *Producer) SendMessage(ctx context.Context, topic string,
	key []byte, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Key:   key,
		Value: valueBytes,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	p.logger.Info().
		Str("topic", topic).
		Str("key", string(key)).
		Msg("message sent successfully")

	return nil
}

// Close -.
func (p *Producer) Close() error {
	return p.writer.Close()
}
