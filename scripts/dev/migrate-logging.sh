#!/bin/bash

# 事件驱动日志系统 - 数据库迁移执行脚本
# 用法：./scripts/dev/migrate-logging.sh

set -e

echo "======================================"
echo "🚀 执行日志系统数据库迁移"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查环境变量
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}⚠️  警告：DATABASE_URL 未设置，使用默认值${NC}"
    export DATABASE_URL="postgres://localhost:5432/go_ddd_scaffold?sslmode=disable"
fi

echo "📊 数据库连接信息:"
echo "   DATABASE_URL: $DATABASE_URL"
echo ""

# 检查 migrate 工具
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}❌ 错误：migrate 工具未安装${NC}"
    echo ""
    echo "请安装 migrate 工具："
    echo "  macOS: brew install golang-migrate"
    echo "  Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.$(uname -s)-$(uname -m).tar.gz | tar xvz"
    echo "  Go:    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

echo "✅ migrate 工具已安装"
echo ""

# 进入 migrations 目录
cd "$(dirname "$0")/../../backend/migrations"

echo "📁 当前目录：$(pwd)"
echo ""

# 显示待执行的迁移文件
echo "📋 待执行的迁移文件:"
ls -lh 005_*.sql 006_*.sql 2>/dev/null || echo "  未找到迁移文件"
echo ""

# 询问是否继续
read -p "是否继续执行迁移？(y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}⚠️  迁移已取消${NC}"
    exit 0
fi

echo ""
echo "======================================"
echo "🔄 开始执行向上迁移"
echo "======================================"
echo ""

# 执行向上迁移
migrate -path . -database "$DATABASE_URL" -verbose up

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ 迁移成功完成！${NC}"
else
    echo ""
    echo -e "${RED}❌ 迁移执行失败${NC}"
    exit 1
fi

echo ""
echo "======================================"
echo "📊 验证迁移结果"
echo "======================================"
echo ""

# 验证表是否存在
psql "$DATABASE_URL" -c "\dt audit_logs" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ audit_logs 表创建成功${NC}"
else
    echo -e "${YELLOW}⚠️  无法验证 audit_logs 表（可能是权限问题）${NC}"
fi

psql "$DATABASE_URL" -c "\dt activity_logs" 2>/dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ activity_logs 表创建成功${NC}"
else
    echo -e "${YELLOW}⚠️  无法验证 activity_logs 表（可能是权限问题）${NC}"
fi

echo ""
echo "======================================"
echo "📝 下一步操作"
echo "======================================"
echo ""
echo "1. 启动 API 服务进行测试:"
echo "   cd backend/cmd/api && go run main.go"
echo ""
echo "2. 启动 Worker 服务:"
echo "   cd backend/cmd/worker && go run main.go"
echo ""
echo "3. 测试登录接口:"
echo "   curl -X POST http://localhost:8080/api/v1/auth/login \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"email\":\"test@example.com\",\"password\":\"password123\"}'"
echo ""
echo "4. 检查审计日志:"
echo "   psql \$DATABASE_URL -c \"SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;\""
echo ""
echo -e "${GREEN}✓ 所有步骤已完成！${NC}"
