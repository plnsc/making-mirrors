# Making Mirrors for Git Repositories

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

A Go application for creating mirrors of Git repositories, built and managed with Nix flakes. This flake is designed to be easily used as an input in other Nix projects.

## Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
- On macOS: Ensure you have the latest Nix with flakes support

## Quick Start

### Using as a Flake Input

To use this project in another Nix flake:

```nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    making-mirrors.url = "github:plnsc/making-mirrors";
  };

  outputs = { nixpkgs, making-mirrors, ... }: {
    # Option 1: Use the package directly
    packages.x86_64-linux.my-package = making-mirrors.packages.x86_64-linux.default;

    # Option 2: Use the overlay
    packages.x86_64-linux.my-env = let
      pkgs = import nixpkgs {
        system = "x86_64-linux";
        overlays = [ making-mirrors.overlays.default ];
      };
    in pkgs.making-mirrors;
  };
}
```

### Local Development

### 1. Enter the Development Environment

```bash
nix develop
```

This will:

- Install the Go toolchain (latest stable)
- Provide development tools (gopls, golangci-lint, gotools, air)
- Set up the shell with helpful commands

### 2. Build the Application

Using Nix (recommended for production builds):

```bash
nix build
```

Using Go (for development):

```bash
go build
```

### 3. Run the Application

Using Nix:

```bash
nix run
```

Using Go (in development shell):

```bash
go run .
```

## Development Workflow

### Enter Development Shell

```bash
nix develop
```

Once in the development shell, you have access to:

### Basic Go Commands

```bash
go build            # Build the project
go run .            # Run the project
go test             # Run tests
go mod tidy         # Clean up dependencies
go fmt              # Format code
go vet              # Examine code for issues
```

### Development Tools

```bash
air                 # Live reload development server
golangci-lint run   # Run comprehensive linting
gopls               # Language server (integrated with editors)
```

### Example: Adding Dependencies

```bash
go get github.com/gin-gonic/gin  # Add a dependency
go mod tidy                      # Clean up go.mod and go.sum
```

## Building and Distribution

### Build Optimized Binary

```bash
nix build
```

The built binary will be available at `./result/bin/making-mirrors`

### Build for Different Targets

The Nix flake supports multiple systems:

- `x86_64-linux` (Intel/AMD Linux)
- `aarch64-linux` (ARM64 Linux)
- `x86_64-darwin` (Intel Mac)
- `aarch64-darwin` (Apple Silicon Mac)

To build for a specific system:

```bash
nix build .#packages.x86_64-linux.default    # Build for Intel/AMD Linux
nix build .#packages.aarch64-linux.default   # Build for ARM64 Linux
nix build .#packages.x86_64-darwin.default   # Build for Intel Mac
nix build .#packages.aarch64-darwin.default  # Build for Apple Silicon Mac
```

You can also use Go's built-in cross-compilation:

```bash
GOOS=linux GOARCH=amd64 go build .     # Build for Linux
GOOS=windows GOARCH=amd64 go build .   # Build for Windows
```

### Install Globally

```bash
nix profile install .
```

## Project Structure

```text
making-mirrors/
├── flake.nix          # Nix flake configuration
├── flake.lock         # Locked dependencies
├── go.mod             # Go module definition
├── go.sum             # Go module checksums (auto-generated)
├── main.go            # Main application code
├── .gitignore         # Git ignore rules
└── README.md          # This file
```

## Nix Flake Features

This flake provides several outputs:

- **`packages.default`**: The built Go application
- **`packages.making-mirrors`**: Alternative name for the same package
- **`devShells.default`**: Development environment with Go toolchain
- **`apps.default`**: Direct application runner

### Key Features

- ✅ **Go toolchain**: Latest stable Go compiler and tools
- ✅ **Development tools**: gopls, golangci-lint, air for live reload
- ✅ **Cross-platform builds**: Support for multiple architectures
- ✅ **Clean development environment**: Isolated and reproducible

### Using Different Outputs

```bash
nix build .#making-mirrors    # Build the package
nix develop .#default         # Enter dev shell
nix run .#default            # Run the application
```

## Troubleshooting

### Development Shell Issues

If the development shell exits immediately, this has been fixed in the current version. The shell should now stay open and display the welcome message.

### First Build Issues

On the first `nix build`, you may see an error about `vendorHash`. This is expected if you add dependencies! Nix will show you the correct hash. Copy it and update the `vendorHash` value in `flake.nix`.

### Updating Dependencies

After modifying `go.mod`, you may need to update the `vendorHash` in `flake.nix`:

1. Delete the current hash (set it to an empty string or wrong hash)
2. Run `nix build`
3. Copy the correct hash from the error message
4. Update `flake.nix` with the new hash

### Apple Silicon Macs

If you're on Apple Silicon (M1/M2/M3), you may want to change the `system` in `flake.nix` from `x86_64-darwin` to `aarch64-darwin` for optimal performance.

### Updating Flake Inputs

```bash
nix flake update    # Update all inputs
nix flake lock      # Update lock file
```

## Contributing

1. Make your changes
2. Test with `go test`
3. Format with `go fmt`
4. Lint with `golangci-lint run`
5. Build with `nix build` to ensure Nix compatibility
6. Test the development shell with `nix develop`

## What's Working

- ✅ Nix flake builds successfully without warnings
- ✅ Development shell stays open and provides Go toolchain
- ✅ Application runs with `nix run`
- ✅ Clean build process with proper dependencies
- ✅ macOS compatibility with `libiconv`

## License

MIT License (add your license file as needed)
