let 
    pkgs = import <nixpkgs> {};
in pkgs.mkShell {
    buildInputs = with pkgs; [
        dbus
        openssl
    ];
    nativeBuildInputs = with pkgs; [
        pkg-config
        atk
        pango
        gdk-pixbuf
        gtk3
        cairo
        webkitgtk_4_1
    ];
    dbus = pkgs.dbus;
}