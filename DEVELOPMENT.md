# Development Guide

This document contains development-specific information for the making-mirrors project.

## Project Structure

```text
making-mirrors/
├── main.go            # Main application
├── main_test.go       # Tests
├── go.mod             # Go dependencies
├── docs/              # Documentation
│   └── unreleased/    # Unreleased documentation
│       ├── MIGRATION.md  # Make to Nix migration guide
│       ├── MIGRATION_SUMMARY.md  # Migration summary
│       ├── DOCUMENTATION_UPDATE_SUMMARY.md  # Doc updates
│       └── SET_VERSION_FIX.md  # Set-version implementation fix
├── flake.nix          # Nix flake (build automation)
├── CONTRIBUTING.md    # Contribution guidelines
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Development Environment

### Requirements

- [Nix](https://nixos.org/download.html) with flakes enabled (strongly recommended)
- Git (for testing repository operations)
- Go 1.22 or later (only if not using Nix - not recommended)

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

   > **Why Nix?** Zero configuration, reproducible environment, all tools included automatically.

3. **Alternative: Manual Go setup (not recommended):**

   > **Note:** Manual setup requires dependency management and may lead to inconsistent environments.

   ```bash
   # Install dependencies manually
   go mod download
   ```

4. Run tests:

   ```bash
   # Using Nix (recommended)
   nix flake check

   # Manual Go testing (if not using Nix)
   go test ./...
   ```

## Development Workflow

### Using Nix (Strongly Recommended)

This project uses Nix as the primary build system for superior developer experience. All development tasks are available using idiomatic Nix commands:

```bash
 # Build for current platform
 nix build

 # Run tests
 nix flake check

 # Enter development shell
 nix develop

 # Install globally
 nix profile install



 # Format code (Go)
 go fmt ./...

 # Run linter (Go)
 golangci-lint run

# Install globally
nix run .#install

# Enter development shell
nix develop
```

## Nix-Based Development Benefits

### Why Choose Nix for Development?

Using Nix provides several advantages over traditional development approaches:

- **Zero Configuration**: Run `nix develop` and everything is ready
- **Reproducible Environments**: Identical setup across all team members and CI/CD
- **No Dependency Hell**: All tools and versions managed automatically
- **Cross-Platform Consistency**: Same experience on Linux, macOS, and Windows
- **Instant Onboarding**: New developers are productive immediately
- **Automated Tooling**: Pre-configured linting, formatting, and development tools

### Enhanced Development Experience

The Nix development environment provides a superior workflow:

```bash
# Single command to get fully configured environment
nix develop

# Rich welcome message with command reference
# Automatic tool availability (Go, golangci-lint, air)
# Consistent versions across team members

# Enhanced build commands with clear feedback
nix run .#build  # Includes status indicators and clear output
nix run .#test   # Better formatted test results
nix run .#clean  # Thorough cleanup including build artifacts
```

### Version Management

The project maintains version consistency across multiple files. To update the version:

```bash
# Using Nix (recommended approach)
nix run .#set-version 1.0.0
```

This command automatically updates:

- `VERSION` file
- `main.go` version constant
- `flake.nix` version field
- `main_test.go` test expectations

### Cross-Platform Builds

#### Automated Release Builds

Create binaries for all supported platforms:

```bash
# Using Nix (recommended and only supported method)
nix run .#release
```

This creates binaries for:

- Linux (x86_64, aarch64)
- macOS (x86_64, aarch64)
- Windows (x86_64, aarch64)

Plus checksums and compressed archives in the `result-release/` directory.

#### Manual Cross-Compilation

If you need individual platform builds, use Nix for consistent results:

```bash
nix build .#packages.x86_64-linux.default    # Intel/AMD Linux
nix build .#packages.aarch64-linux.default   # ARM64 Linux
nix build .#packages.x86_64-darwin.default   # Intel Mac
nix build .#packages.aarch64-darwin.default  # Apple Silicon Mac
```

> **Note:** While Go's built-in cross-compilation is available, Nix provides better reproducibility and dependency management.

## Nix Development

### Nix Flake Features

The project includes a comprehensive Nix flake providing:

- **Cross-platform build support** for all major architectures
- **Rich development shell** with all dependencies pre-configured
- **Automated release system** with checksums and packaging
- **Reproducible builds** guaranteed across environments
- **Zero external dependencies** - everything managed by Nix

### Development Commands

```bash
# Enter development shell (recommended first step)
nix develop

# Build for current platform
nix build

# Run the application directly
nix run

# Create comprehensive release packages
nix run .#release

# Run specific development tasks
nix run .#test     # Run test suite
nix run .#fmt      # Format code
nix run .#lint     # Run linter
nix run .#clean    # Clean artifacts
```

## Testing

### Running Tests

```bash
# Recommended: Using Nix
nix run .#test

# Alternative: Direct Go testing (if in nix develop shell)
go test ./...

# Verbose output (using development shell)
nix develop -c go test -v ./...

# Coverage analysis (using development shell)
nix develop -c go test -cover ./...
```

### Test Structure

- `main_test.go`: Contains unit tests for the main application
- Tests validate version consistency and core functionality
- Version tests ensure the `GoVersion` constant matches the VERSION file

## Code Quality

### Formatting

```bash
# Recommended: Using Nix
nix run .#fmt

# Alternative: Direct Go formatting (if in development shell)
go fmt ./...
```

### Linting

```bash
# Recommended: Using Nix (includes golangci-lint automatically)
nix run .#lint

# Alternative: Manual linting (if golangci-lint is available)
golangci-lint run
```

## Project Structure

```text
making-mirrors/
├── main.go            # Main application code
├── main_test.go       # Tests
├── go.mod             # Go module definition
├── go.sum             # Go module checksums
├── flake.nix          # Nix flake configuration with comprehensive build system
├── flake.lock         # Nix dependencies (locked versions)
├── LICENSE.md         # MIT license
├── README.md          # User documentation
├── DEVELOPMENT.md     # This development guide
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Troubleshooting Development Issues

### Build Issues

- **Recommended**: Use `nix develop` to ensure all dependencies are available
- For manual setups: Ensure Go 1.22+ is installed and run `go mod download`
- Check that Git is available in PATH for repository operations

### Cross-Platform Build Issues

- Use `nix run .#release` for consistent cross-platform builds
- Ensure sufficient disk space for all platform binaries
- For Nix: Ensure Nix is properly installed with flakes enabled
