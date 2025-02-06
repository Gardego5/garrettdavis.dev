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
            sha256 = "sha256-h3l3PsfK/uWwjOLxhj4XukWVzcydLuv60TKTc15YA3A=";
          };
          doCheck = false;
          vendorHash = "sha256-aWS13hx7ZVJGArBS381GJTvhd8Kl6WtbMIGEGV/iChY=";
        };

        tailwind = let
          throwSystem = reason:
            throw
            "tailwindcss has not been packaged for ${system} yet (${reason}).";

          plat = {
            aarch64-darwin = "macos-arm64";
            aarch64-linux = "linux-arm64-musl";
            x86_64-darwin = "macos-x64";
            x86_64-linux = "linux-x64-musl";
          }.${system} or (throwSystem "missing platform");

          hash = {
            aarch64-darwin =
              "sha256-hH9+h6jtXS9uT5mujDRTtjRM2onG8ZQsexOlMaIoXv4=";
            aarch64-linux =
              "sha256-yDKmh7DZl62FCZrSzKe14tK4dg5Cqxijx/O/E6+NiLY=";
            x86_64-linux =
              "sha256-tSn4e+i+48/WikqK7VHBYPdbU1ypxrlMxTJi6ZzH1gw=";
            x86_64-darwin =
              "sha256-iMPHW3snWY9nWgRv6+0IS3Zh29LC0kYmzfwOcJM8xN0=";
          }.${system} or (throwSystem "missing hash");

        in pkgs.tailwindcss.overrideAttrs (final: prev: rec {
          version = "4.0.4";
          buildInputs = [ ];
          src = pkgs.fetchurl {
            url =
              "https://github.com/tailwindlabs/tailwindcss/releases/download/v${version}/tailwindcss-${plat}";
            inherit hash;
          };
          sourceRoot = ".";
          installPhase = ''
            runHook preInstall
            install -m755 -D $src $out/bin/tailwindcss
            runHook postInstall
          '';
        });

        css = pkgs.stdenv.mkDerivation {
          name = "css";
          nativeBuildInputs = [ tailwind ];
          inherit src;
          installPhase = ''
            mkdir -p $out/share
            (cd $src && tailwindcss -i input.css -o $out/share/css/style.css)
          '';
        };

        staticFiles = pkgs.stdenv.mkDerivation {
          name = "staticFiles";
          src = ./static;
          phases = [ "installPhase" ];
          installPhase = ''
            mkdir -p $out/share
            cp -r $src/* $out/share
          '';
        };

        rsyncDerivations = (name: drvs:
          pkgs.stdenv.mkDerivation {
            inherit name;
            nativeBuildInputs = [ pkgs.rsync ];
            phases = [ "installPhase" ];
            installPhase = builtins.concatStringsSep "\n" (builtins.map (drv:
              "rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r '${drv}/${name}' $out")
              drvs);
          });

        build = rsyncDerivations "share" [ font staticFiles css ];

        app = let
          module = rec {
            inherit src;
            pname = "github.com/Gardego5/garrettdavis.dev";
            version = "v0.0.1";
            nativeBuildInputs = [ pkgs.rsync msgp-go ];
            preBuild = ''
              # generate static files
              go generate -tags ${builtins.concatStringsSep "," tags} ./...

              # copy static files that are generated with nix
              rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r "${build}/share" .
              mv share build
            '';
            ldflags = [ ];
            vendorHash = "sha256-hQWoaVdZRqmRp2yveZCpUqiBqiXLZ7FMDDFX0FWzdNE=";
            tags = [ "fonts" "static" ];
          };
          cacheId = builtins.hashString "md5" (builtins.toJSON module);
        in crossPkgs.buildGo123Module (module // {
          ldflags = module.ldflags ++ [ "-X 'main.CacheID=${cacheId}'" ];
        });

        dockerImage = crossPkgs.dockerTools.buildImage {
          name = "registry.fly.io/${name}";
          tag = "latest";
          created = "now";
          copyToRoot = [ pkgs.curl pkgs.cacert ];
          config = {
            Expose = 3000;
            Cmd = [ "${app}/bin/garrettdavis.dev" ];
            Env = [ "PORT=3000" ];
          };
        };

      in {
        packages = { inherit app build css dockerImage font tailwind; };

        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              fblog
              flyctl
              go_1_23
              gopls
              just
              msgp-go
              tailwind
              redis
              turso-cli
              wire
            ];
          };
          cicd = pkgs.mkShell { packages = with pkgs; [ docker flyctl just ]; };
        };
      });
}
