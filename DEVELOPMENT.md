# Development Guide

This document contains development-specific information for the making-mirrors project.

## Development Environment

### Requirements

- [Nix](https://nixos.org/download.html) with flakes enabled (recommended)
- Go 1.22 or later (if not using Nix)
- Git (for testing repository operations)

### Setting Up

1. Clone the repository:

   ```bash
   git clone https://github.com/plnsc/making-mirrors.git
   cd making-mirrors
   ```

2. **Recommended: Use Nix for development:**

   ```bash
   # Enter development shell with all dependencies
   nix develop
   ```

3. **Alternative: Manual Go setup:**

   ```bash
   # Install dependencies manually
   go mod download
   ```

4. Run tests:

   ```bash
   # Using Nix
   nix run .#test

   # Or using Go directly
   go test ./...
   ```

## Development Workflow

### Using Nix (Recommended)

The project has migrated from Make to Nix for improved reproducibility and cross-platform consistency. All previous Make functionality is now available as Nix apps:

```bash
# Build for current platform
nix run .#build

# Run tests
nix run .#test

# Clean build artifacts
nix run .#clean

# Format code
nix run .#fmt

# Run linter
nix run .#lint

# Show version
nix run .#version

# Set a new version across all files
nix run .#set-version 1.0.0

# Create release builds for all platforms
nix run .#release

# Install globally
nix run .#install

# Enter development shell
nix develop
```

### Legacy Make Support

The project includes a comprehensive Makefile for backward compatibility. See [docs/unreleased/MIGRATION.md](docs/unreleased/MIGRATION.md) for the complete command mapping from Make to Nix.

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

## Migration from Make to Nix

### Implementation Details

The migration from Make to Nix was implemented by:

1. **Converting Makefile targets to Nix apps**: Each Make target became a Nix app in `flake.nix`
2. **Preserving all functionality**: Version management, cross-platform builds, and development workflows
3. **Enhanced tooling**: Added golangci-lint, air, and other development tools to the Nix environment
4. **Improved error handling**: Better error messages and path resolution in Nix scripts

### Benefits Achieved

- **Reproducible Builds**: Nix ensures identical builds across different machines and environments
- **Zero External Dependencies**: No need for Make, Perl, or other system tools to be pre-installed
- **Cross-Platform Consistency**: Same development experience on Linux, macOS, and Windows (WSL)
- **Integrated Development Environment**: All tools available in a single `nix develop` command
- **Better Caching**: Nix's content-addressed storage provides efficient build caching
- **Atomic Operations**: Nix ensures builds either complete successfully or fail cleanly

### Migration Strategy

The migration maintains full backward compatibility:

1. **Dual Support**: Both Make and Nix commands work simultaneously
2. **Gradual Migration**: Teams can adopt Nix incrementally
3. **Documentation**: Complete migration guide in `docs/unreleased/MIGRATION.md`
4. **Command Mapping**: One-to-one mapping between Make and Nix commands

### Development Workflow Improvements

The Nix-based workflow provides several enhancements:

```bash
# Single command to get fully configured environment
nix develop

# Rich welcome message with command reference
# Automatic tool availability (Go, golangci-lint, air)
# Consistent versions across team members

# Enhanced build commands with better output
nix run .#build  # Includes emoji feedback and clear status
nix run .#test   # Better formatted test output
nix run .#clean  # More thorough cleanup including Nix artifacts
```

### Version Management

The project maintains version consistency across multiple files. To update the version:

```bash
# Using Nix (recommended)
nix run .#set-version 1.0.0

# Or using Make (legacy)
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
# Using Nix (recommended)
nix run .#release

# Using Make (legacy)
make release
```

Both methods create binaries for:

- Linux (x86_64, aarch64)
- macOS (x86_64, aarch64)
- Windows (x86_64, aarch64)

Plus checksums and compressed archives.

The Nix method creates a `result-release/` symlink while Make creates a `dist/` directory.

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
