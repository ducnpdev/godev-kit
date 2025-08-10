package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent"
	"github.com/ducnpdev/godev-kit/pkg/kafka"
	"github.com/rs/zerolog"
)

// PaymentUseCase represents payment use case
type PaymentUseCase struct {
	paymentRepo *persistent.PaymentRepo
	kafkaProd   *kafka.Producer
	logger      *zerolog.Logger
}

// NewPaymentUseCase creates new payment use case
func NewPaymentUseCase(paymentRepo *persistent.PaymentRepo, kafkaProd *kafka.Producer, logger *zerolog.Logger) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		kafkaProd:   kafkaProd,
		logger:      logger,
	}
}

// RegisterPayment registers a new payment and sends to Kafka
func (uc *PaymentUseCase) RegisterPayment(ctx context.Context, req *entity.PaymentRequest) (*entity.PaymentResponse, error) {
	// Create payment entity
	payment := &entity.Payment{
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentType:   req.PaymentType,
		Status:        entity.PaymentStatusPending,
		MeterNumber:   req.MeterNumber,
		CustomerCode:  req.CustomerCode,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
	}

	// Save to database
	err := uc.paymentRepo.Create(ctx, payment)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to create payment in database")
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// Create payment event for Kafka
	paymentEvent := &entity.PaymentEvent{
		ID:            payment.ID,
		EventType:     entity.PaymentCreatedEvent,
		UserID:        payment.UserID,
		PaymentID:     payment.ID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		MeterNumber:   payment.MeterNumber,
		CustomerCode:  payment.CustomerCode,
		Description:   payment.Description,
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		Timestamp:     time.Now(),
	}

	// Send to Kafka if producer is available
	if uc.kafkaProd != nil {
		err = uc.kafkaProd.SendMessage(ctx, "payment-events", []byte(paymentEvent.TransactionID), paymentEvent)
		if err != nil {
			uc.logger.Error().Err(err).Msg("Failed to send payment event to Kafka")
			// Note: In production, you might want to handle this differently
			// For now, we'll just log the error but still return success
		}
	} else {
		uc.logger.Info().Msg("Kafka producer is disabled, skipping payment event publishing")
	}

	uc.logger.Info().
		Int64("payment_id", payment.ID).
		Int64("user_id", payment.UserID).
		Float64("amount", payment.Amount).
		Str("status", string(payment.Status)).
		Msg("Payment registered successfully")

	return &entity.PaymentResponse{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		MeterNumber:   payment.MeterNumber,
		CustomerCode:  payment.CustomerCode,
		Description:   payment.Description,
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
	}, nil
}

// ProcessPayment processes payment from Kafka message
func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, paymentEvent *entity.PaymentEvent) error {
	uc.logger.Info().
		Int64("payment_id", paymentEvent.PaymentID).
		Str("event_type", paymentEvent.EventType).
		Msg("Processing payment from Kafka")

	// Update status to processing
	err := uc.paymentRepo.UpdateStatus(ctx, paymentEvent.PaymentID, entity.PaymentStatusProcessing)
	if err != nil {
		uc.logger.Error().Err(err).Int64("payment_id", paymentEvent.PaymentID).Msg("Failed to update payment status to processing")
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	// Simulate payment processing
	// In real implementation, this would call external payment gateway
	success := uc.simulatePaymentProcessing(paymentEvent)

	var newStatus entity.PaymentStatus
	if success {
		newStatus = entity.PaymentStatusCompleted
		uc.logger.Info().Int64("payment_id", paymentEvent.PaymentID).Msg("Payment processed successfully")
	} else {
		newStatus = entity.PaymentStatusFailed
		uc.logger.Error().Int64("payment_id", paymentEvent.PaymentID).Msg("Payment processing failed")
	}

	// Update final status
	err = uc.paymentRepo.UpdateStatus(ctx, paymentEvent.PaymentID, newStatus)
	if err != nil {
		uc.logger.Error().Err(err).Int64("payment_id", paymentEvent.PaymentID).Msg("Failed to update payment final status")
		return fmt.Errorf("failed to update payment final status: %w", err)
	}

	// Get updated payment for history
	payment, err := uc.paymentRepo.GetByID(ctx, paymentEvent.PaymentID)
	if err != nil {
		uc.logger.Error().Err(err).Int64("payment_id", paymentEvent.PaymentID).Msg("Failed to get payment for history")
		return fmt.Errorf("failed to get payment for history: %w", err)
	}

	// Create history record
	err = uc.paymentRepo.CreateHistory(ctx, payment)
	if err != nil {
		uc.logger.Error().Err(err).Int64("payment_id", paymentEvent.PaymentID).Msg("Failed to create payment history")
		return fmt.Errorf("failed to create payment history: %w", err)
	}

	uc.logger.Info().
		Int64("payment_id", paymentEvent.PaymentID).
		Str("status", string(newStatus)).
		Msg("Payment processing completed")

	return nil
}

// GetPaymentByID gets payment by ID
func (uc *PaymentUseCase) GetPaymentByID(ctx context.Context, id int64) (*entity.PaymentResponse, error) {
	payment, err := uc.paymentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}

	return &entity.PaymentResponse{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentType:   payment.PaymentType,
		Status:        payment.Status,
		MeterNumber:   payment.MeterNumber,
		CustomerCode:  payment.CustomerCode,
		Description:   payment.Description,
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
	}, nil
}

// GetPaymentsByUserID gets payments by user ID
func (uc *PaymentUseCase) GetPaymentsByUserID(ctx context.Context, userID int64) ([]*entity.PaymentResponse, error) {
	payments, err := uc.paymentRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	responses := make([]*entity.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = &entity.PaymentResponse{
			ID:            payment.ID,
			UserID:        payment.UserID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentType:   payment.PaymentType,
			Status:        payment.Status,
			MeterNumber:   payment.MeterNumber,
			CustomerCode:  payment.CustomerCode,
			Description:   payment.Description,
			TransactionID: payment.TransactionID,
			PaymentMethod: payment.PaymentMethod,
			CreatedAt:     payment.CreatedAt,
		}
	}

	return responses, nil
}

// simulatePaymentProcessing simulates payment processing
// In real implementation, this would call external payment gateway
func (uc *PaymentUseCase) simulatePaymentProcessing(paymentEvent *entity.PaymentEvent) bool {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	// Simulate 90% success rate
	// In real implementation, this would be based on actual payment gateway response
	return true // For demo purposes, always return success
}
