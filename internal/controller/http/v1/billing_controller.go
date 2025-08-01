package v1

import (
	"net/http"
	"path/filepath"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// BillingController represents billing HTTP controller
type BillingController struct {
	billingUseCase *billing.UseCase
	logger         *zerolog.Logger
}

// NewBillingController creates new billing controller
func NewBillingController(billingUseCase *billing.UseCase, logger *zerolog.Logger) *BillingController {
	return &BillingController{
		billingUseCase: billingUseCase,
		logger:         logger,
	}
}

// GenerateInvoicePDF generates a PDF invoice
// @Summary Generate Invoice PDF
// @Description Generate a billing payment PDF
// @Tags billing
// @Accept json
// @Produce json
// @Param request body request.GenerateInvoicePDFRequest true "Invoice data"
// @Success 200 {object} response.GenerateInvoicePDFResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/billing/invoice [post]
func (c *BillingController) GenerateInvoicePDF(ctx *gin.Context) {
	var req request.GenerateInvoicePDFRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to bind billing request")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Map request to usecase.InvoiceData
	ucData := billing.InvoiceData{
		Number:      req.Number,
		Date:        req.Date,
		BilledTo:    req.BilledTo,
		CompanyInfo: req.CompanyInfo,
		Items:       make([]billing.InvoiceItem, len(req.Items)),
		Subtotal:    req.Subtotal,
		Discount:    req.Discount,
		TaxRate:     req.TaxRate,
		Tax:         req.Tax,
		Total:       req.Total,
		Terms:       req.Terms,
		BankDetails: req.BankDetails,
	}
	for i, item := range req.Items {
		ucData.Items[i] = billing.InvoiceItem{
			Description: item.Description,
			UnitCost:    item.UnitCost,
			Qty:         item.Qty,
			Amount:      item.Amount,
		}
	}

	// Generate file path (could use a UUID or timestamp for uniqueness)
	fileName := "invoice_" + ucData.Number + ".pdf"
	outputPath := filepath.Join("./", fileName)

	if err := c.billingUseCase.GenerateInvoicePDF(ucData, outputPath); err != nil {
		c.logger.Error().Err(err).Msg("Failed to generate PDF")
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response.GenerateInvoicePDFResponse{
		FilePath: outputPath,
	})
} 