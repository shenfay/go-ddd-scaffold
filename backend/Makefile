# Go DDD Scaffold Makefile

.PHONY: help lint test build run clean swag

# 默认目标
help:
	@echo "Go DDD Scaffold 开发工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make lint     - 运行代码规范检查"
	@echo "  make test     - 运行单元测试"
	@echo "  make build    - 编译项目"
	@echo "  make run      - 运行服务器"
	@echo "  make clean    - 清理编译产物"

# 代码规范检查
lint:
	@echo "🔍 运行代码规范检查..."
	python3 ../tools/code-lint.py

# 单元测试
test:
	@echo "🧪 运行单元测试..."
	go test -v ./...

# 编译项目
build:
	@echo "🔨 编译项目..."
	go build -o server cmd/server/main.go

# 运行服务器
run: build
	@echo "🚀 启动服务器..."
	./server

# 清理编译产物
clean:
	@echo "🧹 清理编译产物..."
	rm -f server
	rm -f *.test

# 安装开发依赖
install-deps:
	@echo "📦 安装开发依赖..."
	go install golang.org/x/lint/golint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# 生成API文档
swag:
	@echo "📚 生成Swagger API文档..."
	@$$(go env GOPATH)/bin/swag init -g docs/main.go -o docs --parseInternal