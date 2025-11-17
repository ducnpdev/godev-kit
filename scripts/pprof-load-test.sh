#!/bin/bash

# Load Testing with pprof Profiling Script
# This script generates load while collecting pprof profiles to identify performance bottlenecks

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
PROFILES_DIR="load-test-profiles"
CPU_DURATION=60
LOAD_DURATION=120
CONCURRENT_USERS=10
REQUESTS_PER_SECOND=50

# Check if service is running
check_service() {
    if ! curl -s "$BASE_URL/healthz" > /dev/null; then
        echo -e "${RED}❌ Service is not running at $BASE_URL${NC}"
        echo "Please start your service first:"
        echo "  go run cmd/app/main.go"
        exit 1
    fi
    echo -e "${GREEN}✅ Service is running${NC}"
}

# Create profiles directory
setup_directories() {
    mkdir -p "$PROFILES_DIR"
    echo -e "${BLUE}📁 Load test profiles will be saved to: $PROFILES_DIR${NC}"
}

# Get baseline metrics
get_baseline_metrics() {
    echo -e "${BLUE}📊 Getting baseline metrics...${NC}"
    
    if stats=$(curl -s "$BASE_URL/debug/stats" 2>/dev/null); then
        baseline_goroutines=$(echo "$stats" | jq -r '.goroutines // 0')
        baseline_heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
        baseline_heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
        
        echo -e "Baseline Goroutines: ${YELLOW}$baseline_goroutines${NC}"
        echo -e "Baseline Heap Allocated: ${YELLOW}$(numfmt --to=iec $baseline_heap_alloc)${NC}"
        echo -e "Baseline Heap System: ${YELLOW}$(numfmt --to=iec $baseline_heap_sys)${NC}"
    else
        echo -e "${RED}Failed to get baseline metrics${NC}"
    fi
}

# Collect baseline profiles
collect_baseline_profiles() {
    echo -e "${BLUE}📊 Collecting baseline profiles...${NC}"
    
    # CPU profile (shorter duration for baseline)
    local timestamp=$(date +%Y%m%d_%H%M%S)
    curl -s -o "${PROFILES_DIR}/baseline_cpu_${timestamp}.prof" "$BASE_URL/debug/pprof/profile?seconds=10" &
    local cpu_pid=$!
    
    # Memory profile
    curl -s -o "${PROFILES_DIR}/baseline_heap_${timestamp}.prof" "$BASE_URL/debug/pprof/heap" &
    local heap_pid=$!
    
    # Goroutine profile
    curl -s -o "${PROFILES_DIR}/baseline_goroutine_${timestamp}.prof" "$BASE_URL/debug/pprof/goroutine" &
    local goroutine_pid=$!
    
    # Wait for all profiles to complete
    wait $cpu_pid $heap_pid $goroutine_pid
    
    echo -e "${GREEN}✅ Baseline profiles collected${NC}"
}

# Generate load
generate_load() {
    echo -e "${BLUE}🚀 Starting load generation...${NC}"
    echo -e "Duration: ${YELLOW}${LOAD_DURATION}s${NC}"
    echo -e "Concurrent Users: ${YELLOW}${CONCURRENT_USERS}${NC}"
    echo -e "Requests/Second: ${YELLOW}${REQUESTS_PER_SECOND}${NC}"
    
    # Calculate delay between requests
    local delay=$(echo "scale=3; 1 / $REQUESTS_PER_SECOND" | bc)
    
    # Start load generation in background
    (
        local start_time=$(date +%s)
        local end_time=$((start_time + LOAD_DURATION))
        
        while [ $(date +%s) -lt $end_time ]; do
            # Start concurrent requests
            for i in $(seq 1 $CONCURRENT_USERS); do
                # Random endpoint selection
                case $((RANDOM % 5)) in
                    0) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                    1) curl -s "$BASE_URL/healthz" > /dev/null & ;;
                    2) curl -s "$BASE_URL/metrics" > /dev/null & ;;
                    3) curl -s "$BASE_URL/debug/stats" > /dev/null & ;;
                    4) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                esac
            done
            
            # Wait for requests to complete
            wait
            
            # Delay between request batches
            sleep "$delay"
        done
    ) &
    
    local load_pid=$!
    echo -e "${GREEN}✅ Load generation started (PID: $load_pid)${NC}"
    return $load_pid
}

# Monitor metrics during load
monitor_metrics() {
    local load_pid=$1
    local interval=5
    local count=0
    
    echo -e "${BLUE}📈 Monitoring metrics during load...${NC}"
    
    while kill -0 $load_pid 2>/dev/null; do
        if stats=$(curl -s "$BASE_URL/debug/stats" 2>/dev/null); then
            goroutines=$(echo "$stats" | jq -r '.goroutines // 0')
            heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
            heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
            
            echo -e "[$((count * interval))s] Goroutines: ${YELLOW}$goroutines${NC}, Heap: ${YELLOW}$(numfmt --to=iec $heap_alloc)${NC}"
        fi
        
        sleep $interval
        ((count++))
    done
}

# Collect profiles during load
collect_load_profiles() {
    local load_pid=$1
    
    echo -e "${BLUE}📊 Collecting profiles during load...${NC}"
    
    # Wait a bit for load to stabilize
    sleep 10
    
    # Start CPU profiling
    local timestamp=$(date +%Y%m%d_%H%M%S)
    echo -e "${BLUE}📊 Starting CPU profile (${CPU_DURATION}s)...${NC}"
    
    curl -s -o "${PROFILES_DIR}/load_cpu_${timestamp}.prof" "$BASE_URL/debug/pprof/profile?seconds=${CPU_DURATION}" &
    local cpu_pid=$!
    
    # Wait for CPU profile to complete
    wait $cpu_pid
    
    echo -e "${GREEN}✅ CPU profile completed${NC}"
    
    # Collect memory and goroutine profiles
    curl -s -o "${PROFILES_DIR}/load_heap_${timestamp}.prof" "$BASE_URL/debug/pprof/heap" &
    curl -s -o "${PROFILES_DIR}/load_goroutine_${timestamp}.prof" "$BASE_URL/debug/pprof/goroutine" &
    
    wait
    
    echo -e "${GREEN}✅ All load profiles collected${NC}"
}

# Collect final profiles
collect_final_profiles() {
    echo -e "${BLUE}📊 Collecting final profiles...${NC}"
    
    # Wait for system to stabilize
    sleep 10
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    
    # Collect final profiles
    curl -s -o "${PROFILES_DIR}/final_heap_${timestamp}.prof" "$BASE_URL/debug/pprof/heap" &
    curl -s -o "${PROFILES_DIR}/final_goroutine_${timestamp}.prof" "$BASE_URL/debug/pprof/goroutine" &
    
    wait
    
    echo -e "${GREEN}✅ Final profiles collected${NC}"
}

# Get final metrics
get_final_metrics() {
    echo -e "${BLUE}📊 Getting final metrics...${NC}"
    
    if stats=$(curl -s "$BASE_URL/debug/stats" 2>/dev/null); then
        final_goroutines=$(echo "$stats" | jq -r '.goroutines // 0')
        final_heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
        final_heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
        
        echo -e "Final Goroutines: ${YELLOW}$final_goroutines${NC}"
        echo -e "Final Heap Allocated: ${YELLOW}$(numfmt --to=iec $final_heap_alloc)${NC}"
        echo -e "Final Heap System: ${YELLOW}$(numfmt --to=iec $final_heap_sys)${NC}"
        
        # Calculate differences
        if [ -n "$baseline_goroutines" ] && [ -n "$final_goroutines" ]; then
            local goroutine_diff=$((final_goroutines - baseline_goroutines))
            local heap_diff=$((final_heap_alloc - baseline_heap_alloc))
            
            echo -e "Goroutine Change: ${YELLOW}$goroutine_diff${NC}"
            echo -e "Heap Change: ${YELLOW}$(numfmt --to=iec $heap_diff)${NC}"
        fi
    else
        echo -e "${RED}Failed to get final metrics${NC}"
    fi
}

# Generate analysis report
generate_report() {
    echo -e "${BLUE}📋 Generating analysis report...${NC}"
    
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local report_file="${PROFILES_DIR}/load_test_report_${timestamp}.txt"
    
    {
        echo "Load Test Report"
        echo "================"
        echo "Date: $(date)"
        echo "Duration: ${LOAD_DURATION}s"
        echo "Concurrent Users: ${CONCURRENT_USERS}"
        echo "Requests/Second: ${REQUESTS_PER_SECOND}"
        echo ""
        echo "Baseline Metrics:"
        echo "  Goroutines: $baseline_goroutines"
        echo "  Heap Allocated: $(numfmt --to=iec $baseline_heap_alloc)"
        echo "  Heap System: $(numfmt --to=iec $baseline_heap_sys)"
        echo ""
        echo "Final Metrics:"
        echo "  Goroutines: $final_goroutines"
        echo "  Heap Allocated: $(numfmt --to=iec $final_heap_alloc)"
        echo "  Heap System: $(numfmt --to=iec $final_heap_sys)"
        echo ""
        echo "Profiles Generated:"
        ls -la "${PROFILES_DIR}"/*.prof 2>/dev/null || echo "No profiles found"
        echo ""
        echo "Analysis Commands:"
        echo "  # CPU analysis:"
        echo "  go tool pprof ${PROFILES_DIR}/load_cpu_*.prof"
        echo ""
        echo "  # Memory analysis:"
        echo "  go tool pprof ${PROFILES_DIR}/load_heap_*.prof"
        echo ""
        echo "  # Goroutine analysis:"
        echo "  go tool pprof ${PROFILES_DIR}/load_goroutine_*.prof"
        echo ""
        echo "  # Compare baseline vs load:"
        echo "  go tool pprof -base ${PROFILES_DIR}/baseline_heap_*.prof ${PROFILES_DIR}/load_heap_*.prof"
    } > "$report_file"
    
    echo -e "${GREEN}✅ Report saved: $report_file${NC}"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -d, --duration SECONDS    Load test duration (default: 120)"
    echo "  -c, --concurrent USERS    Concurrent users (default: 10)"
    echo "  -r, --rate RPS            Requests per second (default: 50)"
    echo "  -p, --cpu-duration SEC    CPU profile duration (default: 60)"
    echo "  -o, --output DIR          Output directory (default: load-test-profiles)"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                        # Run with default settings"
    echo "  $0 -d 300 -c 20 -r 100   # 5min test, 20 users, 100 RPS"
    echo "  $0 -p 30                 # CPU profile for 30 seconds"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--duration)
                LOAD_DURATION="$2"
                shift 2
                ;;
            -c|--concurrent)
                CONCURRENT_USERS="$2"
                shift 2
                ;;
            -r|--rate)
                REQUESTS_PER_SECOND="$2"
                shift 2
                ;;
            -p|--cpu-duration)
                CPU_DURATION="$2"
                shift 2
                ;;
            -o|--output)
                PROFILES_DIR="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                echo -e "${RED}Unknown option: $1${NC}"
                show_usage
                exit 1
                ;;
        esac
    done
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    echo -e "${GREEN}Load Testing with pprof Profiling${NC}"
    echo "====================================="
    
    # Check prerequisites
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq is required but not installed. Please install jq first.${NC}"
        exit 1
    fi
    
    if ! command -v bc &> /dev/null; then
        echo -e "${RED}❌ bc is required but not installed. Please install bc first.${NC}"
        exit 1
    fi
    
    # Check service and setup
    check_service
    setup_directories
    
    echo -e "${BLUE}🚀 Starting load test with profiling...${NC}"
    echo ""
    
    # Phase 1: Baseline
    echo -e "${YELLOW}Phase 1: Baseline Collection${NC}"
    echo "================================"
    get_baseline_metrics
    collect_baseline_profiles
    echo ""
    
    # Phase 2: Load Generation
    echo -e "${YELLOW}Phase 2: Load Generation${NC}"
    echo "================================"
    generate_load
    local load_pid=$!
    
    # Phase 3: Monitoring and Profiling
    echo -e "${YELLOW}Phase 3: Monitoring and Profiling${NC}"
    echo "=========================================="
    monitor_metrics $load_pid &
    local monitor_pid=$!
    
    collect_load_profiles $load_pid
    
    # Wait for load to complete
    wait $load_pid
    kill $monitor_pid 2>/dev/null || true
    
    echo ""
    
    # Phase 4: Final Analysis
    echo -e "${YELLOW}Phase 4: Final Analysis${NC}"
    echo "================================"
    get_final_metrics
    collect_final_profiles
    generate_report
    
    echo ""
    echo -e "${GREEN}✅ Load test completed!${NC}"
    echo -e "${BLUE}📁 Profiles saved in: $PROFILES_DIR${NC}"
    echo -e "${BLUE}📋 Report generated: ${PROFILES_DIR}/load_test_report_*.txt${NC}"
}

# Run main function
main "$@"
