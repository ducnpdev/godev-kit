// Package http implements routing paths. Each services in own file.
package http

import (
	"net/http"

	"github.com/ducnpdev/godev-kit/config"
	_ "github.com/ducnpdev/godev-kit/docs" // Swagger docs
	"github.com/ducnpdev/godev-kit/internal/controller/http/middleware"
	v1 "github.com/ducnpdev/godev-kit/internal/controller/http/v1"
	"github.com/ducnpdev/godev-kit/internal/usecase"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
	"github.com/ducnpdev/godev-kit/internal/usecase/payment"
	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Go Dev Kit Template API
// @description Using a translation service and user management as examples
// @version     1.0
// @host        localhost:10000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func NewRouter(app *gin.Engine, cfg *config.Config, t usecase.Translation, u usecase.User, k usecase.Kafka, r usecase.Redis, n usecase.Nats, v usecase.VietQR, billing usecase.Billing, l logger.Interface, shipperLocation usecase.ShipperLocation, paymentUseCase *payment.PaymentUseCase, billingUseCase *billing.UseCase) {
	// Middleware
	app.Use(middleware.Logger(l))
	// app.Use(middleware.Recovery(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		registerAt := func() string {
			if cfg.Metrics.Path != "" {
				return cfg.Metrics.Path
			}
			return "/metrics"
		}

		app.GET(registerAt(), gin.WrapH(promhttp.Handler()))
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// K8s probe
	app.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Create V1 controller
	v1Controller := v1.NewV1(l, t, u, k, r, n, v, billing, shipperLocation, paymentUseCase, billingUseCase)

	// Routers
	apiV1Group := app.Group("/v1")
	{
		v1.NewTranslationRoutes(apiV1Group, t, l)
		v1.NewUserRoutes(apiV1Group, u, l)
		v1.NewKafkaRoutes(apiV1Group, k, l)
		v1.NewRedisRoutes(apiV1Group, r, l, shipperLocation)
		v1.NewNatsRoutes(apiV1Group, n, l)
		v1.NewVietQRRoutes(apiV1Group, v, l)

		// Payment routes
		v1Controller.RegisterPaymentRoutes(apiV1Group)

		// Billing routes
		v1Controller.RegisterBillingRoutes(apiV1Group)
	}
}
