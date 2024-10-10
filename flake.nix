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

        app = crossPkgs.buildGo123Module rec {
          pname = "github.com/Gardego5/garrettdavis.dev";
          version = "v0.0.1";
          nativeBuildInputs = [ pkgs.tailwindcss pkgs.rsync pkgs.nix ];
          preBuild = ''
            # generate static files
            go generate -tags ${builtins.concatStringsSep "," tags} ./...
            # copy static files that are generated with nix
            mkdir -p ./build/share
            rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r "${font}/share/fonts" ./build/share
          '';
          inherit src;
          vendorHash = "sha256-fuyTpQ9n9Idr5vk8dakEXmgOrzO+cnq3RiLVftnHNwQ=";
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
        packages = { inherit app dockerImage font; };

        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go_1_23
              gopls
              just
              tailwindcss
              flyctl
              turso-cli
            ];
          };
          cicd = pkgs.mkShell { packages = with pkgs; [ docker flyctl just ]; };
        };
      });
}
