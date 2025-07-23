{
  description = "Making Mirrors for Git Repositories";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };
  outputs =
    inputs@{
      self,
      nixpkgs,
      ...
    }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
        "x86_64-windows"
        "aarch64-windows"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      makingMirrorsForSystem =
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        pkgs.buildGoModule {
          pname = "making-mirrors";
          version = "0.0.2-alpha";
          src = ./.;

          vendorHash = null;

          subPackages = [
            "."
          ];

          ldflags = [
            "-s"
            "-w"
          ];

          meta = with pkgs.lib; {
            description = "Making Mirrors for Git Repositories";
            license = licenses.mit;
            maintainers = [ ];
          };
        };
    in
    {
      # Packages for all supported systems
      packages = forAllSystems (
        system:
        let
          makingMirrors = makingMirrorsForSystem system;
        in
        {
          default = makingMirrors;
          making-mirrors = makingMirrors;

          # Cross-platform release package
          release =
            let
              pkgs = import nixpkgs { inherit system; };
            in
            pkgs.stdenv.mkDerivation {
              pname = "making-mirrors-release";
              version = builtins.replaceStrings [ "\n" ] [ "" ] (builtins.readFile ./VERSION);

              nativeBuildInputs = with pkgs; [
                go
                gnutar
                gzip
                coreutils
              ];
              src = ./.;

              buildPhase = ''
                export GOCACHE=$TMPDIR/go-cache
                export GOPATH=$TMPDIR/go
                export HOME=$TMPDIR

                mkdir -p dist

                # Cross-compile for different platforms
                GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/making-mirrors-x86_64-linux -ldflags "-X main.version=$(cat VERSION)" .
                GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/making-mirrors-aarch64-linux -ldflags "-X main.version=$(cat VERSION)" .
                GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/making-mirrors-x86_64-darwin -ldflags "-X main.version=$(cat VERSION)" .
                GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dist/making-mirrors-aarch64-darwin -ldflags "-X main.version=$(cat VERSION)" .
                GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/making-mirrors-windows-amd64.exe -ldflags "-X main.version=$(cat VERSION)" .
                GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build -o dist/making-mirrors-windows-arm64.exe -ldflags "-X main.version=$(cat VERSION)" .

                # Create checksums
                cd dist && sha256sum * > checksums.txt

                # Create tarball
                cd .. && tar -czf dist/making-mirrors-$(cat VERSION).tar.gz -C dist --exclude="*.tar.gz" .
              '';

              installPhase = ''
                mkdir -p $out
                cp -r dist/* $out/
              '';

              meta = with pkgs.lib; {
                description = "Cross-platform release of Making Mirrors";
                license = licenses.mit;
              };
            };
        }
      );

      # Overlay for easy integration into other flakes
      overlays.default = final: prev: {
        making-mirrors = makingMirrorsForSystem final.stdenv.hostPlatform.system;
      };

      # Development shell for all supported systems
      devShells = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            name = "making-mirrors-dev-shell";
            packages = with pkgs; [
              # Go toolchain
              go

              # Development tools
              gopls
              golangci-lint
              gotools
              air # Live reload for Go apps
            ];

            shellHook = ''
              echo "🐹 Go development environment loaded!"
              echo "Available Nix commands (replacing Make):"
              echo "  nix run .#build         - Build the application with Go"
              echo "  nix run .#test          - Run tests"
              echo "  nix run .#clean         - Clean build artifacts"
              echo "  nix run .#fmt           - Format code"
              echo "  nix run .#lint          - Run linter"
              echo "  nix run .#version       - Show current version"
              echo "  nix run .#set-version   - Set version (usage: nix run .#set-version x.y.z)"
              echo "  nix run .#install       - Install globally with Nix"
              echo "  nix run .#release       - Create cross-platform release"
              echo ""
              echo "Traditional Go commands:"
              echo "  go build       - Build the project"
              echo "  go run .       - Run the project"
              echo "  go test        - Run tests"
              echo "  go mod tidy    - Tidy up dependencies"
              echo ""
              echo "Development tools:"
              echo "  air            - Live reload development server"
              echo "  golangci-lint  - Go linter"
              echo ""
            '';
          };
        }
      );

      # App definition for all systems
      apps = forAllSystems (
        system:
        let
          makingMirrors = makingMirrorsForSystem system;
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = {
            type = "app";
            program = "${makingMirrors}/bin/making-mirrors";
            meta = {
              description = "Run the Making Mirrors application";
            };
          };

          # Build app (equivalent to 'make build')
          build = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "build" ''
                echo "🔨 Building making-mirrors..."
                go build -o making-mirrors
                echo "✅ Build complete: ./making-mirrors"
              ''
            );
            meta = {
              description = "Build the application with Go";
            };
          };

          # Test app (equivalent to 'make test')
          test = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "test" ''
                echo "🧪 Running tests..."
                go test -v ./...
              ''
            );
            meta = {
              description = "Run tests";
            };
          };

          # Clean app (equivalent to 'make clean')
          clean = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "clean" ''
                echo "🧹 Cleaning build artifacts..."
                rm -f making-mirrors
                rm -rf result
                rm -rf result-release
                rm -rf dist
                echo "✅ Clean complete"
              ''
            );
            meta = {
              description = "Clean build artifacts";
            };
          };

          # Format app (equivalent to 'make fmt')
          fmt = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "fmt" ''
                echo "🎨 Formatting code..."
                go fmt ./...
                echo "✅ Code formatting complete"
              ''
            );
            meta = {
              description = "Format code";
            };
          };

          # Lint app (equivalent to 'make lint')
          lint = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "lint" ''
                echo "🔍 Running linter..."
                if command -v golangci-lint >/dev/null 2>&1; then
                  golangci-lint run
                else
                  echo "⚠️  golangci-lint not found. Please install it or use the dev shell."
                  echo "To enter dev shell: nix develop"
                  exit 1
                fi
              ''
            );
            meta = {
              description = "Run linter (requires golangci-lint)";
            };
          };

          # Version app (equivalent to 'make version')
          version = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "version" ''
                if [ -f VERSION ]; then
                  echo "making-mirrors v$(cat VERSION)"
                else
                  echo "VERSION file not found"
                  exit 1
                fi
              ''
            );
            meta = {
              description = "Show version";
            };
          };

          # Set version app (equivalent to 'make set-version VERSION=x.y.z')
          set-version = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "set-version" ''
                if [ -z "$1" ]; then
                  echo "Error: VERSION is required. Usage: nix run .#set-version x.y.z"
                  exit 1
                fi

                VERSION="$1"
                echo "Setting version to $VERSION in all files..."

                # Update VERSION file
                echo "$VERSION" > VERSION

                # Update version in main.go using perl for better regex handling
                ${pkgs.perl}/bin/perl -i -pe "s/AppVersion = \"[^\"]*\"/AppVersion = \"$VERSION\"/" main.go

                # Update version in flake.nix
                ${pkgs.perl}/bin/perl -i -pe "s/version = \"[^\"]*\";/version = \"$VERSION\";/" flake.nix

                # Update version in main_test.go - be more specific to avoid replacing wrong fields
                ${pkgs.perl}/bin/perl -i -pe "s/\\{\"AppVersion\", AppVersion, \"[^\"]*\"\\}/\\{\"AppVersion\", AppVersion, \"$VERSION\"\\}/" main_test.go
                ${pkgs.perl}/bin/perl -i -pe "s/Version:\\s+\"[^\"]*\",/Version:   \"$VERSION\",/" main_test.go
                ${pkgs.perl}/bin/perl -i -pe "s/info\\.Version != \"[^\"]*\"/info.Version != \"$VERSION\"/" main_test.go
                ${pkgs.perl}/bin/perl -i -pe "s/(Version = %q, want %q\", info\\.Version, \")[^\"]*\"(\\))/\$1$VERSION\"\$2/" main_test.go

                echo "Version $VERSION has been set in all files"
                echo "Updated files:"
                echo "  - VERSION"
                echo "  - main.go"
                echo "  - flake.nix"
                echo "  - main_test.go"
              ''
            );
            meta = {
              description = "Set version in all files (usage: nix run .#set-version x.y.z)";
            };
          };

          # Install app (equivalent to 'make install')
          install = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "install" ''
                echo "📦 Installing making-mirrors globally with Nix..."
                nix profile install .
                echo "✅ Installation complete"
              ''
            );
            meta = {
              description = "Install globally with Nix";
            };
          };

          # Release builder app
          release = {
            type = "app";
            program = toString (
              pkgs.writeShellScript "build-release" ''
                echo "🚀 Building cross-platform release with Nix..."

                # Build the release package
                nix build .#release --out-link result-release

                echo "✅ Release build complete!"
                echo "Release artifacts are in: $(readlink -f result-release)"
                echo ""
                echo "Contents:"
                ls -la result-release/
                echo ""
                if [ -f result-release/checksums.txt ]; then
                  echo "Checksums:"
                  cat result-release/checksums.txt
                fi
              ''
            );
            meta = {
              description = "Build cross-platform release of Making Mirrors";
            };
          };
        }
      );
    };
}
