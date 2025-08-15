# pprof Quick Start Guide

This guide provides immediate steps to start using pprof for performance profiling in your GoDev Kit application.

## Prerequisites

1. **Install required tools**:
   ```bash
   # Install jq for JSON processing
   brew install jq  # macOS
   # or
   sudo apt-get install jq  # Ubuntu/Debian
   
   # Install bc for calculations
   brew install bc  # macOS
   # or
   sudo apt-get install bc  # Ubuntu/Debian
   ```

2. **Ensure Go is installed**:
   ```bash
   go version
   ```

## Quick Start Steps

### 1. Start Your Application

```bash
# Start the application
go run cmd/app/main.go
```

Your application will be available at `http://localhost:10000`

### 2. Verify Profiling is Enabled

```bash
# Check if profiling endpoints are available
curl http://localhost:10000/debug/pprof/
```

You should see a list of available pprof endpoints.

### 3. Basic Profiling Commands

#### CPU Profiling
```bash
# Generate a 30-second CPU profile
curl -o cpu.prof http://localhost:10000/debug/pprof/profile?seconds=30

# Analyze the profile
go tool pprof cpu.prof
```

#### Memory Profiling
```bash
# Get current memory profile
curl -o heap.prof http://localhost:10000/debug/pprof/heap

# Analyze the profile
go tool pprof heap.prof
```

#### Goroutine Profiling
```bash
# Get goroutine profile
curl -o goroutine.prof http://localhost:10000/debug/pprof/goroutine

# Analyze the profile
go tool pprof goroutine.prof
```

### 4. Use the Provided Scripts

#### Continuous Profiling
```bash
# Start continuous profiling (collects profiles every 5 minutes)
./scripts/pprof-continuous.sh

# Customize interval and duration
./scripts/pprof-continuous.sh -i 60 -d 15  # Every minute, 15s CPU profiles
```

#### Load Testing with Profiling
```bash
# Run load test with profiling
./scripts/pprof-load-test.sh

# Customize load test parameters
./scripts/pprof-load-test.sh -d 300 -c 20 -r 100  # 5min, 20 users, 100 RPS
```

#### Memory Leak Detection
```bash
# Detect memory leaks
./scripts/pprof-memory-leak.sh

# Trigger GC before analysis
./scripts/pprof-memory-leak.sh -g
```

#### Analyze Existing Profiles
```bash
# Analyze all profiles in current directory
./scripts/pprof-analyze.sh -a

# Analyze specific profile
./scripts/pprof-analyze.sh -t memory heap.prof

# Compare two profiles
./scripts/pprof-analyze.sh -c baseline.prof current.prof
```

## Interactive pprof Commands

Once you're in the pprof interactive mode, use these commands:

```bash
# Show top functions by CPU/memory usage
(pprof) top

# Show top 10 functions
(pprof) top10

# Show cumulative usage
(pprof) top -cum

# Show source code for specific function
(pprof) list <function_name>

# Show call graph
(pprof) web

# Generate PDF report
(pprof) pdf

# Generate SVG report
(pprof) svg

# Show allocation traces (for memory profiles)
(pprof) traces

# Exit pprof
(pprof) quit
```

## Common Analysis Workflows

### 1. CPU Bottleneck Analysis
```bash
# 1. Generate load
./scripts/generate-load.sh

# 2. Collect CPU profile
curl -o cpu.prof http://localhost:10000/debug/pprof/profile?seconds=30

# 3. Analyze
go tool pprof cpu.prof
(pprof) top
(pprof) list <hot_function>
```

### 2. Memory Leak Detection
```bash
# 1. Get baseline
curl -o baseline.prof http://localhost:10000/debug/pprof/heap

# 2. Generate load
./scripts/generate-load.sh

# 3. Get after-load profile
curl -o after.prof http://localhost:10000/debug/pprof/heap

# 4. Compare
go tool pprof -base baseline.prof after.prof
```

### 3. Goroutine Analysis
```bash
# 1. Get goroutine profile
curl -o goroutine.prof http://localhost:10000/debug/pprof/goroutine

# 2. Analyze
go tool pprof goroutine.prof
(pprof) top
(pprof) traces
```

## Monitoring Dashboard

Use the built-in monitoring tool:

```bash
# Start monitoring
./scripts/monitor-functions.sh
```

This provides real-time metrics and profiling data.

## Web Interface

For visual analysis, use the web interface:

```bash
# Start web interface for CPU profile
go tool pprof -http=:8080 cpu.prof

# Start web interface for memory profile
go tool pprof -http=:8080 heap.prof
```

Then open `http://localhost:8080` in your browser.

## Troubleshooting

### Profile Shows "No samples"
- Increase profiling duration: `?seconds=60`
- Generate more load during profiling
- Check if profiling is enabled in config

### Cannot connect to service
- Ensure application is running: `go run cmd/app/main.go`
- Check if port 10000 is available
- Verify config.yaml settings

### Scripts fail
- Install required tools: `jq`, `bc`
- Make scripts executable: `chmod +x scripts/*.sh`
- Check if Go is installed: `go version`

## Next Steps

1. **Read the full guide**: `docs/PPROF_GUIDE.md`
2. **Explore advanced features**: Use `-web` flag for visual analysis
3. **Set up continuous monitoring**: Use the continuous profiling script
4. **Integrate with CI/CD**: Add profiling to your build pipeline

## Quick Reference

| Command | Description |
|---------|-------------|
| `curl -o cpu.prof http://localhost:10000/debug/pprof/profile?seconds=30` | CPU profile |
| `curl -o heap.prof http://localhost:10000/debug/pprof/heap` | Memory profile |
| `curl -o goroutine.prof http://localhost:10000/debug/pprof/goroutine` | Goroutine profile |
| `go tool pprof cpu.prof` | Analyze CPU profile |
| `go tool pprof -web cpu.prof` | Web interface |
| `go tool pprof -base baseline.prof current.prof` | Compare profiles |
| `./scripts/pprof-continuous.sh` | Continuous profiling |
| `./scripts/pprof-load-test.sh` | Load testing |
| `./scripts/pprof-memory-leak.sh` | Memory leak detection |
| `./scripts/pprof-analyze.sh -a` | Analyze all profiles |
