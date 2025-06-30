package persistent

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/pkg/postgres"
)

type VietQRRepo struct {
	pg *postgres.Postgres
}

func NewVietQRRepo(pg *postgres.Postgres) *VietQRRepo {
	return &VietQRRepo{pg}
}

func (r *VietQRRepo) Store(ctx context.Context, qr entity.VietQR) error {
	sql, args, err := r.pg.Builder.
		Insert("vietqr").
		Columns("id", "status", "content").
		Values(qr.ID, qr.Status, qr.Content).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	return err
}
