REPO_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(REPO_DIR)/bin/fractal-func

.PHONY: help build run lint format
.DEFAULT_GOAL := help

help: ## 💬 This help message :)
	@figlet $@ 2> /dev/null || echo "***** Running $@ *****"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

build: ## 🔨 Build server/listener executable
	@figlet $@ 2> /dev/null || echo "***** Running $@ *****"
	rm -f $(BINARY_NAME)
	@go build -o $(BINARY_NAME) fractal-bot-go/src

run: build ## ⚡ Build + run with local Function host
	@figlet $@ 2> /dev/null || echo "***** Running $@ *****"
	func start

lint: ## 🧹 Lint the code
	@figlet $@ 2> /dev/null || echo "***** Running $@ *****"
	golangci-lint run -E revive,gofmt,misspell

format: ## 📝 Format the code
	@figlet $@ 2> /dev/null || echo "***** Running $@ *****"
	gofmt -w -s $(REPO_DIR)/src