package redis

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
)

// RedisUseCase -.
type RedisUseCase struct {
	repo repo.RedisRepo
}

// NewRedisUseCase -.
func NewRedisUseCase(r repo.RedisRepo) *RedisUseCase {
	return &RedisUseCase{
		repo: r,
	}
}

// SetValue -.
func (uc *RedisUseCase) SetValue(ctx context.Context, value entity.RedisValue) error {
	return uc.repo.SetValue(ctx, value)
}

// GetValue -.
func (uc *RedisUseCase) GetValue(ctx context.Context, key string) (entity.RedisValue, error) {
	return uc.repo.GetValue(ctx, key)
}
