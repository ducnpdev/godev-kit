package entity

import (
	"time"
)

// UserEvent -.
type UserEvent struct {
	ID        int64     `json:"id"`
	EventType string    `json:"event_type"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

// TranslationEvent -.
type TranslationEvent struct {
	ID         int64     `json:"id"`
	EventType  string    `json:"event_type"`
	UserID     int64     `json:"user_id"`
	Source     string    `json:"source"`
	Target     string    `json:"target"`
	Original   string    `json:"original"`
	Translated string    `json:"translated"`
	Timestamp  time.Time `json:"timestamp"`
}

// Event types
const (
	UserCreatedEvent          = "user.created"
	UserUpdatedEvent          = "user.updated"
	UserDeletedEvent          = "user.deleted"
	TranslationRequestEvent   = "translation.requested"
	TranslationCompletedEvent = "translation.completed"
)
