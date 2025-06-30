package vietqr

import (
	"context"

	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/vietqr"
)

// VietQRRepo is the interface for the vietqr repository.
type VietQRRepo interface {
	GenerateQR(ctx context.Context, req entity.VietQRGenerateRequest) (string, error)
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
