{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    naersk = {
      url = "github:nix-community/naersk";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { nixpkgs, naersk, fenix, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "garrettdavis.dev";
        src = ./app;

        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        toolchain = with fenix.packages.${system};
          combine [
            minimal.rustc
            minimal.cargo
            targets.aarch64-unknown-linux-gnu.latest.rust-std
            targets.aarch64-unknown-linux-musl.latest.rust-std
          ];

        naersk' = naersk.lib.${system}.override {
          cargo = toolchain;
          rustc = toolchain;
        };

        app = naersk'.buildPackage {
          inherit src;
          strictDeps = true;
        };

        css = pkgs.stdenv.mkDerivation {
          inherit src;
          name = "style.css";
          phases = [ "buildCommand" ];
          buildCommand = ''
            cd $src
            ${pkgs.tailwindcss}/bin/tailwindcss \
              -c $src/tailwind.config.js \
              -i $src/input.css -o $out --minify
          '';
        };

        font = let
          webfontIosevka = pkgs.iosevka.overrideAttrs {
            buildPhase = ''
              export HOME=$TMPDIR
              runHook preBuild
              npm run build --no-update-notifier --targets webfont::$pname -- --jCmd=$NIX_BUILD_CORES --verbose=9
              runHook postBuild
            '';
            installPhase = ''
              runHook preInstall
              fontdir="$out/share/fonts/truetype"
              mkdir -p "$fontdir"
              cp -r "dist/$pname"/* "$fontdir"
              runHook postInstall
            '';
          };
        in webfontIosevka.override {
          set = "GarrettDavisDev";
          privateBuildPlan = ''
            [buildPlans.IosevkaGarrettDavisDev]
            family = "Iosevka GarrettDavisDev"
            spacing = "normal"
            serifs = "sans"
            noCvSs = true
            exportGlyphNames = false

            [buildPlans.IosevkaGarrettDavisDev.variants]
            inherits = "ss20"

            [buildPlans.IosevkaGarrettDavisDev.weights.Regular]
            shape = 400
            menu = 400
            css = 400

            [buildPlans.IosevkaGarrettDavisDev.weights.Heavy]
            shape = 900
            menu = 900
            css = 900

            [buildPlans.IosevkaGarrettDavisDev.widths.Normal]
            shape = 500
            menu = 5
            css = "normal"

            [buildPlans.IosevkaGarrettDavisDev.slopes.Upright]
            angle = 0
            shape = "upright"
            menu = "upright"
            css = "normal"

            [buildPlans.IosevkaGarrettDavisDev.slopes.Italic]
            angle = 9.4
            shape = "italic"
            menu = "italic"
            css = "italic"
          '';
        };

        staticDir = pkgs.stdenv.mkDerivation {
          name = "bundle";
          phases = [ "installPhase" ];
          installPhase = ''
            # copy in static folder
            mkdir -p $out
            cp -r ${./static}/* $out
            chmod +w -R $out # permit adding new files

            # copy additional built files
            cp ${css} $out/style.css
            mkdir -p $out/3p/font
            cp -r ${font}/share/fonts/truetype/* $out/3p/font
          '';
        };

        dockerImage = pkgs.dockerTools.buildImage {
          inherit name;
          tag = "latest";
          created = "now";
          copyToRoot = [ pkgs.curl ];
          config = {
            Expose = 3000;
            Cmd = [ "${app}/bin/garrettdavis-dev" ];
            Env = [ "STATIC_DIR=${staticDir}" ];
          };
        };

      in {
        packages = { inherit app css dockerImage font staticDir; };

        devShells.default = pkgs.mkShell {
          packages = [
            toolchain
            fenix.packages.${system}.rust-analyzer

            pkgs.cargo-outdated
            pkgs.cargo-release
            pkgs.cargo-watch
            pkgs.rustfmt

            pkgs.terraform
            pkgs.tailwindcss
          ];
          RUST_BACKTRACE = 1;
          STATIC_DIR = toString staticDir;
        };
      });
}
