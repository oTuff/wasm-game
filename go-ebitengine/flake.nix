{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        devShell =
          with pkgs;
          mkShell {
            nativeBuildInputs = [
              go
              # gopls
              gotools
              mesa.drivers
              mesa
              gcc
              libGL
              xorg.libX11
              xorg.libXcursor
              xorg.libXext
              xorg.libXi
              xorg.libXinerama
              xorg.libXrandr
              xorg.libXxf86vm
              # glfw
              # alsa-lib
              # pkg-config
            ];

            env.LD_LIBRARY_PATH = lib.makeLibraryPath [
              libGL
              stdenv.cc.cc.lib
            ];
          };
      }
    );
}
