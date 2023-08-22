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
	@$(GO) build -o ./bin/$(BINARY_NAME) ./

run:
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run ./

# Run unit tests
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

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

build-images:
	docker build -t plugin-registry ./tests/registry
	docker build -t assessment-runtime -f tests/runtime/Dockerfile .

start:
	docker compose -d -f ./tests/docker-compose.yml up --build

stop:
	docker compose -f ./tests/docker-compose.yml down

start-local:
	mkdir -p bin/plugins/sample/1.0.0
	go build -o bin/plugins/sample/1.0.0/sample ./tests/sampleplugin/main.go
	chmod +x bin/plugins/sample/1.0.0/sample
	@$(GO) build -o ./bin/$(BINARY_NAME) ./
	cp ./tests/runtime/config.yaml ./bin/config.yaml
