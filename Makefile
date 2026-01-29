# Variables
BINARY_NAME=server
WASM_OUT=web/app.wasm
SERVER_DIR=./cmd/server
WASM_DIR=./cmd/wasm

.PHONY: all build wasm server clean run

# Default target: build everything
all: build

build: wasm server

# 1. Build the Frontend (WebAssembly)
# Note: GOOS=js and GOARCH=wasm are required for go-app to run in the browser
wasm:
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) $(WASM_DIR)

# 2. Build the Backend (Server)
server:
	@echo "Building Server..."
	go build -o $(BINARY_NAME) $(SERVER_DIR)

# Run the application
run: build
	@echo "Starting server at http://localhost:8080"
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f $(WASM_OUT)

# Helpful for development (requires 'air' to be installed)
dev:
	@echo "Starting hot-reload development mode..."
	air