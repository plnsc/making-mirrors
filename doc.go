// Package making-mirrors provides utilities for creating mirrors of Git repositories.
//
// # Overview
//
// Making Mirrors is a command-line tool designed to help you create and maintain
// local mirrors of Git repositories from various hosting providers. It's particularly
// useful for backup purposes, offline development, or creating a local cache of
// important repositories.
//
// # Features
//
//   - Concurrent repository mirroring for improved performance
//   - Support for multiple Git hosting providers
//   - Configurable input and output directories
//   - Built with Nix flakes for reproducible builds
//   - Cross-platform support (Linux, macOS, Windows)
//
// # Installation
//
// Using Nix flakes (recommended):
//
//	nix profile install github:plnsc/making-mirrors
//
// Using Go:
//
//	go install github.com/plnsc/making-mirrors@latest
//
// # Usage
//
// Basic usage with default settings:
//
//	making-mirrors
//
// Specify custom input and output directories:
//
//	making-mirrors -input ./my-repos.txt -output ./my-mirrors
//
// # Registry File Format
//
// The registry file should contain repository information that the tool can parse.
// Each line typically represents a repository to be mirrored.
//
// # Development
//
// This project uses Nix flakes for development and building. To get started:
//
//	nix develop
//
// This will provide you with a development shell containing Go toolchain and
// necessary development tools.
//
// # Contributing
//
// 1. Make your changes
// 2. Test with `go test`
// 3. Format with `go fmt`
// 4. Lint with `golangci-lint run`
// 5. Build with `nix build` to ensure Nix compatibility
//
// # License
//
// MIT License - see LICENSE.md for details.
package main
