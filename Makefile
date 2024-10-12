BINARY_NAME = secretsanta
BUILD_DIR = bin
CONFIG_FILE = configs/secretsanta.config.template
IMG_TAG ?= latest

ifeq ($(OS),Windows_NT)
	BINARY_NAME := $(BINARY_NAME).exe
endif

all: build copy-config

build:
	@echo "Building the project..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/cli

copy-config:
	@echo "Copying config template..."
	$(CP) $(CONFIG_FILE) $(BUILD_DIR)/

docker-build:
	@echo "Building Docker image..."
	docker build -t secretsanta:$(IMG_TAG) -f docker/Dockerfile .

.PHONY: all build copy-config docker-build