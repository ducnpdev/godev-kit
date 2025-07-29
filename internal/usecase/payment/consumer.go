package payment

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/pkg/kafka"
	"github.com/rs/zerolog"
)

// PaymentConsumer represents payment Kafka consumer
type PaymentConsumer struct {
	consumer *kafka.Consumer
	useCase  *PaymentUseCase
	logger   *zerolog.Logger
}

// NewPaymentConsumer creates new payment consumer
func NewPaymentConsumer(brokers []string, groupID string, useCase *PaymentUseCase, logger *zerolog.Logger) *PaymentConsumer {
	handler := func(ctx context.Context, key, value []byte) error {
		pc := &PaymentConsumer{
			useCase: useCase,
			logger:  logger,
		}
		return pc.handlePaymentEvent(ctx, key, value)
	}

	consumer := kafka.NewConsumer(brokers, "payment-events", groupID, handler, *logger)
	return &PaymentConsumer{
		consumer: consumer,
		useCase:  useCase,
		logger:   logger,
	}
}

// Start starts the payment consumer
func (pc *PaymentConsumer) Start(ctx context.Context) error {
	pc.logger.Info().Msg("Starting payment consumer")
	return pc.consumer.Start(ctx)
}

// handlePaymentEvent handles payment events from Kafka
func (pc *PaymentConsumer) handlePaymentEvent(ctx context.Context, key, value []byte) error {
	pc.logger.Info().
		Str("key", string(key)).
		Str("value", string(value)).
		Msg("Received payment event from Kafka")

	// Parse payment event
	var paymentEvent entity.PaymentEvent
	err := json.Unmarshal(value, &paymentEvent)
	if err != nil {
		pc.logger.Error().Err(err).Msg("Failed to unmarshal payment event")
		return fmt.Errorf("failed to unmarshal payment event: %w", err)
	}

	// Process payment
	err = pc.useCase.ProcessPayment(ctx, &paymentEvent)
	if err != nil {
		pc.logger.Error().Err(err).Msg("Failed to process payment")
		return fmt.Errorf("failed to process payment: %w", err)
	}

	pc.logger.Info().
		Int64("payment_id", paymentEvent.PaymentID).
		Str("event_type", paymentEvent.EventType).
		Msg("Payment event processed successfully")

	return nil
}

// Stop stops the payment consumer
func (pc *PaymentConsumer) Stop() error {
	pc.logger.Info().Msg("Stopping payment consumer")
	return pc.consumer.Close()
}
