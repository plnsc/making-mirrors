# Makefile for making-mirrors
.PHONY: help build test clean fmt lint nix-build nix-run install release

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
	@echo "  release   - Create a release with cross-platform builds"

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
	rm -rf dist

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
	@echo "making-mirrors v$(shell cat VERSION)"

# Create a release with cross-platform builds
release: clean test
	@echo "Creating release for version $(shell cat VERSION)"
	@mkdir -p dist

	# Build for different platforms
	GOOS=linux GOARCH=amd64 go build -o dist/making-mirrors-linux-amd64 -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=linux GOARCH=arm64 go build -o dist/making-mirrors-linux-arm64 -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=darwin GOARCH=amd64 go build -o dist/making-mirrors-darwin-amd64 -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=darwin GOARCH=arm64 go build -o dist/making-mirrors-darwin-arm64 -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=windows GOARCH=amd64 go build -o dist/making-mirrors-windows-amd64.exe -ldflags "-X main.version=$(shell cat VERSION)"

	# Create checksums
	cd dist && sha256sum * > checksums.txt

	# Create tarball
	tar -czf dist/making-mirrors-$(shell cat VERSION).tar.gz -C dist --exclude="*.tar.gz" .

	@echo "Release artifacts created in dist/ directory"
	@echo "Version: $(shell cat VERSION)"
