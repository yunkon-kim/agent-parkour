.PHONY: all build test clean install audit-example

BINARY_NAME=token-hop
ALIAS_NAME=thop
BIN_DIR=bin

all: test build

build:
	@mkdir -p $(BIN_DIR)
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/token-hop
	@ln -sf $(BINARY_NAME) $(BIN_DIR)/$(ALIAS_NAME)
	@echo "✅ Build complete: $(BIN_DIR)/$(BINARY_NAME) and alias $(BIN_DIR)/$(ALIAS_NAME)"

test:
	@echo "🧪 Running unit & integration tests..."
	@go test -v ./...

audit-example: build
	@echo "📊 Running token budget audit on cm-beetle fixture..."
	@./$(BIN_DIR)/$(ALIAS_NAME) audit --input test/fixtures/cm-beetle/.github

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) test/output /tmp/token-hop-*
	@echo "✅ Clean complete."

install: build
	@echo "📦 Installing $(BINARY_NAME) to $(GOBIN)..."
	@cp $(BIN_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/
	@cp $(BIN_DIR)/$(ALIAS_NAME) $(shell go env GOPATH)/bin/
	@echo "✅ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"
