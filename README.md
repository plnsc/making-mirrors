# Making Mirrors for Git Repositories

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
![Version](https://img.shields.io/badge/version-0.0.1--alpha-blue)

A Go command-line application for creating and maintaining mirrors of Git repositories. It reads a registry of repositories and creates local bare Git mirrors with concurrent processing for efficient operations.

## Features

- **Concurrent Processing**: Uses all available CPU cores for fast mirroring
- **Multiple Providers**: Supports GitHub, GitLab, and Bitbucket repositories
- **Incremental Updates**: Updates existing mirrors without re-cloning
- **Flexible Configuration**: Customizable input and output directories
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Prerequisites

- [Git](https://git-scm.com/) installed and available in PATH
- [Go](https://golang.org/dl/) 1.21+ (for building from source)
- [Nix](https://nixos.org/download.html) with flakes enabled (optional, for Nix-based workflow)

## Quick Start

### Installation

#### Option 1: Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/plnsc/making-mirrors/releases).

#### Option 2: Install with Go

```bash
go install github.com/plnsc/making-mirrors@latest
```

#### Option 3: Build from Source

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
go build
```

#### Option 4: Using Nix

```bash
nix run github:plnsc/making-mirrors
```

### Basic Usage

1. **Create a registry file** (`registry.txt`) with repositories to mirror:

   ```text
   github:torvalds/linux
   github:golang/go
   gitlab:gitlab-org/gitlab
   bitbucket:atlassian/localstack
   ```

2. **Run the application**:

   ```bash
   making-mirrors
   ```

   Or with custom paths:

   ```bash
   making-mirrors -input ./my-repos.txt -output ./my-mirrors
   ```

### Registry File Format

The registry file contains one repository per line in the format:

```text
provider:owner/repository
```

Supported providers:

- `github` - GitHub repositories
- `gitlab` - GitLab repositories
- `bitbucket` - Bitbucket repositories

Example registry file:

```text
# Core development tools
github:git/git
github:golang/go
github:rust-lang/rust

# Container ecosystem
github:docker/docker
github:kubernetes/kubernetes

# GitLab projects
gitlab:gitlab-org/gitlab
gitlab:gitlab-org/gitaly

# Bitbucket repositories
bitbucket:atlassian/stash
```

Lines starting with `#` are treated as comments and ignored.

## Command Line Options

```text
making-mirrors [flags]

Flags:
  -input string
        Path to the registry file (default "$HOME/Code/mirrors/registry.txt")
  -output string
        Directory to store mirrors (default "$HOME/Code/mirrors")
  -version
        Show version information
```

## Examples

### Mirror to Custom Directory

```bash
making-mirrors -output /backup/git-mirrors
```

### Use Custom Registry File

```bash
making-mirrors -input ./important-repos.txt -output ./mirrors
```

### Environment Variable Expansion

Both input and output paths support environment variable expansion:

```bash
export MIRROR_DIR="/data/mirrors"
making-mirrors -output "$MIRROR_DIR"
```

## How It Works

The application creates **bare Git mirrors** using `git clone --mirror`, which:

- Downloads all branches and tags
- Maintains exact copies of the remote repositories
- Stores repositories in a structured directory format: `provider/owner/repository`
- Supports incremental updates with `git remote update`

### Directory Structure

Mirrors are organized as follows:

```text
mirrors/
├── github/
│   ├── torvalds/
│   │   └── linux/          # Bare Git repository
│   └── golang/
│       └── go/             # Bare Git repository
├── gitlab/
│   └── gitlab-org/
│       └── gitlab/         # Bare Git repository
└── bitbucket/
    └── atlassian/
        └── stash/          # Bare Git repository
```

## Development

### Using Make

The project includes a comprehensive Makefile for development and release management:

```bash
# Available targets
make help          # Show available commands
make build         # Build the application
make test          # Run tests
make clean         # Clean build artifacts
make fmt           # Format code
make lint          # Run linter
make release       # Create cross-platform release
make set-version   # Set version across all files
```

#### Version Management

Update version across all files with a single command:

```bash
make set-version VERSION=1.2.3
```

This automatically updates:

- `VERSION` file
- `main.go` AppVersion constant
- `flake.nix` version field
- `main_test.go` test expectations

#### Cross-Platform Release

Create a complete release with binaries for all platforms:

```bash
make release
```

This generates:

- Binaries for Linux, macOS, and Windows (x86_64 and ARM64)
- SHA256 checksums for all binaries
- Compressed tarball with all artifacts

### Building from Source

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
go mod download
go build
```

### Running Tests

```bash
go test ./...
```

### Using with Nix (Optional)

If you prefer using Nix for development:

#### Enter Development Environment

```bash
nix develop
```

#### Build with Nix

```bash
nix build
```

#### Run with Nix

```bash
nix run
```

#### Nix Release System

Create cross-platform releases using Nix:

```bash
nix run .#release
```

The Nix release system provides:

- Reproducible builds in a clean sandbox
- Automatic dependency management
- Cross-compilation for all supported platforms
- Identical output to the Makefile release system

### Cross-Platform Builds

#### Automated Release (Recommended)

Use the built-in release systems for complete cross-platform builds:

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

## Use Cases

- **Backup Strategy**: Create local backups of important repositories
- **Offline Development**: Work with repositories when internet is limited
- **Repository Analysis**: Bulk analysis of multiple repositories
- **CI/CD Mirroring**: Maintain local copies for build systems
- **Research**: Academic research requiring repository data

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
├── README.md          # This documentation
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Performance

- **Concurrent Processing**: Utilizes all CPU cores for parallel operations
- **Incremental Updates**: Only fetches changes for existing repositories
- **Efficient Storage**: Bare repositories use minimal disk space
- **Progress Tracking**: Real-time status updates during operations

## Troubleshooting

### Common Issues

#### Repository Clone Fails

- Ensure Git is installed and accessible
- Check network connectivity to the Git provider
- Verify repository URLs are correct and accessible

#### Permission Denied

- Ensure write permissions to the output directory
- For private repositories, configure Git credentials

#### Out of Disk Space

- Monitor available disk space before mirroring large repositories
- Consider using a different output directory with more space

### Getting Help

- Check the repository [issues](https://github.com/plnsc/making-mirrors/issues) for known problems
- View debug output by running with verbose Git operations
- Ensure all dependencies (Git, Go) are properly installed

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

### Development Requirements

- Go 1.21 or later
- Git (for testing repository operations)
- Make (optional, for using Makefile commands)

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.

## Author

Paulo Nascimento <paulornasc@gmail.com> - [GitHub](https://github.com/plnsc)
