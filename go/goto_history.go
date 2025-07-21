// goto_history.go - History management functions
// This file contains functions for managing usage history,
// including displaying, sorting, and formatting history data.

package main

import (
	"fmt"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ShowHistory displays the usage history with timestamps and paths
func ShowHistory(customConfigFile, customHistoryFile string) {
	// Get configuration file path
	var tomlFile string
	if customConfigFile != "" {
		tomlFile = customConfigFile
	} else {
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
			return
		}
		tomlFile = filepath.Join(usr.HomeDir, ".goto.toml")
	}

	// Load configuration
	config, err := loadConfig(tomlFile)
	if err != nil {
		fmt.Printf("%s %v\n", messages.ErrorReadingConfig, err)
		return
	}

	// Get history file path
	var historyFile string
	if customHistoryFile != "" {
		historyFile = customHistoryFile
	} else {
		var err error
		historyFile, err = getHistoryFilePath()
		if err != nil {
			fmt.Printf("%s %v\n", messages.ErrorGettingUser, err)
			return
		}
	}

	// Load history
	history, err := loadHistory(historyFile)
	if err != nil {
		fmt.Println(messages.NoUsageHistoryFound)
		return
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
		if dest, exists := config[hist.Label]; exists {
			pathStr = fmt.Sprintf(" → %s", expandPath(dest.Path))
		}

		fmt.Printf("%2d. %s%s\n", i+1, hist.Label, pathStr)
		fmt.Printf("    📅 %s\n", timeStr)

		if i < len(sortedHistory)-1 {
			fmt.Println()
		}
	}
}

func UpdateHistory(tomlFile string, label string, customHistoryFile string) error {
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
