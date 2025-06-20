package v1

import (
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/go-playground/validator/v10"
)

// V1 -.
type V1 struct {
	l logger.Interface
	v *validator.Validate
	//
	t     usecase.Translation
	user  usecase.User
	kafka usecase.Kafka
}
