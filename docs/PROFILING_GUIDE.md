# Function-Level Resource Monitoring Guide

This guide explains how to monitor and profile your service's resource usage at the function level using the comprehensive profiling system built into godev-kit.

## 🎯 Overview

The profiling system provides:
- **Function-level performance tracking** (duration, memory usage)
- **Automatic HTTP request profiling** via middleware
- **Custom Prometheus metrics** for monitoring
- **pprof integration** for detailed profiling
- **Real-time statistics** and dashboards
- **Error tracking** and performance analysis

## 🚀 Quick Start

### 1. Enable Profiling

Profiling is enabled by default in your configuration:

```yaml
PROFILING:
  ENABLED: true
  PATH: "/debug"
  CPU_PROFILE_DURATION: 30
  MEMORY_PROFILE_INTERVAL: 60
```

### 2. Access Profiling Endpoints

Once your application is running, you can access:

- **Function Profiles**: `http://localhost:10000/debug/profiles`
- **Runtime Stats**: `http://localhost:10000/debug/stats`
- **Memory Info**: `http://localhost:10000/debug/memory`
- **Goroutines**: `http://localhost:10000/debug/goroutines`
- **Prometheus Metrics**: `http://localhost:10000/metrics`

## 📊 Available Endpoints

### Function Profiles

#### Get All Function Profiles
```bash
curl http://localhost:10000/debug/profiles
```

Response:
```json
{
  "profiles": {
    "GET_/v1/users": {
      "name": "GET_/v1/users",
      "call_count": 150,
      "total_duration": "2.5s",
      "max_duration": "150ms",
      "min_duration": "5ms",
      "avg_duration": "16.7ms",
      "total_memory": 1048576,
      "max_memory": 8192,
      "last_call_time": "2024-01-15T10:30:00Z",
      "error_count": 2
    }
  },
  "count": 1
}
```

#### Get Specific Function Profile
```bash
curl http://localhost:10000/debug/profiles/GET_/v1/users
```

### Runtime Statistics

```bash
curl http://localhost:10000/debug/stats
```

Response:
```json
{
  "goroutines": 25,
  "memory": {
    "alloc": 10485760,
    "sys": 20971520,
    "idle": 8388608,
    "inuse": 2097152,
    "released": 0,
    "objects": 50000,
    "total_alloc": 52428800,
    "num_gc": 10
  },
  "gc": {
    "num_gc": 10,
    "pause_total_ns": 5000000,
    "pause_ns": 500000
  }
}
```

### Memory Information

```bash
curl http://localhost:10000/debug/memory
```

### Trigger Garbage Collection

```bash
curl -X POST http://localhost:10000/debug/gc
```

## 🔧 Manual Function Profiling

### Basic Function Profiling

```go
import "github.com/ducnpdev/godev-kit/pkg/profiling"

// Initialize profiler
profiler := profiling.NewProfiler(logger, true, "/debug")

// Profile a simple function
err := profiler.ProfileFunction("my_function", "my_package", func() error {
    // Your function logic here
    time.Sleep(100 * time.Millisecond)
    return nil
})
```

### Function with Context

```go
err := profiler.ProfileFunctionWithContext(ctx, "my_function", "my_package", func(ctx context.Context) error {
    // Your function logic here
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(100 * time.Millisecond):
        return nil
    }
})
```

### In HTTP Handlers (Gin)

```go
func (c *gin.Context) {
    // Use the profiling middleware helper
    err := middleware.ProfileFunction(c, "process_payment", "payment", func() error {
        // Your payment processing logic
        return processPayment()
    })
    
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"status": "success"})
}
```

## 📈 Prometheus Metrics

The profiling system automatically exposes these Prometheus metrics:

### Function Metrics
- `function_duration_seconds` - Duration of function calls
- `function_calls_total` - Total number of function calls
- `function_errors_total` - Total number of function errors
- `function_memory_bytes` - Memory usage per function call

### Runtime Metrics
- `goroutines_total` - Current number of goroutines
- `heap_alloc_bytes` - Current heap memory usage
- `heap_sys_bytes` - Total heap memory from system
- `heap_idle_bytes` - Idle heap memory
- `heap_inuse_bytes` - In-use heap memory
- `heap_released_bytes` - Released heap memory
- `heap_objects_total` - Total number of heap objects

### Example Prometheus Query

```promql
# Average function duration by function name
rate(function_duration_seconds_sum[5m]) / rate(function_duration_seconds_count[5m])

# Function error rate
rate(function_errors_total[5m]) / rate(function_calls_total[5m])

# Memory usage by function
function_memory_bytes_sum / function_memory_bytes_count
```

## 🔍 pprof Integration

The system includes pprof endpoints for detailed profiling:

### CPU Profiling
```bash
# 30-second CPU profile
curl -o cpu.prof http://localhost:10000/debug/pprof/profile?seconds=30

# Analyze with go tool
go tool pprof cpu.prof
```

### Memory Profiling
```bash
# Get memory profile
curl -o memory.prof http://localhost:10000/debug/pprof/heap

# Analyze with go tool
go tool pprof memory.prof
```

### Goroutine Profiling
```bash
# Get goroutine stack traces
curl http://localhost:10000/debug/pprof/goroutine?debug=1
```

## 📊 Monitoring Dashboard

### Grafana Dashboard

Create a Grafana dashboard with these panels:

1. **Function Performance**
   - Average response time by function
   - Request rate by function
   - Error rate by function

2. **Memory Usage**
   - Heap allocation over time
   - Memory usage by function
   - GC frequency and duration

3. **System Resources**
   - Number of goroutines
   - CPU usage
   - Memory allocation

### Example Grafana Queries

```promql
# Top 10 slowest functions
topk(10, rate(function_duration_seconds_sum[5m]) / rate(function_duration_seconds_count[5m]))

# Memory usage trend
rate(function_memory_bytes_sum[5m]) / rate(function_memory_bytes_count[5m])

# Error rate by function
rate(function_errors_total[5m]) / rate(function_calls_total[5m])
```

## 🛠️ Best Practices

### 1. Function Naming

Use descriptive function names for better monitoring:

```go
// Good
profiler.ProfileFunction("process_user_payment", "payment_service", fn)

// Avoid
profiler.ProfileFunction("fn", "pkg", fn)
```

### 2. Package Organization

Group related functions by package:

```go
profiler.ProfileFunction("create_user", "user_service", fn)
profiler.ProfileFunction("update_user", "user_service", fn)
profiler.ProfileFunction("delete_user", "user_service", fn)
```

### 3. Error Handling

Always handle errors properly:

```go
err := profiler.ProfileFunction("critical_operation", "core", func() error {
    // Your logic here
    if err := someOperation(); err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
})

if err != nil {
    // Handle error appropriately
    log.Printf("Critical operation failed: %v", err)
}
```

### 4. Performance Thresholds

Set up alerts for performance issues:

```yaml
# Example alerting rules
groups:
  - name: function_performance
    rules:
      - alert: SlowFunction
        expr: rate(function_duration_seconds_sum[5m]) / rate(function_duration_seconds_count[5m]) > 1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Function {{ $labels.function }} is slow"
          description: "Average duration is {{ $value }}s"

      - alert: HighErrorRate
        expr: rate(function_errors_total[5m]) / rate(function_calls_total[5m]) > 0.05
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "High error rate for {{ $labels.function }}"
          description: "Error rate is {{ $value | humanizePercentage }}"
```

## 🔧 Configuration Options

### Profiling Configuration

```yaml
PROFILING:
  ENABLED: true                    # Enable/disable profiling
  PATH: "/debug"                   # Base path for profiling endpoints
  CPU_PROFILE_DURATION: 30         # CPU profiling duration in seconds
  MEMORY_PROFILE_INTERVAL: 60      # Memory profiling interval in seconds
```

### Metrics Configuration

```yaml
METRICS:
  ENABLED: true
  SKIP_PATHS: "/swagger/*;/metrics;/debug/*"  # Paths to exclude from metrics
  PATH: "/metrics"
```

## 📋 Troubleshooting

### Common Issues

1. **High Memory Usage**
   - Check `/debug/memory` endpoint
   - Look for memory leaks in function profiles
   - Trigger garbage collection with `/debug/gc`

2. **Slow Functions**
   - Review function duration metrics
   - Check for blocking operations
   - Analyze CPU profiles

3. **High Error Rates**
   - Monitor error metrics by function
   - Check application logs
   - Review error handling logic

4. **Goroutine Leaks**
   - Check `/debug/goroutines` endpoint
   - Look for unbounded goroutine creation
   - Review context cancellation

### Debug Commands

```bash
# Check all function profiles
curl http://localhost:10000/debug/profiles | jq

# Monitor specific function
curl http://localhost:10000/debug/profiles/GET_/v1/users | jq

# Check runtime stats
curl http://localhost:10000/debug/stats | jq

# Get memory information
curl http://localhost:10000/debug/memory | jq

# Trigger garbage collection
curl -X POST http://localhost:10000/debug/gc
```

## 📚 Examples

See `examples/profiling_demo.go` for comprehensive examples of:
- Basic function profiling
- Context-aware profiling
- Error handling
- Multiple function calls
- Memory-intensive operations
- Database operations
- API calls
- File operations

## 🎯 Performance Optimization

### Based on Profiling Data

1. **Identify Slow Functions**
   - Look for functions with high average duration
   - Focus on frequently called slow functions

2. **Memory Optimization**
   - Identify memory-intensive functions
   - Look for memory leaks (increasing memory usage over time)

3. **Error Reduction**
   - Focus on functions with high error rates
   - Improve error handling and validation

4. **Resource Management**
   - Monitor goroutine count
   - Check for resource leaks

## 🔗 Related Documentation

- [PostgreSQL Optimization Guide](POSTGRES_OPTIMIZATION.md)
- [Payment System Documentation](PAYMENT_SYSTEM.md)
- [Swagger API Documentation](SWAGGER_GUIDE.md)

## 🤝 Support

For issues with the profiling system:
1. Check the troubleshooting section above
2. Review application logs
3. Use the debugging endpoints
4. Check Prometheus metrics
5. Analyze pprof profiles

The profiling system provides comprehensive insights into your application's performance, helping you identify bottlenecks and optimize resource usage effectively. 