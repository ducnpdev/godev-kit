package middleware

import (
	"strconv"
	"strings"

	"github.com/ducnpdev/godev-kit/pkg/logger"
	"github.com/gin-gonic/gin"
)

func buildRequestMessage(c *gin.Context) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.URL.String())
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Writer.Status()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(c.Writer.Size()))

	return result.String()
}

func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log after request is processed
		// l.Info(buildRequestMessage(c))
	}
}
