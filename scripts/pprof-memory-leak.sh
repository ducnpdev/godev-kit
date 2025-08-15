#!/bin/bash

# Memory Leak Detection with pprof Script
# This script helps detect memory leaks by collecting and comparing memory profiles

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
PROFILES_DIR="memory-leak-profiles"
TEST_DURATION=300  # 5 minutes
PROFILE_INTERVAL=30  # Profile every 30 seconds
LOAD_DURATION=60   # Generate load for 60 seconds

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
    echo -e "${BLUE}📁 Memory leak profiles will be saved to: $PROFILES_DIR${NC}"
}

# Get memory statistics
get_memory_stats() {
    if stats=$(curl -s "$BASE_URL/debug/stats" 2>/dev/null); then
        heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
        heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
        heap_idle=$(echo "$stats" | jq -r '.memory.idle // 0')
        heap_inuse=$(echo "$stats" | jq -r '.memory.inuse // 0')
        heap_objects=$(echo "$stats" | jq -r '.memory.objects // 0')
        
        echo "$heap_alloc $heap_sys $heap_idle $heap_inuse $heap_objects"
    else
        echo "0 0 0 0 0"
    fi
}

# Collect memory profile
collect_memory_profile() {
    local timestamp=$1
    local filename="${PROFILES_DIR}/heap_${timestamp}.prof"
    
    if curl -s -o "$filename" "$BASE_URL/debug/pprof/heap"; then
        local size=$(du -h "$filename" | cut -f1)
        echo -e "${GREEN}✅ Memory profile saved: $filename (${size})${NC}"
        return 0
    else
        echo -e "${RED}❌ Failed to collect memory profile${NC}"
        return 1
    fi
}

# Generate load to trigger potential memory leaks
generate_load() {
    echo -e "${BLUE}🚀 Generating load to trigger memory leaks...${NC}"
    echo -e "Duration: ${YELLOW}${LOAD_DURATION}s${NC}"
    
    # Start load generation in background
    (
        local start_time=$(date +%s)
        local end_time=$((start_time + LOAD_DURATION))
        
        while [ $(date +%s) -lt $end_time ]; do
            # Make concurrent requests to different endpoints
            for i in {1..20}; do
                case $((RANDOM % 6)) in
                    0) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                    1) curl -s "$BASE_URL/healthz" > /dev/null & ;;
                    2) curl -s "$BASE_URL/metrics" > /dev/null & ;;
                    3) curl -s "$BASE_URL/debug/stats" > /dev/null & ;;
                    4) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                    5) curl -s "$BASE_URL/debug/memory" > /dev/null & ;;
                esac
            done
            
            # Wait for requests to complete
            wait
            
            # Small delay
            sleep 0.1
        done
    ) &
    
    local load_pid=$!
    echo -e "${GREEN}✅ Load generation started (PID: $load_pid)${NC}"
    return $load_pid
}

# Monitor memory usage
monitor_memory() {
    local interval=10
    local count=0
    local profiles=()
    
    echo -e "${BLUE}📈 Monitoring memory usage...${NC}"
    echo -e "Test Duration: ${YELLOW}${TEST_DURATION}s${NC}"
    echo -e "Profile Interval: ${YELLOW}${PROFILE_INTERVAL}s${NC}"
    echo ""
    
    # Get baseline
    local baseline_stats=$(get_memory_stats)
    local baseline_alloc=$(echo $baseline_stats | cut -d' ' -f1)
    local baseline_sys=$(echo $baseline_stats | cut -d' ' -f2)
    
    echo -e "${YELLOW}Baseline Memory:${NC}"
    echo -e "  Heap Allocated: ${YELLOW}$(numfmt --to=iec $baseline_alloc)${NC}"
    echo -e "  Heap System: ${YELLOW}$(numfmt --to=iec $baseline_sys)${NC}"
    echo ""
    
    # Collect baseline profile
    local baseline_timestamp=$(date +%Y%m%d_%H%M%S)
    collect_memory_profile $baseline_timestamp
    profiles+=("$baseline_timestamp")
    
    # Start monitoring loop
    local start_time=$(date +%s)
    local end_time=$((start_time + TEST_DURATION))
    local last_profile_time=$start_time
    
    while [ $(date +%s) -lt $end_time ]; do
        local current_time=$(date +%s)
        local elapsed=$((current_time - start_time))
        
        # Get current memory stats
        local current_stats=$(get_memory_stats)
        local current_alloc=$(echo $current_stats | cut -d' ' -f1)
        local current_sys=$(echo $current_stats | cut -d' ' -f2)
        local current_idle=$(echo $current_stats | cut -d' ' -f3)
        local current_inuse=$(echo $current_stats | cut -d' ' -f4)
        local current_objects=$(echo $current_stats | cut -d' ' -f5)
        
        # Calculate changes
        local alloc_diff=$((current_alloc - baseline_alloc))
        local sys_diff=$((current_sys - baseline_sys))
        
        # Display current status
        echo -e "[${elapsed}s] Heap: ${YELLOW}$(numfmt --to=iec $current_alloc)${NC} (${alloc_diff:+-}${alloc_diff:+$(numfmt --to=iec $alloc_diff)}) | Objects: ${YELLOW}$current_objects${NC}"
        
        # Collect profile at intervals
        if [ $((current_time - last_profile_time)) -ge $PROFILE_INTERVAL ]; then
            local timestamp=$(date +%Y%m%d_%H%M%S)
            collect_memory_profile $timestamp
            profiles+=("$timestamp")
            last_profile_time=$current_time
        fi
        
        sleep $interval
        ((count++))
    done
    
    # Return profile timestamps
    echo "${profiles[@]}"
}

# Analyze memory profiles
analyze_profiles() {
    local profiles=("$@")
    local profile_count=${#profiles[@]}
    
    if [ $profile_count -lt 2 ]; then
        echo -e "${RED}❌ Not enough profiles for analysis${NC}"
        return 1
    fi
    
    echo -e "${BLUE}🔍 Analyzing memory profiles...${NC}"
    echo -e "Total profiles: ${YELLOW}$profile_count${NC}"
    echo ""
    
    # Get first and last profile files
    local first_profile="${PROFILES_DIR}/heap_${profiles[0]}.prof"
    local last_profile="${PROFILES_DIR}/heap_${profiles[-1]}.prof"
    
    if [ ! -f "$first_profile" ] || [ ! -f "$last_profile" ]; then
        echo -e "${RED}❌ Profile files not found${NC}"
        return 1
    fi
    
    echo -e "${YELLOW}Comparing profiles:${NC}"
    echo -e "  Baseline: $first_profile"
    echo -e "  Final: $last_profile"
    echo ""
    
    # Run pprof comparison
    echo -e "${BLUE}📊 Memory growth analysis:${NC}"
    go tool pprof -base "$first_profile" "$last_profile" 2>/dev/null | head -20 || echo "Could not analyze profiles"
    
    echo ""
    echo -e "${BLUE}📊 Top memory consumers in final profile:${NC}"
    go tool pprof -top "$last_profile" 2>/dev/null | head -10 || echo "Could not analyze final profile"
}

# Generate memory leak report
generate_report() {
    local profiles=("$@")
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local report_file="${PROFILES_DIR}/memory_leak_report_${timestamp}.txt"
    
    echo -e "${BLUE}📋 Generating memory leak report...${NC}"
    
    {
        echo "Memory Leak Detection Report"
        echo "============================"
        echo "Date: $(date)"
        echo "Test Duration: ${TEST_DURATION}s"
        echo "Profile Interval: ${PROFILE_INTERVAL}s"
        echo "Load Duration: ${LOAD_DURATION}s"
        echo ""
        echo "Profiles Collected:"
        for profile in "${profiles[@]}"; do
            echo "  - heap_${profile}.prof"
        done
        echo ""
        echo "Analysis Commands:"
        echo "  # Compare first and last profiles:"
        echo "  go tool pprof -base ${PROFILES_DIR}/heap_${profiles[0]}.prof ${PROFILES_DIR}/heap_${profiles[-1]}.prof"
        echo ""
        echo "  # Analyze final profile:"
        echo "  go tool pprof ${PROFILES_DIR}/heap_${profiles[-1]}.prof"
        echo ""
        echo "  # Show memory allocation traces:"
        echo "  go tool pprof -traces ${PROFILES_DIR}/heap_${profiles[-1]}.prof"
        echo ""
        echo "  # Generate web interface:"
        echo "  go tool pprof -web ${PROFILES_DIR}/heap_${profiles[-1]}.prof"
        echo ""
        echo "Memory Leak Indicators:"
        echo "  - Growing heap allocation over time"
        echo "  - Increasing number of objects"
        echo "  - Objects not being garbage collected"
        echo "  - Large memory allocations in specific functions"
    } > "$report_file"
    
    echo -e "${GREEN}✅ Report saved: $report_file${NC}"
}

# Trigger garbage collection and check
trigger_gc() {
    echo -e "${BLUE}🗑️  Triggering garbage collection...${NC}"
    
    if result=$(curl -s -X POST "$BASE_URL/debug/gc" 2>/dev/null); then
        local before_alloc=$(echo "$result" | jq -r '.before.heap_alloc // 0')
        local after_alloc=$(echo "$result" | jq -r '.after.heap_alloc // 0')
        local freed=$((before_alloc - after_alloc))
        
        echo -e "Before GC: ${YELLOW}$(numfmt --to=iec $before_alloc)${NC}"
        echo -e "After GC: ${YELLOW}$(numfmt --to=iec $after_alloc)${NC}"
        echo -e "Freed: ${GREEN}$(numfmt --to=iec $freed)${NC}"
    else
        echo -e "${RED}Failed to trigger garbage collection${NC}"
    fi
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -d, --duration SECONDS    Test duration (default: 300)"
    echo "  -i, --interval SECONDS    Profile collection interval (default: 30)"
    echo "  -l, --load SECONDS        Load generation duration (default: 60)"
    echo "  -o, --output DIR          Output directory (default: memory-leak-profiles)"
    echo "  -g, --gc                  Trigger GC before analysis"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                        # Run with default settings"
    echo "  $0 -d 600 -i 60          # 10min test, profile every minute"
    echo "  $0 -g                    # Trigger GC before analysis"
}

# Parse command line arguments
parse_args() {
    local trigger_gc_flag=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--duration)
                TEST_DURATION="$2"
                shift 2
                ;;
            -i|--interval)
                PROFILE_INTERVAL="$2"
                shift 2
                ;;
            -l|--load)
                LOAD_DURATION="$2"
                shift 2
                ;;
            -o|--output)
                PROFILES_DIR="$2"
                shift 2
                ;;
            -g|--gc)
                trigger_gc_flag=true
                shift
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
    
    # Set global variable
    TRIGGER_GC=$trigger_gc_flag
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    echo -e "${GREEN}Memory Leak Detection with pprof${NC}"
    echo "================================="
    
    # Check prerequisites
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq is required but not installed. Please install jq first.${NC}"
        exit 1
    fi
    
    # Check service and setup
    check_service
    setup_directories
    
    echo -e "${BLUE}🚀 Starting memory leak detection...${NC}"
    echo ""
    
    # Trigger GC if requested
    if [ "$TRIGGER_GC" = true ]; then
        trigger_gc
        echo ""
    fi
    
    # Phase 1: Baseline collection
    echo -e "${YELLOW}Phase 1: Baseline Collection${NC}"
    echo "================================"
    local baseline_stats=$(get_memory_stats)
    local baseline_alloc=$(echo $baseline_stats | cut -d' ' -f1)
    echo -e "Baseline Heap: ${YELLOW}$(numfmt --to=iec $baseline_alloc)${NC}"
    echo ""
    
    # Phase 2: Load generation and monitoring
    echo -e "${YELLOW}Phase 2: Load Generation and Monitoring${NC}"
    echo "============================================="
    generate_load
    local load_pid=$!
    
    # Monitor memory during load
    local profiles=($(monitor_memory))
    
    # Wait for load to complete
    wait $load_pid 2>/dev/null || true
    
    echo ""
    
    # Phase 3: Analysis
    echo -e "${YELLOW}Phase 3: Analysis${NC}"
    echo "====================="
    analyze_profiles "${profiles[@]}"
    generate_report "${profiles[@]}"
    
    echo ""
    echo -e "${GREEN}✅ Memory leak detection completed!${NC}"
    echo -e "${BLUE}📁 Profiles saved in: $PROFILES_DIR${NC}"
    echo -e "${BLUE}📋 Report generated: ${PROFILES_DIR}/memory_leak_report_*.txt${NC}"
    echo ""
    echo -e "${YELLOW}💡 Next steps:${NC}"
    echo "  1. Review the generated report"
    echo "  2. Analyze profiles with: go tool pprof ${PROFILES_DIR}/heap_*.prof"
    echo "  3. Look for growing memory patterns"
    echo "  4. Check for objects not being garbage collected"
}

# Run main function
main "$@"
