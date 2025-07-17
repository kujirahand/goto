# Release v1.0.0

## 🚀 Features

### Core Navigation
- **Fast directory navigation** with interactive menu and command-line arguments
- **Multiple input methods**: Navigate by number, label name, or shortcut key
- **Current directory addition**: Add your current location with `[+]` shortcut
- **New shell sessions**: Opens fresh shell environment in target directory

### Configuration
- **TOML-based configuration** stored in `~/.goto.toml`
- **Automatic config creation** with sensible defaults on first run
- **Flexible destination setup** with path, shortcut, and optional commands
- **Home directory expansion** with `~/` support

### Command-Line Interface
- **Interactive mode**: Visual menu with numbered destinations
- **Direct navigation**: `goto <destination>` for immediate jumps
- **Help system**: `goto --help` for usage information
- **Version display**: `goto --version` for version information

### Shell Integration
- **Tab completion** for bash and zsh shells
- **Easy installation** with `make install-all`
- **Smart completion** showing destination labels
- **Cross-platform support** for Linux, macOS, and Windows

## 📦 Installation

### Quick Install
```bash
# Download the binary for your platform from the release assets below
# Make it executable and place in your PATH

# Or build from source:
git clone https://github.com/kujirahand/goto.git
cd goto
make install-all  # Includes tab completion
```

### Platform-specific Binaries

- **Linux amd64**: `goto-linux-amd64`
- **Linux arm64**: `goto-linux-arm64`
- **macOS Intel**: `goto-darwin-amd64`
- **macOS Apple Silicon**: `goto-darwin-arm64`
- **Windows amd64**: `goto-windows-amd64.exe`
- **Windows arm64**: `goto-windows-arm64.exe`

## 🛠 Usage Examples

```bash
# Interactive mode
goto

# Direct navigation
goto Home
goto 1
goto h

# Add current directory
goto
# Select [+] and follow prompts

# With tab completion (after installation)
goto <TAB>        # Shows all destinations
goto H<TAB>       # Completes to "Home"
```

## 📋 Configuration Example

```toml
[Home]
path = "~/"
shortcut = "h"

[Project]
path = "~/workspace/my-project"
shortcut = "p"
command = "git status && ls -la"

[Downloads]
path = "~/Downloads"
shortcut = "d"
```

## 🔧 Technical Details

- **Language**: Go 1.21+
- **Dependencies**: github.com/BurntSushi/toml
- **Config file**: `~/.goto.toml`
- **Binary size**: ~3MB (single file, no runtime dependencies)
- **Performance**: Instant startup and navigation

## 🏗 Build Information

Built with Go cross-compilation for maximum compatibility across platforms.

## 🎯 What's Next

This is the initial stable release. Future enhancements may include:
- Destination management commands (edit/delete)
- Import/export configurations
- Enhanced command execution features
