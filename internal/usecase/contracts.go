// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		// Translate translates text from one language to another
		Translate(ctx context.Context, t entity.Translation) (entity.Translation, error)
		// History gets translation history
		History(ctx context.Context) (entity.TranslationHistory, error)
	}

	// User -.
	User interface {
		// Create creates a new user
		Create(ctx context.Context, user entity.User) (entity.User, error)
		// GetByID gets user by ID
		GetByID(ctx context.Context, id int64) (entity.User, error)
		// Update updates user
		Update(ctx context.Context, user entity.User) error
		// Delete deletes user
		Delete(ctx context.Context, id int64) error
		// List gets all users
		List(ctx context.Context) (entity.UserHistory, error)
		// Login authenticates a user and returns a JWT token
		Login(ctx context.Context, email, password string) (string, entity.User, error)
	}

	// Redis -.
	Redis interface {
		// SetValue sets a value in Redis
		SetValue(ctx context.Context, value entity.RedisValue) error
		// GetValue gets a value from Redis
		GetValue(ctx context.Context, key string) (entity.RedisValue, error)
	}

	// Nats -.
	Nats interface {
		Publish(subject string, data []byte) error
		Subscribe(subject string, handler func(msg []byte)) (unsubscribe func() error, err error)
	}

	// VietQR -.
	VietQR interface {
		GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (*entity.VietQR, error)
		InquiryQR(ctx context.Context, id string) (*entity.VietQR, error)
		UpdateStatus(ctx context.Context, id string, status entity.VietQRStatus) error
	}

	// Billing -.
	Billing interface {
		GenerateInvoicePDF(data billing.InvoiceData, outputPath string) error
	}

	// ShipperLocation -.
	ShipperLocation interface {
		UpdateLocation(ctx context.Context, loc entity.ShipperLocation) error
		GetLocation(ctx context.Context, shipperID string) (entity.ShipperLocation, error)
	}
)
