package vietqr

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/repo/externalapi/vietqr"
)

// VietQRUseCase is the interface for the vietqr use case.
type VietQRUseCase interface {
	GenerateQR(ctx context.Context) (*entity.VietQR, error)
	InquiryQR(ctx context.Context, id string) (*entity.VietQR, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type vietQRUseCase struct {
	repo vietqr.VietQRRepo
}

// NewVietQRUseCase creates a new vietqr use case.
func NewVietQRUseCase(repo vietqr.VietQRRepo) VietQRUseCase {
	return &vietQRUseCase{repo: repo}
}

func (uc *vietQRUseCase) GenerateQR(ctx context.Context) (*entity.VietQR, error) {
	return uc.repo.GenerateQR(ctx)
}

func (uc *vietQRUseCase) InquiryQR(ctx context.Context, id string) (*entity.VietQR, error) {
	return uc.repo.InquiryQR(ctx, id)
}

func (uc *vietQRUseCase) UpdateStatus(ctx context.Context, id, status string) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}
