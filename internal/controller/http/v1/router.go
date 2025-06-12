package v1

import (
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group fiber.Router, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.Get("/history", r.history)
		translationGroup.Post("/do-translate", r.doTranslate)
	}
}

// NewUserRoutes -.
func NewUserRoutes(apiV1Group fiber.Router, u usecase.User, l logger.Interface) {
	r := &V1{user: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	userGroup := apiV1Group.Group("/user")

	{
		userGroup.Post("", r.createUser)
		userGroup.Get("", r.listUsers)
		userGroup.Get("/:id", r.getUser)
		userGroup.Put("/:id", r.updateUser)
		userGroup.Delete("/:id", r.deleteUser)
	}
}
