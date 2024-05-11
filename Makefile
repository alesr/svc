PHONY: all
all: lint tools test

PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@1.58.1

.PHONY: test
test:
	go test -race  ./...
