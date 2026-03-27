#!/bin/bash

# Redis 事件系统初始化脚本
# 用于创建事件去重和监控相关的 Redis 键空间

REDIS_HOST=${REDIS_HOST:-localhost}
REDIS_PORT=${REDIS_PORT:-6379}
REDIS_PASSWORD=${REDIS_PASSWORD:-""}

echo "🔄 初始化 Redis 事件系统键空间..."

# 连接 Redis
REDIS_CLI="redis-cli -h ${REDIS_HOST} -p ${REDIS_PORT}"
if [ -n "${REDIS_PASSWORD}" ]; then
    REDIS_CLI="${REDIS_CLI} -a ${REDIS_PASSWORD}"
fi

# 1. 创建事件去重相关键空间配置
echo "1. 配置事件去重键空间..."
${REDIS_CLI} CONFIG SET notify-keyspace-events KEA

# 2. 创建事件监控相关键前缀
echo "2. 初始化监控键前缀..."

# 事件处理统计键
${REDIS_CLI} SET events:stats:last_updated "$(date +%s)" EX 86400

# 队列深度监控键
${REDIS_CLI} SET events:queue:high_priority:depth 0 EX 3600
${REDIS_CLI} SET events:queue:default:depth 0 EX 3600  
${REDIS_CLI} SET events:queue:low_priority:depth 0 EX 3600

# 处理器健康状态键
${REDIS_CLI} SET events:handlers:health_status "healthy" EX 300

echo "3. 配置过期策略..."
# 为事件去重键设置合理的过期时间
${REDIS_CLI} CONFIG SET maxmemory-policy allkeys-lru

echo "✅ Redis 事件系统初始化完成!"

# 验证配置
echo "🔍 验证配置..."
echo "Keyspace notifications: $(${REDIS_CLI} CONFIG GET notify-keyspace-events)"
echo "Maxmemory policy: $(${REDIS_CLI} CONFIG GET maxmemory-policy)"

echo "📊 当前事件相关键数量:"
${REDIS_CLI} KEYS "events:*" | wc -l

echo "🚀 事件系统已准备好接收流量!"