{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  packages = with pkgs; [
    wails
    bun
    go
  ];
  buildInputs = with pkgs; [
    webkitgtk_4_1
  ];
  nativeBuildInputs = with pkgs; [
    pkg-config
  ];
}
