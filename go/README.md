# Go Version of goto

This directory contains the Go implementation of the `goto` command with identical functionality to the Python version.

## Features

- 🚀 Interactive directory navigation menu
- 📝 TOML configuration file support
- 🔢 Navigate by number (e.g., `goto 1`)
- 🏷️  Navigate by label name (e.g., `goto Home`)  
- ⌨️  Navigate by shortcut key (e.g., `goto h`)
- ➕ Add current directory with `[+]` key
- ⚡ Execute commands after navigation
- 🐚 Opens new shell sessions (preserves current shell)

## Build

To build the Go version:

```sh
cd go
go build -o goto goto.go config.go
```

## Install

After building, you can install the Go version by copying the binary to your PATH:

```sh
cp goto /usr/local/bin/goto-go
# or
cp goto ~/bin/goto-go
```

## Usage

The Go version has identical usage to the Python version:

```sh
# Interactive mode
./goto

# Direct navigation
./goto 1              # Navigate to 1st destination
./goto Home           # Navigate to 'Home' destination  
./goto h              # Navigate using shortcut 'h'

# Help
./goto --help
```

## Configuration

Uses the same `~/.goto.toml` configuration file as the Python version. Example:

```toml
[Home]
path = "~/"
shortcut = "h"

[Desktop] 
path = "~/Desktop"
shortcut = "d"

[MyProject]
path = "~/workspace/my-project"
shortcut = "p"
command = "ls -la && git status"
```

## Performance

The Go version offers:
- ⚡ Faster startup time
- 📦 Single binary with no dependencies
- 🔧 Cross-platform compilation support
- 💾 Lower memory usage

## Dependencies

- Go 1.21+ 
- `github.com/BurntSushi/toml` for TOML parsing

Dependencies are automatically managed via `go.mod`.
