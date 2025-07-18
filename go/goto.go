package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
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

func main() {
	// Get configuration file path
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("❌ Error getting current user: %v\n", err)
		os.Exit(1)
	}

	tomlFile := filepath.Join(usr.HomeDir, ".goto.toml")

	// Create default config if it doesn't exist
	if _, err := os.Stat(tomlFile); os.IsNotExist(err) {
		createDefaultConfig(tomlFile)
	}

	// Load configuration with history for sorting
	configWithHistory, err := loadConfigWithHistory(tomlFile)
	if err != nil {
		fmt.Printf("❌ Error reading configuration file: %v\n", err)
		os.Exit(1)
	}

	// Get entries sorted by history and shortcuts
	entries := getEntriesWithHistory(configWithHistory)
	shortcutMap := buildShortcutMap(entries)

	if len(entries) == 0 {
		fmt.Println("⚠️  No destinations configured in ~/.goto.toml")
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

		// Find destination by argument
		targetDir, command, label := findDestinationByArg(arg, entries, shortcutMap)

		if targetDir == "" {
			fmt.Printf("❌ Destination '%s' not found.\n", arg)
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

		fmt.Printf("🎯 Found destination: %s\n", label)
		success := openNewShell(targetDir, command, label)
		if success {
			// Update history
			if label != "" {
				err := updateHistory(tomlFile, label)
				if err != nil {
					fmt.Printf("⚠️  Warning: Failed to update history: %v\n", err)
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
		fmt.Println("ℹ️  No directory selected or operation cancelled.")
		os.Exit(0)
	}

	success := openNewShell(targetDir, command, label)
	if success {
		// Update history
		if label != "" {
			err := updateHistory(tomlFile, label)
			if err != nil {
				fmt.Printf("⚠️  Warning: Failed to update history: %v\n", err)
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
		fmt.Printf("❌ Error creating default configuration: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created default configuration file: %s\n", tomlFile)
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

	// Create a map for quick lookup of destination data
	destMap := make(map[string]Destination)
	for label, dest := range config.Destinations {
		destMap[label] = dest
	}

	// Create a map for quick lookup of history
	historyMap := make(map[string]time.Time)
	for _, hist := range config.History {
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
	// Load current config with history
	config, err := loadConfigWithHistory(tomlFile)
	if err != nil {
		return err
	}

	// Update or add history entry
	now := time.Now()
	found := false

	for i, hist := range config.History {
		if hist.Label == label {
			config.History[i].LastUsed = now
			found = true
			break
		}
	}

	if !found {
		config.History = append(config.History, HistoryEntry{
			Label:    label,
			LastUsed: now,
		})
	}

	// Save updated config
	return saveConfigWithHistory(tomlFile, config)
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

func getUserChoice(entries []Entry, shortcutMap map[string]int, tomlFile string) (string, string, string) {
	fmt.Println("👉 Available destinations:")
	for i, entry := range entries {
		expandedPath := expandPath(entry.Path)
		shortcutStr := ""
		if entry.Shortcut != "" {
			shortcutStr = fmt.Sprintf(" (%s)", entry.Shortcut)
		}
		fmt.Printf("%d. %s → %s%s\n", i+1, entry.Label, expandedPath, shortcutStr)
	}

	fmt.Println("\n➕ [+] Add current directory")
	fmt.Println("\nPlease enter the number, shortcut key, label name, or [+] to add current directory:")

	fmt.Print("Enter number, shortcut key, label name, or [+]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\nOperation cancelled.")
		return "", "", ""
	}

	choice = strings.TrimSpace(choice)

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

	fmt.Println("Invalid input.")
	return "", "", ""
}

func openNewShell(targetDir, command, label string) bool {
	// Check if directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("❌ Directory does not exist: %s\n", targetDir)
		return false
	}

	fmt.Printf("🚀 Opening new shell in: %s\n", targetDir)
	if label != "" {
		fmt.Printf("📍 Destination: %s\n", label)
	}

	// Get user's preferred shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	if command != "" {
		fmt.Printf("⚡ Will execute: %s\n", command)
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
			fmt.Printf("❌ Error opening shell: %v\n", err)
			return false
		}
	} else {
		// Simply open shell in the target directory
		fmt.Println("💡 Type 'exit' to return to previous shell")
		fmt.Println(strings.Repeat("=", 50))

		fmt.Printf("✅ You are now in: %s\n", targetDir)

		// Start new shell with the target directory as working directory
		cmd := exec.Command(shell)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = targetDir // Set working directory for the new shell

		err := cmd.Run()
		if err != nil {
			fmt.Printf("❌ Error opening shell: %v\n", err)
			return false
		}
	}

	return true
}

func createTempScript(targetDir, command, shell string) string {
	tempFile, err := os.CreateTemp("", "goto_*.sh")
	if err != nil {
		fmt.Printf("❌ Error creating temp file: %v\n", err)
		os.Exit(1)
	}

	scriptContent := fmt.Sprintf(`#!/bin/sh
cd "%s"
echo "📍 Current directory: $(pwd)"
echo "⚡ Executing: %s"
echo "%s"
%s
echo "%s"
echo "✅ Command completed. You are now in: $(pwd)"
echo "💡 Type 'exit' to return to previous shell"
exec "%s"
`, targetDir, command, strings.Repeat("-", 40), command, strings.Repeat("-", 40), shell)

	_, err = tempFile.WriteString(scriptContent)
	if err != nil {
		fmt.Printf("❌ Error writing temp script: %v\n", err)
		os.Exit(1)
	}

	tempFile.Close()

	// Make executable
	err = os.Chmod(tempFile.Name(), 0755)
	if err != nil {
		fmt.Printf("❌ Error making script executable: %v\n", err)
		os.Exit(1)
	}

	return tempFile.Name()
}

func addCurrentPathToConfig(tomlFile string) bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		return false
	}

	fmt.Printf("📍 Current directory: %s\n", currentDir)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter a label for this directory: ")
	label, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\n❌ Operation cancelled.")
		return false
	}

	label = strings.TrimSpace(label)
	if label == "" {
		fmt.Println("❌ Label cannot be empty.")
		return false
	}

	fmt.Print("Enter a shortcut key (optional, press Enter to skip): ")
	shortcut, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("\n❌ Operation cancelled.")
		return false
	}

	shortcut = strings.TrimSpace(shortcut)

	// Generate TOML entry
	tomlEntry := fmt.Sprintf("\n[%s]\n", label)
	tomlEntry += fmt.Sprintf("path = \"%s\"\n", currentDir)
	if shortcut != "" {
		tomlEntry += fmt.Sprintf("shortcut = \"%s\"\n", shortcut)
	}

	// Append to TOML file
	file, err := os.OpenFile(tomlFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Error opening config file: %v\n", err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(tomlEntry)
	if err != nil {
		fmt.Printf("❌ Error writing to config file: %v\n", err)
		return false
	}

	fmt.Printf("✅ Added '%s' → %s\n", label, currentDir)
	if shortcut != "" {
		fmt.Printf("🔑 Shortcut: %s\n", shortcut)
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

	fmt.Println("🚀 goto - Navigate directories quickly")
	fmt.Printf("\nConfiguration file: %s\n", configPath)
	fmt.Println("\nUsage:")
	fmt.Println("  goto                 Show interactive menu")
	fmt.Println("  goto <number>        Go to destination by number (e.g., goto 1)")
	fmt.Println("  goto <label>         Go to destination by label name")
	fmt.Println("  goto <shortcut>      Go to destination by shortcut key")
	fmt.Println("  goto -h, --help      Show this help message")
	fmt.Println("  goto -v, --version   Show version information")
	fmt.Println("  goto --complete      Show completion candidates (for shell completion)")
	fmt.Println("\nExamples:")
	fmt.Println("  goto 1              # Navigate to 1st destination")
	fmt.Println("  goto Home           # Navigate to 'Home' destination")
	fmt.Println("  goto h              # Navigate using shortcut 'h'")
	fmt.Println("  goto                # Show interactive menu")
}

func showCompletions(entries []Entry) {
	// Output only labels for completion
	for _, entry := range entries {
		fmt.Println(entry.Label)
	}
}
