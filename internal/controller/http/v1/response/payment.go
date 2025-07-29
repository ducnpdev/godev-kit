package response

import "time"

// PaymentResponse represents payment response
type PaymentResponse struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentType   string    `json:"payment_type"`
	Status        string    `json:"status"`
	MeterNumber   string    `json:"meter_number"`
	CustomerCode  string    `json:"customer_code"`
	Description   string    `json:"description"`
	TransactionID string    `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
}
