.PHONY: build docker run test watch config
build:
	go build -o build/gerry ./cmd/gerry

docker: build
	@docker build -t ghcr.io/DistroByte/gerry:latest .

run:
	go run cmd/gerry/main.go start

test:
	go test -v ./...

watch:
	air start

config: build
	./build/gerry confgen config.yaml
	