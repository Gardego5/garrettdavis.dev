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
    terranix = {
      url = "github:terranix/terranix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { nixpkgs, naersk, fenix, terranix, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        app_name = "garrettdavis_dev";

        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
        lib = pkgs.lib;

        getDir = dir:
          builtins.mapAttrs (file: type:
            if type == "directory" then getDir "${dir}/${file}" else type)
          (builtins.readDir dir);
        getFiles = dir:
          lib.collect builtins.isString (lib.mapAttrsRecursive
            (path: type: builtins.concatStringsSep "/" path) (getDir dir));

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

        naerskBuildPackage = target: args:
          let
            crossPkgs = import nixpkgs {
              inherit system;
              crossSystem = { config = target; };
            };
            cc =
              "${crossPkgs.stdenv.cc}/bin/${crossPkgs.stdenv.cc.targetPrefix}cc";
            rustTargetPlatform =
              crossPkgs.rust.toRustTarget crossPkgs.stdenv.targetPlatform;
            rustTargetPlatformUpper = lib.toUpper
              (builtins.replaceStrings [ "-" ] [ "_" ] rustTargetPlatform);
          in naersk'.buildPackage (args // {
            "CARGO_BUILD_TARGET" = target;
            "CC_${rustTargetPlatform}_LINKER" = "${cc}";
            "CARGO_TARGET_${rustTargetPlatformUpper}_LINKER" = "${cc}";
            "RUSTFLAGS" = "-Zlocation-detail=none";

            depsBuildBuild = [ crossPkgs.stdenv.cc ];
            src = ./.;
            strictDeps = true;
          });

        css = pkgs.stdenv.mkDerivation {
          name = "output.css";
          buildCommand = let
            config = ./tailwind.config.js;
            input = ./src/input.css;
          in ''
            cp -r ${./src} ./src
            ${pkgs.tailwindcss}/bin/tailwindcss -c ${config} -i ${input} -o $out --minify
          '';
        };

        lambdaBinNames = let
          isRustFile = name: dirEntryType:
            dirEntryType == "regular" && lib.hasSuffix ".rs" name;
          toRustBinaryName = name: lib.removeSuffix ".rs" name;
        in map toRustBinaryName (builtins.attrNames
          (lib.attrsets.filterAttrs isRustFile (builtins.readDir ./src/bin)));
        lambdaBinaries = naerskBuildPackage "aarch64-unknown-linux-musl" { };
        lambdaPackages = lib.listToAttrs (map (name: {
          inherit name;
          value = pkgs.stdenv.mkDerivation {
            inherit name;
            buildCommand = ''
              mkdir $out
              cp ${lambdaBinaries}/bin/${name} $out/bootstrap
            '';
          };
        }) lambdaBinNames);

        lambdas = {
          hello_world = { source_dir = toString lambdaPackages.hello_world; };
          resume = { source_dir = toString lambdaPackages.resume; };
        };

        endpoints = [
          {
            lambda = "hello_world";
            method = "GET";
            path = "/hello-world";
          }
          {
            lambda = "resume";
            method = "GET";
            path = "/resume";
          }
        ];

        extToMime = {
          "js" = "text/javascript";
          "css" = "text/css";
        };

        static = {
          "style.css".content_type = extToMime.css;
          "style.css".source = toString css;
          "3p/js/htmx.min.js".content_type = extToMime.js;
          "3p/js/htmx.min.js".source = builtins.fetchurl {
            url =
              "https://github.com/bigskysoftware/htmx/releases/download/v1.9.12/htmx.min.js";
            sha256 =
              "sha256:0lm4lbsgjmgcmi6w54f7qjcs1hwmw68ljqfv22ar87l8wynig4s4";
          };
          "3p/js/alpinejs.min.js".content_type = extToMime.js;
          "3p/js/alpinejs.min.js".source = builtins.fetchurl {
            url =
              "https://cdn.jsdelivr.net/npm/alpinejs@3.14.0/dist/cdn.min.js";
            sha256 =
              "sha256:1llddh6qyip60nvyk0yzg2sdz6ydxlgfz23sglaxmyilcf88r61x";
          };
          "3p/js/alpinejs-morph.min.js".content_type = extToMime.js;
          "3p/js/alpinejs-morph.min.js".source = builtins.fetchurl {
            url =
              "https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js";
            sha256 =
              "sha256:08vm298my7c9ssbp1bxy66pkdl8dcd5a8nk9khvwxh2b49ykps4v";
          };
          "3p/js/iconify-icon.min.js".content_type = extToMime.js;
          "3p/js/iconify-icon.min.js".source = builtins.fetchurl {
            url =
              "https://code.iconify.design/iconify-icon/2.1.0/iconify-icon.min.js";
            sha256 =
              "sha256:00hgal6fhdwzk3njx1pyqyrwydny42hm82zbjzmzvjmhin1r93bm";
          };
        };

        #staticFiles = lib.listToAttrs (map (name:
        #  let
        #    ext = lib.lists.findFirst (ext: lib.hasSuffix ext name) null
        #      (builtins.attrNames extToMime);
        #    source = toString (pkgs.stdenv.mkDerivation {
        #      inherit name;
        #      buildCommand = "cp ${./static}/${name} $out";
        #    });
        #  in {
        #    inherit name;
        #    value = if ext == null then {
        #      inherit source;
        #    } else {
        #      inherit source;
        #      content_type = extToMime.${ext};
        #    };
        #  }) (getFiles ./static));

        infrastructure = terranix.lib.terranixConfiguration {
          inherit system;
          modules = [
            ./infra/api.nix
            ./infra/cdn.nix
            ./infra/dns.nix
            ./infra/lambda.nix
            ./infra/ses.nix
            ./infra/static.nix
            {
              provider.aws = { region = "us-west-2"; };
              terraform.backend.s3 = {
                bucket = "tf-state-20230722071359242500000001";
                key = "state/garrettdavis_dev_2";
                region = "us-west-2";
                encrypt = true;
                kms_key_id = "alias/terraform-state";
                dynamodb_table = "tf-state-20230722071359242500000001";
              };
              inherit app_name endpoints lambdas static;
              zone_id = "Z08935742SJWOUKHZGOC5";
              certificate_arn =
                "arn:aws:acm:us-east-1:256777061266:certificate/14831d98-9205-4a73-863b-35fafa5b2373";
            }
          ];
        };

      in {
        packages = lambdaPackages // {
          inherit css infrastructure;
          default = infrastructure;
        };

        apps = {
          terraform = {
            type = "app";
            program = toString (pkgs.writers.writeBash "terraform" ''
              if [[ -e config.tf.json ]]; then rm -f config.tf.json; fi;
              cp ${infrastructure} config.tf.json \
              && ${pkgs.terraform}/bin/terraform init \
              && ${pkgs.terraform}/bin/terraform "$@"
            '');
          };

          local = {
            type = "app";
            program = toString (pkgs.writers.writeBash "local" ''
              if [[ -e config.tf.json ]]; then rm -f config.tf.json; fi;
              cp ${infrastructure} config.tf.json \
              && ${pkgs.aws-sam-cli}/bin/sam build --hook-name terraform \
              && ${pkgs.aws-sam-cli}/bin/sam local start-api
            '');
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [
            (toolchain)
            fenix.packages.${system}.rust-analyzer

            # pkgs.cargo-audit
            # pkgs.cargo-bloat
            pkgs.cargo-outdated
            pkgs.cargo-release
            pkgs.cargo-watch
            pkgs.rustfmt

            pkgs.terranix
            pkgs.terraform
            pkgs.aws-sam-cli

            pkgs.tailwindcss
          ];
          RUST_BACKTRACE = 1;
        };
      });
}
