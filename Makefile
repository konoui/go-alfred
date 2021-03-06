GOLANGCI_LINT_VERSION := v1.30.0
export GO111MODULE=on

## Format source codes
fmt:
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint:
	@(if ! type golangci-lint >/dev/null 2>&1; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin ${GOLANGCI_LINT_VERSION} ;fi)
	golangci-lint run ./...


## Run tests for my project
test:
	. scripts/setup.sh; \
	go test -v ./...

bench:
	. scripts/setup.sh; \
	go test -benchmem -run="^$$" -bench "^(Benchmark.+)$$" -benchtime 1x -count 2

build-examples:
	examples/build-test.sh

generate:
	go generate ./...

cover:
	. scripts/setup.sh; \
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

## Show help
help:
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)
	@make2help $(MAKEFILE_LIST)

.PHONY: test lint fmt help
