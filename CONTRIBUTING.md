# Contributing to Making Mirrors

Thank you for your interest in contributing to Making Mirrors! This document provides guidelines and information for contributors.

## Development Setup

### Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
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

   This will provide you with:

   - Go toolchain (latest stable)
   - Development tools (gopls, golangci-lint, gotools, air)
   - All necessary dependencies

3. **Verify your setup:**
   ```bash
   go version
   go test
   ```

## Development Workflow

### Making Changes

1. **Create a new branch:**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following the coding standards below

3. **Test your changes:**

   ```bash
   go test ./...
   go build
   ./making-mirrors --help  # Test the CLI
   ```

4. **Format and lint:**

   ```bash
   go fmt ./...
   golangci-lint run
   ```

5. **Test with Nix:**
   ```bash
   nix build
   nix run
   ```

### Coding Standards

- **Go formatting:** Use `go fmt` to format your code
- **Linting:** Code must pass `golangci-lint run` without warnings
- **Documentation:** Add comments for exported functions and types
- **Error handling:** Always handle errors appropriately
- **Testing:** Add tests for new functionality

### Commit Guidelines

- Use clear, descriptive commit messages
- Follow conventional commit format when possible:
  - `feat: add new feature`
  - `fix: resolve bug`
  - `docs: update documentation`
  - `refactor: improve code structure`
  - `test: add or update tests`

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests in verbose mode
go test -v ./...
```

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
   go test ./...
   go fmt ./...
   golangci-lint run
   nix build
   ```

2. **Update documentation** if necessary

3. **Create a pull request** with:

   - Clear title and description
   - Reference any related issues
   - List of changes made
   - Any breaking changes

4. **Respond to review feedback** promptly

## Release Process

Releases are managed by maintainers. The process includes:

1. Update `VERSION` file
2. Update `CHANGELOG.md`
3. Create a git tag
4. Build and test across all supported platforms
5. Create a GitHub release

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
