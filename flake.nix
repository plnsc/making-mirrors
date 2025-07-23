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
          version = "0.1.0";
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
        in
        {
          default = {
            type = "app";
            program = "${makingMirrors}/bin/making-mirrors";
            meta = {
              description = "Run the Making Mirrors application";
            };
          };
        }
      );
    };
}
