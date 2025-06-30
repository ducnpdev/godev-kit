package request

// UpdateVietQRStatus represents the request body for updating a VietQR status.
type UpdateVietQRStatus struct {
	Status string `json:"status" binding:"required"`
}
