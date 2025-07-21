# Makefile for goto command

# Go source files
GO_SOURCES = goto.go goto_config_default.go goto_config.go goto_history.go goto_print.go goto_version.go locale.go utils.go

# Build platforms
PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all build-go install-go install-completion clean test help build-release

# Default target
all: build-go

# Build Go version
build-go:
	@echo "Building Go version..."
	cd go && go build -o goto $(GO_SOURCES)
	@echo "✅ Go version built successfully: go/goto"

# Build release binaries for multiple platforms
build-release:
	@echo "Building release binaries for multiple platforms..."
	@mkdir -p releases
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d/ -f1); \
		GOARCH=$$(echo $$platform | cut -d/ -f2); \
		OUTPUT="goto-$$GOOS-$$GOARCH"; \
		if [ "$$GOOS" = "windows" ]; then OUTPUT="$$OUTPUT.exe"; fi; \
		echo "Building for $$GOOS $$GOARCH..."; \
		(cd go && GOOS=$$GOOS GOARCH=$$GOARCH go build -o ../releases/$$OUTPUT $(GO_SOURCES)); \
	done
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
	@echo "goto Makefile - Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  all              - Build Go version (default target)"
	@echo "  build-go         - Build Go version for current platform"
	@echo "  build-release    - Build release binaries for all supported platforms"
	@echo ""
	@echo "Installation targets:"
	@echo "  install-go       - Install Go version to /usr/local/bin"
	@echo "  install-completion - Install shell completion scripts"
	@echo "  install-all      - Install binary and completion scripts"
	@echo ""
	@echo "Utility targets:"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Test Go version (build and run --help)"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Source files: $(GO_SOURCES)"
	@echo "Supported platforms: $(PLATFORMS)"
