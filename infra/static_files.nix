{ pkgs, lib, ... }:
let
  getDir = dir:
    builtins.mapAttrs
    (file: type: if type == "directory" then getDir "${dir}/${file}" else type)
    (builtins.readDir dir);

  getFiles = dir:
    lib.collect builtins.isString
    (lib.mapAttrsRecursive (path: type: builtins.concatStringsSep "/" path)
      (getDir dir));
in rec {
  extToMime = {
    "js" = "text/javascript";
    "css" = "text/css";
  };

  css = pkgs.stdenv.mkDerivation {
    name = "output.css";
    buildCommand = let
      config = ../tailwind.config.js;
      input = ../src/input.css;
    in ''
      cp -r ${../src} ./src
      ${pkgs.tailwindcss}/bin/tailwindcss -c ${config} -i ${input} -o $out --minify
    '';
  };

  static = {
    "style.css".content_type = extToMime.css;
    "style.css".source = toString css;
    "3p/js/htmx.min.js".content_type = extToMime.js;
    "3p/js/htmx.min.js".source = builtins.fetchurl {
      url =
        "https://github.com/bigskysoftware/htmx/releases/download/v1.9.12/htmx.min.js";
      sha256 = "sha256:0lm4lbsgjmgcmi6w54f7qjcs1hwmw68ljqfv22ar87l8wynig4s4";
    };
    "3p/js/alpinejs.min.js".content_type = extToMime.js;
    "3p/js/alpinejs.min.js".source = builtins.fetchurl {
      url = "https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/cdn.min.js";
      sha256 = "sha256:1llddh6qyip60nvyk0yzg2sdz6ydxlgfz23sglaxmyilcf88r61x";
    };
    "3p/js/alpinejs-morph.min.js".content_type = extToMime.js;
    "3p/js/alpinejs-morph.min.js".source = builtins.fetchurl {
      url =
        "https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js";
      sha256 = "sha256:08vm298my7c9ssbp1bxy66pkdl8dcd5a8nk9khvwxh2b49ykps4v";
    };
    "3p/js/iconify-icon.min.js".content_type = extToMime.js;
    "3p/js/iconify-icon.min.js".source = builtins.fetchurl {
      url =
        "https://code.iconify.design/iconify-icon/2.1.0/iconify-icon.min.js";
      sha256 = "sha256:00hgal6fhdwzk3njx1pyqyrwydny42hm82zbjzmzvjmhin1r93bm";
    };
  };

  staticFilesDirectory = pkgs.stdenv.mkDerivation {
    name = "staticFilesDirectory";
    src = ../.;
    installPhase = lib.strings.concatLines (lib.foldlAttrs (acc: name:
      { source, ... }:
      let file = "$out/static/${name}";
      in acc ++ [ "mkdir -p $(dirname ${file}) && cp ${source} ${file}" ]) [ ]
      static);
  };
}
