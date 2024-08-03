# Makefile for clai project

# Directories
SRC_DIR := ./cmd

# Binaries
BIN_DIR := ./bin
CLAI_BIN := $(BIN_DIR)/clai-cli

# Go source directories for the binaries
CLAI_SRC := $(SRC_DIR)/cli

# Go build command
GO_BUILD := go build -o

# Default target
all: build

# Build both binaries
build: $(CLAI_BIN)

# Build clai-cli binary
$(CLAI_BIN):
	@echo "Building clai-cli binary..."
	@mkdir -p $(BIN_DIR)
	$(GO_BUILD) $(CLAI_BIN) $(CLAI_SRC)

# Clean the binaries
clean:
	@echo "Cleaning up binaries..."
	@rm -rf $(BIN_DIR)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy

.PHONY: all build clean deps