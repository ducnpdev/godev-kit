# PostgreSQL Optimization Guide for Mac (8 CPU, 16GB RAM)

This guide provides optimized PostgreSQL configuration for Mac systems with 8 CPU cores and 16GB RAM, specifically tailored for the godev-kit application.

## 🎯 Optimization Goals

- **Performance**: Maximize query performance and throughput
- **Memory Efficiency**: Optimize memory usage for 16GB RAM
- **Concurrency**: Support multiple concurrent connections
- **Monitoring**: Enable comprehensive logging and monitoring
- **Stability**: Ensure reliable operation under load

## 📊 Hardware Specifications

- **CPU**: 8 cores
- **RAM**: 16GB
- **Storage**: SSD (assumed for I/O optimizations)
- **OS**: macOS

## 🔧 Configuration Files

### 1. Application Configuration (`config/postgres-optimized.yaml`)

```yaml
PG:
  # Connection Pool Settings (optimized for your hardware)
  POOL_MAX: 32                    # Increased from 10 - allows more concurrent connections
  POOL_MIN: 8                     # Increased from 2 - maintains more idle connections
  MAX_CONN_LIFETIME: 1h           # Increased from 30m - reduces connection churn
  MAX_CONN_IDLE_TIME: 15m         # Increased from 10m - keeps connections alive longer
  HEALTH_CHECK_PERIOD: 30s        # Decreased from 1m - more frequent health checks
  URL: "postgres://godevkit:1@localhost:5433/godevkit?sslmode=disable"
  
  # Additional PostgreSQL-specific optimizations
  STATEMENT_TIMEOUT: 30s          # Query timeout
  IDLE_IN_TRANSACTION_TIMEOUT: 10m # Idle transaction timeout
  LOCK_TIMEOUT: 5s                # Lock acquisition timeout
```

### 2. PostgreSQL Server Configuration (`config/postgresql.conf`)

#### Memory Settings
```conf
# Memory Configuration (optimized for 16GB RAM)
shared_buffers = 4GB                     # 25% of RAM (4GB)
effective_cache_size = 12GB              # 75% of RAM (12GB)
maintenance_work_mem = 1GB               # For maintenance operations
work_mem = 16MB                          # Per operation memory
wal_buffers = 16MB                       # WAL buffer size
```

#### Connection Settings
```conf
max_connections = 200                    # Increased for your hardware
max_worker_processes = 8                 # Match your CPU cores
max_parallel_workers_per_gather = 4      # 50% of CPU cores
max_parallel_workers = 8                 # Match your CPU cores
max_parallel_maintenance_workers = 2     # For maintenance operations
```

#### I/O Optimization (SSD)
```conf
random_page_cost = 1.1                   # SSD optimization
effective_io_concurrency = 200           # SSD optimization
checkpoint_completion_target = 0.9       # Spread checkpoint writes
```

## 🚀 Setup Instructions

### 1. Automatic Setup

Run the optimization setup script:

```bash
# Make the script executable
chmod +x scripts/setup-postgres-optimized.sh

# Run the setup
./scripts/setup-postgres-optimized.sh
```

### 2. Manual Setup

If you prefer manual setup:

#### Step 1: Backup Current Configuration
```bash
# Find your PostgreSQL data directory
pg_config --sysconfdir

# Backup current configuration
cp /path/to/postgresql.conf /path/to/postgresql.conf.backup.$(date +%Y%m%d_%H%M%S)
```

#### Step 2: Apply Optimized Configuration
```bash
# Copy the optimized configuration
cp config/postgresql.conf /path/to/postgresql.conf

# Set proper permissions
chmod 600 /path/to/postgresql.conf
```

#### Step 3: Create Log Directory
```bash
# Create log directory
mkdir -p /path/to/postgresql/log
chown $(whoami) /path/to/postgresql/log
```

#### Step 4: Restart PostgreSQL
```bash
# For Homebrew installations
brew services restart postgresql@15

# For system installations
sudo systemctl restart postgresql
```

## 📈 Performance Monitoring

### 1. Built-in Monitoring Script

Use the provided monitoring script:

```bash
./scripts/monitor-postgres.sh
```

This script provides:
- Connection information
- Active connections count
- Database size
- Cache hit ratio
- Slow query analysis

### 2. Manual Monitoring Queries

#### Check Current Settings
```sql
SELECT 
    name,
    setting,
    unit,
    context
FROM pg_settings 
WHERE name IN (
    'shared_buffers',
    'effective_cache_size',
    'work_mem',
    'maintenance_work_mem',
    'max_connections',
    'max_worker_processes'
);
```

#### Monitor Connection Usage
```sql
SELECT 
    count(*) as active_connections,
    count(*) * 100.0 / current_setting('max_connections')::int as usage_percent
FROM pg_stat_activity 
WHERE state = 'active';
```

#### Check Cache Performance
```sql
SELECT 
    schemaname,
    tablename,
    heap_blks_read,
    heap_blks_hit,
    round(100.0 * heap_blks_hit / (heap_blks_hit + heap_blks_read), 2) as cache_hit_ratio
FROM pg_statio_user_tables
WHERE heap_blks_hit + heap_blks_read > 0
ORDER BY cache_hit_ratio DESC;
```

#### Monitor Slow Queries
```sql
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC 
LIMIT 10;
```

## 🔍 Key Optimizations Explained

### 1. Memory Configuration

#### Shared Buffers (4GB)
- **Purpose**: PostgreSQL's main memory area for caching data
- **Rationale**: 25% of 16GB RAM provides optimal balance
- **Impact**: Reduces disk I/O for frequently accessed data

#### Effective Cache Size (12GB)
- **Purpose**: Tells the query planner about available system cache
- **Rationale**: 75% of RAM accounts for OS cache and other processes
- **Impact**: Better query planning and index usage decisions

#### Work Memory (16MB)
- **Purpose**: Memory per operation (sorts, joins, etc.)
- **Rationale**: Balanced for concurrent operations
- **Impact**: Prevents excessive memory usage per query

### 2. Connection Pool Optimization

#### Pool Size (32 max, 8 min)
- **Purpose**: Manage database connections efficiently
- **Rationale**: Supports concurrent requests without overwhelming the database
- **Impact**: Reduces connection overhead and improves response times

#### Connection Lifetime (1 hour)
- **Purpose**: How long connections stay alive
- **Rationale**: Reduces connection churn while preventing stale connections
- **Impact**: Balances resource usage and connection reuse

### 3. I/O Optimization

#### Random Page Cost (1.1)
- **Purpose**: Cost estimate for random page access
- **Rationale**: SSDs have much lower random access costs than HDDs
- **Impact**: Better query planning for SSD storage

#### Effective I/O Concurrency (200)
- **Purpose**: Number of concurrent I/O operations
- **Rationale**: SSDs can handle many concurrent operations
- **Impact**: Improved I/O throughput

### 4. Parallel Processing

#### Worker Processes (8)
- **Purpose**: Maximum parallel workers
- **Rationale**: Matches your 8 CPU cores
- **Impact**: Better utilization of multi-core processors

#### Parallel Workers per Gather (4)
- **Purpose**: Workers per parallel query
- **Rationale**: 50% of CPU cores for balanced load
- **Impact**: Optimal parallel query execution

## 🛠️ Troubleshooting

### Common Issues

#### 1. Memory Issues
**Symptoms**: Out of memory errors, slow performance
**Solutions**:
```bash
# Check memory usage
ps aux | grep postgres

# Reduce shared_buffers if needed
# shared_buffers = 2GB
```

#### 2. Connection Issues
**Symptoms**: "Too many connections" errors
**Solutions**:
```bash
# Check current connections
psql -c "SELECT count(*) FROM pg_stat_activity;"

# Increase max_connections if needed
# max_connections = 300
```

#### 3. Performance Issues
**Symptoms**: Slow queries, high CPU usage
**Solutions**:
```bash
# Enable query logging
# log_min_duration_statement = 1000

# Check for slow queries
psql -c "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### Performance Tuning

#### 1. Index Optimization
```sql
-- Find missing indexes
SELECT 
    schemaname,
    tablename,
    attname,
    n_distinct,
    correlation
FROM pg_stats 
WHERE schemaname = 'public'
ORDER BY n_distinct DESC;
```

#### 2. Table Statistics
```sql
-- Update table statistics
ANALYZE;

-- Check table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### 3. Vacuum Optimization
```sql
-- Check table bloat
SELECT 
    schemaname,
    tablename,
    n_dead_tup,
    n_live_tup,
    round(100.0 * n_dead_tup / (n_dead_tup + n_live_tup), 2) as dead_percent
FROM pg_stat_user_tables
WHERE n_dead_tup > 0
ORDER BY dead_percent DESC;
```

## 📊 Performance Benchmarks

### Expected Improvements

With these optimizations, you should see:

1. **Query Performance**: 20-40% improvement for complex queries
2. **Concurrent Users**: Support for 50-100 concurrent users
3. **Memory Usage**: Efficient use of available RAM
4. **I/O Performance**: Optimized for SSD storage
5. **Connection Handling**: Reduced connection overhead

### Monitoring Metrics

Track these key metrics:

- **Cache Hit Ratio**: Target > 95%
- **Connection Usage**: Keep below 80%
- **Query Response Time**: Monitor slow queries
- **Memory Usage**: Monitor shared_buffers efficiency
- **I/O Wait**: Should be minimal on SSD

## 🔄 Maintenance

### Regular Maintenance Tasks

1. **Weekly**:
   ```bash
   # Run monitoring script
   ./scripts/monitor-postgres.sh
   
   # Check for slow queries
   psql -c "SELECT * FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
   ```

2. **Monthly**:
   ```sql
   -- Update statistics
   ANALYZE;
   
   -- Check for table bloat
   SELECT schemaname, tablename, n_dead_tup FROM pg_stat_user_tables;
   ```

3. **Quarterly**:
   ```bash
   # Review and adjust configuration
   # Monitor performance trends
   # Consider hardware upgrades if needed
   ```

## 📚 Additional Resources

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [PostgreSQL Performance Tuning](https://www.postgresql.org/docs/current/runtime-config-query.html)
- [pg_stat_statements Extension](https://www.postgresql.org/docs/current/pgstatstatements.html)
- [PostgreSQL Monitoring](https://www.postgresql.org/docs/current/monitoring.html)

## 🤝 Support

If you encounter issues with this configuration:

1. Check the troubleshooting section above
2. Review PostgreSQL logs in the configured log directory
3. Use the monitoring script to identify bottlenecks
4. Consider reverting to the backup configuration if needed

Remember to test thoroughly in your specific environment before deploying to production! 