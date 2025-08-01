package persistent

import (
	"context"
	"fmt"
	"time"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/persistent/models"
	"github.com/ducnpdev/godev-kit/pkg/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PaymentRepo represents payment repository
type PaymentRepo struct {
	*postgres.Postgres
}

// NewPaymentRepo creates new payment repository
func NewPaymentRepo(pg *postgres.Postgres) *PaymentRepo {
	return &PaymentRepo{pg}
}

// Create creates new payment
func (r *PaymentRepo) Create(ctx context.Context, payment *entity.Payment) error {
	// Set timestamps and transaction ID
	now := time.Now()
	transactionID := uuid.New().String()
	payment.CreatedAt = now
	payment.UpdatedAt = now
	payment.TransactionID = transactionID

	sql, args, err := r.Builder.
		Insert("payments").
		Columns("user_id, amount, currency, payment_type, status, meter_number, customer_code, description, transaction_id, payment_method, created_at, updated_at").
		Values(payment.UserID, payment.Amount, payment.Currency, payment.PaymentType, payment.Status, payment.MeterNumber, payment.CustomerCode, payment.Description, transactionID, payment.PaymentMethod, now, now).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return fmt.Errorf("PaymentRepo - Create - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&payment.ID)
	if err != nil {
		return fmt.Errorf("PaymentRepo - Create - r.Pool.QueryRow: %w", err)
	}

	return nil
}

// GetByID gets payment by ID
func (r *PaymentRepo) GetByID(ctx context.Context, id int64) (*entity.Payment, error) {
	sql, args, err := r.Builder.
		Select("id, user_id, amount, currency, payment_type, status, meter_number, customer_code, description, transaction_id, payment_method, created_at, updated_at").
		From("payments").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PaymentRepo - GetByID - r.Builder: %w", err)
	}

	var payment models.Payment
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentType,
		&payment.Status,
		&payment.MeterNumber,
		&payment.CustomerCode,
		&payment.Description,
		&payment.TransactionID,
		&payment.PaymentMethod,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("PaymentRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return r.toEntity(&payment), nil
}

// GetByUserID gets payments by user ID
func (r *PaymentRepo) GetByUserID(ctx context.Context, userID int64) ([]*entity.Payment, error) {
	sql, args, err := r.Builder.
		Select("id, user_id, amount, currency, payment_type, status, meter_number, customer_code, description, transaction_id, payment_method, created_at, updated_at").
		From("payments").
		Where("user_id = ?", userID).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PaymentRepo - GetByUserID - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("PaymentRepo - GetByUserID - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.Amount,
			&payment.Currency,
			&payment.PaymentType,
			&payment.Status,
			&payment.MeterNumber,
			&payment.CustomerCode,
			&payment.Description,
			&payment.TransactionID,
			&payment.PaymentMethod,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("PaymentRepo - GetByUserID - rows.Scan: %w", err)
		}
		payments = append(payments, r.toEntity(&payment))
	}

	return payments, nil
}

// UpdateStatus updates payment status
func (r *PaymentRepo) UpdateStatus(ctx context.Context, id int64, status entity.PaymentStatus) error {
	sql, args, err := r.Builder.
		Update("payments").
		Set("status", status).
		Set("updated_at", time.Now()).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("PaymentRepo - UpdateStatus - r.Builder: %w", err)
	}

	result, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PaymentRepo - UpdateStatus - r.Pool.Exec: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("payment with id %d not found", id)
	}

	return nil
}

// CreateHistory creates payment history record
func (r *PaymentRepo) CreateHistory(ctx context.Context, payment *entity.Payment) error {
	sql, args, err := r.Builder.
		Insert("payment_history").
		Columns("payment_id, user_id, status, amount, currency, payment_type, meter_number, customer_code, description, transaction_id, payment_method, created_at").
		Values(payment.ID, payment.UserID, payment.Status, payment.Amount, payment.Currency, payment.PaymentType, payment.MeterNumber, payment.CustomerCode, payment.Description, payment.TransactionID, payment.PaymentMethod, time.Now()).
		ToSql()
	if err != nil {
		return fmt.Errorf("PaymentRepo - CreateHistory - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("PaymentRepo - CreateHistory - r.Pool.Exec: %w", err)
	}

	return nil
}

// toEntity converts database model to entity
func (r *PaymentRepo) toEntity(payment *models.Payment) *entity.Payment {
	return &entity.Payment{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		PaymentType:   entity.PaymentType(payment.PaymentType),
		Status:        entity.PaymentStatus(payment.Status),
		MeterNumber:   payment.MeterNumber,
		CustomerCode:  payment.CustomerCode,
		Description:   payment.Description,
		TransactionID: payment.TransactionID,
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	}
}
