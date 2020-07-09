PROJECT_NAME := "go-tezos"
VERSION := "v2.10.0-alpha"
PKG := "github.com/goat-systems/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: dep clean test coverage coverhtml lint

checks: fmt lint staticcheck race test-integration ## Runs all quality checks

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

staticcheck: ## Static check the files
	@staticcheck ${PKG_LIST}

fmt: ## Static check the files
	@go fmt ${PKG_LIST}

test: ## Run unittests
	@go test -v ${PKG_LIST}

test-integration: ## Run unit tests and integration tests
	@go test -v --tags=integration ${PKG_LIST}

race: ## Run data race detector
	@go test -race -v ${PKG_LIST}

dep: ## Get the dependencies
	@go get -u golang.org/x/lint/golint
	@go get -u honnef.co/go/tools/...

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
