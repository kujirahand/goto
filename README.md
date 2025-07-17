# goto command

`goto` command for navigating directories quickly.

This repository contains two implementations:
- **Python version** (`bin/` directory) - Full-featured with rich dependencies
- **Go version** (`go/` directory) - Fast, single binary with no dependencies

## Install

Please install `goto` command by following the steps below.

### Clone and Install

```sh
# Clone repository
git clone https://github.com/kujirahand/goto.git
# Install dependencies
cd goto/bin
pip3 install -r requirement.txt
```

### Add to PATH

Add the `goto/bin` directory to your PATH by adding the following line to your shell configuration file (`.bashrc`, `.zshrc`, etc.):

```sh
export PATH="$PATH:/path/to/goto/bin"
```

For example, if you cloned to your home directory:

```sh
export PATH="$PATH:$HOME/goto/bin"
```

After adding to PATH, reload your shell configuration:

```sh
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc
```

### Go Version (Alternative)

For faster performance and no dependencies:

```sh
# Build the Go version
cd go
go build -o goto goto.go config.go

# Copy to your PATH
cp goto /usr/local/bin/goto-go
```

The Go version has identical functionality but offers:

- ⚡ Faster startup time
- 📦 Single binary with no dependencies  
- 🔧 Cross-platform compilation support

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



