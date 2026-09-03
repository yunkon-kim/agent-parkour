.PHONY: all build test test-e2e clean install audit-example

BINARY_NAME=agent-parkour
CLI_NAME=parkour
ALIAS_NAME=pk
BIN_DIR=bin

all: test build

build:
	@mkdir -p $(BIN_DIR)
	@echo "🔨 Building $(CLI_NAME)..."
	@go build -o $(BIN_DIR)/$(CLI_NAME) ./cmd/parkour
	@ln -sf $(CLI_NAME) $(BIN_DIR)/$(BINARY_NAME)
	@ln -sf $(CLI_NAME) $(BIN_DIR)/$(ALIAS_NAME)
	@ln -sf $(CLI_NAME) $(BIN_DIR)/thop
	@ln -sf $(CLI_NAME) $(BIN_DIR)/token-hop
	@echo "✅ Build complete: $(BIN_DIR)/$(CLI_NAME) (aliases: $(BINARY_NAME), $(ALIAS_NAME), thop)"

test:
	@echo "🧪 Running unit & integration tests..."
	@go test -v ./...

test-e2e:
	@echo "🏃 Running live E2E AI refinement test with Antigravity CLI (agy)..."
	@go test -v -tags=e2e -count=1 -timeout 180s ./test -run TestE2E_AntigravityRefinementWithAgy

audit-example: build
	@echo "📊 Running token audit on cm-beetle fixture..."
	@./$(BIN_DIR)/$(CLI_NAME) audit --input test/fixtures/cm-beetle/.github

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) test/output /tmp/parkour-* /tmp/token-hop-*
	@echo "✅ Clean complete."

install: build
	@echo "📦 Installing $(CLI_NAME) to $(shell go env GOPATH)/bin/..."
	@cp $(BIN_DIR)/$(CLI_NAME) $(shell go env GOPATH)/bin/
	@cp $(BIN_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/
	@cp $(BIN_DIR)/$(ALIAS_NAME) $(shell go env GOPATH)/bin/
	@echo "✅ Installed $(CLI_NAME) and $(ALIAS_NAME) to $(shell go env GOPATH)/bin/"
