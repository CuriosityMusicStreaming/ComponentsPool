export APP_CMD_NAME = componentspool
export APP_PROTO_FILES = \
	api/contentservice/contentservice.proto
export DOCKER_IMAGE_NAME = vadimmakerov/$(APP_CMD_NAME):master

all: build check test

.PHONY: build
build: modules
	bin/go-build.sh "cmd" "bin/$(APP_CMD_NAME)" $(APP_CMD_NAME)

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	golangci-lint run