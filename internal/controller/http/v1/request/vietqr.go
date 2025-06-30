package request

// GenerateQR represents the request body for generating a VietQR code.
type GenerateQR struct {
	AccountNo    string `json:"accountNo" binding:"required"`
	Amount       string `json:"amount" binding:"required"`
	Description  string `json:"description"`
	MCC          string `json:"mcc"`
	ReceiverName string `json:"receiverName"`
}

// UpdateVietQRStatus represents the request body for updating a VietQR status.
type UpdateVietQRStatus struct {
	Status string `json:"status" binding:"required"`
}
