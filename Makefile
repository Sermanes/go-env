## ─── Project settings ─────────────────────────────────────────────────────────
MODULE = github.com/Sermanes/go-env

## ─── Go tooling ───────────────────────────────────────────────────────────────
GOCMD        = go
GOTEST       = $(GOCMD) test
GOVET        = $(GOCMD) vet
GOFMT        = $(GOCMD) fmt
GOMOD        = $(GOCMD) mod

COVERAGE_OUT  = coverage.out
COVERAGE_HTML = coverage.html

.PHONY: all clean \
        test test-race test-coverage bench \
        fmt vet lint lint-revive \
        tidy ci before-push \
        install-tools install-hooks \
        help

## ─── Default ──────────────────────────────────────────────────────────────────
all: ci

clean: ## Remove build artefacts
	rm -rf $(COVERAGE_OUT) $(COVERAGE_HTML)

## ─── Test ─────────────────────────────────────────────────────────────────────
test: ## Run unit tests (fast, no cache)
	go clean -testcache
	$(GOTEST) -v -failfast -short ./...

test-race: ## Run tests with the race detector
	go clean -testcache
	$(GOTEST) -race -v -failfast -short ./...

test-coverage: ## Generate HTML coverage report
	$(GOTEST) -coverprofile=$(COVERAGE_OUT) -covermode=atomic ./...
	$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report → $(COVERAGE_HTML)"

bench: ## Run benchmarks
	$(GOTEST) -bench=. -benchmem -run=^$$ ./...

## ─── Code quality ─────────────────────────────────────────────────────────────
fmt: ## Apply gofmt to all packages
	$(GOFMT) ./...

vet: ## Run go vet
	$(GOVET) ./...

lint: ## Run golangci-lint (full suite)
	golangci-lint run --config .golangci.yml ./...

lint-revive: ## Run revive linter only
	revive -formatter stylish -config linting.toml ./...

tidy: ## Tidy and verify go.mod / go.sum
	$(GOMOD) tidy
	$(GOMOD) verify

## ─── CI pipeline ──────────────────────────────────────────────────────────────
ci: tidy fmt vet test lint ## Full quality gate (used in CI)

before-push: fmt vet test test-race lint ## Gate to run before every git push

## ─── Tooling ──────────────────────────────────────────────────────────────────
install-tools: ## Install all required developer tools
	$(GOCMD) install github.com/mgechev/revive@latest
	$(GOCMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest

install-hooks: ## Copy git hooks from scripts/hooks/ into .git/hooks/
	@chmod +x scripts/hooks/pre-commit scripts/hooks/pre-push
	@cp scripts/hooks/pre-commit .git/hooks/pre-commit
	@cp scripts/hooks/pre-push   .git/hooks/pre-push
	@echo "Git hooks installed."

## ─── Help ─────────────────────────────────────────────────────────────────────
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
