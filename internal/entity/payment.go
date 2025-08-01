package entity

import (
	"time"
)

// PaymentStatus represents payment status
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

// PaymentType represents payment type
type PaymentType string

const (
	PaymentTypeElectric PaymentType = "electric"
	PaymentTypeWater    PaymentType = "water"
	PaymentTypeGas      PaymentType = "gas"
)

// Payment represents payment entity
type Payment struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PaymentType   PaymentType   `json:"payment_type"`
	Status        PaymentStatus `json:"status"`
	MeterNumber   string        `json:"meter_number"`
	CustomerCode  string        `json:"customer_code"`
	Description   string        `json:"description"`
	TransactionID string        `json:"transaction_id"`
	PaymentMethod string        `json:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// PaymentEvent represents payment event for Kafka
type PaymentEvent struct {
	ID            int64         `json:"id"`
	EventType     string        `json:"event_type"`
	UserID        int64         `json:"user_id"`
	PaymentID     int64         `json:"payment_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PaymentType   PaymentType   `json:"payment_type"`
	Status        PaymentStatus `json:"status"`
	MeterNumber   string        `json:"meter_number"`
	CustomerCode  string        `json:"customer_code"`
	Description   string        `json:"description"`
	TransactionID string        `json:"transaction_id"`
	PaymentMethod string        `json:"payment_method"`
	Timestamp     time.Time     `json:"timestamp"`
}

// PaymentRequest represents payment request from API
type PaymentRequest struct {
	UserID        int64       `json:"user_id" binding:"required"`
	Amount        float64     `json:"amount" binding:"required"`
	Currency      string      `json:"currency" binding:"required"`
	PaymentType   PaymentType `json:"payment_type" binding:"required"`
	MeterNumber   string      `json:"meter_number" binding:"required"`
	CustomerCode  string      `json:"customer_code" binding:"required"`
	Description   string      `json:"description"`
	PaymentMethod string      `json:"payment_method" binding:"required"`
}

// PaymentResponse represents payment response
type PaymentResponse struct {
	ID            int64         `json:"id"`
	UserID        int64         `json:"user_id"`
	Amount        float64       `json:"amount"`
	Currency      string        `json:"currency"`
	PaymentType   PaymentType   `json:"payment_type"`
	Status        PaymentStatus `json:"status"`
	MeterNumber   string        `json:"meter_number"`
	CustomerCode  string        `json:"customer_code"`
	Description   string        `json:"description"`
	TransactionID string        `json:"transaction_id"`
	PaymentMethod string        `json:"payment_method"`
	CreatedAt     time.Time     `json:"created_at"`
}

// Event types for payment
const (
	PaymentCreatedEvent   = "payment.created"
	PaymentProcessedEvent = "payment.processed"
	PaymentCompletedEvent = "payment.completed"
	PaymentFailedEvent    = "payment.failed"
)
