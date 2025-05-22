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
      {
        devShell =
          with pkgs;
          mkShell {
            nativeBuildInputs = [
              (python3.withPackages (
                p: with p; [
                  pandas
                  matplotlib
                  seaborn
                  pyqt5
                ]
              ))
            ];
            env.LD_LIBRARY_PATH = lib.makeLibraryPath [
              # stdenv.cc.cc.lib
              # libz
              qt5.qtbase
              # libGL
            ];
            QT_QPA_PLATFORM_PLUGIN_PATH = "${qt5.qtbase.bin}/lib/qt-${qt5.qtbase.version}/plugins";
            # QT_PLUGIN_PATH = "${qt5.qtbase.bin}/lib/qt-${qt5.qtbase.version}/plugins";
          };
      }
    );

}
