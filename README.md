# Making Mirrors for Git Repositories

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

A Rust application for creating mirrors of Git repositories, built and managed with Nix flakes.

## Prerequisites

- [Nix](https://nixos.org/download.html) with flakes enabled
- On macOS: Ensure you have the latest Nix with flakes support

## Quick Start

### 1. Enter the Development Environment

```bash
nix develop
```

This will:

- Install the Rust toolchain (latest stable)
- Provide development tools (cargo-watch, cargo-edit, rust-analyzer)
- Set up the shell with helpful commands

### 2. Build the Application

Using Nix (recommended for production builds):

```bash
nix build
```

Using Cargo (for development):

```bash
cargo build
```

### 3. Run the Application

Using Nix:

```bash
nix run
```

Using Cargo (in development shell):

```bash
cargo run
```

## Development Workflow

### Enter Development Shell

```bash
nix develop
```

Once in the development shell, you have access to:

### Basic Cargo Commands

```bash
cargo build          # Build the project
cargo run            # Run the project
cargo test           # Run tests
cargo check          # Check for errors without building
cargo clippy         # Run linting
cargo fmt            # Format code
```

### Development Tools

```bash
cargo watch -x run   # Auto-rebuild and run on file changes
cargo watch -x test  # Auto-run tests on changes
cargo edit           # Add/remove dependencies easily
```

### Example: Adding Dependencies

```bash
cargo add serde      # Add serde dependency
cargo add --dev tokio-test  # Add development dependency
```

## Building and Distribution

### Build Optimized Binary

```bash
nix build
```

The built binary will be available at `./result/bin/making-mirrors`

### Build for Different Targets

The Nix flake currently targets `x86_64-darwin` (Intel Mac). To build for other systems, modify the `system` variable in `flake.nix`.

### Install Globally

```bash
nix profile install .
```

## Project Structure

```text
making-mirrors/
├── flake.nix          # Nix flake configuration
├── flake.lock         # Locked dependencies
├── Cargo.toml         # Rust package manifest
├── Cargo.lock         # Locked Rust dependencies
├── src/
│   └── main.rs        # Main application code
├── .gitignore         # Git ignore rules
└── README.md          # This file
```

## Nix Flake Features

This flake provides several outputs:

- **`packages.default`**: The built Rust application
- **`packages.making-mirrors`**: Alternative name for the same package
- **`devShells.default`**: Development environment with Rust toolchain
- **`apps.default`**: Direct application runner

### Key Improvements

- ✅ **No Apple SDK warnings**: Uses `libiconv` instead of deprecated framework stubs
- ✅ **Fixed development shell**: No longer exits unexpectedly
- ✅ **Correct cargo hash**: Properly configured for the current project
- ✅ **Clean build**: No deprecation warnings or errors

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

On the first `nix build`, you may see an error about `cargoHash`. This is expected! Nix will show you the correct hash. Copy it and update the `cargoHash` value in `flake.nix`.

### Updating Dependencies

After modifying `Cargo.toml`, you may need to update the `cargoHash` in `flake.nix`:

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
2. Test with `cargo test`
3. Format with `cargo fmt`
4. Lint with `cargo clippy`
5. Build with `nix build` to ensure Nix compatibility
6. Test the development shell with `nix develop`

## What's Working

- ✅ Nix flake builds successfully without warnings
- ✅ Development shell stays open and provides Rust toolchain
- ✅ Application runs with `nix run`
- ✅ Clean build process with proper dependencies
- ✅ macOS compatibility with `libiconv`

## License

MIT License (add your license file as needed)
