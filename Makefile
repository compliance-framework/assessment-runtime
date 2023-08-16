# Variables
BINARY_NAME=ar
CONFIG=configs/config.yaml
GO=go

.PHONY: all build test clean

# Default target to build the application
all: build

# Build the Go application
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o ./bin/$(BINARY_NAME) ./cmd/

# Run unit tests
test:
	@echo "Running tests..."
	@$(GO) test -v ./pkg/...

# Clean up binaries and test caches
clean:
	@echo "Cleaning up..."
	@$(GO) clean
	@rm -rf ./cmd/bin
	@rm -rf test-results

fmt:
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet:
	@echo "Running vet..."
	@$(GO) vet ./...

lint:
	@echo "Running lint..."
	@golint ./...
