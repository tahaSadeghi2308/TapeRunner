# Project settings
PROJECT_NAME=tape-runner
PKG=./...
MAIN=cmd/main.go
BUILD_DIR=build

# Git info
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Docker setting
TEST_TAG=dev

# Default target
run: ## Run the app
	@clear
	@go run $(MAIN)

all: build ## Build the project (default)

compile: ## Compile (without installing)
	@go build -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN)

build: clean compile
	./$(BUILD_DIR)/$(PROJECT_NAME)

test: ## Run tests
	go test $(PKG)

dep: ## Download and vendor dependencies
	go mod tidy

swag: ## Generate Swagger docs (requires github.com/swaggo/swag)
	@which swag >/dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g internal/delivery/web/setup.go -o docs
	@echo "Swagger docs generated in ./docs"

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

push: ## Push current branch to origin
	git push origin $(BRANCH)

pull: ## Pull current branch from origin
	git pull origin $(BRANCH)

help: ## Show this help
	@echo "Usage: make [target]"
	@awk 'BEGIN {FS = ":.*##"; printf "\nAvailable targets:\n"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)