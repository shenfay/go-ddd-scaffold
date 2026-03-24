# 快速参考手册

本文档提供 Go DDD Scaffold 项目的快速参考，帮助开发者快速查找常用信息。

## 📋 目录结构速查

```
backend/
├── cmd/                          # 应用入口
│   ├── api/main.go              # HTTP API
│   ├── worker/main.go           # Worker 服务
│   └── cli/main.go              # CLI 工具
│
├── internal/                     # 内部实现
│   ├── domain/                  # 领域层
│   │   ├── {module}/            # 限界上下文
│   │   │   ├── aggregate/       # 聚合根
│   │   │   ├── valueobject/     # 值对象
│   │   │   ├── event/           # 领域事件
│   │   │   ├── service/         # 领域服务
│   │   │   └── repository/      # 仓储接口
│   │   └── shared/kernel/       # 核心抽象
│   │
│   ├── application/             # 应用层
│   │   ├── ports/               # ⭐ Ports 接口
│   │   └── {module}/            # 应用服务
│   │
│   ├── infrastructure/          # 基础设施层
│   │   ├── persistence/         # 持久化
│   │   │   ├── dao/             # GORM DAO
│   │   │   └── repository/      # 仓储实现
│   │   └── platform/            # 服务实现
│   │
│   ├── interfaces/              # 接口层
│   │   └── http/                # HTTP 接口
│   │       ├── {module}/        # Handler + Routes
│   │       └── middleware/      # 中间件
│   │
│   └── module/                  # ⭐ 组合根
│       ├── auth.go
│       └── user.go
│
└── pkg/                         # 公共库
    ├── response/                # 统一响应
    └── util/                    # 工具函数
```

---

## 🎯 常见任务快速开始

### 1. 创建新模块

```bash
# 创建目录结构
mkdir -p backend/internal/domain/notification/{aggregate,valueobject,event,service,repository}
mkdir -p backend/internal/application/notification
mkdir -p backend/internal/application/ports/notification
mkdir -p backend/internal/infrastructure/platform/notification
mkdir -p backend/internal/interfaces/http/notification

# 在 module/notification.go 中注册
```

**详细指南：** [Module 开发指南](guides/module-development-guide.md)

---

### 2. 定义聚合根

```go
// domain/user/aggregate/user.go
type User struct {
    *kernel.Entity
    username vo.Username
    email    vo.Email
    password vo.Password
    status   vo.UserStatus
}

func NewUser(username, email, password string) (*User, error) {
    // 验证并创建值对象
    uName, _ := vo.NewUsername(username)
    uEmail, _ := vo.NewEmail(email)
    uPassword, _ := vo.NewPassword(password)
    
    user := &User{
        Entity:   kernel.NewEntity(),
        username: uName,
        email:    uEmail,
        password: uPassword,
        status:   vo.UserStatusPending,
    }
    
    // 发布领域事件
    user.RecordEvent(&event.UserRegistered{...})
    
    return user, nil
}

// 业务方法
func (u *User) Login(password string) error {
    if !u.password.Verify(password) {
        return kernel.ErrInvalidCredentials
    }
    // ...
}
```

**详细指南：** [开发规范](specifications/development-spec.md)

---

### 3. 定义 Repository

```go
// domain/user/repository/user_repository.go
type UserRepository interface {
    FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
    FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
    Save(ctx context.Context, user *aggregate.User) error
}

// infrastructure/persistence/repository/user_repository.go
type userRepositoryImpl struct {
    db       *gorm.DB
    daoQuery *dao_query.Query
}

func NewUserRepository(db *gorm.DB, daoQuery *dao_query.Query) repository.UserRepository {
    return &userRepositoryImpl{db, daoQuery}
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error) {
    dao, err := r.daoQuery.User.WithContext(ctx).Where(r.daoQuery.User.ID.Eq(id.Value())).First()
    if err != nil {
        return nil, err
    }
    return r.toDomain(dao)
}
```

**详细指南：** [Repository 指南](guides/repository-guide.md)

---

### 4. 创建应用服务

```go
// application/auth/service.go
type AuthServiceImpl struct {
    logger       *zap.Logger
    userRepo     repository.UserRepository
    tokenService ports.TokenService
}

func NewAuthService(logger *zap.Logger, userRepo repository.UserRepository, ...) *AuthServiceImpl {
    return &AuthServiceImpl{logger, userRepo, ...}
}

func (s *AuthServiceImpl) AuthenticateUser(ctx context.Context, cmd *AuthenticateUserCommand) (*AuthResult, error) {
    // 1. 查找用户
    user, err := s.userRepo.FindByEmail(ctx, cmd.Email)
    if err != nil {
        return nil, ErrUserNotFound
    }
    
    // 2. 调用领域方法
    err = user.Login(cmd.Password)
    if err != nil {
        return nil, err
    }
    
    // 3. 保存用户
    err = s.userRepo.Save(ctx, user)
    if err != nil {
        return nil, err
    }
    
    // 4. 生成令牌
    pair, err := s.tokenService.GenerateTokenPair(user.ID().Value(), ...)
    
    return &AuthResult{AccessToken: pair.AccessToken}, nil
}
```

**详细指南：** [架构规范](specifications/architecture-spec.md)

---

### 5. 创建 HTTP Handler

```go
// interfaces/http/auth/handler.go
type Handler struct {
    authService *app_auth.AuthServiceImpl
    respHandler *http_shared.ResponseHandler
}

func (h *Handler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.respHandler.Error(c, http.StatusBadRequest, err)
        return
    }
    
    result, err := h.authService.AuthenticateUser(c.Request.Context(), &cmd)
    if err != nil {
        h.respHandler.Error(c, http.StatusUnauthorized, err)
        return
    }
    
    h.respHandler.Success(c, result)
}
```

**详细指南：** [API 设计规范](specifications/api-spec.md)

---

### 6. 组装 Module

```go
// module/auth.go
func NewAuthModule(infra *bootstrap.Infra) *AuthModule {
    // 1. 创建基础设施
    jwtSvc := auth.NewJWTService(...)
    
    // 2. 创建适配器 ⭐
    tokenServiceAdapter := auth.NewTokenServiceAdapter(jwtSvc)
    
    // 3. 创建应用服务
    authSvc := authApp.NewAuthService(
        infra.Logger.Named("auth"),
        userRepo,
        tokenServiceAdapter,  // ← 使用 Port
        ...
    )
    
    // 4. 创建路由
    handler := authHTTP.NewHandler(authSvc, respHandler)
    routes := authHTTP.NewRoutes(handler, jwtSvc)
    
    return &AuthModule{routes: routes}
}
```

**详细指南：** [Ports 模式详解](design/ports-pattern-design.md)

---

## 🔧 常用命令

### 数据库迁移

```bash
# 安装 migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建迁移
migrate create -ext sql -dir migrations -seq create_users_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down 1
```

### 代码生成

```bash
# 生成 DAO（GORM Gen）
cd backend && go generate ./internal/infrastructure/persistence/dao

# 或使用 Makefile
make generate-dao
```

### 测试

```bash
# 运行所有测试
make test

# 运行特定包测试
go test ./internal/domain/user/... -v

# 查看测试覆盖率
make coverage
```

### 代码检查

```bash
# 运行 linter
make lint

# 格式化代码
make fmt

# 检查依赖
go mod tidy
```

### 运行项目

```bash
# 开发模式运行 API
make run

# 生产模式构建
make build

# 运行 Worker
make worker
```

---

## 📊 错误码快速查询

### 通用错误码

| Code | Message | HTTP Status |
|------|---------|-------------|
| 0 | success | 200 OK |
| 1 | internal error | 500 Internal Server Error |
| 2 | invalid params | 400 Bad Request |
| 3 | unauthorized | 401 Unauthorized |
| 4 | forbidden | 403 Forbidden |
| 5 | not found | 404 Not Found |

### 用户相关错误码

| Code | Message | HTTP Status |
|------|---------|-------------|
| 1001 | user not found | 404 Not Found |
| 1002 | username exists | 409 Conflict |
| 1003 | email exists | 409 Conflict |
| 1004 | invalid credentials | 401 Unauthorized |
| 1005 | user locked | 403 Forbidden |
| 1006 | user inactive | 403 Forbidden |

### Token 相关错误码

| Code | Message | HTTP Status |
|------|---------|-------------|
| 2001 | token expired | 401 Unauthorized |
| 2002 | token invalid | 401 Unauthorized |
| 2003 | token missing | 401 Unauthorized |

**完整列表：** [错误处理规范](specifications/error-handling-spec.md)

---

## 🏷️ 命名规范速查

### 包命名

```go
✅ 正确：
package user
package auth
package repository

❌ 错误：
package User      // 不应该大写
package user_repo // 应该用 user 或 repository
```

### 类型命名

```go
✅ 正确：
type UserRepository interface {}  // 接口：名词 + er
type UserService struct {}        // 服务：名词 + Service
type UserAggregate struct {}      // 聚合根：名词 + Aggregate
type CreateUserCommand struct {}  // 命令：动词 + 名词 + Command
type UserResponse struct {}       // 响应：名词 + Response
```

### 变量和常量

```go
✅ 正确：
var userID int64
var isActive bool
const MaxRetryCount = 3
```

### 错误变量

```go
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidPassword = errors.New("invalid password")
)

const (
    CodeUserNotFound = 1001
    CodeInvalidToken = 2001
)
```

**详细指南：** [开发规范](specifications/development-spec.md)

---

## 🎯 架构依赖规则

### 允许的依赖

```
✓ Interfaces → Application
✓ Application → Domain
✓ Infrastructure → Domain (通过适配器)
✓ Bootstrap → All
```

### 禁止的依赖

```
✗ Application → Infrastructure
✗ Domain → Application
✗ Infrastructure → Application
✗ Interfaces → Domain
```

**图示：**

```
┌─────────────┐
│ Interfaces  │
└──────┬──────┘
       ↓
┌─────────────┐
│ Application │
└──────┬──────┘
       ↓
┌─────────────┐
│   Domain    │ ← 最内层
└──────↑──────┘
       │
┌──────┴──────┐
│Infrastructure│
└─────────────┘
```

**详细指南：** [架构规范](specifications/architecture-spec.md)

---

## 📦 配置项速查

### 环境变量 (.env)

```bash
# 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=scaffold
DB_PASSWORD=scaffold
DB_NAME=go_ddd_scaffold

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_EXPIRE=7200
JWT_REFRESH_EXPIRE=604800

# 服务器
SERVER_PORT=8080
SERVER_MODE=debug
```

### 配置文件 (config.yaml)

```yaml
server:
  port: 8080
  mode: debug

jwt:
  secret: "your-secret-key"
  access_expire: 2h
  refresh_expire: 168h

database:
  driver: postgres
  host: localhost
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  db: 0
```

**详细指南：** [快速开始](guides/quickstart.md)

---

## 🔍 常见问题快速诊断

### Q1: 编译失败 - import path not found

**原因：** 依赖未安装或路径错误  
**解决：**
```bash
go mod tidy
go mod download
```

### Q2: 数据库连接失败

**原因：** 数据库未启动或配置错误  
**解决：**
```bash
# 检查 PostgreSQL
docker ps | grep postgres

# 检查配置
cat configs/.env
```

### Q3: 循环依赖错误

**原因：** 违反了架构依赖规则  
**解决：** 检查 import 语句，确保依赖指向内层

### Q4: 测试失败 - connection refused

**原因：** 测试数据库未连接  
**解决：**
```bash
# 启动测试数据库
docker-compose up -d postgres

# 设置测试环境变量
export DB_HOST=localhost
```

---

## 📚 文档索引

### 规范文档

- [开发规范](specifications/development-spec.md)
- [架构规范](specifications/architecture-spec.md)
- [API 设计规范](specifications/api-spec.md)
- [数据库规范](specifications/database-spec.md)
- [错误处理规范](specifications/error-handling-spec.md)
- [安全规范](specifications/security-spec.md)

### 使用指南

- [快速开始](guides/quickstart.md)
- [Module 开发指南](guides/module-development-guide.md)
- [Repository 指南](guides/repository-guide.md)

### 设计文档

- [架构总览](design/architecture-overview.md)
- [Ports 模式详解](design/ports-pattern-design.md)
- [Clean Architecture](design/clean-architecture-spec.md)

### 参考文档

- [技术债务与优化方案](reference/technical-debt-and-optimization.md)

---

## 🚀 学习路径推荐

### 新手开发者

```
Day 1: 快速开始 → guides/quickstart.md
Day 2: 开发规范 → specifications/development-spec.md
Day 3: Module 开发指南 → guides/module-development-guide.md
Week 2: 实战练习 → 完成第一个功能模块
```

### 资深开发者

```
Day 1: 架构总览 → design/architecture-overview.md
Day 2: 架构规范 → specifications/architecture-spec.md
Day 3: Ports 模式 → design/ports-pattern-design.md
Week 2: 深入理解 → Clean Architecture + DDD
```

### 架构师

```
Day 1: 架构总览 + 架构规范
Day 2: Clean Architecture + Ports 模式
Day 3: 技术债务与优化方案
Week 2: 架构演进规划
```

---

## 📞 获取帮助

### 文档找不到答案？

1. **搜索全文** - `grep -r "关键词" docs/`
2. **查看进度** - [DOCUMENTATION_PROGRESS.md](DOCUMENTATION_PROGRESS.md)
3. **提交 Issue** - GitHub Issues
4. **发起讨论** - GitHub Discussions

### 想贡献代码？

1. **Fork 项目**
2. **创建分支** - `git checkout -b feature/feature-name`
3. **提交更改** - `git commit -m 'Add some feature'`
4. **推送到分支** - `git push origin feature/feature-name`
5. **创建 Pull Request**

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
