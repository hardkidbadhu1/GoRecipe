SHELL := /bin/bash

APP_NAME := go-recipe-api
BUILD_DIR := build

.PHONY: build run swagger

build:
	@echo "==> Building the application..."
	go mod tidy
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	@echo "==> Running the application..."
	./$(BUILD_DIR)/$(APP_NAME)

swagger:
	@echo "==> Generating swagger docs..."
	swag init -g main.go

clean:
	@echo "==> Cleaning up..."
	rm -rf $(BUILD_DIR)
