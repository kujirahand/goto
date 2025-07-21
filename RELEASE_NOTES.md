# Release v1.1.2

## 🚀 Enhancements

- **Improved Build Process**: The `Makefile` has been enhanced to automate the creation of ZIP archives for releases, streamlining the packaging process.
- **Added Quick Start Guide**: A "Quick Start" section has been added to the `README.md` to help new users get started more easily.

## 🐛 Bug Fixes

- **Fixed Cursor Mode Bug**: Fixed a bug in interactive cursor mode where input was not handled correctly.

## 🧹 Maintenance

- **Removed Binary Artifacts**: Pre-compiled binary files have been removed from the Git repository to reduce its size. They are now provided exclusively through release assets.
- **Updated Documentation**: All `README-*.md` files have been updated to reflect the latest changes.

---

## Release v1.1.1

## 🔧 Code Improvements & Bug Fixes

### Code Structure Refactoring
- **Extracted utility functions**: Moved URL and file utility functions to dedicated `utils.go` file
- **Improved function naming**: Standardized function names (e.g., `showHistory` → `ShowHistory`)
- **Better code organization**: Cleaner separation of concerns across modules

### History Management Fixes
- **Fixed history update timing**: History now updates before opening shell instead of after
- **Improved display formatting**: Changed item numbering from "-" to "." for items 10 and above
- **Enhanced test coverage**: Added comprehensive tests for history functionality

### Development & Testing
- **Enhanced test infrastructure**: Improved test helper with better resource management
- **Added test automation**: New `test_all.sh` script for running all tests
- **Build system updates**: Added `utils.go` to Makefile build sources

---

# Release v1.1.0

## 🚀 What's New

### Enhanced History Management
- **Separate history storage**: History data moved from `~/.goto.toml` to `~/.goto.history.json`
- **Smart history sorting**: Destinations sorted by recent usage, with most frequently used items at the top
- **Automatic history limits**: Maximum of 100 history entries to prevent file bloat and maintain performance
- **Backward compatibility**: Seamless migration from old TOML-based history

### Improved Error Handling
- **Detailed error messages**: Clear error descriptions with file paths and suggested fixes
- **Multi-language error support**: Error messages localized for all supported languages
- **Config file validation**: Better handling of corrupted or outdated configuration files

### Complete Localization
- **Interactive cursor mode messages**: All hardcoded messages now properly localized
- **5 language support**: Japanese, English, Chinese, Korean, and Spanish
- **System language detection**: Automatically adapts to user's system language settings
- **Consistent UI**: Unified white background styling for better visual consistency

### User Experience Improvements
- **Enhanced visual feedback**: Improved terminal output with consistent styling
- **Better navigation hints**: Context-sensitive help messages for different interaction modes
- **Optimized performance**: Faster startup and reduced memory usage with history limits

## 🔧 Technical Improvements
- **JSON-based history**: More efficient and flexible history storage format
- **Asynchronous history saving**: Non-blocking history updates for better responsiveness
- **Memory optimization**: Automatic cleanup of excessive history entries
- **Code organization**: Better separation of concerns with dedicated locale management

---

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
