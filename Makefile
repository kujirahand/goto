# Makefile for goto command

.PHONY: all build-go install-go install-completion clean test help build-release

# Default target
all: build-go

# Build Go version
build-go:
	@echo "Building Go version..."
	cd go && go build -o goto goto.go config.go goto_version.go
	@echo "✅ Go version built successfully: go/goto"

# Build release binaries for multiple platforms
build-release:
	@echo "Building release binaries for multiple platforms..."
	@mkdir -p releases
	@echo "Building for Linux amd64..."
	cd go && GOOS=linux GOARCH=amd64 go build -o ../releases/goto-linux-amd64 goto.go config.go goto_version.go
	@echo "Building for Linux arm64..."
	cd go && GOOS=linux GOARCH=arm64 go build -o ../releases/goto-linux-arm64 goto.go config.go goto_version.go
	@echo "Building for macOS amd64 (Intel)..."
	cd go && GOOS=darwin GOARCH=amd64 go build -o ../releases/goto-darwin-amd64 goto.go config.go goto_version.go
	@echo "Building for macOS arm64 (Apple Silicon)..."
	cd go && GOOS=darwin GOARCH=arm64 go build -o ../releases/goto-darwin-arm64 goto.go config.go goto_version.go
	@echo "Building for Windows amd64..."
	cd go && GOOS=windows GOARCH=amd64 go build -o ../releases/goto-windows-amd64.exe goto.go config.go goto_version.go
	@echo "Building for Windows arm64..."
	cd go && GOOS=windows GOARCH=arm64 go build -o ../releases/goto-windows-arm64.exe goto.go config.go goto_version.go
	@echo "✅ All release binaries built successfully in releases/ directory"

# Install Go version to /usr/local/bin
install-go: build-go
	@echo "Installing Go version to /usr/local/bin/goto..."
	sudo cp go/goto /usr/local/bin/goto
	@echo "✅ Go version installed as 'goto'"

# Install completion scripts
install-completion:
	@echo "Installing completion scripts..."
	@# Create completion directories if they don't exist
	@mkdir -p ~/.bash_completion.d
	@mkdir -p ~/.zsh/completions
	@# Install bash completion
	@cp completion/goto-completion.bash ~/.bash_completion.d/
	@echo "📝 Bash completion installed to ~/.bash_completion.d/"
	@# Install zsh completion
	@cp completion/_goto ~/.zsh/completions/
	@echo "📝 Zsh completion installed to ~/.zsh/completions/"
	@echo ""
	@echo "To enable completion, add these lines to your shell config:"
	@echo ""
	@echo "For bash (~/.bashrc or ~/.bash_profile):"
	@echo "  source ~/.bash_completion.d/goto-completion.bash"
	@echo ""
	@echo "For zsh (~/.zshrc):"
	@echo "  fpath=(~/.zsh/completions \$$fpath)"
	@echo "  autoload -U compinit && compinit"
	@echo ""
	@echo "Then restart your shell or run: source ~/.bashrc (or ~/.zshrc)"

# Install everything (binary + completion)
install-all: install-go install-completion
	@echo "✅ Complete installation finished!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f go/goto
	rm -rf releases/
	@echo "✅ Clean completed"

# Test Go version
test: build-go
	@echo "Testing Go version..."
	cd go && ./goto --help
	@echo "✅ Go version working"

# Show help
help:
	@echo "Available targets:"
	@echo "  all              - Build Go version (default)"
	@echo "  build-go         - Build Go version"
	@echo "  build-release    - Build release binaries for multiple platforms"
	@echo "  install-go       - Install Go version to /usr/local/bin"
	@echo "  install-completion - Install shell completion scripts"
	@echo "  install-all      - Install binary and completion scripts"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Test Go version"
	@echo "  help             - Show this help"
