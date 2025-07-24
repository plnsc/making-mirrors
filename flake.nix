{
  description = "Making Mirrors: A Go CLI for mirroring Git repositories";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        go = pkgs.go;
        golangci-lint = pkgs.golangci-lint;
        src = ./.;
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "making-mirrors";
          version = builtins.readFile ./VERSION;
          src = src;
          vendorHash = null;
          subPackages = [ "." ];
          meta = {
            description = "A Go CLI for creating and maintaining mirrors of Git repositories.";
            license = pkgs.lib.licenses.mit;
            maintainers = [ "plnsc" ];
          };
        };
        checks = {
          go-tests = pkgs.buildGoModule {
            pname = "making-mirrors-tests";
            version = builtins.readFile ./VERSION;
            src = src;
            vendorHash = null;
            subPackages = [ "." ];
            doCheck = true;
            checkPhase = "go test ./...";
            meta = {
              description = "Run go tests for making-mirrors";
            };
          };
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/making-mirrors";
          meta = {
            description = "A Go CLI for creating and maintaining mirrors of Git repositories.";
            license = pkgs.lib.licenses.mit;
            maintainers = [ "plnsc" ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            go
            golangci-lint
            pkgs.git
          ];
          shellHook = ''
            echo "Welcome to the Making Mirrors development shell!"
            echo "Go version: $(go version)"
          '';
        };
      }
    );
}
