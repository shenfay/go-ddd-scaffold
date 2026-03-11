# Go DDD Scaffold API设计规范文档

## 文档概述

本文档定义了 go-ddd-scaffold 项目的API设计规范，包括RESTful设计原则、接口命名规范、请求响应格式标准以及错误处理机制。

## RESTful API设计原则

### 核心设计哲学
遵循REST架构风格，以资源为中心，通过标准HTTP方法操作资源：

```
资源标识 → HTTP方法 → 操作语义
/users     GET        获取用户列表
/users     POST       创建新用户
/users/123 GET        获取特定用户
/users/123 PUT        更新用户信息
/users/123 DELETE     删除用户
```

### 资源层次结构
```
/api/v1/
├── /users              # 用户资源
│   ├── /{id}          # 单个用户
│   ├── /{id}/profile  # 用户档案
│   └── /{id}/tenants  # 用户所属租户
├── /tenants           # 租户资源
│   ├── /{id}          # 单个租户
│   ├── /{id}/members  # 租户成员
│   └── /{id}/roles    # 租户角色
├── /auth              # 认证资源
│   ├── /login         # 用户登录
│   ├── /logout        # 用户登出
│   └── /refresh       # 刷新令牌
└── /permissions       # 权限资源
    └── /{id}          # 单个权限
```

## 接口命名规范

### URI设计规范

**基本规则：**
- 使用名词复数形式表示资源集合
- 使用小写字母和连字符分隔单词
- 避免动词，使用HTTP方法表达操作
- 版本号放在URI开头

**良好示例：**
```
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/123
PUT    /api/v1/users/123
DELETE /api/v1/users/123/profile
```

**避免的示例：**
```
GET    /api/v1/getUsers          # 不要使用动词
POST   /api/v1/create-user       # 不要使用动作词汇
GET    /api/V1/Users             # 不要使用大写字母
GET    /api/v1/user_management   # 不要使用下划线
```

### HTTP方法语义

| 方法 | 语义 | 幂等性 | 安全性 |
|------|------|--------|--------|
| GET | 获取资源 | ✓ | ✓ |
| POST | 创建资源 | ✗ | ✗ |
| PUT | 更新完整资源 | ✓ | ✗ |
| PATCH | 部分更新资源 | ✗ | ✗ |
| DELETE | 删除资源 | ✓ | ✗ |
| HEAD | 获取资源元信息 | ✓ | ✓ |
| OPTIONS | 获取支持的方法 | ✓ | ✓ |

### 状态码规范

**成功状态码：**
```http
200 OK                    # 请求成功
201 Created               # 资源创建成功
202 Accepted              # 请求已接受，异步处理中
204 No Content            # 请求成功但无返回内容
```

**客户端错误：**
```http
400 Bad Request           # 请求参数错误
401 Unauthorized          # 未认证
403 Forbidden             # 权限不足
404 Not Found             # 资源不存在
405 Method Not Allowed    # HTTP方法不允许
409 Conflict              # 资源冲突
422 Unprocessable Entity  # 语义错误（如验证失败）
```

**服务器错误：**
```http
500 Internal Server Error # 服务器内部错误
502 Bad Gateway           # 网关错误
503 Service Unavailable   # 服务不可用
504 Gateway Timeout       # 网关超时
```

## 请求响应格式标准

### 请求格式规范

**JSON请求体结构：**
```json
{
  "field1": "value1",
  "field2": "value2",
  "nested_object": {
    "sub_field": "sub_value"
  },
  "array_field": ["item1", "item2"]
}
```

**查询参数规范：**
```
GET /api/v1/users?page=1&limit=20&sort=-created_at&status=active
```

**分页参数：**
- `page`: 页码（从1开始）
- `limit`: 每页条数
- `sort`: 排序字段（负号表示降序）

### 响应格式规范

**统一响应结构：**
```go
type APIResponse struct {
    Code    int         `json:"code"`     // 业务状态码
    Message string      `json:"message"`  // 响应消息
    Data    interface{} `json:"data"`     // 响应数据
    Meta    *MetaInfo   `json:"meta"`     // 元信息（分页等）
}

type MetaInfo struct {
    Page       int `json:"page,omitempty"`
    Limit      int `json:"limit,omitempty"`
    Total      int `json:"total,omitempty"`
    TotalPages int `json:"total_pages,omitempty"`
}
```

**成功响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1234567890123456789,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**分页响应示例：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1234567890123456789,
      "username": "user1",
      "email": "user1@example.com"
    },
    {
      "id": 1234567890123456790,
      "username": "user2",
      "email": "user2@example.com"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

## 核心API接口设计

### 1. 用户管理接口

#### 创建用户
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "profile": {
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890"
  }
}
```

**成功响应：**
```http
HTTP/1.1 201 Created
Content-Type: application/json

{
  "code": 0,
  "message": "User created successfully",
  "data": {
    "id": 1234567890123456789,
    "username": "john_doe",
    "email": "john@example.com",
    "status": 1,
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 获取用户列表
```http
GET /api/v1/users?page=1&limit=20&status=active
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1234567890123456789,
      "username": "john_doe",
      "email": "john@example.com",
      "status": 1,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

#### 获取用户详情
```http
GET /api/v1/users/1234567890123456789
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1234567890123456789,
    "username": "john_doe",
    "email": "john@example.com",
    "profile": {
      "first_name": "John",
      "last_name": "Doe",
      "phone": "+1234567890"
    },
    "status": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

#### 更新用户
```http
PUT /api/v1/users/1234567890123456789
Content-Type: application/json

{
  "email": "newemail@example.com",
  "profile": {
    "first_name": "John",
    "last_name": "Smith"
  }
}
```

#### 删除用户
```http
DELETE /api/v1/users/1234567890123456789
```

### 2. 认证接口

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
    "expires_in": 1800,
    "token_type": "Bearer",
    "user": {
      "id": 1234567890123456789,
      "username": "john_doe",
      "email": "john@example.com"
    }
  }
}
```

#### 刷新令牌
```http
POST /api/v1/auth/refresh
Authorization: Bearer <refresh_token>
```

#### 用户登出
```http
POST /api/v1/auth/logout
Authorization: Bearer <access_token>
```

### 3. 租户管理接口

#### 创建租户
```http
POST /api/v1/tenants
Content-Type: application/json

{
  "code": "COMPANY-A",
  "name": "Company A",
  "config": {
    "max_members": 100,
    "storage_limit": "100GB"
  }
}
```

#### 获取租户成员
```http
GET /api/v1/tenants/1234567890123456789/members
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "user_id": 1234567890123456789,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "admin",
      "joined_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### 4. 权限管理接口

#### 获取用户权限
```http
GET /api/v1/users/1234567890123456789/permissions
Authorization: Bearer <access_token>
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "resource": "user",
      "action": "read",
      "description": "查看用户信息"
    },
    {
      "resource": "user",
      "action": "write",
      "description": "修改用户信息"
    }
  ]
}
```

## 错误处理规范

### 统一错误响应格式
```go
type ErrorResponse struct {
    Code    int         `json:"code"`     // 业务错误码
    Message string      `json:"message"`  // 错误消息
    Details interface{} `json:"details"`  // 详细错误信息
    TraceID string      `json:"trace_id"` // 请求追踪ID
}
```

### 错误码定义
```go
const (
    // 通用错误码
    SuccessCode           = 0
    UnknownErrorCode      = 1000
    InvalidParamCode      = 1001
    UnauthorizedCode      = 1002
    ForbiddenCode         = 1003
    NotFoundCode          = 1004
    ConflictCode          = 1005
    
    // 业务错误码
    UserNotFoundCode      = 2001
    UserAlreadyExistsCode = 2002
    InvalidPasswordCode   = 2003
    AccountLockedCode     = 2004
    
    // 系统错误码
    DatabaseErrorCode     = 5001
    CacheErrorCode        = 5002
    ExternalServiceCode   = 5003
)
```

### 错误响应示例

**参数验证错误：**
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "code": 1001,
  "message": "Invalid request parameters",
  "details": {
    "username": "用户名长度必须在3-20个字符之间",
    "email": "邮箱格式不正确",
    "password": "密码必须包含大小写字母和数字"
  },
  "trace_id": "abc123def456"
}
```

**认证失败：**
```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "code": 1002,
  "message": "Authentication failed",
  "details": "Invalid username or password",
  "trace_id": "xyz789uvw012"
}
```

**权限不足：**
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "code": 1003,
  "message": "Insufficient permissions",
  "details": "You don't have permission to perform this action",
  "trace_id": "def456ghi789"
}
```

## 安全规范

### 认证授权
- 所有API接口（除登录外）都需要JWT认证
- 敏感操作需要额外的权限验证
- 实施合理的令牌过期和刷新机制

### 输入验证
```go
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=128"`
    Profile  UserProfile `json:"profile" validate:"required"`
}

type UserProfile struct {
    FirstName string `json:"first_name" validate:"required,max=50"`
    LastName  string `json:"last_name" validate:"required,max=50"`
    Phone     string `json:"phone" validate:"omitempty,e164"`
}
```

### 速率限制
- 普通接口：1000次/小时/IP
- 登录接口：10次/分钟/IP
- 敏感接口：100次/小时/IP

### CORS配置
```yaml
cors:
  allowed_origins:
    - "https://app.example.com"
    - "https://admin.example.com"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Origin"
    - "Content-Type"
    - "Accept"
    - "Authorization"
  exposed_headers:
    - "Content-Length"
    - "Content-Type"
  max_age: 3600
```

这个API设计规范文档为项目提供了统一的接口设计标准和实现指南。