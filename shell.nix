with import <nixpkgs> { };

stdenv.mkDerivation {
  name = "go";
  buildInputs = [
    delve
    gcc
    go
    libcap
    (libreoffice-still-unwrapped.overrideAttrs (oA: {
      patches = oA.patches or [] ++ [
        ./libreoffice.patch
      ];
    }))
  ];
  shellHook = ''
    export GOPATH=$PWD/gopath
  '';
}
