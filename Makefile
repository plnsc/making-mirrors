# Makefile for making-mirrors
.PHONY: help build test clean fmt lint nix-build nix-run install

# Default target
help:
	@echo "Available targets:"
	@echo "  build     - Build the application with Go"
	@echo "  test      - Run tests"
	@echo "  clean     - Clean build artifacts"
	@echo "  fmt       - Format code"
	@echo "  lint      - Run linter"
	@echo "  nix-build - Build with Nix"
	@echo "  nix-run   - Run with Nix"
	@echo "  install   - Install globally with Nix"

# Build with Go
build:
	go build -o making-mirrors

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f making-mirrors
	rm -rf result

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Build with Nix
nix-build:
	nix build

# Run with Nix
nix-run:
	nix run

# Install globally with Nix
install:
	nix profile install .

# Development shell
dev:
	nix develop

# Show version
version:
	@echo "making-mirrors v0.1.0"
