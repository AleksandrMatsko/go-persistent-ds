GO_PATH := $(shell go env GOPATH)
GOLANGCI_LINT_VERSION := ""

.PHONY: default
default: build

.PHONY: build
build:
	go build .

.PHONY: test
test:
	go test -v -race ./...

.PHONY: install-lint
install-lint:
	# The recommended way to install golangci-lint into CI/CD
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GO_PATH}/bin ${GOLANGCI_LINT_VERSION}

.PHONY: lint
lint:
	golangci-lint run