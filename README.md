# goto command

`goto` command for navigating directories quickly.

## Install

Please install `goto` command by following the steps below.

### Clone and Install

```sh
# Clone repository
git clone https://github.com/kujirahand/goto.git
# Install package
cd goto/bin
pip3 install -r requirement.txt
```

### Install as Shell Function

The recommended way to install `goto` is as a shell function, which allows proper directory navigation:

```sh
# Run the install script
./install.sh
```

This will:
1. Copy the necessary Python files to `~/.local/bin`
2. Add the `goto` function to your shell configuration file (`.zshrc` or `.bashrc`)
3. Make the function available in your current shell

After installation, reload your shell configuration:

```sh
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc
```

### Alternative: Add to PATH (deprecated)

You can also add the `goto` command to your PATH, but this method has limitations with directory changing:

```sh
export PATH="$PATH:/path/to/goto/bin"
```

### Config file - `~/.goto.toml`

You can create a configuration file at `~/.goto.toml` to customize the behavior of the `goto` command. Here is an example configuration:



