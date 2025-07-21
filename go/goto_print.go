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

	"golang.org/x/term"
	"golang.org/x/text/width"
)

// getDisplayWidth calculates the display width of a string considering multi-byte characters
func getDisplayWidth(s string) int {
	displayWidth := 0
	for _, r := range s {
		switch width.LookupRune(r).Kind() {
		case width.EastAsianFullwidth, width.EastAsianWide:
			displayWidth += 2
		case width.EastAsianHalfwidth, width.EastAsianNarrow, width.Neutral:
			displayWidth += 1
		case width.EastAsianAmbiguous:
			// Ambiguous characters are typically displayed as 1 width in most terminals
			displayWidth += 1
		}
	}
	return displayWidth
}

// shortenPathMiddle truncates a path in the middle with ellipsis
func shortenPathMiddle(path string, maxLen int) string {
	// 表示幅を計算
	currentWidth := getDisplayWidth(path)
	if currentWidth <= maxLen {
		return path
	}

	// 省略が必要な場合
	r := []rune(path)
	ellipsis := "..."
	ellipsisWidth := getDisplayWidth(ellipsis)

	// 利用可能な幅から省略記号の幅を引く
	availableWidth := maxLen - ellipsisWidth
	if availableWidth < 6 {
		// 省略しすぎないように、最低限の文字数を確保
		if len(r) > maxLen {
			return string(r[:maxLen])
		}
		return path
	}

	// 前半と後半に分ける
	halfWidth := availableWidth / 2

	// 前半部分を取得
	var head []rune
	headWidth := 0
	for _, char := range r {
		charWidth := getDisplayWidth(string(char))
		if headWidth+charWidth > halfWidth {
			break
		}
		head = append(head, char)
		headWidth += charWidth
	}

	// 後半部分を取得
	var tail []rune
	tailWidth := 0
	for i := len(r) - 1; i >= 0; i-- {
		char := r[i]
		charWidth := getDisplayWidth(string(char))
		if tailWidth+charWidth > availableWidth-headWidth {
			break
		}
		tail = append([]rune{char}, tail...)
		tailWidth += charWidth
	}

	return string(head) + ellipsis + string(tail)
}

// PrintWhiteBackgroundLine prints a line with white background
func PrintWhiteBackgroundLine(text string) {
	// ターミナル横幅取得
	termWidth := 80
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		termWidth = w
	}

	// テキストの表示幅を計算
	textDisplayWidth := getDisplayWidth(text)

	// パディングが必要な文字数を計算
	paddingWidth := termWidth - textDisplayWidth
	if paddingWidth < 0 {
		paddingWidth = 0
	}

	// 一行全部を白背景にして表示
	BOW := "\033[47;30m"
	EOW := "\033[0m"
	fmt.Printf(
		"%s%s%*s%s\n",
		BOW,
		text,
		paddingWidth, "",
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
	
	// Get executable path to help locate README files
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)
	
	possiblePaths := []string{
		readmeFile,                                         // Current directory
		filepath.Join("..", readmeFile),                    // Parent directory
		filepath.Join("..", "..", readmeFile),              // Two levels up
		filepath.Join(execDir, readmeFile),                 // Same directory as executable
		filepath.Join(execDir, "..", readmeFile),           // Parent directory of executable
		filepath.Join(execDir, "..", "..", readmeFile),     // Two levels up from executable
		filepath.Join("/usr/local/share/goto", readmeFile), // System installation
		// Add absolute paths for common development locations
		filepath.Join("/Users/kujirahand/repos/goto", readmeFile),
		"/Users/kujirahand/repos/goto/README.md",           // Fallback to English README
		"/Users/kujirahand/repos/goto/README-ja.md",        // Fallback to Japanese README
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			readmePath = path
			break
		}
	}

	// If no README found, show debug information
	if readmePath == "" {
		fmt.Printf("\n❌ README file not found. Searched paths:\n")
		for i, path := range possiblePaths {
			if i < 5 { // Show first 5 paths to avoid too much output
				fmt.Printf("  - %s\n", path)
			}
		}
		fmt.Printf("  ... and %d more paths\n", len(possiblePaths)-5)
		fmt.Printf("📂 Current working directory: %s\n", func() string {
			if wd, err := os.Getwd(); err == nil {
				return wd
			}
			return "unknown"
		}())
		fmt.Printf("📂 Executable path: %s\n", execPath)
		showBasicHelp()
		return
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
	fmt.Println("  1. Label (shortcut) → Path    - For items 1-9")
	fmt.Println("  -. Label (shortcut) → Path    - For items 10+")
	fmt.Println("  0. Exit                       - Exit application")
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
