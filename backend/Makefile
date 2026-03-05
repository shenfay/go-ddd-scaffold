# Makefile for DDD-Scaffold Backend

.PHONY: help build test run dev swagger migrate clean deps lint fmt tidy generate mock profile trace all prod-build release install-tools dao-generate dao-clean dao-install-deps

# 获取当前时间
TIMESTAMP := $(shell date +%Y%m%d_%H%M%S)
VERSION ?= 1.0.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

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
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	@echo "Binary built: bin/server"

build-local: ## 本地编译（当前平台）
	@echo "Building for local platform..."
	mkdir -p bin
	go build -ldflags="-s -w" -o bin/server cmd/server/main.go
	@echo "Local binary built: bin/server"

test: ## 运行所有测试
	@echo "Running tests..."
	go test -v ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

cover: ## 查看覆盖率报告
	@echo "Opening coverage report..."
	go tool cover -html=coverage.out

bench: ## 运行基准测试
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

race: ## 运行竞态检测
	@echo "Running race detector..."
	go test -race ./...

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
	SWAG_BIN=$(shell go env GOPATH)/bin/swag; \
	if [ ! -f "$$SWAG_BIN" ]; then \
		echo "swag not found, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		SWAG_BIN=$(shell go env GOPATH)/bin/swag; \
	fi; \
	$$SWAG_BIN init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal
	@echo "Swagger docs generated successfully!"
	@echo "Access at: http://localhost:8080/swagger/index.html"

swagger-validate: ## 验证 Swagger 文档
	@echo "Validating Swagger documentation..."
	@if ! command -v swag &> /dev/null; then \
		echo "swag not found, installing..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	fi
	swag fmt
	@echo "Swagger validation completed!"

swagger-serve: ## 启动 Swagger UI 服务器
	@echo "Starting Swagger UI server..."
	swag serve --dir ./docs --port 8081

migrate-up: ## 执行数据库迁移（向上）
	@echo "Running database migrations (up)..."
	@if [ -f "cmd/migrate/main.go" ]; then \
		go run cmd/migrate/main.go up; \
	else \
		echo "Migration command not found. Using goose directly..."; \
		goose -dir ./migrations/sql postgres "host=localhost user=shenfay password=postgres dbname=go_ddd_scaffold sslmode=disable" up; \
	fi

migrate-down: ## 回滚数据库迁移（向下）
	@echo "Rolling back database migrations (down)..."
	@if [ -f "cmd/migrate/main.go" ]; then \
		go run cmd/migrate/main.go down; \
	else \
		echo "Migration command not found. Using goose directly..."; \
		goose -dir ./migrations/sql postgres "host=localhost user=shenfay password=postgres dbname=go_ddd_scaffold sslmode=disable" down; \
	fi

migrate-status: ## 查看迁移状态
	@echo "Checking migration status..."
	@if [ -f "cmd/migrate/main.go" ]; then \
		go run cmd/migrate/main.go status; \
	else \
		echo "Migration command not found. Using goose directly..."; \
		goose -dir ./migrations/sql postgres "host=localhost user=shenfay password=postgres dbname=go_ddd_scaffold sslmode=disable" status; \
	fi

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
	LINTER_BIN=$(shell go env GOPATH)/bin/golangci-lint; \
	if [ ! -f "$$LINTER_BIN" ]; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		LINTER_BIN=$(shell go env GOPATH)/bin/golangci-lint; \
	fi; \
	$$LINTER_BIN run ./...

lint-fix: ## 运行代码检查并自动修复
	@echo "Running linter with auto-fix..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run --fix ./...

fmt: ## 格式化代码
	@echo "Formatting code..."
	go fmt ./...
	@if ! command -v goimports &> /dev/null; then \
		echo "goimports not found, installing..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	$(shell go env GOPATH)/bin/goimports -w .

vet: ## 运行 go vet
	@echo "Running go vet..."
	go vet ./...

tidy: ## 清理依赖
	@echo "Tidying dependencies..."
	go mod tidy
	go mod verify

deps: ## 下载依赖
	@echo "Downloading dependencies..."
	go mod download

generate: ## 生成代码（Wire 等）
	@echo "Generating code..."
	WIRE_BIN=$(shell go env GOPATH)/bin/wire; \
	if [ ! -f "$$WIRE_BIN" ]; then \
		echo "wire not found, installing..."; \
		go install github.com/google/wire/cmd/wire@latest; \
		WIRE_BIN=$(shell go env GOPATH)/bin/wire; \
	fi; \
	$$WIRE_BIN gen ./internal/infrastructure/wire
	@echo "Code generated successfully!"

go-generate: ## 运行 go generate
	@echo "Running go generate..."
	go generate ./...

dao-generate: ## 生成 DAO 层代码
	@echo "Generating DAO layer code..."
	@cd tools/dao-gen && go run main.go
	@echo "DAO generation completed!"

dao-clean: ## 清理生成的 DAO 文件
	@echo "Cleaning generated DAO files..."
	rm -rf internal/infrastructure/persistence/gorm/dao
	rm -rf internal/infrastructure/persistence/gorm/model
	@echo "DAO files cleaned!"

dao-install-deps: ## 安装 DAO 生成工具依赖
	@echo "Installing DAO generator dependencies..."
	@cd tools/dao-gen && go mod tidy
	@echo "DAO dependencies installed!"

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
	rm -f *.log
	@echo "Clean completed!"

clean-all: ## 清理所有生成文件
	@echo "Cleaning all generated files..."
	make clean
	go clean -cache -modcache -i -r
	@echo "All generated files cleaned!"

all: deps fmt vet lint test build swagger dao-generate ## 执行完整构建流程
	@echo "Build completed successfully!"

check: fmt vet lint ## 代码质量检查
	@echo "Code quality check completed!"

validate: check test ## 完整验证
	@echo "Validation completed successfully!"

# 生产环境构建
prod-build: ## 生产环境编译
	@echo "Building for production..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build \
		-ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" \
		-o bin/server cmd/server/main.go
	@echo "Production binary built: bin/server"

prod-build-all: ## 构建多个平台的生产包
	@echo "Building for multiple platforms..."
	mkdir -p bin
	# Linux AMD64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/server-linux-amd64 cmd/server/main.go
	# Linux ARM64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/server-linux-arm64 cmd/server/main.go
	# Darwin AMD64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/server-darwin-amd64 cmd/server/main.go
	# Darwin ARM64
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'" -o bin/server-darwin-arm64 cmd/server/main.go
	@echo "Multi-platform binaries built successfully!"

# 发布版本
release: prod-build-all ## 创建发布版本
	@echo "Creating release $(VERSION)..."
	mkdir -p releases/$(VERSION)
	# 打包各平台二进制文件
	tar -czf releases/$(VERSION)/ddd-scaffold-$(VERSION)-linux-amd64.tar.gz -C bin server-linux-amd64
	tar -czf releases/$(VERSION)/ddd-scaffold-$(VERSION)-linux-arm64.tar.gz -C bin server-linux-arm64
	tar -czf releases/$(VERSION)/ddd-scaffold-$(VERSION)-darwin-amd64.tar.gz -C bin server-darwin-amd64
	tar -czf releases/$(VERSION)/ddd-scaffold-$(VERSION)-darwin-arm64.tar.gz -C bin server-darwin-arm64
	# 生成校验和
	cd releases/$(VERSION) && sha256sum *.tar.gz > checksums.txt
	@echo "Release $(VERSION) created in releases/$(VERSION)/"

install-tools: ## 安装开发工具
	@echo "Installing development tools..."
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing wire..."
	go install github.com/google/wire/cmd/wire@latest
	@echo "Installing goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Installing mockgen..."
	go install go.uber.org/mock/mockgen@latest
	@echo "Verifying installations..."
	@which $(shell go env GOPATH)/bin/swag > /dev/null && echo "✓ swag installed" || echo "✗ swag installation failed"
	@which $(shell go env GOPATH)/bin/golangci-lint > /dev/null && echo "✓ golangci-lint installed" || echo "✗ golangci-lint installation failed"
	@which $(shell go env GOPATH)/bin/wire > /dev/null && echo "✓ wire installed" || echo "✗ wire installation failed"
	@which $(shell go env GOPATH)/bin/goimports > /dev/null && echo "✓ goimports installed" || echo "✗ goimports installation failed"
	@which $(shell go env GOPATH)/bin/mockgen > /dev/null && echo "✓ mockgen installed" || echo "✗ mockgen installation failed"
	@echo "Development tools installation completed!"

info: ## 显示项目信息
	@echo "=== Project Info ==="
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell go version)"
	@echo "Go Mod Path: $(shell go list -m)"
	@echo "==================="
