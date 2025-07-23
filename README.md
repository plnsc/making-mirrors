# Making Mirrors for Git Repositories

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)
![Version](https://img.shields.io/badge/version-0.0.1--alpha-blue)

A Go command-line application for creating and maintaining mirrors of Git repositories. It reads a registry of repositories and creates local bare Git mirrors with concurrent processing.

## Features

- **Concurrent Processing**: Uses all available CPU cores for fast mirroring
- **Multiple Providers**: Supports GitHub, GitLab, and Bitbucket repositories
- **Incremental Updates**: Updates existing mirrors without re-cloning
- **Flexible Configuration**: Customizable input and output directories
- **Cross-Platform**: Works on Linux, macOS, and Windows

## Prerequisites

- [Git](https://git-scm.com/) installed and available in PATH
- [Nix](https://nixos.org/download.html) with flakes enabled (recommended for best experience)
- [Go](https://golang.org/dl/) 1.22+ (only if building from source without Nix)

## Quick Start

### Installation

#### Option 1: Using Nix (Recommended)

**Direct run without installation:**

```bash
nix run github:plnsc/making-mirrors
```

**Install globally:**

```bash
nix profile install github:plnsc/making-mirrors
```

**For development:**

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
nix develop  # Enter development environment
nix run .#build  # Build the project
```

> **Why Nix?** Nix provides reproducible builds, zero dependency management, and works identically across all platforms. No need to install Go, Make, or manage toolchains manually.

#### Option 2: Download Pre-built Binary

Download the latest release from the [releases page](https://github.com/plnsc/making-mirrors/releases).

#### Option 3: Install with Go

```bash
go install github.com/plnsc/making-mirrors@latest
```

#### Option 4: Build from Source (Not Recommended)

> **Note:** Building from source requires manual dependency management. Consider using Nix instead for a better experience.

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
go build
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

For development instructions, build automation, cross-platform compilation, and contribution guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md).

**Recommended Development Setup**: Use Nix for the best development experience with zero configuration. See [docs/unreleased/MIGRATION.md](docs/unreleased/MIGRATION.md) for the complete command reference and migration benefits.

### Build System Migration

This project uses Nix as the primary build system, providing superior developer experience compared to traditional approaches:

- **Zero Dependencies**: No need to install Go, Make, or manage toolchains
- **Reproducible Builds**: Identical environments guaranteed across all platforms
- **Rich Development Environment**: Pre-configured with Go, linters, and development tools
- **Cross-Platform Consistency**: Perfect experience on Linux, macOS, and Windows
- **Instant Setup**: Just run `nix develop` and you're ready to contribute

#### Quick Command Reference

| Task        | Nix Command                   |
| ----------- | ----------------------------- |
| Build       | `nix run .#build`             |
| Test        | `nix run .#test`              |
| Clean       | `nix run .#clean`             |
| Format      | `nix run .#fmt`               |
| Lint        | `nix run .#lint`              |
| Version     | `nix run .#version`           |
| Set Version | `nix run .#set-version 1.0.0` |
| Release     | `nix run .#release`           |
| Dev Shell   | `nix develop`                 |

For complete migration details, see [docs/unreleased/MIGRATION.md](docs/unreleased/MIGRATION.md).

## Use Cases

- **Backup Strategy**: Create local backups of important repositories
- **Offline Development**: Work with repositories when internet is limited
- **Repository Analysis**: Bulk analysis of multiple repositories
- **CI/CD Mirroring**: Maintain local copies for build systems
- **Research**: Academic research requiring repository data

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
