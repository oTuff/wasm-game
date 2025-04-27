{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    nixgl.url = "github:nix-community/nixGL";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      nixpkgs,
      nixgl,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ nixgl.overlay ]; # Add nixGL overlay
        };
      in
      {
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            pkgs.nixgl.nixGLIntel
            go
            libGL
            mesa
            xorg.libXrandr
            xorg.libXcursor
            xorg.libXinerama
            xorg.libXi
            xorg.libXxf86vm
            alsa-lib
            pkg-config
            tinygo
            twiggy
            binaryen
          ];

          shellHook = with pkgs; ''
            export LIBGL_ALWAYS_SOFTWARE=1
            export LD_LIBRARY_PATH=${lib.getLib libGL}/lib:$LD_LIBRARY_PATH
          '';
        };
      }
    );
}
