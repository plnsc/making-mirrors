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
- [Go](https://golang.org/dl/) 1.22+ (for building from source)
- [Nix](https://nixos.org/download.html) with flakes enabled (optional, for Nix-based workflow)

## Quick Start

### Installation

#### Option 1: Install with Go

```bash
go install github.com/plnsc/making-mirrors@latest
```

#### Option 2: Using Nix

```bash
nix run github:plnsc/making-mirrors
```

#### Option 3: Build from Source

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
go build
```

#### Option 4: Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/plnsc/making-mirrors/releases).

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

For development instructions, build automation, cross-platform compilation, and contribution guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md).

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
├── DEVELOPMENT.md     # Development guide
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

## License

MIT License - see [LICENSE.md](LICENSE.md) for details.

## Author

Paulo Nascimento <paulornasc@gmail.com> - [GitHub](https://github.com/plnsc)
