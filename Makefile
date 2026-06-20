.DEFAULT_GOAL := help
PROJECT_NAME := svc

.PHONY: help
help: ## Display this help.
	@echo "------------------------------------------------------------------------"
	@echo "${PROJECT_NAME}"
	@echo "------------------------------------------------------------------------"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Code Quality

.PHONY: tidy
tidy: ## Remove unused dependencies and add missing ones.
	go mod tidy

.PHONY: fmt
fmt: ## Format all Go source files.
	go fmt ./...

.PHONY: vet
vet: ## Report likely mistakes in Go source code.
	go vet ./...

.PHONY: lint
lint: ## Run staticcheck linter against all source files.
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...

.PHONY: vulncheck
vulncheck: ## Scan for known vulnerabilities in dependencies.
	go run golang.org/x/vuln/cmd/govulncheck ./...

##@ Test

.PHONY: test
test: ## Run tests with race detector, randomized order, no caching.
	go test -race -shuffle=on -v -cover -count=1 ./...

##@ Examples

.PHONY: run-minimal
run-minimal: ## Run the minimal example.
	go run ./examples/minimal

.PHONY: run-liveliness
run-liveliness: ## Run the liveliness example.
	go run ./examples/liveliness
