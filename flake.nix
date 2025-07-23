{
  description = "Making Mirrors for Git Repositories";
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };
  };
  outputs =
    inputs@{
      self,
      nixpkgs,
      ...
    }:
    let
      system = "x86_64-darwin";
      pkgs = import nixpkgs { inherit system; };

      makingMirrors = pkgs.buildGoModule {
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
      # Default package
      packages.${system} = {
        default = makingMirrors;
        making-mirrors = makingMirrors;
      };

      # Development shell
      devShells.${system}.default = pkgs.mkShell {
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

      # App definition for easy running
      apps.${system}.default = {
        type = "app";
        program = "${makingMirrors}/bin/making-mirrors";
        meta = {
          description = "Run the Making Mirrors application";
        };
      };
    };
}
