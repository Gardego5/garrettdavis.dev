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
            rustTargetPlatformUpper = pkgs.lib.toUpper
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

        lambdaBinaries = [ "hello_world" ];

        lambdaPackages = pkgs.lib.listToAttrs (builtins.map (name: {
          inherit name;
          value = naerskBuildPackage "aarch64-unknown-linux-musl" {
            pname = name;
            postInstall = "mv $out/bin/${name} $out/bin/bootstrap ";
          };
        }) lambdaBinaries);

        infrastructure = terranix.lib.terranixConfiguration {
          inherit system;
          modules = [
            ./infra/lambda.nix
            {
              config = {
                provider.aws = { region = "us-west-2"; };

                app_name = "garrettdavis_dev";

                lambdas = pkgs.lib.attrsets.concatMapAttrs (name: package: {
                  ${name} = { source_dir = "${package}/bin"; };
                }) lambdaPackages;
              };
            }
          ];
        };

      in rec {
        packages = lambdaPackages // { default = infrastructure; };

        apps = {
          terraform = {
            type = "app";
            program = builtins.toString (pkgs.writers.writeBash "apply" ''
              if [[ -e config.tf.json ]]; then rm -f config.tf.json; fi;
              cp ${packages.default} config.tf.json \
              && ${pkgs.terraform}/bin/terraform init \
              && ${pkgs.terraform}/bin/terraform "$@"
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

            pkgs.terranix
            pkgs.terraform
          ];
          RUST_BACKTRACE = 1;
        };
      });
}
