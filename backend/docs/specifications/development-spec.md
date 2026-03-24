# 开发规范

本文档定义了 Go DDD Scaffold 项目的开发规范和编码标准，所有贡献者必须遵守。

## 📋 目录结构规范

### 标准目录布局

```
backend/
├── cmd/                           # 应用入口
│   ├── api/main.go               # HTTP API 启动入口
│   ├── worker/main.go            # Worker 服务启动入口
│   └── cli/main.go               # CLI 工具启动入口
│
├── internal/                      # 内部实现（私有代码）
│   ├── domain/                   # 领域层（业务核心）
│   │   ├── shared/               # 共享领域模型
│   │   │   └── kernel/           # 核心概念（Entity, ValueObject, Aggregate）
│   │   ├── user/                 # 用户限界上下文
│   │   │   ├── aggregate/        # 聚合根
│   │   │   ├── valueobject/      # 值对象
│   │   │   ├── event/            # 领域事件
│   │   │   ├── service/          # 领域服务
│   │   │   └── repository/       # 仓储接口（Domain 定义）
│   │   └── tenant/               # 租户限界上下文
│   │
│   ├── application/              # 应用层（用例编排 + Ports 定义）
│   │   ├── ports/                # ⭐ Ports 接口定义
│   │   │   ├── auth/             # 认证相关 Ports
│   │   │   ├── idgen/            # ID 生成器 Ports
│   │   │   ├── cache/            # 缓存 Ports
│   │   │   └── email/            # 邮件 Ports
│   │   ├── user/                 # 用户应用服务
│   │   │   ├── service.go        # 应用服务实现
│   │   │   └── dto.go            # 数据传输对象
│   │   └── auth/                 # 认证应用服务
│   │
│   ├── infrastructure/           # 基础设施层（Adapters 实现）
│   │   ├── persistence/          # 数据持久化
│   │   │   ├── dao/              # GORM 生成的 DAO
│   │   │   └── repository/       # Repository 适配器实现
│   │   ├── platform/             # 平台服务
│   │   │   ├── auth/             # 认证基础设施
│   │   │   │   ├── jwt_service.go
│   │   │   │   └── token_service_adapter.go  # ⭐ 适配器
│   │   │   └── snowflake/        # Snowflake ID 生成器
│   │   ├── cache/redis/          # Redis 缓存实现
│   │   └── email/                # 邮件服务实现
│   │
│   ├── interfaces/               # 接口层（协议适配）
│   │   ├── http/                 # HTTP 接口
│   │   │   ├── auth/
│   │   │   │   ├── handler.go
│   │   │   │   └── routes.go
│   │   │   └── user/
│   │   └── grpc/                 # gRPC 接口
│   │
│   ├── module/                   # ⭐ 组合根（Module 组装）
│   │   ├── auth.go               # AuthModule
│   │   └── user.go               # UserModule
│   │
│   └── bootstrap/                # 启动引导
│       ├── infra.go              # Infra 容器
│       └── module.go             # Module 管理
│
├── pkg/                          # 公共库（可复用工具）
│   ├── response/                 # 统一响应格式
│   ├── util/                     # 工具函数
│   └── useragent/                # UserAgent 解析
│
├── configs/                      # 配置文件
│   ├── config.yaml
│   ├── .env.example
│   └── .env
│
├── migrations/                   # 数据库迁移脚本
├── docs/                         # 文档
└── tools/                        # 工具脚本
    ├── generator/                # 代码生成器
    ├── migrator/                 # 迁移工具
    └── core-flow-test.sh         # 集成测试
```

---

## 📝 Go 编码规范

### 命名规范

#### 1. 包命名

```go
// ✅ 正确：使用小写，避免驼峰
package user
package auth
package repository

// ❌ 错误：大写字母
package User
package AUTH

// ❌ 错误：下划线
package user_repo  // 应该用 user 或 repository
```

#### 2. 类型命名

```go
// ✅ 正确：驼峰命名，首字母大写（导出）
type UserRepository interface { ... }
type TokenService interface { ... }
type UserAggregate struct { ... }

// ✅ 正确：首字母小写（私有）
type userRepository struct { ... }
type jwtService struct { ... }
```

#### 3. 接口命名

```go
// ✅ 正确：单方法接口用 -er 后缀
type Reader interface { Read() []byte }
type Writer interface { Write([]byte) error }

// ✅ 正确：多方法接口用完整名称
type UserRepository interface {
    FindByID(id UserID) (*User, error)
    Save(user *User) error
}

// ✅ 正确：两个方法的接口可以用动名词
type Closer interface { Close() error }
```

#### 4. 变量和常量命名

```go
// ✅ 正确：见名知意
var (
    userID     int64
    username   string
    isActive   bool
    users      []*User
)

const (
    MaxRetryCount = 3
    DefaultTimeout = 30 * time.Second
)

// ❌ 错误：无意义命名
var d int  // 太短，无意义
var data interface{}  // 太泛
```

#### 5. 错误命名

```go
// ✅ 正确：错误变量以 Err 开头
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidPassword  = errors.New("invalid password")
)

// ✅ 正确：错误码以 Code 开头
const (
    CodeUserNotFound = 1001
    CodeInvalidToken = 2001
)
```

---

### 注释规范

#### 1. 包注释

```go
// ✅ 正确：每个包必须有注释
// Package user 提供用户相关的领域模型和服务
package user

// ✅ 正确：多行注释
// Package auth 提供认证和授权功能
// 
// 包含以下核心组件：
// - JWTService: JWT 令牌生成和验证
// - TokenServiceAdapter: 适配器模式实现
package auth
```

#### 2. 类型注释

```go
// ✅ 正确：导出的类型必须有注释
// User 用户聚合根
// 
// 封装用户的核心业务逻辑：
// - 注册、登录、登出
// - 密码修改、资料更新
// - 账户状态管理（激活、锁定）
type User struct {
    *kernel.Entity
    username Username
    email    Email
    password Password
}

// ✅ 正确：方法注释
// Login 用户登录
// 
// 参数:
//   - password: 密码明文
//   - ip: 登录 IP 地址
//   - userAgent: User-Agent 字符串
//
// 返回:
//   - error: 登录失败原因
func (u *User) Login(password string, ip string, userAgent string) error {
    // ...
}
```

#### 3. 函数注释

```go
// ✅ 正确：复杂函数需要详细注释
// GenerateTokenPair 生成 JWT 令牌对
//
// 参数:
//   - userID: 用户 ID
//   - username: 用户名
//   - email: 邮箱地址
//
// 返回:
//   - *TokenPair: 包含 AccessToken 和 RefreshToken
//   - error: 生成失败的错误信息
//
// 示例:
//   pair, err := svc.GenerateTokenPair(123, "john", "john@example.com")
func (s *JWTService) GenerateTokenPair(userID int64, username string, email string) (*TokenPair, error) {
    // ...
}
```

---

### 错误处理规范

#### 1. 错误处理原则

```go
// ✅ 正确：必须检查所有错误
result, err := someFunction()
if err != nil {
    return nil, err  // 或者包装错误后返回
}

// ❌ 错误：忽略错误
result, _ := someFunction()  // 除非明确知道可以忽略

// ❌ 错误：延迟检查
result, err := someFunction()
// ... 很多行代码后 ...
if err != nil {  // 太晚了
    return err
}
```

#### 2. 错误包装

```go
// ✅ 正确：使用 fmt.Errorf 包装错误
user, err := repo.FindByID(id)
if err != nil {
    return nil, fmt.Errorf("find user by id %d: %w", id, err)
}

// ✅ 正确：使用 errors.Wrap（github.com/pkg/errors）
if err != nil {
    return errors.Wrap(err, "authenticate user failed")
}

// ❌ 错误：简单拼接
return errors.New("error: " + err.Error())
```

#### 3. 自定义错误类型

```go
// ✅ 正确：定义业务错误类型
type BusinessError struct {
    Code    int
    Message string
    Cause   error
}

func (e *BusinessError) Error() string {
    return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func (e *BusinessError) Unwrap() error {
    return e.Cause
}

// 使用
return &BusinessError{
    Code:    CodeUserNotFound,
    Message: "用户不存在",
    Cause:   kernel.ErrAggregateNotFound,
}
```

#### 4. 错误日志

```go
// ✅ 正确：记录足够的上下文
logger.Error("failed to authenticate user",
    zap.String("email", email),
    zap.String("ip", ip),
    zap.Error(err),
)

// ❌ 错误：信息不足
logger.Error("auth failed", zap.Error(err))
```

---

### 代码组织规范

#### 1. 文件组织

```go
// ✅ 正确：标准文件结构
package user

import (
    // 标准库
    "context"
    "time"
    
    // 第三方库
    "go.uber.org/zap"
    
    // 项目内部
    "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
    vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// 常量定义
const (
    MaxUsernameLength = 50
)

// 变量定义
var (
    ErrUserNotFound = errors.New("user not found")
)

// 类型定义
type User struct { ... }

// 构造函数
func NewUser(...) (*User, error) { ... }

// 公有方法
func (u *User) Login(...) error { ... }

// 私有方法
func (u *User) validate() error { ... }

// 接口实现
var _ kernel.Entity = (*User)(nil)
```

#### 2. 函数长度

```go
// ✅ 推荐：函数不超过 50 行
func (s *UserService) CreateUser(ctx context.Context, cmd *CreateUserCommand) (*User, error) {
    // 1. 参数验证 (5 行)
    // 2. 业务逻辑 (20 行)
    // 3. 保存数据 (10 行)
    // 4. 发布事件 (10 行)
    // 总计：45 行 ✅
}

// ❌ 避免：函数过长
func (s *UserService) CreateUser(...) (*User, error) {
    // 100+ 行代码 ❌
    // 应该拆分为多个小函数
}
```

#### 3. 圈复杂度

```go
// ✅ 推荐：圈复杂度 < 10
func ProcessOrder(order *Order) error {
    if order.Status != StatusPending {
        return ErrInvalidStatus
    }
    
    // 简单的条件判断
    if order.Amount > 1000 {
        return s.processLargeOrder(order)
    }
    
    return s.processNormalOrder(order)
}

// ❌ 避免：过度嵌套
func ProcessOrder(order *Order) error {
    if order.Status == StatusPending {
        if order.Amount > 1000 {
            if order.VIP {
                // ... 多层嵌套
            }
        }
    }
    return nil
}
```

---

## 🏗️ 架构规范

### 分层依赖规则

```
允许依赖方向:
- Domain ← Application ← Interfaces
- Domain ← Application ← Infrastructure

禁止依赖方向:
- Application → Interfaces (❌)
- Infrastructure → Application (❌)
- Domain → Application (❌)
```

### Ports 使用规范

```go
// ✅ 正确：Application 层定义 Port
package ports

type TokenService interface {
    GenerateTokenPair(userID int64, username, email string) (*TokenPair, error)
}

// ✅ 正确：Infrastructure 层实现适配器
type TokenServiceAdapter struct {
    service *JWTService  // 具体实现
}

func (a *TokenServiceAdapter) GenerateTokenPair(...) (*ports.TokenPair, error) {
    // 类型转换
}

// ✅ 正确：Application Service 依赖 Port
type AuthServiceImpl struct {
    tokenService ports.TokenService  // 依赖接口，不是实现
}
```

### Module 组装规范

```go
// ✅ 正确：Module 作为组合根
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 注入到应用服务
    authSvc := authApp.NewAuthService(
        tokenServiceAdapter,  // 使用适配器
        ...
    )
    
    return &AuthModule{...}
}
```

---

## 📊 代码质量指标

### 测试覆盖率要求

| 层级 | 覆盖率要求 | 说明 |
|------|-----------|------|
| Domain | ≥ 90% | 业务逻辑必须充分测试 |
| Application | ≥ 80% | 核心流程需要测试 |
| Infrastructure | ≥ 60% | 重点测试适配器逻辑 |
| Interfaces | ≥ 40% | 主要测试 Handler |

### 代码审查清单

提交代码前请检查：

- [ ] 遵循命名规范
- [ ] 所有导出类型有注释
- [ ] 错误被正确处理
- [ ] 日志包含足够上下文
- [ ] 函数长度合理 (< 50 行)
- [ ] 圈复杂度低 (< 10)
- [ ] 单元测试通过
- [ ] 无 lint 警告

---

## 🔧 工具配置

### golangci-lint 配置

```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - gosimple
    - staticcheck
    - unused
    - ineffassign
    - misspell
    - errcheck

linters-settings:
  gofmt:
    simplify: true
  
  misspell:
    locale: US

run:
  timeout: 5m
```

### VS Code 设置

```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v", "-race"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

---

## 📚 参考资源

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Clean Architecture](../../design/clean-architecture-spec.md)
- [Ports 模式](../../design/ports-pattern-design.md)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
