# Makefile for DDD-Scaffold Backend

.PHONY: help build test run docker-build docker-run swagger migrate clean

# 默认目标
.DEFAULT_GOAL := help

# 显示帮助信息
help: ## 显示帮助信息
	@echo "DDD-Scaffold Backend - Makefile Commands"
	@echo ""
	@echo "Usage: make [command]"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# 开发环境
dev: ## 启动开发环境（热重载）
	@echo "Starting development server with hot reload..."
	go run github.com/githubnemo/CompileDaemon -build="go build -o main cmd/server/main.go" -command="./main"

run: ## 运行服务
	@echo "Running server..."
	go run cmd/server/main.go

build: ## 编译二进制文件
	@echo "Building binary..."
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	@echo "Binary built: bin/server"

test: ## 运行所有测试
	@echo "Running tests..."
	go test -v ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-unit: ## 运行单元测试
	@echo "Running unit tests..."
	go test -v ./internal/domain/... -race

test-integration: ## 运行集成测试
	@echo "Running integration tests..."
	go test -v ./tests/integration/... -race

test-bench: ## 运行性能基准测试
	@echo "Running benchmark tests..."
	go test -bench=. -benchmem ./internal/infrastructure/auth/... -run=^$

swagger: ## 生成 Swagger API 文档
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o ./docs
	@echo "Swagger docs generated successfully!"
	@echo "Access at: http://localhost:8080/swagger/index.html"

swagger-serve: ## 启动 Swagger UI 服务器
	@echo "Starting Swagger UI server..."
	swag serve --dir ./docs --port 8081

migrate-up: ## 执行数据库迁移（向上）
	@echo "Running database migrations (up)..."
	go run cmd/migrate/main.go up

migrate-down: ## 回滚数据库迁移（向下）
	@echo "Rolling back database migrations (down)..."
	go run cmd/migrate/main.go down

migrate-status: ## 查看迁移状态
	@echo "Checking migration status..."
	go run cmd/migrate/main.go status

migrate-create: ## 创建新的迁移文件
ifndef name
	$(error "name" parameter is required. Usage: make migrate-create name=20240101_create_table)
endif
	@echo "Creating migration: $(name)..."
	go run cmd/migrate/main.go create $(name)

docker-build: ## 构建 Docker 镜像
	@echo "Building Docker image..."
	docker build -t ddd-scaffold-backend:latest .
	@echo "Docker image built successfully!"

docker-run: ## 运行 Docker 容器
	@echo "Running Docker container..."
	docker-compose -f deployments/docker/docker-compose.yml up -d
	@echo "Docker containers started!"

docker-stop: ## 停止 Docker 容器
	@echo "Stopping Docker containers..."
	docker-compose -f deployments/docker/docker-compose.yml down
	@echo "Docker containers stopped!"

docker-logs: ## 查看 Docker 日志
	@echo "Viewing Docker logs..."
	docker-compose -f deployments/docker/docker-compose.yml logs -f backend

docker-clean: ## 清理 Docker 资源
	@echo "Cleaning Docker resources..."
	docker-compose -f deployments/docker/docker-compose.yml down -v
	docker system prune -f
	@echo "Docker resources cleaned!"

lint: ## 运行代码检查
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

tidy: ## 清理依赖
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

deps: ## 下载依赖
	@echo "Downloading dependencies..."
	go mod download

generate: ## 生成代码（Wire 等）
	@echo "Generating code..."
	wire gen ./internal/infrastructure/wire
	@echo "Code generated successfully!"

mock: ## 生成 Mock 文件
	@echo "Generating mock files..."
	mockgen -source=./internal/domain/user/repository/user_repository.go -destination=./internal/domain/user/repository/mock/user_repository_mock.go
	@echo "Mock files generated!"

profile: ## 性能分析
	@echo "Starting profiling..."
	go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

trace: ## 链路追踪
	@echo "Starting trace..."
	go tool trace http://localhost:8080/debug/trace

clean: ## 清理构建产物
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf docs/swagger.json docs/swagger.yaml
	@echo "Clean completed!"

all: lint test build swagger ## 执行完整构建流程
	@echo "Build completed successfully!"

# 生产环境构建
prod-build: ## 生产环境编译
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build \
		-ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)'" \
		-o bin/server cmd/server/main.go
	@echo "Production binary built: bin/server"

VERSION ?= 1.0.0

# 发布版本
release: prod-build ## 创建发布版本
	@echo "Creating release $(VERSION)..."
	tar -czf ddd-scaffold-$(VERSION)-linux-amd64.tar.gz -C bin server
	sha256sum ddd-scaffold-$(VERSION)-linux-amd64.tar.gz > ddd-scaffold-$(VERSION)-checksum.txt
	@echo "Release $(VERSION) created!"
