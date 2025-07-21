package main

import (
	"fmt"
	"os"

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

func PrintWhiteBackgroundLine(text string) {
	// ターミナル横幅を固定で80に設定してテスト
	termWidth := 80

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_multibyte.go <text>")
		return
	}

	text := os.Args[1]
	fmt.Printf("テストテキスト: '%s'\n", text)
	fmt.Printf("表示幅: %d 文字\n", getDisplayWidth(text))
	fmt.Println("白背景での表示:")
	PrintWhiteBackgroundLine(text)
	fmt.Println("通常表示:")
	fmt.Println(text)
}
