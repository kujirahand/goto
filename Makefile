# Makefile for goto command

.PHONY: all build-go install-go install-completion clean test help

# Default target
all: build-go

# Build Go version
build-go:
	@echo "Building Go version..."
	cd go && go build -o goto goto.go config.go goto_version.go
	@echo "✅ Go version built successfully: go/goto"

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
	@echo "  install-go       - Install Go version to /usr/local/bin"
	@echo "  install-completion - Install shell completion scripts"
	@echo "  install-all      - Install binary and completion scripts"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Test Go version"
	@echo "  help             - Show this help"
