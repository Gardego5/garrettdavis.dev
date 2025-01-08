go:
	go generate ./...

	rm -rf build
	rsync -q -av --no-o --no-g --chmod=Du=rwx,Dg=rx,Do=rx,Fu=rw,Fg=r,Fo=r "$(nix build .#build --no-link --print-out-paths)/share" .
	mv share build

	go run -ldflags "-X 'main.CacheID=$(git rev-parse --short HEAD)'" .

docker-load:
	docker load < $(nix build .#dockerImage --print-out-paths --no-link)

docker-push: docker-load
	docker image push registry.fly.io/garrettdavis-dev:latest

docker-run: docker-load
	docker run -p 3000:3000 garrettdavis.dev:latest
