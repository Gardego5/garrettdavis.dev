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
        inherit (pkgs) lib;

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
          blog = { source_dir = toString lambdaPackages.blog; };
          contact = { source_dir = toString lambdaPackages.contact; };
          notes = { source_dir = toString lambdaPackages.notes; };
          resume = { source_dir = toString lambdaPackages.resume; };
          index = { source_dir = toString lambdaPackages.index; };
        };

        endpoints = {
          "GET /blog" = "blog";
          "GET /contact" = "contact";
          "POST /contact" = "contact";
          "GET /notes" = "notes";
          "GET /notes/{id+}" = "notes";
          "GET /resume" = "resume";
          "$default" = "index";
        };

        inherit (import ./infra/static_files.nix { inherit pkgs lib; })
          css font static staticFilesDirectory;

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
          inherit css font infrastructure staticFilesDirectory;
          default = infrastructure;
        };

        apps = rec {
          terraform = {
            type = "app";
            program = toString (pkgs.writers.writeBash "terraform" ''
              if [[ -e config.tf.json ]]; then rm -f config.tf.json; fi;
              cp ${infrastructure} config.tf.json \
              && ${pkgs.terraform}/bin/terraform init \
              && ${pkgs.terraform}/bin/terraform "$@"
            '');
          };

          local-refresh = {
            type = "app";
            program = toString (pkgs.writers.writeBash "local-refresh" ''
              if [[ -e config.tf.json ]]; then rm -f config.tf.json; fi;
              cp ${infrastructure} config.tf.json
              ${pkgs.aws-sam-cli}/bin/sam build --hook-name terraform
            '');
          };

          local = let
            caddyfile = pkgs.writers.writeText "Caddyfile" ''
              http://localhost:8000

              reverse_proxy localhost:3000

              file_server /static/* {
                root ${staticFilesDirectory}
              }
            '';
          in {
            type = "app";
            program = toString (pkgs.writers.writeBash "devserver" ''
              set -o pipefail
              trap 'kill %1; kill %2' SIGINT
              ${pkgs.aws-sam-cli}/bin/sam local start-api --warm-containers lazy & \
              ${pkgs.caddy}/bin/caddy run --config ${caddyfile} --adapter caddyfile
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
