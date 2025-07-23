{
  description = "Making Mirrors for Git Repositories";
  inputs = {
    nixpkgs = {
      url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    };
    rust-overlay = {
      url = "github:oxalica/rust-overlay";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };
  outputs =
    inputs@{
      self,
      nixpkgs,
      rust-overlay,
      ...
    }:
    let
      system = "x86_64-darwin";
      overlays = [ (import rust-overlay) ];
      pkgs = import nixpkgs { inherit system overlays; };

      # Rust toolchain
      rustToolchain = pkgs.rust-bin.stable.latest.default.override {
        extensions = [
          "rust-src"
          "clippy"
          "rustfmt"
        ];
      };

      # Build the Rust package
      makingMirrors = pkgs.rustPlatform.buildRustPackage {
        pname = "making-mirrors";
        version = "0.1.0";
        src = ./.;

        cargoHash = "sha256-k1I6R35GaiICuWfypFyuY5cMMKlpUMFR1pJVL9HyKXA=";

        nativeBuildInputs = with pkgs; [
          rustToolchain
        ];

        buildInputs =
          with pkgs;
          lib.optionals stdenv.isDarwin [
            libiconv
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
          # Rust toolchain
          rustToolchain

          # Development tools
          cargo-watch
          cargo-edit
          rust-analyzer

          # Build tools
          pkg-config

          # Python (keeping from original)
          (python311.withPackages (pypkgs: [
            pypkgs.pip
            pypkgs.distutils
          ]))
        ];

        buildInputs =
          with pkgs;
          lib.optionals stdenv.isDarwin [
            libiconv
          ];

        shellHook = ''
          echo "🦀 Rust development environment loaded!"
          echo "Available commands:"
          echo "  cargo build    - Build the project"
          echo "  cargo run      - Run the project"
          echo "  cargo test     - Run tests"
          echo "  cargo watch -x run - Auto-rebuild on changes"
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
