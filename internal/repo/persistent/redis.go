package persistent

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
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

// SetShipperLocation stores the latest shipper location in Redis.
func (r *RedisRepo) SetShipperLocation(ctx context.Context, loc entity.ShipperLocation) error {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()
	data, err := json.Marshal(loc)
	if err != nil {
		return err
	}
	key := "shipper:location:" + loc.ShipperID
	return r.r.Client().Set(ctx, key, data, 0).Err()
}

// GetShipperLocation retrieves the latest shipper location from Redis.
func (r *RedisRepo) GetShipperLocation(ctx context.Context, shipperID string) (entity.ShipperLocation, error) {
	ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
	defer cancel()
	key := "shipper:location:" + shipperID
	val, err := r.r.Client().Get(ctx, key).Result()
	if err != nil {
		return entity.ShipperLocation{}, err
	}
	var loc entity.ShipperLocation
	if err := json.Unmarshal([]byte(val), &loc); err != nil {
		return entity.ShipperLocation{}, err
	}
	return loc, nil
}

// Convert entity.ShipperLocation to models.ShipperLocationModel
func ToShipperLocationModel(loc entity.ShipperLocation) models.ShipperLocationModel {
	return models.ShipperLocationModel{
		ShipperID: loc.ShipperID,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timestamp: loc.Timestamp,
	}
}

// Convert models.ShipperLocationModel to entity.ShipperLocation
func ToShipperLocationEntity(model models.ShipperLocationModel) entity.ShipperLocation {
	return entity.ShipperLocation{
		ShipperID: model.ShipperID,
		Latitude:  model.Latitude,
		Longitude: model.Longitude,
		Timestamp: model.Timestamp,
	}
}
