package entity

// VietQRStatus represents the status of a VietQR code.
type VietQRStatus string

const (
	VietQRStatusGenerated VietQRStatus = "generated"
	VietQRStatusInProcess VietQRStatus = "in-process"
	VietQRStatusPaid      VietQRStatus = "paid"
	VietQRStatusFail      VietQRStatus = "fail"
	VietQRStatusTimeout   VietQRStatus = "timeout"
)

// VietQR represents the vietqr entity.
type VietQR struct {
	ID      string       `json:"id"`
	Status  VietQRStatus `json:"status"`
	Content string       `json:"content"`
}

// VietQRGenerateRequest represents the data needed to generate a VietQR code.
type VietQRGenerateRequest struct {
	AccountNo    string
	Amount       string
	Description  string
	MCC          string
	ReceiverName string
}
