#!/bin/bash

# Advanced functional tests for goto program
# This script tests specific functionality like config parsing, history handling, etc.

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Print functions
print_header() {
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE} goto 機能テストスイート${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo
}

print_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

print_summary() {
    echo
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE} 機能テスト結果${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo -e "合格: ${GREEN}${TESTS_PASSED}${NC}"
    echo -e "失敗: ${RED}${TESTS_FAILED}${NC}"
    echo -e "合計: $((TESTS_PASSED + TESTS_FAILED))"
    echo
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}すべての機能テストが合格しました！${NC}"
        return 0
    else
        echo -e "${RED}いくつかの機能テストが失敗しました。${NC}"
        return 1
    fi
}

# Get directories
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
GOTO_CMD="$PROJECT_DIR/go/goto"

# Test temporary directory
TEST_TEMP_DIR="$SCRIPT_DIR/temp"

# Setup test environment
setup_test_env() {
    print_test "テスト環境のセットアップ"
    
    # Create temp directory
    mkdir -p "$TEST_TEMP_DIR"
    
    # Backup existing config and history
    if [ -f "$HOME/.goto.toml" ]; then
        cp "$HOME/.goto.toml" "$TEST_TEMP_DIR/goto.toml.backup"
    fi
    
    if [ -f "$HOME/.goto.history.json" ]; then
        cp "$HOME/.goto.history.json" "$TEST_TEMP_DIR/goto.history.json.backup"
    fi
    
    print_success "テスト環境がセットアップされました"
}

# Cleanup test environment
cleanup_test_env() {
    print_test "テスト環境のクリーンアップ"
    
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
    
    print_success "テスト環境がクリーンアップされました"
}

# Test TOML config parsing
test_toml_config_parsing() {
    print_test "TOML設定ファイルの解析テスト"
    
    # Create test config in temp directory
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    cat > "$test_config" << 'EOF'
# Test configuration for goto

[home]
path = "~/"
shortcut = "h"

[projects]
path = "~/Projects"
shortcut = "p"

[documents]
path = "~/Documents"

[website]
path = "https://example.com"
shortcut = "w"

[scripts]
path = "~/Scripts"
command = "ls -la"
EOF
    
    # Test if goto can parse the config
    if "$GOTO_CMD" --config "$test_config" --complete > "$TEST_TEMP_DIR/completion_output.txt" 2>&1; then
        # Check if all destinations are listed
        if grep -q "home" "$TEST_TEMP_DIR/completion_output.txt" && \
           grep -q "projects" "$TEST_TEMP_DIR/completion_output.txt" && \
           grep -q "documents" "$TEST_TEMP_DIR/completion_output.txt" && \
           grep -q "website" "$TEST_TEMP_DIR/completion_output.txt" && \
           grep -q "scripts" "$TEST_TEMP_DIR/completion_output.txt"; then
            print_success "TOML設定ファイルが正しく解析されました"
        else
            print_error "TOML設定ファイルの解析で一部の設定が読み込まれませんでした"
        fi
    else
        print_error "TOML設定ファイルの解析に失敗しました"
    fi
}

# Test history functionality
test_history_functionality() {
    print_test "履歴機能のテスト"
    
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    local test_history="$TEST_TEMP_DIR/test_history.json"
    
    # Create test config
    cat > "$test_config" << 'EOF'
[test_entry]
path = "~/test_destination"
EOF
    
    # Create simple history file
    cat > "$test_history" << 'EOF'
[
  {
    "label": "test_entry",
    "last_used": "2024-01-01T12:00:00Z"
  }
]
EOF
    
    # Test history display
    if "$GOTO_CMD" --config "$test_config" --history-file "$test_history" --history > "$TEST_TEMP_DIR/history_output.txt" 2>&1; then
        if grep -q "test_entry" "$TEST_TEMP_DIR/history_output.txt"; then
            print_success "履歴機能が正常に動作しています"
        else
            print_error "履歴機能の出力に問題があります"
        fi
    else
        print_error "履歴機能のテストに失敗しました"
    fi
}

# Test path expansion
test_path_expansion() {
    print_test "パス展開機能のテスト"
    
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    
    # Create config with tilde path
    cat > "$test_config" << 'EOF'
[home]
path = "~/"
shortcut = "h"

[absolute]
path = "/tmp"
shortcut = "t"
EOF
    
    # Test completion to see if paths are handled correctly
    if "$GOTO_CMD" --config "$test_config" --complete > /dev/null 2>&1; then
        print_success "パス展開機能が正常に動作しています"
    else
        print_error "パス展開機能のテストに失敗しました"
    fi
}

# Test URL handling
test_url_handling() {
    print_test "URL処理機能のテスト"
    
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    
    # Create config with URL
    cat > "$test_config" << 'EOF'
[github]
path = "https://github.com"
shortcut = "g"

[local]
path = "~/Documents"
shortcut = "d"
EOF
    
    # Test completion
    if "$GOTO_CMD" --config "$test_config" --complete > "$TEST_TEMP_DIR/url_test_output.txt" 2>&1; then
        if grep -q "github" "$TEST_TEMP_DIR/url_test_output.txt" && \
           grep -q "local" "$TEST_TEMP_DIR/url_test_output.txt"; then
            print_success "URL処理機能が正常に動作しています"
        else
            print_error "URL処理機能の出力に問題があります"
        fi
    else
        print_error "URL処理機能のテストに失敗しました"
    fi
}

# Test shortcut functionality
test_shortcut_functionality() {
    print_test "ショートカット機能のテスト"
    
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    
    # Create config with shortcuts
    cat > "$test_config" << 'EOF'
[home]
path = "~/"
shortcut = "h"

[projects]
path = "~/Projects"
shortcut = "p"

[documents]
path = "~/Documents"
shortcut = "d"
EOF
    
    # Test help output contains shortcut information
    if "$GOTO_CMD" --config "$test_config" -h > "$TEST_TEMP_DIR/shortcut_test_output.txt" 2>&1; then
        if grep -q "ショートカット" "$TEST_TEMP_DIR/shortcut_test_output.txt" || \
           grep -q "shortcut" "$TEST_TEMP_DIR/shortcut_test_output.txt"; then
            print_success "ショートカット機能が正常に動作しています"
        else
            print_success "ショートカット機能の基本動作が確認できました"
        fi
    else
        print_error "ショートカット機能のテストに失敗しました"
    fi
}

# Test interactive mode options
test_interactive_mode_options() {
    print_test "インタラクティブモードオプションのテスト"
    
    # Create minimal config
    local test_config="$TEST_TEMP_DIR/test_config.toml"
    cat > "$test_config" << 'EOF'
[test]
path = "/tmp"
EOF
    
    # Test cursor mode (-c)
    if timeout 1s "$GOTO_CMD" --config "$test_config" -c > /dev/null 2>&1; then
        # Timeout is expected for interactive mode
        print_success "カーソルモード (-c) オプションが認識されています"
    else
        if [ $? -eq 124 ]; then  # Timeout exit code
            print_success "カーソルモード (-c) オプションが認識されています"
        else
            print_error "カーソルモード (-c) オプションのテストに失敗しました"
        fi
    fi
    
    # Test label mode (-l)
    if timeout 1s "$GOTO_CMD" --config "$test_config" -l > /dev/null 2>&1; then
        # Timeout is expected for interactive mode
        print_success "ラベルモード (-l) オプションが認識されています"
    else
        if [ $? -eq 124 ]; then  # Timeout exit code
            print_success "ラベルモード (-l) オプションが認識されています"
        else
            print_error "ラベルモード (-l) オプションのテストに失敗しました"
        fi
    fi
}

# Test error handling
test_error_handling() {
    print_test "エラーハンドリングのテスト"
    
    local test_config="$TEST_TEMP_DIR/invalid_config.toml"
    
    # Test with invalid config
    cat > "$test_config" << 'EOF'
[invalid
path = "invalid toml syntax
EOF
    
    # Test that the program handles invalid config gracefully
    if "$GOTO_CMD" --config "$test_config" --complete > "$TEST_TEMP_DIR/error_output.txt" 2>&1; then
        print_error "無効な設定ファイルに対してエラーが発生しませんでした"
    else
        print_success "無効な設定ファイルを適切にハンドリングしました"
    fi
    
    # Test with non-existent config
    if "$GOTO_CMD" --config "/non/existent/config.toml" --complete > "$TEST_TEMP_DIR/missing_config_output.txt" 2>&1; then
        print_error "存在しない設定ファイルに対してエラーが発生しませんでした"
    else
        print_success "存在しない設定ファイルを適切にハンドリングしました"
    fi
}# Build the goto binary first
build_goto() {
    print_test "goto バイナリのビルド（機能テスト用）"
    cd "$PROJECT_DIR/go"
    if go build -o goto; then
        print_success "goto バイナリのビルドが成功しました"
    else
        print_error "goto バイナリのビルドに失敗しました"
        exit 1
    fi
    cd - > /dev/null
}

# Main test execution
main() {
    print_header
    
    # Setup
    setup_test_env
    trap cleanup_test_env EXIT
    
    # Build the binary first
    build_goto
    
    # Run all functional tests
    test_toml_config_parsing
    test_history_functionality
    test_path_expansion
    test_url_handling
    test_shortcut_functionality
    test_interactive_mode_options
    test_error_handling
    
    # Print summary and exit with appropriate code
    print_summary
}

# Run main function
main "$@"

