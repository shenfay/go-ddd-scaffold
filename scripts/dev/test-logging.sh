#!/bin/bash

# 事件驱动日志系统 - 快速测试脚本
# 用法：./scripts/dev/test-logging.sh

set -e

echo "======================================"
echo "🧪 测试事件驱动日志系统"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API 地址
API_URL="${API_URL:-http://localhost:8080}"

echo -e "${BLUE}📍 API 地址：$API_URL${NC}"
echo ""

# 测试登录接口
echo "📝 步骤 1: 测试登录接口"
echo "----------------------------------------"

LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}')

echo "$LOGIN_RESPONSE" | jq '.' 2>/dev/null || echo "$LOGIN_RESPONSE"

# 检查是否成功
if echo "$LOGIN_RESPONSE" | grep -q "access_token"; then
    echo -e "${GREEN}✅ 登录成功！${NC}"
    
    # 提取 token
    ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token' 2>/dev/null || echo "")
    
    if [ -n "$ACCESS_TOKEN" ]; then
        echo ""
        echo "🎫 Access Token: ${ACCESS_TOKEN:0:50}..."
        echo ""
    fi
else
    echo -e "${RED}❌ 登录失败，请检查用户名密码或 API 服务状态${NC}"
    echo ""
    echo "提示："
    echo "  1. 确保 API 服务已启动：cd backend/cmd/api && go run main.go"
    echo "  2. 确保数据库迁移已执行：./scripts/dev/migrate-logging.sh"
    exit 1
fi

echo ""
echo "======================================"
echo "📊 步骤 2: 检查审计日志（需要数据库权限）"
echo "======================================"
echo ""

# 检查 DATABASE_URL
if [ -z "$DATABASE_URL" ]; then
    echo -e "${YELLOW}⚠️  DATABASE_URL 未设置，跳过数据库检查${NC}"
    echo ""
else
    # 查询最新的审计日志
    AUDIT_LOGS=$(psql "$DATABASE_URL" -t -c "SELECT action, status, created_at FROM audit_logs ORDER BY created_at DESC LIMIT 5;" 2>/dev/null)
    
    if [ $? -eq 0 ] && [ -n "$AUDIT_LOGS" ]; then
        echo -e "${GREEN}✅ 审计日志记录成功：${NC}"
        echo "$AUDIT_LOGS"
    else
        echo -e "${YELLOW}⚠️  无法查询审计日志（可能是权限问题或暂无数据）${NC}"
    fi
fi

echo ""
echo "======================================"
echo "📝 下一步操作"
echo "======================================"
echo ""
echo "1. 手动验证审计日志表:"
echo "   psql \$DATABASE_URL -c \"SELECT * FROM audit_logs ORDER BY created_at DESC LIMIT 10;\""
echo ""
echo "2. 查看 Worker 日志确认任务处理:"
echo "   tail -f backend/logs/worker.log"
echo ""
echo "3. 性能压测（可选）:"
echo "   wrk -t12 -c400 -d30s $API_URL/api/v1/auth/login \\"
echo "     -H 'Content-Type: application/json' \\"
echo "     -d '{\"email\":\"test@example.com\",\"password\":\"Test123456\"}'"
echo ""
echo -e "${GREEN}✓ 基本功能测试完成！${NC}"
