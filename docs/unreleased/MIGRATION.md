# Migration from Make to Nix

This project has been migrated to use Nix instead of Make for build automation. All the functionality from the Makefile has been integrated into `flake.nix`.

## Command Mapping

| Old Make Command                 | New Nix Command               | Description                   |
| -------------------------------- | ----------------------------- | ----------------------------- |
| `make build`                     | `nix run .#build`             | Build the application with Go |
| `make test`                      | `nix run .#test`              | Run tests                     |
| `make clean`                     | `nix run .#clean`             | Clean build artifacts         |
| `make fmt`                       | `nix run .#fmt`               | Format code                   |
| `make lint`                      | `nix run .#lint`              | Run linter                    |
| `make version`                   | `nix run .#version`           | Show current version          |
| `make set-version VERSION=x.y.z` | `nix run .#set-version x.y.z` | Set version in all files      |
| `make install`                   | `nix run .#install`           | Install globally with Nix     |
| `make release`                   | `nix run .#release`           | Create cross-platform release |
| `make dev`                       | `nix develop`                 | Enter development shell       |

## Benefits of the Migration

1. **Reduced Dependencies**: No longer need Make installed on the system
2. **Reproducible Builds**: Nix ensures consistent builds across different environments
3. **Cross-platform**: Works identically on Linux, macOS, and Windows (with WSL)
4. **Isolated Environment**: All dependencies are managed by Nix
5. **Better Tooling**: Access to the full Nix ecosystem

## Quick Start

1. **Enter development environment**:

   ```bash
   nix develop
   ```

2. **Build the project**:

   ```bash
   nix run .#build
   ```

3. **Run tests**:

   ```bash
   nix run .#test
   ```

4. **Create a release**:

   ```bash
   nix run .#release
   ```

## Development Workflow

The development shell (`nix develop`) provides all necessary tools including:

- Go toolchain
- golangci-lint for linting
- air for live reload during development
- All other development dependencies

## Legacy Support

The original Makefile is still present for compatibility, but it's recommended to use the Nix commands for new development.

## Installation

To install making-mirrors globally:

```bash
# Using Nix profile
nix run .#install

# Or directly from the flake
nix profile install .

# Or run without installing
nix run .
```

## Setting Version

To update the version across all files:

```bash
nix run .#set-version 1.0.0
```

This will update:

- VERSION file
- main.go (AppVersion constant)
- flake.nix (version field)
- main_test.go (test expectations)

**Note**: The set-version command has been enhanced with improved regex patterns to ensure all version references are properly updated. See [SET_VERSION_FIX.md](SET_VERSION_FIX.md) for technical details about the implementation and fixes applied.
