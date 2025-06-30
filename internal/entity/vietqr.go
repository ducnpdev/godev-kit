package entity

// VietQR represents the vietqr entity.
type VietQR struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}

// VietQRGenerateRequest represents the data needed to generate a VietQR code.
type VietQRGenerateRequest struct {
	AccountNo    string
	Amount       string
	Description  string
	MCC          string
	ReceiverName string
}
