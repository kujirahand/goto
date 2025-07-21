package main

import (
	"os"
	"os/exec"
	"strings"
)

// URLかどうかを判定する関数
func IsURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// URLをデフォルトブラウザで開く関数
func OpenURL(url string) error {
	var cmd *exec.Cmd

	// OSに応じたコマンドを設定
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case FileExists("/usr/bin/open"): // macOS
		cmd = exec.Command("open", url)
	default: // Linux and others
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

// ファイルが存在するかチェックする関数
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
