package v1

import (
	"net/http"

	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/request"
	"github.com/ducnpdev/godev-kit/internal/controller/http/v1/response"
	"github.com/ducnpdev/godev-kit/internal/entity"
	"github.com/gin-gonic/gin"
)

// @Summary     Set value
// @Description Set a key-value pair in Redis
// @ID          set-value
// @Tags  	    redis
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body request.RedisValue true "Set value"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/redis/set [post]
func (v1 *V1) setValue(c *gin.Context) {
	var req request.RedisValue
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.l.Error(err, "http - v1 - setValue - c.ShouldBindJSON")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := v1.redis.SetValue(c.Request.Context(), entity.RedisValue{
		Key:   req.Key,
		Value: req.Value,
	})
	if err != nil {
		v1.l.Error(err, "http - v1 - setValue - r.r.SetValue")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, response.Success{Message: "success"})
}

// @Summary     Get value
// @Description Get a value from Redis by key
// @ID          get-value
// @Tags  	    redis
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       key path string true "Key"
// @Success     200 {object} response.RedisValue
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/redis/get/{key} [get]
func (v1 *V1) getValue(c *gin.Context) {
	key := c.Param("key")

	val, err := v1.redis.GetValue(c.Request.Context(), key)
	if err != nil {
		v1.l.Error(err, "http - v1 - getValue - r.r.GetValue")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, response.NewRedisValue(val))
}

// @Summary     Update shipper location
// @Description Update the latest location of a shipper in Redis
// @ID          update-shipper-location
// @Tags  	    redis
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body request.ShipperLocation true "Shipper location"
// @Success     200 {object} response.Success
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/redis/shipper/location [post]
func (v1 *V1) UpdateShipperLocation(c *gin.Context) {
	var req request.ShipperLocation
	if err := c.ShouldBindJSON(&req); err != nil {
		v1.l.Error(err, "http - v1 - UpdateShipperLocation - c.ShouldBindJSON")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := v1.v.Struct(req); err != nil {
		v1.l.Error(err, "http - v1 - UpdateShipperLocation - v1.v.Struct")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	// Compose entity
	loc := entity.ShipperLocation{
		ShipperID: req.ShipperID,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: req.Timestamp,
	}

	err := v1.shipperLocation.UpdateLocation(c.Request.Context(), loc)
	if err != nil {
		v1.l.Error(err, "http - v1 - UpdateShipperLocation - usecase.UpdateLocation")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, response.Success{Message: "shipper location updated"})
}

// GetShipperLocation retrieves the latest location of a shipper
// @Summary     Get shipper location
// @Description Get the latest location of a shipper from Redis
// @ID          get-shipper-location
// @Tags   redis
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       shipper_id path string true "Shipper ID"
// @Success     200 {object} entity.ShipperLocation
// @Failure     400 {object} response.Error
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /v1/redis/shipper/location/{shipper_id} [get]
func (v1 *V1) GetShipperLocation(c *gin.Context) {
	shipperID := c.Param("shipper_id")
	if shipperID == "" {
		errorResponse(c, http.StatusBadRequest, "shipper_id is required")
		return
	}
	loc, err := v1.shipperLocation.GetLocation(c.Request.Context(), shipperID)
	if err != nil {
		v1.l.Error(err, "http - v1 - GetShipperLocation - usecase.GetLocation")
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, loc)
}
