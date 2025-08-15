# pprof Profiling Guide for GoDev Kit

This guide covers how to use pprof for performance profiling in your GoDev Kit application.

## Table of Contents

1. [Overview](#overview)
2. [Available pprof Endpoints](#available-pprof-endpoints)
3. [Basic pprof Usage](#basic-pprof-usage)
4. [Advanced Profiling Techniques](#advanced-profiling-techniques)
5. [Performance Analysis Workflow](#performance-analysis-workflow)
6. [Troubleshooting Common Issues](#troubleshooting-common-issues)
7. [Best Practices](#best-practices)

## Overview

Your application has comprehensive profiling capabilities built-in with pprof integration. The profiling system provides:

- **CPU Profiling**: Identify CPU bottlenecks
- **Memory Profiling**: Detect memory leaks and high memory usage
- **Goroutine Profiling**: Analyze goroutine usage patterns
- **Block Profiling**: Find blocking operations
- **Mutex Profiling**: Identify lock contention

## Available pprof Endpoints

Your application exposes these pprof endpoints at `http://localhost:10000/debug/pprof/`:

| Endpoint | Description | Use Case |
|----------|-------------|----------|
| `/debug/pprof/profile` | CPU profile | CPU bottlenecks |
| `/debug/pprof/heap` | Memory profile | Memory leaks, high usage |
| `/debug/pprof/goroutine` | Goroutine profile | Goroutine analysis |
| `/debug/pprof/block` | Block profile | Blocking operations |
| `/debug/pprof/mutex` | Mutex profile | Lock contention |
| `/debug/pprof/trace` | Execution trace | Detailed execution flow |

## Basic pprof Usage

### 1. CPU Profiling

```bash
# Generate a 30-second CPU profile
curl -o cpu.prof http://localhost:10000/debug/pprof/profile?seconds=30

# Or use go tool pprof directly
go tool pprof http://localhost:10000/debug/pprof/profile?seconds=30
```

### 2. Memory Profiling

```bash
# Get current memory profile
curl -o heap.prof http://localhost:10000/debug/pprof/heap

# Or use go tool pprof directly
go tool pprof http://localhost:10000/debug/pprof/heap
```

### 3. Goroutine Profiling

```bash
# Get goroutine profile
curl -o goroutine.prof http://localhost:10000/debug/pprof/goroutine

# Or use go tool pprof directly
go tool pprof http://localhost:10000/debug/pprof/goroutine
```

## Advanced Profiling Techniques

### 1. Continuous Profiling

Use the provided scripts for continuous monitoring:

```bash
# Start continuous profiling
./scripts/pprof-continuous.sh

# Monitor specific endpoints
./scripts/pprof-monitor.sh
```

### 2. Load Testing with Profiling

```bash
# Generate load while profiling
./scripts/pprof-load-test.sh
```

### 3. Memory Leak Detection

```bash
# Baseline memory profile
curl -o baseline.prof http://localhost:10000/debug/pprof/heap

# Generate load
./scripts/generate-load.sh

# After load memory profile
curl -o after-load.prof http://localhost:10000/debug/pprof/heap

# Compare profiles
go tool pprof -base baseline.prof after-load.prof
```

## Performance Analysis Workflow

### Step 1: Identify the Problem

1. **Monitor application metrics**:
   ```bash
   ./scripts/monitor-functions.sh
   ```

2. **Check runtime statistics**:
   ```bash
   curl http://localhost:10000/debug/stats | jq
   ```

### Step 2: Generate Profiles

1. **CPU Profile** (if CPU-bound):
   ```bash
   go tool pprof http://localhost:10000/debug/pprof/profile?seconds=30
   ```

2. **Memory Profile** (if memory-bound):
   ```bash
   go tool pprof http://localhost:10000/debug/pprof/heap
   ```

3. **Goroutine Profile** (if goroutine-bound):
   ```bash
   go tool pprof http://localhost:10000/debug/pprof/goroutine
   ```

### Step 3: Analyze Profiles

#### CPU Profile Analysis

```bash
go tool pprof cpu.prof

# Interactive commands:
(pprof) top                    # Show top functions by CPU usage
(pprof) top10                  # Show top 10 functions
(pprof) list <function_name>   # Show source code for function
(pprof) web                    # Open web interface (requires graphviz)
(pprof) pdf                    # Generate PDF report
(pprof) svg                    # Generate SVG report
```

#### Memory Profile Analysis

```bash
go tool pprof heap.prof

# Interactive commands:
(pprof) top                    # Show top functions by memory usage
(pprof) top10 -cum            # Show cumulative memory usage
(pprof) list <function_name>   # Show source code for function
(pprof) web                    # Open web interface
(pprof) traces                 # Show allocation traces
```

#### Goroutine Profile Analysis

```bash
go tool pprof goroutine.prof

# Interactive commands:
(pprof) top                    # Show top goroutine types
(pprof) traces                 # Show goroutine stack traces
(pprof) web                    # Open web interface
```

### Step 4: Optimize and Re-test

1. **Implement optimizations** based on profile analysis
2. **Re-run profiling** to measure improvements
3. **Compare before/after profiles**:
   ```bash
   go tool pprof -base before.prof after.prof
   ```

## Troubleshooting Common Issues

### 1. Profile Shows "No samples"

**Cause**: Not enough load or profiling duration too short
**Solution**: 
- Increase profiling duration: `?seconds=60`
- Generate more load during profiling
- Check if profiling is enabled in config

### 2. High Memory Usage

**Analysis Steps**:
```bash
# Get memory profile
go tool pprof http://localhost:10000/debug/pprof/heap

# Check for memory leaks
(pprof) top -cum
(pprof) traces
```

**Common Causes**:
- Unclosed resources (files, connections)
- Large objects not being garbage collected
- Memory fragmentation

### 3. High CPU Usage

**Analysis Steps**:
```bash
# Get CPU profile
go tool pprof http://localhost:10000/debug/pprof/profile?seconds=30

# Analyze hot paths
(pprof) top
(pprof) list <hot_function>
```

**Common Causes**:
- Inefficient algorithms
- Excessive logging
- Tight loops
- Blocking operations

### 4. Goroutine Leaks

**Analysis Steps**:
```bash
# Get goroutine profile
go tool pprof http://localhost:10000/debug/pprof/goroutine

# Check goroutine counts over time
curl http://localhost:10000/debug/stats | jq '.goroutines'
```

**Common Causes**:
- Unclosed channels
- Blocked goroutines
- Infinite loops

## Best Practices

### 1. Profiling in Production

```yaml
# config.yaml
PROFILING:
  ENABLED: true
  PATH: "/debug"
  CPU_PROFILE_DURATION: 30
  MEMORY_PROFILE_INTERVAL: 60
```

### 2. Profile Collection Strategy

- **Development**: Continuous profiling
- **Staging**: Load testing with profiling
- **Production**: Periodic profiling during peak loads

### 3. Profile Storage and Analysis

```bash
# Create timestamped profiles
timestamp=$(date +%Y%m%d_%H%M%S)
curl -o "profiles/cpu_${timestamp}.prof" http://localhost:10000/debug/pprof/profile?seconds=30
curl -o "profiles/heap_${timestamp}.prof" http://localhost:10000/debug/pprof/heap
```

### 4. Automated Profiling

Use the provided scripts for automated profiling:

```bash
# Continuous monitoring
./scripts/pprof-continuous.sh &

# Load testing with profiling
./scripts/pprof-load-test.sh

# Memory leak detection
./scripts/pprof-memory-leak.sh
```

### 5. Profile Comparison

```bash
# Compare profiles
go tool pprof -base baseline.prof current.prof

# Generate diff report
go tool pprof -diff_base baseline.prof current.prof
```

## Integration with Monitoring

Your application integrates with Prometheus for metrics collection:

```bash
# View Prometheus metrics
curl http://localhost:10000/metrics | grep function_

# Key metrics to monitor:
# - function_duration_seconds
# - function_calls_total
# - function_errors_total
# - function_memory_bytes
# - goroutines_total
# - heap_alloc_bytes
```

## Example Analysis Session

```bash
# 1. Start your application
go run cmd/app/main.go

# 2. Generate load
./scripts/generate-load.sh

# 3. Get CPU profile
go tool pprof http://localhost:10000/debug/pprof/profile?seconds=30

# 4. Analyze in interactive mode
(pprof) top
(pprof) list <function_name>
(pprof) web

# 5. Get memory profile
go tool pprof http://localhost:10000/debug/pprof/heap

# 6. Check goroutines
go tool pprof http://localhost:10000/debug/pprof/goroutine
```

## Additional Resources

- [Go pprof Documentation](https://golang.org/pkg/net/http/pprof/)
- [Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
- [Go Performance Best Practices](https://golang.org/doc/effective_go.html#performance)

## Scripts Reference

- `scripts/pprof-continuous.sh`: Continuous profiling
- `scripts/pprof-monitor.sh`: Real-time monitoring
- `scripts/pprof-load-test.sh`: Load testing with profiling
- `scripts/pprof-memory-leak.sh`: Memory leak detection
- `scripts/pprof-analyze.sh`: Automated profile analysis
