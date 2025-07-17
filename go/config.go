package main

// DefaultConfig contains the default TOML configuration
const DefaultConfig = `[Home]
path = "~/"
shortcut = "h"

[Desktop]
path = "~/Desktop"
shortcut = "d"

[Downloads]
path = "~/Downloads"
shortcut = "b"

[goto_app]
path = "~/repos/goto"
shortcut = "a"

[test_with_command]
path = "~/Desktop"
shortcut = "t"
command = "ls -la && echo 'Welcome to test directory!'"

[test_dir]
path = "/tmp"
shortcut = "t"
`
