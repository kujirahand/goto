package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

// Constants
const (
	maxHistoryEntries = 100 // Maximum number of history entries to keep
)

// Destination represents a goto destination
type Destination struct {
	Path     string `toml:"path"`
	Shortcut string `toml:"shortcut"`
	Command  string `toml:"command"`
}

// HistoryEntry represents a history entry with timestamp
type HistoryEntry struct {
	Label    string    `json:"label"`
	LastUsed time.Time `json:"last_used"`
}

// Config represents the TOML configuration
type Config map[string]Destination

// History represents the JSON history data
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// AppConfig holds application configuration
type AppConfig struct {
	ConfigFile      string
	HistoryFile     string
	InteractiveMode string
	FilteredArgs    []string
}

func main() {
	// Initialize language support
	initializeLanguage()

	// Parse command line arguments and get configuration
	appConfig := parseCommandLineArgs()

	// Get configuration file path
	tomlFile := getConfigFilePath(appConfig.ConfigFile)

	// Load and validate configuration
	entries, shortcutMap := loadAndValidateConfig(tomlFile, appConfig.HistoryFile)

	// Handle command line arguments
	if len(appConfig.FilteredArgs) > 0 {
		handleCommandLineArguments(appConfig.FilteredArgs, entries, shortcutMap, tomlFile, appConfig.HistoryFile)
		return
	}

	// Run interactive mode
	runInteractiveMode(entries, shortcutMap, tomlFile, appConfig.InteractiveMode)
}

// initializeLanguage initializes language support
func initializeLanguage() {
	currentLanguage = detectLanguage()
	messages = getMessages(currentLanguage)
}

// parseCommandLineArgs parses command line arguments and returns configuration
func parseCommandLineArgs() AppConfig {
	config := AppConfig{
		InteractiveMode: "auto", // auto, cursor, label
	}

	if len(os.Args) <= 1 {
		return config
	}

	args := os.Args[1:]
	config.FilteredArgs = []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--config-file" && i+1 < len(args) {
			config.ConfigFile = args[i+1]
			i++ // Skip the next argument as it's the file path
		} else if arg == "--history-file" && i+1 < len(args) {
			config.HistoryFile = args[i+1]
			i++ // Skip the next argument as it's the file path
		} else if arg == "-c" {
			config.InteractiveMode = "cursor"
		} else if arg == "-l" {
			config.InteractiveMode = "label"
		} else {
			config.FilteredArgs = append(config.FilteredArgs, arg)
		}
	}

	return config
}

// getConfigFilePath returns the configuration file path
func getConfigFilePath(customConfigFile string) string {
	if customConfigFile != "" {
		return customConfigFile
	}

	usr, err := user.Current()
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
		os.Exit(1)
	}

	return filepath.Join(usr.HomeDir, ".goto.toml")
}

// loadAndValidateConfig loads and validates configuration, returns entries and shortcut map
func loadAndValidateConfig(tomlFile, customHistoryFile string) ([]Entry, map[string]int) {
	// Create default config if it doesn't exist
	if _, err := os.Stat(tomlFile); os.IsNotExist(err) {
		createDefaultConfig(tomlFile)
	}

	// Load configuration
	config, err := loadConfig(tomlFile)
	if err != nil {
		fmt.Printf("%s\n", messages.ErrorReadingConfig)
		fmt.Printf("📁 %s: %s\n", messages.ConfigFile, tomlFile)
		fmt.Printf("🔍 %s: %v\n", messages.ErrorDetails, err)
		fmt.Printf("💡 %s\n", messages.ConfigFixSuggestion)
		os.Exit(1)
	}

	// Get entries sorted by history and shortcuts
	entries := getEntriesFromConfig(config, customHistoryFile)
	shortcutMap := buildShortcutMap(entries)

	if len(entries) == 0 {
		fmt.Println(messages.NoDestinationsConfigured)
		os.Exit(1)
	}

	return entries, shortcutMap
}

// handleCommandLineArguments processes command line arguments
func handleCommandLineArguments(filteredArgs []string, entries []Entry, shortcutMap map[string]int, tomlFile, customHistoryFile string) {
	arg := filteredArgs[0]

	// Handle help option
	if arg == "-h" || arg == "--help" || arg == "help" {
		showHelp()
		os.Exit(0)
	}

	// Handle version option
	if arg == "-v" || arg == "--version" || arg == "version" {
		showVersion()
		os.Exit(0)
	}

	// Handle completion option for bash/zsh tab completion
	if arg == "--complete" {
		showCompletions(entries)
		os.Exit(0)
	}

	// Handle history option
	if arg == "--history" {
		ShowHistory(tomlFile, customHistoryFile)
		os.Exit(0)
	}

	// Handle list option
	if arg == "--list" {
		showList(entries)
		os.Exit(0)
	}

	// Handle list-label option
	if arg == "--list-label" {
		showListLabel(entries)
		os.Exit(0)
	}

	// Handle add option
	if arg == "--add" {
		success := addCurrentPathToConfig(tomlFile)
		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Find destination by argument
	handleDestinationNavigation(arg, entries, shortcutMap, tomlFile, customHistoryFile)
}

// handleDestinationNavigation handles navigation to a specific destination
func handleDestinationNavigation(arg string, entries []Entry, shortcutMap map[string]int, tomlFile, customHistoryFile string) {
	targetDir, command, label := findDestinationByArg(arg, entries, shortcutMap)

	if targetDir == "" {
		fmt.Printf(messages.DestinationNotFound, arg)
		fmt.Println("\n📋 Available destinations:")
		for _, entry := range entries {
			shortcutStr := ""
			if entry.Shortcut != "" {
				shortcutStr = fmt.Sprintf(" (%s)", entry.Shortcut)
			}
			expandedPath := expandPath(entry.Path)
			fmt.Printf("  • %s%s → %s\n", entry.Label, shortcutStr, expandedPath)
		}
		os.Exit(1)
	}

	fmt.Printf("%s %s\n", messages.FoundDestination, label)
	success := openNewShell(targetDir, command, label)
	if success {
		// Update history
		if label != "" {
			err := UpdateHistory(tomlFile, label, customHistoryFile)
			if err != nil {
				fmt.Printf("%s %v\n", messages.WarningFailedToUpdateHistory, err)
			}
		}
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// runInteractiveMode runs the interactive mode
func runInteractiveMode(entries []Entry, shortcutMap map[string]int, tomlFile, interactiveMode string) {
	targetDir, command, label := getUserChoice(entries, shortcutMap, tomlFile, interactiveMode)

	if targetDir == "ADD_CURRENT" {
		success := addCurrentPathToConfig(tomlFile)
		if success {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	if targetDir == "" {
		fmt.Println(messages.NoDirectorySelected)
		os.Exit(0)
	}

	// Update history
	if label != "" {
		err := UpdateHistory(tomlFile, label, "")
		if err != nil {
			fmt.Printf("%s %v\n", messages.WarningFailedToUpdateHistory, err)
		}
	}

	// Open the selected destination
	success := openNewShell(targetDir, command, label)
	if success {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

// Entry represents a configuration entry with label
type Entry struct {
	Label    string
	Path     string
	Shortcut string
	Command  string
}

func buildShortcutMap(entries []Entry) map[string]int {
	shortcutMap := make(map[string]int)
	for i, entry := range entries {
		if entry.Shortcut != "" {
			shortcutMap[entry.Shortcut] = i + 1
		}
	}
	return shortcutMap
}

func expandPath(path string) string {
	// URLの場合はそのまま返す
	if IsURL(path) {
		return path
	}

	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

func getUserChoice(entries []Entry, shortcutMap map[string]int, tomlFile string, interactiveMode string) (string, string, string) {
	// インタラクティブモードに基づいて分岐
	switch interactiveMode {
	case "cursor":
		return getUserChoiceCursorMode(entries, shortcutMap, tomlFile)
	case "label":
		return getUserChoiceCmdMode(entries, shortcutMap, tomlFile)
	default: // "auto"
		return getUserChoiceCursorMode(entries, shortcutMap, tomlFile) // デフォルトでカーソルモード
	}
}

// 共通のエントリー表示処理
func displayEntries(entries []Entry, selectedIndex int, cursorMode bool) {
	// ターミナル横幅取得
	termWidth := 80
	termHeight := 24
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termWidth = w
		termHeight = h
	}

	// カーソルモードの場合、画面に収まる行数を計算
	maxDisplayEntries := len(entries)
	if cursorMode {
		// ヘッダー(2行) + フッター(3行) + Exit(1行) + マージン(2行) = 8行を除く
		availableLines := termHeight - 8
		if availableLines < 3 {
			availableLines = 3 // 最低3行は確保
		}
		if len(entries) > availableLines {
			maxDisplayEntries = availableLines
		}
	}

	displayStart := 0
	displayEnd := maxDisplayEntries

	// カーソルモードで選択項目が表示範囲外の場合、表示範囲を調整
	if cursorMode && maxDisplayEntries < len(entries) {
		if selectedIndex >= maxDisplayEntries {
			// 選択項目が表示範囲の下にある場合
			displayStart = selectedIndex - maxDisplayEntries + 1
			if displayStart < 0 {
				displayStart = 0
			}
			displayEnd = displayStart + maxDisplayEntries
			if displayEnd > len(entries) {
				displayEnd = len(entries)
				displayStart = displayEnd - maxDisplayEntries
			}
		}
	}

	// エントリーの表示
	for i := displayStart; i < displayEnd; i++ {
		entry := entries[i]
		expandedPath := expandPath(entry.Path)
		shortcutStr := ""
		if entry.Shortcut != "" {
			shortcutStr = fmt.Sprintf(" (%s)", entry.Shortcut)
		}

		// 表示番号の決定（10以上は"-"で表示）
		var numStr string
		if i+1 < 10 {
			numStr = fmt.Sprintf("%d", i+1)
		} else {
			numStr = " "
		}

		// フォーマット: 数字 ラベル (ショートカットキー) → パス
		// ラベルを20文字に左寄せ
		labelWithShortcut := entry.Label + shortcutStr
		formattedLabel := fmt.Sprintf("%-20s", labelWithShortcut)
		if len([]rune(labelWithShortcut)) > 20 {
			// 20文字を超える場合は切り詰める
			runes := []rune(labelWithShortcut)
			formattedLabel = string(runes[:20])
		}

		prefix := fmt.Sprintf("%s %s → ", numStr, formattedLabel)
		maxPathLen := termWidth - len([]rune(prefix))
		pathStr := expandedPath
		if maxPathLen > 8 && len([]rune(expandedPath)) > maxPathLen {
			pathStr = shortenPathMiddle(expandedPath, maxPathLen)
		}

		// カーソルモードの場合、選択中の項目をハイライト
		if cursorMode && i == selectedIndex {
			fmt.Printf("\033[47;30m%s%s\033[0m\n", prefix, pathStr) // 白背景でハイライト
		} else {
			fmt.Printf("%s%s\n", prefix, pathStr)
		}
	}

	// 省略表示の情報
	if cursorMode && maxDisplayEntries < len(entries) {
		omittedCount := len(entries) - maxDisplayEntries
		fmt.Printf("... (%d more entries hidden)\n", omittedCount)
	}

	// print "0 Exit"
	exitPrefix := "0 Exit"
	exitIndex := len(entries) // Exitは最後のインデックス
	if cursorMode && selectedIndex == exitIndex {
		fmt.Printf("\033[47;30m%s\033[0m\n", exitPrefix) // 白背景でハイライト
	} else {
		fmt.Printf("%s\n", exitPrefix)
	}
}

// 共通の入力解析処理
func parseUserInput(choice string, entries []Entry, shortcutMap map[string]int) (string, string, string) {
	// Check if user wants to exit
	if choice == "0" || choice == "exit" || choice == "quit" {
		return "EXIT", "", ""
	}

	// Check if user wants to add current directory
	if choice == "+" {
		return "ADD_CURRENT", "", ""
	}

	// Determine input type and get corresponding entry
	index := 0

	// Check if it's a number
	if num, err := strconv.Atoi(choice); err == nil {
		index = num
	} else if shortcutIndex, exists := shortcutMap[choice]; exists {
		// Check if it's a shortcut
		index = shortcutIndex
	} else {
		// Check if it's a label name (case-insensitive)
		for i, entry := range entries {
			if strings.EqualFold(entry.Label, choice) {
				index = i + 1
				break
			}
		}
	}

	if index >= 1 && index <= len(entries) {
		entry := entries[index-1]
		expandedPath := expandPath(entry.Path)
		return expandedPath, entry.Command, entry.Label
	}

	return "", "", "" // Invalid input
}

// カーソルモードでのユーザー選択
func getUserChoiceCursorMode(entries []Entry, shortcutMap map[string]int, tomlFile string) (string, string, string) {
	selectedIndex := 0
	inputBuffer := "" // 複数文字入力用のバッファ

	// 初期表示
	redrawCursorMode(entries, selectedIndex)
	for {
		// Raw modeで入力を読み取り
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Printf("Error entering raw mode: %v\n", err)
			return "", "", ""
		}
		buffer := make([]byte, 4)
		n, err := os.Stdin.Read(buffer)
		term.Restore(int(os.Stdin.Fd()), oldState)

		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return "", "", ""
		}

		redraw := false

		// キー入力を解析
		if n == 1 {
			switch buffer[0] {
			case 13: // Enter
				if selectedIndex == len(entries) {
					// Exitが選択された場合
					fmt.Printf("\n%s\n", messages.OperationCancelled)
					return "", "", ""
				}
				entry := entries[selectedIndex]
				expandedPath := expandPath(entry.Path)
				return expandedPath, entry.Command, entry.Label
			case 27: // Escape
				// ラベル入力モードに切り替え
				return getUserChoiceCmdMode(entries, shortcutMap, tomlFile)
			case '+':
				return "ADD_CURRENT", "", ""
			case '0': // 0キーでExit
				fmt.Printf("\n%s\n", messages.OperationCancelled)
				return "", "", ""
			case '?': // ?キーでヘルプ表示
				showInteractiveHelp()
				// 画面をクリアして再表示
				redraw = true
				continue
			case 'j': // j キーで下移動 (Vim風)
				inputBuffer = "" // バッファをクリア
				if selectedIndex < len(entries) {
					selectedIndex++
					redraw = true
				}
			case 'k': // k キーで上移動 (Vim風)
				inputBuffer = "" // バッファをクリア
				if selectedIndex > 0 {
					selectedIndex--
					redraw = true
				}
			default:
				// 数字キー (0-9) またはアルファベットキーの場合
				if (buffer[0] >= '0' && buffer[0] <= '9') || (buffer[0] >= 'a' && buffer[0] <= 'z') || (buffer[0] >= 'A' && buffer[0] <= 'Z') {
					inputChar := string(buffer[0])

					// j/k は上で処理済みなのでスキップ
					if inputChar == "j" || inputChar == "k" {
						break
					}

					// 数字の場合、バッファに追加
					if buffer[0] >= '0' && buffer[0] <= '9' {
						inputBuffer += inputChar
						// 入力された数字が有効な範囲内かチェック
						if num, err := strconv.Atoi(inputBuffer); err == nil {
							if num >= 1 && num <= len(entries) {
								// 有効な番号の場合、即座に決定
								entry := entries[num-1]
								expandedPath := expandPath(entry.Path)
								return expandedPath, entry.Command, entry.Label
							} else if num > len(entries) {
								// 範囲外の場合、バッファをクリア
								inputBuffer = ""
							}
						}
					} else {
						// ショートカットキーの場合、バッファをクリアして即座に実行
						inputBuffer = ""
						if shortcutIndex, exists := shortcutMap[inputChar]; exists {
							entry := entries[shortcutIndex-1]
							expandedPath := expandPath(entry.Path)
							return expandedPath, entry.Command, entry.Label
						}
					}
				} else {
					// その他のキーが押された場合、バッファをクリア
					inputBuffer = ""
				}
			}
		} else if n >= 3 && buffer[0] == 27 && buffer[1] == '[' {
			switch buffer[2] {
			case 'A': // Up arrow
				inputBuffer = "" // バッファをクリア
				if selectedIndex > 0 {
					selectedIndex--
					redraw = true
				}
			case 'B': // Down arrow
				inputBuffer = "" // バッファをクリア
				if selectedIndex < len(entries) {
					selectedIndex++
					redraw = true
				}
			}
		}

		// 画面の再描画
		if redraw {
			redrawCursorMode(entries, selectedIndex)
		}
	}
}

// カーソルモードの画面再描画
func redrawCursorMode(entries []Entry, selectedIndex int) {
	// より効率的な再描画: 画面全体をクリアしてから再描画
	fmt.Print("\033[2J\033[H") // 画面をクリアしてカーソルを左上に移動

	// ヘッダーを表示
	PrintWhiteBackgroundLine(messages.AvailableDestinations)
	fmt.Println() // ヘッダーの後に改行

	// エントリーリストを再表示
	displayEntries(entries, selectedIndex, true)

	// 横線を表示
	PrintHorzontalLine("-")
	// フッターメッセージを更新
	fmt.Println(messages.InteractiveHelp)
	fmt.Printf("%s\n", messages.CursorModeHint)
}

// コマンド（ラベル）入力モードでのユーザー選択
func getUserChoiceCmdMode(entries []Entry, shortcutMap map[string]int, tomlFile string) (string, string, string) {
	for {
		// 画面をクリア
		fmt.Print("\033[2J\033[H")
		PrintWhiteBackgroundLine(messages.AvailableDestinations)
		fmt.Println()
		displayEntries(entries, 0, false)
		PrintWhiteBackgroundLine(messages.InteractiveHelp)
		fmt.Println()
		fmt.Printf("%s\n", messages.EnterChoice)
		fmt.Printf("%s\n", messages.BackToCursorModeHint)
		fmt.Printf("%s ", messages.EnterChoicePrompt)

		// 通常の入力モード
		reader := bufio.NewReader(os.Stdin)
		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n%s\n", messages.OperationCancelled)
			return "", "", ""
		}

		choice = strings.TrimSpace(choice)

		// 空の入力の場合、カーソルモードに切り替え
		if choice == "" {
			return getUserChoiceCursorMode(entries, shortcutMap, tomlFile)
		}

		// Check if user wants to show help
		if choice == "?" {
			showInteractiveHelp()
			continue
		}

		// 入力を解析
		targetDir, command, label := parseUserInput(choice, entries, shortcutMap)

		// Exit選択の場合
		if targetDir == "EXIT" {
			fmt.Printf("\n%s\n", messages.OperationCancelled)
			return "", "", ""
		}

		// ADD_CURRENT選択の場合
		if targetDir == "ADD_CURRENT" {
			return "ADD_CURRENT", "", ""
		}

		// 無効な入力の場合
		if targetDir == "" && label == "" && command == "" {
			fmt.Println(messages.InvalidInput)
			continue
		}

		return targetDir, command, label
	}
}

func openNewShell(targetDir, command, label string) bool {
	// URLの場合はブラウザで開く
	if IsURL(targetDir) {
		fmt.Printf("%s %s\n", messages.OpeningShell, targetDir)
		if label != "" {
			fmt.Printf("%s %s\n", messages.Destination, label)
		}

		err := OpenURL(targetDir)
		if err != nil {
			fmt.Printf("Error opening URL: %v\n", err)
			return false
		}

		fmt.Printf("✅ Opened URL in default browser: %s\n", targetDir)
		return true
	}

	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("%s %s\n", messages.DirectoryNotExist, targetDir)
		return false
	}

	openShellMessage := fmt.Sprintf("%s %s", messages.OpeningShell, targetDir)
	PrintWhiteBackgroundLine(openShellMessage)
	fmt.Println()
	if label != "" {
		fmt.Printf("%s %s\n", messages.Destination, label)
	}

	// Get user's preferred shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	if command != "" {
		fmt.Printf("%s %s\n", messages.WillExecute, command)
		fmt.Println(strings.Repeat("=", 50))

		// Create a temporary startup script
		tempScript := createTempScript(targetDir, command, shell)
		defer os.Remove(tempScript)

		// Execute the temporary script
		cmd := exec.Command("/bin/sh", tempScript)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = targetDir

		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s %v\n", messages.ErrorOpeningShell, err)
			return false
		}
	} else {
		// Simply open shell in the target directory
		fmt.Println(messages.TypeExitToReturn)
		fmt.Println(strings.Repeat("=", 50))

		fmt.Printf("%s %s\n", messages.YouAreNowIn, targetDir)

		// Start new shell with the target directory as working directory
		cmd := exec.Command(shell)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = targetDir // Set working directory for the new shell

		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s %v\n", messages.ErrorOpeningShell, err)
			return false
		}
	}

	return true
}

func createTempScript(targetDir, command, shell string) string {
	tempFile, err := os.CreateTemp("", "goto_*.sh")
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorCreatingTempFile, err)
		os.Exit(1)
	}

	scriptContent := fmt.Sprintf(`#!/bin/sh
cd "%s"
echo "%s $(pwd)"
echo "%s %s"
echo "%s"
%s
echo "%s"
echo "%s $(pwd)"
echo "%s"
exec "%s"
`, targetDir, messages.CurrentDirectory, messages.ExecutingCommand, command,
		strings.Repeat("-", 40), command, strings.Repeat("-", 40),
		messages.CommandCompleted, messages.TypeExitToReturn, shell)

	_, err = tempFile.WriteString(scriptContent)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorWritingTempScript, err)
		os.Exit(1)
	}

	tempFile.Close()

	// Make executable
	err = os.Chmod(tempFile.Name(), 0755)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorMakingExecutable, err)
		os.Exit(1)
	}

	return tempFile.Name()
}

func addCurrentPathToConfig(tomlFile string) bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorGettingCurrentDir, err)
		return false
	}

	fmt.Printf("%s %s\n", messages.CurrentDirectory, currentDir)

	// 既存の設定を読み込んでショートカットマップを作成
	config, err := loadConfig(tomlFile)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorReadingConfig, err)
		return false
	}

	entries := getEntriesFromConfig(config, "")
	shortcutMap := buildShortcutMap(entries)

	// フォルダ名をデフォルトラベルとして取得
	defaultLabel := filepath.Base(currentDir)

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [%s]: ", messages.EnterLabel, defaultLabel)
	label, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n%s\n", messages.OperationCancelled)
		return false
	}

	label = strings.TrimSpace(label)
	if label == "" {
		label = defaultLabel // ラベルが空の場合、フォルダ名をデフォルトとして使用
	}

	var shortcut string
	for {
		fmt.Printf("%s ", messages.EnterShortcutOptional)
		shortcutInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n%s\n", messages.OperationCancelled)
			return false
		}

		shortcut = strings.TrimSpace(shortcutInput)

		// ショートカットが空の場合は問題なし
		if shortcut == "" {
			break
		}

		// 既存のショートカットと重複していないかチェック
		if _, exists := shortcutMap[shortcut]; exists {
			fmt.Printf(messages.ShortcutAlreadyExists, shortcut)
			continue
		}

		// 重複していなければループを抜ける
		break
	}

	// Generate TOML entry
	tomlEntry := fmt.Sprintf("\n[%s]\n", label)
	tomlEntry += fmt.Sprintf("path = \"%s\"\n", currentDir)
	if shortcut != "" {
		tomlEntry += fmt.Sprintf("shortcut = \"%s\"\n", shortcut)
	}

	// Append to TOML file
	file, err := os.OpenFile(tomlFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorOpeningConfigFile, err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(tomlEntry)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorWritingConfigFile, err)
		return false
	}

	fmt.Printf("%s '%s' → %s\n", messages.Added, label, currentDir)
	if shortcut != "" {
		fmt.Printf("%s %s\n", messages.Shortcut, shortcut)
	}

	return true
}

func findDestinationByArg(arg string, entries []Entry, shortcutMap map[string]int) (string, string, string) {
	// Check if it's a number
	if num, err := strconv.Atoi(arg); err == nil {
		if num >= 1 && num <= len(entries) {
			entry := entries[num-1]
			expandedPath := expandPath(entry.Path)
			return expandedPath, entry.Command, entry.Label
		}
		return "", "", ""
	}

	// Check if it's a shortcut
	if index, exists := shortcutMap[arg]; exists {
		entry := entries[index-1]
		expandedPath := expandPath(entry.Path)
		return expandedPath, entry.Command, entry.Label
	}

	// Check if it's a label (case-insensitive)
	for _, entry := range entries {
		if strings.EqualFold(entry.Label, arg) {
			expandedPath := expandPath(entry.Path)
			return expandedPath, entry.Command, entry.Label
		}
	}

	return "", "", ""
}

func showVersion() {
	fmt.Printf("%s version %s\n", AppName, Version)
}

func showHelp() {
	// Get configuration file path for display
	usr, err := user.Current()
	configPath := "~/.goto.toml"
	if err == nil {
		configPath = filepath.Join(usr.HomeDir, ".goto.toml")
	}

	fmt.Println(messages.NavigateDirectoriesQuickly)
	fmt.Printf("\n%s %s\n", messages.ConfigurationFile, configPath)
	fmt.Printf("\n%s\n", messages.Usage)
	fmt.Printf("  goto                 %s\n", messages.ShowInteractiveMenu)
	fmt.Printf("  goto -c              %s\n", "カーソル移動モードでインタラクティブメニューを表示")
	fmt.Printf("  goto -l              %s\n", "ラベル入力モードでインタラクティブメニューを表示")
	fmt.Printf("  goto --config-file FILE %s\n", "指定した設定ファイルを使用")
	fmt.Printf("  goto --history-file FILE %s\n", "指定した履歴ファイルを使用")
	fmt.Printf("  goto <number>        %s\n", messages.GoToDestinationByNumber)
	fmt.Printf("  goto <label>         %s\n", messages.GoToDestinationByLabel)
	fmt.Printf("  goto <shortcut>      %s\n", messages.GoToDestinationByShortcut)
	fmt.Printf("  goto -h, --help      %s\n", messages.ShowHelpMessage)
	fmt.Printf("  goto -v, --version   %s\n", messages.ShowVersionInfo)
	fmt.Printf("  goto --complete      %s\n", messages.ShowCompletionCandidates)
	fmt.Printf("  goto --history       %s\n", messages.ShowRecentUsageHistory)
	fmt.Printf("  goto --list          %s\n", "履歴順でディレクトリ一覧を表示")
	fmt.Printf("  goto --list-label    %s\n", "履歴順でラベル一覧を表示")
	fmt.Printf("  goto --add           %s\n", messages.AddCurrentDirectoryToConfig)
	fmt.Printf("\n%s\n", messages.Examples)
	fmt.Printf("  goto 1              %s\n", messages.NavigateToFirstDest)
	fmt.Printf("  goto Home           %s\n", messages.NavigateToHomeDest)
	fmt.Printf("  goto h              %s\n", messages.NavigateUsingShortcut)
	fmt.Printf("  goto                %s\n", messages.ShowInteractiveMenuExample)
}

// showList displays all destinations sorted by history
func showList(entries []Entry) {
	for i, entry := range entries {
		// Format: number. label (shortcut) → path
		shortcutStr := ""
		if entry.Shortcut != "" {
			shortcutStr = fmt.Sprintf(" (%s)", entry.Shortcut)
		}

		expandedPath := expandPath(entry.Path)
		fmt.Printf("%2d. %s%s → %s\n", i+1, entry.Label, shortcutStr, expandedPath)
	}
}

// showListLabel displays only labels sorted by history
func showListLabel(entries []Entry) {
	for _, entry := range entries {
		fmt.Println(entry.Label)
	}
}

func showCompletions(entries []Entry) {
	// Output only labels for completion
	for _, entry := range entries {
		fmt.Println(entry.Label)
	}
}
