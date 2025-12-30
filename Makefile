GO_CMD ?= go

DOCKER_BUILDX_CMD ?= docker buildx
IMAGE_BUILD_CMD ?= $(DOCKER_BUILDX_CMD) build

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
BIN_DIR ?= $(PROJECT_DIR)/bin

.PHONY: build-manager
build-manager:
	$(GO_CMD) build -o bin/sandcastle-manager cmd/sandcastle-manager/main.go

.PHONY: build-worker
build-worker:
	$(GO_CMD) build -o bin/sandcastle-worker cmd/sandcastle-worker/main.go
