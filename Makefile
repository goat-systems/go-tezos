library = GoTezos
Go_Tezos_VERSION = v2.0.0
object = $(library)
package = github.com/DefinitelyNotAGoat/go-tezos

GO ?= GO111MODULE=on go
GOTEST_FLAGS ?=

GO_SOURCES := $(shell find $(PWD) -path $(PWD)/vendor -prune -o -path ./test -o -name '*.go' \! -name "*_test.go" -print)
GO_TEST_SOURCES := $(shell find $(PWD) -path $(PWD)/vendor -prune -o -name '*_test.go' -print)

test: $(GO_SOURCES) $(GO_TEST_SOURCES)
	$(GO) test $(GOTEST_FLAGS) -p=1 -cover ./...

fmt:
	gofmt -l -w -e $(GO_SOURCES) $(GO_TEST_SOURCES)

vet: 
	$(GO) vet -mod=vendor ./...

staticcheck:
	staticcheck ./...

checks: vet staticcheck