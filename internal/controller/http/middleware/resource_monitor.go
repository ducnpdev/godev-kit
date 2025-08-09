package middleware

import (
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ResourceMonitor tracks active requests and goroutines
type ResourceMonitor struct {
	activeRequests map[string]*RequestInfo
	mutex          sync.RWMutex
}

// RequestInfo holds information about an active request
type RequestInfo struct {
	StartTime   time.Time
	Path        string
	Method      string
	GoroutineID int
	ContextDone <-chan struct{}
}

var (
	globalMonitor = &ResourceMonitor{
		activeRequests: make(map[string]*RequestInfo),
	}
)

// ResourceMonitorMiddleware creates middleware to monitor resource usage
func ResourceMonitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID(c)

		// Track request start
		globalMonitor.TrackRequest(requestID, &RequestInfo{
			StartTime:   time.Now(),
			Path:        c.Request.URL.Path,
			Method:      c.Request.Method,
			GoroutineID: runtime.NumGoroutine(),
			ContextDone: c.Request.Context().Done(),
		})

		// Clean up after request
		defer func() {
			globalMonitor.UntrackRequest(requestID)
		}()

		c.Next()
	}
}

// TrackRequest adds a request to monitoring
func (rm *ResourceMonitor) TrackRequest(id string, info *RequestInfo) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	rm.activeRequests[id] = info
}

// UntrackRequest removes a request from monitoring
func (rm *ResourceMonitor) UntrackRequest(id string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()
	delete(rm.activeRequests, id)
}

// GetActiveRequests returns current active requests
func (rm *ResourceMonitor) GetActiveRequests() map[string]*RequestInfo {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	result := make(map[string]*RequestInfo)
	for k, v := range rm.activeRequests {
		result[k] = v
	}
	return result
}

// CheckForLeaks identifies potentially leaked resources
func (rm *ResourceMonitor) CheckForLeaks(maxDuration time.Duration) []string {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	var leaks []string
	now := time.Now()

	for id, info := range rm.activeRequests {
		if now.Sub(info.StartTime) > maxDuration {
			// Check if context is cancelled but request still active
			select {
			case <-info.ContextDone:
				leaks = append(leaks, "Request "+id+" context cancelled but still active after "+now.Sub(info.StartTime).String())
			default:
				if now.Sub(info.StartTime) > maxDuration*2 {
					leaks = append(leaks, "Request "+id+" running too long: "+now.Sub(info.StartTime).String())
				}
			}
		}
	}

	return leaks
}

// StartLeakDetector starts a background goroutine to detect leaks
func StartLeakDetector(interval time.Duration, maxRequestDuration time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			leaks := globalMonitor.CheckForLeaks(maxRequestDuration)
			if len(leaks) > 0 {
				for _, leak := range leaks {
					// Log leak detection
					// You can integrate with your logging system
					println("[RESOURCE_LEAK_DETECTED]", leak)
				}
			}

			// Log current goroutine count
			goroutineCount := runtime.NumGoroutine()
			activeRequests := len(globalMonitor.GetActiveRequests())

			if goroutineCount > 100 || activeRequests > 50 { // Adjust thresholds as needed
				println("[RESOURCE_WARNING] Goroutines:", goroutineCount, "Active requests:", activeRequests)
			}
		}
	}()
}

// GetResourceStats returns current resource statistics
func GetResourceStats() map[string]interface{} {
	var m runtime.MemStats
	runtime.GC() // Force garbage collection for accurate stats
	runtime.ReadMemStats(&m)

	activeRequests := globalMonitor.GetActiveRequests()

	return map[string]interface{}{
		"goroutines":               runtime.NumGoroutine(),
		"active_requests":          len(activeRequests),
		"memory_alloc_mb":          float64(m.Alloc) / 1024 / 1024,
		"memory_total_alloc_mb":    float64(m.TotalAlloc) / 1024 / 1024,
		"memory_sys_mb":            float64(m.Sys) / 1024 / 1024,
		"gc_runs":                  m.NumGC,
		"longest_request_duration": getLongestRequestDuration(activeRequests),
	}
}

func getLongestRequestDuration(requests map[string]*RequestInfo) time.Duration {
	var longest time.Duration
	now := time.Now()

	for _, info := range requests {
		duration := now.Sub(info.StartTime)
		if duration > longest {
			longest = duration
		}
	}

	return longest
}

func generateRequestID(c *gin.Context) string {
	return c.Request.Method + "-" + c.Request.URL.Path + "-" + time.Now().Format("20060102150405.000000")
}
