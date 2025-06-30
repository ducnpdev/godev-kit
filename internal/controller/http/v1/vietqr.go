package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/entity"
)

// @Summary     Generate QR Code
// @Description Generate a new VietQR code
// @ID          generate-qr
// @Tags  	    vietqr
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} entity.VietQR
// @Failure     500 {object} response.Error
// @Router      /v1/vietqr/gen [post]
func (v1 *V1) generateQR(c *gin.Context) {
	var req request.GenerateQR
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.l.Error(err, "http - v1 - generateQR - c.ShouldBindJSON")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	qr, err := v1.vietqr.GenerateQR(c.Request.Context(), entity.VietQRGenerateRequest{
		AccountNo:    req.AccountNo,
		Amount:       req.Amount,
		Description:  req.Description,
		MCC:          req.MCC,
		ReceiverName: req.ReceiverName,
	})
	if err != nil {
		v1.l.Error(err, "http - v1 - generateQR - v1.vietqr.GenerateQR")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, qr)
}

// @Summary     Inquiry QR Status
// @Description Get the status of a VietQR code
// @ID          inquiry-qr
// @Tags  	    vietqr
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "QR ID"
// @Success     200 {object} entity.VietQR
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/vietqr/inquiry/{id} [get]
func (v1 *V1) inquiryQR(c *gin.Context) {
	id := c.Param("id")
	qr, err := v1.vietqr.InquiryQR(c.Request.Context(), id)
	if err != nil {
		v1.l.Error(err, "http - v1 - inquiryQR - v1.vietqr.InquiryQR")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, qr)
}

// @Summary     Update QR Status
// @Description Update the status of a VietQR code
// @ID          update-qr-status
// @Tags  	    vietqr
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id path string true "QR ID"
// @Param       request body request.UpdateVietQRStatus true "Update status"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/vietqr/update/{id} [put]
func (v1 *V1) updateStatus(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateVietQRStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.l.Error(err, "http - v1 - updateStatus - c.ShouldBindJSON")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	status := entity.VietQRStatus(req.Status)
	switch status {
	case entity.VietQRStatusInProcess, entity.VietQRStatusPaid, entity.VietQRStatusFail, entity.VietQRStatusTimeout:
		// valid status
	default:
		errorResponse(c, http.StatusBadRequest, "invalid status")
		return
	}

	err := v1.vietqr.UpdateStatus(c.Request.Context(), id, status)
	if err != nil {
		v1.l.Error(err, "http - v1 - updateStatus - v1.vietqr.UpdateStatus")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
