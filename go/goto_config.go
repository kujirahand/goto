// goto_config.go - Configuration file I/O operations
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

// getEntriesFromConfig converts config map to sorted Entry slice
func getEntriesFromConfig(config map[string]Destination, customHistoryFile string) []Entry {
	var entries []Entry

	// Load history from separate JSON file for sorting
	var historyFile string
	var err error

	if customHistoryFile != "" {
		historyFile = customHistoryFile
	} else {
		historyFile, err = getHistoryFilePath()
		if err != nil {
			// If we can't get history file path, proceed without history sorting
			for label, dest := range config {
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
	}

	history, err := loadHistory(historyFile)
	if err != nil {
		// If history file doesn't exist or has error, proceed without history sorting
		for label, dest := range config {
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
	for label, dest := range config {
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

func updateHistory(tomlFile string, label string, customHistoryFile string) error {
	// Get history file path
	var historyFile string
	var err error

	if customHistoryFile != "" {
		historyFile = customHistoryFile
	} else {
		historyFile, err = getHistoryFilePath()
		if err != nil {
			return err
		}
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
