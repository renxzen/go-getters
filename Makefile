GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOPATHBIN=$(shell go env GOPATH)/bin
BINARY_NAME=go-getters
MAIN_PATH=./cmd/go-getters
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html


help:
	@echo "This is a helper makefile for go-getters"
	@echo "Targets:"
	@echo "    build:          build the binary"
	@echo "    install:        install the binary"
	@echo "    test:           run all tests"
	@echo "    test/update:    update golden test files"
	@echo "    generate:       run go generate to update examples"
	@echo "    generate/local: run go generate with local build to update examples"
	@echo "    tidy:           tidy go mod"
	@echo "    lint:           lint the project"
	@echo "    clean:          clean build artifacts"

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v2.2.2

.PHONY: tools
tools: $(GOBIN)/golangci-lint

build:
	@mkdir -p bin
	go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_PATH)

install:
	go install $(MAIN_PATH)

uninstall:
	rm -rf $(GOPATHBIN)/$(BINARY_NAME)

lint: tools
	$(GOBIN)/golangci-lint run ./...

.PHONY: test
test:
	@echo "Running tests..."
	go test -race -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test/update:
	go test -v ./test -update

generate: build
	go generate ./...

generate/local: build
	@export PATH="$$(pwd)/bin:$$PATH" && go generate ./...

tidy:
	go mod tidy

clean:
	go clean
	@rm -rf bin/
