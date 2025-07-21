// goto_print.go - Display and printing utility functions
// This file contains functions for formatting and displaying text output,
// including terminal formatting, help display, and text manipulation utilities.

package main

import (
	"fmt"
	"os"
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

// showInteractiveHelp displays help information same as goto -h
func showInteractiveHelp() {
	// Clear screen and show help
	fmt.Print("\033[2J\033[H")

	// Call the same help function as goto -h
	showHelp()

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
