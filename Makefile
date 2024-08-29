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

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o bin/gerry .

.PHONY: debug-build
debug-build:
	go build -ldflags "$(LDFLAGS)" -gcflags=all="-N -l" -o bin/gerry .

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

.PHONY: docker
docker: build
	@docker build -t ghcr.io/distrobyte/gerry:$(GIT_COMMIT) .

.PHONY: test
test:
	go test -v -bench=. ./...

.PHONY: watch
watch:
	air start
