# Contributing to Making Mirrors

Thank you for your interest in contributing to Making Mirrors! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled (strongly recommended for best experience)
- Git

### Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/plnsc/making-mirrors.git
   cd making-mirrors
   ```

2. **Enter the development environment:**

   ```bash
   nix develop
   ```

   > **Why Nix?** Zero configuration setup, reproducible environment across all contributors, and all development tools pre-installed automatically.

   This will provide you with:

   - Go toolchain (latest stable)
   - Development tools (gopls, golangci-lint, gotools, air)
   - All necessary dependencies
   - Consistent environment across all team members

3. **Verify your setup:**

   ```bash
   nix run .#version
   nix build
   nix flake check
   ```

## Development Workflow

### Recommended Nix-Based Workflow

Using Nix provides the best development experience with zero configuration and perfect reproducibility.

### Step-by-Step Development Flow

1. **Setup environment**: `nix develop` (zero configuration required)
2. **Create feature branch**: `git checkout -b feature/your-feature-name`
3. **Make changes** using Nix development tools
4. **Test changes**: `nix flake check`
5. **Format code**: `nix run .#fmt`
6. **Run linting**: `golangci-lint run`
7. **Build verification**: `nix run .#build`
8. **Commit changes**: `git commit -am 'Add feature'`
9. **Push to branch**: `git push origin feature-name`
10. **Create Pull Request**

### Alternative Manual Flow (Not Recommended)

If you must work without Nix:

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

### Detailed Development Steps

#### Making Changes

1. **Enter development environment** (if not already done):

   ```bash
   nix develop
   ```

2. **Test your changes thoroughly:**

   ```bash
   nix flake check
   nix run .#build
   nix run .#default -- --help  # Test the CLI
   ```

3. **Format and lint your code:**

   ```bash
   go fmt ./...
   golangci-lint run
   ```

4. **Final verification:**

   ```bash
   nix build
   nix run .#default
   ```

### Commit Guidelines

- Use clear, descriptive commit messages
- Follow conventional commit format when possible:
  - `feat: add new feature`
  - `fix: resolve bug`
  - `docs: update documentation`
  - `refactor: improve code structure`
  - `test: add or update tests`

## Code Quality Standards

### Coding Standards

**Go formatting:** Use `go fmt ./...` for consistent code formatting
**Linting:** Code must pass `golangci-lint run` without warnings (golangci-lint included)
**Documentation:** Add comments for exported functions and types
**Error handling:** Always handle errors appropriately
**Testing:** Add tests for new functionality using `nix flake check`

> **Note:** All code quality tools are pre-configured in the Nix environment for consistent results across contributors.

### Code Review Guidelines

## Testing

### Running Tests

```bash
# Run all tests (recommended)
nix flake check

# Advanced testing options (using development shell)
nix develop -c go test -cover ./...  # With coverage
nix develop -c go test -v ./...      # Verbose output
```

> **Recommendation:** Use `nix run .#test` for standard testing as it provides the most consistent environment.

### Writing Tests

- Place test files alongside the code they test (e.g., `main_test.go`)
- Use table-driven tests when testing multiple scenarios
- Mock external dependencies when necessary
- Aim for good test coverage of new functionality

## Documentation

### Code Documentation

- Add package-level documentation to `doc.go`
- Document all exported functions, types, and constants
- Use examples in documentation when helpful

### README Updates

If your changes affect usage or installation:

- Update the README.md file
- Ensure all examples are tested and working
- Update the feature list if applicable

## Pull Request Process

1. **Ensure your code passes all checks:**

   ```bash
   nix flake check
   nix build
   golangci-lint run
   ```

2. **Update documentation** if necessary
3. **Create a pull request** with:

   - Clear title and description
   - Reference any related issues
   - List of changes made
   - Any breaking changes

4. **Respond to review feedback** promptly

## Release Process

### Nix-Based Release Workflow

Releases are managed by maintainers using the Nix build system:

1. **Update version**: `nix run .#set-version x.y.z`
2. **Update CHANGELOG.md** with new version details
3. **Run full test suite**: `nix run .#test`
4. **Create release builds**: `nix run .#release`
5. **Verify build artifacts** in `result-release/`
6. **Commit changes** and create git tag
7. **Push to GitHub** and create release

### Release Artifacts

The automated release system creates:

- **Cross-platform binaries** for all supported architectures
- **SHA256 checksums** for integrity verification
- **Compressed archives** for distribution
- **Organized directory structure** for easy deployment

> **Benefits:** Nix ensures reproducible, cross-platform releases with automated checksums and packaging.

## Getting Help

- **Questions:** Open a GitHub Discussion
- **Bugs:** Open a GitHub Issue
- **Feature Requests:** Open a GitHub Issue with the enhancement label

## Code of Conduct

Be respectful and inclusive. We follow the standard open-source community guidelines:

- Be welcoming to newcomers
- Be respectful of differing viewpoints
- Focus on what is best for the community
- Show empathy towards other community members

Thank you for contributing to Making Mirrors!
