GOLANGCI_LINT_VERSION := v1.48.0
export GO111MODULE=on

## Lint
lint:
	@(if ! type golangci-lint >/dev/null 2>&1; then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION) ;fi)
	$$(go env GOPATH)/bin/golangci-lint run ./...


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

.PHONY: test lint fmt help
