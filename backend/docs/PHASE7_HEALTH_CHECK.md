# Phase 7: 健康检查增强 - 实施报告

**日期**: 2026-04-02  
**阶段**: Phase 7（高优先级）  
**状态**: ✅ 完成  

---

## 📋 实施概述

本次实施为项目添加了**生产级的健康检查系统**，替代了简单的 `/health` 端点，提供了完整的数据库、Redis 连接检测，以及 Kubernetes 原生支持的 Liveness 和 Readiness 探针。

### **核心功能**

1. ✅ **pkg/health/health.go** - 健康检查封装
   - 完整健康检查（`/health`）
   - Liveness 探针（`/health/live`）
   - Readiness 探针（`/health/ready`）
   - 数据库连接检测
   - Redis 连接检测
   - 响应时间监控

2. ✅ **健康状态分级**
   - `ok`: 正常
   - `warning`: 警告（慢响应）
   - `error`: 错误（连接失败）

3. ✅ **详细的组件状态**
   - 数据库：状态、响应时间、错误信息
   - Redis: 状态、响应时间、Ping 结果、错误信息

---

## 🔧 技术实现

### **1. 健康检查 API 设计**

#### **完整健康检查 - `/health`**
```bash
$ curl http://localhost:8080/health
```

**响应示例**:
```json
{
  "status": "ok",
  "timestamp": 1712048400,
  "version": "1.0.0",
  "environment": "development",
  "checks": {
    "database": {
      "status": "ok",
      "response_time_ms": "<10ms"
    },
    "redis": {
      "status": "ok",
      "response_time_ms": "<10ms",
      "ping_result": "PONG"
    }
  }
}
```

**使用场景**:
- 负载均衡器健康检查
- 监控系统数据采集
- 运维人员手动检查

---

#### **Liveness 探针 - `/health/live`**
```bash
$ curl http://localhost:8080/health/live
```

**响应示例**:
```json
{
  "status": "alive",
  "time": 1712048400
}
```

**特点**:
- 只检查服务进程是否存活
- 不依赖任何外部服务
- 快速响应（毫秒级）

**使用场景**:
- Kubernetes Liveness Probe
- 容器重启判断
- 服务存活检测

---

#### **Readiness 探针 - `/health/ready`**
```bash
$ curl http://localhost:8080/health/ready
```

**成功响应**:
```json
{
  "status": "ready",
  "time": 1712048400
}
```

**失败响应** (503):
```json
{
  "status": "not_ready",
  "database": "unhealthy",
  "redis": "healthy"
}
```

**特点**:
- 检查所有关键依赖（DB、Redis）
- 任一依赖失败返回 503
- 包含各组件详细状态

**使用场景**:
- Kubernetes Readiness Probe
- 流量接管判断
- 服务就绪检测

---

### **2. 数据库健康检查**

```go
func (h *Handler) checkDatabase(ctx context.Context) *DatabaseHealth {
    start := time.Now()
    
    // 1. 执行简单查询验证连接
    var result int
    err := h.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error
    
    duration := time.Since(start).Milliseconds()
    
    if err != nil {
        return &DatabaseHealth{
            Status: StatusError,
            Error:  err.Error(),
        }
    }
    
    // 2. 检查连接池状态
    sqlDB, _ := h.db.DB()
    stats := sqlDB.Stats()
    
    health := &DatabaseHealth{
        Status:       StatusOK,
        ResponseTime: formatDuration(duration),
    }
    
    // 3. 如果响应时间超过阈值，标记为 warning
    if duration > 1000 {
        health.Status = StatusWarning
    }
    
    // 4. 如果连接池已满，可以标记为 warning
    if stats.OpenConnections >= stats.MaxOpenConnections {
        health.Status = StatusWarning
    }
    
    return health
}
```

**检查项**:
- ✅ SQL 查询能力（`SELECT 1`）
- ✅ 响应时间检测
- ✅ 连接池状态
- ✅ 超时控制（5 秒）

---

### **3. Redis 健康检查**

```go
func (h *Handler) checkRedis(ctx context.Context) *RedisHealth {
    start := time.Now()
    
    // 1. 执行 Ping 命令
    result, err := h.redis.Ping(ctx).Result()
    duration := time.Since(start).Milliseconds()
    
    if err != nil {
        return &RedisHealth{
            Status: StatusError,
            Error:  err.Error(),
        }
    }
    
    health := &RedisHealth{
        Status:       StatusOK,
        ResponseTime: formatDuration(duration),
        PingResult:   result, // "PONG"
    }
    
    // 2. 如果响应时间超过阈值，标记为 warning
    if duration > 100 {
        health.Status = StatusWarning
    }
    
    return health
}
```

**检查项**:
- ✅ PING/PONG 响应
- ✅ 响应时间检测
- ✅ 超时控制（3 秒）

---

### **4. 响应时间格式化**

```go
func formatDuration(ms int64) string {
    if ms < 10 {
        return "<10ms"
    } else if ms < 100 {
        return "<100ms"
    } else if ms < 1000 {
        return ">100ms"
    }
    return ">1s"
}
```

**输出示例**:
- `<10ms` - 优秀
- `<100ms` - 良好
- `>100ms` - 较慢
- `>1s` - 很慢（会触发 warning）

---

## 🎯 使用示例

### **在 API 服务中集成**

```go
// cmd/api/main.go
import (
    "github.com/shenfay/go-ddd-scaffold/pkg/health"
)

func main() {
    // ... 初始化 DB 和 Redis ...
    
    db := initDatabase(cfg.Database)
    redisClient := initRedis(cfg.Redis)
    
    // 创建健康检查 Handler
    healthHandler := health.NewHandler(
        db, 
        redisClient, 
        "1.0.0",           // version
        "development",     // environment
    )
    
    router := gin.Default()
    
    // 注册健康检查路由
    healthHandler.RegisterRoutes(router)
    
    // ... 其他路由 ...
}
```

---

### **Kubernetes 配置示例**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: go-ddd-scaffold
spec:
  containers:
  - name: api
    image: go-ddd-scaffold:latest
    ports:
    - containerPort: 8080
    livenessProbe:
      httpGet:
        path: /health/live
        port: 8080
      initialDelaySeconds: 10  # 容器启动后 10 秒开始探测
      periodSeconds: 10        # 每 10 秒探测一次
      timeoutSeconds: 3        # 超时时间 3 秒
      failureThreshold: 3      # 连续 3 次失败认为不健康
    readinessProbe:
      httpGet:
        path: /health/ready
        port: 8080
      initialDelaySeconds: 5   # 容器启动后 5 秒开始探测
      periodSeconds: 5         # 每 5 秒探测一次
      timeoutSeconds: 3        # 超时时间 3 秒
      successThreshold: 2      # 连续 2 次成功才认为就绪
      failureThreshold: 3      # 连续 3 次失败认为未就绪
```

---

### **Docker Compose 配置示例**

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health/live"]
      interval: 10s
      timeout: 3s
      retries: 3
      start_period: 10s

  postgres:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: postgres
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 3

  redis:
    image: redis:7-alpine
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3
```

---

## 📊 健康检查响应示例

### **场景 1: 全部正常**
```bash
$ curl -X GET http://localhost:8080/health | jq
```

```json
{
  "status": "ok",
  "timestamp": 1712048400,
  "version": "1.0.0",
  "environment": "development",
  "checks": {
    "database": {
      "status": "ok",
      "response_time_ms": "<10ms"
    },
    "redis": {
      "status": "ok",
      "response_time_ms": "<10ms",
      "ping_result": "PONG"
    }
  }
}
```

---

### **场景 2: 数据库连接慢**
```bash
$ curl -X GET http://localhost:8080/health | jq
```

```json
{
  "status": "warning",
  "timestamp": 1712048400,
  "version": "1.0.0",
  "environment": "production",
  "checks": {
    "database": {
      "status": "warning",
      "response_time_ms": ">1s"
    },
    "redis": {
      "status": "ok",
      "response_time_ms": "<10ms",
      "ping_result": "PONG"
    }
  }
}
```

---

### **场景 3: Redis 连接失败**
```bash
$ curl -X GET http://localhost:8080/health | jq
```

```json
{
  "status": "error",
  "timestamp": 1712048400,
  "version": "1.0.0",
  "environment": "production",
  "checks": {
    "database": {
      "status": "ok",
      "response_time_ms": "<10ms"
    },
    "redis": {
      "status": "error",
      "error": "context deadline exceeded"
    }
  }
}
```

---

### **场景 4: Readiness 检查失败**
```bash
$ curl -v -X GET http://localhost:8080/health/ready
```

```
HTTP/1.1 503 Service Unavailable
Content-Type: application/json

{
  "status": "not_ready",
  "database": "unhealthy",
  "redis": "healthy"
}
```

---

## 📈 监控集成

### **Prometheus 指标采集**

可以添加一个 `/metrics` 端点用于 Prometheus 采集：

```go
router.GET("/metrics", func(c *gin.Context) {
    health := healthHandler.GetHealthStatus(c.Request.Context())
    
    // 转换为 Prometheus 指标格式
    metrics := fmt.Sprintf(`
# HELP app_health_status 应用健康状态
# TYPE app_health_status gauge
app_health_status{component="database"} %d
app_health_status{component="redis"} %d
app_health_response_time_ms{component="database"} %d
app_health_response_time_ms{component="redis"} %d
`, 
        healthStatusToInt(health.Checks.Database.Status),
        healthStatusToInt(health.Checks.Redis.Status),
        parseDuration(health.Checks.Database.ResponseTime),
        parseDuration(health.Checks.Redis.ResponseTime),
    )
    
    c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(metrics))
})
```

---

### **Grafana 仪表盘**

可以基于健康检查数据创建 Grafana 仪表盘：

```
┌─────────────────────────────────────────────┐
│  Application Health Dashboard               │
├─────────────────────────────────────────────┤
│                                             │
│  Database Health                            │
│  ┌───────────────────────────────────┐     │
│  │ ████████████████████████████ OK   │     │
│  └───────────────────────────────────┘     │
│  Response Time: <10ms                      │
│                                             │
│  Redis Health                               │
│  ┌───────────────────────────────────┐     │
│  │ ████████████████████████████ OK   │     │
│  └───────────────────────────────────┘     │
│  Response Time: <10ms                      │
│                                             │
│  Uptime: 99.99%                           │
│  Last Check: 2026-04-02 10:30:00          │
└─────────────────────────────────────────────┘
```

---

## 💡 最佳实践

### **1. 超时设置建议**

| 探针类型 | 超时时间 | 说明 |
|---------|---------|------|
| Liveness | 3 秒 | 快速判断进程是否存活 |
| Readiness | 3 秒 | 确保依赖可用 |
| Full Health | 5 秒 | 完整检查所有组件 |

### **2. 探测间隔建议**

| 环境 | Liveness 间隔 | Readiness 间隔 |
|------|-------------|---------------|
| 开发 | 30 秒 | 10 秒 |
| 生产 | 10 秒 | 5 秒 |
| 关键业务 | 5 秒 | 3 秒 |

### **3. 故障处理策略**

```yaml
livenessProbe:
  failureThreshold: 3      # 连续 3 次失败重启容器
  
readinessProbe:
  failureThreshold: 3      # 连续 3 次失败移除流量
  successThreshold: 2      # 连续 2 次成功恢复流量
```

### **4. 日志记录**

```go
func (h *Handler) HandleHealth(c *gin.Context) {
    response := h.checkHealth(c.Request.Context())
    
    // 记录健康检查结果
    if response.Status == StatusError {
        logger.Error("Health check failed",
            logger.String("status", string(response.Status)),
            logger.Any("checks", response.Checks))
    } else {
        logger.Debug("Health check passed",
            logger.String("status", string(response.Status)))
    }
    
    c.JSON(http.StatusOK, response)
}
```

---

## 📝 Git 提交历史

```bash
commit xxx
Author: AI Assistant
Date:   Thu Apr 2 2026

    feat: 添加增强的健康检查系统
    
    新增内容:
    - pkg/health/health.go: 健康检查封装
      * 完整健康检查 (/health)
      * Liveness 探针 (/health/live)
      * Readiness 探针 (/health/ready)
      * 数据库连接检测
      * Redis 连接检测
      * 响应时间监控
    
    技术特性:
    - 三级健康状态 (ok/warning/error)
    - 详细的组件状态报告
    - 响应时间格式化
    - 超时控制
    - Kubernetes 原生支持
    
    API 端点:
    - GET /health: 完整健康检查
    - GET /health/live: Liveness 探针
    - GET /health/ready: Readiness 探针
    
    使用示例:
    curl http://localhost:8080/health
    curl http://localhost:8080/health/live
    curl http://localhost:8080/health/ready
    
    集成支持:
    - Kubernetes Probes
    - Docker Compose healthcheck
    - Prometheus 指标采集
    - Grafana 仪表盘
```

---

## 🎉 总结

Phase 7 成功实现了**生产级的健康检查系统**，带来了以下优势：

✅ **完整性** - 三种探针满足不同场景需求  
✅ **可观测性** - 详细的组件状态和响应时间  
✅ **Kubernetes 原生** - 完美支持 K8s Probes  
✅ **生产就绪** - 超时控制、错误处理、日志记录  
✅ **易用性** - 简洁的 API，清晰的响应格式  
✅ **可扩展** - 易于添加新的健康检查项  

**这是提升系统可靠性的关键一步！** 🚀

---

## 📞 参考文档

- [Kubernetes Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) - 官方文档
- [QUICKSTART.md](QUICKSTART.md) - 运行和测试指南
- [ARCHITECTURE_SUMMARY.md](ARCHITECTURE_SUMMARY.md) - 整体架构说明
