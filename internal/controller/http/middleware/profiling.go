package middleware

import (
	"context"
	"runtime"
	"time"

	"github.com/ducnpdev/godev-kit/pkg/profiling"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// ProfilingMiddleware creates middleware for automatic function profiling
func ProfilingMiddleware(profiler *profiling.Profiler, logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if profiler == nil {
			c.Next()
			return
		}

		start := time.Now()
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		startMem := m.HeapAlloc

		// Create function name from request
		functionName := c.Request.Method + "_" + c.FullPath()
		if functionName == "" {
			functionName = c.Request.Method + "_" + c.Request.URL.Path
		}

		// Execute the handler
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		runtime.ReadMemStats(&m)
		endMem := m.HeapAlloc
		memoryUsed := int64(endMem - startMem)

		// Log performance metrics
		logger.Debug().
			Str("function", functionName).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", c.Writer.Status()).
			Dur("duration", duration).
			Int64("memory_used", memoryUsed).
			Str("memory_formatted", profiling.FormatBytes(memoryUsed)).
			Str("duration_formatted", profiling.FormatDuration(duration)).
			Msg("Function performance metrics")

		// Update profiler metrics if enabled
		if profiler != nil {
			profiler.ProfileFunction(functionName, "http_handler", func() error {
				// This is just for metrics recording, the actual work is already done
				return nil
			})
		}
	}
}

// ProfilingContextMiddleware adds profiling context to gin context
func ProfilingContextMiddleware(profiler *profiling.Profiler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if profiler != nil {
			c.Set("profiler", profiler)
		}
		c.Next()
	}
}

// ProfileFunction is a helper function to profile any function from gin context
func ProfileFunction(c *gin.Context, functionName, packageName string, fn func() error) error {
	if profilerInterface, exists := c.Get("profiler"); exists {
		if profiler, ok := profilerInterface.(*profiling.Profiler); ok {
			return profiler.ProfileFunction(functionName, packageName, fn)
		}
	}
	return fn()
}

// ProfileFunctionWithContext is a helper function to profile any function with context
func ProfileFunctionWithContext(c *gin.Context, functionName, packageName string, fn func(context.Context) error) error {
	if profilerInterface, exists := c.Get("profiler"); exists {
		if profiler, ok := profilerInterface.(*profiling.Profiler); ok {
			return profiler.ProfileFunctionWithContext(c.Request.Context(), functionName, packageName, fn)
		}
	}
	return fn(c.Request.Context())
}
