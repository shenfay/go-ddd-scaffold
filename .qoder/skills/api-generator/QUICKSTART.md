# API Generator 快速开始指南

## 5 分钟快速上手

### 第一步：安装 Skill

```bash
npx skills install api-generator
```

### 第二步：生成 API 端点

#### 方式 1: 基于已有的 DDD 聚合

假设你已经有了 `UserAggregate`：

```bash
/api-generator --aggregate UserAggregate --with-validation --auth jwt
```

这会生成：
- ✅ CRUD 完整端点
- ✅ JWT 认证中间件集成
- ✅ 参数验证逻辑
- ✅ Swagger 文档注解
- ✅ 统一响应格式

#### 方式 2: 批量生成多个聚合

```bash
/api-generator \
  --aggregates User,Order,Product \
  --with-validation \
  --auth jwt \
  --output ./my-api
```

### 第三步：注册路由

在 `cmd/server/main.go` 中注册生成的路由：

```go
import (
    "your-project/internal/interfaces/http/user"
    "your-project/internal/interfaces/http/order"
    "your-project/internal/interfaces/http/product"
)

func main() {
    r := gin.Default()
    
    v1 := r.Group("/api/v1")
    {
        // 注册用户路由
        userHandler := http.NewUserHandler(userService)
        http.RegisterUserRoutes(v1, userHandler)
        
        // 注册订单路由
        orderHandler := http.NewOrderHandler(orderService)
        http.RegisterOrderRoutes(v1, orderHandler)
        
        // 注册商品路由
        productHandler := http.NewProductHandler(productService)
        http.RegisterProductRoutes(v1, productHandler)
    }
    
    r.Run(":8080")
}
```

### 第四步：生成 Swagger 文档

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 API 文档
swag init -g cmd/server/main.go

# 启动应用查看文档
make run
```

访问 `http://localhost:8080/swagger/index.html` 查看 API 文档

## 生成的代码结构

```
your-project/
├── internal/
│   ├── interfaces/http/
│   │   └── user/
│   │       ├── user_handler.go      # HTTP Handler
│   │       └── user_router.go       # 路由配置
│   └── application/user/dto/
│       └── user_dto.go              # 数据传输对象
└── docs/
    └── swagger.json                 # API 文档
```

## 示例：完整的 User API

### 生成的端点

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/users` | 创建用户 | ✅ JWT |
| GET | `/api/v1/users/:id` | 获取用户详情 | ✅ JWT |
| GET | `/api/v1/users` | 用户列表（分页） | ✅ JWT |
| PUT | `/api/v1/users/:id` | 更新用户 | ✅ JWT |
| DELETE | `/api/v1/users/:id` | 删除用户 | ✅ JWT |

### 请求示例

#### 创建用户

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

#### 获取用户列表

```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&pageSize=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 响应示例

成功响应：
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "uuid-here",
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2026-02-25T10:00:00Z",
    "updated_at": "2026-02-25T10:00:00Z"
  }
}
```

错误响应：
```json
{
  "code": 400,
  "message": "Validation failed",
  "errors": [
    {
      "field": "email",
      "message": "invalid email format"
    }
  ]
}
```

## 常用命令

```bash
# 生成单个聚合的 API
/api-generator --aggregate UserAggregate

# 启用 JWT 认证和验证
/api-generator --aggregate UserAggregate --with-validation --auth jwt

# 批量生成
/api-generator --aggregates User,Order,Product

# 指定输出目录
/api-generator --aggregate Product --output ./custom/path

# 自定义 API 前缀
/api-generator --aggregate Order --prefix /api/v2
```

## 配置 JWT 认证

### 1. 设置环境变量

```bash
export JWT_SECRET="your-secret-key-here"
```

### 2. 添加 JWT 中间件

在 `internal/infrastructure/web/middleware/jwt_middleware.go`：

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "time"
)

func JWTAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }
        
        // 移除 "Bearer " 前缀
        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        
        // 解析 Token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        // 提取 Claims
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("user_id", claims["user_id"])
        }
        
        c.Next()
    }
}
```

### 3. 注册中间件

```go
r := gin.Default()
r.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))
```

## 参数验证规则

支持的验证规则：

| 规则 | 说明 | 示例 |
|------|------|------|
| `required` | 必填字段 | `binding:"required"` |
| `min=2` | 最小长度 2 | `binding:"min=2"` |
| `max=50` | 最大长度 50 | `binding:"max=50"` |
| `email` | 邮箱格式 | `binding:"email"` |
| `url` | URL 格式 | `binding:"url"` |
| `uuid` | UUID 格式 | `binding:"uuid"` |
| `oneof=a b c` | 枚举值 | `binding:"oneof=a b c"` |

## 下一步

### 1. 修改业务逻辑

编辑生成的 Handler 文件，添加具体的业务逻辑：

```go
// internal/interfaces/http/user/user_handler.go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ... 现有代码 ...
    
    // TODO: 添加你的业务逻辑
    // 例如：发送欢迎邮件、初始化用户配置等
}
```

### 2. 添加自定义端点

在生成的 Router 文件中添加额外路由：

```go
// internal/interfaces/http/user/user_router.go
func RegisterUserRoutes(r *gin.RouterGroup, handler *UserHandler) {
    group := r.Group("/users")
    {
        // 标准 CRUD
        group.POST("", handler.CreateUser)
        group.GET("/:id", handler.GetUser)
        // ...
        
        // 自定义端点
        group.POST("/login", handler.Login)
        group.POST("/logout", handler.Logout)
        group.GET("/me", handler.GetCurrentUser)
    }
}
```

### 3. 编写测试

如果使用了 `--with-tests` 选项，会生成测试模板：

```go
// tests/interfaces/user/user_handler_test.go
func TestCreateUser(t *testing.T) {
    // TODO: 实现测试逻辑
}
```

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：咨询 API Design Agent

## 推荐学习路径

1. ✅ 完成本快速开始（5 分钟）
2. 📚 阅读 [REFERENCE.md](./REFERENCE.md) 了解详细配置（15 分钟）
3. 🔍 研究生成的代码结构（30 分钟）
4. 🎯 参考 [EXAMPLES.md](./EXAMPLES.md) 添加自定义端点（1 小时）
5. 🤖 咨询 API Design Agent 优化 API 设计（按需）

祝你开发顺利！🚀
