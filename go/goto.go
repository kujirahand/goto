package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/BurntSushi/toml"
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
	Label    string    `toml:"label"`
	LastUsed time.Time `toml:"last_used"`
}

// ConfigWithHistory represents the TOML configuration with history
type ConfigWithHistory struct {
	Destinations map[string]Destination `toml:",inline"`
	History      []HistoryEntry         `toml:"history"`
}

// Config represents the TOML configuration
type Config map[string]Destination

// History represents the JSON history data
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// Get history file path
func getHistoryFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".goto.history.json"), nil
}

func main() {
	// Initialize language support
	currentLanguage = detectLanguage()
	messages = getMessages(currentLanguage)

	// Get configuration file path
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
		os.Exit(1)
	}

	tomlFile := filepath.Join(usr.HomeDir, ".goto.toml")
	historyFile, err := getHistoryFilePath()
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
		os.Exit(1)
	}

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

	// Load history
	history, err := loadHistory(historyFile)
	if err != nil {
		// If history file doesn't exist or has an error, create an empty history
		history = History{Entries: []HistoryEntry{}}
	}

	// Create ConfigWithHistory for backward compatibility
	configWithHistory := ConfigWithHistory{
		Destinations: config,
		History:      history.Entries,
	}

	// Get entries sorted by history and shortcuts
	entries := getEntriesWithHistory(configWithHistory)
	shortcutMap := buildShortcutMap(entries)

	if len(entries) == 0 {
		fmt.Println(messages.NoDestinationsConfigured)
		os.Exit(1)
	}

	// Handle command line arguments
	if len(os.Args) > 1 {
		arg := os.Args[1]

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
			showHistory(configWithHistory)
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
				err := updateHistory(tomlFile, label)
				if err != nil {
					fmt.Printf("%s %v\n", messages.WarningFailedToUpdateHistory, err)
				}
			}
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	// Interactive mode
	targetDir, command, label := getUserChoice(entries, shortcutMap, tomlFile)

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

	success := openNewShell(targetDir, command, label)
	if success {
		// Update history
		if label != "" {
			err := updateHistory(tomlFile, label)
			if err != nil {
				fmt.Printf("%s %v\n", messages.WarningFailedToUpdateHistory, err)
			}
		}
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

func createDefaultConfig(tomlFile string) {
	err := os.WriteFile(tomlFile, []byte(DefaultConfig), 0644)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorWritingConfigFile, err)
		os.Exit(1)
	}
	fmt.Printf("%s %s\n", messages.CreatedDefaultConfig, tomlFile)
}

func loadConfig(tomlFile string) (map[string]Destination, error) {
	var config map[string]Destination
	_, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func loadHistory(historyFile string) (History, error) {
	var history History

	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return History{Entries: []HistoryEntry{}}, nil
	}

	// Read and parse history file
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return History{Entries: []HistoryEntry{}}, err
	}

	err = json.Unmarshal(data, &history)
	if err != nil {
		return History{Entries: []HistoryEntry{}}, err
	}

	// Limit history to the latest 100 entries when loading
	if len(history.Entries) > maxHistoryEntries {
		// Sort by most recent first
		sort.Slice(history.Entries, func(i, j int) bool {
			return history.Entries[i].LastUsed.After(history.Entries[j].LastUsed)
		})

		// Keep only the latest 100 entries
		history.Entries = history.Entries[:maxHistoryEntries]

		// Save the trimmed history back to file
		go func() {
			// Save asynchronously to avoid blocking the main operation
			saveHistory(historyFile, history)
		}()
	}

	return history, nil
}

func saveHistory(historyFile string, history History) error {
	// Limit history to the latest entries
	if len(history.Entries) > maxHistoryEntries {
		// Sort by most recent first
		sort.Slice(history.Entries, func(i, j int) bool {
			return history.Entries[i].LastUsed.After(history.Entries[j].LastUsed)
		})

		// Keep only the latest 100 entries
		history.Entries = history.Entries[:maxHistoryEntries]
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

func loadConfigWithHistory(tomlFile string) (ConfigWithHistory, error) {
	var rawConfig map[string]interface{}
	_, err := toml.DecodeFile(tomlFile, &rawConfig)
	if err != nil {
		return ConfigWithHistory{}, err
	}

	config := ConfigWithHistory{
		Destinations: make(map[string]Destination),
		History:      []HistoryEntry{},
	}

	// Parse destinations
	for key, value := range rawConfig {
		if key == "history" {
			// Parse history section
			if historyData, ok := value.([]map[string]interface{}); ok {
				for _, histItem := range historyData {
					if label, ok := histItem["label"].(string); ok {
						if lastUsedStr, ok := histItem["last_used"].(string); ok {
							if lastUsed, err := time.Parse(time.RFC3339, lastUsedStr); err == nil {
								config.History = append(config.History, HistoryEntry{
									Label:    label,
									LastUsed: lastUsed,
								})
							}
						}
					}
				}
			}
		} else {
			// Parse destination
			if destData, ok := value.(map[string]interface{}); ok {
				dest := Destination{}
				if path, ok := destData["path"].(string); ok {
					dest.Path = path
				}
				if shortcut, ok := destData["shortcut"].(string); ok {
					dest.Shortcut = shortcut
				}
				if command, ok := destData["command"].(string); ok {
					dest.Command = command
				}
				config.Destinations[key] = dest
			}
		}
	}

	return config, nil
}

func getEntriesWithHistory(config ConfigWithHistory) []Entry {
	var entries []Entry

	// Load history from separate JSON file
	historyFile, err := getHistoryFilePath()
	if err != nil {
		// If we can't get history file path, proceed without history sorting
		for label, dest := range config.Destinations {
			entries = append(entries, Entry{
				Label:    label,
				Path:     dest.Path,
				Shortcut: dest.Shortcut,
				Command:  dest.Command,
			})
		}
		// Sort alphabetically if no history available
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Label < entries[j].Label
		})
		return entries
	}

	history, err := loadHistory(historyFile)
	if err != nil {
		// If history file doesn't exist or has error, proceed without history sorting
		for label, dest := range config.Destinations {
			entries = append(entries, Entry{
				Label:    label,
				Path:     dest.Path,
				Shortcut: dest.Shortcut,
				Command:  dest.Command,
			})
		}
		// Sort alphabetically if no history available
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Label < entries[j].Label
		})
		return entries
	}

	// Create a map for quick lookup of history
	historyMap := make(map[string]time.Time)
	for _, hist := range history.Entries {
		historyMap[hist.Label] = hist.LastUsed
	}

	// Collect all entries
	for label, dest := range config.Destinations {
		entries = append(entries, Entry{
			Label:    label,
			Path:     dest.Path,
			Shortcut: dest.Shortcut,
			Command:  dest.Command,
		})
	}

	// Sort entries by history (most recent first)
	sort.Slice(entries, func(i, j int) bool {
		timeI, hasI := historyMap[entries[i].Label]
		timeJ, hasJ := historyMap[entries[j].Label]

		// If both have history, sort by time (most recent first)
		if hasI && hasJ {
			return timeI.After(timeJ)
		}

		// If only one has history, prioritize it
		if hasI && !hasJ {
			return true
		}
		if !hasI && hasJ {
			return false
		}

		// If neither has history, sort alphabetically
		return entries[i].Label < entries[j].Label
	})

	return entries
}

func saveConfigWithHistory(tomlFile string, config ConfigWithHistory) error {
	// Create a map that includes both destinations and history
	configMap := make(map[string]interface{})

	// Add destinations
	for label, dest := range config.Destinations {
		configMap[label] = dest
	}

	// Add history if it exists
	if len(config.History) > 0 {
		historyEntries := make([]map[string]interface{}, len(config.History))
		for i, hist := range config.History {
			historyEntries[i] = map[string]interface{}{
				"label":     hist.Label,
				"last_used": hist.LastUsed.Format(time.RFC3339),
			}
		}
		configMap["history"] = historyEntries
	}

	// Convert to TOML
	var buf strings.Builder
	encoder := toml.NewEncoder(&buf)
	err := encoder.Encode(configMap)
	if err != nil {
		return err
	}

	return os.WriteFile(tomlFile, []byte(buf.String()), 0644)
}

func updateHistory(tomlFile string, label string) error {
	// Get history file path
	historyFile, err := getHistoryFilePath()
	if err != nil {
		return err
	}

	// Load history
	history, err := loadHistory(historyFile)
	if err != nil {
		// If error loading history, create a new one
		history = History{Entries: []HistoryEntry{}}
	}

	// Update or add history entry
	now := time.Now()
	found := false

	for i, hist := range history.Entries {
		if hist.Label == label {
			history.Entries[i].LastUsed = now
			found = true
			break
		}
	}

	if !found {
		history.Entries = append(history.Entries, HistoryEntry{
			Label:    label,
			LastUsed: now,
		})
	}

	// Save updated history
	return saveHistory(historyFile, history)
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
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return path
		}
		return filepath.Join(usr.HomeDir, path[2:])
	}
	return path
}

// パスの途中を省略して...で表示する
func shortenPathMiddle(path string, maxLen int) string {
	r := []rune(path)
	if len(r) <= maxLen {
		return path
	}
	// 先頭3文字 + ... + 末尾(maxLen-6)文字
	keep := maxLen - 3
	if keep < 6 {
		// 省略しすぎないように
		return string(r[:maxLen])
	}
	head := keep / 2
	tail := keep - head
	return string(r[:head]) + "..." + string(r[len(r)-tail:])
}

// PrintWhiteBackgroundLine prints a line with white background
func PrintWhiteBackgroundLine(text string) {
	// ターミナル横幅取得
	termWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termWidth = w
	}

	// 一行全部を白背景にして表示
	BOW := "\033[47;30m"
	EOW := "\033[0m"
	fmt.Printf(
		"%s%-*s%s\n",
		BOW,
		termWidth-utf8.RuneCountInString(text),
		text,
		EOW,
	)
}

func getUserChoice(entries []Entry, shortcutMap map[string]int, tomlFile string) (string, string, string) {
	// ターミナル横幅取得
	termWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termWidth = w
	}

	// 初期表示
	PrintWhiteBackgroundLine(messages.AvailableDestinations)

	// 桁数を計算
	numWidth := 1
	if len(entries) >= 10 {
		numWidth = len(fmt.Sprintf("%d", len(entries)))
	}

	displayEntries := func(selectedIndex int, cursorMode bool) {
		for i, entry := range entries {
			expandedPath := expandPath(entry.Path)
			shortcutStr := ""
			if entry.Shortcut != "" {
				shortcutStr = fmt.Sprintf("(%s)", entry.Shortcut)
			}
			// 右寄せで桁揃え
			prefix := fmt.Sprintf("[%*d]%s %s → ", numWidth, i+1, shortcutStr, entry.Label)
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
	}

	// カーソル選択モードかどうかを判定
	selectedIndex := 0
	cursorMode := true // デフォルトでカーソルモードを有効にする
	inputBuffer := ""  // 複数文字入力用のバッファ

	// 初期表示（カーソルモードで開始）
	displayEntries(selectedIndex, true)
	fmt.Printf("%s\n", messages.AddCurrentDirectory)
	fmt.Printf("%s\n", messages.CursorModeHint)

	for {
		if !cursorMode {
			// 通常の入力モード表示
			displayEntries(selectedIndex, false)
			fmt.Printf("%s\n", messages.AddCurrentDirectory)
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
				cursorMode = true
				// 画面をクリア
				fmt.Print("\033[2J\033[H")
				PrintWhiteBackgroundLine(messages.AvailableDestinations)
				displayEntries(selectedIndex, true)
				fmt.Printf("%s\n", messages.AddCurrentDirectory)
				fmt.Printf("%s\n", messages.CursorModeHint)
				continue
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

			fmt.Println(messages.InvalidInput)
			continue
		} else {
			// カーソルモード
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
					entry := entries[selectedIndex]
					expandedPath := expandPath(entry.Path)
					return expandedPath, entry.Command, entry.Label
				case 27: // Escape
					cursorMode = false
					inputBuffer = ""           // バッファをクリア
					fmt.Print("\033[2J\033[H") // 画面クリア
					PrintWhiteBackgroundLine(messages.AvailableDestinations)
					continue
				case '+':
					return "ADD_CURRENT", "", ""
				case 'j': // j キーで下移動 (Vim風)
					inputBuffer = "" // バッファをクリア
					if selectedIndex < len(entries)-1 {
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
									// 有効な番号の場合、少し待ってから決定
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
					if selectedIndex < len(entries)-1 {
						selectedIndex++
						redraw = true
					}
				}
			}

			// 画面の再描画
			if redraw {
				// カーソルを最初の行に移動
				fmt.Printf("\033[%dA", len(entries)+2)
				displayEntries(selectedIndex, true)
				fmt.Printf("%s\n", messages.AddCurrentDirectory)
				fmt.Printf("%s\n", messages.CursorNavigationHint)
			}
		}
	}
}

func openNewShell(targetDir, command, label string) bool {
	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("%s %s\n", messages.DirectoryNotExist, targetDir)
		return false
	}

	openShellMessage := fmt.Sprintf("%s %s", messages.OpeningShell, targetDir)
	PrintWhiteBackgroundLine(openShellMessage)
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
	configWithHistory, err := loadConfigWithHistory(tomlFile)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorReadingConfig, err)
		return false
	}

	entries := getEntriesWithHistory(configWithHistory)
	shortcutMap := buildShortcutMap(entries)

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s ", messages.EnterLabel)
	label, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n%s\n", messages.OperationCancelled)
		return false
	}

	label = strings.TrimSpace(label)
	if label == "" {
		fmt.Println(messages.LabelCannotBeEmpty)
		return false
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
	fmt.Printf("  goto <number>        %s\n", messages.GoToDestinationByNumber)
	fmt.Printf("  goto <label>         %s\n", messages.GoToDestinationByLabel)
	fmt.Printf("  goto <shortcut>      %s\n", messages.GoToDestinationByShortcut)
	fmt.Printf("  goto -h, --help      %s\n", messages.ShowHelpMessage)
	fmt.Printf("  goto -v, --version   %s\n", messages.ShowVersionInfo)
	fmt.Printf("  goto --complete      %s\n", messages.ShowCompletionCandidates)
	fmt.Printf("  goto --history       %s\n", messages.ShowRecentUsageHistory)
	fmt.Printf("  goto --add           %s\n", messages.AddCurrentDirectoryToConfig)
	fmt.Printf("\n%s\n", messages.Examples)
	fmt.Printf("  goto 1              %s\n", messages.NavigateToFirstDest)
	fmt.Printf("  goto Home           %s\n", messages.NavigateToHomeDest)
	fmt.Printf("  goto h              %s\n", messages.NavigateUsingShortcut)
	fmt.Printf("  goto                %s\n", messages.ShowInteractiveMenuExample)
}

func showCompletions(entries []Entry) {
	// Output only labels for completion
	for _, entry := range entries {
		fmt.Println(entry.Label)
	}
}

func showHistory(config ConfigWithHistory) {
	// Get history file path
	historyFile, err := getHistoryFilePath()
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
		return
	}

	// Load history
	history, err := loadHistory(historyFile)
	if err != nil {
		// If there was an error loading history, try using the old format from config
		if len(config.History) == 0 {
			fmt.Println(messages.NoUsageHistoryFound)
			return
		}
		history.Entries = config.History
	}

	if len(history.Entries) == 0 {
		fmt.Println(messages.NoUsageHistoryFound)
		return
	}

	fmt.Println(messages.RecentUsageHistory)
	fmt.Println(strings.Repeat("=", 50))

	// Sort history by most recent first
	sortedHistory := make([]HistoryEntry, len(history.Entries))
	copy(sortedHistory, history.Entries)
	sort.Slice(sortedHistory, func(i, j int) bool {
		return sortedHistory[i].LastUsed.After(sortedHistory[j].LastUsed)
	})

	for i, hist := range sortedHistory {
		// Format timestamp for display
		timeStr := hist.LastUsed.Format("2006-01-02 15:04:05")

		// Get destination path if exists
		pathStr := ""
		if dest, exists := config.Destinations[hist.Label]; exists {
			pathStr = fmt.Sprintf(" → %s", expandPath(dest.Path))
		}

		fmt.Printf("%2d. %s%s\n", i+1, hist.Label, pathStr)
		fmt.Printf("    📅 %s\n", timeStr)

		if i < len(sortedHistory)-1 {
			fmt.Println()
		}
	}
}
