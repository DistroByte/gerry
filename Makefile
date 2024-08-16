build:
	@GOOS=linux CGO_ENABLED=0 go build -o build/gerry ./cmd/gerry
	@docker build -t ghcr.io/DistroByte/gerry:latest .

run:
	@go run cmd/gerry/main.go

test:
	@go test -v ./...

watch:
	@air