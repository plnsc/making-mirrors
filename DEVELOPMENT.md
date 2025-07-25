# Development Guide

I'm keeping only the [README.md](README.md) and this [DEVELOPMENT.md](DEVELOPMENT.md) to hold documentation.

## Project structure

```text
making-mirrors/
├── main.go            # Main application
├── main_test.go       # Tests
├── go.mod             # Go dependencies
├── flake.nix          # Nix flake (dev environment)
├── flake.lock         # Nix flake lock file
├── CONTRIBUTING.md    # Contribution guidelines
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Development environment

### Requirements

- [Nix](https://nixos.org/download.html) with flakes enabled
- Go 1.22 or later (only if not using Nix)
- Git

### Setting up

1. Clone the repository:

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
```

2. **Use Nix for development (Recommended):**

| Task      | Nix Command                        |
| --------- | ---------------------------------- |
| Build     | `nix build`                        |
| Test      | `nix flake check`                  |
| Dev Shell | `nix develop`                      |
| Format    | `nix develop -c go fmt ./...`      |
| Lint      | `nix develop -c golangci-lint run` |
| Install   | `nix profile install`              |

> **Why Nix?** Zero configuration, reproducible environment, all tools included automatically. **Also** I use Nix in the definition of my home environment, so I need this to be plug-and-play in these "production nodes".

3. **Manual Go setup (Alternative):**

> **Note:** Manual setup requires dependency management and may lead to inconsistent environments.

| Task   | Nix Command         |
| ------ | ------------------- |
| Build  | `go build`          |
| Test   | `go test ./...`     |
| Format | `go fmt ./...`      |
| Lint   | `golangci-lint run` |

## Contributing

It would be awesome to hear that you using this and want to contribute. So open an [issue](https://github.com/plnsc/making-mirrors/issues) with any question or problem and lets start there. You can also follow me on [GitHub](https://github.com/plnsc) and visit my freshly started [personal blog](https://taboza.dev).