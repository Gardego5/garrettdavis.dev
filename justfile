run: (build-static "css") (build-static "font") (build-static "staticDir")
	cargo -C app -Z unstable-options run

build-static derivation:
	nix build .#{{derivation}} --out-link app/target/static_dir

docker-load:
	docker load < $(nix build .#dockerImage --print-out-paths --no-link)

docker-push: docker-load
	docker image push registry.fly.io/garrettdavis-dev:latest

docker-run: docker-load
	docker run -p 3000:3000 garrettdavis.dev:latest
