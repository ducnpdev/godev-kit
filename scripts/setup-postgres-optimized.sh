#!/bin/bash

# PostgreSQL Optimization Setup Script for Mac (8 CPU, 16GB RAM)
# This script helps configure PostgreSQL for optimal performance

set -e

echo "🚀 PostgreSQL Optimization Setup for Mac (8 CPU, 16GB RAM)"
echo "=========================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[SETUP]${NC} $1"
}

# Check if PostgreSQL is installed
check_postgres() {
    print_header "Checking PostgreSQL installation..."
    
    if ! command -v psql &> /dev/null; then
        print_error "PostgreSQL is not installed. Please install it first:"
        echo "  brew install postgresql@15"
        echo "  brew services start postgresql@15"
        exit 1
    fi
    
    POSTGRES_VERSION=$(psql --version | awk '{print $3}' | cut -d. -f1,2)
    print_status "PostgreSQL version: $POSTGRES_VERSION"
}

# Get PostgreSQL data directory
get_postgres_data_dir() {
    print_header "Finding PostgreSQL data directory..."
    
    # Try different methods to find the data directory
    if command -v pg_config &> /dev/null; then
        DATA_DIR=$(pg_config --sysconfdir)
        if [ -d "$DATA_DIR" ]; then
            print_status "PostgreSQL config directory: $DATA_DIR"
        fi
    fi
    
    # For Homebrew installations
    if [ -d "/opt/homebrew/var/postgresql@15" ]; then
        DATA_DIR="/opt/homebrew/var/postgresql@15"
    elif [ -d "/usr/local/var/postgresql@15" ]; then
        DATA_DIR="/usr/local/var/postgresql@15"
    elif [ -d "/opt/homebrew/var/postgresql" ]; then
        DATA_DIR="/opt/homebrew/var/postgresql"
    elif [ -d "/usr/local/var/postgresql" ]; then
        DATA_DIR="/usr/local/var/postgresql"
    else
        print_warning "Could not automatically detect PostgreSQL data directory"
        print_warning "Please manually specify the path to your PostgreSQL data directory"
        read -p "Enter PostgreSQL data directory path: " DATA_DIR
    fi
    
    print_status "Using data directory: $DATA_DIR"
}

# Backup current configuration
backup_config() {
    print_header "Backing up current PostgreSQL configuration..."
    
    if [ -f "$DATA_DIR/postgresql.conf" ]; then
        BACKUP_FILE="$DATA_DIR/postgresql.conf.backup.$(date +%Y%m%d_%H%M%S)"
        cp "$DATA_DIR/postgresql.conf" "$BACKUP_FILE"
        print_status "Configuration backed up to: $BACKUP_FILE"
    else
        print_warning "No existing postgresql.conf found"
    fi
}

# Apply optimized configuration
apply_config() {
    print_header "Applying optimized PostgreSQL configuration..."
    
    # Copy the optimized configuration
    if [ -f "config/postgresql.conf" ]; then
        cp "config/postgresql.conf" "$DATA_DIR/postgresql.conf"
        print_status "Optimized configuration applied"
    else
        print_error "Optimized configuration file not found at config/postgresql.conf"
        exit 1
    fi
    
    # Set proper permissions
    chmod 600 "$DATA_DIR/postgresql.conf"
    print_status "Configuration file permissions set"
}

# Create log directory
setup_logging() {
    print_header "Setting up logging directory..."
    
    LOG_DIR="$DATA_DIR/log"
    mkdir -p "$LOG_DIR"
    chown $(whoami) "$LOG_DIR"
    print_status "Log directory created: $LOG_DIR"
}

# Restart PostgreSQL service
restart_postgres() {
    print_header "Restarting PostgreSQL service..."
    
    if command -v brew &> /dev/null; then
        print_status "Restarting PostgreSQL via Homebrew..."
        brew services restart postgresql@15 2>/dev/null || brew services restart postgresql 2>/dev/null
    else
        print_warning "Please restart PostgreSQL manually:"
        echo "  sudo systemctl restart postgresql"
        echo "  or"
        echo "  brew services restart postgresql@15"
    fi
    
    # Wait for PostgreSQL to start
    print_status "Waiting for PostgreSQL to start..."
    sleep 5
    
    # Test connection
    if pg_isready -h localhost -p 5433 >/dev/null 2>&1; then
        print_status "PostgreSQL is running and accepting connections"
    else
        print_warning "PostgreSQL may not be running. Please check manually:"
        echo "  pg_isready -h localhost -p 5433"
    fi
}

# Create monitoring script
create_monitoring_script() {
    print_header "Creating PostgreSQL monitoring script..."
    
    cat > scripts/monitor-postgres.sh << 'EOF'
#!/bin/bash

# PostgreSQL Monitoring Script
# Usage: ./scripts/monitor-postgres.sh

echo "📊 PostgreSQL Performance Monitoring"
echo "==================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if PostgreSQL is running
if ! pg_isready -h localhost -p 5433 >/dev/null 2>&1; then
    echo -e "${RED}❌ PostgreSQL is not running${NC}"
    exit 1
fi

echo -e "${GREEN}✅ PostgreSQL is running${NC}"

# Connection info
echo -e "\n${YELLOW}Connection Information:${NC}"
psql -h localhost -p 5433 -U godevkit -d godevkit -c "
SELECT 
    version() as postgres_version,
    current_setting('max_connections') as max_connections,
    current_setting('shared_buffers') as shared_buffers,
    current_setting('effective_cache_size') as effective_cache_size,
    current_setting('work_mem') as work_mem,
    current_setting('maintenance_work_mem') as maintenance_work_mem;
" 2>/dev/null || echo "Could not connect to database"

# Active connections
echo -e "\n${YELLOW}Active Connections:${NC}"
psql -h localhost -p 5433 -U godevkit -d godevkit -c "
SELECT 
    count(*) as active_connections,
    count(*) * 100.0 / current_setting('max_connections')::int as connection_usage_percent
FROM pg_stat_activity 
WHERE state = 'active';
" 2>/dev/null || echo "Could not get connection info"

# Database size
echo -e "\n${YELLOW}Database Size:${NC}"
psql -h localhost -p 5433 -U godevkit -d godevkit -c "
SELECT 
    pg_size_pretty(pg_database_size(current_database())) as database_size;
" 2>/dev/null || echo "Could not get database size"

# Cache hit ratio
echo -e "\n${YELLOW}Cache Hit Ratio:${NC}"
psql -h localhost -p 5433 -U godevkit -d godevkit -c "
SELECT 
    round(100.0 * sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)), 2) as cache_hit_ratio
FROM pg_statio_user_tables;
" 2>/dev/null || echo "Could not get cache hit ratio"

# Slow queries (if any)
echo -e "\n${YELLOW}Recent Slow Queries (>1s):${NC}"
psql -h localhost -p 5433 -U godevkit -d godevkit -c "
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    rows
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC 
LIMIT 5;
" 2>/dev/null || echo "pg_stat_statements extension not available"

echo -e "\n${GREEN}Monitoring complete${NC}"
EOF

    chmod +x scripts/monitor-postgres.sh
    print_status "Monitoring script created: scripts/monitor-postgres.sh"
}

# Main setup function
main() {
    check_postgres
    get_postgres_data_dir
    backup_config
    apply_config
    setup_logging
    restart_postgres
    create_monitoring_script
    
    echo -e "\n${GREEN}✅ PostgreSQL optimization setup complete!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}"
    echo "1. Test your application with the new configuration"
    echo "2. Monitor performance using: ./scripts/monitor-postgres.sh"
    echo "3. Check logs in: $DATA_DIR/log/"
    echo "4. If issues occur, restore backup: cp $BACKUP_FILE $DATA_DIR/postgresql.conf"
    
    echo -e "\n${BLUE}Key optimizations applied:${NC}"
    echo "• Increased connection pool (32 max, 8 min)"
    echo "• Optimized memory settings (4GB shared_buffers, 12GB effective_cache)"
    echo "• SSD-optimized I/O settings"
    echo "• Enhanced logging and monitoring"
    echo "• Improved autovacuum settings"
}

# Run main function
main "$@" 