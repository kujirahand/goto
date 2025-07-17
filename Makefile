# Makefile for goto command

.PHONY: all build-go install-go clean test help

# Default target
all: build-go

# Build Go version
build-go:
	@echo "Building Go version..."
	cd go && go build -o goto goto.go config.go
	@echo "✅ Go version built successfully: go/goto"

# Install Go version to /usr/local/bin
install-go: build-go
	@echo "Installing Go version to /usr/local/bin/goto..."
	sudo cp go/goto /usr/local/bin/goto
	@echo "✅ Go version installed as 'goto'"

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
	@echo "  all          - Build Go version (default)"
	@echo "  build-go     - Build Go version"
	@echo "  install-go   - Install Go version to /usr/local/bin"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Test Go version"
	@echo "  help         - Show this help"
