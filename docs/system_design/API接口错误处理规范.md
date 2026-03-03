# MathFun API 接口错误处理规范

## 1. 概述

本文档定义了 MathFun 后端 API 的错误处理标准，包括错误响应格式、错误码体系、日志规范等。

## 2. 统一响应格式

### 2.1 成功响应

```json
{
  "code": "Success",
  "message": "操作成功",
  "data": {
    // 业务数据
  },
  "requestId": "20260203123456-abc123def",
  "timestamp": "2026-02-03T12:34:56Z"
}
```

### 2.2 错误响应

```json
{
  "code": "InvalidParameter",
  "message": "无效请求，请检查输入参数",
  "error": {
    "code": "InvalidParameter",
    "message": "无效请求，请检查输入参数",
    "details": {
      "fieldName": "详细错误信息"
    }
  },
  "requestId": "20260203123456-abc123def",
  "timestamp": "2026-02-03T12:34:56Z"
}
```

### 2.3 分页响应

```json
{
  "code": "Success",
  "message": "操作成功",
  "data": {
    "items": [],
    "total": 100,
    "page": 1,
    "pageSize": 20,
    "totalPages": 5
  },
  "requestId": "20260203123456-abc123def",
  "timestamp": "2026-02-03T12:34:56Z"
}
```

## 3. 错误码体系

### 3.1 错误码规范

错误码格式：
- **通用错误**：直接使用语义化英文，如 `Success`, `InvalidParameter`
- **业务错误**：`[模块].[具体错误]`，如 `KG.Domain.NotFound`

**层级前缀：**
- `KG.Domain`：知识图谱-领域
- `KG.Trunk`：知识图谱-主线
- `KG.Node`：知识图谱-节点
- `KG.Relationship`：知识图谱-关系
- `System`：系统错误

**设计原则：**
- 通用错误简洁无前缀
- 业务错误分层清晰，易于定位问题
- 支持国际化（message 可翻译）

### 3.2 错误码列表

| 错误码 | 说明 | HTTP 状态码 |
|-------|------|------------|
| **Success** | 操作成功 | 200 |
| **InvalidParameter** | 无效请求 | 400 |
| **MissingParameter** | 缺少必要参数 | 400 |
| **Unauthorized** | 未授权 | 401 |
| **Forbidden** | 禁止访问 | 403 |
| **NotFound** | 资源不存在 | 404 |
| **MethodNotAllowed** | 不支持的请求方法 | 405 |
| **TooManyRequests** | 请求过于频繁 | 429 |
| **ValidationFailed** | 参数校验失败 | 400 |
| **ResourceConflict** | 资源冲突 | 409 |
| **UnsupportedMediaType** | 不支持的媒体类型 | 415 |
| **KG.Domain.NotFound** | 知识领域不存在 | 404 |
| **KG.Domain.AlreadyExists** | 知识领域已存在 | 409 |
| **KG.Domain.InvalidData** | 知识领域数据无效 | 400 |
| **KG.Trunk.NotFound** | 知识主线不存在 | 404 |
| **KG.Trunk.AlreadyExists** | 知识主线已存在 | 409 |
| **KG.Trunk.NotInDomain** | 知识主线不属于指定领域 | 400 |
| **KG.Trunk.InvalidData** | 知识主线数据无效 | 400 |
| **KG.Node.NotFound** | 知识节点不存在 | 404 |
| **KG.Node.AlreadyExists** | 知识节点已存在 | 409 |
| **KG.Node.InvalidType** | 无效的节点类型 | 400 |
| **KG.Node.NotInTrunk** | 知识节点不属于指定主线 | 400 |
| **KG.Node.InvalidData** | 知识节点数据无效 | 400 |
| **KG.Relationship.NotFound** | 知识关系不存在 | 404 |
| **KG.Relationship.AlreadyExists** | 知识关系已存在 | 409 |
| **KG.Relationship.InvalidType** | 无效的关系类型 | 400 |
| **KG.Relationship.CycleDetected** | 检测到循环引用 | 400 |
| **KG.Relationship.InvalidData** | 知识关系数据无效 | 400 |
| **System.InternalError** | 系统内部错误 | 500 |
| **System.DatabaseError** | 数据库操作失败 | 500 |
| **System.CacheUnavailable** | 缓存服务不可用 | 503 |
| **System.ExternalServiceError** | 外部服务调用失败 | 502 |
| **System.Timeout** | 请求超时 | 504 |

## 4. DDD 分层错误处理规范

### 4.1 分层原则

| 层级 | 错误类型 | 是否使用预定义错误码 | 说明 |
|------|---------|-------------------|------|
| **Repository** | 技术错误 | ❌ 不使用 | 返回 gorm.ErrRecordNotFound 等底层错误 |
| **Service** | 业务错误 | ✅ 使用 | 将技术错误转换为预定义错误码 |
| **Handler** | 响应构建 | ✅ 使用 | 使用预定义错误码构建统一响应 |

### 4.2 错误码定义位置

所有预定义错误码统一位于：

```
backend/internal/pkg/errors/
├── codes.go       # 所有错误码定义
├── context.go     # 上下文相关
├── factory.go     # 错误工厂
└── types.go       # 错误类型定义
```

### 4.3 使用示例

**Repository 层 - 返回技术错误：**

```go
func (r *UserRepo) GetByID(id uuid.UUID) (*User, error) {
    var user User
    if err := r.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, err  // 返回底层技术错误
        }
        return nil, err
    }
    return &user, nil
}
```

**Service 层 - 转换为业务错误：**

```go
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.ErrUserNotFound  // 转换为预定义业务错误
        }
        return nil, errors.InternalError.Wrap(err)  // 包装系统错误
    }
    return user, nil
}
```

**Handler 层 - 构建响应：**

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        switch err {
        case errors.ErrUserNotFound:
            c.JSON(http.StatusNotFound, response.Fail(ctx, errors.ErrUserNotFound))
        default:
            c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
        }
        return
    }
    c.JSON(http.StatusOK, response.OK(ctx, user))
}
```

### 4.4 错误码添加规范

新增错误码时：

1. 在 `backend/internal/pkg/errors/codes.go` 中添加
2. 按模块用注释分隔
3. 使用 `NewCategorized` 创建

```go
// ============================================
// [模块名] 模块错误码
// ============================================
var (
    ErrXxx = NewCategorized("[模块]", "Module.Xxx", "错误描述")
)
```

## 5. 日志规范

### 4.1 结构化日志格式

所有 API 请求生成 JSON 格式日志：

```json
{
  "level": "info",
  "ts": "2026-02-03T12:34:56.789Z",
  "caller": "logger.go:123",
  "msg": "request completed successfully",
  "requestId": "20260203123456-abc123def",
  "method": "GET",
  "path": "/api/knowledge/domains",
  "query": "",
  "status": 200,
  "latency": 1.234567s,
  "clientIp": "127.0.0.1",
  "userAgent": "Mozilla/5.0...",
  "bodySize": 1024
}
```

### 4.2 日志字段说明

| 字段 | 类型 | 说明 |
|-----|------|-----|
| requestId | string | 请求唯一标识 |
| method | string | HTTP 方法 |
| path | string | 请求路径 |
| query | string | 查询参数 |
| status | int | HTTP 状态码 |
| latency | duration | 响应延迟 |
| clientIp | string | 客户端 IP |
| userAgent | string | 用户代理 |
| bodySize | int | 响应体大小 |
| errors | array | 错误列表（仅错误时） |

### 4.3 日志级别

- **DEBUG**: 详细的调试信息
- **INFO**: 常规请求日志
- **WARN**: 警告（4xx 状态码）
- **ERROR**: 错误（5xx 状态码或业务错误）

## 5. 中间件使用

### 5.1 请求 ID 中间件

```go
import "mathfun/internal/infrastructure/middleware"

r.Use(middleware.RequestIDMiddleware("X-Request-ID"))
```

功能：
- 从请求头获取或生成请求 ID
- 设置到上下文和响应头
- 用于请求追踪

### 5.2 日志中间件

```go
r.Use(middleware.Logger(middleware.DefaultLoggerConfig()))
```

功能：
- 记录请求/响应日志
- 跳过 /health 等监控端点
- 支持 JSON 和 Console 格式

### 5.3 恢复中间件

```go
logger, _ := zap.NewProduction()
r.Use(middleware.RecoveryWithLogger(logger))
```

功能：
- 捕获 panic
- 记录错误堆栈
- 返回 500 错误

### 5.4 参数校验中间件

```go
// 初始化校验器
middleware.InitValidator()

// 校验查询参数
r.GET("/domains", middleware.ValidateQuery(&QueryParams{}), handler)

// 校验请求体
r.POST("/domains", middleware.ValidateJSON(&CreateDomainRequest{}), handler)
```

## 6. 健康检查接口

### 6.1 基础健康检查

```
GET /health
```

响应：
```json
{
  "status": "up",
  "app": "mathfun",
  "timestamp": "2026-02-03T12:34:56Z"
}
```

### 6.2 详细健康检查

```
GET /health/detail
```

响应：
```json
{
  "status": "healthy",
  "app": "mathfun",
  "database": "connected",
  "databaseHost": "localhost",
  "timestamp": "2026-02-03T12:34:56Z"
}
```

## 7. 文件结构

```
backend/
├── cmd/server/main.go              # 服务入口
├── internal/
│   ├── application/knowledge/
│   │   ├── service/                # 应用服务
│   │   │   └── service.go
│   │   └── dto/                    # 数据传输对象
│   │       └── dto.go
│   ├── domain/knowledge/
│   │   ├── entity/                 # 领域实体
│   │   ├── aggregate/              # 聚合根
│   │   ├── repository/             # 仓储接口
│   │   └── service/                # 领域服务
│   ├── infrastructure/
│   │   ├── middleware/
│   │   │   ├── logger.go           # 日志中间件
│   │   │   ├── validator.go        # 校验中间件
│   │   │   └── cors.go             # CORS 配置
│   │   └── persistence/
│   │       └── gorm/repo/          # 数据仓储实现
│   ├── interfaces/
│   │   └── http/                   # HTTP 处理器
│   │       └── knowledge_handler.go
│   └── pkg/
│       └── errors/
│           └── errors.go           # 错误码定义
└── docs/
    └── system_design/
        └── API接口错误处理规范.md  # 本文档
```

## 8. 待实现功能

### 8.1 幂等性设计（优先级：中）
- 使用 `Idempotency-Key` 请求头
- Redis 缓存响应结果
- 过期时间：24 小时

### 8.2 缓存策略（优先级：中）
- HTTP 缓存头（ETag、Last-Modified）
- Redis 热点数据缓存

## 9. 变更历史

| 版本 | 日期 | 说明 |
|-----|------|------|
| v1.1 | 2026-02-14 | 新增 DDD 分层错误处理规范 |
| v1.0 | 2026-02-03 | 初始版本，包含错误码体系、日志规范、校验中间件 |
