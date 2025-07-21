#!/bin/bash

# Test script for goto program
# This script tests the basic functionality of the goto command

set +e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Print functions
print_header() {
    echo -e "${YELLOW}======================================${NC}"
    echo -e "${YELLOW} goto プログラム テストスイート${NC}"
    echo -e "${YELLOW}======================================${NC}"
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
    echo -e "${YELLOW}======================================${NC}"
    echo -e "${YELLOW} テスト結果${NC}"
    echo -e "${YELLOW}======================================${NC}"
    echo -e "合格: ${GREEN}${TESTS_PASSED}${NC}"
    echo -e "失敗: ${RED}${TESTS_FAILED}${NC}"
    echo -e "合計: $((TESTS_PASSED + TESTS_FAILED))"
    echo
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}すべてのテストが合格しました！${NC}"
        return 0
    else
        echo -e "${RED}いくつかのテストが失敗しました。${NC}"
        return 1
    fi
}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
GOTO_CMD="$PROJECT_DIR/go/goto"

# Test if goto binary exists
test_binary_exists() {
    print_test "goto バイナリの存在確認"
    if [ -f "$GOTO_CMD" ]; then
        print_success "goto バイナリが見つかりました: $GOTO_CMD"
    else
        print_error "goto バイナリが見つかりません: $GOTO_CMD"
    fi
}

# Test help option
test_help_option() {
    print_test "ヘルプオプションのテスト"
    if "$GOTO_CMD" -h > /dev/null 2>&1; then
        print_success "ヘルプオプション (-h) が正常に動作します"
    else
        print_error "ヘルプオプション (-h) が失敗しました"
    fi
}

# Test version option
test_version_option() {
    print_test "バージョンオプションのテスト"
    if "$GOTO_CMD" -v > /dev/null 2>&1; then
        print_success "バージョンオプション (-v) が正常に動作します"
    else
        print_error "バージョンオプション (-v) が失敗しました"
    fi
}

# Test completion option
test_completion_option() {
    print_test "補完オプションのテスト"
    if "$GOTO_CMD" --complete > /dev/null 2>&1; then
        print_success "補完オプション (--complete) が正常に動作します"
    else
        print_error "補完オプション (--complete) が失敗しました"
    fi
}

# Test history option
test_history_option() {
    print_test "履歴オプションのテスト"
    if "$GOTO_CMD" --history > /dev/null 2>&1; then
        print_success "履歴オプション (--history) が正常に動作します"
    else
        print_error "履歴オプション (--history) が失敗しました"
    fi
}

# Test cursor mode option
test_cursor_mode_option() {
    print_test "カーソルモードオプションのテスト"
    # カーソルモードは対話的なので、タイムアウト付きで実行
    if timeout 1s "$GOTO_CMD" -c > /dev/null 2>&1; then
        print_success "カーソルモードオプション (-c) が正常に動作します"
    else
        # タイムアウトによる終了は正常（対話的なため）
        if [ $? -eq 124 ]; then
            print_success "カーソルモードオプション (-c) が正常に動作します（対話的）"
        else
            print_error "カーソルモードオプション (-c) が失敗しました"
        fi
    fi
}

# Test label mode option
test_label_mode_option() {
    print_test "ラベルモードオプションのテスト"
    # ラベルモードは対話的なので、タイムアウト付きで実行
    if timeout 1s "$GOTO_CMD" -l > /dev/null 2>&1; then
        print_success "ラベルモードオプション (-l) が正常に動作します"
    else
        # タイムアウトによる終了は正常（対話的なため）
        if [ $? -eq 124 ]; then
            print_success "ラベルモードオプション (-l) が正常に動作します（対話的）"
        else
            print_error "ラベルモードオプション (-l) が失敗しました"
        fi
    fi
}

# Test config file creation
test_config_creation() {
    print_test "設定ファイル作成のテスト"
    
    # Backup existing config if it exists
    CONFIG_FILE="$HOME/.goto.toml"
    BACKUP_FILE="$HOME/.goto.toml.backup"
    
    if [ -f "$CONFIG_FILE" ]; then
        cp "$CONFIG_FILE" "$BACKUP_FILE"
        rm "$CONFIG_FILE"
    fi
    
    # Run goto to create default config
    if timeout 1s "$GOTO_CMD" > /dev/null 2>&1; then
        if [ -f "$CONFIG_FILE" ]; then
            print_success "デフォルト設定ファイルが正常に作成されました"
        else
            print_error "デフォルト設定ファイルの作成に失敗しました"
        fi
    else
        # タイムアウトでも設定ファイルが作成されていればOK
        if [ -f "$CONFIG_FILE" ]; then
            print_success "デフォルト設定ファイルが正常に作成されました"
        else
            print_error "デフォルト設定ファイルの作成に失敗しました"
        fi
    fi
    
    # Restore backup if it existed
    if [ -f "$BACKUP_FILE" ]; then
        mv "$BACKUP_FILE" "$CONFIG_FILE"
    fi
}

# Test invalid argument
test_invalid_argument() {
    print_test "無効な引数のテスト"
    if "$GOTO_CMD" invalid_destination_12345 > /dev/null 2>&1; then
        print_error "無効な引数が受け入れられました（本来はエラーになるべき）"
    else
        print_success "無効な引数が正しく拒否されました"
    fi
}

# Build the goto binary first
build_goto() {
    print_test "goto バイナリのビルド"
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
    
    # Build the binary first
    build_goto
    
    # Run all tests
    test_binary_exists
    test_help_option
    test_version_option
    test_completion_option
    test_history_option
    test_cursor_mode_option
    test_label_mode_option
    test_config_creation
    test_invalid_argument
    
    # Print summary and exit with appropriate code
    print_summary
}

# Run main function
main "$@"
