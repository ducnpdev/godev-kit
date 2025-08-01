package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterPaymentRoutes registers payment routes
func (v *V1) RegisterPaymentRoutes(api *gin.RouterGroup) {
	payments := api.Group("/payments")
	{
		payments.POST("", v.paymentController.RegisterPayment)
		payments.GET("/:id", v.paymentController.GetPaymentByID)
	}

	users := api.Group("/users")
	{
		users.GET("/:user_id/payments", v.paymentController.GetPaymentsByUserID)
	}
}
