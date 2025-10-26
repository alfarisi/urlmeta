.PHONY: test test-coverage test-race lint fmt vet build clean install-tools help examples bench

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-15s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## test-race: Run tests with race detector
test-race:
	@echo "Running tests with race detector..."
	go test -v -race ./...

## bench: Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	~/go/bin/golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## build: Build the project
build:
	@echo "Building..."
	go build -v ./...

## build-examples: Build all examples
build-examples:
	@echo "Building examples..."
	cd examples/basic && go build -v .
	cd examples/advanced && go build -v .
	cd examples/batch && go build -v .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	go clean
	rm -f coverage.out coverage.html
	rm -f examples/basic/basic
	rm -f examples/advanced/advanced
	rm -f examples/batch/batch

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	go mod tidy
	go mod verify

## mod-update: Update dependencies
mod-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## check: Run all checks (fmt, vet, lint, test)
check: fmt vet lint test
	@echo "All checks passed!"

## ci: Run CI pipeline locally
ci: check test-coverage
	@echo "CI pipeline completed!"
