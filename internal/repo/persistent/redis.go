package persistent

import (
	"context"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/pkg/redis"
)

const _defaultTimeout = 5 * time.Second

// RedisRepo -.
type RedisRepo struct {
	r *redis.Redis
}

// NewRedisRepo -.
func NewRedisRepo(r *redis.Redis) *RedisRepo {
	return &RedisRepo{r: r}
}

// SetValue -.
func (r *RedisRepo) SetValue(ctx context.Context, value entity.RedisValue) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	return r.r.Client().Set(ctx, value.Key, value.Value, 0).Err()
}

// GetValue -.
func (r *RedisRepo) GetValue(ctx context.Context, key string) (entity.RedisValue, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()

	val, err := r.r.Client().Get(ctx, key).Result()
	if err != nil {
		return entity.RedisValue{}, err
	}

	return entity.RedisValue{Key: key, Value: val}, nil
}
