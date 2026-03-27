# Enterprise DDD Monolith Template Makefile

.PHONY: help clean build build-worker test run run-worker dev setup install-deps

# Variables
APP_NAME = go-ddd-scaffold
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT = $(shell git rev-parse HEAD)
GO_LDFLAGS = -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# CLI Tool
CLI_BINARY=go-ddd-scaffold
CLI_CMD=cmd/cli/main.go

# Help target
help:
	@echo "Go DDD Scaffold - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make run          - Start API server (development mode)"
	@echo "  make run-worker   - Start Worker (development mode)"
	@echo "  make setup        - Setup development environment"
	@echo "  make install-deps - Install dependencies"
	@echo ""
	@echo "asynqmon (Task Monitoring):"
	@echo "  make asynqmon-install   - Install asynqmon CLI tool"
	@echo "  make asynqmon           - Start asynqmon UI (default port 8080)"
	@echo "  make asynqmon-ui        - Open asynqmon UI in browser"
	@echo ""
	@echo "CLI Tool:"
	@echo "  make cli          - Build CLI tool for current OS"
	@echo "  make cli-linux    - Build CLI tool for Linux"
	@echo "  make cli-darwin   - Build CLI tool for macOS"
	@echo "  make cli-windows  - Build CLI tool for Windows"
	@echo "  make cli-install  - Install CLI tool to GOPATH/bin"
	@echo "  make cli-test     - Test CLI tool with sample command"
	@echo ""
	@echo "Build:"
	@echo "  make build        - Build application for current OS"
	@echo "  make build-worker - Build Worker for current OS"
	@echo "  make build-linux  - Build for Linux (production)"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Test:"
	@echo "  make test         - Run all tests"
	@echo "  make test-short   - Run short tests (skip integration)"
	@echo "  make coverage     - Generate test coverage report"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo "  make lint         - Run code linting"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo ""
	@echo "asynqmon (Task Monitoring):"
	@echo "  make asynqmon-install   - Install asynqmon CLI tool"
	@echo "  make asynqmon           - Start asynqmon UI (default port 8080)"
	@echo "  make asynqmon-port PORT=8081 - Start asynqmon on custom port"
	@echo "  make asynqmon-ui        - Open asynqmon UI in browser (macOS)"
	@echo ""
	@echo "Utilities:"
	@echo "  make health       - Check application health"
	@echo ""
	@echo "Documentation:"
	@echo "  make swagger-gen  - Generate Swagger documentation"
	@echo "  make swagger-serve - Generate and serve Swagger UI"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@latest
	@echo "Development environment ready!"

# ==============================================================================
# CLI Tool Targets
# ==============================================================================

.PHONY: cli cli-linux cli-darwin cli-windows cli-install cli-clean cli-test

# Build CLI tool for current OS
cli:
	@echo "Building CLI tool..."
	go build -o bin/$(CLI_BINARY) $(CLI_CMD)
	@echo "CLI built: bin/$(CLI_BINARY)"

# Build CLI for Linux
cli-linux:
	@echo "Building CLI for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/$(CLI_BINARY)-linux-amd64 $(CLI_CMD)
	@echo "CLI built: bin/$(CLI_BINARY)-linux-amd64"

# Build CLI for macOS
cli-darwin:
	@echo "Building CLI for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/$(CLI_BINARY)-darwin-amd64 $(CLI_CMD)
	@echo "CLI built: bin/$(CLI_BINARY)-darwin-amd64"

# Build CLI for Windows
cli-windows:
	@echo "Building CLI for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/$(CLI_BINARY)-windows-amd64.exe $(CLI_CMD)
	@echo "CLI built: bin/$(CLI_BINARY)-windows-amd64.exe"

# Install CLI to GOPATH/bin
cli-install: cli
	@echo "Installing CLI to GOPATH/bin..."
	cp bin/$(CLI_BINARY) $(GOPATH)/bin/
	@echo "CLI installed. You can now use '$(CLI_BINARY)' command"

# Test CLI
cli-test: cli
	@echo "Testing CLI..."
	./bin/$(CLI_BINARY) --version
	./bin/$(CLI_BINARY) generate dao --help
	@echo "CLI test passed!"

# Clean CLI binaries
cli-clean:
	@echo "Cleaning CLI binaries..."
	rm -f bin/$(CLI_BINARY) bin/$(CLI_BINARY)-*
	@echo "CLI binaries cleaned"

# Start API server (development mode)
run:
	@echo "Starting API server..."
	go run ./cmd/api/

# Start Worker (development mode)
run-worker:
	@echo "Starting Asynq Worker..."
	go run ./cmd/worker/

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build application for current OS
build:
	@echo "Building application..."
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) ./cmd/api

# Build Worker for current OS
build-worker:
	@echo "Building Worker..."
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-worker ./cmd/worker
	@echo "Worker built: bin/$(APP_NAME)-worker"

# Build for Linux (production)
build-linux:
	@echo "Building for Linux (production)..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-linux ./cmd/api

# Build Worker for Linux (production)
build-worker-linux:
	@echo "Building Worker for Linux (production)..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-worker-linux ./cmd/worker

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf coverage.out
	rm -rf coverage.html

# Run all tests
test:
	@echo "Running all tests..."
	go test -v ./...

# Run short tests (skip integration tests)
test-short:
	@echo "Running short tests..."
	go test -v -short ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code formatting
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Code vetting
vet:
	@echo "Running go vet..."
	go vet ./...

# Code linting
lint:
	@echo "Running linter..."
	golangci-lint run

# Swagger documentation
swagger-gen:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
	$(GOBIN)/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
	~/go/bin/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal
	@echo "Swagger documentation generated at api/swagger/"
	@echo "Start server and visit: http://localhost:8080/swagger/index.html"

swagger-serve: swagger-gen
	@echo "Starting Swagger UI server..."
	go run ./cmd/api/
	@echo "Swagger UI available at http://localhost:8080/swagger/index.html"

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	go run ./cmd/cli migrate up

migrate-down:
	@echo "Rolling back database migrations..."
	go run ./cmd/cli migrate down

# Check application health
health:
	@echo "Checking application health..."
	@curl -s -f http://localhost:8080/health || (echo "\nApplication is not running or unhealthy" && exit 1)
	@echo "\n✓ Application is healthy"

# === asynqmon Commands ===

# Install asynqmon CLI tool
asynqmon-install:
	@echo "Installing asynqmon..."
	go install github.com/hibiken/asynqmon/cmd/asynqmon@latest
	@echo "asynqmon installed successfully!"
	@echo "Run 'make asynqmon' to start the UI"

# Start asynqmon UI
asynqmon:
	@echo "Starting asynqmon UI on http://localhost:8080..."
	@echo "Press Ctrl+C to stop"
	$$HOME/go/bin/asynqmon --redis-addr=localhost:6379 || (echo "asynqmon not found. Run 'make asynqmon-install' first." && exit 1)

# Start asynqmon on custom port
asynqmon-port:
	@echo "Starting asynqmon UI on http://localhost:$(PORT)..."
	$$HOME/go/bin/asynqmon --redis-addr=localhost:6379 --port=$(PORT) || (echo "asynqmon not found. Run 'make asynqmon-install' first." && exit 1)

# Open asynqmon UI in browser (macOS only)
asynqmon-ui:
	@echo "Opening asynqmon UI in browser..."
	open http://localhost:8080

.DEFAULT_GOAL := help