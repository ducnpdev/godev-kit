package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig represents the configuration for timeout middleware
type TimeoutConfig struct {
	Timeout time.Duration
	// Optional: Custom timeout response
	TimeoutResponse gin.H
	// Optional: Skip timeout for specific paths
	SkipPaths []string
}

// DefaultTimeoutConfig returns a default timeout configuration
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second,
		TimeoutResponse: gin.H{
			"error":   "Request timeout",
			"code":    http.StatusRequestTimeout,
			"message": "The server took too long to process your request",
		},
		SkipPaths: []string{"/health", "/metrics"},
	}
}

// TimeoutMiddleware creates a middleware that enforces request timeouts
func TimeoutMiddleware(config ...TimeoutConfig) gin.HandlerFunc {
	cfg := DefaultTimeoutConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return func(c *gin.Context) {
		// Check if this path should skip timeout
		for _, path := range cfg.SkipPaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.Timeout)
		defer cancel()

		// Replace the request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		finished := make(chan struct{})
		panicChan := make(chan interface{}, 1)

		// Run the request in a goroutine
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			c.Next()
			finished <- struct{}{}
		}()

		// Wait for either completion or timeout
		select {
		case <-finished:
			// Request completed normally
			return

		case p := <-panicChan:
			// Request panicked
			panic(p)

		case <-ctx.Done():
			// Request timed out
			c.Header("Connection", "close")
			c.JSON(http.StatusRequestTimeout, cfg.TimeoutResponse)
			c.Abort()

			// Log the timeout
			fmt.Printf("[TIMEOUT] Request to %s timed out after %v\n",
				c.Request.URL.Path, cfg.Timeout)

			return
		}
	}
}

// TimeoutWithHandler creates a timeout middleware with custom handler
func TimeoutWithHandler(timeout time.Duration, handler gin.HandlerFunc) gin.HandlerFunc {
	return TimeoutMiddleware(TimeoutConfig{
		Timeout: timeout,
		TimeoutResponse: gin.H{
			"error":   "Request timeout",
			"timeout": timeout.String(),
		},
	})
}
