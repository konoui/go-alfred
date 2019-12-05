
export GO111MODULE=on


## Setup
setup:
	#installing golint
	@(if ! type golint >/dev/null 2>&1; then go get -u golang.org/x/lint/golint ;fi)
	#installing golangci-lint
	@(if ! type golangci-lint >/dev/null 2>&1; then go get -u github.com/golangci/golangci-lint/cmd/golangci-lint ;fi)
	#installing goimports
	@(if ! type goimports >/dev/null 2>&1; then go get -u golang.org/x/tools/cmd/goimports ;fi)
	#installing ghr
	@(if ! type ghr >/dev/null 2>&1; then go get -u github.com/tcnksm/ghr ;fi)
	#installing make2help
	@(if ! type make2help >/dev/null 2>&1; then go get -u github.com/Songmu/make2help/cmd/make2help ;fi)

## Format source codes
fmt: setup
	goimports -w $$(go list -f {{.Dir}} ./... | grep -v /vendor/)

## Lint
lint: setup
	golangci-lint run ./...


## Run tests for my project
test: setup
	go test -v ./...


## Initialize directory
init:
	@(if [ ! -e ${SRC_DIR} ]; then mkdir ${SRC_DIR}; fi)
	@(if [ ! -e ${BIN_DIR} ]; then mkdir ${BIN_DIR}; fi)
	@(if [ ! -e go.mod ]; then go mod init; fi)


## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: build setup test lint fmt linux init clean help
