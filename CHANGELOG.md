# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial package documentation and metadata
- Comprehensive README with usage instructions
- Nix flake configuration for reproducible builds
- Development environment with Go toolchain

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
- **Nix release system**: Equivalent cross-platform release build using `nix run .#release`
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
