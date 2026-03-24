# 错误处理规范

本文档定义了 Go DDD Scaffold 项目的错误处理规范和最佳实践。

## 📋 设计原则

### 核心原则

1. **明确性** - 错误信息清晰易懂
2. **可追溯性** - 完整的错误堆栈
3. **可恢复性** - 区分可恢复和不可恢复错误
4. **安全性** - 不泄露敏感信息

---

## 🎯 错误分类

### 按层级分类

```
┌─────────────────────────────────────┐
│         HTTP Layer Errors           │  ← HTTP 状态码相关
├─────────────────────────────────────┤
│      Application Layer Errors       │  ← 业务逻辑错误
├─────────────────────────────────────┤
│        Domain Layer Errors          │  ← 领域规则违反
├─────────────────────────────────────┤
│    Infrastructure Layer Errors      │  ← 技术基础设施错误
└─────────────────────────────────────┘
```

### 按严重程度分类

| 级别 | 说明 | 处理方式 | 示例 |
|------|------|----------|------|
| 🔴 Critical | 系统崩溃 | 立即告警，自动回滚 | 数据库宕机 |
| 🟠 Error | 功能失败 | 记录错误，返回用户 | 密码错误 |
| 🟡 Warning | 可降级 | 记录警告，继续执行 | 缓存失效 |
| 🟢 Info | 正常异常 | 记录日志，无需处理 | 资源不存在 |

---

## 🏗️ 错误类型定义

### 1. 基础错误类型

#### Kernel 错误（Domain 层）

```go
// domain/shared/kernel/errors.go
package kernel

import "errors"

// 通用领域错误
var (
    ErrAggregateNotFound   = errors.New("aggregate not found")
    ErrEntityNotFound      = errors.New("entity not found")
    ErrValueObjectInvalid  = errors.New("value object invalid")
    ErrDomainRuleViolated  = errors.New("domain rule violated")
)

// BusinessError 业务错误
type BusinessError struct {
    Code      int
    Message   string
    Cause     error
    Metadata  map[string]interface{}
}

func (e *BusinessError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

func (e *BusinessError) Unwrap() error {
    return e.Cause
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message string, opts ...ErrorOption) *BusinessError {
    err := &BusinessError{
        Code:     code,
        Message:  message,
        Metadata: make(map[string]interface{}),
    }
    
    for _, opt := range opts {
        opt(err)
    }
    
    return err
}

// ErrorOption 错误选项函数
type ErrorOption func(*BusinessError)

// WithCause 设置原因错误
func WithCause(cause error) ErrorOption {
    return func(e *BusinessError) {
        e.Cause = cause
    }
}

// WithMetadata 设置元数据
func WithMetadata(key string, value interface{}) ErrorOption {
    return func(e *BusinessError) {
        e.Metadata[key] = value
    }
}
```

#### 应用层错误

```go
// application/shared/errors.go
package shared

import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// 错误码分段定义
const (
    // 通用错误 (0-999)
    CodeOK                  = 0
    CodeInternalError       = 1
    CodeInvalidParams       = 2
    CodeUnauthorized          = 3
    CodeForbidden           = 4
    CodeNotFound            = 5
    
    // 用户相关 (1000-1999)
    CodeUserNotFound        = 1001
    CodeUsernameExists      = 1002
    CodeEmailExists         = 1003
    CodeInvalidCredentials  = 1004
    CodeUserLocked          = 1005
    CodeUserInactive        = 1006
    
    // 认证相关 (2000-2999)
    CodeTokenExpired        = 2001
    CodeTokenInvalid        = 2002
    CodeTokenMissing        = 2003
    CodeRefreshTokenInvalid = 2004
    
    // 租户相关 (3000-3999)
    CodeTenantNotFound      = 3001
    CodeTenantInactive      = 3002
    CodeTenantMemberExists  = 3003
    CodeNotTenantMember     = 3004
    
    // 权限相关 (4000-4999)
    CodeRoleNotFound        = 4001
    CodePermissionDenied    = 4002
    CodeRoleInUse           = 4003
)

// 便捷错误创建函数
var (
    ErrUserNotFound = kernel.NewBusinessError(
        CodeUserNotFound,
        "用户不存在",
    )
    
    ErrInvalidCredentials = kernel.NewBusinessError(
        CodeInvalidCredentials,
        "用户名或密码错误",
    )
    
    ErrTokenExpired = kernel.NewBusinessError(
        CodeTokenExpired,
        "令牌已过期",
    )
)
```

### 2. 验证错误

```go
// domain/shared/kernel/validation.go
package kernel

import (
    "fmt"
    "strings"
)

// ValidationError 验证错误
type ValidationError struct {
    Field   string
    Message string
    Value   interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// ValidationErrors 验证错误集合
type ValidationErrors []*ValidationError

func (e ValidationErrors) Error() string {
    msgs := make([]string, len(e))
    for i, err := range e {
        msgs[i] = err.Error()
    }
    return strings.Join(msgs, "; ")
}

// FieldError 创建字段验证错误
func FieldError(field, message string, value interface{}) *ValidationError {
    return &ValidationError{
        Field:   field,
        Message: message,
        Value:   value,
    }
}
```

### 3. 基础设施错误

```go
// infrastructure/shared/errors.go
package shared

import "fmt"

// RepositoryError 仓储错误
type RepositoryError struct {
    Operation string
    EntityType string
    Cause     error
}

func (e *RepositoryError) Error() string {
    return fmt.Sprintf("repository %s on %s: %v", e.Operation, e.EntityType, e.Cause)
}

func (e *RepositoryError) Unwrap() error {
    return e.Cause
}

// DatabaseError 数据库错误
type DatabaseError struct {
    Query   string
    Params  []interface{}
    Cause   error
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database error: %v", e.Cause)
}

func (e *DatabaseError) Unwrap() error {
    return e.Cause
}

// CacheError 缓存错误
type CacheError struct {
    Key     string
    Operation string
    Cause   error
}

func (e *CacheError) Error() string {
    return fmt.Sprintf("cache %s for key '%s': %v", e.Operation, e.Key, e.Cause)
}
```

---

## 🔄 错误处理策略

### 分层处理策略

#### 1. Domain 层 - 抛出领域错误

```go
// domain/user/aggregate/user.go
func (u *User) Login(password string, ip string, userAgent string) error {
    // 验证用户状态
    if u.Status() == vo.UserStatusLocked {
        return kernel.NewBusinessError(
            CodeUserLocked,
            "账户已被锁定",
            kernel.WithMetadata("user_id", u.ID().String()),
        )
    }
    
    if u.Status() != vo.UserStatusActive {
        return kernel.NewBusinessError(
            CodeUserInactive,
            "账户未激活",
        )
    }
    
    // 验证密码
    if !u.password.Verify(password) {
        // 记录失败尝试
        u.failedLoginAttempts++
        
        if u.failedLoginAttempts >= MaxLoginAttempts {
            u.Lock()
            return kernel.NewBusinessError(
                CodeUserLocked,
                "多次密码错误，账户已被锁定",
            )
        }
        
        return kernel.NewBusinessError(
            CodeInvalidCredentials,
            "用户名或密码错误",
        )
    }
    
    // 登录成功
    u.ResetFailedAttempts()
    u.LastLoginAt = time.Now()
    u.LastLoginIP = ip
    
    return nil
}
```

#### 2. Application 层 - 包装和转换错误

```go
// application/auth/service.go
func (s *AuthServiceImpl) AuthenticateUser(
    ctx context.Context, 
    cmd *AuthenticateUserCommand,
) (*AuthResult, error) {
    
    // 查找用户
    var user *aggregate.User
    var err error
    
    if strings.Contains(cmd.Identifier, "@") {
        user, err = s.userRepo.FindByEmail(ctx, cmd.Identifier)
    } else {
        user, err = s.userRepo.FindByUsername(ctx, cmd.Identifier)
    }
    
    if err != nil {
        if errors.Is(err, kernel.ErrAggregateNotFound) {
            // 转换为业务错误
            return nil, application_shared.ErrUserNotFound
        }
        // 包装错误，添加上下文
        return nil, fmt.Errorf("find user failed: %w", err)
    }
    
    // 调用领域方法
    err = user.Login(cmd.Password, cmd.IP, cmd.UserAgent)
    if err != nil {
        var bizErr *kernel.BusinessError
        if errors.As(err, &bizErr) {
            // 直接返回业务错误
            return nil, err
        }
        // 包装未知错误
        return nil, fmt.Errorf("user login failed: %w", err)
    }
    
    // 保存用户
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("save user failed: %w", err)
    }
    
    // 发布事件
    s.eventPublisher.Publish(&event.UserLoggedIn{
        UserID:    user.ID().String(),
        Email:     user.Email().String(),
        IP:        cmd.IP,
        UserAgent: cmd.UserAgent,
    })
    
    return &AuthResult{
        UserID:    user.ID().String(),
        Username:  user.Username().String(),
        Email:     user.Email().String(),
    }, nil
}
```

#### 3. Infrastructure 层 - 记录和包装技术错误

```go
// infrastructure/persistence/repository/user_repository.go
func (r *UserRepository) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    dao, err := r.daoQuery.User().WithContext(ctx).Where(r.idEq(id.Value())).First()
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, kernel.ErrAggregateNotFound
        }
        
        // 包装为数据库错误
        return nil, &infra_shared.DatabaseError{
            Query:  "SELECT * FROM users WHERE id = ?",
            Params: []interface{}{id.Value()},
            Cause:  err,
        }
    }
    
    return r.toDomain(dao)
}

func (r *UserRepository) Save(ctx context.Context, user *aggregate.User) error {
    // 开始事务
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 保存或更新
    err := r.saveUser(tx, user)
    if err != nil {
        tx.Rollback()
        return &infra_shared.RepositoryError{
            Operation:  "save",
            EntityType: "User",
            Cause:      err,
        }
    }
    
    // 保存领域事件
    events := user.ReleaseEvents()
    for _, event := range events {
        err = r.saveEvent(tx, event)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("save domain event failed: %w", err)
        }
    }
    
    return tx.Commit()
}
```

#### 4. Interfaces 层 - 转换为 HTTP 响应

```go
// interfaces/http/shared/response_handler.go
type ResponseHandler struct {
    errorMapper *response.ErrorMapper
    logger      *zap.Logger
}

func (h *ResponseHandler) Error(c *gin.Context, err error) {
    // 记录错误
    h.logError(c, err)
    
    // 映射为 HTTP 响应
    resp := h.errorMapper.MapError(err)
    
    // 返回响应
    c.JSON(resp.HTTPStatus, resp.Body)
}

func (h *ResponseHandler) logError(c *gin.Context, err error) {
    var bizErr *kernel.BusinessError
    
    if errors.As(err, &bizErr) {
        // 业务错误 - INFO 级别
        h.logger.Info("business error",
            zap.Int("code", bizErr.Code),
            zap.String("message", bizErr.Message),
            zap.Any("metadata", bizErr.Metadata),
            zap.String("path", c.Request.URL.Path),
            zap.String("method", c.Request.Method),
        )
    } else {
        // 系统错误 - ERROR 级别
        h.logger.Error("system error",
            zap.Error(err),
            zap.String("path", c.Request.URL.Path),
            zap.String("method", c.Request.Method),
            zap.String("client_ip", c.ClientIP()),
        )
    }
}
```

---

## 🗺️ 错误映射

### HTTP 状态码映射

```go
// pkg/response/error_mapper.go
type ErrorMapper struct {
    codeToStatus map[int]int
}

func NewErrorMapper() *ErrorMapper {
    return &ErrorMapper{
        codeToStatus: map[int]int{
            // 成功
            response.CodeOK: http.StatusOK,
            
            // 客户端错误
            response.CodeInvalidParams:    http.StatusBadRequest,
            response.CodeUnauthorized:     http.StatusUnauthorized,
            response.CodeForbidden:        http.StatusForbidden,
            response.CodeNotFound:         http.StatusNotFound,
            response.CodeConflict:         http.StatusConflict,
            
            // 用户相关
            response.CodeUserNotFound:     http.StatusNotFound,
            response.CodeInvalidCredentials: http.StatusUnauthorized,
            response.CodeUserLocked:       http.StatusForbidden,
            
            // Token 相关
            response.CodeTokenExpired:     http.StatusUnauthorized,
            response.CodeTokenInvalid:     http.StatusUnauthorized,
            response.CodeTokenMissing:     http.StatusUnauthorized,
            
            // 服务器错误
            response.CodeInternalError:    http.StatusInternalServerError,
        },
    }
}

type MappedError struct {
    HTTPStatus int
    Body       *ErrorResponse
}

func (m *ErrorMapper) MapError(err error) *MappedError {
    // 默认内部错误
    status := http.StatusInternalServerError
    code := response.CodeInternalError
    message := "内部错误"
    
    var bizErr *kernel.BusinessError
    if errors.As(err, &bizErr) {
        code = bizErr.Code
        message = bizErr.Message
        
        if s, ok := m.codeToStatus[code]; ok {
            status = s
        }
    }
    
    return &MappedError{
        HTTPStatus: status,
        Body: &ErrorResponse{
            Code:    code,
            Message: message,
            Data:    nil,
        },
    }
}
```

---

## 📝 错误日志

### 结构化日志

```go
// 使用 Zap 记录错误
logger.Error("authenticate user failed",
    zap.String("email", email),
    zap.String("ip", ip),
    zap.String("user_agent", userAgent),
    zap.Error(err),
    zap.Int64("duration_ms", duration.Milliseconds()),
)

// 带上下文的日志
logger.Error("database operation failed",
    zap.String("operation", "save"),
    zap.String("entity", "User"),
    zap.String("query", query),
    zap.Any("params", params),
    zap.Error(err),
)
```

### 错误追踪（Sentry）

```go
// 捕获错误到 Sentry
if err != nil {
    sentry.CaptureException(err)
}

// 带上下文的错误报告
hub := sentry.CurrentHub().Clone()
hub.Scope().SetContext("user", map[string]interface{}{
    "id":    userID,
    "email": email,
})
hub.Scope().SetTag("operation", "login")
hub.CaptureException(err)
```

---

## ✅ 最佳实践

### 1. 错误处理 Do's

✅ **总是检查错误**
```go
result, err := someFunction()
if err != nil {
    return nil, err
}
```

✅ **包装错误，添加上下文**
```go
user, err := repo.FindByID(id)
if err != nil {
    return nil, fmt.Errorf("find user by id %d: %w", id, err)
}
```

✅ **使用 errors.Is 和 errors.As**
```go
var bizErr *kernel.BusinessError
if errors.As(err, &bizErr) {
    // 处理业务错误
}

if errors.Is(err, kernel.ErrAggregateNotFound) {
    // 处理特定错误
}
```

✅ **记录足够的上下文**
```go
logger.Error("operation failed",
    zap.String("operation", name),
    zap.Any("params", params),
    zap.Error(err),
)
```

### 2. 错误处理 Don'ts

❌ **忽略错误**
```go
result, _ := someFunction()  // 除非明确知道可以忽略
```

❌ **简单拼接错误信息**
```go
return errors.New("error: " + err.Error())  // 丢失堆栈
```

❌ **暴露内部细节**
```go
// ❌ 不要返回原始 SQL 错误
return errors.New(fmt.Sprintf("SQL Error: %s", err.Error()))

// ✅ 应该返回友好的错误
return kernel.NewBusinessError(CodeInternalError, "操作失败")
```

❌ **过度嵌套**
```go
// ❌ 错误处理嵌套过深
if err != nil {
    if bizErr, ok := err.(*BusinessError); ok {
        if bizErr.Code == CodeUserNotFound {
            // ...
        }
    }
}

// ✅ 使用 errors.As 更清晰
var bizErr *BusinessError
if errors.As(err, &bizErr) && bizErr.Code == CodeUserNotFound {
    // ...
}
```

---

## 📊 错误码字典

### 完整错误码列表

详见：[错误码字典](../reference/error-code-dictionary.md)

---

## 📚 参考资源

- [Effective Go - Errors](https://golang.org/doc/effective_go.html#errors)
- [Dave Cheney - Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- [pkg/errors 包文档](https://github.com/pkg/errors)
- [Go 错误处理最佳实践](https://github.com/golang-standards/project-layout/blob/master/docs/error-handling.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
