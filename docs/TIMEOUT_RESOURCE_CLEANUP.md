# Timeout Middleware & Resource Cleanup Improvements

## 🚨 Issues Fixed

### 1. **Goroutine Leak in Timeout Middleware**

**Problem:** The original implementation could leak goroutines when requests timed out because the spawned goroutine would continue running even after the main function returned.

**Fix:** Added proper goroutine lifecycle tracking:
```go
// Track goroutine completion for cleanup
goroutineDone := make(chan struct{})

go func() {
    defer func() {
        close(goroutineDone) // Signal goroutine completion
        // ... rest of cleanup
    }()
    // ... request processing
}()

// Wait for goroutine cleanup in all exit paths
<-goroutineDone
```

### 2. **External API Context Cancellation**

**Problem:** External APIs didn't respect context cancellation, potentially continuing to run after timeout.

**Fix:** Updated `TranslationWebAPI.Translate()` to support context cancellation:
```go
func (t *TranslationWebAPI) Translate(ctx context.Context, translation entity.Translation) (entity.Translation, error) {
    // Run in goroutine with context cancellation support
    select {
    case result := <-resultChan:
        return result.result, result.err
    case <-ctx.Done():
        return entity.Translation{}, fmt.Errorf("context cancelled: %w", ctx.Err())
    }
}
```

### 3. **Resource Monitoring System**

**Added:** Comprehensive resource monitoring to detect memory leaks:
- Active request tracking
- Goroutine count monitoring  
- Memory usage statistics
- Leak detection for long-running requests

## 🛠️ Implementation Details

### Timeout Middleware Improvements

#### **Buffered Channels**
```go
// Prevents blocking if main routine already returned
finished := make(chan struct{}, 1)
panicChan := make(chan interface{}, 1)
```

#### **Non-blocking Channel Operations**
```go
select {
case finished <- struct{}{}:
default: // Prevent blocking if main routine already returned
}
```

#### **Cleanup Goroutine**
```go
// Start cleanup goroutine to wait for the original goroutine
go func() {
    <-goroutineDone
    fmt.Printf("[CLEANUP] Goroutine for %s request cleanup completed\n", c.Request.URL.Path)
}()
```

### Resource Monitoring Features

#### **Active Request Tracking**
- Tracks all active HTTP requests
- Monitors request duration
- Detects context cancellation vs active requests

#### **Memory Leak Detection**
```go
func (rm *ResourceMonitor) CheckForLeaks(maxDuration time.Duration) []string {
    // Identifies requests running too long
    // Detects cancelled contexts still active
    // Returns leak reports
}
```

#### **Real-time Statistics**
- Goroutine count
- Active requests count
- Memory allocation (heap, total, system)
- Garbage collection runs
- Longest running request duration

## 🔧 Configuration

### 1. **Enable Resource Monitoring**

Add to your router setup:
```go
// In your main router setup
middleware.StartLeakDetector(30*time.Second, 60*time.Second) // Check every 30s, max request 60s

// Add resource monitoring middleware
app.Use(middleware.ResourceMonitorMiddleware())

// Add monitoring endpoints
v1.NewMonitoringRoutes(apiV1Group, logger)
```

### 2. **Monitoring Endpoints**

**GET /api/v1/monitoring/resources**
```json
{
  "status": "success",
  "data": {
    "goroutines": 45,
    "active_requests": 3,
    "memory_alloc_mb": 12.5,
    "memory_total_alloc_mb": 156.7,
    "memory_sys_mb": 25.8,
    "gc_runs": 12,
    "longest_request_duration": "2.5s"
  }
}
```

### 3. **Leak Detection Alerts**

The system automatically logs warnings:
```
[RESOURCE_LEAK_DETECTED] Request POST-/api/users-20240115123045.123456 context cancelled but still active after 35s
[RESOURCE_WARNING] Goroutines: 150 Active requests: 75
[CLEANUP] Goroutine for /api/users request cleanup completed
```

## 📊 Context Propagation Verification

### Database Operations ✅
```go
// All database operations properly use context
err = r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
```

### Redis Operations ✅  
```go
// Redis operations have proper timeouts
ctx, cancel := context.WithTimeout(ctx, _defaultTimeout)
defer cancel()
return r.r.Client().Set(ctx, key, value, 0).Err()
```

### External APIs ✅
```go
// Now properly supports context cancellation
translation, err := uc.webAPI.Translate(ctx, t)
```

## 🚀 Benefits

1. **Memory Leak Prevention:** Goroutines are properly cleaned up
2. **Resource Monitoring:** Real-time visibility into system resources
3. **Early Detection:** Automatic leak detection and alerting
4. **Context Propagation:** All layers respect timeout cancellation
5. **Graceful Degradation:** Proper error responses instead of connection drops

## 🔍 Testing Resource Cleanup

### 1. **Load Testing**
```bash
# Send concurrent requests to test goroutine cleanup
for i in {1..100}; do
  curl -X POST http://localhost:10000/api/v1/translation/do-translate \
    -H "Content-Type: application/json" \
    -d '{"source":"en","destination":"vi","original":"hello"}' &
done
```

### 2. **Monitor Resources**
```bash
# Check resource stats
curl http://localhost:10000/api/v1/monitoring/resources

# Watch for leak detection logs
tail -f logs/app.log | grep "RESOURCE_LEAK\|CLEANUP"
```

### 3. **Timeout Testing**
```bash
# Test timeout behavior with slow endpoints
curl -X POST http://localhost:10000/api/v1/slow-endpoint \
  --max-time 35  # Should timeout at 30s with proper cleanup
```

## 🎯 Next Steps

1. **Integrate with Metrics:** Export to Prometheus/Grafana
2. **Alerting:** Set up alerts for resource thresholds
3. **Dashboard:** Create resource monitoring dashboard
4. **Profiling:** Add pprof endpoints for detailed analysis
5. **Circuit Breaker:** Implement circuit breaker for external APIs

## 🚨 Production Recommendations

1. **Monitor Goroutine Count:** Alert if > 1000 goroutines
2. **Track Memory Growth:** Alert on sustained memory increase
3. **Request Duration:** Alert on requests > 60s
4. **Active Requests:** Alert if > 100 concurrent requests
5. **Regular Cleanup:** Schedule periodic resource cleanup checks