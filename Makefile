.PHONY: help build install clean lint test fmt vet doctor

BINARY_NAME=envoy-ai-installer
VERSION?=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags="-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)"

help:
	@echo "Envoy AI Unified Installer - Build Commands"
	@echo ""
	@echo "Available targets:"
	@echo "  make build        - Build CLI binary"
	@echo "  make install      - Build and install binary"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Run golangci-lint"
	@echo "  make vet          - Run go vet"
	@echo "  make test         - Run tests"
	@echo "  make doctor       - Run system health check"
	@echo ""
	@echo "Development targets:"
	@echo "  make dev          - Build with debug info"
	@echo "  make release      - Build release binary"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION=$(VERSION)"
	@echo "  GIT_COMMIT=$(GIT_COMMIT)"
	@echo "  BUILD_TIME=$(BUILD_TIME)"

build:
	@echo "Building $(BINARY_NAME)..."
	@cd cli && go build $(LDFLAGS) -o ../$(BINARY_NAME)
	@echo "✓ Binary created: ./$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo install -Dm755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installed to: /usr/local/bin/$(BINARY_NAME)"

dev:
	@echo "Building debug binary..."
	@cd cli && go build -race -o ../$(BINARY_NAME)
	@echo "✓ Debug binary created: ./$(BINARY_NAME)"

release:
	@echo "Building release binary..."
	@cd cli && go build -ldflags "-s -w $(LDFLAGS)" -o ../$(BINARY_NAME)
	@echo "✓ Release binary created: ./$(BINARY_NAME)"
	@du -h ./$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@cd cli && go clean
	@echo "✓ Clean complete"

fmt:
	@echo "Formatting Go code..."
	@cd cli && go fmt ./...
	@echo "✓ Format complete"

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@cd cli && golangci-lint run ./...
	@echo "✓ Lint complete"

vet:
	@echo "Running go vet..."
	@cd cli && go vet ./...
	@echo "✓ Vet complete"

test:
	@echo "Running tests..."
	@cd cli && go test -v -race -coverprofile=coverage.out ./...
	@echo "✓ Tests complete"

doctor: build
	@echo "Running health check..."
	@./$(BINARY_NAME) doctor

version: build
	@./$(BINARY_NAME) version

.PHONY: all
all: clean fmt vet lint test build
	@echo "✓ All checks passed!"
