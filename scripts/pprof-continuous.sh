#!/bin/bash

# Continuous pprof Profiling Script
# This script continuously collects pprof profiles for ongoing performance monitoring

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
PROFILES_DIR="profiles"
CPU_DURATION=30
INTERVAL=300  # 5 minutes
MAX_PROFILES=100  # Keep last 100 profiles

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
    echo -e "${BLUE}📁 Profiles will be saved to: $PROFILES_DIR${NC}"
}

# Clean old profiles
cleanup_old_profiles() {
    local profile_type=$1
    local count=$(ls -1 "$PROFILES_DIR"/${profile_type}_*.prof 2>/dev/null | wc -l)
    
    if [ "$count" -gt "$MAX_PROFILES" ]; then
        echo -e "${YELLOW}🧹 Cleaning old ${profile_type} profiles...${NC}"
        ls -t "$PROFILES_DIR"/${profile_type}_*.prof | tail -n +$((MAX_PROFILES + 1)) | xargs rm -f
    fi
}

# Collect CPU profile
collect_cpu_profile() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local filename="${PROFILES_DIR}/cpu_${timestamp}.prof"
    
    echo -e "${BLUE}📊 Collecting CPU profile (${CPU_DURATION}s)...${NC}"
    
    if curl -s -o "$filename" "$BASE_URL/debug/pprof/profile?seconds=${CPU_DURATION}"; then
        local size=$(du -h "$filename" | cut -f1)
        echo -e "${GREEN}✅ CPU profile saved: $filename (${size})${NC}"
        
        # Show top functions
        echo -e "${YELLOW}🔍 Top CPU consumers:${NC}"
        go tool pprof -top "$filename" 2>/dev/null | head -10 || echo "Could not analyze profile"
    else
        echo -e "${RED}❌ Failed to collect CPU profile${NC}"
    fi
    
    cleanup_old_profiles "cpu"
}

# Collect memory profile
collect_memory_profile() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local filename="${PROFILES_DIR}/heap_${timestamp}.prof"
    
    echo -e "${BLUE}💾 Collecting memory profile...${NC}"
    
    if curl -s -o "$filename" "$BASE_URL/debug/pprof/heap"; then
        local size=$(du -h "$filename" | cut -f1)
        echo -e "${GREEN}✅ Memory profile saved: $filename (${size})${NC}"
        
        # Show top memory consumers
        echo -e "${YELLOW}🔍 Top memory consumers:${NC}"
        go tool pprof -top "$filename" 2>/dev/null | head -10 || echo "Could not analyze profile"
    else
        echo -e "${RED}❌ Failed to collect memory profile${NC}"
    fi
    
    cleanup_old_profiles "heap"
}

# Collect goroutine profile
collect_goroutine_profile() {
    local timestamp=$(date +%Y%m%d_%H%M%S)
    local filename="${PROFILES_DIR}/goroutine_${timestamp}.prof"
    
    echo -e "${BLUE}🔄 Collecting goroutine profile...${NC}"
    
    if curl -s -o "$filename" "$BASE_URL/debug/pprof/goroutine"; then
        local size=$(du -h "$filename" | cut -f1)
        echo -e "${GREEN}✅ Goroutine profile saved: $filename (${size})${NC}"
        
        # Show goroutine count
        local goroutines=$(go tool pprof -raw "$filename" 2>/dev/null | grep -c "goroutine" || echo "0")
        echo -e "${YELLOW}🔍 Active goroutines: $goroutines${NC}"
    else
        echo -e "${RED}❌ Failed to collect goroutine profile${NC}"
    fi
    
    cleanup_old_profiles "goroutine"
}

# Get runtime statistics
get_runtime_stats() {
    echo -e "${BLUE}📈 Runtime Statistics${NC}"
    echo "=================="
    
    if stats=$(curl -s "$BASE_URL/debug/stats" 2>/dev/null); then
        goroutines=$(echo "$stats" | jq -r '.goroutines // 0')
        heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
        heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
        
        echo -e "Goroutines: ${YELLOW}$goroutines${NC}"
        echo -e "Heap Allocated: ${YELLOW}$(numfmt --to=iec $heap_alloc)${NC}"
        echo -e "Heap System: ${YELLOW}$(numfmt --to=iec $heap_sys)${NC}"
    else
        echo -e "${RED}Failed to get runtime statistics${NC}"
    fi
}

# Generate load for better profiling
generate_load() {
    echo -e "${BLUE}🚀 Generating load for profiling...${NC}"
    
    # Make concurrent requests
    for i in {1..10}; do
        curl -s "$BASE_URL/v1/users" > /dev/null &
        curl -s "$BASE_URL/healthz" > /dev/null &
        curl -s "$BASE_URL/metrics" > /dev/null &
    done
    
    wait
    echo -e "${GREEN}✅ Load generation completed${NC}"
}

# Main profiling loop
main_profiling_loop() {
    local cycle=1
    
    echo -e "${GREEN}🔄 Starting continuous profiling${NC}"
    echo -e "Interval: ${YELLOW}${INTERVAL}s${NC}"
    echo -e "CPU Duration: ${YELLOW}${CPU_DURATION}s${NC}"
    echo -e "Max Profiles: ${YELLOW}${MAX_PROFILES}${NC}"
    echo ""
    
    while true; do
        echo -e "${BLUE}🔄 Profiling Cycle #$cycle${NC}"
        echo "=================="
        
        # Get runtime stats
        get_runtime_stats
        
        # Generate some load
        generate_load
        
        # Collect profiles
        collect_cpu_profile
        collect_memory_profile
        collect_goroutine_profile
        
        echo ""
        echo -e "${YELLOW}⏰ Next profile collection in ${INTERVAL}s...${NC}"
        echo "Press Ctrl+C to stop"
        echo ""
        
        sleep "$INTERVAL"
        ((cycle++))
    done
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -i, --interval SECONDS    Profile collection interval (default: 300)"
    echo "  -d, --duration SECONDS    CPU profile duration (default: 30)"
    echo "  -m, --max-profiles COUNT  Maximum profiles to keep (default: 100)"
    echo "  -o, --output DIR          Output directory (default: profiles)"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                        # Start with default settings"
    echo "  $0 -i 60 -d 15           # Profile every 60s, CPU for 15s"
    echo "  $0 -o /tmp/profiles      # Save profiles to /tmp/profiles"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -i|--interval)
                INTERVAL="$2"
                shift 2
                ;;
            -d|--duration)
                CPU_DURATION="$2"
                shift 2
                ;;
            -m|--max-profiles)
                MAX_PROFILES="$2"
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

# Signal handler
cleanup() {
    echo ""
    echo -e "${YELLOW}🛑 Stopping continuous profiling...${NC}"
    echo -e "${GREEN}✅ Profiles saved in: $PROFILES_DIR${NC}"
    exit 0
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    # Setup signal handler
    trap cleanup SIGINT SIGTERM
    
    echo -e "${GREEN}Continuous pprof Profiling Tool${NC}"
    echo "====================================="
    
    # Check prerequisites
    if ! command -v jq &> /dev/null; then
        echo -e "${RED}❌ jq is required but not installed. Please install jq first.${NC}"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ go tool is required but not installed.${NC}"
        exit 1
    fi
    
    # Check service and setup
    check_service
    setup_directories
    
    # Start profiling loop
    main_profiling_loop
}

# Run main function
main "$@"
