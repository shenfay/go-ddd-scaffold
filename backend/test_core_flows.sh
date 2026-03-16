#!/bin/bash

# 核心功能测试脚本
# 测试注册、登录、获取个人信息三个核心流程

BASE_URL="http://localhost:8080/api/v1"

echo "======================================"
echo "DDD-Scaffold 核心功能测试"
echo "======================================"
echo ""

# 1. 用户注册
echo "📝 1. 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }')

echo "注册响应:"
echo "$REGISTER_RESPONSE" | jq .
echo ""

# 提取 token（用于后续请求）
ACCESS_TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.data.token // empty')

if [ -z "$ACCESS_TOKEN" ]; then
  echo "⚠️  注册可能失败或没有返回 token，尝试使用登录获取 token..."
else
  echo "✅ 注册成功，Access Token: ${ACCESS_TOKEN:0:20}..."
fi
echo ""

# 2. 用户登录
echo "🔐 2. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }')

echo "登录响应:"
echo "$LOGIN_RESPONSE" | jq .
echo ""

# 提取登录后的 token
ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // empty')
REFRESH_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.refresh_token // empty')

if [ -z "$ACCESS_TOKEN" ]; then
  echo "❌ 登录失败，请检查用户名密码是否正确"
  exit 1
fi

echo "✅ 登录成功"
echo "Access Token: ${ACCESS_TOKEN:0:20}..."
echo "Refresh Token: ${REFRESH_TOKEN:0:20}..."
echo ""

# 3. 获取当前用户信息（使用 auth/me 端点）
echo "👤 3. 测试获取当前用户信息..."
ME_RESPONSE=$(curl -s -X GET "${BASE_URL}/auth/me" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "当前用户信息:"
echo "$ME_RESPONSE" | jq .
echo ""

# 4. 获取指定用户信息（使用 users/:id 端点）
echo "📋 4. 测试获取指定用户信息..."
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user_id // empty')

if [ -n "$USER_ID" ]; then
  USER_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/${USER_ID}" \
    -H "Authorization: Bearer $ACCESS_TOKEN")
  
  echo "用户详情:"
  echo "$USER_RESPONSE" | jq .
  echo ""
else
  echo "⚠️  无法获取 User ID，跳过此测试"
fi
echo ""

# 5. 刷新 Token
echo "🔄 5. 测试刷新 Token..."
REFRESH_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "刷新 Token 响应:"
echo "$REFRESH_RESPONSE" | jq .
echo ""

# 6. 登出
echo "🚪 6. 测试用户登出..."
LOGOUT_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "登出响应:"
echo "$LOGOUT_RESPONSE" | jq .
echo ""

echo "======================================"
echo "✅ 所有测试完成！"
echo "======================================"
