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
          version = "0.0.1-alpha";
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
          release = let
            pkgs = import nixpkgs { inherit system; };
          in pkgs.stdenv.mkDerivation {
            pname = "making-mirrors-release";
            version = builtins.replaceStrings ["\n"] [""] (builtins.readFile ./VERSION);
            
            nativeBuildInputs = with pkgs; [ go gnutar gzip coreutils ];
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
              echo "Available commands:"
              echo "  go build       - Build the project"
              echo "  go run .       - Run the project"
              echo "  go test        - Run tests"
              echo "  go mod tidy    - Tidy up dependencies"
              echo "  air            - Live reload development server"
              echo "  nix run .#release - Create cross-platform release"
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
          
          # Release builder app
          release = {
            type = "app";
            program = toString (pkgs.writeShellScript "build-release" ''
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
            '');
            meta = {
              description = "Build cross-platform release of Making Mirrors";
            };
          };
        }
      );
    };
}
