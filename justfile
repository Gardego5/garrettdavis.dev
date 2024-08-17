run: (build-static "css") (build-static "font") (build-static "staticDir")
	cargo -C app -Z unstable-options run

build-static derivation:
	nix build .#{{derivation}} --out-link app/target/static_dir
