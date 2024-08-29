all: build

GIT_COMMIT:=$(shell git rev-list -1 HEAD)
GIT_LAST_TAG:=$(shell git describe --abbrev=0 --tags)
GIT_EXACT_TAG:=$(shell git name-rev --name-only --tags HEAD)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

COMMANDS_PATH:=github.com/distrobyte/gerry/internal/config
LDFLAGS:=-X ${COMMANDS_PATH}.GitCommit=${GIT_COMMIT} \
	-X ${COMMANDS_PATH}.GitLastTag=${GIT_LAST_TAG} \
	-X ${COMMANDS_PATH}.GitExactTag=${GIT_EXACT_TAG} \
	-X ${COMMANDS_PATH}.BuildDate=${BUILD_DATE}

.PHONY: build docker run test watch config
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" .

.PHONY: debug-build
debug-build:
	go build -ldflags "$(LDFLAGS)" -gcflags=all="-N -l" .

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

docker: build
	@docker build -t ghcr.io/distrobyte/gerry:$(GIT_COMMIT) .

test:
	go test -v -bench=. ./...

watch:
	air start
