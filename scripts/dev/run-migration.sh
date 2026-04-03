#!/bin/bash

# 自动执行数据库迁移（无需确认）

set -e

echo "======================================"
echo "🚀 自动执行数据库迁移"
echo "======================================"
echo ""

# 设置默认 DATABASE_URL（如果未设置）
export DATABASE_URL="${DATABASE_URL:-postgres://localhost:5432/go_ddd_scaffold?sslmode=disable}"

echo "📊 数据库连接信息:"
echo "   DATABASE_URL: $DATABASE_URL"
echo ""

# 检查 psql 是否可用
if ! command -v psql &> /dev/null; then
    echo "❌ 错误：psql 未安装，请确保 PostgreSQL 已正确安装"
    exit 1
fi

# 进入 migrations 目录
cd "$(dirname "$0")/../../backend/migrations"

echo "📁 当前目录：$(pwd)"
echo ""

# 显示待执行的迁移文件
echo "📋 待执行的迁移文件:"
ls -lh 005_*.sql 006_*.sql 2>/dev/null || {
    echo "❌ 未找到迁移文件"
    exit 1
}
echo ""

# 测试数据库连接
echo "🔍 测试数据库连接..."
if ! psql "$DATABASE_URL" -c "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ 无法连接到数据库，请检查 DATABASE_URL 配置"
    exit 1
fi
echo "✅ 数据库连接成功"
echo ""

# 执行迁移
echo "======================================"
echo "🔄 开始执行向上迁移"
echo "======================================"
echo ""

# 执行 005 审计日志表迁移
echo "📝 [1/2] 创建 audit_logs 表..."
if psql "$DATABASE_URL" < 005_create_audit_logs_table.up.sql > /dev/null 2>&1; then
    echo "✅ audit_logs 表创建成功"
else
    echo "⚠️  audit_logs 表可能已存在或创建失败"
fi
echo ""

# 执行 006 活动日志表迁移
echo "📝 [2/2] 创建 activity_logs 表..."
if psql "$DATABASE_URL" < 006_create_activity_logs_table.up.sql > /dev/null 2>&1; then
    echo "✅ activity_logs 表创建成功"
else
    echo "⚠️  activity_logs 表可能已存在或创建失败"
fi
echo ""

# 验证迁移结果
echo "======================================"
echo "📊 验证迁移结果"
echo "======================================"
echo ""

# 检查表是否存在
echo "检查 audit_logs 表..."
TABLE_EXISTS=$(psql "$DATABASE_URL" -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'audit_logs');" 2>/dev/null | xargs)
if [ "$TABLE_EXISTS" = "t" ]; then
    echo "✅ audit_logs 表存在"
else
    echo "❌ audit_logs 表不存在"
fi

echo "检查 activity_logs 表..."
TABLE_EXISTS=$(psql "$DATABASE_URL" -t -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'activity_logs');" 2>/dev/null | xargs)
if [ "$TABLE_EXISTS" = "t" ]; then
    echo "✅ activity_logs 表存在"
else
    echo "❌ activity_logs 表不存在"
fi

echo ""

# 显示表结构
echo "======================================"
echo "📋 表结构概览"
echo "======================================"
echo ""

echo "audit_logs 表字段:"
psql "$DATABASE_URL" -c "\d+ audit_logs" 2>/dev/null | head -20 || echo "无法查看表结构（可能是权限问题）"
echo ""

echo "activity_logs 表字段:"
psql "$DATABASE_URL" -c "\d+ activity_logs" 2>/dev/null | head -15 || echo "无法查看表结构（可能是权限问题）"
echo ""

echo "======================================"
echo "✅ 数据库迁移完成！"
echo "======================================"
echo ""
echo "下一步操作："
echo "1. 启动 API 服务：cd backend/cmd/api && go run main.go"
echo "2. 启动 Worker 服务：cd backend/cmd/worker && go run main.go"
echo "3. 测试登录接口：curl -X POST http://localhost:8080/api/v1/auth/login ..."
echo "4. 检查审计日志：psql \$DATABASE_URL -c \"SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;\""
echo ""
