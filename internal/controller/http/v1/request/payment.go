package request

// PaymentRequest represents payment request
type PaymentRequest struct {
	UserID        int64   `json:"user_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Currency      string  `json:"currency" binding:"required"`
	PaymentType   string  `json:"payment_type" binding:"required"`
	MeterNumber   string  `json:"meter_number" binding:"required"`
	CustomerCode  string  `json:"customer_code" binding:"required"`
	Description   string  `json:"description"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
}
