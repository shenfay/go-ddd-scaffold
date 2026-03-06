# 5 分钟快速体验

## 🎯 目标

在 5 分钟内创建一个完整的用户管理 API。

---

## ⏱️ 步骤概览

1. 创建领域模块（1 分钟）
2. 创建数据库迁移（1 分钟）
3. 实现业务逻辑（2 分钟）
4. 运行和测试（1 分钟）

---

## 📝 详细步骤

### Step 1: 创建领域模块

```bash
# 使用脚手架工具（待实现）
make create-domain name=user
```

**自动生成**:
```
backend/internal/domain/user/
├── entity/user.go          # 用户实体
├── valueobject/email.go    # 邮箱值对象
├── repository/repository.go # 仓储接口
└── event/user_events.go    # 领域事件
```

**手动创建**（如果脚手架未实现）:

创建 `entity/user.go`:

```go
package entity

import (
    "time"
    "github.com/google/uuid"
)

// User 用户聚合根
type User struct {
    ID        uuid.UUID
    Email     string
    Password  string
    Nickname  string
    Status    string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewUser 创建新用户
func NewUser(email, password, nickname string) *User {
    return &User{
        ID:        uuid.New(),
        Email:     email,
        Password:  password, // 实际应该加密
        Nickname:  nickname,
        Status:    "active",
        CreatedAt: time.Now(),
    }
}
```

---

### Step 2: 创建数据库迁移

```bash
# 创建迁移文件
make migration-create name=create_users_table
```

编辑生成的迁移文件 `migrations/sql/TIMESTAMP_create_users.sql`:

```sql
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

执行迁移:

```bash
make migrate-up
```

---

### Step 3: 实现业务逻辑

#### 3.1 创建应用服务

创建 `application/user/service/user_service.go`:

```go
package service

import (
    "context"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/domain/user/entity"
    "go-ddd-scaffold/internal/domain/user/repository"
)

type UserService struct {
    userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, email, password, nickname string) (*entity.User, error) {
    user := entity.NewUser(email, password, nickname)
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    return s.userRepo.GetByID(ctx, id)
}
```

---

#### 3.2 创建 HTTP Handler

创建 `interfaces/http/user/handler.go`:

```go
package http

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "go-ddd-scaffold/internal/application/user/service"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
    return &UserHandler{userService: svc}
}

// CreateUser godoc
// @Summary 创建用户
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "用户信息"
// @Success 201 {object} entity.User
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
        Nickname string `json:"nickname"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.userService.CreateUser(
        c.Request.Context(),
        req.Email,
        req.Password,
        req.Nickname,
    )
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary 获取用户
// @Tags users
// @Param id path string true "用户 ID"
// @Success 200 {object} entity.User
// @Router /api/users/:id [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效 ID"})
        return
    }
    
    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
        return
    }
    
    c.JSON(http.StatusOK, user)
}
```

---

#### 3.3 注册路由

编辑 `interfaces/http/router.go`:

```go
func SetupRouter() *gin.Engine {
    router := gin.Default()
    
    // 注册路由
    userHandler := http.NewUserHandler(serviceInstance)
    
    api := router.Group("/api")
    {
        api.POST("/users", userHandler.CreateUser)
        api.GET("/users/:id", userHandler.GetUser)
    }
    
    return router
}
```

---

### Step 4: 运行和测试

```bash
# 启动服务
make run

# 在另一个终端测试 API
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "nickname": "Test User"
  }'

# 获取用户
curl http://localhost:8080/api/users/{USER_ID}
```

---

## ✅ 完成！

你已经创建了一个完整的用户管理 API，包含：

- ✅ 领域模型（Entity）
- ✅ 数据持久化（Repository + Migration）
- ✅ 应用服务（Service）
- ✅ HTTP 接口（Handler）
- ✅ API 文档（Swagger）

---

## 📚 下一步学习

- [安装指南](installation.md) - 环境配置
- [配置说明](configuration.md) - 配置文件详解
- [创建领域模块](../guides/create-domain-module.md) - 深入学习 DDD 建模
