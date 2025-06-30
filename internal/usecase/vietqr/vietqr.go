package vietqr

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/externalapi/vietqr"
	"github.com/google/uuid"
)

// VietQRUseCase is the interface for the vietqr use case.
type VietQRUseCase interface {
	GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (*entity.VietQR, error)
	InquiryQR(ctx context.Context, id string) (*entity.VietQR, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

// VietQRPersistentRepo is the interface for the vietqr persistent repository.
type VietQRPersistentRepo interface {
	Store(ctx context.Context, qr entity.VietQR) error
}

type vietQRUseCase struct {
	repo           vietqr.VietQRRepo
	persistentRepo VietQRPersistentRepo
}

// NewVietQRUseCase creates a new vietqr use case.
func NewVietQRUseCase(repo vietqr.VietQRRepo, persistentRepo VietQRPersistentRepo) VietQRUseCase {
	return &vietQRUseCase{repo: repo, persistentRepo: persistentRepo}
}

func (uc *vietQRUseCase) GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (*entity.VietQR, error) {
	content, err := uc.repo.GenerateQR(ctx, req)
	if err != nil {
		return nil, err
	}

	qrEntity := &entity.VietQR{
		ID:      uuid.NewString(),
		Status:  "generated",
		Content: content,
	}

	if err := uc.persistentRepo.Store(ctx, *qrEntity); err != nil {
		return nil, err
	}

	return qrEntity, nil
}

func (uc *vietQRUseCase) InquiryQR(ctx context.Context, id string) (*entity.VietQR, error) {
	return uc.repo.InquiryQR(ctx, id)
}

func (uc *vietQRUseCase) UpdateStatus(ctx context.Context, id, status string) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}
