BINARY_NAME = secretsanta
BUILD_DIR = bin
CONFIG_FILE = configs/secretsanta.config.template
IMG_TAG ?= latest
DOCKER_COMPOSE = docker-compose
GO_VERSION = 1.24

# Multi-architecture build variables
PLATFORMS ?= linux/amd64,linux/arm64,linux/arm/v7
BUILDER_NAME ?= secretsanta-builder
REGISTRY ?= 
IMAGE_NAME ?= $(REGISTRY)secretsanta

# Build metadata
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
VCS_REF := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

ifeq ($(OS),Windows_NT)
	BINARY_NAME := $(BINARY_NAME).exe
	CP = copy
else
	CP = cp
endif

all: build copy-config

# Build targets
build:
	@echo "Building the CLI..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cli

build-web:
	@echo "Building the web server..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/secretsanta-web ./cmd/web

build-all: build build-web

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/cli
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/secretsanta-web-linux ./cmd/web

run-web:
	@echo "Starting web server..."
	go run ./cmd/web

copy-config:
	@echo "Copying config template..."
	@mkdir -p $(BUILD_DIR)
	$(CP) $(CONFIG_FILE) $(BUILD_DIR)/

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t secretsanta:$(IMG_TAG) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg VERSION="$(VERSION)" \
		--build-arg VCS_REF="$(VCS_REF)" \
		-f docker/Dockerfile .

docker-build-notifier:
	@echo "Building notifier service from external repo..."
	@if [ -d "../notifier" ]; then \
		cd ../notifier && docker build -t igodwin/notifier:latest .; \
	else \
		echo "Notifier repository not found at ../notifier"; \
		echo "Please clone the notifier service or pull the image: docker pull igodwin/notifier:latest"; \
	fi

# Multi-architecture Docker targets
docker-buildx-setup:
	@echo "Setting up Docker Buildx..."
	@docker buildx inspect $(BUILDER_NAME) >/dev/null 2>&1 || \
		docker buildx create --name $(BUILDER_NAME) --driver docker-container --bootstrap
	@docker buildx use $(BUILDER_NAME)

docker-buildx-build: docker-buildx-setup
	@echo "Building multi-architecture Docker image for platforms: $(PLATFORMS)"
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg VERSION="$(VERSION)" \
		--build-arg VCS_REF="$(VCS_REF)" \
		-t $(IMAGE_NAME):$(IMG_TAG) \
		-t $(IMAGE_NAME):latest \
		-f docker/Dockerfile \
		--push \
		.

docker-buildx-build-local: docker-buildx-setup
	@echo "Building multi-architecture Docker image locally (no push)"
	docker buildx build \
		--platform $(PLATFORMS) \
		--build-arg BUILD_DATE="$(BUILD_DATE)" \
		--build-arg VERSION="$(VERSION)" \
		--build-arg VCS_REF="$(VCS_REF)" \
		-t $(IMAGE_NAME):$(IMG_TAG) \
		-t $(IMAGE_NAME):latest \
		-f docker/Dockerfile \
		--load \
		.

docker-buildx-inspect:
	@echo "Inspecting multi-architecture image..."
	docker buildx imagetools inspect $(IMAGE_NAME):$(IMG_TAG)

docker-buildx-cleanup:
	@echo "Cleaning up Docker Buildx builder..."
	docker buildx rm $(BUILDER_NAME) || true

# Docker Compose targets
compose-up:
	@echo "Starting services with docker-compose..."
	$(DOCKER_COMPOSE) up -d notifier

compose-up-dev:
	@echo "Starting development environment..."
	$(DOCKER_COMPOSE) --profile dev up -d

compose-run:
	@echo "Running Secret Santa with notifier service..."
	$(DOCKER_COMPOSE) --profile run up --build secretsanta

compose-down:
	@echo "Stopping all services..."
	$(DOCKER_COMPOSE) down

compose-logs:
	@echo "Showing logs..."
	$(DOCKER_COMPOSE) logs -f

compose-logs-notifier:
	@echo "Showing notifier logs..."
	$(DOCKER_COMPOSE) logs -f notifier

# Development targets
test:
	@echo "Running tests..."
	go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

format:
	@echo "Formatting code..."
	go fmt ./...

mod-tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

clean-docker:
	@echo "Cleaning Docker containers and images..."
	$(DOCKER_COMPOSE) down --rmi all --volumes --remove-orphans

# Help target
help:
	@echo "Available targets:"
	@echo "  build              - Build the CLI binary"
	@echo "  build-web          - Build the web server"
	@echo "  build-all          - Build both CLI and web server"
	@echo "  build-linux        - Build for Linux (cross-compile)"
	@echo "  run-web            - Run web server in development mode"
	@echo "  copy-config        - Copy config template to build directory"
	@echo "  docker-build       - Build Docker image for secretsanta"
	@echo "  docker-build-notifier - Build notifier service Docker image"
	@echo "  docker-buildx-build - Build multi-architecture images"
	@echo "  compose-up         - Start notifier service only"
	@echo "  compose-up-dev     - Start development environment"
	@echo "  compose-run        - Run secretsanta with notifier service"
	@echo "  compose-down       - Stop all services"
	@echo "  compose-logs       - Show logs from all services"
	@echo "  compose-logs-notifier - Show logs from notifier service"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  lint               - Run linter"
	@echo "  format             - Format code"
	@echo "  mod-tidy           - Tidy go modules"
	@echo "  clean              - Clean build artifacts"
	@echo "  clean-docker       - Clean Docker containers and images"
	@echo "  help               - Show this help message"

.PHONY: all build build-web build-all build-linux run-web copy-config docker-build docker-build-notifier \
        docker-buildx-setup docker-buildx-build docker-buildx-build-local docker-buildx-inspect docker-buildx-cleanup \
        compose-up compose-up-dev compose-run compose-down compose-logs compose-logs-notifier \
        test test-coverage lint format mod-tidy clean clean-docker help