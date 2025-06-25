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
