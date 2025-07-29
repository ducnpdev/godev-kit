package response

import "time"

// PaymentResponse represents payment response
// @Description Payment response with all payment details
type PaymentResponse struct {
	ID            int64     `json:"id" example:"1"`
	UserID        int64     `json:"user_id" example:"1"`
	Amount        float64   `json:"amount" example:"500000"`
	Currency      string    `json:"currency" example:"VND"`
	PaymentType   string    `json:"payment_type" example:"electric"`
	Status        string    `json:"status" example:"pending"`
	MeterNumber   string    `json:"meter_number" example:"EVN001234567"`
	CustomerCode  string    `json:"customer_code" example:"CUST001"`
	Description   string    `json:"description" example:"Thanh toán tiền điện tháng 12/2024"`
	TransactionID string    `json:"transaction_id" example:"uuid-here"`
	PaymentMethod string    `json:"payment_method" example:"bank_transfer"`
	CreatedAt     time.Time `json:"created_at" example:"2024-12-20T10:30:00Z"`
}
