go:
	go generate ./...
	rm -rf ./build
	mkdir -p ./build/share
	rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r "$(nix build .#font --no-link --print-out-paths)/share/fonts" ./build/share
	CACHE_ID=$(git rev-parse --short HEAD) go run . | fblog -d --

docker-load:
	docker load < $(nix build .#dockerImage --print-out-paths --no-link)

docker-push: docker-load
	docker image push registry.fly.io/garrettdavis-dev:latest

docker-run: docker-load
	docker run -p 3000:3000 garrettdavis.dev:latest
