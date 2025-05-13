# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOVET=$(GOCMD) vet
BINARY_NAME=owlify
COVERAGE_DIR=coverage
COVERAGE_FILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

TAG = ${NEW_RELEASE_TAG}
TARGET = novelist

check-tag:
ifndef NEW_RELEASE_TAG
	$(error Please set the NEW_RELEASE_TAG env variable)
	exit 1
endif

check-token:
ifndef GITHUB_TOKEN
	$(error Please set the GITHUB_TOKEN env variable)
	exit 1
endif

# Default target
.PHONY: all
all: test build

# Build the application
.PHONY: build
build: check-tag
	@echo "Building..."
	goreleaser build --clean

.PHONY: publish
publish: tidy tag check-token
	go install github.com/goreleaser/goreleaser@latest
	goreleaser release --clean

tag: check-tag
	git tag -a "$(TAG)" -m "$(TAG)"
	git push origin $(TAG)

tidy:
	go mod tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

# Run tests with race detection
.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race ./...

# Run linter (golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(COVERAGE_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v -t -d ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Run all checks (test, lint, vet)
.PHONY: check
check: test lint vet

# Help command
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all          - Run tests and build the application"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make test-race    - Run tests with race detection"
	@echo "  make lint         - Run linter"
	@echo "  make vet          - Run go vet"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Install dependencies"
	@echo "  make fmt          - Format code"
	@echo "  make check        - Run all checks (test, lint, vet)"
	@echo "  make help         - Show this help message" 