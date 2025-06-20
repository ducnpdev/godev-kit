// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, models.TranslationModel) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}

	// UserRepo -.
	UserRepo interface {
		Create(context.Context, models.UserModel) (entity.User, error)
		GetByID(context.Context, int64) (entity.User, error)
		GetByEmail(context.Context, string) (entity.User, error)
		Update(context.Context, models.UserModel) error
		Delete(context.Context, int64) error
		List(context.Context) ([]entity.User, error)
		// Database access methods
		GetBuilder() squirrel.StatementBuilderType
		GetPool() *pgxpool.Pool
	}

	// KafkaRepo -.
	KafkaRepo interface {
		SendMessage(ctx context.Context, topic string, key []byte, value interface{}) error
		AddConsumer(topic, groupID string, handler func(ctx context.Context, key, value []byte) error) error
		StartConsumer(ctx context.Context, topic string) error
		StartAllConsumers(ctx context.Context)
		Close() error
	}
)
