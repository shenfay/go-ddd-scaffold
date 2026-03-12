# 统一错误处理规范

## 文档概述

本文档描述了 go-ddd-scaffold 项目中统一的错误处理机制，包括错误码体系、响应格式、错误映射规则以及在代码中的使用方法。

## 错误码体系

### 编码规则

项目采用 **5位数字错误码**，格式为 `AABBB`：
- `AA` - 大模块编号（10, 20, 30, 40... = 10的倍数）
- `BBB` - 具体错误编号（000-999）
- 成功统一为 `0`

### 模块分配

| 错误码范围 | 模块 | 说明 |
|-----------|------|------|
| 0 | 成功 | 统一成功码 |
| 10000-10999 | 系统级通用 | 通用错误、系统内部错误 |
| 20000-29999 | 用户模块 | 用户基础、认证、资料、关系 |
| 30000-39999 | 租户模块 | 租户基础、成员、角色、限制 |
| 40000-49999 | 认证授权模块 | Token、权限、验证码、安全 |
| 50000-59999 | 内容资源模块 | 内容、文件、媒体（预留） |
| 60000-69999 | 交易订单模块 | 订单、支付（预留） |
| 70000-79999 | 消息通知模块 | 消息、通知（预留） |
| 80000-89999 | 工作流模块 | 工作流、审批（预留） |

### 子模块细分示例

以用户模块为例：
```
20000-20999 - 用户基础
21000-21999 - 用户认证 (子模块21)
22000-22999 - 用户资料 (子模块22)
23000-23999 - 用户关系 (子模块23)
```

### 常用错误码速查

```go
// 成功
CodeSuccess = 0

// 系统级通用 (10000+)
CodeUnknownError    = 10000  // 未知错误
CodeInvalidParam    = 10001  // 参数无效
CodeNotFound        = 10002  // 资源不存在
CodeConflict        = 10003  // 资源冲突
CodeUnauthorized    = 10004  // 未授权
CodeForbidden       = 10005  // 禁止访问
CodeInternalError   = 10010  // 内部错误
CodeDatabaseError   = 10011  // 数据库错误
CodeConcurrency    = 10014  // 并发冲突

// 用户模块 (20000+)
CodeUserNotFound     = 20001  // 用户不存在
CodeUserExists       = 20002  // 用户已存在
CodeInvalidPassword  = 21001  // 密码错误
CodeAccountLocked    = 21004  // 账户已锁定
CodeInvalidEmail     = 22002  // 邮箱格式无效
CodeEmailExists      = 22004  // 邮箱已存在

// 租户模块 (30000+)
CodeTenantNotFound   = 30001  // 租户不存在
CodeNotTenantMember  = 31001  // 不是租户成员
CodeNotTenantOwner   = 32003  // 不是租户所有者

// 认证授权 (40000+)
CodeTokenExpired     = 40001  // Token已过期
CodeTokenInvalid     = 40002  // Token无效
CodePermissionDenied = 41001  // 权限不足
```

## 统一响应格式

### 成功响应

```json
{
    "code": 0,
    "message": "success",
    "data": {...},
    "trace_id": "optional-trace-id",
    "timestamp": 1234567890
}
```

### 分页响应

```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [...],
        "total": 100,
        "page": 1,
        "page_size": 20,
        "total_page": 5
    },
    "timestamp": 1234567890
}
```

### 错误响应

```json
{
    "code": 20001,
    "message": "用户不存在",
    "details": {
        "error_code": "USER_NOT_FOUND",
        "field": "user_id"
    },
    "trace_id": "optional-trace-id",
    "timestamp": 1234567890
}
```

### HTTP状态码映射

| 错误码范围 | HTTP状态码 |
|-----------|-----------|
| 0 (成功) | 200 OK |
| 10001-10009 | 400 Bad Request |
| 10002 | 404 Not Found |
| 10003 | 409 Conflict |
| 10004 | 401 Unauthorized |
| 10005 | 403 Forbidden |
| 10010-10999 | 500 Internal Server Error |
| 4xxxx (业务错误) | 根据具体错误码映射 |

## 代码使用指南

### 1. 错误类型定义

项目定义了三种主要错误类型：

```go
// BusinessError 业务错误
err := apperrors.NewBusinessError(CodeUserNotFound, "用户不存在")

// ValidationErrors 验证错误集合
validationErrs := &apperrors.ValidationErrors{}
validationErrs.Add("username", "用户名不能为空", nil)
err := validationErrs.ToBusinessError()

// ConcurrencyError 并发错误
err := apperrors.NewConcurrencyError(userID, expectedVersion, actualVersion, "版本冲突")
```

### 2. 错误映射器

```go
// 创建错误映射器
mapper := apperrors.NewErrorMapper()

// 映射错误到HTTP响应信息
httpStatus, code, message, details := mapper.Map(err)
// httpStatus: HTTP状态码
// code: 5位错误码
// message: 错误消息
// details: 详细错误信息
```

### 3. HTTP 中间件使用

```go
// 初始化
mapper := apperrors.NewErrorMapper()
logger := createLogger() // 开发环境支持彩色日志

// 注册中间件（按正确顺序）
router := gin.New()
router.Use(
    middleware.TraceIDMiddleware(),        // ① TraceID 追踪中间件
    gin.Logger(),                          // ② Gin 默认彩色日志
    middleware.Recovery(logger),           // ③ Panic 恢复中间件
    middleware.Error(mapper, logger),      // ④ 错误处理中间件
    middleware.LoggerWithTrace(logger),    // ⑤ 带 TraceID 的自定义日志
)
```

**中间件说明：**

| 中间件 | 作用 | 说明 |
|--------|------|------|
| `TraceIDMiddleware()` | 请求追踪 | 为每个请求生成唯一 TraceID，注入到 Header 和 Context |
| `gin.Logger()` | Gin 默认日志 | 彩色文本格式：`[GIN] 2026/03/12 - 10:18:31 \| 200 \| 571.158µs \| ::1 \| GET "/health"` |
| `Recovery(logger)` | Panic 恢复 | 捕获 panic 并返回友好错误响应 |
| `Error(mapper, logger)` | 错误处理 | 统一映射业务错误到 HTTP 响应 |
| `LoggerWithTrace(logger)` | 自定义日志 | 记录带 TraceID 的结构化日志 |

**TraceID 流转：**

- **Header**: `X-Trace-ID: abc-123-def`
- **Body**: `{ "trace_id": "abc-123-def", ... }`
- **日志**: `{"trace_id":"abc-123-def", ...}`

### 4. 控制器中使用

```go
type Handler struct {
    errorMapper *apperrors.ErrorMapper
}

func NewHandler(mapper *apperrors.ErrorMapper) *Handler {
    return &Handler{errorMapper: mapper}
}

// 成功响应
func (h *Handler) Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, response.NewResponse(data))
}

// 错误响应
func (h *Handler) Error(c *gin.Context, err error) {
    httpStatus, code, message, details := h.errorMapper.Map(err)
    c.JSON(httpStatus, response.NewErrorResponse(code, message, details))
}

// 分页响应
func (h *Handler) Page(c *gin.Context, items interface{}, total int64, page, pageSize int) {
    c.JSON(http.StatusOK, response.NewPageResponse(items, total, page, pageSize))
}

// 控制器中使用
func (ctrl *UserController) GetUser(c *gin.Context) {
    user, err := ctrl.service.GetUser(c.Request.Context(), id)
    if err != nil {
        ctrl.handler.Error(c, err)
        return
    }
    ctrl.handler.Success(c, user)
}
```

### 5. 领域层错误使用

```go
// 在领域层中使用
func (u *User) Activate() error {
    if u.status != UserStatusPending {
        return apperrors.NewBusinessError(
            apperrors.CodeInvalidParam,
            "用户不在待激活状态",
        ).WithField("status")
    }
    // ...
}
```

## 错误码扩展指南

### 添加新的业务错误码

在 `shared/errors/codes.go` 中添加：

```go
// 例如：在用户模块添加新错误
// 用户资料 (22000-22999) - 子模块22
const (
    CodeInvalidUsername    = 22001 // 用户名无效
    // 添加新错误码
    CodeProfileIncomplete  = 22010 // 用户资料不完整
)

// 在错误码映射表中添加
var errorCodes = map[int]CodeInfo{
    // ...
    CodeProfileIncomplete: {Code: 22010, HTTPStatus: 400, Message: "用户资料不完整", Module: "user"},
}
```

### 添加新的模块

1. 在 `codes.go` 中添加新的常量块（60000+）
2. 在 `errorCodes` 映射表中添加对应条目
3. 在 `mapper.go` 中添加必要的错误类型处理

## 文件结构

```
backend/
├── shared/
│   ├── errors/
│   │   ├── codes.go           # 错误码常量定义
│   │   ├── business_error.go  # 业务错误类型
│   │   └── mapper.go          # 错误映射器
│   └── response/
│       └── response.go        # 统一响应结构
└── internal/
    └── interfaces/
        └── http/
            ├── middleware/
            │   └── error_handler.go  # HTTP中间件
            └── response.go            # 响应辅助函数
```

## 设计原则

1. **分层清晰**：领域错误 → 应用错误 → HTTP响应，逐层转换
2. **错误码统一**：5位数字，按模块分类，便于监控和告警
3. **HTTP语义正确**：200/201/400/401/403/404/409/500 状态码语义正确
4. **中间件统一处理**：避免每个控制器重复处理错误
5. **全链路追踪**：通过 TraceID 实现请求全链路追踪
6. **日志分级**：4xx警告日志，5xx错误日志

## 注意事项

1. 所有新增错误码必须遵循 `AABBB` 格式
2. 错误消息应简洁明确，便于用户理解
3. 生产环境不建议返回过于详细的错误信息
4. 内部错误应记录详细日志，对外返回统一错误
