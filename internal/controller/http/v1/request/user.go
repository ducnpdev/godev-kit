package request

// CreateUser represents create user request
type CreateUser struct {
	Email    string `json:"email"     validate:"required,email" example:"user@example.com"`
	Username string `json:"username"  validate:"required,min=3"  example:"johndoe"`
	Password string `json:"password"  validate:"required,min=6"  example:"password123"`
}

// UpdateUser represents update user request
type UpdateUser struct {
	Email    string `json:"email"     validate:"omitempty,email" example:"user@example.com"`
	Username string `json:"username"  validate:"omitempty,min=3"  example:"johndoe"`
	Password string `json:"password"  validate:"omitempty,min=6"  example:"password123"`
}

// LoginUser represents login user request
type LoginUser struct {
	Email    string `json:"email"     validate:"required,email" example:"user@example.com"`
	Password string `json:"password"  validate:"required,min=6"  example:"password123"`
}
