{
  pkgs ? import <nixpkgs> { },
}:
pkgs.mkShell {
  packages = with pkgs; [
    wails
    bun
    go
    nodePackages.vercel
  ];
  buildInputs = with pkgs; [
    webkitgtk_4_1
  ];
  nativeBuildInputs = with pkgs; [
    pkg-config
  ];
}
