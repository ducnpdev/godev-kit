package persistent

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/ducnpdev/godev-kit/pkg/postgres"
)

type ShipperLocationRepo struct {
	pg *postgres.Postgres
}

func NewShipperLocationRepo(pg *postgres.Postgres) *ShipperLocationRepo {
	return &ShipperLocationRepo{pg: pg}
}

// Store inserts a new shipper location record into the DB
func (r *ShipperLocationRepo) Store(ctx context.Context, loc entity.ShipperLocation) error {
	model := models.ShipperLocationModel{
		ShipperID: loc.ShipperID,
		Latitude:  loc.Latitude,
		Longitude: loc.Longitude,
		Timestamp: loc.Timestamp,
	}
	sql, args, err := r.pg.Builder.
		Insert("shipper_locations").
		Columns("shipper_id", "latitude", "longitude", "timestamp").
		Values(model.ShipperID, model.Latitude, model.Longitude, model.Timestamp).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	return err
}
