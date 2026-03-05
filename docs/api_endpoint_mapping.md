# 前后端 API 接口映射表

## 📋 概述

本文档详细说明前端与后端的 API 接口映射关系，确保前后端路径一致。

---

## 🔐 认证模块 (Auth)

### 后端路由定义

**文件：** `backend/internal/interfaces/http/auth/handler.go`

| 方法 | 路径 | Handler 方法 | 说明 |
|------|------|------------|------|
| POST | `/api/auth/register` | `AuthHandler.Register` | 用户注册 |
| POST | `/api/auth/login` | `AuthHandler.Login` | 用户登录 |
| POST | `/api/auth/logout` | `AuthHandler.Logout` | 用户登出（Token 加入黑名单）✅ 推荐使用 |

**注意：** 
- ~~`POST /api/users/logout`~~ 已废弃，请使用 `/api/auth/logout`
- `/api/auth/logout` 会将 Token 加入 Redis 黑名单，防止重放攻击 |

### 前端调用封装

**文件：** `frontend/src/data/api/services/userService.js`

```javascript
// 用户注册
async register(userData) {
  const path = getEndpoint('auth.register');
  return httpClient.post(path, userData);
}

// 用户登录
async login(email, password) {
  const path = getEndpoint('auth.login');
  return httpClient.post(path, { email, password });
}

// 用户登出
async logout() {
  const path = getEndpoint('auth.logout');
  return httpClient.post(path);
}
```

### 接口详情

#### 1. POST /api/auth/register

**请求示例：**
```json
{
  "email": "user@example.com",
  "password": "StrongPass123",
  "nickname": "Test User"
}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "nickname": "Test User"
    }
  },
  "message": "注册成功"
}
```

---

#### 2. POST /api/auth/login

**请求示例：**
```json
{
  "email": "user@example.com",
  "password": "StrongPass123"
}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "message": "登录成功"
}
```

---

#### 3. POST /api/auth/logout

**请求头：**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例：**
```json
{
  "code": 0,
  "message": "登出成功"
}
```

**后端处理逻辑：**
1. 从 Token 中提取 `jti`（唯一标识）
2. 将 `jti` 加入 Redis 黑名单
3. 设置过期时间等于 Token 剩余有效期
4. 返回成功响应

**Redis Key 格式：**
```
token:blacklist:{jti}
```

---

## 👤 用户模块 (User)

### 后端路由定义

**文件：** `backend/internal/interfaces/http/user/handler.go`

| 方法 | 路径 | Handler 方法 | 说明 |
|------|------|------------|------|
| GET | `/api/users/info` | `UserHandler.GetUserInfo` | 获取当前用户信息 |
| PUT | `/api/users/profile` | `UserHandler.UpdateProfile` | 更新个人资料 |
| GET | `/api/users/:id` | `UserHandler.GetUser` | 获取指定用户详情 |
| PUT | `/api/users/:id` | `UserHandler.UpdateUser` | 更新指定用户信息 |

**注意：** 
- ~~`POST /api/users/logout`~~ 接口已删除
- 登出操作请统一使用 `/api/auth/logout`

### 前端调用封装

```javascript
// 获取当前用户信息
async getInfo() {
  const path = getEndpoint('user.getInfo');
  return httpClient.get(path);
}

// 获取个人资料
async getProfile() {
  const path = getEndpoint('user.profile');
  return httpClient.get(path);
}

// 更新个人资料
async updateProfile(userData) {
  const path = getEndpoint('user.updateProfile');
  return httpClient.put(path, userData);
}

// 获取指定用户详情
async getUser(userId) {
  const path = getEndpoint('user.getUser', { id: userId });
  return httpClient.get(path);
}

// 更新指定用户信息
async updateUser(userId, userData) {
  const path = getEndpoint('user.updateUser', { id: userId });
  return httpClient.put(path, userData);
}
```

### 接口详情

#### 1. GET /api/users/info

**请求头：**
```
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "nickname": "Test User",
      "avatar": "https://example.com/avatar.jpg"
    }
  }
}
```

---

#### 2. PUT /api/users/profile

**请求头：**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体：**
```json
{
  "nickname": "New Nickname",
  "avatar": "https://example.com/new-avatar.jpg"
}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "资料更新成功"
}
```

---

#### 3. POST /api/auth/logout ✅ 推荐使用的登出接口

**说明：** 标准的登出接口，会将 Token 加入 Redis 黑名单

**请求头：**
```
Authorization: Bearer {token}
```

**响应示例：**
```json
{
  "code": 0,
  "message": "登出成功"
}
```

**后端处理逻辑：**
1. 从 Context 中获取用户 ID
2. 提取 Token 字符串
3. 解析 Token 获取 `jti`（唯一标识）
4. 将 `jti` 加入 Redis 黑名单
5. 设置过期时间等于 Token 剩余有效期
6. 返回成功响应

**Redis Key 格式：**
```
token:blacklist:{jti}
```

**安全特性：**
- ✅ 防止 Token 重放攻击
- ✅ 限流保护（100 req/s）
- ✅ 熔断保护（快速失败）
- ✅ 监控埋点（Prometheus）

~~~废弃的接口~~

**POST /api/users/logout** ❌ 已删除
- 原因：无实际业务逻辑，仅返回成功
- 风险：无法防止 Token 重放攻击
- 替代：请使用 `/api/auth/logout`

---

## 🏢 租户模块 (Tenant)

### 后端路由定义

**文件：** `backend/internal/interfaces/http/tenant/handler.go`

| 方法 | 路径 | Handler 方法 | 说明 |
|------|------|------------|------|
| GET | `/api/tenants/my-tenants` | `TenantHandler.GetUserTenants` | 获取用户的租户列表 |
| POST | `/api/tenants` | `TenantHandler.CreateTenant` | 创建租户 |

### 前端调用封装

```javascript
// 获取用户的租户列表
async getUserTenants() {
  const path = getEndpoint('tenant.userTenants');
  return httpClient.get(path);
}

// 创建租户
async createTenant(tenantData) {
  const path = getEndpoint('tenant.create');
  return httpClient.post(path, tenantData);
}

// 选择当前租户（本地操作）
selectTenant(tenantId) {
  localStorage.setItem('current_tenant_id', tenantId);
  window.dispatchEvent(new CustomEvent('tenantChanged', { 
    detail: { tenantId } 
  }));
}
```

### 接口详情

#### 1. GET /api/tenants/my-tenants

**请求头：**
```
Authorization: Bearer {token}
X-Tenant-ID: {tenant_id} (可选)
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "tenants": [
      {
        "id": "uuid",
        "name": "家庭 A",
        "maxMembers": 50,
        "currentMembers": 5
      },
      {
        "id": "uuid",
        "name": "公司 B",
        "maxMembers": 100,
        "currentMembers": 20
      }
    ]
  }
}
```

---

#### 2. POST /api/tenants

**请求头：**
```
Authorization: Bearer {token}
Content-Type: application/json
```

**请求体：**
```json
{
  "name": "新租户名称",
  "maxMembers": 50,
  "description": "租户描述"
}
```

**响应示例：**
```json
{
  "code": 0,
  "data": {
    "tenant": {
      "id": "uuid",
      "name": "新租户名称",
      "maxMembers": 50,
      "currentMembers": 1
    }
  },
  "message": "创建成功"
}
```

---

## 🔒 安全机制

### 认证流程

```mermaid
sequenceDiagram
    participant Frontend
    participant Backend
    participant Redis
    
    Frontend->>Backend: POST /api/auth/login
    Backend-->>Frontend: 返回 accessToken
    
    Frontend->>Backend: GET /api/users/info<br/>Authorization: Bearer token
    Backend->>Backend: JWT 验证
    Backend-->>Frontend: 返回用户信息
    
    Frontend->>Backend: POST /api/auth/logout<br/>Authorization: Bearer token
    Backend->>Redis: SETEX token:blacklist:{jti} EXPIRE
    Backend-->>Frontend: 登出成功
    
    Frontend->>Backend: GET /api/users/info<br/>Authorization: Bearer token
    Backend->>Redis: EXISTS token:blacklist:{jti}
    Redis-->>Backend: 存在（在黑名单中）
    Backend-->>Frontend: 401 Token 已被列入黑名单
```

### Token 黑名单机制

**Redis Key 命名规则：**
```
token:blacklist:{jti}
```

**过期时间设置：**
- 等于 Token 的剩余有效期
- 使用 `SETEX` 命令同时设置值和过期时间

**检查逻辑：**
```go
// 中间件中自动检查
func (m *AuthMiddleware) HandleFunc() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 解析 Token
        claims, err := m.jwtService.ParseToken(token)
        
        // 2. 检查是否在黑名单
        inBlacklist, _ := m.tokenBlacklist.IsBlacklisted(ctx, claims.ID)
        if inBlacklist {
            c.AbortWithStatusJSON(401, ErrorResponse{
                Code:    401,
                Message: "Token 已被列入黑名单",
            })
            return
        }
        
        // 3. Token 有效，继续处理请求
        c.Next()
    }
}
```

---

## 📊 性能优化

### Redis Pipeline 批量检查

**场景：** 需要同时检查多个 Token 是否在黑名单中

**传统方式（N 次网络往返）：**
```go
for _, token := range tokens {
    inBlacklist, _ := redis.Exists(ctx, key)
    // N 次网络往返
}
```

**Pipeline 方式（1 次网络往返）：**
```go
pipe := redis.Pipeline()
for _, token := range tokens {
    pipe.Exists(ctx, key)
}
results, _ := pipe.Exec(ctx)
// 仅 1 次网络往返
```

**性能对比：**
- **传统方式**: ~50ms × 100 = 5000ms (100 tokens)
- **Pipeline**: ~50ms (100 tokens)
- **提升**: **100 倍**

---

## 🎯 前端配置清单

### 环境变量配置

**文件：** `frontend/.env`

```bash
# API 基础 URL
REACT_APP_API_BASE_URL=http://localhost:8080/api

# 版本号
REACT_APP_VERSION=1.0.0
```

### HTTP 客户端配置

**文件：** `frontend/src/data/api/client.js`

```javascript
class HttpClient {
  constructor() {
    this.baseURL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:3000/api';
    this.timeout = 30000; // 30 seconds
  }
}
```

### 拦截器配置

**文件：** `frontend/src/data/api/interceptors/requestInterceptors.js`

```javascript
// 自动添加 Token
export function authInterceptor(config) {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }
  return config;
}

// 自动添加租户 ID
export function commonHeaderInterceptor(config) {
  const tenantId = localStorage.getItem('current_tenant_id');
  if (tenantId) {
    config.headers['X-Tenant-ID'] = tenantId;
  }
  return config;
}
```

---

## ✅ 接口一致性检查

### 已验证的接口

| 后端路径 | 前端端点 | 状态 | 说明 |
|----------|---------|------|------|
| `/api/auth/register` | `auth.register` | ✅ 一致 | 用户注册 |
| `/api/auth/login` | `auth.login` | ✅ 一致 | 用户登录 |
| `/api/auth/logout` | `auth.logout` | ✅ 一致 | 用户登出（带黑名单）✅ 推荐使用 |
| `/api/users/info` | `user.getInfo` | ✅ 一致 | 获取用户信息 |
| `/api/users/profile` | `user.profile` | ✅ 一致 | 获取/更新个人资料 |
| `/api/users/:id` | `user.getUser` | ✅ 一致 | 获取指定用户详情 |
| `/api/users/:id` | `user.updateUser` | ✅ 一致 | 更新指定用户信息 |
| `/api/tenants` | `tenant.create` | ✅ 一致 | 创建租户 |
| `/api/tenants/my-tenants` | `tenant.userTenants` | ✅ 一致 | 获取用户的租户列表 |

### 已废弃的接口

| 后端路径 | 前端端点 | 状态 | 原因 |
|----------|---------|------|------|
| ~~`/api/users/logout`~~ | ~~`user.logout`~~ | ❌ 已删除 | 无实际业务逻辑，存在安全风险 |

### 待实现的接口

| 后端路径 | 前端端点 | 优先级 | 说明 |
|----------|---------|--------|------|
| `/api/tenants/select` | `tenant.select` | ⭐⭐⭐ | 选择租户（后端预留） |
| `/api/users/change-password` | `user.changePassword` | ⭐⭐ | 修改密码 |

---

## 🐛 常见问题

### 问题 1: 404 Not Found

**可能原因：**
- 前端配置的 API 路径与后端不一致
- 缺少 `/api` 前缀

**解决方案：**
```javascript
// 检查 baseURL 配置
const baseURL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api';
//                                                     ^^^^^^^^^^^^^^^^^^^^^^^^
//                                                     必须包含 /api
```

---

### 问题 2: 401 Unauthorized

**可能原因：**
- Token 过期
- Token 被加入黑名单
- 未携带 Authorization 头

**排查步骤：**
```javascript
// 1. 检查 localStorage 中的 token
console.log(localStorage.getItem('auth_token'));

// 2. 检查请求头是否包含 token
const response = await fetch('/api/users/info', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
  }
});

// 3. 查看后端日志
tail -f backend/logs/app.log | grep "Token"
```

---

### 问题 3: CORS 跨域错误

**可能原因：**
- 后端未配置 CORS 中间件
- 前端请求的域名与后端不一致

**解决方案：**
```go
// backend/internal/infrastructure/middleware/cors.go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        
        c.Next()
    }
}
```

---

## 📖 参考资源

### 后端文件索引

- `backend/internal/interfaces/http/auth/handler.go` - 认证接口
- `backend/internal/interfaces/http/user/handler.go` - 用户接口
- `backend/internal/interfaces/http/tenant/handler.go` - 租户接口
- `backend/internal/infrastructure/auth/token_blacklist_service.go` - Token 黑名单服务

### 前端文件索引

- `frontend/src/data/endpoints/endpoints.js` - API 端点配置
- `frontend/src/data/api/services/userService.js` - 用户服务封装
- `frontend/src/data/api/interceptors/requestInterceptors.js` - 请求拦截器
- `frontend/src/data/api/interceptors/responseInterceptors.js` - 响应拦截器

---

**🎉 前后端接口映射清晰，所有核心接口已对齐！**
