package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterBillingRoutes registers billing routes
func (v *V1) RegisterBillingRoutes(api *gin.RouterGroup) {
	billing := api.Group("/billing")
	{
		billing.POST("/invoice", v.billingController.GenerateInvoicePDF)
	}
}
