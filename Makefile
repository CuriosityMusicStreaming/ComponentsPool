all: check test

.PHONY: modules
modules:
	go mod tidy

.PHONY: test
test: modules
	go test ./...

.PHONY: check
check:
	golangci-lint run