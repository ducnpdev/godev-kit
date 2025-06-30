package persistent

import (
	"context"

	"github.com/Masterminds/squirrel"
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

func (r *VietQRRepo) FindByID(ctx context.Context, id string) (entity.VietQR, error) {
	sql, args, err := r.pg.Builder.
		Select("id", "status", "content").
		From("vietqr").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return entity.VietQR{}, err
	}

	var qr entity.VietQR
	err = r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&qr.ID, &qr.Status, &qr.Content)
	if err != nil {
		return entity.VietQR{}, err
	}

	return qr, nil
}

func (r *VietQRRepo) UpdateStatus(ctx context.Context, id string, status entity.VietQRStatus) error {
	sql, args, err := r.pg.Builder.
		Update("vietqr").
		Set("status", status).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Pool.Exec(ctx, sql, args...)
	return err
}
