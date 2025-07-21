package main

// DefaultConfig contains the default TOML configuration
const DefaultConfig = `# Configuration for the "goto" command using TOML format
[Home]
path = "~/"
shortcut = "h"
command = "ls -la && echo 'Welcome to home directory!'"

[Desktop]
path = "~/Desktop"
shortcut = "d"

[Downloads]
path = "~/Downloads"
shortcut = "b"

[Documents]
path = "~/Documents"
shortcut = "D"

[goto-config]
path = "~/"
shortcut = "e"
command = "vi ~/.goto.toml"

["goto-web"]
path = "https://github.com/kujirahand/goto"

["kujirahand.com"]
path = "https://kujirahand.com"
shortcut = "K"
`
