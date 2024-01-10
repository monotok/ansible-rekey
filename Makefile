GO := go
ifdef GO_BIN
	GO = $(GO_BIN)
endif

GOLANGCI_LINT_VERSION := v1.18.0
BIN_DIR := $(GOPATH)/bin

all: test lint

tidy:
	$(GO) mod tidy -v

fmt:
	gofmt -s -w .

build:
	$(GO) build ./...

ci-build:
	$(GO) build -v -ldflags '-X $(VERSION_PACKAGE).GitHash=$(GIT_COMMIT) -X $(VERSION_PACKAGE).GitTag=$(GIT_TAG) -X $(VERSION_PACKAGE).GitBranch=$(GIT_BRANCH) -X $(VERSION_PACKAGE).BuildTime=$(BUILD_TIME) -X $(VERSION_PACKAGE).GitCommitMessage=$(GIT_COMMIT_MESSAGE)'

test: build
	$(GO) test -cover -race -v ./...

test-coverage:
	$(GO) test ./... -race -coverprofile=.testCoverage.txt && $(GO) tool cover -html=.testCoverage.txt

ci-test: ci-build
	$(GO) test -race $$($(GO) list ./...) -v -coverprofile .testCoverage.txt

lint: $(GOLANGCI_LINT)
	golangci-lint run --fast

$(GOLANGCI_LINT):
	GO111MODULE=on $(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)