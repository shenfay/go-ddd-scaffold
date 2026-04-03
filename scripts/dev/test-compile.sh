#!/bin/bash

# 测试重构后的编译状态

set -e

echo "======================================"
echo "🔍 测试重构后编译状态"
echo "======================================"
echo ""

cd /Users/shenfay/Projects/ddd-scaffold/backend

echo "📦 步骤 1: 清理模块缓存..."
go clean -modcache 2>/dev/null || true
echo "✅ 完成"
echo ""

echo "📦 步骤 2: 整理依赖..."
go mod tidy 2>&1 | head -20 || echo "⚠️  go mod tidy 有警告（可能是正常的）"
echo ""

echo "📦 步骤 3: 编译 API 服务..."
if go build -o /tmp/api_test ./cmd/api 2>&1; then
    echo "✅ API 编译成功"
else
    echo "❌ API 编译失败，错误信息："
    go build ./cmd/api 2>&1 | head -30
    exit 1
fi
echo ""

echo "📦 步骤 4: 编译 Worker 服务..."
if go build -o /tmp/worker_test ./cmd/worker 2>&1; then
    echo "✅ Worker 编译成功"
else
    echo "❌ Worker 编译失败，错误信息："
    go build ./cmd/worker 2>&1 | head -30
    exit 1
fi
echo ""

echo "======================================"
echo "✅ 所有编译测试通过！"
echo "======================================"
echo ""
echo "下一步操作："
echo "1. 启动 API 服务：cd cmd/api && go run main.go"
echo "2. 启动 Worker 服务：cd cmd/worker && go run main.go"
echo "3. 测试登录功能"
echo ""
