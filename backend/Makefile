# Enterprise DDD Monolith Template Makefile

.PHONY: help clean build test run dev setup install-deps

# Variables
APP_NAME = go-ddd-scaffold
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT = $(shell git rev-parse HEAD)
GO_LDFLAGS = -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Help target
help:
	@echo "Go DDD Scaffold - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make run          - Start API server (development mode)"
	@echo "  make setup        - Setup development environment"
	@echo "  make install-deps - Install dependencies"
	@echo ""
	@echo "Build:"
	@echo "  make build        - Build application for current OS"
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

# Start API server (development mode)
# Note: Uses 'go run ./cmd/api/' to include domains.go for route auto-registration
run:
	@echo "Starting API server..."
	@echo "Note: Using 'go run ./cmd/api/' to include all package files (main.go + domains.go)"
	go run ./cmd/api/

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Build application for current OS
build:
	@echo "Building application..."
	go build -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) ./cmd/api

# Build for Linux (production)
build-linux:
	@echo "Building for Linux (production)..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME)-linux ./cmd/api

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
	swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal || \
	$(GOBIN)/swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal || \
	~/go/bin/swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal
	@echo "Swagger documentation generated at docs/swagger/"
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

.DEFAULT_GOAL := help