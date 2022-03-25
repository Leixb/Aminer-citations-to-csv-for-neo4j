{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
    };

  };

  outputs = { self, nixpkgs, flake-utils, }:

  flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = (import nixpkgs { inherit system; });
    in
    rec {
      devShell = pkgs.mkShellNoCC {

        buildInputs = with pkgs; [
          go_1_18
          gopls
          gotools
        ];

      };

      packages.parser = pkgs.buildGo118Module {
        pname = "parser";
        version = "1";
        src = ./.;
        vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
      };

      defaultPackage = packages.parser;
    }
  );
}
