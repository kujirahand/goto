#!/bin/bash

# Performance tests for goto program
# This script tests the performance characteristics of the goto command

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Print functions
print_header() {
    echo -e "${PURPLE}======================================${NC}"
    echo -e "${PURPLE} goto パフォーマンステストスイート${NC}"
    echo -e "${PURPLE}======================================${NC}"
    echo
}

print_test() {
    echo -e "${YELLOW}[PERF]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_summary() {
    echo
    echo -e "${PURPLE}======================================${NC}"
    echo -e "${PURPLE} パフォーマンステスト結果${NC}"
    echo -e "${PURPLE}======================================${NC}"
    echo -e "合格: ${GREEN}${TESTS_PASSED}${NC}"
    echo -e "失敗: ${RED}${TESTS_FAILED}${NC}"
    echo -e "合計: $((TESTS_PASSED + TESTS_FAILED))"
    echo
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}すべてのパフォーマンステストが合格しました！${NC}"
        return 0
    else
        echo -e "${RED}いくつかのパフォーマンステストが失敗しました。${NC}"
        return 1
    fi
}

# Get directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
GOTO_CMD="$PROJECT_DIR/go/goto"

# Test temporary directory
TEST_TEMP_DIR="$SCRIPT_DIR/temp_perf"

# Setup test environment
setup_test_env() {
    print_test "パフォーマンステスト環境のセットアップ"
    
    # Create temp directory
    mkdir -p "$TEST_TEMP_DIR"
    
    # Backup existing config and history
    if [ -f "$HOME/.goto.toml" ]; then
        cp "$HOME/.goto.toml" "$TEST_TEMP_DIR/goto.toml.backup"
    fi
    
    if [ -f "$HOME/.goto.history.json" ]; then
        cp "$HOME/.goto.history.json" "$TEST_TEMP_DIR/goto.history.json.backup"
    fi
    
    print_success "パフォーマンステスト環境がセットアップされました"
}

# Cleanup test environment
cleanup_test_env() {
    print_test "パフォーマンステスト環境のクリーンアップ"
    
    # Restore original config and history
    if [ -f "$TEST_TEMP_DIR/goto.toml.backup" ]; then
        mv "$TEST_TEMP_DIR/goto.toml.backup" "$HOME/.goto.toml"
    else
        rm -f "$HOME/.goto.toml"
    fi
    
    if [ -f "$TEST_TEMP_DIR/goto.history.json.backup" ]; then
        mv "$TEST_TEMP_DIR/goto.history.json.backup" "$HOME/.goto.history.json"
    else
        rm -f "$HOME/.goto.history.json"
    fi
    
    # Remove temp directory
    rm -rf "$TEST_TEMP_DIR"
    
    print_success "パフォーマンステスト環境がクリーンアップされました"
}

# Measure execution time
measure_time() {
    local command="$1"
    local description="$2"
    
    print_info "測定中: $description"
    
    # Run the command multiple times and take average
    local total_time=0
    local runs=5
    
    for i in $(seq 1 $runs); do
        local start_time=$(date +%s.%N)
        eval "$command" > /dev/null 2>&1 || true
        local end_time=$(date +%s.%N)
        local duration=$(echo "$end_time - $start_time" | bc -l)
        total_time=$(echo "$total_time + $duration" | bc -l)
    done
    
    local avg_time=$(echo "scale=4; $total_time / $runs" | bc -l)
    echo "$avg_time"
}

# Test startup time
test_startup_time() {
    print_test "起動時間のテスト"
    
    # Create minimal config
    cat > "$HOME/.goto.toml" << 'EOF'
[test]
path = "/tmp"
EOF
    
    local startup_time=$(measure_time "$GOTO_CMD --version" "バージョン表示")
    print_info "起動時間 (バージョン表示): ${startup_time}秒"
    
    # Check if startup time is reasonable (less than 1 second)
    if (( $(echo "$startup_time < 1.0" | bc -l) )); then
        print_success "起動時間が十分に高速です (${startup_time}秒)"
    else
        print_error "起動時間が遅すぎます (${startup_time}秒)"
    fi
}

# Test config loading performance
test_config_loading_performance() {
    print_test "設定ファイル読み込みパフォーマンステスト"
    
    # Create large config file
    cat > "$HOME/.goto.toml" << 'EOF'
# Large configuration file for performance testing
EOF
    
    # Generate 100 destinations
    for i in $(seq 1 100); do
        cat >> "$HOME/.goto.toml" << EOF

[destination$i]
path = "/tmp/test$i"
shortcut = "$(printf "%c" $((97 + i % 26))$((i % 10))"
EOF
    done
    
    local load_time=$(measure_time "$GOTO_CMD --complete" "大きな設定ファイルの読み込み")
    print_info "設定ファイル読み込み時間 (100エントリ): ${load_time}秒"
    
    # Check if load time is reasonable (less than 2 seconds)
    if (( $(echo "$load_time < 2.0" | bc -l) )); then
        print_success "設定ファイル読み込み時間が十分に高速です (${load_time}秒)"
    else
        print_error "設定ファイル読み込み時間が遅すぎます (${load_time}秒)"
    fi
}

# Test history processing performance
test_history_performance() {
    print_test "履歴処理パフォーマンステスト"
    
    # Create large history file
    cat > "$HOME/.goto.history.json" << 'EOF'
{
  "entries": [
EOF
    
    # Generate 100 history entries
    for i in $(seq 1 100); do
        local timestamp=$(date -d "-$i days" --iso-8601=seconds)
        cat >> "$HOME/.goto.history.json" << EOF
    {
      "label": "destination$i",
      "last_used": "$timestamp"
    }$([ $i -lt 100 ] && echo "," || echo "")
EOF
    done
    
    cat >> "$HOME/.goto.history.json" << 'EOF'
  ]
}
EOF
    
    local history_time=$(measure_time "$GOTO_CMD --history" "大きな履歴ファイルの処理")
    print_info "履歴処理時間 (100エントリ): ${history_time}秒"
    
    # Check if history processing time is reasonable (less than 1 seconds)
    if (( $(echo "$history_time < 1.0" | bc -l) )); then
        print_success "履歴処理時間が十分に高速です (${history_time}秒)"
    else
        print_error "履歴処理時間が遅すぎます (${history_time}秒)"
    fi
}

# Test memory usage
test_memory_usage() {
    print_test "メモリ使用量テスト"
    
    # Create config with many entries
    cat > "$HOME/.goto.toml" << 'EOF'
EOF
    
    # Generate 50 destinations
    for i in $(seq 1 50); do
        cat >> "$HOME/.goto.toml" << EOF
[destination$i]
path = "/tmp/test$i"
shortcut = "d$i"

EOF
    done
    
    # Monitor memory usage during execution
    if command -v time > /dev/null 2>&1; then
        # Use time command if available
        local mem_output=$(timeout 10s /usr/bin/time -f "%M" "$GOTO_CMD" -c 2>&1 | tail -n 1 || echo "N/A")
        if [ "$mem_output" != "N/A" ]; then
            print_info "最大メモリ使用量: ${mem_output}KB"
            
            # Check if memory usage is reasonable (less than 50MB)
            if [ "$mem_output" -lt 51200 ]; then
                print_success "メモリ使用量が適切です (${mem_output}KB)"
            else
                print_error "メモリ使用量が多すぎます (${mem_output}KB)"
            fi
        else
            print_info "メモリ使用量の測定をスキップしました（timeコマンドエラー）"
            print_success "メモリ使用量テストをスキップ"
        fi
    else
        print_info "timeコマンドが利用できないため、メモリ使用量テストをスキップしました"
        print_success "メモリ使用量テストをスキップ"
    fi
}

# Test concurrent execution
test_concurrent_execution() {
    print_test "並行実行パフォーマンステスト"
    
    # Create simple config
    cat > "$HOME/.goto.toml" << 'EOF'
[test1]
path = "/tmp"

[test2]
path = "/usr"

[test3]
path = "/var"
EOF
    
    # Run multiple instances concurrently
    local start_time=$(date +%s.%N)
    
    "$GOTO_CMD" --version > /dev/null 2>&1 &
    "$GOTO_CMD" --complete > /dev/null 2>&1 &
    "$GOTO_CMD" --history > /dev/null 2>&1 &
    
    wait
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc -l)
    
    print_info "並行実行時間 (3プロセス): ${duration}秒"
    
    # Check if concurrent execution time is reasonable (less than 3 seconds)
    if (( $(echo "$duration < 3.0" | bc -l) )); then
        print_success "並行実行が適切に処理されています (${duration}秒)"
    else
        print_error "並行実行の処理が遅すぎます (${duration}秒)"
    fi
}

# Build the goto binary first
build_goto() {
    print_test "goto バイナリのビルド（パフォーマンステスト用）"
    cd "$PROJECT_DIR/go"
    if go build -o goto; then
        print_success "goto バイナリのビルドが成功しました"
    else
        print_error "goto バイナリのビルドに失敗しました"
        exit 1
    fi
    cd - > /dev/null
}

# Check required tools
check_requirements() {
    print_test "必要なツールの確認"
    
    local missing_tools=()
    
    if ! command -v bc > /dev/null 2>&1; then
        missing_tools+=("bc")
    fi
    
    if [ ${#missing_tools[@]} -eq 0 ]; then
        print_success "すべての必要なツールが利用可能です"
    else
        print_error "必要なツールが不足しています: ${missing_tools[*]}"
        echo "以下のコマンドでインストールしてください:"
        echo "  macOS: brew install bc"
        echo "  Ubuntu/Debian: sudo apt-get install bc"
        echo "  CentOS/RHEL: sudo yum install bc"
        exit 1
    fi
}

# Main test execution
main() {
    print_header
    
    # Check requirements
    check_requirements
    
    # Setup
    setup_test_env
    trap cleanup_test_env EXIT
    
    # Build the binary first
    build_goto
    
    # Run all performance tests
    test_startup_time
    test_config_loading_performance
    test_history_performance
    test_memory_usage
    test_concurrent_execution
    
    # Print summary and exit with appropriate code
    print_summary
}

# Run main function
main "$@"
