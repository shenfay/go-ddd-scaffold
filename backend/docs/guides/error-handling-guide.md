# 错误处理最佳实践指南

本指南提供 Go DDD Scaffold 项目中错误处理的最佳实践和具体示例。

## 📋 目录

1. [核心原则](#核心原则)
2. [分层错误处理](#分层错误处理)
3. [常见场景示例](#常见场景示例)
4. [错误码使用规范](#错误码使用规范)
5. [错误日志记录](#错误日志记录)
6. [常见问题](#常见问题)

---

## 🎯 核心原则

### 1. 错误分类

**业务错误（BusinessError）：**
- 违反业务规则
- 资源不存在
- 参数验证失败
- 权限不足

```go
// ✅ 正确
err := kernel.NewBusinessError(
    kernel.CodeUserNotFound,
    "用户不存在",
).WithDetails(map[string]interface{}{
    "user_id": userID,
})
```

**技术错误（Infrastructure Error）：**
- 数据库连接失败
- Redis 超时
- 第三方服务调用失败

```go
// ✅ 正确：转换为业务错误
err := kernel.NewBusinessError(
    kernel.CodeDatabaseError,
    "数据库操作失败",
).WithCause(originalErr)
```

---

### 2. 错误处理三不原则

1. **不吞掉错误** - 必须处理或向上传递
2. **不包装过度** - 保持错误的原始语义
3. **不暴露细节** - 对用户隐藏技术细节

```go
// ❌ 错误：吞掉错误
func GetUser(id int64) *User {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil // ← 错误被吞掉
    }
    return user
}

// ❌ 错误：包装过度
func GetUser(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("get user failed: %w", 
            fmt.Errorf("find by id failed: %w", 
                fmt.Errorf("database error: %w", err)))
    }
    return user, nil
}

// ✅ 正确
func GetUser(id int64) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, aggregate.ErrUserNotFound
        }
        return nil, kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "查询用户失败",
        ).WithCause(err)
    }
    return user, nil
}
```

---

## 🏗️ 分层错误处理

### Domain 层

**职责：** 定义错误类型，抛出业务错误

**原则：**
- ✅ 使用预定义的业务错误
- ✅ 不记录日志（保持领域纯净）
- ✅ 错误信息清晰明确

```go
// domain/user/aggregate/user.go
package aggregate

import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

const (
    CodeInvalidPassword = 21001
    CodeAccountLocked   = 21004
)

var (
    ErrInvalidPassword = kernel.NewBusinessError(CodeInvalidPassword, "密码错误")
    ErrAccountLocked   = kernel.NewBusinessError(CodeAccountLocked, "账户已锁定")
)

// Login 用户登录
func (u *User) Login(password string, ip string) error {
    // 检查账户状态
    if u.status == valueobject.StatusDisabled {
        return kernel.NewBusinessError(
            CodeAccountDisabled,
            "账户已禁用",
        ).WithDetails(map[string]interface{}{
            "user_id": u.ID().Value(),
        })
    }
    
    // 检查账户锁定
    if u.isLocked() {
        return ErrAccountLocked.WithDetails(map[string]interface{}{
            "locked_until": u.lockedUntil,
        })
    }
    
    // 验证密码
    if !u.password.Verify(password) {
        u.failedLoginAttempts++
        
        // 检查是否达到锁定阈值
        if u.failedLoginAttempts >= MaxLoginAttempts {
            u.Lock()
            return ErrAccountLocked
        }
        
        return ErrInvalidPassword
    }
    
    // 登录成功，重置失败计数
    u.failedLoginAttempts = 0
    u.lastLoginAt = time.Now()
    
    return nil
}
```

---

### Application 层

**职责：** 包装错误，添加应用层上下文

**原则：**
- ✅ 保留原始错误码
- ✅ 添加额外上下文信息
- ✅ 记录业务日志

```go
// application/auth/service.go
package auth

import (
    "context"
    "go.uber.org/zap"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
)

func (s *AuthServiceImpl) AuthenticateUser(
    ctx context.Context, 
    cmd *AuthenticateUserCommand,
) (*AuthResult, error) {
    // 1. 查找用户
    user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
    if err != nil {
        // ✅ 包装错误并添加上下文
        if be := kernel.AsBusinessError(err); be != nil {
            s.logger.Warn("authenticate user failed",
                zap.String("email", cmd.Email),
                zap.String("reason", "user not found"),
            )
            return nil, be.WithDetails(map[string]interface{}{
                "email":     cmd.Email,
                "operation": "authenticate",
            })
        }
        
        // ✅ 基础设施错误转换为业务错误
        s.logger.Error("authenticate user failed",
            zap.String("email", cmd.Email),
            zap.Error(err),
        )
        return nil, kernel.NewBusinessError(
            kernel.CodeInternalError,
            "认证服务异常",
        ).WithCause(err)
    }
    
    // 2. 验证密码
    if err := user.Login(cmd.Password, cmd.IP); err != nil {
        s.logger.Warn("password verification failed",
            zap.Int64("user_id", user.ID().Value()),
            zap.String("error", err.Error()),
        )
        
        // ✅ 直接返回领域错误
        return nil, err
    }
    
    // 3. 保存聚合
    if err := s.userRepo.Save(ctx, user); err != nil {
        s.logger.Error("save user failed",
            zap.Int64("user_id", user.ID().Value()),
            zap.Error(err),
        )
        return nil, kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "更新用户状态失败",
        ).WithCause(err)
    }
    
    // 4. 生成 Token
    tokenPair, err := s.tokenService.GenerateTokenPair(user.ID().Value())
    if err != nil {
        s.logger.Error("generate token failed",
            zap.Int64("user_id", user.ID().Value()),
            zap.Error(err),
        )
        return nil, kernel.NewBusinessError(
            kernel.CodeTokenGenerationFailed,
            "Token 生成失败",
        ).WithCause(err)
    }
    
    return &AuthResult{
        UserID:       user.ID().Value(),
        AccessToken:  tokenPair.AccessToken,
        RefreshToken: tokenPair.RefreshToken,
    }, nil
}
```

---

### Infrastructure 层

**职责：** 将技术错误转换为业务错误

**原则：**
- ✅ 不直接返回基础设施错误
- ✅ 转换为对应的业务错误
- ✅ 记录技术日志

```go
// infrastructure/persistence/repository/user_repository.go
package repository

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "go.uber.org/zap"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    dao_query "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao/query"
)

type userRepositoryImpl struct {
    db       *gorm.DB
    daoQuery *dao_query.Query
    logger   *zap.Logger
}

func (r *userRepositoryImpl) FindByID(
    ctx context.Context, 
    userID int64,
) (*aggregate.User, error) {
    // 1. 查询 DAO
    daoUser, err := r.daoQuery.User.Where(func(db *dao.User) bool {
        return db.ID.Eq(userID)
    }).First()
    
    if err != nil {
        // ✅ GORM 错误转换为业务错误
        if errors.Is(err, gorm.ErrRecordNotFound) {
            r.logger.Debug("user not found",
                zap.Int64("user_id", userID),
            )
            return nil, aggregate.ErrUserNotFound
        }
        
        // ✅ 记录技术日志
        r.logger.Error("database query failed",
            zap.Int64("user_id", userID),
            zap.Error(err),
            zap.Stack("stack"),
        )
        
        return nil, kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "查询用户失败",
        ).WithCause(err)
    }
    
    // 2. 转换为领域对象
    return r.toDomain(daoUser), nil
}

func (r *userRepositoryImpl) Save(
    ctx context.Context, 
    user *aggregate.User,
) error {
    // 保存逻辑...
    // 如果失败，同样转换为业务错误
    if err := tx.Save(daoUser).Error; err != nil {
        r.logger.Error("save user to database failed",
            zap.Int64("user_id", user.ID().Value()),
            zap.Error(err),
        )
        return kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "保存用户失败",
        ).WithCause(err)
    }
    
    return nil
}
```

---

### Interfaces 层

**职责：** 统一错误响应

**原则：**
- ✅ 所有错误通过 `respHandler.Error()` 处理
- ✅ 不直接使用 `c.JSON()` 或 `c.AbortWithStatusJSON()`
- ✅ 在 Handler 层完成错误响应

```go
// interfaces/http/auth/handler.go
package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/shenfay/go-ddd-scaffold/internal/application/auth"
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/shared"
)

type Handler struct {
    authService   *auth.AuthService
    respHandler   *httpShared.Handler
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
    // 1. 绑定请求参数
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // ✅ 验证错误统一处理
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "请求参数无效",
        ).WithDetails(err))
        return
    }
    
    // 2. 参数验证
    if req.Email == "" {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "邮箱不能为空",
        ).WithField("email"))
        return
    }
    
    if req.Password == "" {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "密码不能为空",
        ).WithField("password"))
        return
    }
    
    // 3. 调用应用服务
    cmd := &auth.AuthenticateUserCommand{
        Email:     req.Email,
        Password:  req.Password,
        IP:        c.ClientIP(),
        UserAgent: c.Request.UserAgent(),
    }
    
    result, err := h.authService.AuthenticateUser(c.Request.Context(), cmd)
    if err != nil {
        // ✅ 业务错误统一处理
        h.respHandler.Error(c, err)
        return
    }
    
    // 4. 返回成功响应
    h.respHandler.Success(c, result)
}

// RefreshToken 刷新 Token
func (h *Handler) RefreshToken(c *gin.Context) {
    var req RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "请求参数无效",
        ).WithDetails(err))
        return
    }
    
    if req.RefreshToken == "" {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "Refresh Token 不能为空",
        ).WithField("refresh_token"))
        return
    }
    
    result, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
    if err != nil {
        // ✅ 业务错误统一处理
        h.respHandler.Error(c, err)
        return
    }
    
    h.respHandler.Success(c, result)
}
```

---

## 🔧 常见场景示例

### 场景 1：参数验证错误

```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    
    // 1. JSON 绑定错误
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "请求参数格式错误",
        ).WithDetails(map[string]interface{}{
            "error": err.Error(),
        }))
        return
    }
    
    // 2. 必填字段验证
    if req.Username == "" {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "用户名不能为空",
        ).WithField("username"))
        return
    }
    
    // 3. 邮箱格式验证
    if !isValidEmail(req.Email) {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "邮箱格式不正确",
        ).WithField("email").WithDetails(req.Email))
        return
    }
    
    // 4. 密码强度验证
    if len(req.Password) < 8 {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "密码长度至少为 8 位",
        ).WithField("password"))
        return
    }
    
    // ... 继续处理
}
```

---

### 场景 2：资源不存在

```go
func (h *Handler) GetUser(c *gin.Context) {
    // 1. 解析路径参数
    userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodeInvalidParam,
            "无效的用户 ID",
        ).WithField("id"))
        return
    }
    
    // 2. 查询用户
    user, err := h.userService.GetUser(c.Request.Context(), userID)
    if err != nil {
        // ✅ 资源不存在返回 404
        h.respHandler.Error(c, err)
        return
    }
    
    h.respHandler.Success(c, user)
}

// Application Service
func (s *UserService) GetUser(ctx context.Context, userID int64) (*UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        if kernel.IsBusinessError(err) {
            // ✅ 保留原始错误码（CodeUserNotFound）
            return nil, err
        }
        return nil, kernel.NewBusinessError(
            kernel.CodeInternalError,
            "获取用户失败",
        ).WithCause(err)
    }
    
    return &UserResponse{
        ID:       user.ID().Value(),
        Username: user.Username().String(),
        Email:    user.Email().String(),
    }, nil
}

// Repository
func (r *userRepositoryImpl) FindByID(ctx context.Context, userID int64) (*aggregate.User, error) {
    daoUser, err := r.daoQuery.User.Where(func(db *dao.User) bool {
        return db.ID.Eq(userID)
    }).First()
    
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // ✅ 返回预定义错误（CodeUserNotFound）
            return nil, aggregate.ErrUserNotFound
        }
        return nil, kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "查询用户失败",
        ).WithCause(err)
    }
    
    return r.toDomain(daoUser), nil
}
```

---

### 场景 3：并发冲突

```go
func (s *UserService) UpdateUser(
    ctx context.Context, 
    cmd *UpdateUserCommand,
) (*UserResponse, error) {
    user, err := s.userRepo.FindByID(ctx, cmd.UserID)
    if err != nil {
        return nil, err
    }
    
    // 更新用户信息
    if err := user.Update(cmd.Username, cmd.Email); err != nil {
        return nil, err
    }
    
    // 保存时检查版本冲突
    if err := s.userRepo.Save(ctx, user); err != nil {
        if kernel.IsConcurrencyError(err) {
            // ✅ 并发冲突错误
            return nil, kernel.NewBusinessError(
                kernel.CodeConcurrency,
                "数据已被其他用户修改，请刷新后重试",
            ).WithDetails(map[string]interface{}{
                "user_id": cmd.UserID,
            })
        }
        return nil, err
    }
    
    return &UserResponse{...}, nil
}
```

---

### 场景 4：权限检查

```go
func (h *Handler) DeleteUser(c *gin.Context) {
    userID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
    currentUser := GetCurrentUser(c) // 从上下文获取当前用户
    
    // 权限检查
    if !currentUser.CanDeleteUser(userID) {
        h.respHandler.Error(c, kernel.NewBusinessError(
            kernel.CodePermissionDenied,
            "权限不足",
        ).WithDetails(map[string]interface{}{
            "required_permission": "user.delete",
            "current_user_role":   currentUser.Role,
        }))
        return
    }
    
    // 执行删除...
}
```

---

## 🎯 错误码使用规范

### 错误码范围

| 范围 | 用途 | 示例 |
|------|------|------|
| 0 | 成功 | CodeSuccess |
| 10000-10099 | 通用错误 | CodeInvalidParam, CodeNotFound |
| 20000-29999 | 用户模块 | CodeUserNotFound, CodeUserExists |
| 30000-39999 | 租户模块 | CodeTenantNotFound |
| 40000-49999 | 认证授权 | CodeTokenExpired, CodePermissionDenied |

### 使用示例

```go
// ✅ 正确：使用模块专属错误码
const (
    // 用户模块错误码 (20000-29999)
    CodeUserNotFound     = 20001
    CodeUserExists       = 20002
    CodeInvalidPassword  = 21001
    CodeAccountLocked    = 21004
)

// ❌ 错误：混用错误码
const (
    CodeUserNotFound    = 10003 // ← 使用了通用错误码
    CodeInvalidPassword = 10001 // ← 使用了通用错误码
)
```

---

## 📝 错误日志记录

### 分层日志策略

```go
// Domain 层：不记录日志
func (u *User) Login(password string) error {
    if !u.password.Verify(password) {
        return ErrInvalidPassword // ← 干净，无日志
    }
    return nil
}

// Application 层：业务日志（Warn 级别）
func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthResult, error) {
    user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
    if err != nil {
        s.logger.Warn("user login failed", // ← 业务警告
            zap.String("email", cmd.Email),
            zap.String("reason", "user not found"),
        )
        return nil, aggregate.ErrUserNotFound
    }
    // ...
}

// Infrastructure 层：技术日志（Error 级别）
func (r *userRepositoryImpl) FindByID(ctx context.Context, userID int64) (*aggregate.User, error) {
    daoUser, err := r.daoQuery.User.First()
    if err != nil {
        r.logger.Error("database query failed", // ← 技术错误
            zap.Int64("user_id", userID),
            zap.Error(err),
            zap.Stack("stack"),
        )
        return nil, aggregate.ErrUserNotFound
    }
    // ...
}
```

### 日志级别选择

| 情况 | 级别 | 示例 |
|------|------|------|
| 预期内的业务错误 | Warn | 用户登录失败、参数验证失败 |
| 技术错误 | Error | 数据库连接失败、Redis 超时 |
| 系统异常 | Error + Stack | Panic 恢复、未捕获的错误 |
| 调试信息 | Debug | SQL 查询详情、缓存命中/未命中 |

---

## ❓ 常见问题

### Q1: 什么时候应该包装错误？

**A:** 需要添加额外上下文时

```go
// ✅ 正确
user, err := s.userRepo.FindByID(ctx, userID)
if err != nil {
    return nil, kernel.NewBusinessError(
        kernel.CodeInternalError,
        "获取用户失败",
    ).WithDetails(map[string]interface{}{
        "user_id": userID,
        "operation": "get_user",
    }).WithCause(err)
}

// ❌ 错误：没有添加上下文
user, err := s.userRepo.FindByID(ctx, userID)
if err != nil {
    return nil, err
}
```

---

### Q2: 如何处理第三方服务错误？

**A:** 转换为业务错误并记录详细日志

```go
// 发送短信
err := smsService.Send(phone, code)
if err != nil {
    logger.Error("send sms failed",
        zap.String("phone", phone),
        zap.Error(err),
    )
    return kernel.NewBusinessError(
        kernel.CodeInternalError,
        "短信发送失败",
    ).WithCause(err)
}
```

---

### Q3: 是否需要为每个错误都创建 BusinessError？

**A:** 不是，基础设施错误可以统一包装

```go
// ✅ 正确：统一包装
func (r *repo) Save(ctx context.Context, entity Entity) error {
    if err := r.db.Save(entity).Error; err != nil {
        return kernel.NewBusinessError(
            kernel.CodeDatabaseError,
            "保存失败",
        ).WithCause(err)
    }
    return nil
}

// ❌ 错误：为每个 SQL 错误创建 BusinessError
func (r *repo) Save(ctx context.Context, entity Entity) error {
    if err := r.db.Save(entity).Error; err != nil {
        if strings.Contains(err.Error(), "foreign key") {
            return kernel.NewBusinessError(...)
        }
        if strings.Contains(err.Error(), "unique") {
            return kernel.NewBusinessError(...)
        }
        // ... 无限判断
    }
    return nil
}
```

---

### Q4: 如何区分客户端错误和服务器错误？

**A:** 根据 HTTP 状态码判断

```go
// 客户端错误 (4xx)
- CodeInvalidParam      → 400 Bad Request
- CodeNotFound          → 404 Not Found
- CodeUnauthorized      → 401 Unauthorized
- CodeForbidden         → 403 Forbidden
- CodeConflict          → 409 Conflict

// 服务器错误 (5xx)
- CodeInternalError     → 500 Internal Server Error
- CodeDatabaseError     → 500 Internal Server Error
- CodeCacheError        → 500 Internal Server Error
```

---

## 📚 参考资源

- [错误处理规范](specifications/error-handling-spec.md)
- [开发规范](specifications/development-spec.md)
- [架构规范](specifications/architecture-spec.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
