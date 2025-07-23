# Makefile for making-mirrors
.PHONY: help build test clean fmt lint nix-build nix-run install release set-version

# Default target
help:
	@echo "Available targets:"
	@echo "  build       - Build the application with Go"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  fmt         - Format code"
	@echo "  lint        - Run linter"
	@echo "  nix-build   - Build with Nix"
	@echo "  nix-run     - Run with Nix"
	@echo "  install     - Install globally with Nix"
	@echo "  release     - Create a release with cross-platform builds"
	@echo "  set-version - Set version in all files (usage: make set-version VERSION=x.y.z)"

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

# Set version in all places (usage: make set-version VERSION=x.y.z)
set-version:
	@if [ -z "$(VERSION)" ]; then \
		echo "Error: VERSION is required. Usage: make set-version VERSION=x.y.z"; \
		exit 1; \
	fi
	@echo "Setting version to $(VERSION) in all files..."
	
	# Update VERSION file
	echo "$(VERSION)" > VERSION
	
	# Update version in main.go using perl for better regex handling
	perl -i -pe 's/AppVersion = "[^"]*"/AppVersion = "$(VERSION)"/' main.go
	
	# Update version in flake.nix
	perl -i -pe 's/version = "[^"]*";/version = "$(VERSION)";/' flake.nix
	
	# Update version in main_test.go
	perl -i -pe 's/\{"AppVersion", AppVersion, "[^"]*"\}/{"AppVersion", AppVersion, "$(VERSION)"}/' main_test.go
	perl -i -pe 's/Version:\s+"[^"]*",/Version:   "$(VERSION)",/' main_test.go
	perl -i -pe 's/info\.Version != "[^"]*"/info.Version != "$(VERSION)"/' main_test.go
	perl -i -pe 's/(Version = %q, want %q", info\.Version, ")[^"]*"/$$1$(VERSION)"/' main_test.go
	
	@echo "Version $(VERSION) has been set in all files"
	@echo "Updated files:"
	@echo "  - VERSION"
	@echo "  - main.go"
	@echo "  - flake.nix"
	@echo "  - main_test.go"

# Create a release with cross-platform builds
release: clean test
	@echo "Creating release for version $(shell cat VERSION)"
	@mkdir -p dist

	# Build for different platforms
	GOOS=linux GOARCH=amd64 go build -o dist/making-mirrors-x86_64-linux -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=linux GOARCH=arm64 go build -o dist/making-mirrors-aarch64-linux -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=darwin GOARCH=amd64 go build -o dist/making-mirrors-x86_64-darwin -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=darwin GOARCH=arm64 go build -o dist/making-mirrors-aarch64-darwin -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=windows GOARCH=amd64 go build -o dist/making-mirrors-windows-amd64.exe -ldflags "-X main.version=$(shell cat VERSION)"
	GOOS=windows GOARCH=arm64 go build -o dist/making-mirrors-windows-arm64.exe -ldflags "-X main.version=$(shell cat VERSION)"

	# Create checksums
	cd dist && sha256sum * > checksums.txt

	# Create tarball
	tar -czf dist/making-mirrors-$(shell cat VERSION).tar.gz -C dist --exclude="*.tar.gz" .

	@echo "Release artifacts created in dist/ directory"
	@echo "Version: $(shell cat VERSION)"
