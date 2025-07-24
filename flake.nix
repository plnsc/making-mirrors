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
      # For each supported system, set up the environment
      system:
      let
        pkgs = import nixpkgs { inherit system; }; # Import the package set for the current system
        src = ./.; # Project source directory
        description = "A Go CLI for creating and maintaining mirrors of Git repositories.";
      in
      {
        packages.default = pkgs.buildGoModule {
          # The main package: builds the making-mirrors Go binary
          pname = "making-mirrors";
          version = builtins.readFile ./VERSION;
          src = src;
          vendorHash = null;
          subPackages = [ "." ];
          meta = {
            # Package metadata
            description = description;
            license = pkgs.lib.licenses.mit;
            maintainers = [ "plnsc" ];
          };
        };
        checks = {
          # Run Go tests as a Nix check
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
          # Expose the CLI as a Nix flake app
          type = "app";
          program = "${self.packages.${system}.default}/bin/making-mirrors";
          meta = {
            description = description;
            license = pkgs.lib.licenses.mit;
            maintainers = [ "plnsc" ];
          };
        };

        devShells.default = pkgs.mkShell {
          # Development shell with Go, linter, and Git
          buildInputs = [
            pkgs.go
            pkgs.golangci-lint
            pkgs.git
          ];
          shellHook = ''
            # Print a welcome message and Go version on shell startup
            echo "Welcome to the Making Mirrors development shell!"
            echo "Go version: $(go version)"
          '';
        };
      }
    );
}
