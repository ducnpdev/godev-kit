#!/bin/bash

# Load Testing Script for GoDev Kit - User Endpoint
# Usage: ./scripts/load-test.sh [test-type]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
USER_ENDPOINT="/v1/user"
TEST_DATA_FILE="test-data.json"
RESULTS_DIR="load-test-results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Create test data if not exists
create_test_data() {
    cat > $TEST_DATA_FILE << EOF
{
    "text": "Hello world, this is a test message for translation",
    "source": "en",
    "target": "vi"
}
EOF
    echo -e "${GREEN}✓ Test data created${NC}"
}

# Create results directory
mkdir -p $RESULTS_DIR   

# Health check
health_check() {
    echo -e "${BLUE}🔍 Checking application health...${NC}"
    if curl -s "$BASE_URL/health" > /dev/null; then
        echo -e "${GREEN}✓ Application is healthy${NC}"
    else
        echo -e "${RED}✗ Application is not responding${NC}"
        exit 1
    fi
}

# Warm up test
warm_up() {
    echo -e "${BLUE}🔥 Warming up application...${NC}"
    ab -n 1000 -c 10 -k "$BASE_URL$USER_ENDPOINT" > "$RESULTS_DIR/warmup_$TIMESTAMP.txt" 2>&1
    echo -e "${GREEN}✓ Warm up completed${NC}"
}

# 10k concurrency test
test_10k_concurrency() {
    echo -e "${BLUE}📊 Running 10k concurrency test...${NC}"
    echo "Test: 100,000 requests with 10,000 concurrent users"
    
    ab -n 100000 -c 10000 -k \
       -H "Content-Type: application/json" \
       "$BASE_URL$USER_ENDPOINT" > "$RESULTS_DIR/10k_concurrency_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ 10k concurrency test completed${NC}"
}

# 50k concurrency test
test_50k_concurrency() {
    echo -e "${BLUE}📊 Running 50k concurrency test...${NC}"
    echo "Test: 500,000 requests with 50,000 concurrent users"
    
    ab -n 500000 -c 50000 -k \
       -H "Content-Type: application/json" \
       "$BASE_URL$USER_ENDPOINT" > "$RESULTS_DIR/50k_concurrency_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ 50k concurrency test completed${NC}"
}

# 100k concurrency test
test_100k_concurrency() {
    echo -e "${BLUE}📊 Running 100k concurrency test...${NC}"
    echo "Test: 1,000,000 requests with 100,000 concurrent users"
    
    ab -n 1000000 -c 100000 -k \
       -H "Content-Type: application/json" \
       "$BASE_URL$USER_ENDPOINT" > "$RESULTS_DIR/100k_concurrency_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ 100k concurrency test completed${NC}"
}

# Basic load test (original)
basic_load_test() {
    echo -e "${BLUE}📊 Running basic load test...${NC}"
    echo "Test: 10,000 requests with 100 concurrent users"
    
    ab -n 10000 -c 100 -k \
       -H "Content-Type: application/json" \
       -p $TEST_DATA_FILE \
       "$BASE_URL/api/v1/translate" > "$RESULTS_DIR/basic_load_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ Basic load test completed${NC}"
}

# Stress test (original)
stress_test() {
    echo -e "${BLUE}💪 Running stress test...${NC}"
    echo "Test: 100,000 requests with 200 concurrent users"
    
    ab -n 100000 -c 200 -k \
       -H "Content-Type: application/json" \
       -p $TEST_DATA_FILE \
       "$BASE_URL/api/v1/translate" > "$RESULTS_DIR/stress_test_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ Stress test completed${NC}"
}

# Spike test (original)
spike_test() {
    echo -e "${BLUE}⚡ Running spike test...${NC}"
    echo "Test: 50,000 requests with 500 concurrent users"
    
    ab -n 50000 -c 500 -k \
       -H "Content-Type: application/json" \
       -p $TEST_DATA_FILE \
       "$BASE_URL/api/v1/translate" > "$RESULTS_DIR/spike_test_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ Spike test completed${NC}"
}

# Endurance test (original)
endurance_test() {
    echo -e "${BLUE}⏰ Running endurance test...${NC}"
    echo "Test: 1,000,000 requests with 100 concurrent users over time"
    
    ab -n 1000000 -c 100 -k \
       -H "Content-Type: application/json" \
       -p $TEST_DATA_FILE \
       "$BASE_URL/api/v1/translate" > "$RESULTS_DIR/endurance_test_$TIMESTAMP.txt" 2>&1
    
    echo -e "${GREEN}✓ Endurance test completed${NC}"
}

# Artillery test (if available)
artillery_test() {
    if command -v artillery &> /dev/null; then
        echo -e "${BLUE}🎯 Running Artillery test...${NC}"
        
        cat > artillery-config.yml << EOF
config:
  target: '$BASE_URL'
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
  - name: "User API Load Test"
    requests:
      - get:
          url: "$USER_ENDPOINT"
EOF
        
        artillery run artillery-config.yml > "$RESULTS_DIR/artillery_test_$TIMESTAMP.txt" 2>&1
        echo -e "${GREEN}✓ Artillery test completed${NC}"
    else
        echo -e "${YELLOW}⚠ Artillery not installed, skipping artillery test${NC}"
    fi
}

# Monitor system resources
monitor_resources() {
    echo -e "${BLUE}📈 Monitoring system resources...${NC}"
    
    # Start monitoring in background
    (
        while true; do
            echo "$(date '+%H:%M:%S') - $(ps aux | grep main | grep -v grep | awk '{print $3, $4}' | head -1)" >> "$RESULTS_DIR/cpu_memory_$TIMESTAMP.txt"
            echo "$(date '+%H:%M:%S') - Goroutines: $(curl -s $BASE_URL/debug/pprof/goroutine?debug=1 | grep -c 'goroutine')" >> "$RESULTS_DIR/goroutines_$TIMESTAMP.txt"
            sleep 5
        done
    ) &
    MONITOR_PID=$!
    
    # Store PID for cleanup
    echo $MONITOR_PID > "$RESULTS_DIR/monitor_pid.txt"
}

# Cleanup monitoring
cleanup_monitoring() {
    if [ -f "$RESULTS_DIR/monitor_pid.txt" ]; then
        MONITOR_PID=$(cat "$RESULTS_DIR/monitor_pid.txt")
        kill $MONITOR_PID 2>/dev/null || true
        rm -f "$RESULTS_DIR/monitor_pid.txt"
    fi
}

# Generate summary report
generate_report() {
    echo -e "${BLUE}📋 Generating summary report...${NC}"
    
    cat > "$RESULTS_DIR/summary_$TIMESTAMP.md" << EOF
# Load Test Results - $(date)

## Test Configuration
- Base URL: $BASE_URL
- Endpoint: $USER_ENDPOINT
- Timestamp: $TIMESTAMP
- Test Data: $TEST_DATA_FILE

## Results Summary

### 10k Concurrency Test (100k requests, 10k concurrent)
\`\`\`
$(tail -20 "$RESULTS_DIR/10k_concurrency_$TIMESTAMP.txt" | grep -E "(Requests per second|Time per request|Transfer rate|Failed requests)")
\`\`\`

### 50k Concurrency Test (500k requests, 50k concurrent)
\`\`\`
$(tail -20 "$RESULTS_DIR/50k_concurrency_$TIMESTAMP.txt" | grep -E "(Requests per second|Time per request|Transfer rate|Failed requests)")
\`\`\`

### 100k Concurrency Test (1M requests, 100k concurrent)
\`\`\`
$(tail -20 "$RESULTS_DIR/100k_concurrency_$TIMESTAMP.txt" | grep -E "(Requests per second|Time per request|Transfer rate|Failed requests)")
\`\`\`

## Performance Metrics
- CPU Usage: Check cpu_memory_$TIMESTAMP.txt
- Memory Usage: Check cpu_memory_$TIMESTAMP.txt  
- Goroutines: Check goroutines_$TIMESTAMP.txt

## Recommendations
1. Review failed requests
2. Monitor memory usage patterns
3. Check database connection pool
4. Analyze response time distribution
5. Consider connection limits and timeouts
EOF

    echo -e "${GREEN}✓ Summary report generated: $RESULTS_DIR/summary_$TIMESTAMP.md${NC}"
}

# Main execution
main() {
    echo -e "${GREEN}🚀 Starting Load Testing Suite for User Endpoint${NC}"
    echo "Timestamp: $TIMESTAMP"
    echo "Target Endpoint: $BASE_URL$USER_ENDPOINT"
    echo "Results will be saved to: $RESULTS_DIR"
    echo ""
    
    # Setup
    create_test_data
    health_check
    
    # Start monitoring
    monitor_resources
    
    # Run tests based on argument
    case "${1:-all}" in
        "10k")
            warm_up
            test_10k_concurrency
            ;;
        "50k")
            warm_up
            test_50k_concurrency
            ;;
        "100k")
            warm_up
            test_100k_concurrency
            ;;
        "user-all")
            warm_up
            test_10k_concurrency
            test_50k_concurrency
            test_100k_concurrency
            ;;
        "basic")
            warm_up
            basic_load_test
            ;;
        "stress")
            warm_up
            stress_test
            ;;
        "spike")
            warm_up
            spike_test
            ;;
        "endurance")
            warm_up
            endurance_test
            ;;
        "artillery")
            artillery_test
            ;;
        "all")
            warm_up
            test_10k_concurrency
            test_50k_concurrency
            test_100k_concurrency
            basic_load_test
            stress_test
            spike_test
            endurance_test
            artillery_test
            ;;
        *)
            echo -e "${RED}Invalid test type. Use: 10k, 50k, 100k, user-all, basic, stress, spike, endurance, artillery, or all${NC}"
            echo -e "${YELLOW}Available options:${NC}"
            echo -e "  ${GREEN}10k${NC}     - 10,000 concurrent users test"
            echo -e "  ${GREEN}50k${NC}     - 50,000 concurrent users test"
            echo -e "  ${GREEN}100k${NC}    - 100,000 concurrent users test"
            echo -e "  ${GREEN}user-all${NC} - All user endpoint tests (10k, 50k, 100k)"
            echo -e "  ${GREEN}all${NC}     - All tests including original ones"
            exit 1
            ;;
    esac
    
    # Cleanup and report
    cleanup_monitoring
    generate_report
    
    echo -e "${GREEN}🎉 Load testing completed!${NC}"
    echo -e "${BLUE}📁 Check results in: $RESULTS_DIR${NC}"
}

# Trap to cleanup on exit
trap cleanup_monitoring EXIT

# Run main function
main "$@" 