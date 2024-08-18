.PHONY: build docker run test watch config
build:
	CGO_ENABLED=0 go build -o build/gerry ./cmd/gerry

docker:
	@docker build -t ghcr.io/distrobyte/gerry:latest .

run:
	go run cmd/gerry/main.go start

test:
	go test -v ./...

watch:
	air start

config: build
	./build/gerry confgen config.yaml
