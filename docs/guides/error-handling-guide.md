# 错误处理最佳实践

## 📋 概述

本项目采用统一的错误处理机制，确保：
- ✅ 所有错误都使用 `AppError`
- ✅ 统一的 HTTP 状态码映射
- ✅ 一致的错误响应格式
- ✅ 简洁的 Handler 代码（无重复错误处理）

---

## 🎯 核心组件

### 1. AppError 错误类型

```go
// internal/pkg/errors/types.go
type AppError struct {
    CodeVal    string  // 错误码
    MsgVal     string  // 错误消息
    CatVal     string  // 分类
    DetailsVal any     // 详情
    CauseVal   error   // 根本原因
}
```

**接口方法**:
- `GetCode()` - 获取错误码
- `GetMessage()` - 获取错误消息
- `GetCategory()` - 获取分类
- `GetDetails()` - 获取详情
- `GetCause()` - 获取根本原因
- `WithDetails(d any)` - 添加详情
- `WithCause(err error)` - 添加原因

---

### 2. 错误码定义

```go
// internal/pkg/errors/codes.go

// 通用错误
InvalidParameter = NewCategorized("Common", "Common.InvalidParameter", "无效请求，请检查输入参数")
NotFound         = NewCategorized("Common", "Common.NotFound", "请求的资源不存在")

// 用户模块错误
ErrUserExists        = NewCategorized("User", "User.Exists", "用户已存在")
ErrUserNotFound      = NewCategorized("User", "User.NotFound", "用户不存在")
ErrInvalidPassword   = NewCategorized("User", "User.InvalidPassword", "密码错误")
```

**命名规范**:
- 通用错误：`InvalidParameter`, `NotFound`, `Unauthorized`
- 模块错误：`Err + 模块名 + 具体错误` (如 `ErrUserExists`)

---

### 3. HTTP 状态码映射

```go
// internal/pkg/errors/http_mapper.go

// 自动映射规则
func GetHTTPStatus(err error) (int, string) {
    appErr, ok := err.(*AppError)
    if !ok {
        return http.StatusInternalServerError, "Internal Server Error"
    }
    
    switch appErr.GetCategory() {
    case "Common":
        return mapCommonErrorToHTTP(appErr)
    case "User":
        return mapUserErrorToHTTP(appErr)
    // ...
    }
}
```

**映射表**:

| 错误类型 | HTTP 状态码 | 示例 |
|---------|-----------|------|
| InvalidParameter | 400 | 参数验证失败 |
| Unauthorized | 401 | 未授权 |
| NotFound | 404 | 资源不存在 |
| User.Exists | 409 | 用户已存在 |
| System.InternalError | 500 | 内部错误 |

---

### 4. 统一错误处理中间件

```go
// internal/interfaces/http/middleware/error_handler.go

// ErrorHandler 统一错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        // 捕获所有通过 c.Error() 记录的错误
        if len(c.Errors) > 0 {
            for _, err := range c.Errors {
                handleError(c, err)
            }
        }
    }
}

func handleError(c *gin.Context, ginErr *gin.Error) {
    ctx := c.Request.Context()
    
    // 尝试转换为 AppError
    if appErr, ok := ginErr.Err.(*errors.AppError); ok {
        // 自动映射 HTTP 状态码
        statusCode, _ := errors.GetHTTPStatus(appErr)
        
        // 返回统一响应
        c.JSON(statusCode, response.Fail(ctx, appErr))
        return
    }
    
    // 非 AppError，视为内部错误
    c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
}
```

---

## 🚀 使用指南

### 在 Handler 中使用

#### ❌ 错误方式（旧代码）

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        // ❌ 重复的错误处理逻辑
        c.JSON(http.StatusBadRequest, response.Fail(...))
        return
    }
    
    user, err := h.userService.GetUser(ctx, userID)
    if err != nil {
        if err == errors.ErrUserNotFound {
            // ❌ 每个 Handler 都要判断状态码
            c.JSON(http.StatusNotFound, response.Fail(...))
            return
        }
        h.logger.Error("获取用户失败", zap.Error(err))
        c.JSON(http.StatusInternalServerError, response.ServerErr(ctx))
        return
    }
    
    c.JSON(http.StatusOK, response.OK(ctx, user))
}
```

#### ✅ 正确方式（新标准）

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    userID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        // ✅ 使用 c.Error() 记录错误，交给中间件处理
        c.Error(errors.InvalidParameter.WithDetails("无效的用户 ID 格式"))
        return
    }
    
    user, err := h.userService.GetUser(ctx, userID)
    if err != nil {
        // ✅ 直接传递 AppError，中间件会自动处理
        c.Error(err)
        return
    }
    
    c.JSON(http.StatusOK, response.OK(ctx, user))
}
```

**优势**:
- ✅ Handler 代码更简洁（减少 60% 错误处理代码）
- ✅ 无需手动判断 HTTP 状态码
- ✅ 无需重复编写响应逻辑
- ✅ 更容易维护

---

### 在 Application Service 中使用

#### ✅ 正确示例

```go
// internal/application/user/service/user_service.go

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        // ✅ 直接返回 AppError
        return nil, err
    }
    
    return s.assembler.ToResponse(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // 1. 验证参数
    email, err := valueobject.NewEmail(req.Email)
    if err != nil {
        return nil, errors.InvalidParameter.WithDetails(err.Error())
    }
    
    // 2. 检查唯一性
    existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
    if existing != nil {
        return nil, errors.ErrUserExists
    }
    
    // 3. 创建用户
    user := entity.NewUser(...)
    
    // 4. 持久化
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, errors.Wrap(err, "CREATE_USER_FAILED", "创建用户失败")
    }
    
    return s.assembler.ToResponse(user), nil
}
```

---

### 在 Domain Service 中使用

#### ✅ 正确示例

```go
// internal/domain/user/service/user_registration_service.go

func (s *UserRegistrationService) RegisterUser(ctx context.Context, cmd RegisterCommand) (*entity.User, error) {
    // 1. 验证邮箱唯一性
    if err := s.validateEmailUnique(ctx, cmd.Email); err != nil {
        return nil, err  // 返回 AppError
    }
    
    // 2. 验证密码强度
    if err := s.validatePasswordStrength(cmd.Password); err != nil {
        return nil, errors.ErrInvalidPassword.WithDetails(err.Error())
    }
    
    // 3. 创建实体
    user := entity.NewUser(...)
    
    return user, nil
}

func (s *UserRegistrationService) validateEmailUnique(ctx context.Context, email string) error {
    existing, err := s.userRepo.GetByEmail(ctx, email)
    if err != nil {
        if err == errors.ErrUserNotFound {
            return nil  // 邮箱可用
        }
        return err
    }
    
    if existing != nil {
        return errors.ErrUserExists  // 返回 AppError
    }
    
    return nil
}
```

---

## 📊 错误处理流程

```
┌─────────────────┐
│   Controller    │
│  (HTTP Handler) │
└────────┬────────┘
         │
         │ 调用
         ↓
┌─────────────────┐
│ Application     │
│   Service       │
└────────┬────────┘
         │
         │ 调用
         ↓
┌─────────────────┐
│  Domain         │
│   Service       │
└────────┬────────┘
         │
         │ 抛出 AppError
         ↓
┌─────────────────┐
│  ErrorHandler   │ ←── 中间件捕获 c.Error()
│   Middleware    │
└────────┬────────┘
         │
         │ 自动映射 HTTP 状态码
         │ 格式化响应
         ↓
┌─────────────────┐
│  HTTP Response  │
│  {              │
│    "code": "",  │
│    "message": "",│
│    "details": {}│
│  }              │
└─────────────────┘
```

---

## 🎯 最佳实践

### 1. 错误包装

```go
// ✅ 推荐：包装底层错误，提供上下文
func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*dto.UserResponse, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        // 包装错误，添加技术上下文
        return nil, errors.Wrap(err, "GET_USER_FAILED", "从数据库获取用户失败")
    }
    return s.assembler.ToResponse(user), nil
}
```

### 2. 错误详情

```go
// ✅ 推荐：添加错误详情，便于调试
c.Error(errors.InvalidParameter.WithDetails(map[string]string{
    "field": "email",
    "value": email,
    "reason": "invalid format",
}))
```

### 3. 日志记录

```go
// ✅ 推荐：在中间件中统一记录日志
func handleError(c *gin.Context, ginErr *gin.Error) {
    if appErr, ok := ginErr.Err.(*errors.AppError); ok {
        statusCode, _ := errors.GetHTTPStatus(appErr)
        
        // 只记录服务器内部错误
        if statusCode >= 500 {
            logger.Error("Internal error", 
                zap.String("code", appErr.GetCode()),
                zap.String("message", appErr.GetMessage()),
                zap.Any("details", appErr.GetDetails()),
            )
        }
    }
}
```

### 4. 客户端 vs 服务端错误

```go
// ✅ 区分客户端错误和服务端错误
if errors.IsClientError(err) {
    // 客户端错误（4xx），记录警告日志
    logger.Warn("Client error", zap.Error(err))
} else if errors.IsServerError(err) {
    // 服务端错误（5xx），记录错误日志
    logger.Error("Server error", zap.Error(err))
}
```

---

## 🐛 常见错误

### ❌ 错误 1: 直接使用 fmt.Errorf

```go
// ❌ 禁止
return fmt.Errorf("user not found")

// ✅ 正确
return errors.ErrUserNotFound
```

### ❌ 错误 2: 裸 error 传递

```go
// ❌ 禁止
if err != nil {
    return nil, err  // 可能是 fmt.Errorf
}

// ✅ 正确
if err != nil {
    if appErr, ok := err.(*errors.AppError); ok {
        return nil, appErr
    }
    return nil, errors.Wrap(err, "UNKNOWN_ERROR", "未知错误")
}
```

### ❌ 错误 3: 在 Handler 中判断状态码

```go
// ❌ 禁止
if err == errors.ErrUserNotFound {
    c.JSON(http.StatusNotFound, response.Fail(...))
    return
}

// ✅ 正确
c.Error(err)  // 中间件会自动判断
```

### ❌ 错误 4: 忽略错误

```go
// ❌ 禁止
user, _ := s.userRepo.GetByID(ctx, id)  // 忽略错误

// ✅ 正确
user, err := s.userRepo.GetByID(ctx, id)
if err != nil {
    return nil, err
}
```

---

## 📋 检查清单

在提交代码前，请确认：

- [ ] 所有错误都使用 `AppError`
- [ ] 没有使用 `fmt.Errorf`
- [ ] Handler 中使用 `c.Error(err)` 而非 `c.JSON(..., response.Fail(...))`
- [ ] 错误消息清晰明确
- [ ] 必要时添加了错误详情
- [ ] 敏感信息已脱敏（密码、Token 等）

---

## 🔧 工具函数

### 错误判断

```go
// 是否客户端错误
errors.IsClientError(err)

// 是否服务端错误
errors.IsServerError(err)

// 是否网络错误
errors.IsNetworkError(err)

// 是否应该重试
errors.ShouldRetry(err)
```

### 错误包装

```go
// 包装错误
errors.Wrap(err, "CODE", "message")

// 添加详情
errors.ErrUserNotFound.WithDetails(userID)

// 添加原因
errors.InternalError.WithCause(originalErr)
```

---

## 📚 相关文档

- [代码规范](standards/code-style.md) - 错误处理规范
- [DDD 实现规范](standards/ddd-implementation.md) - 分层职责
- [添加 API 端点指南](guides/add-api-endpoint.md) - Handler 编写

---

**版本**: v1.0  
**生效日期**: 2026-03-06  
**下次回顾**: 2026-03-20
