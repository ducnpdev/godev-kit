package v1

import (
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group *gin.RouterGroup, t usecase.Translation, l logger.Interface) {
	r := &V1{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	translationGroup := apiV1Group.Group("/translation")

	{
		translationGroup.GET("/history", r.history)
		translationGroup.POST("/do-translate", r.doTranslate)
	}
}

// NewUserRoutes -.
func NewUserRoutes(apiV1Group *gin.RouterGroup, u usecase.User, l logger.Interface) {
	r := &V1{user: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	userGroup := apiV1Group.Group("/user")

	{
		userGroup.POST("", r.CreateUser)
		userGroup.GET("", r.ListUsers)
		userGroup.GET("/:id", r.GetUser)
		userGroup.PUT("/:id", r.UpdateUser)
		userGroup.DELETE("/:id", r.DeleteUser)
	}
	apiV1Group.POST("auth/login", r.LoginUser)
}

// // NewUserRoutes -.
// func NewAuthRoutes(apiV1Group *gin.RouterGroup, u usecase.User, l logger.Interface) {
// 	r := &V1{user: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

// 	userGroup := apiV1Group.Group("/user")

// 	{
// 		userGroup.POST("", r.CreateUser)
// 		userGroup.GET("", r.ListUsers)
// 		userGroup.GET("/:id", r.GetUser)
// 		userGroup.PUT("/:id", r.UpdateUser)
// 		userGroup.DELETE("/:id", r.DeleteUser)
// 	}
// }
