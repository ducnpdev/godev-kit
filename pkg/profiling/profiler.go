package profiling

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

// Profiler provides comprehensive function-level resource monitoring
type Profiler struct {
	logger    zerolog.Logger
	enabled   bool
	debugPath string
	metrics   *FunctionMetrics
	mu        sync.RWMutex
	profiles  map[string]*FunctionProfile
}

// FunctionProfile tracks resource usage for a specific function
type FunctionProfile struct {
	Name          string
	CallCount     int64
	TotalDuration time.Duration
	MaxDuration   time.Duration
	MinDuration   time.Duration
	AvgDuration   time.Duration
	TotalMemory   int64
	MaxMemory     int64
	LastCallTime  time.Time
	ErrorCount    int64
	mu            sync.RWMutex
}

// FunctionMetrics holds Prometheus metrics for function monitoring
type FunctionMetrics struct {
	FunctionDuration *prometheus.HistogramVec
	FunctionCalls    *prometheus.CounterVec
	FunctionErrors   *prometheus.CounterVec
	FunctionMemory   *prometheus.HistogramVec
	Goroutines       prometheus.Gauge
	HeapAlloc        prometheus.Gauge
	HeapSys          prometheus.Gauge
	HeapIdle         prometheus.Gauge
	HeapInuse        prometheus.Gauge
	HeapReleased     prometheus.Gauge
	HeapObjects      prometheus.Gauge
}

// NewProfiler creates a new profiler instance
func NewProfiler(logger zerolog.Logger, enabled bool, debugPath string) *Profiler {
	p := &Profiler{
		logger:    logger,
		enabled:   enabled,
		debugPath: debugPath,
		profiles:  make(map[string]*FunctionProfile),
	}

	if enabled {
		p.metrics = p.createMetrics()
		p.startMetricsCollection()
	}

	return p
}

// createMetrics initializes Prometheus metrics
func (p *Profiler) createMetrics() *FunctionMetrics {
	return &FunctionMetrics{
		FunctionDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "function_duration_seconds",
				Help:    "Duration of function calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"function", "package"},
		),
		FunctionCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "function_calls_total",
				Help: "Total number of function calls",
			},
			[]string{"function", "package"},
		),
		FunctionErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "function_errors_total",
				Help: "Total number of function errors",
			},
			[]string{"function", "package"},
		),
		FunctionMemory: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "function_memory_bytes",
				Help:    "Memory usage per function call in bytes",
				Buckets: []float64{1024, 10240, 102400, 1048576, 10485760, 104857600},
			},
			[]string{"function", "package"},
		),
		Goroutines: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "goroutines_total",
				Help: "Current number of goroutines",
			},
		),
		HeapAlloc: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_alloc_bytes",
				Help: "Current heap memory usage in bytes",
			},
		),
		HeapSys: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_sys_bytes",
				Help: "Total heap memory from system in bytes",
			},
		),
		HeapIdle: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_idle_bytes",
				Help: "Idle heap memory in bytes",
			},
		),
		HeapInuse: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_inuse_bytes",
				Help: "In-use heap memory in bytes",
			},
		),
		HeapReleased: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_released_bytes",
				Help: "Released heap memory in bytes",
			},
		),
		HeapObjects: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "heap_objects_total",
				Help: "Total number of heap objects",
			},
		),
	}
}

// startMetricsCollection starts periodic metrics collection
func (p *Profiler) startMetricsCollection() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			p.collectRuntimeMetrics()
		}
	}()
}

// collectRuntimeMetrics collects runtime statistics
func (p *Profiler) collectRuntimeMetrics() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	p.metrics.Goroutines.Set(float64(runtime.NumGoroutine()))
	p.metrics.HeapAlloc.Set(float64(m.HeapAlloc))
	p.metrics.HeapSys.Set(float64(m.HeapSys))
	p.metrics.HeapIdle.Set(float64(m.HeapIdle))
	p.metrics.HeapInuse.Set(float64(m.HeapInuse))
	p.metrics.HeapReleased.Set(float64(m.HeapReleased))
	p.metrics.HeapObjects.Set(float64(m.HeapObjects))
}

// ProfileFunction wraps a function to track its performance
func (p *Profiler) ProfileFunction(functionName, packageName string, fn func() error) error {
	if !p.enabled {
		return fn()
	}

	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startMem := m.HeapAlloc

	// Record function call
	p.metrics.FunctionCalls.WithLabelValues(functionName, packageName).Inc()

	// Execute function
	err := fn()

	// Calculate metrics
	duration := time.Since(start)
	runtime.ReadMemStats(&m)
	endMem := m.HeapAlloc
	memoryUsed := int64(endMem - startMem)

	// Record metrics
	p.metrics.FunctionDuration.WithLabelValues(functionName, packageName).Observe(duration.Seconds())
	p.metrics.FunctionMemory.WithLabelValues(functionName, packageName).Observe(float64(memoryUsed))

	if err != nil {
		p.metrics.FunctionErrors.WithLabelValues(functionName, packageName).Inc()
	}

	// Update internal profile
	p.updateFunctionProfile(functionName, duration, memoryUsed, err != nil)

	return err
}

// ProfileFunctionWithContext wraps a function with context
func (p *Profiler) ProfileFunctionWithContext(ctx context.Context, functionName, packageName string, fn func(context.Context) error) error {
	if !p.enabled {
		return fn(ctx)
	}

	start := time.Now()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startMem := m.HeapAlloc

	// Record function call
	p.metrics.FunctionCalls.WithLabelValues(functionName, packageName).Inc()

	// Execute function
	err := fn(ctx)

	// Calculate metrics
	duration := time.Since(start)
	runtime.ReadMemStats(&m)
	endMem := m.HeapAlloc
	memoryUsed := int64(endMem - startMem)

	// Record metrics
	p.metrics.FunctionDuration.WithLabelValues(functionName, packageName).Observe(duration.Seconds())
	p.metrics.FunctionMemory.WithLabelValues(functionName, packageName).Observe(float64(memoryUsed))

	if err != nil {
		p.metrics.FunctionErrors.WithLabelValues(functionName, packageName).Inc()
	}

	// Update internal profile
	p.updateFunctionProfile(functionName, duration, memoryUsed, err != nil)

	return err
}

// updateFunctionProfile updates internal function profile
func (p *Profiler) updateFunctionProfile(functionName string, duration time.Duration, memoryUsed int64, hasError bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	profile, exists := p.profiles[functionName]
	if !exists {
		profile = &FunctionProfile{
			Name:        functionName,
			MinDuration: duration,
			MaxMemory:   memoryUsed,
		}
		p.profiles[functionName] = profile
	}

	profile.mu.Lock()
	defer profile.mu.Unlock()

	profile.CallCount++
	profile.TotalDuration += duration
	profile.TotalMemory += memoryUsed
	profile.LastCallTime = time.Now()

	if duration > profile.MaxDuration {
		profile.MaxDuration = duration
	}
	if duration < profile.MinDuration {
		profile.MinDuration = duration
	}
	if memoryUsed > profile.MaxMemory {
		profile.MaxMemory = memoryUsed
	}

	profile.AvgDuration = profile.TotalDuration / time.Duration(profile.CallCount)

	if hasError {
		profile.ErrorCount++
	}
}

// GetFunctionProfile returns profile for a specific function
func (p *Profiler) GetFunctionProfile(functionName string) *FunctionProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if profile, exists := p.profiles[functionName]; exists {
		profile.mu.RLock()
		defer profile.mu.RUnlock()
		return profile
	}
	return nil
}

// GetAllProfiles returns all function profiles
func (p *Profiler) GetAllProfiles() map[string]*FunctionProfile {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]*FunctionProfile)
	for name, profile := range p.profiles {
		profile.mu.RLock()
		result[name] = profile
		profile.mu.RUnlock()
	}
	return result
}

// SetupRoutes sets up profiling and monitoring routes
func (p *Profiler) SetupRoutes(router *gin.Engine) {
	if !p.enabled {
		return
	}

	debugGroup := router.Group(p.debugPath)
	{
		// pprof endpoints
		debugGroup.GET("/pprof/", gin.WrapF(pprof.Index))
		debugGroup.GET("/pprof/cmdline", gin.WrapF(pprof.Cmdline))
		debugGroup.GET("/pprof/profile", gin.WrapF(pprof.Profile))
		debugGroup.GET("/pprof/symbol", gin.WrapF(pprof.Symbol))
		debugGroup.GET("/pprof/trace", gin.WrapF(pprof.Trace))
		debugGroup.GET("/pprof/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
		debugGroup.GET("/pprof/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
		debugGroup.GET("/pprof/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
		debugGroup.GET("/pprof/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
		debugGroup.GET("/pprof/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))

		// Custom profiling endpoints
		debugGroup.GET("/profiles", p.getProfilesHandler)
		debugGroup.GET("/profiles/:function", p.getFunctionProfileHandler)
		debugGroup.GET("/stats", p.getStatsHandler)
		debugGroup.GET("/memory", p.getMemoryHandler)
		debugGroup.GET("/goroutines", p.getGoroutinesHandler)
		debugGroup.POST("/gc", p.triggerGCHandler)
	}
}

// getProfilesHandler returns all function profiles
func (p *Profiler) getProfilesHandler(c *gin.Context) {
	profiles := p.GetAllProfiles()
	c.JSON(http.StatusOK, gin.H{
		"profiles": profiles,
		"count":    len(profiles),
	})
}

// getFunctionProfileHandler returns profile for a specific function
func (p *Profiler) getFunctionProfileHandler(c *gin.Context) {
	functionName := c.Param("function")
	profile := p.GetFunctionProfile(functionName)

	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Function profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// getStatsHandler returns runtime statistics
func (p *Profiler) getStatsHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"goroutines": runtime.NumGoroutine(),
		"memory": gin.H{
			"alloc":       m.HeapAlloc,
			"sys":         m.HeapSys,
			"idle":        m.HeapIdle,
			"inuse":       m.HeapInuse,
			"released":    m.HeapReleased,
			"objects":     m.HeapObjects,
			"total_alloc": m.TotalAlloc,
			"num_gc":      m.NumGC,
		},
		"gc": gin.H{
			"num_gc":         m.NumGC,
			"pause_total_ns": m.PauseTotalNs,
			"pause_ns":       m.PauseNs[(m.NumGC+255)%256],
		},
	})
}

// getMemoryHandler returns detailed memory information
func (p *Profiler) getMemoryHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"heap": gin.H{
			"alloc":    m.HeapAlloc,
			"sys":      m.HeapSys,
			"idle":     m.HeapIdle,
			"inuse":    m.HeapInuse,
			"released": m.HeapReleased,
			"objects":  m.HeapObjects,
		},
		"stack": gin.H{
			"inuse": m.StackInuse,
			"sys":   m.StackSys,
		},
		"mspan": gin.H{
			"inuse": m.MSpanInuse,
			"sys":   m.MSpanSys,
		},
		"mcache": gin.H{
			"inuse": m.MCacheInuse,
			"sys":   m.MCacheSys,
		},
		"buck_hash_sys": m.BuckHashSys,
		"gc_sys":        m.GCSys,
		"other_sys":     m.OtherSys,
		"next_gc":       m.NextGC,
		"last_gc":       m.LastGC,
	})
}

// getGoroutinesHandler returns goroutine information
func (p *Profiler) getGoroutinesHandler(c *gin.Context) {
	// Get stack traces for all goroutines
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)

	c.Data(http.StatusOK, "text/plain", buf[:n])
}

// triggerGCHandler triggers garbage collection
func (p *Profiler) triggerGCHandler(c *gin.Context) {
	before := new(runtime.MemStats)
	runtime.ReadMemStats(before)

	runtime.GC()

	after := new(runtime.MemStats)
	runtime.ReadMemStats(after)

	c.JSON(http.StatusOK, gin.H{
		"message": "Garbage collection completed",
		"before": gin.H{
			"heap_alloc": before.HeapAlloc,
			"heap_sys":   before.HeapSys,
		},
		"after": gin.H{
			"heap_alloc": after.HeapAlloc,
			"heap_sys":   after.HeapSys,
		},
		"freed": gin.H{
			"heap_alloc": before.HeapAlloc - after.HeapAlloc,
			"heap_sys":   before.HeapSys - after.HeapSys,
		},
	})
}

// StartCPUProfile starts CPU profiling
func (p *Profiler) StartCPUProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "CPU profiling not implemented in this version"}`))
}

// GetMemoryProfile returns memory profile
func (p *Profiler) GetMemoryProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error": "Memory profiling not implemented in this version"}`))
}

// FormatBytes formats bytes to human readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration formats duration to human readable format
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2f μs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2f s", d.Seconds())
}
