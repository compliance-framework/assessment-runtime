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
	docker build -t plugin-registry ./test/registry
	docker build -t assessment-runtime -f test/runtime/Dockerfile .

up:  ## Run up test environment
	docker compose -f ./test/docker-compose.yml up --build

down:  ## Bring down test environment
	docker compose -f ./test/docker-compose.yml down

build-plugin: ## Build sample plugin and copy to bin along with config and assessment
	rm -rf bin
	mkdir -p bin/plugins/busy/1.0.0
	mkdir -p bin/plugins/hello/1.0.0
	mkdir -p bin/assessments
	cp ./test/runtime/config.yml ./bin/config.yml
	cp ./test/runtime/assessments/assess-1234.yaml ./bin/assessments/assessment-1234.yaml
	@$(GO) build -o bin/plugins/busy/1.0.0/busy ./test/plugins/busy.go
	@$(GO) build -o bin/plugins/hello/1.0.0/hello ./test/plugins/hello.go
	chmod +x bin/plugins/busy/1.0.0/busy
	chmod +x bin/plugins/hello/1.0.0/hello

start: build-plugin
	@$(GO) build -o ./bin/$(BINARY_NAME) ./
	./bin/$(BINARY_NAME)

protoc: ## Generate protobuf files
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative plugins/action.proto

graph: ## Generate dependency graph, requires godepgraph and graphviz
	godepgraph -p google.golang.org,github.com/sirupsen,github.com/hashicorp,github.com/robfig,gopkg.in -s ./internal | dot -Tpng -o godepgraph.png
