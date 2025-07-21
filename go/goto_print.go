// goto_print.go - Display and printing utility functions
// This file contains functions for formatting and displaying text output,
// including terminal formatting, help display, and text manipulation utilities.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// shortenPathMiddle truncates a path in the middle with ellipsis
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

// showInteractiveHelp displays help information using less or direct output
func showInteractiveHelp() {
	// Get current language for README file selection
	lang := detectLanguage()

	// Determine README file name based on language
	var readmeFile string
	switch lang {
	case "ja":
		readmeFile = "README-ja.md"
	case "ko":
		readmeFile = "README-ko.md"
	case "zh":
		readmeFile = "README-zh.md"
	case "es":
		readmeFile = "README-es.md"
	default:
		readmeFile = "README.md"
	}

	// Try to find README file in various locations
	var readmePath string
	possiblePaths := []string{
		readmeFile,                                         // Current directory
		filepath.Join("..", readmeFile),                    // Parent directory
		filepath.Join("..", "..", readmeFile),              // Two levels up
		filepath.Join("/usr/local/share/goto", readmeFile), // System installation
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			readmePath = path
			break
		}
	}

	// Check if less command is available
	lessAvailable := false
	if _, err := exec.LookPath("less"); err == nil {
		lessAvailable = true
	}

	if readmePath != "" && lessAvailable {
		// Use less to display README
		fmt.Printf("\n%s %s\n", "📖 Displaying help from", readmePath)
		fmt.Println("Press 'q' to quit, arrow keys to navigate")
		fmt.Println(strings.Repeat("=", 50))

		cmd := exec.Command("less", readmePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else if readmePath != "" {
		// Display README content directly
		fmt.Printf("\n%s %s\n", "📖 Help from", readmePath)
		fmt.Println(strings.Repeat("=", 50))

		content, err := os.ReadFile(readmePath)
		if err == nil {
			fmt.Println(string(content))
		} else {
			fmt.Printf("Error reading help file: %v\n", err)
			showBasicHelp()
		}

		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("Press any key to continue...")

		// Wait for key press
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			buffer := make([]byte, 1)
			os.Stdin.Read(buffer)
			term.Restore(int(os.Stdin.Fd()), oldState)
		}
	} else {
		// Show basic help if README not found
		showBasicHelp()
	}
}

// showBasicHelp displays basic help information when README is not available
func showBasicHelp() {
	fmt.Println("\n📖 goto - Quick Directory Navigation")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("KEYBOARD SHORTCUTS:")
	fmt.Println("  [1-9]      - Navigate to destination by number")
	fmt.Println("  [a-z]      - Navigate using shortcut key")
	fmt.Println("  [0]        - Exit application")
	fmt.Println("  [+]        - Add current directory")
	fmt.Println("  [?]        - Show this help")
	fmt.Println("  [j] or ↓   - Move cursor down")
	fmt.Println("  [k] or ↑   - Move cursor up")
	fmt.Println("  [Enter]    - Select highlighted option")
	fmt.Println("  [Esc]      - Switch to input mode")
	fmt.Println()
	fmt.Println("DISPLAY FORMAT:")
	fmt.Println("  1.(shortcut) Label → Path    - For items 1-9")
	fmt.Println("  -.(shortcut) Label → Path    - For items 10+")
	fmt.Println("  0. Exit                      - Exit application")
	fmt.Println()
	fmt.Println("COMMAND LINE OPTIONS:")
	fmt.Println("  goto                 - Show interactive menu")
	fmt.Println("  goto <number>        - Go to destination by number")
	fmt.Println("  goto <label>         - Go to destination by label")
	fmt.Println("  goto <shortcut>      - Go to destination by shortcut")
	fmt.Println("  goto --add           - Add current directory")
	fmt.Println("  goto --history       - Show usage history")
	fmt.Println("  goto --help          - Show detailed help")
	fmt.Println("  goto --version       - Show version information")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("Press any key to continue...")

	// Wait for key press
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		buffer := make([]byte, 1)
		os.Stdin.Read(buffer)
		term.Restore(int(os.Stdin.Fd()), oldState)
	}
}
