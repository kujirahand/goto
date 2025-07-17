# Makefile for goto command

.PHONY: all build-go build-python install-python install-go clean test help

# Default target
all: build-go

# Build Go version
build-go:
	@echo "Building Go version..."
	cd go && go build -o goto goto.go config.go
	@echo "✅ Go version built successfully: go/goto"

# Install Python dependencies
install-python:
	@echo "Installing Python dependencies..."
	cd bin && pip3 install --break-system-packages -r requirements.txt
	@echo "✅ Python dependencies installed"

# Install Go version to /usr/local/bin
install-go: build-go
	@echo "Installing Go version to /usr/local/bin/goto-go..."
	sudo cp go/goto /usr/local/bin/goto-go
	@echo "✅ Go version installed as 'goto-go'"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f go/goto
	@echo "✅ Clean completed"

# Test both versions
test: build-go
	@echo "Testing Python version..."
	cd bin && python3 goto.py --help
	@echo "\nTesting Go version..."
	cd go && ./goto --help
	@echo "✅ Both versions working"

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Build Go version (default)"
	@echo "  build-go     - Build Go version"
	@echo "  install-python - Install Python dependencies"
	@echo "  install-go   - Install Go version to /usr/local/bin"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Test both versions"
	@echo "  help         - Show this help"
