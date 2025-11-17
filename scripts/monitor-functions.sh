#!/bin/bash

# Function-Level Resource Monitoring Script
# This script demonstrates how to monitor your service's resource usage per function

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:10000"
DEBUG_PATH="/debug"
METRICS_PATH="/metrics"

# Check if jq is available
if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is required but not installed. Please install jq first.${NC}"
    exit 1
fi

# Function to check if service is running
check_service() {
    if ! curl -s "$BASE_URL/healthz" > /dev/null; then
        echo -e "${RED}❌ Service is not running at $BASE_URL${NC}"
        echo "Please start your service first:"
        echo "  go run cmd/app/main.go"
        exit 1
    fi
    echo -e "${GREEN}✅ Service is running${NC}"
}

# Function to get all function profiles
get_function_profiles() {
    echo -e "\n${BLUE}📊 Function Profiles${NC}"
    echo "=================="
    
    if profiles=$(curl -s "$BASE_URL$DEBUG_PATH/profiles" 2>/dev/null); then
        count=$(echo "$profiles" | jq -r '.count // 0')
        echo -e "Total functions profiled: ${YELLOW}$count${NC}"
        
        if [ "$count" -gt 0 ]; then
            echo "$profiles" | jq -r '.profiles | to_entries[] | "\(.key): \(.value.call_count) calls, avg: \(.value.avg_duration), errors: \(.value.error_count)"'
        else
            echo -e "${YELLOW}No functions have been profiled yet.${NC}"
            echo "Make some API calls to see profiling data."
        fi
    else
        echo -e "${RED}Failed to get function profiles${NC}"
    fi
}

# Function to get runtime statistics
get_runtime_stats() {
    echo -e "\n${BLUE}🔄 Runtime Statistics${NC}"
    echo "=================="
    
    if stats=$(curl -s "$BASE_URL$DEBUG_PATH/stats" 2>/dev/null); then
        goroutines=$(echo "$stats" | jq -r '.goroutines // 0')
        heap_alloc=$(echo "$stats" | jq -r '.memory.alloc // 0')
        heap_sys=$(echo "$stats" | jq -r '.memory.sys // 0')
        heap_objects=$(echo "$stats" | jq -r '.memory.objects // 0')
        num_gc=$(echo "$stats" | jq -r '.memory.num_gc // 0')
        
        echo -e "Goroutines: ${YELLOW}$goroutines${NC}"
        echo -e "Heap Allocated: ${YELLOW}$(numfmt --to=iec $heap_alloc)${NC}"
        echo -e "Heap System: ${YELLOW}$(numfmt --to=iec $heap_sys)${NC}"
        echo -e "Heap Objects: ${YELLOW}$heap_objects${NC}"
        echo -e "GC Count: ${YELLOW}$num_gc${NC}"
    else
        echo -e "${RED}Failed to get runtime statistics${NC}"
    fi
}

# Function to get memory information
get_memory_info() {
    echo -e "\n${BLUE}💾 Memory Information${NC}"
    echo "=================="
    
    if memory=$(curl -s "$BASE_URL$DEBUG_PATH/memory" 2>/dev/null); then
        heap_alloc=$(echo "$memory" | jq -r '.heap.alloc // 0')
        heap_sys=$(echo "$memory" | jq -r '.heap.sys // 0')
        heap_idle=$(echo "$memory" | jq -r '.heap.idle // 0')
        heap_inuse=$(echo "$memory" | jq -r '.heap.inuse // 0')
        stack_inuse=$(echo "$memory" | jq -r '.stack.inuse // 0')
        
        echo -e "Heap Allocated: ${YELLOW}$(numfmt --to=iec $heap_alloc)${NC}"
        echo -e "Heap System: ${YELLOW}$(numfmt --to=iec $heap_sys)${NC}"
        echo -e "Heap Idle: ${YELLOW}$(numfmt --to=iec $heap_idle)${NC}"
        echo -e "Heap In-Use: ${YELLOW}$(numfmt --to=iec $heap_inuse)${NC}"
        echo -e "Stack In-Use: ${YELLOW}$(numfmt --to=iec $stack_inuse)${NC}"
        
        # Calculate memory utilization
        if [ "$heap_sys" -gt 0 ]; then
            utilization=$((heap_inuse * 100 / heap_sys))
            echo -e "Memory Utilization: ${YELLOW}${utilization}%${NC}"
        fi
    else
        echo -e "${RED}Failed to get memory information${NC}"
    fi
}

# Function to get top slow functions
get_slow_functions() {
    echo -e "\n${BLUE}🐌 Slowest Functions${NC}"
    echo "=================="
    
    if profiles=$(curl -s "$BASE_URL$DEBUG_PATH/profiles" 2>/dev/null); then
        echo "$profiles" | jq -r '
            .profiles | to_entries[] | 
            select(.value.call_count > 0) | 
            "\(.value.avg_duration) - \(.key) (\(.value.call_count) calls)"
        ' | sort -r | head -10
    else
        echo -e "${RED}Failed to get function profiles${NC}"
    fi
}

# Function to get functions with errors
get_error_functions() {
    echo -e "\n${BLUE}❌ Functions with Errors${NC}"
    echo "========================"
    
    if profiles=$(curl -s "$BASE_URL$DEBUG_PATH/profiles" 2>/dev/null); then
        echo "$profiles" | jq -r '
            .profiles | to_entries[] | 
            select(.value.error_count > 0) | 
            "\(.value.error_count) errors - \(.key) (\(.value.call_count) total calls)"
        ' | sort -r
    else
        echo -e "${RED}Failed to get function profiles${NC}"
    fi
}

# Function to trigger garbage collection
trigger_gc() {
    echo -e "\n${BLUE}🗑️  Triggering Garbage Collection${NC}"
    echo "================================"
    
    if result=$(curl -s -X POST "$BASE_URL$DEBUG_PATH/gc" 2>/dev/null); then
        echo "$result" | jq -r '.message'
        before_alloc=$(echo "$result" | jq -r '.before.heap_alloc // 0')
        after_alloc=$(echo "$result" | jq -r '.after.heap_alloc // 0')
        freed=$((before_alloc - after_alloc))
        
        if [ "$freed" -gt 0 ]; then
            echo -e "Freed: ${GREEN}$(numfmt --to=iec $freed)${NC}"
        else
            echo -e "No memory freed"
        fi
    else
        echo -e "${RED}Failed to trigger garbage collection${NC}"
    fi
}

# Function to show Prometheus metrics
show_prometheus_metrics() {
    echo -e "\n${BLUE}📈 Prometheus Metrics${NC}"
    echo "==================="
    
    if metrics=$(curl -s "$BASE_URL$METRICS_PATH" 2>/dev/null); then
        echo "$metrics" | grep -E "(function_|goroutines_|heap_)" | head -20
        echo -e "\n${YELLOW}... (showing first 20 metrics)${NC}"
        echo "Full metrics available at: $BASE_URL$METRICS_PATH"
    else
        echo -e "${RED}Failed to get Prometheus metrics${NC}"
    fi
}

# Function to show available endpoints
show_endpoints() {
    echo -e "\n${BLUE}🔗 Available Endpoints${NC}"
    echo "====================="
    echo -e "Function Profiles: ${YELLOW}$BASE_URL$DEBUG_PATH/profiles${NC}"
    echo -e "Runtime Stats: ${YELLOW}$BASE_URL$DEBUG_PATH/stats${NC}"
    echo -e "Memory Info: ${YELLOW}$BASE_URL$DEBUG_PATH/memory${NC}"
    echo -e "Goroutines: ${YELLOW}$BASE_URL$DEBUG_PATH/goroutines${NC}"
    echo -e "Trigger GC: ${YELLOW}$BASE_URL$DEBUG_PATH/gc${NC}"
    echo -e "Prometheus Metrics: ${YELLOW}$BASE_URL$METRICS_PATH${NC}"
    echo -e "pprof CPU: ${YELLOW}$BASE_URL$DEBUG_PATH/pprof/profile${NC}"
    echo -e "pprof Memory: ${YELLOW}$BASE_URL$DEBUG_PATH/pprof/heap${NC}"
}

# Function to generate load for testing
generate_load() {
    echo -e "\n${BLUE}🚀 Generating Load for Testing${NC}"
    echo "============================="
    
    echo "Making API calls to generate profiling data..."
    
    # Make some API calls to generate profiling data
    for i in {1..5}; do
        echo "Call $i: GET /v1/users"
        curl -s "$BASE_URL/v1/users" > /dev/null &
        
        echo "Call $i: GET /healthz"
        curl -s "$BASE_URL/healthz" > /dev/null &
        
        echo "Call $i: GET /metrics"
        curl -s "$BASE_URL/metrics" > /dev/null &
        
        sleep 0.5
    done
    
    wait
    echo -e "${GREEN}Load generation completed${NC}"
}

# Main menu
show_menu() {
    echo -e "\n${BLUE}🔍 Function Resource Monitoring${NC}"
    echo "================================"
    echo "1. Check all function profiles"
    echo "2. Show runtime statistics"
    echo "3. Show memory information"
    echo "4. Show slowest functions"
    echo "5. Show functions with errors"
    echo "6. Trigger garbage collection"
    echo "7. Show Prometheus metrics"
    echo "8. Generate test load"
    echo "9. Show all endpoints"
    echo "0. Exit"
    echo ""
    read -p "Select an option (0-9): " choice
}

# Main function
main() {
    echo -e "${GREEN}Function-Level Resource Monitoring Tool${NC}"
    echo "============================================="
    
    # Check if service is running
    check_service
    
    while true; do
        show_menu
        
        case $choice in
            1)
                get_function_profiles
                ;;
            2)
                get_runtime_stats
                ;;
            3)
                get_memory_info
                ;;
            4)
                get_slow_functions
                ;;
            5)
                get_error_functions
                ;;
            6)
                trigger_gc
                ;;
            7)
                show_prometheus_metrics
                ;;
            8)
                generate_load
                ;;
            9)
                show_endpoints
                ;;
            0)
                echo -e "${GREEN}Goodbye!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}Invalid option. Please try again.${NC}"
                ;;
        esac
        
        echo ""
        read -p "Press Enter to continue..."
    done
}

# Run main function
main "$@" 