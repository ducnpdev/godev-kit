# Performance Tuning Guide - Handle 100k Requests

## Tổng quan
Hướng dẫn này sẽ giúp bạn tuning ứng dụng Go để handle 100k requests với hiệu suất cao và ổn định.

## 1. HTTP Server Tuning

### 1.1 Cấu hình Fiber Server
```yaml
HTTP:
  PORT: 10000
  READ_TIMEOUT: 5s        # Giảm từ 10s xuống 5s
  WRITE_TIMEOUT: 10s      # Giảm từ 30s xuống 10s  
  IDLE_TIMEOUT: 60s       # Tăng từ 30s lên 60s
  SHUTDOWN_TIMEOUT: 10s   # Tăng từ 5s lên 10s
  USE_PREFORK_MODE: true  # Bật prefork mode
  API_TIMEOUT: 3s         # Giảm từ 5s xuống 3s
```

### 1.2 Tối ưu hóa Fiber Config
```go
app := fiber.New(fiber.Config{
    Prefork:               true,           // Sử dụng multiple processes
    ReadTimeout:           5 * time.Second,
    WriteTimeout:          10 * time.Second,
    IdleTimeout:           60 * time.Second,
    DisableStartupMessage: true,           // Tắt startup message
    EnableTrustedProxyCheck: false,        // Tắt proxy check nếu không cần
    ProxyHeader:           "X-Forwarded-For",
    GetOnly:               false,
    ErrorHandler:          customErrorHandler,
    JSONEncoder:           json.Marshal,
    JSONDecoder:           json.Unmarshal,
    // Tối ưu memory
    BodyLimit:             10 * 1024 * 1024, // 10MB limit
    Concurrency:           256 * 1024,        // Tăng concurrency
})
```

## 2. Database Connection Pool Tuning

### 2.1 PostgreSQL Connection Pool
```yaml
PG:
  POOL_MAX: 50              # Tăng từ 10 lên 50
  POOL_MIN: 10              # Tăng từ 2 lên 10
  MAX_CONN_LIFETIME: 15m    # Giảm từ 30m xuống 15m
  MAX_CONN_IDLE_TIME: 5m    # Giảm từ 10m xuống 5m
  HEALTH_CHECK_PERIOD: 30s  # Giảm từ 1m xuống 30s
```

### 2.2 PostgreSQL Server Tuning
```sql
-- postgresql.conf optimizations
max_connections = 200
shared_buffers = 256MB
effective_cache_size = 1GB
work_mem = 4MB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
```

## 3. Redis Tuning

### 3.1 Redis Configuration
```yaml
REDIS:
  URL: "redis://localhost:6379/0"
  POOL_SIZE: 50
  MIN_IDLE_CONNS: 10
  MAX_RETRIES: 3
  DIAL_TIMEOUT: 5s
  READ_TIMEOUT: 3s
  WRITE_TIMEOUT: 3s
```

### 3.2 Redis Server Tuning
```conf
# redis.conf
maxmemory 512mb
maxmemory-policy allkeys-lru
save 900 1
save 300 10
save 60 10000
tcp-keepalive 300
```

## 4. Application Level Optimizations

### 4.1 Goroutine Pool
```go
// internal/pkg/workerpool/workerpool.go
type WorkerPool struct {
    workers    int
    jobQueue   chan Job
    workerPool chan chan Job
    quit       chan bool
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:    workers,
        jobQueue:   make(chan Job, 1000),
        workerPool: make(chan chan Job, workers),
        quit:       make(chan bool),
    }
}
```

### 4.2 Connection Pooling
```go
// Tối ưu connection pooling
var (
    dbPool *sql.DB
    redisPool *redis.Client
)

func initPools() {
    // Database pool
    dbPool, _ = sql.Open("postgres", dbURL)
    dbPool.SetMaxOpenConns(50)
    dbPool.SetMaxIdleConns(10)
    dbPool.SetConnMaxLifetime(15 * time.Minute)
    
    // Redis pool
    redisPool = redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        PoolSize:     50,
        MinIdleConns: 10,
        MaxRetries:   3,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })
}
```

## 5. Memory Management

### 5.1 GC Tuning
```bash
# Environment variables for GC tuning
export GOGC=100        # Trigger GC when heap grows 100%
export GOMEMLIMIT=512MiB  # Memory limit
export GOMAXPROCS=8    # Number of CPU cores
```

### 5.2 Memory Pool
```go
// internal/pkg/mempool/mempool.go
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, 1024)
            },
        },
    }
}

func (bp *BufferPool) Get() []byte {
    return bp.pool.Get().([]byte)
}

func (bp *BufferPool) Put(buf []byte) {
    buf = buf[:0] // Reset slice
    bp.pool.Put(buf)
}
```

## 6. Caching Strategy

### 6.1 Multi-level Caching
```go
type CacheManager struct {
    l1Cache *redis.Client  // Redis - Fast access
    l2Cache *bigcache.BigCache // In-memory - Very fast
}

func (cm *CacheManager) Get(key string) ([]byte, error) {
    // Try L2 cache first (in-memory)
    if data, err := cm.l2Cache.Get(key); err == nil {
        return data, nil
    }
    
    // Try L1 cache (Redis)
    if data, err := cm.l1Cache.Get(ctx, key).Bytes(); err == nil {
        // Store in L2 cache
        cm.l2Cache.Set(key, data)
        return data, nil
    }
    
    return nil, errors.New("not found")
}
```

## 7. Load Testing Scripts

### 7.1 Apache Bench (ab)
```bash
#!/bin/bash
# scripts/load-test.sh

echo "Starting load test..."

# Warm up
ab -n 1000 -c 10 http://localhost:10000/health

# Main test
ab -n 100000 -c 100 -k http://localhost:10000/api/v1/translate \
  -H "Content-Type: application/json" \
  -p test-data.json

# Stress test
ab -n 1000000 -c 200 -k http://localhost:10000/api/v1/translate \
  -H "Content-Type: application/json" \
  -p test-data.json
```

### 7.2 Artillery
```yaml
# scripts/artillery-config.yml
config:
  target: 'http://localhost:10000'
  phases:
    - duration: 60
      arrivalRate: 100
    - duration: 300
      arrivalRate: 500
    - duration: 60
      arrivalRate: 1000
  defaults:
    headers:
      Content-Type: 'application/json'

scenarios:
  - name: "API Load Test"
    requests:
      - post:
          url: "/api/v1/translate"
          json:
            text: "Hello world"
            source: "en"
            target: "vi"
```

## 8. Monitoring & Profiling

### 8.1 Prometheus Metrics
```go
// internal/pkg/metrics/metrics.go
var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)
```

### 8.2 Health Check Endpoint
```go
func healthCheck(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "status": "healthy",
        "timestamp": time.Now(),
        "uptime": time.Since(startTime).String(),
        "goroutines": runtime.NumGoroutine(),
        "memory": getMemoryStats(),
    })
}
```

## 9. Deployment Optimizations

### 9.1 Docker Optimization
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 10000
CMD ["./main"]
```

### 9.2 Kubernetes Resources
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: godev-kit
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: godev-kit
        image: godev-kit:latest
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        env:
        - name: GOGC
          value: "100"
        - name: GOMEMLIMIT
          value: "512MiB"
        - name: GOMAXPROCS
          value: "8"
```

## 10. Performance Checklist

### 10.1 Pre-deployment
- [ ] Enable prefork mode
- [ ] Optimize connection pools
- [ ] Configure proper timeouts
- [ ] Set up monitoring
- [ ] Implement caching strategy
- [ ] Tune GC parameters

### 10.2 During Testing
- [ ] Monitor memory usage
- [ ] Check goroutine count
- [ ] Monitor database connections
- [ ] Track response times
- [ ] Monitor error rates
- [ ] Check CPU usage

### 10.3 Post-deployment
- [ ] Set up alerts
- [ ] Monitor logs
- [ ] Track performance metrics
- [ ] Optimize based on real usage
- [ ] Plan for scaling

## 11. Expected Performance Metrics

Với các tối ưu hóa trên, bạn có thể đạt được:

- **Throughput**: 10,000-15,000 requests/second
- **Latency**: < 50ms (95th percentile)
- **Memory Usage**: < 512MB
- **Goroutines**: < 1000
- **Database Connections**: < 50
- **Error Rate**: < 0.1%

## 12. Troubleshooting

### 12.1 High Memory Usage
```bash
# Check memory usage
go tool pprof http://localhost:10000/debug/pprof/heap

# Check goroutines
go tool pprof http://localhost:10000/debug/pprof/goroutine
```

### 12.2 High CPU Usage
```bash
# CPU profiling
go tool pprof http://localhost:10000/debug/pprof/profile
```

### 12.3 Database Bottlenecks
```sql
-- Check slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC 
LIMIT 10;
```

## Kết luận

Để handle 100k requests hiệu quả, cần tập trung vào:
1. **HTTP Server tuning** với prefork mode
2. **Connection pooling** cho database và Redis
3. **Memory management** và GC tuning
4. **Caching strategy** multi-level
5. **Monitoring** và alerting
6. **Load testing** trước deployment

Áp dụng các tối ưu hóa này sẽ giúp ứng dụng của bạn handle được 100k requests một cách ổn định và hiệu quả. 