package models

import "time"

// UserModel represents user database model
type UserModel struct {
	ID        int64     `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"column:username;uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"column:password;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
} 