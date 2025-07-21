// config_io.go - Configuration file I/O operations
// This file contains functions for reading and writing configuration files,
// including TOML configuration files and JSON history files.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// getHistoryFilePath returns the full path to the history file
func getHistoryFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".goto.history.json"), nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(tomlFile string) {
	err := os.WriteFile(tomlFile, []byte(DefaultConfig), 0644)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorWritingConfigFile, err)
		os.Exit(1)
	}
	fmt.Printf("%s %s\n", messages.CreatedDefaultConfig, tomlFile)
}

// loadConfig loads the TOML configuration file
func loadConfig(tomlFile string) (map[string]Destination, error) {
	var config map[string]Destination
	_, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// loadHistory loads the JSON history file
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

// saveHistory saves the history data to JSON file
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

// loadConfigWithHistory loads TOML configuration including embedded history section (legacy)
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
