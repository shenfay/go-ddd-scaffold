#!/bin/bash
# Grafana 批量导入仪表盘脚本
# 用法: ./scripts/import-all-dashboards.sh [options]
#
# Options:
#   --url URL          Grafana URL (默认: http://localhost:3000)
#   --api-key KEY      Grafana API Key
#   --dir DIR          仪表盘目录 (默认: grafana/dashboards)
#   --dry-run          只显示请求，不实际执行
#   --help             显示帮助信息

set -e

# 默认配置
GRAFANA_URL="http://localhost:3000"
DASHBOARD_DIR="grafana/dashboards"
API_KEY=""
DRY_RUN=false

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        --url)
            GRAFANA_URL="$2"
            shift 2
            ;;
        --api-key)
            API_KEY="$2"
            shift 2
            ;;
        --dir)
            DASHBOARD_DIR="$2"
            shift 2
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --help)
            head -11 "$0" | tail -10 | sed 's/^# //'
            exit 0
            ;;
        *)
            echo -e "${RED}❌ 未知参数: $1${NC}"
            exit 1
            ;;
    esac
done

# 检查 API Key
if [ -z "$API_KEY" ]; then
    echo -e "${YELLOW}⚠️  未提供 API Key，尝试从环境变量读取...${NC}"
    API_KEY=${GRAFANA_API_KEY:-""}
fi

if [ -z "$API_KEY" ]; then
    echo -e "${RED}❌ 错误: 请提供 Grafana API Key${NC}"
    echo -e "${YELLOW}💡 使用方式:${NC}"
    echo -e "   1. 命令行: --api-key YOUR_KEY"
    echo -e "   2. 环境变量: export GRAFANA_API_KEY=YOUR_KEY"
    echo -e "   3. 创建方式: Grafana > Configuration > API Keys > Add API Key (Admin 权限)"
    exit 1
fi

# 检查仪表盘目录
if [ ! -d "$DASHBOARD_DIR" ]; then
    echo -e "${RED}❌ 错误: 仪表盘目录不存在: $DASHBOARD_DIR${NC}"
    exit 1
fi

# 查找所有 JSON 文件
DASHBOARD_FILES=($(find "$DASHBOARD_DIR" -name "*.json" -type f | sort))

if [ ${#DASHBOARD_FILES[@]} -eq 0 ]; then
    echo -e "${RED}❌ 错误: 在 $DASHBOARD_DIR 中未找到 JSON 文件${NC}"
    exit 1
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 Grafana 批量导入仪表盘${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Grafana URL:${NC} $GRAFANA_URL"
echo -e "${GREEN}仪表盘目录:${NC} $DASHBOARD_DIR"
echo -e "${GREEN}找到文件数:${NC} ${#DASHBOARD_FILES[@]}"
echo ""

# 统计
SUCCESS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0

# 逐个导入
for DASHBOARD_FILE in "${DASHBOARD_FILES[@]}"; do
    FILENAME=$(basename "$DASHBOARD_FILE")
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "📄 处理: ${FILENAME}"
    
    # 验证 JSON 格式
    if ! python3 -c "import json; json.load(open('$DASHBOARD_FILE'))" 2>/dev/null; then
        echo -e "${RED}❌ JSON 格式错误，跳过${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
        continue
    fi
    
    # 提取仪表盘信息
    DASH_INFO=$(python3 -c "
import json
with open('$DASHBOARD_FILE') as f:
    data = json.load(f)
print(data.get('title', 'Unknown'))
print(data.get('uid', 'N/A'))
print(len(data.get('panels', [])))
" 2>/dev/null)
    
    DASH_TITLE=$(echo "$DASH_INFO" | sed -n '1p')
    DASH_UID=$(echo "$DASH_INFO" | sed -n '2p')
    PANEL_COUNT=$(echo "$DASH_INFO" | sed -n '3p')
    
    echo -e "  ${GREEN}标题:${NC} $DASH_TITLE"
    echo -e "  ${GREEN}UID:${NC} $DASH_UID"
    echo -e "  ${GREEN}面板数:${NC} $PANEL_COUNT"
    
    # Dry Run 模式
    if [ "$DRY_RUN" = true ]; then
        echo -e "  ${YELLOW}⏭️  Dry Run 模式，跳过导入${NC}"
        SKIP_COUNT=$((SKIP_COUNT + 1))
        continue
    fi
    
    # 准备导入数据（内联处理）
    TEMP_FILE=$(mktemp)
    cat "$DASHBOARD_FILE" | python3 -c "
import sys, json
dashboard = json.load(sys.stdin)
json.dump({'dashboard': dashboard, 'overwrite': True, 'message': 'Imported via script'}, sys.stdout, indent=2)
" > "$TEMP_FILE"
    
    # 导入仪表盘
    HTTP_RESPONSE=$(curl -s -w "\n%{http_code}" \
        -X POST "$GRAFANA_URL/api/dashboards/db" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_KEY" \
        -d @"$TEMP_FILE" \
    )
    
    HTTP_BODY=$(echo "$HTTP_RESPONSE" | sed '$d')
    HTTP_STATUS=$(echo "$HTTP_RESPONSE" | tail -n 1)
    
    rm -f "$TEMP_FILE"
    
    # 检查响应
    if [ "$HTTP_STATUS" = "200" ]; then
        ACTION=$(echo "$HTTP_BODY" | python3 -c 'import sys, json; data=json.load(sys.stdin); print(data.get("status", "success"))' 2>/dev/null || echo "success")
        echo -e "  ${GREEN}✅ 成功: $ACTION${NC}"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    elif [ "$HTTP_STATUS" = "401" ]; then
        echo -e "  ${RED}❌ 认证失败: API Key 无效或过期${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    elif [ "$HTTP_STATUS" = "403" ]; then
        echo -e "  ${RED}❌ 权限不足: API Key 需要 Admin 权限${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    else
        echo -e "  ${RED}❌ 失败 (HTTP $HTTP_STATUS)${NC}"
        echo -e "  ${YELLOW}响应: $HTTP_BODY${NC}"
        FAIL_COUNT=$((FAIL_COUNT + 1))
    fi
    
    echo ""
done

# 汇总
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 导入完成${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ 成功:${NC} $SUCCESS_COUNT"
if [ $FAIL_COUNT -gt 0 ]; then
    echo -e "${RED}❌ 失败:${NC} $FAIL_COUNT"
fi
if [ $SKIP_COUNT -gt 0 ]; then
    echo -e "${YELLOW}⏭️  跳过:${NC} $SKIP_COUNT"
fi
echo ""

if [ $FAIL_COUNT -eq 0 ] && [ "$DRY_RUN" = false ]; then
    echo -e "${GREEN}🎉 所有仪表盘导入成功！${NC}"
    echo -e "${YELLOW}💡 访问 Grafana 查看: $GRAFANA_URL/dashboards${NC}"
    exit 0
elif [ "$DRY_RUN" = true ]; then
    echo -e "${YELLOW}⚠️  Dry Run 完成，使用实际导入请移除 --dry-run 参数${NC}"
    exit 0
else
    echo -e "${RED}⚠️  部分仪表盘导入失败，请检查上方错误信息${NC}"
    exit 1
fi
