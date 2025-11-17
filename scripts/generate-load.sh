#!/bin/bash

# Load Generation Script for Profiling
# This script generates load to help with profiling analysis

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
DURATION=60
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

# Generate load
generate_load() {
    echo -e "${BLUE}🚀 Generating load...${NC}"
    echo -e "Duration: ${YELLOW}${DURATION}s${NC}"
    echo -e "Concurrent Users: ${YELLOW}${CONCURRENT_USERS}${NC}"
    echo -e "Requests/Second: ${YELLOW}${REQUESTS_PER_SECOND}${NC}"
    
    # Calculate delay between requests
    local delay=$(echo "scale=3; 1 / $REQUESTS_PER_SECOND" | bc 2>/dev/null || echo "0.02")
    
    # Start load generation
    local start_time=$(date +%s)
    local end_time=$((start_time + DURATION))
    local request_count=0
    
    while [ $(date +%s) -lt $end_time ]; do
        # Start concurrent requests
        for i in $(seq 1 $CONCURRENT_USERS); do
            # Random endpoint selection
            case $((RANDOM % 6)) in
                0) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                1) curl -s "$BASE_URL/healthz" > /dev/null & ;;
                2) curl -s "$BASE_URL/metrics" > /dev/null & ;;
                3) curl -s "$BASE_URL/debug/stats" > /dev/null & ;;
                4) curl -s "$BASE_URL/v1/users" > /dev/null & ;;
                5) curl -s "$BASE_URL/debug/memory" > /dev/null & ;;
            esac
            ((request_count++))
        done
        
        # Wait for requests to complete
        wait
        
        # Delay between request batches
        sleep "$delay"
    done
    
    echo -e "${GREEN}✅ Load generation completed${NC}"
    echo -e "Total requests: ${YELLOW}$request_count${NC}"
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -d, --duration SECONDS    Load duration (default: 60)"
    echo "  -c, --concurrent USERS    Concurrent users (default: 10)"
    echo "  -r, --rate RPS            Requests per second (default: 50)"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                        # Run with default settings"
    echo "  $0 -d 120 -c 20 -r 100   # 2min, 20 users, 100 RPS"
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--duration)
                DURATION="$2"
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
    
    echo -e "${GREEN}Load Generation for Profiling${NC}"
    echo "==============================="
    
    # Check prerequisites
    if ! command -v bc &> /dev/null; then
        echo -e "${YELLOW}⚠️  bc not found, using default delay${NC}"
    fi
    
    # Check service and generate load
    check_service
    generate_load
}

# Run main function
main "$@"
