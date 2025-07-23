# Making Mirrors for Git Repositories

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

A Rust application built and managed with Nix flakes.

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
- Keep Python support from the original environment

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
├── src/
│   └── main.rs        # Main application code
├── .gitignore         # Git ignore rules
└── README.md          # This file
```

## Nix Flake Outputs

This flake provides several outputs:

- **`packages.default`**: The built Rust application
- **`devShells.default`**: Development environment with Rust toolchain
- **`apps.default`**: Direct application runner

### Using Different Outputs

```bash
nix build .#making-mirrors    # Build the package
nix develop .#default         # Enter dev shell
nix run .#default            # Run the application
```

## Troubleshooting

### First Build Issues

On the first `nix build`, you may see an error about `cargoHash`. This is expected! Nix will show you the correct hash. Copy it and update the `cargoHash` value in `flake.nix`.

### Updating Dependencies

After modifying `Cargo.toml`, you may need to update the `cargoHash` in `flake.nix`:

1. Delete the current hash (set it to an empty string or wrong hash)
2. Run `nix build`
3. Copy the correct hash from the error message
4. Update `flake.nix` with the new hash

### Shell Issues

If the development shell doesn't start properly, try:

```bash
nix develop --impure
```

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

## License

MIT License (add your license file as needed)
