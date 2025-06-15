package response

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       int64  `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"user"`
}
