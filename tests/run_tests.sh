#!/bin/bash

# Main test runner for goto program
# This script runs all test suites and provides a comprehensive test report

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Test suite counter
TOTAL_SUITES=0
PASSED_SUITES=0
FAILED_SUITES=0

print_main_header() {
    echo -e "${WHITE}================================================${NC}"
    echo -e "${WHITE} goto プログラム 統合テストスイート${NC}"
    echo -e "${WHITE}================================================${NC}"
    echo
    echo -e "${CYAN}このテストスイートは以下のテストを実行します:${NC}"
    echo -e "  1. ${YELLOW}基本機能テスト${NC} - コマンドラインオプションと基本動作"
    echo -e "  2. ${BLUE}機能テスト${NC} - 設定解析、履歴処理、パス展開など"
    echo -e "  3. ${PURPLE}パフォーマンステスト${NC} - 起動時間、メモリ使用量など"
    echo
}

print_suite_header() {
    local suite_name="$1"
    echo -e "${WHITE}================================================${NC}"
    echo -e "${WHITE} $suite_name 実行中...${NC}"
    echo -e "${WHITE}================================================${NC}"
}

print_suite_result() {
    local suite_name="$1"
    local result="$2"
    
    ((TOTAL_SUITES++))
    
    if [ "$result" -eq 0 ]; then
        echo -e "${GREEN}✅ $suite_name が成功しました${NC}"
        ((PASSED_SUITES++))
    else
        echo -e "${RED}❌ $suite_name が失敗しました${NC}"
        ((FAILED_SUITES++))
    fi
    echo
}

print_final_summary() {
    echo -e "${WHITE}================================================${NC}"
    echo -e "${WHITE} 最終テスト結果${NC}"
    echo -e "${WHITE}================================================${NC}"
    echo -e "実行されたテストスイート: $TOTAL_SUITES"
    echo -e "成功したテストスイート: ${GREEN}$PASSED_SUITES${NC}"
    echo -e "失敗したテストスイート: ${RED}$FAILED_SUITES${NC}"
    echo
    
    if [ $FAILED_SUITES -eq 0 ]; then
        echo -e "${GREEN}🎉 すべてのテストスイートが成功しました！${NC}"
        echo -e "${GREEN}goto プログラムは正常に動作しています。${NC}"
        return 0
    else
        echo -e "${RED}⚠️  一部のテストスイートが失敗しました。${NC}"
        echo -e "${YELLOW}詳細については上記のテスト結果を確認してください。${NC}"
        return 1
    fi
}

run_test_suite() {
    local script_name="$1"
    local suite_name="$2"
    local script_path="$SCRIPT_DIR/$script_name"
    
    if [ ! -f "$script_path" ]; then
        echo -e "${RED}エラー: テストスクリプト '$script_path' が見つかりません${NC}"
        return 1
    fi
    
    if [ ! -x "$script_path" ]; then
        chmod +x "$script_path"
    fi
    
    print_suite_header "$suite_name"
    
    if "$script_path"; then
        print_suite_result "$suite_name" 0
        return 0
    else
        print_suite_result "$suite_name" 1
        return 1
    fi
}

# Check if Go is installed
check_go_installation() {
    echo -e "${CYAN}Go インストールの確認...${NC}"
    if ! command -v go > /dev/null 2>&1; then
        echo -e "${RED}エラー: Go がインストールされていません${NC}"
        echo "Go をインストールしてからテストを実行してください"
        echo "https://golang.org/dl/"
        exit 1
    fi
    
    local go_version=$(go version)
    echo -e "${GREEN}✅ Go がインストールされています: $go_version${NC}"
    echo
}

# Parse command line arguments
parse_arguments() {
    local run_basic=true
    local run_functional=true
    local run_performance=true
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --basic-only)
                run_functional=false
                run_performance=false
                shift
                ;;
            --functional-only)
                run_basic=false
                run_performance=false
                shift
                ;;
            --performance-only)
                run_basic=false
                run_functional=false
                shift
                ;;
            --no-performance)
                run_performance=false
                shift
                ;;
            -h|--help)
                echo "使用方法: $0 [オプション]"
                echo
                echo "オプション:"
                echo "  --basic-only        基本機能テストのみ実行"
                echo "  --functional-only   機能テストのみ実行"
                echo "  --performance-only  パフォーマンステストのみ実行"
                echo "  --no-performance    パフォーマンステスト以外を実行"
                echo "  -h, --help          このヘルプメッセージを表示"
                echo
                exit 0
                ;;
            *)
                echo -e "${RED}エラー: 不明なオプション '$1'${NC}"
                echo "使用方法については '$0 --help' を参照してください"
                exit 1
                ;;
        esac
    done
    
    echo "$run_basic $run_functional $run_performance"
}

main() {
    # Parse command line arguments
    local args=$(parse_arguments "$@")
    local run_basic=$(echo $args | cut -d' ' -f1)
    local run_functional=$(echo $args | cut -d' ' -f2)
    local run_performance=$(echo $args | cut -d' ' -f3)
    
    print_main_header
    check_go_installation
    
    local overall_result=0
    
    # Run basic tests
    if [ "$run_basic" = "true" ]; then
        if ! run_test_suite "test.sh" "基本機能テスト"; then
            overall_result=1
        fi
    fi
    
    # Run functional tests
    if [ "$run_functional" = "true" ]; then
        if ! run_test_suite "functional_test.sh" "機能テスト"; then
            overall_result=1
        fi
    fi
    
    # Run performance tests
    if [ "$run_performance" = "true" ]; then
        if ! run_test_suite "performance_test.sh" "パフォーマンステスト"; then
            overall_result=1
        fi
    fi
    
    # Print final summary
    if ! print_final_summary; then
        overall_result=1
    fi
    
    exit $overall_result
}

# Run main function with all arguments
main "$@"
