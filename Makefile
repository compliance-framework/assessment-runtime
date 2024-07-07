# Variables
BINARY_NAME=ar
CONFIG=configs/config.yaml
GO=go

.PHONY: help build test clean fmt vet lint build-images run-docker run-local protoc graph

help:  ## Display this help message
	@echo "Help for Makefile: $(MAKEFILE_LIST) in $(dir $(abspath $(lastword $(MAKEFILE_LIST))))"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build:  ## Build the Go application
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o ./bin/$(BINARY_NAME) ./

test:  ## Run unit tests
	@echo "Running tests..."
	@$(GO) test -v ./...

clean:  ## Clean up binaries and test caches
	@echo "Cleaning up..."
	@$(GO) clean
	@rm -rf ./bin
	@rm -rf test-results

fmt:  ## Format the code
	@echo "Formatting code..."
	@$(GO) fmt ./...

vet:  ## Vet the code
	@echo "Running vet..."
	@$(GO) vet ./...

lint:  ## Lint the code
	@echo "Running lint..."
	@golint ./...

build-local:  ## Build image to be used locally
	docker build -t ghcr.io/compliance-framework/assessment-runtime:latest_local -f ./test/runtime/Dockerfile .

build-local-ar:  ## Build assessment-runtime
	docker build -t ghcr.io/compliance-framework/assessment-runtime:latest_local -f ./test/runtime/Dockerfile .

build-images:  ## Build Docker images
	@echo "Building Docker images..."
	docker build -t assessment-runtime -f test/runtime/Dockerfile .

run-docker:  ## Run the test environment using Docker Compose
	@echo "Running Docker Compose..."
	docker compose -p argus -f ./test/docker-compose.yml up --build -d

stop-docker:  ## Stop the test environment
	@echo "Stopping Docker Compose..."
	docker compose -p argus -f ./test/docker-compose.yml down

run-local:   ## Build and run the application locally
	@echo "Building and running $(BINARY_NAME) locally..."
	@$(GO) build -o ./bin/$(BINARY_NAME) ./
	./bin/$(BINARY_NAME)

protoc:  ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative provider/job.proto

graph:  ## Generate dependency graph
	@echo "Generating dependency graph..."
	godepgraph -p google.golang.org,github.com/sirupsen,github.com/hashicorp,github.com/robfig,gopkg.in -s ./internal | dot -Tpng -o godepgraph.png
