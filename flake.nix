{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-21.11";

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
    {
      devShell = pkgs.mkShellNoCC {

        buildInputs = with pkgs; [
          go
          gopls
        ];

      };
    }
  );
}
