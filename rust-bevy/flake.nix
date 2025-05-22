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
          overlays = [ nixgl.overlay ];
        };
      in
      {
        devShell = pkgs.mkShell rec {
          nativeBuildInputs = with pkgs; [
            busybox
            gnumake
            bmake
            pdpmake

            rustc
            cargo
            rust-analyzer
            rustfmt
            wasm-pack
            pkg-config
            pkgs.nixgl.nixVulkanIntel # nixGLIntel
            lld
            wasm-bindgen-cli
            binaryen # for wasm-opt
            mold
            cargo-zigbuild
            twiggy
          ];

          buildInputs = with pkgs; [
            udev
            alsa-lib-with-plugins
            vulkan-loader
            xorg.libX11
            xorg.libXcursor
            xorg.libXi
            xorg.libXrandr
            libxkbcommon
          ];
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath buildInputs;
        };
      }
    );
}
