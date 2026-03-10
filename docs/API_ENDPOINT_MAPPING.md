# 前后端 API 接口映射表

> 基于 Swagger 文档生成的前后端对接指南  
> 最后更新：2026-03-10

## 认证模块 (Auth)

### POST /api/auth/register - 用户注册

**后端 Handler**: `AuthHandler.Register`  
**前端 Service**: `userService.register(userData)`  
**请求参数**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "nickname": "User Nickname",
  "role": "member",
  "tenantId": "uuid-string"
}
```

**响应格式**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "nickname": "User Nickname",
    "status": "ACTIVE",
    "createdAt": "2026-03-10T00:00:00Z"
  },
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 400 Bad Request - 参数验证失败（邮箱格式/密码长度）
- 409 Conflict - 用户已存在
- 500 Internal Server Error - 服务器内部错误

---

### POST /api/auth/login- 用户登录

**后端 Handler**: `AuthHandler.Login`  
**前端 Service**: `userService.login(email, password)`  
**请求参数**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应格式**:
```json
{
  "code": "Success",
  "message": "登录成功",
  "data": {
    "accessToken": "jwt-token-string",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "nickname": "User Nickname",
      "status": "ACTIVE"
    }
  },
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 400 Bad Request - 参数验证失败
- 401 Unauthorized - 用户名或密码错误
- 500 Internal Server Error - 服务器内部错误

---

### POST /api/auth/logout - 用户登出

**后端 Handler**: `AuthHandler.Logout`  
**前端 Service**: `userService.logout()`  
**请求头**: `Authorization: Bearer {token}`  
**响应格式**:
```json
{
  "code": "Success",
  "message": "登出成功",
  "data": null,
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 401 Unauthorized - 未授权
- 500 Internal Server Error - 服务器内部错误

---

## 用户管理模块 (User Management)

### GET /api/users/:id - 获取用户信息

**后端 Handler**: `UserHandler.GetUser`  
**前端 Service**: `userService.getUser(userId)`  
**路径参数**: `id` - 用户 UUID  
**响应格式**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "nickname": "User Nickname",
    "phone": "+86-13800138000",
    "bio": "个人简介",
    "avatar": "https://...",
    "status": "ACTIVE",
    "createdAt": "2026-03-10T00:00:00Z",
    "updatedAt": "2026-03-10T00:00:00Z"
  },
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 400 Bad Request - 无效的用户 ID 格式
- 404 Not Found- 用户不存在
- 500 Internal Server Error - 服务器内部错误

---

### PUT /api/users/:id - 更新用户信息

**后端 Handler**: `UserHandler.UpdateUser`  
**前端 Service**: `userService.updateUser(userId, userData)`  
**路径参数**: `id` - 用户 UUID  
**请求参数**:
```json
{
  "email": "newemail@example.com",
  "password": "newpassword123",
  "status": "ACTIVE"
}
```

**响应格式**:
```json
{
  "code": "NoContent",
  "message": "无内容",
  "data": null,
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**HTTP Status**: 204 No Content

**错误处理**:
- 400 Bad Request - 参数验证失败
- 404 Not Found - 用户不存在
- 500 Internal Server Error - 服务器内部错误

---

## 个人资料模块 (Profile)

### GET /api/users/info - 获取当前用户信息

**后端 Handler**: `ProfileHandler.GetUserInfo`  
**前端 Service**: `userService.getProfile()`  
**请求头**: `Authorization: Bearer {token}`  
**响应格式**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "nickname": "User Nickname",
    "phone": "+86-13800138000",
    "bio": "个人简介",
    "avatar": "https://...",
    "status": "ACTIVE",
    "createdAt": "2026-03-10T00:00:00Z",
    "updatedAt": "2026-03-10T00:00:00Z"
  },
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 401 Unauthorized - 未授权
- 500 Internal Server Error - 服务器内部错误

---

### PUT /api/users/profile- 更新个人资料

**后端 Handler**: `ProfileHandler.UpdateProfile`  
**前端 Service**: `userService.updateProfile(profileData)`  
**请求头**: `Authorization: Bearer {token}`  
**请求参数**:
```json
{
  "nickname": "New Nickname",
  "phone": "+86-13800138000",
  "bio": "Updated bio"
}
```

**响应格式**:
```json
{
  "code": "NoContent",
  "message": "无内容",
  "data": null,
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**HTTP Status**: 204 No Content

**错误处理**:
- 400 Bad Request - 参数验证失败
- 404 Not Found- 用户不存在
- 500 Internal Server Error - 服务器内部错误

---

## 租户管理模块 (Tenant)

### POST /api/tenants - 创建租户

**后端 Handler**: `TenantHandler.CreateTenant`  
**前端 Service**: `userService.createTenant(tenantData)`  
**请求头**: `Authorization: Bearer {token}`  
**请求参数**:
```json
{
  "name": "My Tenant",
  "description": "租户描述",
  "maxMembers": 100
}
```

**响应格式**:
```json
{
  "code": "Created",
  "message": "创建成功",
  "data": {
    "id": "uuid",
    "name": "My Tenant",
    "description": "租户描述",
    "maxMembers": 100,
    "memberCount": 1,
    "createdAt": "2026-03-10T00:00:00Z",
    "updatedAt": "2026-03-10T00:00:00Z"
  },
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 400 Bad Request - 参数验证失败
- 401 Unauthorized - 未授权
- 500 Internal Server Error - 服务器内部错误

---

### GET /api/tenants/my-tenants- 获取用户的租户列表

**后端 Handler**: `TenantHandler.GetUserTenants`  
**前端 Service**: `userService.getUserTenants()`  
**请求头**: `Authorization: Bearer {token}`  
**响应格式**:
```json
{
  "code": "Success",
  "message": "成功",
  "data": [
    {
      "id": "uuid",
      "name": "Tenant 1",
      "description": "描述 1",
      "maxMembers": 100,
      "memberCount": 5,
      "role": "owner",
      "createdAt": "2026-03-10T00:00:00Z",
      "updatedAt": "2026-03-10T00:00:00Z"
    }
  ],
  "requestId": "req-uuid",
  "timestamp": "2026-03-10T00:00:00Z"
}
```

**错误处理**:
- 401 Unauthorized - 未授权
- 500 Internal Server Error - 服务器内部错误

---

## 统一响应格式说明

### 成功响应结构

```typescript
interface SuccessResponse<T> {
  code: string;        // "Success" | "Created" | "NoContent"
  message: string;     // 成功消息
  data: T;            // 数据对象
  requestId: string;   // 请求 ID（用于追踪）
  timestamp: string;   // ISO 8601 时间戳
}
```

### 错误响应结构

```typescript
interface ErrorResponse {
  code: string;       // 错误码
  message: string;    // 错误消息
  error?: {
   code: string;     // 详细错误码
   message: string;  // 详细错误描述
    details?: any;    // 错误详情
  };
  requestId: string;  // 请求 ID
  timestamp: string;  // 时间戳
}
```

### HTTP 状态码语义

| 状态码 | 含义 | 使用场景 |
|--------|------|---------|
| 200 OK | 成功 | GET/POST 查询操作 |
| 201 Created | 创建成功 | POST 创建资源 |
| 204 No Content | 无内容 | PUT 更新操作 |
| 400 Bad Request | 请求参数错误 | 验证失败 |
| 401 Unauthorized | 未授权 | Token 缺失/过期 |
| 403 Forbidden | 禁止访问 | 权限不足 |
| 404 Not Found| 资源不存在 | 找不到数据 |
| 409 Conflict | 冲突 | 资源已存在 |
| 500 Internal Server Error | 服务器内部错误 | 系统异常 |

---

## 前端集成清单

### ✅ 已完成

- [x] API 端点配置 (`endpoints.js`)
- [x] HTTP 客户端封装 (`client.js`)
- [x] 用户服务封装 (`userService.js`)
- [x] 请求拦截器（Token 注入）
- [x] 响应拦截器（统一处理）
- [x] 错误处理器（统一提示）

### 📋 待完成

- [ ] 登录页面调用新接口
- [ ] 注册页面调用新接口
- [ ] 个人中心页面调用 Profile API
- [ ] 租户管理页面调用 Tenant API
- [ ] 全局路由守卫（Token 验证）
- [ ] Token 过期自动刷新（可选）

---

## 调试工具

### 查看 Swagger 文档

```bash
cd backend/docs
cat swagger.yaml
```

### 测试 API 端点

```bash
# 登录测试
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 获取用户信息
curl -X GET http://localhost:8080/api/users/info \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 前端调试模式

在 `.env.development` 中配置：
```bash
REACT_APP_API_BASE_URL=http://localhost:8080/api
REACT_APP_DEBUG=true
```

---

## 常见问题 FAQ

### Q: 401 Unauthorized 如何处理？
A: 检查 Token 是否有效，可在 localStorage 中查看 `access_token`

### Q: 403 Forbidden 是什么意思？
A: 权限不足，检查用户角色和租户权限配置

### Q: 如何查看详细的错误信息？
A: 响应体中的 `error.details` 字段包含详细信息

### Q: Token 有效期多久？
A: 默认 72 小时，可在后端配置调整

### Q: 跨域问题如何解决？
A: 后端已配置 CORS，开发环境使用 proxy 配置
