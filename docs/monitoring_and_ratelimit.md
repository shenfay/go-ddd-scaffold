# 监控与限流熔断使用指南

## 概述

本系统集成了 Prometheus 监控指标和限流熔断机制，确保 Token 黑名单服务的高可用性和稳定性。

---

## 一、Prometheus Metrics 监控指标

### 1.1 Redis 相关指标

```prometheus
# Redis 请求总数
redis_requests_total{operation="exists"}

# Redis 请求延迟（秒）
redis_request_duration_seconds{operation="exists"}

# Redis 错误总数
redis_errors_total{operation="exists"}

# Redis Pipeline 大小
redis_pipeline_size{operation="blacklist_batch"}
```

### 1.2 Token 黑名单指标

```prometheus
# Token 黑名单检查总次数
token_blacklist_checks_total{type="single"}    # 单次检查
token_blacklist_checks_total{type="batch"}     # 批量检查

# Token 黑名单命中次数（在黑名单中）
token_blacklist_hits_total{type="single"}
token_blacklist_hits_total{type="batch"}

# Token 黑名单未命中次数（不在黑名单中）
token_blacklist_miss_total{type="single"}
token_blacklist_miss_total{type="batch"}

# Token 黑名单检查延迟（秒）
token_blacklist_check_duration_seconds{type="single"}
token_blacklist_check_duration_seconds{type="batch"}
```

### 1.3 JWT 相关指标

```prometheus
# JWT 签发总数
jwt_issued_total{type="access_token"}

# JWT 验证总次数
jwt_validated_total{result="success"}
jwt_validated_total{result="failure"}

# JWT 错误总数
jwt_errors_total{error_type="expired"}
jwt_errors_total{error_type="malformed"}

# JWT 操作延迟（秒）
jwt_operation_duration_seconds{operation="validate"}
```

### 1.4 限流熔断指标

```prometheus
# 限流触发总次数
rate_limit_triggered_total{resource="token_blacklist"}

# 熔断器状态（0=关闭，1=打开，2=半开）
circuit_breaker_state{resource="redis"}

# 熔断器跳闸总次数
circuit_breaker_trips_total{resource="redis"}
```

---

## 二、限流熔断配置

### 2.1 限流器配置

```go
// 初始化限流器
rateLimiter := ratelimit.NewRateLimiter(
    100,   // rate: 每秒允许 100 次请求
    200,   // burst: 突发容量 200
    "token_blacklist",  // resource: 资源名称
    metrics,           // metrics: 监控指标
)
```

**参数说明：**
- `rate`: 令牌生成速率（个/秒）
- `burst`: 令牌桶容量（应对突发流量）
- `resource`: 资源标识（用于监控）

### 2.2 熔断器配置

```go
// 初始化熔断器
config := ratelimit.DefaultCircuitBreakerConfig()
config.MaxFailures = 5              // 最大失败次数：5 次
config.ResetTimeout = 30 * time.Second  // 恢复超时：30 秒
config.HalfOpenMaxCall = 3          // 半开状态允许调用数：3 次

circuitBreaker := ratelimit.NewCircuitBreaker(
    "redis",  // name: 熔断器名称
    config,   // config: 配置
    metrics,  // metrics: 监控指标
)
```

**状态转换：**
- **Closed（关闭）**：正常状态，请求正常执行
- **Open（打开）**：熔断状态，直接拒绝所有请求
- **Half-Open（半开）**：尝试恢复，允许有限次数的请求

---

## 三、集成使用示例

### 3.1 Wire 依赖注入

```go
// providers_monitoring.go

// InitializeMetrics 初始化 Prometheus 监控指标
func InitializeMetrics() *metrics.Metrics {
    registry := prometheus.DefaultRegisterer
    return metrics.NewMetrics(registry)
}

// InitializeRateLimiter 初始化限流器
func InitializeRateLimiter(metrics *metrics.Metrics) *ratelimit.RateLimiter {
    return ratelimit.NewRateLimiter(100, 200, "token_blacklist", metrics)
}

// InitializeCircuitBreaker 初始化熔断器
func InitializeCircuitBreaker(metrics *metrics.Metrics) *ratelimit.CircuitBreaker {
    config := ratelimit.DefaultCircuitBreakerConfig()
    config.MaxFailures = 5
    config.ResetTimeout = 30 * time.Second
    config.HalfOpenMaxCall = 3
    
    cb := ratelimit.NewCircuitBreaker("redis", config, metrics)
    
    // 设置状态变化回调
    cb.OnStateChange(func(state ratelimit.CircuitBreakerState) {
        // 添加日志或告警
    })
    
    return cb
}

// InitializeTokenBlacklistService 初始化 Token 黑名单服务
func InitializeTokenBlacklistService(
    rdb *redis.Client,
    metrics *metrics.Metrics,
    rateLimiter *ratelimit.RateLimiter,
    circuitBreaker *ratelimit.CircuitBreaker,
) auth.TokenBlacklistService {
    return auth.NewRedisTokenBlacklistService(
        rdb,
        "token:blacklist:",
        rateLimiter,
        circuitBreaker,
        metrics,
    )
}
```

### 3.2 HTTP Handler 中使用

```go
// handler.go

type AuthHandler struct {
    tokenBlacklist auth.TokenBlacklistService
}

func (h *AuthHandler) Logout(c *gin.Context) {
    ctx := c.Request.Context()
    token := extractToken(c)
    
    // 登出时加入黑名单（自动限流和熔断保护）
    err := h.tokenBlacklist.AddToBlacklist(ctx, token, expireAt)
    if err != nil {
        // 处理错误（可能是限流或熔断）
        if errors.Is(err, ratelimit.ErrRateLimited) {
            // 限流处理
        } else if errors.Is(err, ratelimit.ErrCircuitBreakerOpen) {
            // 熔断处理
        }
    }
}
```

### 3.3 中间件中使用

```go
// middleware/auth.go

func (m *AuthMiddleware) HandlerFunc() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        
        // 检查黑名单（带监控和限流熔断）
        isBlacklisted, err := m.tokenBlacklist.IsBlacklisted(ctx, token)
        
        if err != nil {
            if errors.Is(err, ratelimit.ErrRateLimited) {
                // 限流：记录日志但继续处理
                log.Warn("Token blacklist rate limited")
            } else if errors.Is(err, ratelimit.ErrCircuitBreakerOpen) {
                // 熔断：降级处理（暂时跳过黑名单检查）
                log.Warn("Token blacklist circuit breaker open")
            }
        }
        
        if isBlacklisted {
            c.AbortWithStatusJSON(401, "token has been revoked")
            return
        }
        
        c.Next()
    }
}
```

---

## 四、Grafana 仪表盘配置

### 4.1 Token 黑名单监控面板

```json
{
  "dashboard": {
    "title": "Token Blacklist Monitoring",
    "panels": [
      {
        "title": "QPS",
        "targets": [
          {
            "expr": "rate(token_blacklist_checks_total[1m])"
          }
        ]
      },
      {
        "title": "P99 Latency",
        "targets": [
          {
            "expr": "histogram_quantile(0.99, rate(token_blacklist_check_duration_seconds_bucket[5m]))"
          }
        ]
      },
      {
        "title": "Hit Rate",
        "targets": [
          {
            "expr": "rate(token_blacklist_hits_total[5m]) / rate(token_blacklist_checks_total[5m])"
          }
        ]
      },
      {
        "title": "Circuit Breaker State",
        "targets": [
          {
            "expr": "circuit_breaker_state{resource=\"redis\"}"
          }
        ]
      }
    ]
  }
}
```

### 4.2 告警规则

```yaml
# prometheus_alerts.yml
groups:
  - name: token_blacklist
    rules:
      # 高延迟告警
      - alert: TokenBlacklistHighLatency
        expr: histogram_quantile(0.99, rate(token_blacklist_check_duration_seconds_bucket[5m])) > 0.1
        for: 5m
        annotations:
          summary: "Token 黑名单检查延迟过高"
          
      # 熔断器打开告警
      - alert: CircuitBreakerOpen
        expr: circuit_breaker_state{resource="redis"} == 1
        for: 1m
        annotations:
          summary: "Redis 熔断器已打开"
          
      # 限流频繁触发告警
      - alert: RateLimitFrequent
        expr: rate(rate_limit_triggered_total{resource="token_blacklist"}[5m]) > 10
        for: 5m
        annotations:
          summary: "Token 黑名单限流频繁触发"
```

---

## 五、最佳实践

### 5.1 限流配置建议

| 场景 | Rate | Burst | 说明 |
|------|------|-------|------|
| 开发环境 | 50 | 100 | 低频访问 |
| 生产环境（小型） | 100 | 200 | 日活 < 10 万 |
| 生产环境（中型） | 500 | 1000 | 日活 < 100 万 |
| 生产环境（大型） | 2000 | 5000 | 日活 > 100 万 |

### 5.2 熔断配置建议

| 参数 | 推荐值 | 说明 |
|------|--------|------|
| MaxFailures | 5-10 | 根据业务容忍度调整 |
| ResetTimeout | 30s-2m | 给后端恢复时间 |
| HalfOpenMaxCall | 3-5 | 谨慎试探恢复 |

### 5.3 降级策略

```go
// 当限流或熔断时的降级处理
isBlacklisted, err := tokenBlacklist.IsBlacklisted(ctx, token)

if err != nil {
    if errors.Is(err, ratelimit.ErrRateLimited) {
        // 限流：记录日志，暂时跳过检查
        log.Warn("Rate limited, skipping blacklist check")
        isBlacklisted = false
    } else if errors.Is(err, ratelimit.ErrCircuitBreakerOpen) {
        // 熔断：使用本地缓存降级
        isBlacklisted = checkLocalCache(token)
    }
}
```

---

## 六、监控指标解读

### 6.1 QPS（Queries Per Second）

```prometheus
# 计算每秒请求数
rate(token_blacklist_checks_total[1m])
```

**健康指标：** 
- 平稳期：< 100 QPS
- 高峰期：< 500 QPS
- 突刺：允许短暂超过 1000 QPS

### 6.2 延迟分布

```prometheus
# P50 延迟
histogram_quantile(0.50, rate(token_blacklist_check_duration_seconds_bucket[5m]))

# P95 延迟
histogram_quantile(0.95, rate(token_blacklist_check_duration_seconds_bucket[5m]))

# P99 延迟
histogram_quantile(0.99, rate(token_blacklist_check_duration_seconds_bucket[5m]))
```

**健康指标：**
- P50: < 2ms
- P95: < 10ms
- P99: < 50ms

### 6.3 命中率

```prometheus
# 黑名单命中率
rate(token_blacklist_hits_total[5m]) / rate(token_blacklist_checks_total[5m])
```

**解读：**
- 命中率过高（> 20%）：可能存在大量恶意请求
- 命中率过低（< 1%）：正常状态

### 6.4 Pipeline 效率

```prometheus
# 平均 Pipeline 大小
avg(redis_pipeline_size{operation="blacklist_batch"})
```

**健康指标：**
- 平均值：50-100
- 过小：未充分利用批量优势
- 过大：可能增加单次延迟

---

## 七、故障排查

### 7.1 限流频繁触发

**现象：**
- `rate_limit_triggered_total` 持续增长
- 用户反馈 429 Too Many Requests

**排查步骤：**
1. 检查 QPS 是否突增
2. 调整限流器配置（提高 Rate/Burst）
3. 分析流量来源，识别异常请求

### 7.2 熔断器反复跳闸

**现象：**
- `circuit_breaker_trips_total` 频繁增加
- 服务间歇性不可用

**排查步骤：**
1. 检查 Redis 健康状态
2. 查看 `redis_errors_total` 错误类型
3. 调整熔断器配置（增加 MaxFailures）
4. 优化 Redis 性能（连接池、慢查询）

### 7.3 延迟飙升

**现象：**
- P99 延迟超过阈值
- Grafana 告警

**排查步骤：**
1. 检查 Redis 网络延迟
2. 查看 Pipeline 大小是否合理
3. 分析 Redis 慢查询日志
4. 考虑增加 Redis 副本

---

## 八、总结

通过集成 Prometheus 监控和限流熔断机制，Token 黑名单服务具备以下能力：

✅ **可观测性**：完整的监控指标体系  
✅ **高可用性**：限流防止过载，熔断快速失败  
✅ **自愈能力**：熔断器自动恢复机制  
✅ **性能优化**：Pipeline 批量处理提升吞吐量  

**下一步优化方向：**
- 分布式限流（Redis Lua 脚本）
- 自适应熔断（基于延迟动态调整）
- 多级缓存（Local + Redis）
