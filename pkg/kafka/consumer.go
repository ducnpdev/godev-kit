package kafka

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

// MessageHandler -.
type MessageHandler func(ctx context.Context, key, value []byte) error

// Consumer -.
type Consumer struct {
	reader  *kafka.Reader
	logger  zerolog.Logger
	handler MessageHandler
}

// NewConsumer -.
func NewConsumer(brokers []string, topic, groupID string, handler MessageHandler, logger zerolog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Logger:   kafka.LoggerFunc(logger.Printf),
	})

	return &Consumer{
		reader:  reader,
		logger:  logger,
		handler: handler,
	}
}

// Start -.
func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info().
		Str("topic", c.reader.Config().Topic).
		Str("group_id", c.reader.Config().GroupID).
		Msg("starting kafka consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("stopping kafka consumer")
			return c.reader.Close()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error().Err(err).Msg("failed to read message")
				continue
			}

			c.logger.Debug().
				Str("topic", m.Topic).
				Int("partition", m.Partition).
				Int64("offset", m.Offset).
				Str("key", string(m.Key)).
				Msg("received message")

			if err := c.handler(ctx, m.Key, m.Value); err != nil {
				c.logger.Error().Err(err).Msg("failed to handle message")
			}
		}
	}
}

// Close -.
func (c *Consumer) Close() error {
	return c.reader.Close()
}

// ConsumeMessages -.
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func(key, value []byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to read message: %w", err)
			}

			if err := handler(m.Key, m.Value); err != nil {
				c.logger.Error().Err(err).Msg("failed to process message")
			}
		}
	}
}
