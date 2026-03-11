# Enterprise DDD Monolith Template Makefile

.PHONY: help clean build test run dev-setup deploy

# Variables
APP_NAME = enterprise-monolith
VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_TIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT = $(shell git rev-parse HEAD)
GO_LDFLAGS = -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Help target
help:
	@echo "Enterprise DDD Monolith Template - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make setup        - Setup development environment"
	@echo "  make dev          - Start development server"
	@echo "  make test         - Run all tests"
	@echo "  make test-unit    - Run unit tests"
	@echo "  make test-int     - Run integration tests"
	@echo "  make coverage     - Generate test coverage report"
	@echo ""
	@echo "Build:"
	@echo "  make build        - Build application"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "Deployment:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-push  - Push Docker image"
	@echo "  make deploy       - Deploy to Kubernetes"
	@echo ""
	@echo "Tools:"
	@echo "  make lint         - Run code linting"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Run go vet"
	@echo "  make generate     - Generate code"
	@echo "  make validate     - Validate architecture compliance"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vektra/mockery/v2@latest
	docker-compose -f deployments/docker/docker-compose.dev.yml up -d
	@echo "Development environment ready!"

# Development server
dev:
	@echo "Starting development server..."
	go run cmd/api/main.go

# Build application
build:
	@echo "Building application..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "$(GO_LDFLAGS)" -o bin/$(APP_NAME) cmd/api/main.go

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

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

# Run integration tests
test-int:
	@echo "Running integration tests..."
	go test -v -tags=integration ./...

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

# Code generation
generate:
	@echo "Generating code..."
	go generate ./...
	mockery --all --dir internal/domain --output internal/mocks

# Architecture validation
validate:
	@echo "Validating architecture compliance..."
	go run tools/validator/main.go

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) -f deployments/docker/Dockerfile .

# Docker push
docker-push:
	@echo "Pushing Docker image..."
	docker push $(APP_NAME):$(VERSION)

# Deploy to Kubernetes
deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/kubernetes/

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	go run cmd/cli/main.go migrate up

migrate-down:
	@echo "Rolling back database migrations..."
	go run cmd/cli/main.go migrate down

# Health check
health-check:
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || exit 1

.DEFAULT_GOAL := help