#!/bin/bash

# Automated pprof Profile Analysis Script
# This script analyzes existing pprof profiles and generates comprehensive reports

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROFILES_DIR="."
OUTPUT_DIR="analysis-reports"
ANALYSIS_DEPTH=20

# Check if profile file exists
check_profile() {
    local profile_file=$1
    if [ ! -f "$profile_file" ]; then
        echo -e "${RED}❌ Profile file not found: $profile_file${NC}"
        return 1
    fi
    
    local size=$(du -h "$profile_file" | cut -f1)
    echo -e "${GREEN}✅ Found profile: $profile_file (${size})${NC}"
    return 0
}

# Analyze CPU profile
analyze_cpu_profile() {
    local profile_file=$1
    local output_file=$2
    
    echo -e "${BLUE}📊 Analyzing CPU profile: $profile_file${NC}"
    
    {
        echo "CPU Profile Analysis"
        echo "==================="
        echo "File: $profile_file"
        echo "Date: $(date)"
        echo ""
        
        echo "Top CPU Consumers:"
        echo "-----------------"
        go tool pprof -top "$profile_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not analyze CPU profile"
        echo ""
        
        echo "Top CPU Consumers (Cumulative):"
        echo "-------------------------------"
        go tool pprof -top -cum "$profile_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not analyze CPU profile"
        echo ""
        
        echo "Function Details:"
        echo "----------------"
        go tool pprof -list=. "$profile_file" 2>/dev/null | head -50 || echo "Could not list functions"
        echo ""
        
        echo "Call Graph (Top 10):"
        echo "-------------------"
        go tool pprof -weblist=. "$profile_file" 2>/dev/null | head -20 || echo "Could not generate call graph"
        
    } > "$output_file"
    
    echo -e "${GREEN}✅ CPU analysis saved: $output_file${NC}"
}

# Analyze memory profile
analyze_memory_profile() {
    local profile_file=$1
    local output_file=$2
    
    echo -e "${BLUE}💾 Analyzing memory profile: $profile_file${NC}"
    
    {
        echo "Memory Profile Analysis"
        echo "======================"
        echo "File: $profile_file"
        echo "Date: $(date)"
        echo ""
        
        echo "Top Memory Consumers:"
        echo "-------------------"
        go tool pprof -top "$profile_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not analyze memory profile"
        echo ""
        
        echo "Top Memory Consumers (Cumulative):"
        echo "--------------------------------"
        go tool pprof -top -cum "$profile_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not analyze memory profile"
        echo ""
        
        echo "Memory Allocation Traces:"
        echo "------------------------"
        go tool pprof -traces "$profile_file" 2>/dev/null | head -30 || echo "Could not show allocation traces"
        echo ""
        
        echo "Function Details:"
        echo "----------------"
        go tool pprof -list=. "$profile_file" 2>/dev/null | head -50 || echo "Could not list functions"
        echo ""
        
        echo "Memory Allocation by Function:"
        echo "----------------------------"
        go tool pprof -weblist=. "$profile_file" 2>/dev/null | head -20 || echo "Could not generate allocation details"
        
    } > "$output_file"
    
    echo -e "${GREEN}✅ Memory analysis saved: $output_file${NC}"
}

# Analyze goroutine profile
analyze_goroutine_profile() {
    local profile_file=$1
    local output_file=$2
    
    echo -e "${BLUE}🔄 Analyzing goroutine profile: $profile_file${NC}"
    
    {
        echo "Goroutine Profile Analysis"
        echo "========================="
        echo "File: $profile_file"
        echo "Date: $(date)"
        echo ""
        
        echo "Goroutine Count by Function:"
        echo "---------------------------"
        go tool pprof -top "$profile_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not analyze goroutine profile"
        echo ""
        
        echo "Goroutine Stack Traces:"
        echo "----------------------"
        go tool pprof -traces "$profile_file" 2>/dev/null | head -30 || echo "Could not show goroutine traces"
        echo ""
        
        echo "Function Details:"
        echo "----------------"
        go tool pprof -list=. "$profile_file" 2>/dev/null | head -50 || echo "Could not list functions"
        echo ""
        
        echo "Goroutine Call Graph:"
        echo "-------------------"
        go tool pprof -weblist=. "$profile_file" 2>/dev/null | head -20 || echo "Could not generate call graph"
        
    } > "$output_file"
    
    echo -e "${GREEN}✅ Goroutine analysis saved: $output_file${NC}"
}

# Compare two profiles
compare_profiles() {
    local baseline_file=$1
    local current_file=$2
    local output_file=$3
    
    echo -e "${BLUE}📊 Comparing profiles:${NC}"
    echo -e "  Baseline: $baseline_file"
    echo -e "  Current: $current_file"
    
    {
        echo "Profile Comparison Analysis"
        echo "=========================="
        echo "Baseline: $baseline_file"
        echo "Current: $current_file"
        echo "Date: $(date)"
        echo ""
        
        echo "Memory Growth Analysis:"
        echo "---------------------"
        go tool pprof -base "$baseline_file" "$current_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not compare profiles"
        echo ""
        
        echo "Top Changes:"
        echo "-----------"
        go tool pprof -diff_base "$baseline_file" "$current_file" 2>/dev/null | head -$ANALYSIS_DEPTH || echo "Could not show differences"
        echo ""
        
        echo "Detailed Comparison:"
        echo "------------------"
        go tool pprof -list=. -base "$baseline_file" "$current_file" 2>/dev/null | head -50 || echo "Could not show detailed comparison"
        
    } > "$output_file"
    
    echo -e "${GREEN}✅ Comparison analysis saved: $output_file${NC}"
}

# Generate summary report
generate_summary_report() {
    local profiles_dir=$1
    local output_file=$2
    
    echo -e "${BLUE}📋 Generating summary report...${NC}"
    
    {
        echo "pprof Profile Analysis Summary"
        echo "============================="
        echo "Date: $(date)"
        echo "Profiles Directory: $profiles_dir"
        echo ""
        
        echo "Available Profiles:"
        echo "------------------"
        find "$profiles_dir" -name "*.prof" -type f | while read -r profile; do
            local size=$(du -h "$profile" | cut -f1)
            local modified=$(stat -f "%Sm" "$profile" 2>/dev/null || stat -c "%y" "$profile" 2>/dev/null)
            echo "  - $profile (${size}, modified: $modified)"
        done
        echo ""
        
        echo "Analysis Reports Generated:"
        echo "-------------------------"
        find "$OUTPUT_DIR" -name "*.txt" -type f | while read -r report; do
            local size=$(du -h "$report" | cut -f1)
            echo "  - $report (${size})"
        done
        echo ""
        
        echo "Quick Analysis Commands:"
        echo "----------------------"
        echo "  # CPU analysis:"
        echo "  go tool pprof -web \$PROFILE"
        echo ""
        echo "  # Memory analysis:"
        echo "  go tool pprof -top \$PROFILE"
        echo ""
        echo "  # Goroutine analysis:"
        echo "  go tool pprof -traces \$PROFILE"
        echo ""
        echo "  # Profile comparison:"
        echo "  go tool pprof -base baseline.prof current.prof"
        echo ""
        echo "  # Generate PDF report:"
        echo "  go tool pprof -pdf \$PROFILE > report.pdf"
        echo ""
        echo "  # Generate SVG report:"
        echo "  go tool pprof -svg \$PROFILE > report.svg"
        
    } > "$output_file"
    
    echo -e "${GREEN}✅ Summary report saved: $output_file${NC}"
}

# Find and analyze all profiles
analyze_all_profiles() {
    local profiles_dir=$1
    
    echo -e "${BLUE}🔍 Finding profiles in: $profiles_dir${NC}"
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR"
    
    local cpu_profiles=()
    local memory_profiles=()
    local goroutine_profiles=()
    
    # Find all profile files
    while IFS= read -r -d '' profile; do
        local basename=$(basename "$profile")
        
        if [[ "$basename" == *"cpu"* ]]; then
            cpu_profiles+=("$profile")
        elif [[ "$basename" == *"heap"* ]] || [[ "$basename" == *"memory"* ]]; then
            memory_profiles+=("$profile")
        elif [[ "$basename" == *"goroutine"* ]]; then
            goroutine_profiles+=("$profile")
        fi
    done < <(find "$profiles_dir" -name "*.prof" -type f -print0)
    
    echo -e "${GREEN}Found profiles:${NC}"
    echo -e "  CPU: ${#cpu_profiles[@]}"
    echo -e "  Memory: ${#memory_profiles[@]}"
    echo -e "  Goroutine: ${#goroutine_profiles[@]}"
    echo ""
    
    # Analyze CPU profiles
    for profile in "${cpu_profiles[@]}"; do
        local basename=$(basename "$profile" .prof)
        local output_file="$OUTPUT_DIR/cpu_analysis_${basename}.txt"
        analyze_cpu_profile "$profile" "$output_file"
    done
    
    # Analyze memory profiles
    for profile in "${memory_profiles[@]}"; do
        local basename=$(basename "$profile" .prof)
        local output_file="$OUTPUT_DIR/memory_analysis_${basename}.txt"
        analyze_memory_profile "$profile" "$output_file"
    done
    
    # Analyze goroutine profiles
    for profile in "${goroutine_profiles[@]}"; do
        local basename=$(basename "$profile" .prof)
        local output_file="$OUTPUT_DIR/goroutine_analysis_${basename}.txt"
        analyze_goroutine_profile "$profile" "$output_file"
    done
    
    # Generate comparisons if we have multiple profiles of the same type
    if [ ${#memory_profiles[@]} -ge 2 ]; then
        local first_memory="${memory_profiles[0]}"
        local last_memory="${memory_profiles[-1]}"
        local comparison_file="$OUTPUT_DIR/memory_comparison_$(basename "$first_memory" .prof)_vs_$(basename "$last_memory" .prof).txt"
        compare_profiles "$first_memory" "$last_memory" "$comparison_file"
    fi
    
    # Generate summary report
    local summary_file="$OUTPUT_DIR/analysis_summary.txt"
    generate_summary_report "$profiles_dir" "$summary_file"
}

# Analyze specific profile
analyze_specific_profile() {
    local profile_file=$1
    local profile_type=$2
    
    echo -e "${BLUE}🔍 Analyzing specific profile: $profile_file${NC}"
    
    # Check if profile exists
    if ! check_profile "$profile_file"; then
        return 1
    fi
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR"
    
    local basename=$(basename "$profile_file" .prof)
    local output_file="$OUTPUT_DIR/${profile_type}_analysis_${basename}.txt"
    
    case $profile_type in
        cpu)
            analyze_cpu_profile "$profile_file" "$output_file"
            ;;
        memory)
            analyze_memory_profile "$profile_file" "$output_file"
            ;;
        goroutine)
            analyze_goroutine_profile "$profile_file" "$output_file"
            ;;
        *)
            echo -e "${RED}❌ Unknown profile type: $profile_type${NC}"
            echo "Supported types: cpu, memory, goroutine"
            return 1
            ;;
    esac
}

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] [PROFILE_FILE]"
    echo ""
    echo "Options:"
    echo "  -d, --dir DIRECTORY       Directory containing profiles (default: current)"
    echo "  -o, --output DIRECTORY    Output directory for reports (default: analysis-reports)"
    echo "  -t, --type TYPE           Profile type: cpu, memory, goroutine"
    echo "  -c, --compare FILE1 FILE2 Compare two profiles"
    echo "  -a, --all                 Analyze all profiles in directory"
    echo "  -h, --help               Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -a                     # Analyze all profiles in current directory"
    echo "  $0 -d profiles -a         # Analyze all profiles in 'profiles' directory"
    echo "  $0 -t cpu cpu.prof        # Analyze specific CPU profile"
    echo "  $0 -c baseline.prof current.prof  # Compare two profiles"
    echo "  $0 -o reports -a          # Save reports to 'reports' directory"
}

# Parse command line arguments
parse_args() {
    local analyze_all=false
    local compare_files=()
    local specific_profile=""
    local profile_type=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -d|--dir)
                PROFILES_DIR="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -t|--type)
                profile_type="$2"
                shift 2
                ;;
            -c|--compare)
                compare_files+=("$2" "$3")
                shift 3
                ;;
            -a|--all)
                analyze_all=true
                shift
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            -*)
                echo -e "${RED}Unknown option: $1${NC}"
                show_usage
                exit 1
                ;;
            *)
                if [ -z "$specific_profile" ]; then
                    specific_profile="$1"
                else
                    echo -e "${RED}Too many arguments${NC}"
                    show_usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    # Set global variables
    ANALYZE_ALL=$analyze_all
    COMPARE_FILES=("${compare_files[@]}")
    SPECIFIC_PROFILE=$specific_profile
    PROFILE_TYPE=$profile_type
}

# Main function
main() {
    # Parse arguments
    parse_args "$@"
    
    echo -e "${GREEN}Automated pprof Profile Analysis${NC}"
    echo "================================="
    
    # Check prerequisites
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ go tool is required but not installed.${NC}"
        exit 1
    fi
    
    # Handle different analysis modes
    if [ ${#COMPARE_FILES[@]} -eq 2 ]; then
        # Compare two profiles
        local baseline_file="${COMPARE_FILES[0]}"
        local current_file="${COMPARE_FILES[1]}"
        local comparison_file="$OUTPUT_DIR/profile_comparison_$(basename "$baseline_file" .prof)_vs_$(basename "$current_file" .prof).txt"
        
        mkdir -p "$OUTPUT_DIR"
        compare_profiles "$baseline_file" "$current_file" "$comparison_file"
        
    elif [ "$ANALYZE_ALL" = true ]; then
        # Analyze all profiles in directory
        analyze_all_profiles "$PROFILES_DIR"
        
    elif [ -n "$SPECIFIC_PROFILE" ]; then
        # Analyze specific profile
        if [ -z "$PROFILE_TYPE" ]; then
            echo -e "${RED}❌ Profile type (-t) is required for specific profile analysis${NC}"
            show_usage
            exit 1
        fi
        
        analyze_specific_profile "$SPECIFIC_PROFILE" "$PROFILE_TYPE"
        
    else
        # Default: analyze all profiles in current directory
        echo -e "${YELLOW}No specific mode specified, analyzing all profiles in current directory...${NC}"
        analyze_all_profiles "$PROFILES_DIR"
    fi
    
    echo ""
    echo -e "${GREEN}✅ Analysis completed!${NC}"
    echo -e "${BLUE}📁 Reports saved in: $OUTPUT_DIR${NC}"
    echo ""
    echo -e "${YELLOW}💡 Next steps:${NC}"
    echo "  1. Review the generated reports"
    echo "  2. Use 'go tool pprof -web' for interactive analysis"
    echo "  3. Generate visual reports with -pdf or -svg options"
}

# Run main function
main "$@"
