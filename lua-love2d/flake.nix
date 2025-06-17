{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            nodejs_24 # needed for love.js
            luajit
            love
          ];

          shellHook = with pkgs; ''
            export LD_LIBRARY_PATH=${
              lib.makeLibraryPath [
                mesa
                # egl-wayland
                # libvdpau-va-gl
                # libglvnd
                # egl-wayland
                # libGL
              ]
            }:$LD_LIBRARY_PATH
          '';
        };
      }
    );
}
