# Define variables for output names
SERVER_BIN=bin/server
DIST_DIR=dist/
DIST_CONTENTS=$(DIST_DIR)*
WASM_BIN=dist/main.wasm

.PHONY: all serve clean

# Default target: build everything
all: server wasm

# Build the Server (Native)
server:
	go build -o $(SERVER_BIN) ./cmd/server

# Build the Game (WebAssembly)
wasm:
	GOOS=js GOARCH=wasm go build -o $(WASM_BIN) ./cmd/wasm
	cp ./cmd/wasm/static/* $(DIST_DIR)

# Build and Run
serve: all
	@echo "Starting server..."
	./$(SERVER_BIN)

# Clean up build artifacts
clean:
	rm -f $(SERVER_BIN) $(DIST_CONTENTS)
