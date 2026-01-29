# Variables
BINARY_NAME=server
WASM_OUT=web/app.wasm
SERVER_DIR=./cmd/server
WASM_DIR=./cmd/wasm
APP_PKG=github.com/maxence-charriere/go-app/v10/pkg/app

.PHONY: all build wasm server clean deps run

# Default target: build everything
all: build

# New target to install/update the go-app dependency
deps:
	@echo "Installing go-app v10..."
	go get -u $(APP_PKG)
	go mod tidy

# Cleanup unused dependencies
tidy:
	go mod tidy

# Build the Frontend (WebAssembly)
# Note: GOOS=js and GOARCH=wasm are required for go-app to run in the browser
wasm:
	@echo "Building WebAssembly..."
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) $(WASM_DIR)

# Build the Backend (Server)
server:
	@echo "Building Server..."
	go build -o $(BINARY_NAME) $(SERVER_DIR)

build: wasm server

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