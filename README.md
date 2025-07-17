# goto command

`goto` command for navigating directories quickly.

This is a Go implementation providing fast, dependency-free directory navigation.

## Quick Start

1. **Download** the latest binary for your platform from [Releases](https://github.com/kujirahand/goto/releases)
2. **Make it executable** and place in your PATH
3. **Run** `goto` to see the interactive menu

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

### Add to PATH

Add the `goto/go` directory to your PATH by adding the following line to your shell configuration file (`.bashrc`, `.zshrc`, etc.):

```sh
export PATH="$PATH:/path/to/goto/go"
```

For example, if you cloned to your home directory:

```sh
export PATH="$PATH:$HOME/goto/go"
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

### Config file - `~/.goto.toml`

The `goto` command uses a TOML configuration file located at `~/.goto.toml`. When you first run `goto`, it will automatically create a default configuration file with sample destinations.

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



