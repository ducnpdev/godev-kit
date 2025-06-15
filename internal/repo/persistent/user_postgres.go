package persistent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/ducnpdev/godev-kit/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// const _defaultEntityCap = 64

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// Create -.
func (r *UserRepo) Create(ctx context.Context, user models.UserModel) (entity.User, error) {
	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	sql, args, err := r.Builder.
		Insert("users").
		Columns("email, username, password, created_at, updated_at").
		Values(user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Create - r.Builder: %w", err)
	}

	var id int64
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - Create - r.Pool.QueryRow: %w", err)
	}

	return entity.User{
		ID:        id,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, id int64) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, created_at, updated_at").
		From("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByID - r.Builder: %w", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}

// Update -.
func (r *UserRepo) Update(ctx context.Context, user models.UserModel) error {
	// Update timestamp
	user.UpdatedAt = time.Now()

	builder := r.Builder.Update("users").
		Set("updated_at", user.UpdatedAt)

	if user.Email != "" {
		builder = builder.Set("email", user.Email)
	}
	if user.Username != "" {
		builder = builder.Set("username", user.Username)
	}
	if user.Password != "" {
		builder = builder.Set("password", user.Password)
	}

	sql, args, err := builder.
		Where("id = ?", user.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Update - r.Pool.Exec: %w", err)
	}

	return nil
}

// Delete -.
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	sql, args, err := r.Builder.
		Delete("users").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Delete - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo - Delete - r.Pool.Exec: %w", err)
	}

	return nil
}

// List -.
func (r *UserRepo) List(ctx context.Context) ([]entity.User, error) {
	sql, _, err := r.Builder.
		Select("id, email, username, created_at, updated_at").
		From("users").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - List - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("UserRepo - List - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	users := make([]entity.User, 0, _defaultEntityCap)

	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("UserRepo - List - rows.Scan: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByEmail -.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, email, username, password, created_at, updated_at").
		From("users").
		Where("email = ?", email).
		ToSql()
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Builder: %w", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, errors.New("user not found")
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetByEmail - r.Pool.QueryRow: %w", err)
	}

	return user, nil
}

// GetBuilder returns the statement builder
func (r *UserRepo) GetBuilder() squirrel.StatementBuilderType {
	return r.Builder
}

// GetPool returns the database pool
func (r *UserRepo) GetPool() *pgxpool.Pool {
	return r.Pool
}
