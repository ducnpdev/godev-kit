package request

// PaymentRequest represents payment request
// @Description Payment request for electric bill
type PaymentRequest struct {
	UserID        int64   `json:"user_id" binding:"required" example:"1"`
	Amount        float64 `json:"amount" binding:"required" example:"500000"`
	Currency      string  `json:"currency" binding:"required" example:"VND"`
	PaymentType   string  `json:"payment_type" binding:"required" example:"electric"`
	MeterNumber   string  `json:"meter_number" binding:"required" example:"EVN001234567"`
	CustomerCode  string  `json:"customer_code" binding:"required" example:"CUST001"`
	Description   string  `json:"description" example:"Thanh toán tiền điện tháng 12/2024"`
	PaymentMethod string  `json:"payment_method" binding:"required" example:"bank_transfer"`
}
