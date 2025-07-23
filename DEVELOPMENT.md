# Development Guide

This document contains development-specific information for the making-mirrors project.

## Development Environment

### Requirements

- Go 1.22 or later
- Git (for testing repository operations)
- Make (optional, for using Makefile commands)
- Nix (optional, for reproducible builds)

### Setting Up

1. Clone the repository:

   ```bash
   git clone https://github.com/plnsc/making-mirrors.git
   cd making-mirrors
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run tests:

   ```bash
   go test ./...
   ```

## Development Workflow

### Using Make

The project includes a comprehensive Makefile for build automation:

```bash
# Build for current platform
make build

# Run tests
make test

# Clean build artifacts
make clean

# Set a new version across all files
make set-version VERSION=1.0.0

# Create release builds for all platforms
make release

# Show available targets
make help
```

### Version Management

The project maintains version consistency across multiple files. To update the version:

```bash
# Update version in all relevant files
make set-version VERSION=1.0.0
```

This command updates:

- `VERSION` file
- `main.go` version constant
- `flake.nix` version field
- `main_test.go` test expectations

The version management uses Perl for robust regex replacements across different file formats.

### Cross-Platform Builds

#### Automated Release Builds

Create binaries for all supported platforms:

```bash
# Using Make (creates dist/ directory with all platforms)
make release

# Using Nix (creates result-release/ symlink with all platforms)
nix run .#release
```

Both methods create binaries for:

- Linux (x86_64, aarch64)
- macOS (x86_64, aarch64)
- Windows (x86_64, aarch64)

Plus checksums and compressed archives.

#### Manual Cross-Compilation

If you need individual platform builds:

##### Using Go

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o making-mirrors-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o making-mirrors-windows-amd64.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o making-mirrors-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o making-mirrors-darwin-arm64
```

##### Using Nix

```bash
nix build .#packages.x86_64-linux.default    # Intel/AMD Linux
nix build .#packages.aarch64-linux.default   # ARM64 Linux
nix build .#packages.x86_64-darwin.default   # Intel Mac
nix build .#packages.aarch64-darwin.default  # Apple Silicon Mac
```

## Nix Development

### Nix Flake Features

The project includes a comprehensive Nix flake with:

- Cross-platform build support
- Development shell with all dependencies
- Automated release system
- Reproducible builds

### Using Nix for Development

```bash
# Enter development shell
nix develop

# Build for current platform
nix build

# Run the application
nix run

# Create release packages
nix run .#release
```

### Nix Development Shell

The development shell includes:

- Go compiler and tools
- Git for version control
- Make for build automation
- All project dependencies

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Test Structure

- `main_test.go`: Contains unit tests for the main application
- Tests validate version consistency and core functionality
- Version tests ensure the `GoVersion` constant matches the VERSION file

## Code Quality

### Formatting

```bash
# Format all Go code
go fmt ./...
```

### Linting

If golangci-lint is available:

```bash
golangci-lint run
```

## Project Structure

```text
making-mirrors/
├── main.go            # Main application code
├── main_test.go       # Tests
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── flake.nix          # Nix flake configuration with release system
├── flake.lock         # Nix dependencies
├── Makefile           # Build automation with release targets
├── LICENSE.md         # MIT license
├── README.md          # User documentation
├── DEVELOPMENT.md     # This development guide
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Troubleshooting Development Issues

### Build Issues

- Ensure Go 1.22+ is installed
- Run `go mod download` to fetch dependencies
- Check that Git is available in PATH

### Version Management Issues

- The `set-version` target requires Perl for regex operations
- Ensure all target files exist before running version updates
- Check file permissions if updates fail

### Cross-Platform Build Issues

- Ensure sufficient disk space for all platform binaries
- For Nix builds, ensure Nix is properly installed and flakes are enabled

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for new functionality
5. Run tests: `go test ./...`
6. Format code: `go fmt ./...`
7. Run linting (if available): `golangci-lint run`
8. Commit changes: `git commit -am 'Add feature'`
9. Push to branch: `git push origin feature-name`
10. Create a Pull Request

### Code Review Guidelines

- Ensure all tests pass
- Add appropriate documentation
- Follow Go best practices
- Update CHANGELOG.md for significant changes
- Update version if needed using `make set-version`

## Release Process

1. Update version: `make set-version VERSION=x.y.z`
2. Update CHANGELOG.md with new version details
3. Run tests: `go test ./...`
4. Create release builds: `make release`
5. Commit changes and create git tag
6. Push to GitHub and create release

The automated release system creates:

- Binaries for all supported platforms
- SHA256 checksums
- Compressed archives
- Release directory structure
