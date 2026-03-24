# API 设计规范

本文档定义了 Go DDD Scaffold 项目的 RESTful API 设计规范。

## 📋 设计原则

### RESTful 风格

1. **资源导向** - 所有操作围绕资源进行
2. **无状态** - 每次请求包含所有必要信息
3. **统一接口** - 使用标准 HTTP 方法
4. **可缓存** - 响应明确是否可缓存

### HTTP 方法语义

| 方法 | 用途 | 幂等性 |
|------|------|--------|
| GET | 获取资源 | ✅ 是 |
| POST | 创建资源/执行操作 | ❌ 否 |
| PUT | 更新资源（全量） | ✅ 是 |
| PATCH | 更新资源（部分） | ✅ 是 |
| DELETE | 删除资源 | ✅ 是 |

---

## 🎯 URL 设计规范

### 资源命名

```
✅ 正确：
GET    /api/users              # 获取用户列表
POST   /api/users              # 创建用户
GET    /api/users/{id}         # 获取指定用户
PUT    /api/users/{id}         # 更新用户
DELETE /api/users/{id}         # 删除用户

❌ 错误：
GET    /api/getUsers           # 不应该用动词
POST   /api/createUser         # 动词应该在 HTTP 方法中体现
```

### 嵌套资源

```
✅ 正确：
GET /api/users/{userId}/tenants      # 获取用户的租户列表
GET /api/tenants/{tenantId}/members  # 获取租户的成员列表

❌ 错误：
GET /api/userTenants?userId=xxx      # 避免扁平化查询参数
```

### 复数形式

```
✅ 正确：使用名词复数
/api/users
/api/tenants
/api/roles

❌ 错误：单数形式
/api/user
/api/tenant
```

---

## 📦 请求响应格式

### 标准响应结构

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 实际数据
  },
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-03-23T10:00:00Z"
  }
}
```

### 成功响应

```json
// 单个资源
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 123,
    "username": "john",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}

// 资源列表（带分页）
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [
      {"id": 1, "username": "john"},
      {"id": 2, "username": "jane"}
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### 错误响应

```json
{
  "code": 1001,
  "message": "用户不存在",
  "data": null,
  "meta": {
    "request_id": "req_123456789",
    "timestamp": "2024-03-23T10:00:00Z",
    "details": {
      "field": "user_id",
      "value": "999"
    }
  }
}
```

---

## 🔐 认证授权

### JWT Token 使用

#### 请求头

```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 刷新 Token

```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 权限标识

```go
// RBAC 权限检查
GET /api/admin/users     // 需要 admin 角色
POST /api/users          // 需要 user:create 权限
```

---

## 📊 分页规范

### Query 参数

```
GET /api/users?page=1&page_size=20&sort=-created_at
```

### 参数说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| page | int | 1 | 页码（从 1 开始） |
| page_size | int | 20 | 每页数量（最大 100） |
| sort | string | -created_at | 排序字段（- 表示降序） |

### 响应中的分页信息

```json
{
  "data": {
    "items": [...]
  },
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100,
    "total_pages": 5,
    "has_next": true,
    "has_prev": false
  }
}
```

---

## 🔍 查询和过滤

### 过滤参数

```
GET /api/users?status=active&role=admin&created_after=2024-01-01
```

### 常用过滤器

| 参数 | 说明 | 示例 |
|------|------|------|
| status | 状态过滤 | `?status=active` |
| role | 角色过滤 | `?role=admin` |
| created_after | 创建时间之后 | `?created_after=2024-01-01` |
| created_before | 创建时间之前 | `?created_before=2024-12-31` |
| search | 全文搜索 | `?search=john` |

### 范围查询

```
GET /api/users?age_gte=18&age_lte=60
// gte: >=, lte: <=, gt: >, lt: <
```

---

## ⚠️ 错误处理

### HTTP 状态码

| 状态码 | 含义 | 使用场景 |
|--------|------|----------|
| 200 OK | 成功 | GET/PUT/PATCH 成功 |
| 201 Created | 已创建 | POST 创建资源成功 |
| 204 No Content | 无内容 | DELETE 成功 |
| 400 Bad Request | 请求错误 | 参数验证失败 |
| 401 Unauthorized | 未授权 | Token 无效或过期 |
| 403 Forbidden | 禁止访问 | 权限不足 |
| 404 Not Found | 未找到 | 资源不存在 |
| 409 Conflict | 冲突 | 资源已存在 |
| 422 Unprocessable Entity | 无法处理 | 业务验证失败 |
| 429 Too Many Requests | 请求过多 | 限流 |
| 500 Internal Server Error | 服务器错误 | 系统异常 |

### 错误码体系

```go
// 错误码分段定义
const (
    // 通用错误 (0-999)
    CodeOK                  = 0
    CodeInternalError       = 1
    CodeInvalidParams       = 2
    
    // 用户相关 (1000-1999)
    CodeUserNotFound        = 1001
    CodeUsernameExists      = 1002
    CodeEmailExists         = 1003
    CodeInvalidCredentials  = 1004
    
    // 认证相关 (2000-2999)
    CodeTokenExpired        = 2001
    CodeTokenInvalid        = 2002
    CodeTokenMissing        = 2003
    
    // 租户相关 (3000-3999)
    CodeTenantNotFound      = 3001
    CodeTenantInactive      = 3002
)
```

---

## 📝 版本控制

### URL 版本化

```
✅ 推荐：
/api/v1/users
/api/v2/users

❌ 避免：
/api/users/v1
```

### 请求头版本化

```http
Accept: application/vnd.go-ddd.v1+json
Accept: application/vnd.go-ddd.v2+json
```

---

## 🚀 性能优化

### 字段选择

```
GET /api/users?fields=id,username,email
// 只返回指定字段，减少数据传输
```

### 批量操作

```
POST /api/users/batch
{
  "users": [
    {"username": "user1", "email": "user1@example.com"},
    {"username": "user2", "email": "user2@example.com"}
  ]
}
```

### 缓存控制

```http
# 响应头
Cache-Control: max-age=3600, public
ETag: "abc123"
Last-Modified: Wed, 21 Oct 2024 07:28:00 GMT

# 请求头
If-None-Match: "abc123"
If-Modified-Since: Wed, 21 Oct 2024 07:28:00 GMT
```

---

## 📚 最佳实践

### 1. 使用 DTO

```go
// 请求 DTO
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

// 响应 DTO
type UserResponse struct {
    ID          int64     `json:"id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    DisplayName string    `json:"display_name"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### 2. 输入验证

```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 验证失败
        h.respHandler.Error(c, http.StatusBadRequest, 
            response.NewErrorResponse(CodeInvalidParams, "参数验证失败", err.Error()))
        return
    }
    
    // 业务逻辑
}
```

### 3. 统一响应

```go
type ResponseHandler struct {
    errorMapper *response.ErrorMapper
}

func (h *ResponseHandler) Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data":    data,
    })
}

func (h *ResponseHandler) Error(c *gin.Context, status int, err *response.ErrorResponse) {
    c.Status(status)
    c.JSON(status, err)
}
```

---

## 📖 参考资源

- [RESTful API 最佳实践](https://restfulapi.net/)
- [HTTP 状态码 RFC](https://tools.ietf.org/html/rfc7231)
- [JSON API 规范](https://jsonapi.org/)
- [OpenAPI Specification](https://swagger.io/specification/)

---

**最后更新：** 2024-03-23  
**维护者：** Go DDD Scaffold Team
