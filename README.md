# Making Mirrors: Create mirrors of Git repositories

![GitHub Tag](https://img.shields.io/github/v/tag/plnsc/making-mirrors)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/plnsc/making-mirrors/ci.yml?label=build)
[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

Making Mirrors is a Go command-line tool for creating and maintaining local copies of Git repositories. It does so by using `git clone --mirror` to get local bare Git mirrors of remote repositories in well known providers.

Own local stored mirrors of what is important or most used. Be able to manage a curated list of mirrors for resource limited storage. It provides a copy of interest-specific repositories to create a layer of resilience and increase availability in software deployment.

## Features

- **Concurrent Processing**: Uses all available CPU cores for fast mirroring
- **Multiple Providers**: Supports GitHub, GitLab, Bitbucket, Gitea, AWS CodeCommit, and Azure Repos
- **Incremental Updates**: Updates existing mirrors without re-cloning
- **Flexible Configuration**: Customizable input and output directories
- **Cross-Platform**: Works on Linux, macOS, and Windows

### Future

- Read-only host capabilities enabled. Example: Serve the repos in equivalent servers like `https://unofficial-local-github-mirror/torvalds/linux.git`.
- Accept plain URL as repository input. Currently only the short format is accepted.

### Known issues

- Cloning big repositories\* is a work in progress. \* (Like the ones in the examples)

### Usage

1. **Create a registry file** with repositories to mirror (Default: `~/Code/mirrors/registry.txt`):

   ```text
   github:golang/go
   github:NixOS/nix
   github:NixOS/nixpkgs
   github:torvalds/linux
   gitlab:gitlab-org/gitlab
   ```

2. **Run the application**:

   ```bash
   making-mirrors
   ```

   Or with custom paths:

   > **Note:** Make sure the directories and the registry file exist.

   ```bash
   making-mirrors -input ./Repos/registry.txt -output ./Repos/mirrors
   ```

3. **Artifacts**:

   This will pull the repositories to a default `~/Code/mirrors`. See **Directory structure** for more information.

#### Install with Go

```bash
go install github.com/plnsc/making-mirrors@latest
```

#### Install with Nix

```bash
nix profile install github:plnsc/making-mirrors
```

#### Build from source

> **Note:** Consider using Nix instead for a consistent experience.

```bash
git clone https://github.com/plnsc/making-mirrors.git
cd making-mirrors
go build # or `nix build`
```

## How It Works

The application creates **bare Git mirrors** using `git clone --mirror`, which:

- Downloads all branches and tags
- Maintains exact copies of the remote repositories
- Stores repositories in a structured directory format: `provider/owner/repository`
- Supports incremental updates with `git remote update`

### Command Line Options

```text
making-mirrors [flags]

Flags:
  -input string
        Path to the registry file (default "$HOME/Code/mirrors/registry.txt")
  -output string
        Directory to store mirrors (default "$HOME/Code/mirrors")
  -version
        Show version information
```

### Registry file format

The registry file consists a text file that contains one repository per line. The repositories are written in a short format so the software can expand it to the right targets.

Format:

```text
provider:owner/repository
```

Currently supported providers:

- `github` - GitHub repositories
- `gitlab` - GitLab repositories
- `bitbucket` - Bitbucket repositories
- `gitea` - Gitea repositories
- `codecommit` - AWS CodeCommit repositories (e.g. `codecommit:us-west-2/myrepo`)
- `azure` - Azure Repos (e.g. `azure:org/project`)

Example registry file:

```text
# Core development tools
github:git/git
github:golang/go
github:rust-lang/rust

# Container ecosystem
github:docker/docker
github:kubernetes/kubernetes

# GitLab projects
gitlab:gitlab-org/gitlab
gitlab:gitlab-org/gitaly

# Bitbucket repositories
bitbucket:atlassian/stash

# Gitea repositories
gitea:john/doerepo

# AWS CodeCommit
codecommit:us-west-2/myrepo

# Azure Repos
azure:myorg/myproject
```

Lines starting with `#` are treated as comments and ignored.

### Directory structure

Mirrors are organized as follows:

```text
mirrors/
├── github/
│   ├── torvalds/
│   │   └── linux/          # Bare Git repository
│   └── golang/
│       └── go/             # Bare Git repository
├── gitlab/
│   └── gitlab-org/
│       └── gitlab/         # Bare Git repository
└── bitbucket/
    └── atlassian/
        └── stash/          # Bare Git repository
```

## Development

For development instructions, build automation, cross-platform compilation, and contribution guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md).

### Quick Command Reference

| Task      | Nix Command                        |
| --------- | ---------------------------------- |
| Build     | `nix build`                        |
| Test      | `nix flake check`                  |
| Dev Shell | `nix develop`                      |
| Format    | `nix develop -c go fmt ./...`      |
| Lint      | `nix develop -c golangci-lint run` |
| Install   | `nix profile install`              |

## Use Cases

- **Backup Strategy**: Create local backups of important repositories
- **Offline Development**: Work with repositories when internet is limited
- **Repository Analysis**: Bulk analysis of multiple repositories
- **CI/CD Mirroring**: Maintain local copies for build systems
- **Research**: Academic research requiring repository data

## Project Structure

```text
making-mirrors/
├── main.go            # Main application
├── main_test.go       # Tests
├── go.mod             # Go dependencies
├── docs/              # Documentation
│   └── unreleased/    # Unreleased documentation
│       ├── MIGRATION.md  # Make to Nix migration guide
│       ├── MIGRATION_SUMMARY.md  # Migration summary
│       ├── DOCUMENTATION_UPDATE_SUMMARY.md  # Doc updates
│       └── SET_VERSION_FIX.md  # Set-version implementation fix
├── flake.nix          # Nix flake (build automation)
├── CONTRIBUTING.md    # Contribution guidelines
├── CHANGELOG.md       # Version history
└── VERSION            # Current version (0.0.1-alpha)
```

## Troubleshooting

### Common Issues

#### Repository Clone Fails

- Ensure Git is installed and accessible
- Check network connectivity to the Git provider
- Verify repository URLs are correct and accessible

#### Permission Denied

- Ensure write permissions to the output directory
- For private repositories, configure Git credentials

#### Out of Disk Space

- Monitor available disk space before mirroring large repositories
- Consider using a different output directory with more space

### Getting Help

- Check the repository [issues](https://github.com/plnsc/making-mirrors/issues) for known problems
- View debug output by running with verbose Git operations
- Ensure all dependencies (Git, Go) are properly installed

## Author

Paulo Nascimento. [GitHub](https://github.com/plnsc). [Personal Blog](https://taboza.dev)

## License

 [MIT License](LICENSE.md)
