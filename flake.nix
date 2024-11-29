{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "garrettdavis-dev";
        src = ./.;

        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        crossPkgs = import nixpkgs {
          localSystem = system;
          crossSystem = "x86_64-linux";
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

        msgp-go = pkgs.buildGoModule rec {
          name = "msgp";
          version = "v1.2.2";
          src = pkgs.fetchFromGitHub {
            owner = "tinylib";
            repo = "msgp";
            rev = version;
            sha256 = "sha256-5AXy3/FFNrgqWRf/4t2vFpJeBcGaE6LZCFaKNvLhmH0=";
          };
          doCheck = false;
          vendorHash = "sha256-aWS13hx7ZVJGArBS381GJTvhd8Kl6WtbMIGEGV/iChY=";
        };

        tailwind = let
          throwSystem =
            throw "tailwindcss has not been packaged for ${system} yet.";

          plat = {
            aarch64-darwin = "macos-arm64";
            aarch64-linux = "linux-arm64";
            armv7l-linux = "linux-armv7";
            x86_64-darwin = "macos-x64";
            x86_64-linux = "linux-x64";
          }.${system} or throwSystem;

          sha256 = {
            aarch64-darwin =
              "bce402ef6da7f3da611021a389bec0bf082c8c85a8bed284d8ccd86d9eafff8c";
            aarch64-linux = pkgs.lib.fakeSha256;
            armv7l-linux = pkgs.lib.fakeSha256;
            x86_64-darwin = pkgs.lib.fakeSha256;
            x86_64-linux = pkgs.lib.fakeSha256;
          }.${system} or throwSystem;

        in pkgs.tailwindcss.overrideAttrs (final: prev: rec {
          version = "4.0.0-beta.1";
          src = pkgs.fetchurl {
            url =
              "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-${plat}";
            inherit sha256;
          };
          installPhase = ''
            mkdir -p $out/bin
            cp ${src} $out/bin/tailwindcss
            chmod 755 $out/bin/tailwindcss
          '';
        });

        app = crossPkgs.buildGo123Module rec {
          pname = "github.com/Gardego5/garrettdavis.dev";
          version = "v0.0.1";
          nativeBuildInputs = [ pkgs.tailwindcss pkgs.rsync pkgs.nix msgp-go ];
          preBuild = ''
            # generate static files
            go generate -tags ${builtins.concatStringsSep "," tags} ./...
            # copy static files that are generated with nix
            mkdir -p ./build/share
            rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r "${font}/share/fonts" ./build/share
          '';
          inherit src;
          vendorHash = "sha256-qgOamviwFiwG+1GcEKKOBGNASZpTkj8GU5RmiIZOVQ0=";
          tags = [ "fonts" "static" ];
        };

        dockerImage = crossPkgs.dockerTools.buildImage {
          name = "registry.fly.io/${name}";
          tag = "latest";
          created = "now";
          copyToRoot = [ pkgs.curl pkgs.cacert ];
          config = {
            Expose = 3000;
            Cmd = [ "${app}/bin/garrettdavis.dev" ];
            Env = [ "FONTS_DIR=${font}/share/fonts" "PORT=3000" ];
          };
        };

      in {
        packages = { inherit app dockerImage font tailwind; };

        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              flyctl
              go_1_23
              gopls
              just
              redis
              tailwind
              turso-cli
              wire

              msgp-go
            ];
          };
          cicd = pkgs.mkShell { packages = with pkgs; [ docker flyctl just ]; };
        };
      });
}
