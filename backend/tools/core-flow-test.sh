#!/bin/bash

# 核心功能端到端测试脚本
# 测试注册、登录、获取个人信息三个核心流程
# 
# 用法:
#   ./core-flow-test.sh [选项]
#
# 选项:
#   -u, --username     用户名 (默认：testuser)
#   -e, --email        邮箱 (默认：test@example.com)
#   -p, --password     密码 (默认：password123)
#   -b, --base-url     API 基础 URL (默认：http://localhost:8080/api/v1)
#   -h, --help         显示帮助信息
#
# 示例:
#   ./core-flow-test.sh -u myuser -e myuser@test.com
#   ./core-flow-test.sh --username admin --base-url http://api.example.com/v1

set -e

# 默认值
USERNAME="testuser"
EMAIL="test@example.com"
PASSWORD="password123"
BASE_URL="http://localhost:8080/api/v1"
INTERACTIVE=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -u|--username)
      USERNAME="$2"
      shift 2
      ;;
    -e|--email)
      EMAIL="$2"
      shift 2
      ;;
    -p|--password)
      PASSWORD="$2"
      shift 2
      ;;
    -b|--base-url)
      BASE_URL="$2"
      shift 2
      ;;
    -h|--help)
      echo "用法：$0 [选项]"
      echo ""
      echo "选项:"
      echo "  -u, --username     用户名 (默认：testuser)"
      echo "  -e, --email        邮箱 (默认：test@example.com)"
      echo "  -p, --password     密码 (默认：password123)"
      echo "  -b, --base-url     API 基础 URL (默认：http://localhost:8080/api/v1)"
      echo "  -h, --help         显示帮助信息"
      echo ""
      echo "示例:"
      echo "  ./core-flow-test.sh                              # 使用默认值"
      echo "  ./core-flow-test.sh -u myuser -e myuser@test.com # 自定义用户名和邮箱"
      echo "  ./core-flow-test.sh --interactive                # 交互模式"
      exit 0
      ;;
    -i|--interactive)
      INTERACTIVE=true
      shift
      ;;
    *)
      echo "错误：未知选项 $1"
      echo "使用 -h 或 --help 查看帮助"
      exit 1
      ;;
  esac
done

# 交互模式函数
run_interactive() {
  echo ""
  echo "====================================="
  echo -e "${GREEN}🚀 进入交互模式${NC}"
  echo "====================================="
  echo ""
  
  read -p "请输入用户名 [默认：$USERNAME]: " input_username
  if [ -n "$input_username" ]; then
    USERNAME="$input_username"
  fi
  
  read -p "请输入邮箱 [默认：$EMAIL]: " input_email
  if [ -n "$input_email" ]; then
    EMAIL="$input_email"
  fi
  
  read -s -p "请输入密码 [默认：$PASSWORD]: " input_password
  echo ""
  if [ -n "$input_password" ]; then
    PASSWORD="$input_password"
  fi
  
  read -p "请输入 API 基础 URL [默认：$BASE_URL]: " input_base_url
  if [ -n "$input_base_url" ]; then
    BASE_URL="$input_base_url"
  fi
  
  echo ""
  echo -e "${GREEN}✅ 配置已更新:${NC}"
  echo "  用户名：$USERNAME"
  echo "  邮箱：$EMAIL"
  echo "  基础 URL: $BASE_URL"
  echo ""
}

# 如果没有提供任何参数，询问是否使用交互模式
if [ "$INTERACTIVE" = false ] && [ $# -eq 0 ]; then
  echo ""
  echo -e "${YELLOW}💡 提示：可以使用以下参数自定义测试:${NC}"
  echo "  -u, --username     用户名"
  echo "  -e, --email        邮箱"
  echo "  -p, --password     密码"
  echo "  -b, --base-url     API 基础 URL"
  echo "  -i, --interactive  交互模式"
  echo ""
  read -p "是否使用交互模式？(y/n) [默认：n]: " use_interactive
  echo ""
  if [ "$use_interactive" = "y" ] || [ "$use_interactive" = "Y" ]; then
    run_interactive
  fi
elif [ "$INTERACTIVE" = true ]; then
  run_interactive
fi

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_header() {
  echo ""
  echo "======================================"
  echo -e "${GREEN}DDD-Scaffold 核心功能测试${NC}"
  echo "======================================"
  echo -e "${YELLOW}配置:${NC}"
  echo "  用户名：$USERNAME"
  echo "  邮箱：$EMAIL"
  echo "  基础 URL: $BASE_URL"
  echo "======================================"
  echo ""
}

print_step() {
  echo -e "${GREEN}$1${NC}"
}

print_success() {
  echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
  echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
  echo -e "${RED}❌ $1${NC}"
}

# 检查 jq 是否安装
if ! command -v jq &> /dev/null; then
  print_error "jq 未安装，请先安装 jq"
  exit 1
fi

# 检查 curl 是否安装
if ! command -v curl &> /dev/null; then
  print_error "curl 未安装，请先安装 curl"
  exit 1
fi

print_header

ACCESS_TOKEN=""
REFRESH_TOKEN=""
USER_ID=""

# 1. 用户注册
print_step "📝 1. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"$USERNAME\",
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq .
echo ""

# 检查注册是否成功（通过 code 字段和用户 ID 判断）
REGISTER_CODE=$(echo "$REGISTER_RESPONSE" | jq -r '.code // empty')
REGISTER_USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user_id // empty')

if [ "$REGISTER_CODE" = "0" ] && [ -n "$REGISTER_USER_ID" ]; then
  print_success "注册成功，User ID: $REGISTER_USER_ID"
else
  # 注册可能失败，但继续尝试登录流程
  print_warning "注册响应异常，继续执行登录流程..."
fi
echo ""

# 2. 用户登录
print_step "🔐 2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"username_or_email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq .
echo ""

# 提取登录后的 token
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token // empty')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.refresh_token // empty')
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id // empty')

if [ -z "$ACCESS_TOKEN" ]; then
  print_error "登录失败，请检查用户名密码是否正确"
  exit 1
fi

print_success "登录成功"
echo "Access Token: ${ACCESS_TOKEN:0:20}..."
echo "Refresh Token: ${REFRESH_TOKEN:0:20}..."
echo "User ID: $USER_ID"
echo ""

# 3. 获取当前用户信息（使用 auth/me 端点）
print_step "👤 3. 测试获取当前用户信息..."
ME_RESPONSE=$(curl -s -X GET "${BASE_URL}/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "当前用户信息:"
echo "$ME_RESPONSE" | jq .
echo ""

# 4. 获取指定用户信息（使用 users/:id 端点）
print_step "📋 4. 测试获取指定用户信息..."

if [ -n "$USER_ID" ]; then
  USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/${USER_ID}" \
    -H "Authorization: Bearer $ACCESS_TOKEN")
  
  echo "用户详情:"
  echo "$USER_RESPONSE" | jq .
  echo ""
else
  print_warning "无法获取 User ID，跳过此测试"
fi
echo ""

# 5. 刷新 Token
print_step "🔄 5. 测试刷新 Token..."
REFRESH_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "刷新 Token 响应:"
echo "$REFRESH_RESPONSE" | jq .
echo ""

# 6. 登出
print_step "🚪 6. 测试用户登出..."
LOGOUT_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "登出响应:"
echo "$LOGOUT_RESPONSE" | jq .
echo ""

print_header
print_success "所有测试完成！"
print_header
