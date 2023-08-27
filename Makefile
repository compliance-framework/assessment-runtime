# Variables
BINARY_NAME=ar
CONFIG=configs/config.yaml
GO=go

.PHONY: help build test clean fmt vet lint build-images compose-up compose-down start

help:  ## This help
	@echo "Help for Makefile: $(MAKEFILE_LIST) in $(dir $(abspath $(lastword $(MAKEFILE_LIST))))"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:  ## Build the Go application
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o ./bin/$(BINARY_NAME) ./

run:  ## Run the assessment-runtime binary
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run ./

test:  ## Run unit tests
	@echo "Running tests..."
	@$(GO) test -v ./...

clean:  ## Clean up binaries and test caches
	@echo "Cleaning up..."
	@$(GO) clean
	@rm -rf ./cmd/bin
	@rm -rf test-results

fmt:  ## Format the code
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet:  ## Vet the code
	@echo "Running vet..."
	@$(GO) vet ./...

lint:  ## Lint the code
	@echo "Running lint..."
	@golint ./...

build-images:  ## Build Docker images
	docker build -t plugin-registry ./tests/registry
	docker build -t assessment-runtime -f tests/runtime/Dockerfile .

compose-up:  ## Run up test environment
	docker compose -f ./tests/docker-compose.yml up --build

compose-down:  ## Bring down test environment
	docker compose -f ./tests/docker-compose.yml down

start:
	mkdir -p bin/plugins/sample/1.0.0
	go build -o bin/plugins/sample/1.0.0/sample ./tests/sampleplugin/main.go
	chmod +x bin/plugins/sample/1.0.0/sample
	@$(GO) build -o ./bin/$(BINARY_NAME) ./
	cp ./tests/runtime/config.yml ./bin/config.yml
