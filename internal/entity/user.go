package entity

import "time"

// User represents user entity
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password is not exposed in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserHistory represents user history entity
type UserHistory struct {
	Users []User `json:"users"`
} 