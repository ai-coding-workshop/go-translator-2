# Makefile for Translator Service

# Default target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  format     - Format code with goimports and gofmt"
	@echo "  lint       - Run golangci-lint for static analysis"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run all tests"
	@echo "  bench      - Run benchmarks"
	@echo "  install-tools - Install development tools"

# Format code with goimports and gofmt
.PHONY: format
format:
	@echo "Formatting code..."
	goimports -w .
	gofmt -s -w .

# Lint code with golangci-lint
.PHONY: lint
lint:
	@echo "Running linters..."
	golangci-lint run

# Build the application
.PHONY: build
build:
	@echo "Building application..."
	go build -o bin/translator ./cmd/translator

# Run the application
.PHONY: run
run:
	@echo "Running application..."
	go run ./cmd/translator

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	go test -bench=.

# Install development tools
.PHONY: install-tools
install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
