# Whalio Makefile

# Variables
APP_NAME=whalio
CMD_DIR=./cmd
MAIN_FILE=$(CMD_DIR)/main.go
BINARY_DIR=./bin
STATIC_DIR=./static
TEMPLATES_DIR=./templates
CSS_INPUT=./assets/css/input.css
CSS_OUTPUT=$(STATIC_DIR)/css/output.css

# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run

# Node/bun variables
BUN=bun
BUNX=bunx

# Default target
.PHONY: help
help: ## Show this help message
	@echo "Whalio Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# Setup and dependencies
.PHONY: setup
setup: ## Setup development environment
	@echo "🔧 Setting up development environment..."
	$(BUN) install
	$(GOMOD) tidy
	$(GOGET) github.com/a-h/templ/cmd/templ@latest
	@echo "✅ Setup complete!"

.PHONY: deps
deps: ## Download Go dependencies
	@echo "📦 Downloading Go dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## Update Go dependencies
	@echo "🔄 Updating Go dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Development
.PHONY: dev
dev: ## Run development server with hot reloading
	@echo "🚀 Starting development server..."
	$(BUN) run dev

.PHONY: run
run: build-quick ## Build and run the application
	@echo "🏃 Running application..."
	$(BINARY_DIR)/$(APP_NAME)

.PHONY: dev-go
dev-go: ## Run Go server only (without frontend watching)
	@echo "🐹 Starting Go server..."
	$(GORUN) $(MAIN_FILE)

# Building
.PHONY: build
build: clean deps generate build-css build-go ## Full build (clean + deps + generate + css + binary)
	@echo "✅ Build complete!"

.PHONY: build-quick
build-quick: generate build-css build-go ## Quick build (no clean, no deps)
	@echo "⚡ Quick build complete!"

.PHONY: build-go
build-go: ## Build Go binary
	@echo "🔨 Building Go binary..."
	@mkdir -p $(BINARY_DIR)
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_DIR)/$(APP_NAME) $(MAIN_FILE)

.PHONY: build-css
build-css: ## Build CSS with TailwindCSS
	@echo "🎨 Building CSS..."
	@mkdir -p $(STATIC_DIR)/css
	$(BUNX) tailwindcss -i $(CSS_INPUT) -o $(CSS_OUTPUT) --minify

.PHONY: watch-css
watch-css: ## Watch and rebuild CSS files
	@echo "👀 Watching CSS files..."
	$(BUNX) tailwindcss -i $(CSS_INPUT) -o $(CSS_OUTPUT) --watch

# Templates
.PHONY: generate
generate: ## Generate templ templates
	@echo "📝 Generating templ templates..."
	templ generate

.PHONY: watch-templ
watch-templ: ## Watch and regenerate templ templates
	@echo "👀 Watching templ files..."
	templ generate --watch --proxy=http://localhost:8080

# Testing
.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "🏁 Running tests with race detection..."
	$(GOTEST) -race -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "📈 Coverage report generated: coverage.html"

# Linting and formatting
.PHONY: fmt
fmt: ## Format Go code
	@echo "🎯 Formatting Go code..."
	$(GOCMD) fmt ./...
	templ fmt .

.PHONY: lint
lint: ## Run linters
	@echo "🔍 Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "🔎 Running go vet..."
	$(GOCMD) vet ./...

# Cleaning
.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -rf $(STATIC_DIR)/css/output.css
	rm -f coverage.out coverage.html

.PHONY: clean-deps
clean-deps: ## Clean dependencies
	@echo "🧹 Cleaning dependencies..."
	$(GOCLEAN) -modcache
	rm -rf node_modules

# Production
.PHONY: build-prod
build-prod: clean deps generate ## Build for production
	@echo "🏭 Building for production..."
	@mkdir -p $(BINARY_DIR)
	$(BUNX) tailwindcss -i $(CSS_INPUT) -o $(CSS_OUTPUT) --minify
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags="-w -s -X main.version=$$(git describe --tags --always --dirty)" \
		-o $(BINARY_DIR)/$(APP_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "✅ Production build complete!"

# Docker (if needed)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):latest

# Utilities
.PHONY: check
check: fmt vet lint test ## Run all checks (format, vet, lint, test)
	@echo "✅ All checks passed!"

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	$(GOGET) github.com/a-h/templ/cmd/templ@latest
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "✅ Tools installed!"

.PHONY: serve-static
serve-static: build-css ## Serve static files for development
	@echo "📁 Serving static files on http://localhost:8080..."
	@cd $(STATIC_DIR) && python3 -m http.server 8080

.PHONY: info
info: ## Show project information
	@echo "📋 Project Information:"
	@echo "  Name: $(APP_NAME)"
	@echo "  Go version: $$($(GOCMD) version)"
	@echo "  Node version: $$(node --version 2>/dev/null || echo 'Not installed')"
	@echo "  NPM version: $$(npm --version 2>/dev/null || echo 'Not installed')"
	@echo "  Templ version: $$(templ version 2>/dev/null || echo 'Not installed')"
	@echo ""
	@echo "📂 Project structure:"
	@echo "  Binary: $(BINARY_DIR)/"
	@echo "  Static: $(STATIC_DIR)/"
	@echo "  Templates: $(TEMPLATES_DIR)/"
	@echo "  Main: $(MAIN_FILE)"