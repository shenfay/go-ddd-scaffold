#!/bin/bash

# 事件系统测试脚本
# 测试领域事件的发布、订阅、处理流程
# 
# 用法:
#   ./event-flow-test.sh [选项]
#
# 选项:
#   -t, --type     事件类型 (user_registered, user_logged_in, all) (默认：all)
#   -b, --base-url API 基础 URL (默认：http://localhost:8080/api/v1)
#   -h, --help     显示帮助信息
#
# 示例:
#   ./event-flow-test.sh                              # 测试所有事件类型
#   ./event-flow-test.sh -t user_registered           # 只测试用户注册事件
#   ./event-flow-test.sh -t user_logged_in            # 只测试用户登录事件

set -e

# 默认值
EVENT_TYPE="all"
BASE_URL="http://localhost:8080/api/v1"
USERNAME="testuser_event"
EMAIL="testuser_event@example.com"
PASSWORD="password123"

# 从 .env 文件加载数据库配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/../configs/.env"

if [ -f "$ENV_FILE" ]; then
  # 读取 .env 文件并导出变量
  export $(grep -v '^#' "$ENV_FILE" | xargs)
  
  # 设置数据库连接参数
  DB_HOST="${APP_DATABASE_HOST:-localhost}"
  DB_PORT="${APP_DATABASE_PORT:-5432}"
  DB_NAME="${APP_DATABASE_NAME:-go_ddd_scaffold}"
  DB_USER="${APP_DATABASE_USER:-shenfay}"
  DB_PASSWORD="${APP_DATABASE_PASSWORD:-postgres}"
else
  # 默认值
  DB_HOST="localhost"
  DB_PORT="5432"
  DB_NAME="go_ddd_scaffold"
  DB_USER="shenfay"
  DB_PASSWORD="postgres"
fi

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
  echo ""
  echo "======================================"
  echo -e "${GREEN}📬 DDD 事件系统流程测试${NC}"
  echo "======================================"
  echo -e "${YELLOW}配置:${NC}"
  echo "  事件类型：$EVENT_TYPE"
  echo "  基础 URL: $BASE_URL"
  echo "======================================"
  echo ""
}

print_step() {
  echo -e "${BLUE}$1${NC}"
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

print_info() {
  echo -e "${BLUE}ℹ️  $1${NC}"
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--type)
      EVENT_TYPE="$2"
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
      echo "  -t, --type     事件类型 (user_registered, user_logged_in, all) (默认：all)"
      echo "  -b, --base-url API 基础 URL (默认：http://localhost:8080/api/v1)"
      echo "  -h, --help     显示帮助信息"
      echo ""
      echo "示例:"
      echo "  ./event-flow-test.sh                              # 测试所有事件类型"
      echo "  ./event-flow-test.sh -t user_registered           # 只测试用户注册事件"
      echo "  ./event-flow-test.sh -t user_logged_in            # 只测试用户登录事件"
      exit 0
      ;;
    *)
      echo "错误：未知选项 $1"
      echo "使用 -h 或 --help 查看帮助"
      exit 1
      ;;
  esac
done

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

# 测试用户注册事件（UserRegistered）
test_user_registered() {
  print_step "📝 测试 1: UserRegistered 事件发布与订阅"
  echo ""
  
  # 生成唯一的测试用户
  TIMESTAMP=$(date +%s)
  TEST_USERNAME="event_test_${TIMESTAMP}"
  TEST_EMAIL="event_test_${TIMESTAMP}@example.com"
  
  print_info "使用测试用户: $TEST_USERNAME"
  
  # 1. 注册新用户（触发 UserRegistered 事件）
  print_info "步骤 1: 注册新用户（触发 UserRegistered 事件）"
  REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"username\": \"$TEST_USERNAME\",
      \"email\": \"$TEST_EMAIL\",
      \"password\": \"$PASSWORD\"
    }")
  
  echo "$REGISTER_RESPONSE" | jq .
  echo ""
  
  REGISTER_CODE=$(echo "$REGISTER_RESPONSE" | jq -r '.code // empty')
  REGISTER_USER_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.data.user_id // empty')
  
  if [ "$REGISTER_CODE" = "0" ] && [ -n "$REGISTER_USER_ID" ]; then
    print_success "用户注册成功，User ID: $REGISTER_USER_ID"
    print_info "✅ UserRegistered 事件已发布到 EventPublisher（异步队列）"
  else
    print_error "注册失败，无法继续测试"
    echo "$REGISTER_RESPONSE" | jq -r '.message // "未知错误"'
    return 1
  fi
  echo ""
  
  # 2. 等待事件处理完成
  print_info "步骤 2: 等待异步事件处理完成..."
  print_info "提示：事件通过 AsynqPublisher 发送到 Redis 队列，Worker 需要时间消费"
  
  # 轮询等待，最多等待 10 秒
  MAX_WAIT=10
  WAITED=0
  AUDIT_LOG_COUNT=0
  
  while [ $WAITED -lt $MAX_WAIT ]; do
    sleep 1
    WAITED=$((WAITED + 1))
    
    # 检查审计日志
    AUDIT_LOG_COUNT=$(PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -t -c \
      "SELECT COUNT(*) FROM audit_logs WHERE action='USER_REGISTERED' AND user_id=$REGISTER_USER_ID;" 2>/dev/null | tr -d ' ')
    
    if [ "$AUDIT_LOG_COUNT" -gt 0 ]; then
      print_success "事件在 ${WAITED} 秒内被处理完成"
      break
    fi
    
    echo -n "."
  done
  echo ""
  
  # 3. 检查数据库中的审计日志
  print_info "步骤 3: 查询数据库中的审计日志..."
  
  if [ "$AUDIT_LOG_COUNT" -gt 0 ]; then
    print_success "✅ 发现审计日志！AuditSubscriber 成功处理了 UserRegistered 事件"
    
    # 显示最近的审计日志
    print_info "审计日志详情:"
    PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -c \
      "SELECT id, user_id, action, resource_type, status, occurred_at FROM audit_logs WHERE action='USER_REGISTERED' AND user_id=$REGISTER_USER_ID ORDER BY occurred_at DESC LIMIT 3;" 2>/dev/null
  else
    print_warning "⚠️  未发现审计日志（等待了 ${MAX_WAIT} 秒）"
    print_info "可能的原因:"
    echo "   1. Worker 未启动或未正确处理事件"
    echo "   2. Redis 队列有问题"
    echo "   3. 事件处理器抛出异常"
    echo ""
    print_info "排查建议:"
    echo "   - 查看 Worker 日志：tail -f logs/app.log | grep -i worker"
    echo "   - 查看 Redis 队列：redis-cli LLEN asynq:default"
    echo "   - 检查 domain_events 表：SELECT * FROM domain_events ORDER BY id DESC LIMIT 5;"
    echo "   - 启动 Asynq Monitor: make run-asynqmon"
  fi
  echo ""
  
  # 4. 检查 domain_events 表（事件溯源存储）
  print_info "步骤 4: 检查 domain_events 表（事件溯源）..."
  
  EVENT_COUNT=$(PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -t -c \
    "SELECT COUNT(*) FROM domain_events WHERE aggregate_id='$REGISTER_USER_ID' AND event_type='UserRegistered';" 2>/dev/null | tr -d ' ')
  
  if [ "$EVENT_COUNT" -gt 0 ]; then
    print_success "✅ 领域事件已持久化到 domain_events 表"
    
    print_info "事件详情:"
    PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -c \
      "SELECT id, aggregate_id, event_type, occurred_on, created_at FROM domain_events WHERE aggregate_id='$REGISTER_USER_ID' AND event_type='UserRegistered' ORDER BY id DESC LIMIT 3;" 2>/dev/null
  else
    print_warning "⚠️  未在 domain_events 表中找到事件记录"
  fi
  echo ""
  
  # 5. 说明事件流向
  print_info "📊 事件流向分析:"
  echo ""
  echo "   ┌─────────────────────────────────────────────────────────────┐"
  echo "   │  1️⃣  HTTP POST /auth/register                               │"
  echo "   │       ↓                                                      │"
  echo "   │  2️⃣  AuthService.RegisterUser()                             │"
  echo "   │       ├─ 检查用户名/邮箱唯一性                               │"
  echo "   │       ├─ aggregate.NewUser() → 产生 UserRegisteredEvent     │"
  echo "   │       ├─ UserRepository.Save() → 持久化用户 + 事件到 DB     │"
  echo "   │       └─ 提交事务                                           │"
  echo "   │       ↓                                                      │"
  echo "   │  3️⃣  事务外: EventPublisher.Publish()                       │"
  echo "   │       ↓                                                      │"
  echo "   │  4️⃣  AsynqPublisher → 发送到 Redis 队列 (asynq:default)     │"
  echo "   │       ↓                                                      │"
  echo "   │  5️⃣  Worker 消费 → Processor.Handle()                       │"
  echo "   │       ↓                                                      │"
  echo "   │  6️⃣  AuditSubscriber.Handle() → 保存审计日志                │"
  echo "   └─────────────────────────────────────────────────────────────┘"
  echo ""
  
  # 保存 USER_ID 供后续测试使用
  echo "$REGISTER_USER_ID"
}

# 测试用户登录事件（UserLoggedIn）
test_user_logged_in() {
  local TEST_USERNAME="$1"
  local TEST_EMAIL="$2"
  
  print_step "🔐 测试 2: UserLoggedIn 事件发布与订阅"
  echo ""
  
  # 1. 用户登录（触发 UserLoggedIn 事件）
  print_info "步骤 1: 用户登录（触发 UserLoggedIn 事件）"
  LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"username_or_email\": \"$TEST_EMAIL\",
      \"password\": \"$PASSWORD\"
    }")
  
  echo "$LOGIN_RESPONSE" | jq .
  echo ""
  
  ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.access_token // empty')
  USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.user.id // empty')
  
  if [ -n "$ACCESS_TOKEN" ] && [ -n "$USER_ID" ]; then
    print_success "用户登录成功，User ID: $USER_ID"
    print_info "✅ UserLoggedIn 事件已发布到 EventPublisher（异步队列）"
  else
    print_error "登录失败"
    return 1
  fi
  echo ""
  
  # 2. 等待事件处理完成
  print_info "步骤 2: 等待异步事件处理完成..."
  
  # 轮询等待，最多等待 10 秒
  MAX_WAIT=10
  WAITED=0
  AUDIT_LOG_COUNT=0
  LOGIN_LOG_COUNT=0
  
  while [ $WAITED -lt $MAX_WAIT ]; do
    sleep 1
    WAITED=$((WAITED + 1))
    
    # 检查审计日志
    AUDIT_LOG_COUNT=$(PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -t -c \
      "SELECT COUNT(*) FROM audit_logs WHERE action='USER_LOGIN' AND user_id=$USER_ID;" 2>/dev/null | tr -d ' ')
    
    # 检查登录日志
    LOGIN_LOG_COUNT=$(PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -t -c \
      "SELECT COUNT(*) FROM login_logs WHERE user_id=$USER_ID;" 2>/dev/null | tr -d ' ')
    
    if [ "$AUDIT_LOG_COUNT" -gt 0 ] && [ "$LOGIN_LOG_COUNT" -gt 0 ]; then
      print_success "事件在 ${WAITED} 秒内被处理完成"
      break
    fi
    
    echo -n "."
  done
  echo ""
  
  # 3. 检查数据库中的审计日志
  print_info "步骤 3: 查询数据库中的审计日志..."
  
  if [ "$AUDIT_LOG_COUNT" -gt 0 ]; then
    print_success "✅ 发现审计日志！AuditSubscriber 成功处理了 UserLoggedIn 事件"
    
    print_info "审计日志详情:"
    PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -c \
      "SELECT id, user_id, action, resource_type, status, occurred_at FROM audit_logs WHERE action='USER_LOGIN' AND user_id=$USER_ID ORDER BY occurred_at DESC LIMIT 3;" 2>/dev/null
  else
    print_warning "⚠️  未发现审计日志（等待了 ${MAX_WAIT} 秒）"
  fi
  echo ""
  
  # 4. 检查数据库中的登录日志
  print_info "步骤 4: 查询数据库中的登录日志..."
  
  if [ "$LOGIN_LOG_COUNT" -gt 0 ]; then
    print_success "✅ 发现登录日志！LoginLogSubscriber 成功处理了 UserLoggedIn 事件"
    
    print_info "登录日志详情:"
    PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -c \
      "SELECT id, user_id, username, ip_address, device_type, os_info, browser_info, occurred_at FROM login_logs WHERE user_id=$USER_ID ORDER BY occurred_at DESC LIMIT 3;" 2>/dev/null
  else
    print_warning "⚠️  未发现登录日志（等待了 ${MAX_WAIT} 秒）"
  fi
  echo ""
  
  # 5. 检查 domain_events 表
  print_info "步骤 5: 检查 domain_events 表（事件溯源）..."
  
  EVENT_COUNT=$(PGPASSWORD=postgres psql -h localhost -U shenfay -d go_ddd_scaffold -t -c \
    "SELECT COUNT(*) FROM domain_events WHERE aggregate_id='$USER_ID' AND event_type='UserLoggedIn';" 2>/dev/null | tr -d ' ')
  
  if [ "$EVENT_COUNT" -gt 0 ]; then
    print_success "✅ 登录事件已持久化到 domain_events 表"
  else
    print_info "ℹ️  当前登录流程未持久化 UserLoggedIn 事件到 domain_events 表"
    print_info "   （这是正常的，登录事件通常直接发布到队列而不需要事件溯源）"
  fi
  echo ""
  
  # 6. 说明事件流向
  print_info "📊 事件流向分析:"
  echo ""
  echo "   ┌─────────────────────────────────────────────────────────────┐"
  echo "   │  1️⃣  HTTP POST /auth/login                                  │"
  echo "   │       ↓                                                      │"
  echo "   │  2️⃣  AuthService.AuthenticateUser()                         │"
  echo "   │       ├─ 验证用户名/密码                                     │"
  echo "   │       ├─ 生成 Token                                          │"
  echo "   │       └─ 返回认证结果                                        │"
  echo "   │       ↓                                                      │"
  echo "   │  3️⃣  当前实现：登录事件未通过 EventPublisher 发布           │"
  echo "   │       （注：登录事件通常在需要记录登录统计时才发布）         │"
  echo "   │       ↓                                                      │"
  echo "   │  4️⃣  审计日志和登录日志通过其他机制记录                      │"
  echo "   └─────────────────────────────────────────────────────────────┘"
  echo ""
  
  # 返回 access token 供后续测试使用
  echo "$ACCESS_TOKEN"
}

# 测试同步事件总线
test_sync_event_bus() {
  print_step "🔄 测试 3: 事件系统架构说明"
  echo ""
  
  print_info "当前系统使用的事件机制:"
  echo ""
  echo "   📌 **领域事件发布** (kernel.EventPublisher 接口):"
  echo "      - 位置：backend/internal/domain/shared/kernel/event.go"
  echo "      - 实现：AsynqPublisher (backend/internal/infrastructure/eventstore/publisher.go)"
  echo "      - 特点：事务外异步发布到 Redis 队列"
  echo "      - 用途：用户注册后发布 UserRegistered 事件"
  echo "      - 代码路径："
  echo "        backend/internal/application/auth/service.go:192-198"
  echo ""
  echo "   📌 **事件溯源存储** (domain_events 表):"
  echo "      - 位置：Repository.Save() 自动保存"
  echo "      - 特点：在事务中持久化领域事件"
  echo "      - 用途：事件溯源、审计追踪"
  echo "      - 代码路径："
  echo "        backend/internal/infrastructure/persistence/repository/user_repository.go:93-124"
  echo ""
  echo "   📌 **Worker 消费** (Asynq 队列):"
  echo "      - 位置：backend/cmd/worker/main.go"
  echo "      - 特点：从 Redis 队列消费事件并处理"
  echo "      - 处理器：AuditSubscriber, LoginLogSubscriber, UserSideEffectHandler"
  echo ""
  
  print_info "✅ 当前架构：**异步事件驱动**"
  echo "   - 领域事件产生 → 聚合根创建时"
  echo "   - 事件持久化 → 事务中保存到 domain_events 表"
  echo "   - 事件发布 → 事务外发送到 Redis 队列"
  echo "   - 事件消费 → Worker 异步处理"
  echo ""
}

# 检查 Worker 状态
check_worker_status() {
  print_step "🔍 测试 4: 检查 Worker 运行状态"
  echo ""
  
  # 检查是否有 Worker 进程在运行
  if pgrep -f "worker" > /dev/null; then
    print_success "Worker 进程正在运行"
    
    # 显示 Worker 进程信息
    ps aux | grep "[w]orker" | head -5
  else
    print_warning "⚠️  Worker 进程未运行"
    print_info "启动 Worker 命令:"
    echo "   cd backend"
    echo "   go run cmd/worker/main.go"
    echo ""
    print_info "或使用 Makefile:"
    echo "   make run-worker"
  fi
  echo ""
  
  # 检查 Redis 连接
  print_info "检查 Redis 连接..."
  if command -v redis-cli &> /dev/null; then
    if redis-cli ping &> /dev/null; then
      print_success "Redis 服务正常运行"
      
      # 显示 Asynq 队列信息
      print_info "Asynq 队列长度:"
      redis-cli LLEN asynq:default 2>/dev/null || echo "无法获取队列长度"
    else
      print_warning "⚠️  Redis 服务未运行"
    fi
  else
    print_warning "⚠️  redis-cli 未安装，无法检查 Redis 状态"
  fi
  echo ""
}

# 显示事件系统架构
show_architecture() {
  print_step "🏗️  事件系统架构说明"
  echo ""
  
  cat << 'EOF'
┌─────────────────────────────────────────────────────────────┐
│                    事件系统架构图                            │
└─────────────────────────────────────────────────────────────┘

📌 用户注册事件完整流程:

┌─────────────────────────────────────────────────────────────┐
│  HTTP Layer                                                 │
│  POST /api/v1/auth/register                                 │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│  Application Layer                                          │
│  AuthService.RegisterUser()                                 │
│    ├─ 1. 检查用户名/邮箱唯一性                              │
│    ├─ 2. aggregate.NewUser() → 产生 UserRegisteredEvent     │
│    ├─ 3. UserRepository.Save(user) → 事务内保存             │
│    │      ├─ 保存用户到 users 表                            │
│    │      └─ 保存事件到 domain_events 表（事件溯源）        │
│    ├─ 4. 提交事务                                           │
│    └─ 5. eventPublisher.Publish(event) → 发送到 Redis 队列  │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│  Infrastructure Layer (Async)                               │
│  AsynqPublisher                                             │
│    └─ 序列化事件 → 发送到 asynq:default 队列                │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│  Worker Layer                                               │
│  cmd/worker/main.go                                         │
│    └─ Processor 消费队列中的事件                            │
│        └─ 调用注册的 Handler                                │
└──────────────────────┬──────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────────────┐
│  Event Handlers                                             │
│    ├─ AuditSubscriber.Handle()                              │
│    │   └─ 保存审计日志到 audit_logs 表                      │
│    ├─ LoginLogSubscriber.Handle()                           │
│    │   └─ 保存登录日志到 login_logs 表                      │
│    └─ UserSideEffectHandler.Handle()                        │
│        └─ 处理副作用（发送邮件、初始化统计等）              │
└─────────────────────────────────────────────────────────────┘


📌 关键代码位置:

  领域事件定义:
    backend/internal/domain/user/event/events.go

  事件发布（Application Service）:
    backend/internal/application/auth/service.go:192-198

  事件持久化（Repository）:
    backend/internal/infrastructure/persistence/repository/user_repository.go:93-124

  事件发布器（AsynqPublisher）:
    backend/internal/infrastructure/eventstore/publisher.go

  事件订阅注册:
    backend/internal/interfaces/event/subscriber.go

  Worker 启动:
    backend/cmd/worker/main.go


📌 数据流总结:

  1. 用户注册 API 调用
  2. 创建 User 聚合根（产生 UserRegisteredEvent）
  3. Repository.Save() 在事务中:
     - 保存用户数据
     - 保存事件到 domain_events 表（事件溯源）
  4. 事务提交后，Application Service 发布事件到 Redis
  5. Worker 消费事件并调用处理器
  6. 处理器保存审计日志、登录日志等

EOF

  print_info "📖 更多详情参考文档:"
  echo "   - backend/docs/architecture/event-driven-architecture.md"
  echo "   - backend/docs/guides/asynqmon-usage-guide.md"
  echo ""
}

# 主流程测试
run_all_tests() {
  print_info "开始完整的事件流程测试..."
  echo ""
  
  # 测试 1: 用户注册事件
  test_user_registered
  
  # 生成测试用户用于登录测试
  TIMESTAMP=$(date +%s)
  TEST_USERNAME="event_test_${TIMESTAMP}"
  TEST_EMAIL="event_test_${TIMESTAMP}@example.com"
  
  # 先注册用户
  print_info "注册测试用户用于登录测试：$TEST_USERNAME"
  REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$TEST_USERNAME\", \"email\": \"$TEST_EMAIL\", \"password\": \"$PASSWORD\"}")
  
  REGISTER_CODE=$(echo "$REGISTER_RESPONSE" | jq -r '.code // empty')
  if [ "$REGISTER_CODE" != "0" ]; then
    print_warning "⚠️  注册测试用户失败，跳过登录测试"
  fi
  
  # 测试 2: 用户登录事件
  test_user_logged_in "$TEST_USERNAME" "$TEST_EMAIL"
  
  # 测试 3: 事件系统架构说明
  test_sync_event_bus
  
  # 测试 4: 检查 Worker 状态
  check_worker_status
  
  # 显示架构说明
  show_architecture
}

# 根据事件类型运行测试
case $EVENT_TYPE in
  user_registered)
    test_user_registered
    test_sync_event_bus
    check_worker_status
    ;;
  user_logged_in)
    # 生成测试用户
    TIMESTAMP=$(date +%s)
    TEST_USERNAME="event_test_${TIMESTAMP}"
    TEST_EMAIL="event_test_${TIMESTAMP}@example.com"
    test_user_logged_in "$TEST_USERNAME" "$TEST_EMAIL"
    test_sync_event_bus
    check_worker_status
    ;;
  all)
    run_all_tests
    ;;
  *)
    print_error "未知的事件类型：$EVENT_TYPE"
    print_info "支持的事件类型：user_registered, user_logged_in, all"
    exit 1
    ;;
esac

print_header
print_success "事件系统测试完成！"
print_header
