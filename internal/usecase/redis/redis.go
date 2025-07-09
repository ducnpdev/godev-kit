package redis

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent"
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

// ShipperLocationUseCase implements updating shipper location in Redis and DB
type ShipperLocationUseCase struct {
	redisRepo    *persistent.RedisRepo
	locationRepo *persistent.ShipperLocationRepo
}

func NewShipperLocationUseCase(redisRepo *persistent.RedisRepo, locationRepo *persistent.ShipperLocationRepo) *ShipperLocationUseCase {
	return &ShipperLocationUseCase{
		redisRepo:    redisRepo,
		locationRepo: locationRepo,
	}
}

// UpdateLocation updates the latest location in Redis and appends to DB
func (uc *ShipperLocationUseCase) UpdateLocation(ctx context.Context, loc entity.ShipperLocation) error {
	if err := uc.redisRepo.SetShipperLocation(ctx, loc); err != nil {
		return err
	}
	return uc.locationRepo.Store(ctx, loc)
}

// GetLocation retrieves the latest location of a shipper from Redis (cache-aside pattern)
func (uc *ShipperLocationUseCase) GetLocation(ctx context.Context, shipperID string) (entity.ShipperLocation, error) {
	// Try Redis first
	loc, err := uc.redisRepo.GetShipperLocation(ctx, shipperID)
	if err == nil {
		return loc, nil
	}
	// If not found in Redis, get from DB
	dbLoc, dbErr := uc.locationRepo.GetLatestByShipperID(ctx, shipperID)
	if dbErr != nil {
		return entity.ShipperLocation{}, dbErr
	}
	// Update Redis for next time
	_ = uc.redisRepo.SetShipperLocation(ctx, dbLoc)
	return dbLoc, nil
}
