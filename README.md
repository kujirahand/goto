# goto command

`goto` command for navigating directories quickly.

This is a Go implementation providing fast, dependency-free directory navigation.

- [日本語](README-ja.md) / [中文](README-zh.md) / [한국어](README-ko.md) / [Español](README-es.md)

## Quick Start

1. **Download** the latest binary for your platform from [Releases](https://github.com/kujirahand/goto/releases)
2. **Make it executable** and place in your PATH
3. **Run** `goto` to see the interactive menu

## Key Features

- **Fast Directory Navigation**: Jump to frequently used directories instantly
- **Smart History**: Automatically sorts destinations by most recently used
- **Multiple Input Methods**: Use numbers, labels, or shortcut keys
- **Tab Completion**: Bash and Zsh completion support
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Multilingual Support**: Automatic language detection (English, Japanese, Chinese, Korean)
- **Zero Dependencies**: Single binary with no external dependencies

## Install

Please install `goto` command by following the steps below.

### Download Pre-built Binary (Recommended)

The easiest way to install `goto` is to download a pre-built binary from the GitHub releases page:

1. **Visit the releases page**: <https://github.com/kujirahand/goto/releases>
2. **Download the binary for your platform**:
   - **Linux amd64**: `goto-linux-amd64`
   - **Linux arm64**: `goto-linux-arm64`
   - **macOS Intel**: `goto-darwin-amd64`
   - **macOS Apple Silicon**: `goto-darwin-arm64`
   - **Windows amd64**: `goto-windows-amd64.exe`
   - **Windows arm64**: `goto-windows-arm64.exe`

3. **Make it executable and place in your PATH**:

   **For Linux/macOS**:

   ```sh
   # Download and make executable
   chmod +x goto-*
   
   # Move to a directory in your PATH
   sudo mv goto-* /usr/local/bin/goto
   
   # Or create a local bin directory (if it doesn't exist)
   mkdir -p ~/bin
   mv goto-* ~/bin/goto
   export PATH="$PATH:$HOME/bin"  # Add this to your shell config
   ```

   **For Windows**:
   - Rename the downloaded file to `goto.exe`
   - Place it in a directory that's in your PATH, or create a new directory and add it to PATH

4. **Verify installation**:

   ```sh
   goto --version
   ```

### Clone and Build from Source

```sh
# Clone repository
git clone https://github.com/kujirahand/goto.git
# Build
cd goto
make
```

### Build Release Archives (for developers)

To create release archives for all platforms:

```sh
# Create ZIP archives for all platforms (binaries are auto-cleaned)
make build-release-zip

# Note: Generated ZIP files are in releases/ but excluded from git tracking
```

### Add to PATH

After building, add the compiled `goto` executable to your PATH by adding the following line to your shell configuration file (`.bashrc`, `.zshrc`, etc.):

```sh
export PATH="$PATH:/path/to/goto"
```

For example, if you cloned to your home directory:

```sh
export PATH="$PATH:$HOME/goto"
```

After adding to PATH, reload your shell configuration:

```sh
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc
```

### Install with Tab Completion (Source Build)

If you built from source, you can install both the binary and completion scripts:

```sh
# Build and install everything (requires source code)
make install-all
```

### Manual Tab Completion Setup (For Pre-built Binaries)

If you downloaded a pre-built binary, you can still set up tab completion manually:

1. **Download completion scripts**:

   ```sh
   # Create completion directories
   mkdir -p ~/.bash_completion.d ~/.zsh/completions
   
   # Download bash completion script
   curl -o ~/.bash_completion.d/goto-completion.bash \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/goto-completion.bash
   
   # Download zsh completion script  
   curl -o ~/.zsh/completions/_goto \
     https://raw.githubusercontent.com/kujirahand/goto/main/completion/_goto
   ```

2. **Add to your shell configuration**:

   **For bash** (`~/.bashrc` or `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **For zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. **Restart your shell or reload configuration**:

   ```sh
   source ~/.bashrc   # for bash
   source ~/.zshrc    # for zsh
   ```

### Advanced Installation with Tab Completion (Source Build)

For the best experience when building from source, install both the binary and completion scripts:

```sh
# Build and install everything
make install-all
```

This will:

1. Install the `goto` binary to `/usr/local/bin/`
2. Install shell completion scripts
3. Show instructions for enabling completion

#### Alternative: Manual Completion Setup (Source Build)

If you built from source but prefer to install completion manually:

1. Install completion scripts:

   ```sh
   make install-completion
   ```

2. Add the following to your shell configuration:

   **For bash** (`~/.bashrc` or `~/.bash_profile`):

   ```sh
   source ~/.bash_completion.d/goto-completion.bash
   ```

   **For zsh** (`~/.zshrc`):

   ```sh
   fpath=(~/.zsh/completions $fpath)
   autoload -U compinit && compinit
   ```

3. Restart your shell or reload configuration:

   ```sh
   source ~/.bashrc   # for bash
   source ~/.zshrc    # for zsh
   ```

#### Using Tab Completion

Once enabled, you can use tab completion with the `goto` command:

```sh
goto <TAB>        # Shows all available destinations
goto h<TAB>       # Completes shortcuts starting with 'h'
goto Home<TAB>    # Completes labels starting with 'Home'
goto 1<TAB>       # Shows destinations with numbers starting with '1'
```

## Configuration

### Config Files

The `goto` command uses the following configuration files:

- **`~/.goto.toml`**: Main configuration file containing your destinations
- **`~/.goto.history.json`**: History data storing your recent usage information

When you first run `goto`, it will automatically create a default configuration file with sample destinations.

Example configuration:

```toml
[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"

[Downloads]
path = "~/Downloads"
shortcut = "b"

[MyProject]
path = "~/workspace/my-project"
shortcut = "p"
command = "ls -la && git status"
```

Each destination can have:

- `path` (required): Directory path (supports `~` for home directory)
- `shortcut` (optional): Single character shortcut key
- `command` (optional): Command to execute after changing directory

## Usage

### Basic Usage

Run the `goto` command to see available destinations:

```sh
goto
```

### Command Line Arguments

You can also specify a destination directly as a command line argument:

```sh
# Using number
goto 1
goto 4

# Using label name
goto Home
goto MyProject

# Using shortcut key
goto h
goto p

# View usage history
goto --history

# Show help
goto --help

# Show version
goto --version
```

This is useful for scripting or when you know exactly where you want to go.

### Interactive Mode

When run without arguments, `goto` displays an interactive menu:

Example output:

```text
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
2. Desktop → /Users/username/Desktop (shortcut: d)
3. Downloads → /Users/username/Downloads (shortcut: b)
4. MyProject → /Users/username/workspace/my-project (shortcut: p)

➕ [+] Add current directory

Please enter the number, shortcut key, or [+] to add current directory:
Enter number, shortcut key, or [+]:
```

You can navigate by:

- **Number**: Enter `1`, `2`, `3`, etc.
- **Shortcut**: Enter `h`, `d`, `b`, etc.
- **Add current**: Enter `+` to add current directory

### Adding Current Directory

You can add the current directory to your goto destinations by selecting `[+]`:

```sh
goto
# Select [+] from the menu
# Enter a label for the current directory
# Optionally enter a shortcut key
```

Example:

```text
Enter number, shortcut key, or [+]: +
📍 Current directory: /Users/username/workspace/new-project
Enter a label for this directory: NewProject
Enter a shortcut key (optional, press Enter to skip): n
✅ Added 'NewProject' → /Users/username/workspace/new-project
🔑 Shortcut: n
```

This feature allows you to quickly add frequently used directories to your goto list.

### New Shell Functionality

When you select a destination, `goto` opens a new shell session in the target directory. This means:

- Your current shell session remains unchanged
- You get a fresh shell environment in the new location
- Type `exit` to return to your previous shell
- If a `command` is specified in the configuration, it will be executed automatically

### Usage History

`goto` automatically tracks usage history and displays destinations in order of most recently used. This makes frequently accessed directories appear at the top of the interactive menu.

#### Viewing Usage History

You can view your recent usage history with:

```sh
goto --history
```

Example output:

```text
📈 Recent usage history:
==================================================
 1. Home → /Users/username
    📅 2025-07-18 16:08:38

 2. Desktop → /Users/username/Desktop
    📅 2025-07-18 16:04:40

 3. MyProject → /Users/username/workspace/my-project
    📅 2025-07-18 15:30:15
```

#### How History Works

- **Automatic tracking**: Every time you navigate to a destination, the timestamp is recorded
- **Smart sorting**: In interactive mode, destinations are sorted by most recently used first
- **Persistent storage**: History is stored in the `~/.goto.toml` configuration file
- **No manual maintenance**: History is automatically updated - no need to manually manage it

#### History Storage

Usage history is stored in your `~/.goto.history.json` file in the following format:

```json
{
  "entries": [
    {
      "label": "Home",
      "last_used": "2025-07-18T16:08:38+09:00"
    },
    {
      "label": "Desktop",
      "last_used": "2025-07-18T16:04:40+09:00"
    }
  ]
}
```

This intelligent ordering ensures that your most frequently used directories are always easily accessible.

## Multilingual Support

`goto` automatically detects your system language and displays messages in your preferred language. Currently supported languages:

- **English** (en) - Default
- **Japanese** (ja) - 日本語
- **Chinese** (zh) - 中文
- **Korean** (ko) - 한국어
- **Spanish** (es) - Español

### How Language Detection Works

The application automatically detects your system language by checking the following environment variables in order:

1. `LANG`
2. `LANGUAGE`
3. `LC_ALL`
4. `LC_MESSAGES`

For example, if your system is set to Japanese (`LANG=ja_JP.UTF-8`), `goto` will automatically display all messages in Japanese.

### Example Output in Different Languages

**English:**

```text
🚀 goto - Navigate directories quickly
👉 Available destinations:
1. Home → /Users/username/ (shortcut: h)
📈 Recent usage history:
```

**Japanese:**

```text
🚀 goto - ディレクトリ間を素早く移動
👉 利用可能なディレクトリ:
1. Home → /Users/username/ (shortcut: h)
📈 最近の使用履歴:
```

**Chinese:**

```text
🚀 goto - 快速导航目录
👉 可用目录:
1. Home → /Users/username/ (shortcut: h)
📈 最近使用历史:
```

**Korean:**

```text
🚀 goto - 디렉토리 빠른 탐색
👉 사용 가능한 디렉토리:
1. Home → /Users/username/ (shortcut: h)
📈 최근 사용 기록:
```

**Spanish:**

```text
🚀 goto - Navegar directorios rápidamente
👉 Destinos disponibles:
1. Home → /Users/username/ (shortcut: h)
📈 Historial de uso reciente:
```

### Language Override

If you want to use a specific language regardless of your system settings, you can set the `LANG` environment variable:

```sh
# Use Japanese interface
LANG=ja_JP.UTF-8 goto

# Use English interface
LANG=en_US.UTF-8 goto

# Use Chinese interface
LANG=zh_CN.UTF-8 goto

# Use Korean interface
LANG=ko_KR.UTF-8 goto

# Use Spanish interface
LANG=es_ES.UTF-8 goto
```

### Supported Languages

The multilingual support covers all user interface elements including:

- Interactive menu messages
- Navigation confirmations
- Error messages
- Help text
- History display
- Configuration messages

All messages are automatically localized based on your system language settings, providing a native experience for international users.

### Examples

1. **Navigate using command line argument (number):**

   ```sh
   goto 1
   goto 4
   ```

2. **Navigate using command line argument (label):**

   ```sh
   goto Home
   goto MyProject
   ```

3. **Navigate using command line argument (shortcut):**

   ```sh
   goto h
   goto p
   ```

4. **Interactive navigation:**

   ```sh
   goto
   # Then enter: h (shortcut), 1 (number), or Home (label)
   ```

5. **Add current directory:**

   ```sh
   cd /path/to/important/project
   goto
   # Enter: +
   # Label: ImportantProject
   # Shortcut: i
   ```

6. **View usage history:**

   ```sh
   goto --history
   ```



