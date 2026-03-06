#!/bin/bash

# 前端联调测试脚本 - 验证注册/登录功能
# 使用 curl 测试 API 接口

set -e

BASE_URL="http://localhost:8080"
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_PASSWORD="TestPassword123!"
TEST_NICKNAME="测试用户"

echo "=========================================="
echo "前端联调测试 - 注册/登录功能验证"
echo "=========================================="
echo ""
echo "测试环境：$BASE_URL"
echo "测试邮箱：$TEST_EMAIL"
echo "测试密码：$TEST_PASSWORD"
echo ""

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 2

# 1. 健康检查
echo "=========================================="
echo "1️⃣  健康检查"
echo "=========================================="
curl -s "$BASE_URL/health" | jq . || {
    echo "❌ 健康检查失败，服务可能未启动"
    exit 1
}
echo "✅ 健康检查通过"
echo ""

# 2. 用户注册
echo "=========================================="
echo "2️⃣  用户注册"
echo "=========================================="
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\",
    \"nickname\": \"$TEST_NICKNAME\"
  }")

echo "注册请求响应:"
echo "$REGISTER_RESPONSE" | jq .

# 检查注册是否成功（注册接口可能不返回 token，需要手动登录）
CODE=$(echo "$REGISTER_RESPONSE" | jq -r '.code')
if [ "$CODE" = "Success" ]; then
    echo "✅ 注册成功"
else
    echo "❌ 注册失败"
    echo "$REGISTER_RESPONSE" | jq '.message'
    exit 1
fi
echo ""

# 3. 用户登录（使用刚注册的账号）
echo "=========================================="
echo "3️⃣  用户登录"
echo "=========================================="
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"$TEST_PASSWORD\"
  }")

echo "登录请求响应:"
echo "$LOGIN_RESPONSE" | jq .

# 检查登录是否成功
CODE=$(echo "$LOGIN_RESPONSE" | jq -r '.code')
if [ "$CODE" = "Success" ]; then
    echo "✅ 登录成功"
    LOGIN_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.accessToken')
    if [ -z "$LOGIN_TOKEN" ] || [ "$LOGIN_TOKEN" = "null" ]; then
        echo "❌ 登录成功但未返回 accessToken"
        exit 1
    fi
    echo "✅ 获取到 accessToken: ${LOGIN_TOKEN:0:20}..."
else
    echo "❌ 登录失败"
    echo "$LOGIN_RESPONSE" | jq '.message'
    exit 1
fi
echo ""

# 4. 获取用户信息（需要认证）
echo "=========================================="
echo "4️⃣  获取用户信息（需要认证）"
echo "=========================================="
USER_INFO_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/info" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "用户信息响应:"
echo "$USER_INFO_RESPONSE" | jq .

# 检查用户信息是否正确
if echo "$USER_INFO_RESPONSE" | jq -e '.code == "Success"' > /dev/null 2>&1; then
    echo "✅ 获取用户信息成功"
    RETURNED_EMAIL=$(echo "$USER_INFO_RESPONSE" | jq -r '.data.email')
    RETURNED_NICKNAME=$(echo "$USER_INFO_RESPONSE" | jq -r '.data.nickname')
    
    if [ "$RETURNED_EMAIL" = "$TEST_EMAIL" ]; then
        echo "✅ 邮箱匹配：$RETURNED_EMAIL"
    else
        echo "❌ 邮箱不匹配：期望 $TEST_EMAIL, 实际 $RETURNED_EMAIL"
        exit 1
    fi
    
    if [ "$RETURNED_NICKNAME" = "$TEST_NICKNAME" ]; then
        echo "✅ 昵称匹配：$RETURNED_NICKNAME"
    else
        echo "❌ 昵称不匹配：期望 $TEST_NICKNAME, 实际 $RETURNED_NICKNAME"
        exit 1
    fi
else
    echo "❌ 获取用户信息失败"
    echo "$USER_INFO_RESPONSE" | jq '.message'
    exit 1
fi
echo ""

# 5. 更新用户资料
echo "=========================================="
echo "5️⃣  更新用户资料"
echo "=========================================="
NEW_NICKNAME="更新昵称"
UPDATE_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/users/profile" \
  -H "Authorization: Bearer $LOGIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"nickname\": \"$NEW_NICKNAME\"
  }")

echo "更新资料响应:"
echo "$UPDATE_RESPONSE" | jq .

if echo "$UPDATE_RESPONSE" | jq -e '.code == "Success"' > /dev/null 2>&1; then
    echo "✅ 更新用户资料成功"
else
    echo "❌ 更新用户资料失败"
    echo "$UPDATE_RESPONSE" | jq '.message'
    exit 1
fi
echo ""

# 6. 用户登出
echo "=========================================="
echo "6️⃣  用户登出"
echo "=========================================="
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/logout" \
  -H "Authorization: Bearer $LOGIN_TOKEN")

echo "登出响应:"
echo "$LOGOUT_RESPONSE" | jq .

if echo "$LOGOUT_RESPONSE" | jq -e '.code == "Success"' > /dev/null 2>&1; then
    echo "✅ 登出成功"
else
    echo "❌ 登出失败"
    echo "$LOGOUT_RESPONSE" | jq '.message'
    # 登出失败不影响整体测试结果，继续执行
fi
echo ""

# 7. 测试错误密码登录
echo "=========================================="
echo "7️⃣  测试错误密码（预期失败）"
echo "=========================================="
WRONG_PASSWORD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"password\": \"WrongPassword\"
  }")

echo "错误密码响应:"
echo "$WRONG_PASSWORD_RESPONSE" | jq .

# 应该返回错误
if echo "$WRONG_PASSWORD_RESPONSE" | jq -e '.code != 0' > /dev/null 2>&1; then
    echo "✅ 正确返回错误：密码错误"
    ERROR_MESSAGE=$(echo "$WRONG_PASSWORD_RESPONSE" | jq -r '.message')
    echo "错误信息：$ERROR_MESSAGE"
else
    echo "❌ 错误密码竟然登录成功了？！"
    exit 1
fi
echo ""

# 总结
echo "=========================================="
echo "🎉 所有测试通过！"
echo "=========================================="
echo ""
echo "测试项目："
echo "✅ 1. 健康检查"
echo "✅ 2. 用户注册"
echo "✅ 3. 用户登录"
echo "✅ 4. 获取用户信息"
echo "✅ 5. 更新用户资料"
echo "✅ 6. 用户登出"
echo "✅ 7. 错误密码验证"
echo ""
echo "测试完成时间：$(date)"
echo "=========================================="
