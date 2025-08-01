package models

import (
	"time"
)

// Payment represents payment database model
type Payment struct {
	ID            int64     `db:"id" json:"id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	Amount        float64   `db:"amount" json:"amount"`
	Currency      string    `db:"currency" json:"currency"`
	PaymentType   string    `db:"payment_type" json:"payment_type"`
	Status        string    `db:"status" json:"status"`
	MeterNumber   string    `db:"meter_number" json:"meter_number"`
	CustomerCode  string    `db:"customer_code" json:"customer_code"`
	Description   string    `db:"description" json:"description"`
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

// PaymentHistory represents payment history database model
type PaymentHistory struct {
	ID            int64     `db:"id" json:"id"`
	PaymentID     int64     `db:"payment_id" json:"payment_id"`
	UserID        int64     `db:"user_id" json:"user_id"`
	Status        string    `db:"status" json:"status"`
	Amount        float64   `db:"amount" json:"amount"`
	Currency      string    `db:"currency" json:"currency"`
	PaymentType   string    `db:"payment_type" json:"payment_type"`
	MeterNumber   string    `db:"meter_number" json:"meter_number"`
	CustomerCode  string    `db:"customer_code" json:"customer_code"`
	Description   string    `db:"description" json:"description"`
	TransactionID string    `db:"transaction_id" json:"transaction_id"`
	PaymentMethod string    `db:"payment_method" json:"payment_method"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}
