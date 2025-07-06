package v1

import (
	"net/http"
	"path/filepath"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/usecase/billing"
	"github.com/gin-gonic/gin"
)

// @Summary     Generate Invoice PDF
// @Description Generate a billing payment PDF
// @ID          generate-invoice-pdf
// @Tags   	billing
// @Accept      json
// @Produce     json
// @Param       request body request.GenerateInvoicePDFRequest true "Invoice data"
// @Success     200 {object} response.GenerateInvoicePDFResponse
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/billing/invoice [post]
func (v1 *V1) GenerateInvoicePDF(c *gin.Context) {
	var req request.GenerateInvoicePDFRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request body")
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

	if err := v1.billing.GenerateInvoicePDF(ucData, outputPath); err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to generate PDF")
		return
	}

	c.JSON(http.StatusOK, response.GenerateInvoicePDFResponse{
		FilePath: outputPath,
	})
}
