package response

type Error struct {
	Error string `json:"error" example:"message"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error" example:"error_type"`
	Message string `json:"message" example:"error_message"`
}
