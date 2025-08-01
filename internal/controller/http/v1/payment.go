package v1

import (
	"net/http"
	"strconv"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/ducnpdev/godev-kit/internal/usecase/payment"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// PaymentController represents payment HTTP controller
type PaymentController struct {
	paymentUseCase *payment.PaymentUseCase
	logger         *zerolog.Logger
}

// NewPaymentController creates new payment controller
func NewPaymentController(paymentUseCase *payment.PaymentUseCase, logger *zerolog.Logger) *PaymentController {
	return &PaymentController{
		paymentUseCase: paymentUseCase,
		logger:         logger,
	}
}

// RegisterPayment registers a new payment
// @Summary Register a new payment
// @Description Register a new payment for electric bill and send to Kafka for processing
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body request.PaymentRequest true "Payment request"
// @Success 201 {object} response.PaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/payments [post]
func (c *PaymentController) RegisterPayment(ctx *gin.Context) {
	var req request.PaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Error().Err(err).Msg("Failed to bind payment request")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Convert request to entity
	paymentReq := &entity.PaymentRequest{
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		PaymentType:   entity.PaymentType(req.PaymentType),
		MeterNumber:   req.MeterNumber,
		CustomerCode:  req.CustomerCode,
		Description:   req.Description,
		PaymentMethod: req.PaymentMethod,
	}

	// Register payment
	paymentResp, err := c.paymentUseCase.RegisterPayment(ctx, paymentReq)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to register payment")
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	resp := response.PaymentResponse{
		ID:            paymentResp.ID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentType:   string(paymentResp.PaymentType),
		Status:        string(paymentResp.Status),
		MeterNumber:   paymentResp.MeterNumber,
		CustomerCode:  paymentResp.CustomerCode,
		Description:   paymentResp.Description,
		TransactionID: paymentResp.TransactionID,
		PaymentMethod: paymentResp.PaymentMethod,
		CreatedAt:     paymentResp.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, resp)
}

// GetPaymentByID gets payment by ID
// @Summary Get payment by ID
// @Description Get payment details by ID
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} response.PaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/payments/{id} [get]
func (c *PaymentController) GetPaymentByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.logger.Error().Err(err).Str("id", idStr).Msg("Invalid payment ID")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid payment ID",
			Message: "Payment ID must be a valid integer",
		})
		return
	}

	paymentResp, err := c.paymentUseCase.GetPaymentByID(ctx, id)
	if err != nil {
		c.logger.Error().Err(err).Int64("payment_id", id).Msg("Failed to get payment")
		if err.Error() == "payment not found" {
			ctx.JSON(http.StatusNotFound, response.ErrorResponse{
				Error:   "Payment not found",
				Message: "Payment with the specified ID was not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	resp := response.PaymentResponse{
		ID:            paymentResp.ID,
		UserID:        paymentResp.UserID,
		Amount:        paymentResp.Amount,
		Currency:      paymentResp.Currency,
		PaymentType:   string(paymentResp.PaymentType),
		Status:        string(paymentResp.Status),
		MeterNumber:   paymentResp.MeterNumber,
		CustomerCode:  paymentResp.CustomerCode,
		Description:   paymentResp.Description,
		TransactionID: paymentResp.TransactionID,
		PaymentMethod: paymentResp.PaymentMethod,
		CreatedAt:     paymentResp.CreatedAt,
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetPaymentsByUserID gets payments by user ID
// @Summary Get payments by user ID
// @Description Get all payments for a specific user
// @Tags payments
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {array} response.PaymentResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /v1/users/{user_id}/payments [get]
func (c *PaymentController) GetPaymentsByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.logger.Error().Err(err).Str("user_id", userIDStr).Msg("Invalid user ID")
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Invalid user ID",
			Message: "User ID must be a valid integer",
		})
		return
	}

	payments, err := c.paymentUseCase.GetPaymentsByUserID(ctx, userID)
	if err != nil {
		c.logger.Error().Err(err).Int64("user_id", userID).Msg("Failed to get payments by user ID")
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "Internal server error",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	responses := make([]response.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = response.PaymentResponse{
			ID:            payment.ID,
			UserID:        payment.UserID,
			Amount:        payment.Amount,
			Currency:      payment.Currency,
			PaymentType:   string(payment.PaymentType),
			Status:        string(payment.Status),
			MeterNumber:   payment.MeterNumber,
			CustomerCode:  payment.CustomerCode,
			Description:   payment.Description,
			TransactionID: payment.TransactionID,
			PaymentMethod: payment.PaymentMethod,
			CreatedAt:     payment.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, responses)
}
