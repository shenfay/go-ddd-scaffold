# Enterprise DDD Monolith Template Makefile

.PHONY: help clean build build-worker test run setup install-deps short unit integration e2e smoke up down api worker asynqmon gen serve deps check monitoring status install

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
	@echo "  make run api             - Start API server (development mode)"
	@echo "  make run worker          - Start Worker (development mode)"
	@echo "  make run asynqmon PORT=  - Start asynqmon UI with specified port"
	@echo "  make setup               - Setup development environment"
	@echo "  make install deps        - Install dependencies"
	@echo ""
	@echo "Scripts:"
	@echo "  make start               - Start all services (via scripts)"
	@echo "  make health check        - Check service health"
	@echo "  make start monitoring    - Start Prometheus + Grafana"
	@echo ""
	@echo "Database:"
	@echo "  make migrate up          - Run database migrations"
	@echo "  make migrate down        - Rollback database migrations"
	@echo "  make db status           - Show migration status"
	@echo ""
	@echo "Testing:"
	@echo "  make test                - Run all tests"
	@echo "  make test short          - Run short tests (skip integration)"
	@echo "  make test unit           - Run unit tests only"
	@echo "  make test integration    - Run integration tests (requires Docker)"
	@echo "  make test e2e            - Run E2E core flow tests"
	@echo "  make test smoke          - Run smoke tests (requires running API)"
	@echo "  make coverage            - Generate test coverage report"
	@echo ""
	@echo "CLI Tool:"
	@echo "  make cli                 - Build CLI tool for current OS"
	@echo "  make cli install         - Install CLI tool to GOPATH/bin"
	@echo ""
	@echo "Build:"
	@echo "  make build               - Build application for current OS"
	@echo "  make build worker        - Build Worker for current OS"
	@echo "  make build linux         - Build for Linux (production)"
	@echo "  make clean               - Clean build artifacts"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt                 - Format code"
	@echo "  make vet                 - Run go vet"
	@echo "  make lint                - Run code linting"
	@echo ""
	@echo "Documentation:"
	@echo "  make swagger gen         - Generate Swagger documentation"
	@echo "  make swagger serve       - Generate and serve Swagger UI"
	@echo ""
	@echo "Monitoring:"
	@echo "  make health              - Quick health check"
	@echo "  make health check        - Full health check (all services)"

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
	@if [ -z "$(filter-out cli,$(MAKECMDGOALS))" ]; then \
		echo "Building CLI tool..."; \
		go build -o bin/$(CLI_BINARY) $(CLI_CMD); \
		echo "CLI built: bin/$(CLI_BINARY)"; \
	elif [ "$(filter-out cli,$(MAKECMDGOALS))" = "install" ]; then \
		echo "Building CLI tool..."; \
		go build -o bin/$(CLI_BINARY) $(CLI_CMD); \
		echo "Installing CLI to GOPATH/bin..."; \
		cp bin/$(CLI_BINARY) $(GOPATH)/bin/; \
		echo "CLI installed. You can now use '$(CLI_BINARY)' command"; \
	else \
		echo "Usage: make cli [install]"; \
		exit 1; \
	fi

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
	@if [ "$(filter-out run,$(MAKECMDGOALS))" = "api" ]; then \
		echo "Starting API server..."; \
		go run ./cmd/api/; \
	elif [ "$(filter-out run,$(MAKECMDGOALS))" = "worker" ]; then \
		echo "Starting Asynq Worker..."; \
		go run ./cmd/worker/; \
	elif [ "$(filter-out run,$(MAKECMDGOALS))" = "asynqmon" ]; then \
		if [ -z "$(PORT)" ]; then \
			echo "Usage: make run asynqmon PORT=<port>"; \
			echo "Example: make run asynqmon PORT=8081"; \
			exit 1; \
		fi; \
		echo "Starting asynqmon UI on http://localhost:$(PORT)..."; \
		echo "Press Ctrl+C to stop"; \
		$$HOME/go/bin/asynqmon --redis-addr=localhost:6379 --port=$(PORT) || (echo "asynqmon not found. Run 'make asynqmon-install' first." && exit 1); \
	else \
		echo "Usage: make run <api|worker|asynqmon>"; \
		exit 1; \
	fi

# Install dependencies
install:
	@if [ "$(filter-out install,$(MAKECMDGOALS))" = "deps" ]; then \
		echo "Installing dependencies..."; \
		go mod download; \
		go mod tidy; \
	elif [ "$(filter-out install,$(MAKECMDGOALS))" = "" ]; then \
		echo "Usage: make install deps"; \
		exit 1; \
	fi

# Prevent Make from treating deps as target
deps:
	@:

# Build application for current OS
build:
	@if [ -z "$(filter-out build,$(MAKECMDGOALS))" ]; then \
		echo "Building application..."; \
		go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) ./cmd/api; \
	elif [ "$(filter-out build,$(MAKECMDGOALS))" = "worker" ]; then \
		echo "Building Worker..."; \
		go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-worker ./cmd/worker; \
		echo "Worker built: bin/$(APP_NAME)-worker"; \
	elif [ "$(filter-out build,$(MAKECMDGOALS))" = "linux" ]; then \
		echo "Building for Linux (production)..."; \
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-linux ./cmd/api; \
	else \
		echo "Usage: make build [worker|linux]"; \
		exit 1; \
	fi

# Prevent Make from treating worker/linux as targets
worker linux:
	@:

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf coverage.out
	rm -rf coverage.html

# Run tests
test:
	@if [ -z "$(filter-out test,$(MAKECMDGOALS))" ]; then \
		echo "Running all tests..."; \
		go test -v ./...; \
	elif [ "$(filter-out test,$(MAKECMDGOALS))" = "short" ]; then \
		echo "Running short tests..."; \
		go test -v -short ./...; \
	elif [ "$(filter-out test,$(MAKECMDGOALS))" = "unit" ]; then \
		echo "Running unit tests..."; \
		go test -v ./internal/...; \
	elif [ "$(filter-out test,$(MAKECMDGOALS))" = "integration" ]; then \
		echo "Running integration tests..."; \
		go test -v ./test/integration/...; \
	elif [ "$(filter-out test,$(MAKECMDGOALS))" = "e2e" ]; then \
		echo "Running E2E core flow tests..."; \
		cd ../scripts/testing && bash core-flow-test.sh; \
	elif [ "$(filter-out test,$(MAKECMDGOALS))" = "smoke" ]; then \
		echo "Running smoke tests..."; \
		cd ../scripts/testing && bash core-flow-test.sh; \
	else \
		echo "Usage: make test [short|unit|integration|e2e|smoke]"; \
		exit 1; \
	fi

# Prevent Make from treating test subcommands as targets
short unit integration e2e smoke:
	@:

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

# Database migrations
migrate:
	@if [ "$(filter-out migrate,$(MAKECMDGOALS))" = "up" ]; then \
		echo "Running database migrations..."; \
		go run ./cmd/cli migrate up; \
		echo "✓ Database migrations completed successfully"; \
	elif [ "$(filter-out migrate,$(MAKECMDGOALS))" = "down" ]; then \
		echo "Rolling back database migrations..."; \
		go run ./cmd/cli migrate down; \
		echo "✓ Database rollback completed successfully"; \
	else \
		echo "Usage: make migrate <up|down>"; \
		exit 1; \
	fi

# Prevent Make from treating subcommands as targets
up down:
	@:

api asynqmon:
	@:

gen serve:
	@:

# Swagger documentation
swagger:
	@if [ "$(filter-out swagger,$(MAKECMDGOALS))" = "gen" ]; then \
		echo "Generating Swagger documentation..."; \
		swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
		$(GOBIN)/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
		~/go/bin/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal; \
		echo "Swagger documentation generated at api/swagger/"; \
		echo "Start server and visit: http://localhost:8080/swagger/index.html"; \
	elif [ "$(filter-out swagger,$(MAKECMDGOALS))" = "serve" ]; then \
		echo "Generating Swagger documentation..."; \
		swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
		$(GOBIN)/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal || \
		~/go/bin/swag init -g cmd/api/main.go -o api/swagger --parseDependency --parseInternal; \
		echo "Starting Swagger UI server..."; \
		go run ./cmd/api/; \
	else \
		echo "Usage: make swagger <gen|serve>"; \
		exit 1; \
	fi

# Check application health
health:
	@if [ -z "$(filter-out health,$(MAKECMDGOALS))" ]; then \
		echo "Checking application health..."; \
		curl -s -f http://localhost:8080/health || (echo "\nApplication is not running or unhealthy" && exit 1); \
		echo "\n✓ Application is healthy"; \
	elif [ "$(filter-out health,$(MAKECMDGOALS))" = "check" ]; then \
		cd .. && ./scripts/services/health-check.sh; \
	else \
		echo "Usage: make health [check]"; \
		exit 1; \
	fi

# Start all services
start:
	@if [ -z "$(filter-out start,$(MAKECMDGOALS))" ]; then \
		cd .. && ./scripts/services/start.sh; \
	elif [ "$(filter-out start,$(MAKECMDGOALS))" = "monitoring" ]; then \
		cd .. && ./scripts/monitoring/start-monitoring.sh; \
	else \
		echo "Usage: make start [monitoring]"; \
		exit 1; \
	fi

# Database migration status
db:
	@if [ "$(filter-out db,$(MAKECMDGOALS))" = "status" ]; then \
		cd .. && ./scripts/database/migrate.sh status; \
	else \
		echo "Usage: make db status"; \
		exit 1; \
	fi

# Prevent Make from treating check/monitoring/status as targets
check monitoring status:
	@:

# Run all tests (unit + integration + E2E)
test-all:
	@cd .. && ./scripts/testing/run-all-tests.sh

# Run E2E core flow tests
test-e2e:
	@cd .. && ./scripts/testing/core-flow-test.sh

# Install asynqmon CLI tool
asynqmon-install:
	@echo "Installing asynqmon..."
	go install github.com/hibiken/asynqmon/cmd/asynqmon@latest
	@echo "asynqmon installed successfully!"
	@echo "Run 'make run asynqmon PORT=<port>' to start the UI"

.DEFAULT_GOAL := help