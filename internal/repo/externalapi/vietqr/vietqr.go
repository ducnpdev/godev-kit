package vietqr

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/vietqr"
)

// VietQRRepo is the interface for the vietqr repository.
type VietQRRepo interface {
	GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (string, error)
	InquiryQR(ctx context.Context, id string) (*entity.VietQR, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type vietQRRepo struct {
	// Add any dependencies here, like an HTTP client
}

// NewVietQRRepo creates a new vietqr repository.
func NewVietQRRepo() VietQRRepo {
	return &vietQRRepo{}
}

func (r *vietQRRepo) GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (string, error) {
	qrRequest := vietqr.RequestGenerateViQR{
		MerchantAccountInformation: vietqr.MerchantAccountInformation{
			AccountNo: req.AccountNo,
		},
		TransactionAmount: req.Amount,
		AdditionalDataFieldTemplate: vietqr.AdditionalDataFieldTemplate{
			Description: req.Description,
		},
		Mcc:          req.MCC,
		ReceiverName: req.ReceiverName,
	}

	return vietqr.GenerateViQR(qrRequest), nil
}

func (r *vietQRRepo) InquiryQR(ctx context.Context, id string) (*entity.VietQR, error) {
	// TODO: Implement QR inquiry logic
	return &entity.VietQR{ID: id, Status: "inquired"}, nil
}

func (r *vietQRRepo) UpdateStatus(ctx context.Context, id, status string) error {
	// TODO: Implement status update logic
	return nil
}
