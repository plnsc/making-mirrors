# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Removed

- **Makefile**: Legacy Make build system removed in favor of Nix-only workflow
- **Make Dependencies**: No longer requires Make to be installed

### Added

- **Nix-based Development Workflow**: Complete migration from Make to Nix for build automation
  **Nix-based Build & Test Commands**: All build and test tasks now use idiomatic Nix commands:
  - `nix build` - Build the application
  - `nix flake check` - Run tests
  - `nix develop` - Enter the development shell
  - `nix profile install` - Install globally with Nix
  - `nix build .#release` - Create cross-platform release
  - `go fmt ./...` - Format code
  - `golangci-lint run` - Run linter
- **Enhanced Development Shell**: Rich development environment with welcome message and command reference
- **Migration Documentation**:
  - `docs/unreleased/MIGRATION.md` - Complete guide for migrating from Make to Nix
  - `docs/unreleased/MIGRATION_SUMMARY.md` - Detailed summary of changes and benefits
  - `docs/unreleased/DOCUMENTATION_UPDATE_SUMMARY.md` - Summary of documentation updates
  - `docs/unreleased/SET_VERSION_FIX.md` - Technical details about set-version implementation fix
- **Zero External Dependencies**: No longer requires Make to be installed
- Support for Gitea, AWS CodeCommit, and Azure Repos

### Changed

- **Primary Build System**: Nix is now the recommended build system (Make remains for compatibility)
- **Development Documentation**: Updated `DEVELOPMENT.md` to prioritize Nix workflow
- **README**: Added migration notes and updated project structure
- **Development Shell**: Enhanced with comprehensive tool listing and usage instructions

### Improved

- **Reproducible Builds**: Nix ensures consistent builds across different environments
- **Cross-Platform Consistency**: Identical development experience on Linux, macOS, and Windows (WSL)
- **Developer Experience**: Integrated toolchain with golangci-lint, air, and Go tools
- **Documentation**: Clear migration path and command mapping from Make to Nix

### Migration Technical Details

- All Makefile functionality preserved as Nix flake apps
- Perl-based version management integrated into Nix scripts
- Development shell includes Go toolchain, linting tools, and live reload
- Backward compatibility maintained with existing Makefile

## [0.0.1-alpha] - 2025-07-23

### Initial Release

- Initial release of making-mirrors
- Command-line interface for mirroring Git repositories
- Concurrent repository processing
- Support for custom input and output directories
- Cross-platform build support
- MIT license

### Core Features

- Read repository registry from file
- Create bare Git mirrors
- Configurable paths with environment variable expansion
- Built-in error handling and logging

### Build & Release Automation

- **Makefile with `set-version` target**: Automatically updates version across all files (main.go, flake.nix, main_test.go, VERSION)
- **Makefile `release` target**: Creates cross-platform builds for all supported architectures
  **Nix release system**: Equivalent cross-platform release build using `nix build .#release`
- **Automated checksums**: Both Makefile and Nix generate SHA256 checksums for all binaries
- **Tarball creation**: Compressed release archives with all artifacts

### Cross-Platform Support

- Linux (x86_64, aarch64)
- macOS/Darwin (x86_64, aarch64)
- Windows (x86_64, aarch64)

### Developer Experience

- Enhanced development shell with release command information
- Consistent version injection across build systems
- Reproducible builds via Nix sandbox environment
- Comprehensive README with usage instructions
- Nix flake configuration for reproducible builds
- Development environment with Go toolchain

### Technical Details

- Makefile `set-version` uses perl for robust regex replacement
- Nix release builds in clean sandbox with proper Go environment variables
- Both systems inject version from VERSION file into binaries via ldflags

[Unreleased]: https://github.com/plnsc/making-mirrors/compare/v0.0.1-alpha...HEAD
[0.0.1-alpha]: https://github.com/plnsc/making-mirrors/releases/tag/v0.0.1-alpha
