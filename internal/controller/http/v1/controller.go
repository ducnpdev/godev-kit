package v1

import (
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
	"github.com/ducnpdev/godev-kit/internal/usecase/payment"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	l logger.Interface
	v *validator.Validate
	//
	t                 usecase.Translation
	user              usecase.User
	kafka             usecase.Kafka
	redis             usecase.Redis
	nats              usecase.Nats
	vietqr            usecase.VietQR
	billing           usecase.Billing
	shipperLocation   usecase.ShipperLocation
	paymentController *PaymentController
	billingController *BillingController
}

// NewV1 creates new V1 controller
func NewV1(l logger.Interface, t usecase.Translation, u usecase.User, k usecase.Kafka, r usecase.Redis, n usecase.Nats, v usecase.VietQR, billing usecase.Billing, shipperLocation usecase.ShipperLocation, paymentUseCase *payment.PaymentUseCase, billingUseCase *billing.UseCase) *V1 {
	return &V1{
		l:                 l,
		v:                 validator.New(),
		t:                 t,
		user:              u,
		kafka:             k,
		redis:             r,
		nats:              n,
		vietqr:            v,
		billing:           billing,
		shipperLocation:   shipperLocation,
		paymentController: NewPaymentController(paymentUseCase, l.(*logger.Logger).ZerologPtr()),
		billingController: NewBillingController(billingUseCase, l.(*logger.Logger).ZerologPtr()),
	}
}
