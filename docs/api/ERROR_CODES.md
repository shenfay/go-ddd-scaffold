# API 错误码参考

本文档列出所有 API 错误码及其说明。

## 错误响应格式

所有错误响应都遵循以下格式:

```json
{
  "code": "ERROR_CODE",
  "message": "人类可读的错误描述",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "请求追踪ID"
}
```

## 通用错误码

| 错误码 | HTTP 状态码 | 说明 | 示例场景 |
|--------|------------|------|---------|
| `INVALID_ARGUMENT` | 400 | 请求参数无效 | 邮箱格式错误、密码不符合要求 |
| `UNAUTHENTICATED` | 401 | 未认证或认证失败 | Token 无效、Token 过期 |
| `PERMISSION_DENIED` | 403 | 权限不足 | 访问其他用户资源 |
| `NOT_FOUND` | 404 | 资源不存在 | 用户不存在 |
| `ALREADY_EXISTS` | 409 | 资源已存在 | 邮箱已注册 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 | 数据库连接失败 |

## 认证相关错误

### 注册

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `INVALID_ARGUMENT` | 邮箱格式无效 | 使用有效的邮箱格式 |
| `INVALID_ARGUMENT` | 密码不符合要求 | 密码需包含大小写字母、数字和特殊字符 |
| `ALREADY_EXISTS` | 邮箱已被注册 | 使用其他邮箱或登录 |

**示例**:
```json
{
  "code": "ALREADY_EXISTS",
  "message": "Email already exists",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

### 登录

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `INVALID_ARGUMENT` | 邮箱或密码错误 | 检查凭据是否正确 |
| `PERMISSION_DENIED` | 账户已锁定 | 联系管理员解锁 |
| `PERMISSION_DENIED` | 邮箱未验证 | 先验证邮箱 |

**示例**:
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "Invalid email or password",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

**安全说明**: 登录失败时不区分邮箱错误还是密码错误,防止枚举攻击。

### Token 刷新

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `INVALID_ARGUMENT` | Refresh Token 无效 | 重新登录 |
| `INVALID_ARGUMENT` | Refresh Token 已过期 | 重新登录 |

### 登出

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `UNAUTHENTICATED` | 未提供 Access Token | 在请求头中添加 Token |
| `UNAUTHENTICATED` | Access Token 无效 | 重新登录 |

### 邮箱验证

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `INVALID_ARGUMENT` | 验证 Token 无效或过期 | 请求重新发送验证邮件 |
| `ALREADY_EXISTS` | 邮箱已验证 | 无需重复验证 |

### 设备管理

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `NOT_FOUND` | 设备不存在 | 检查设备 ID 是否正确 |
| `PERMISSION_DENIED` | 无权操作该设备 | 只能管理自己的设备 |

## 用户管理错误

### 获取用户

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `NOT_FOUND` | 用户不存在 | 检查用户 ID |
| `PERMISSION_DENIED` | 无权查看该用户 | 只能查看自己的信息(普通用户) |

### 创建用户

| 错误码 | 说明 | 解决方案 |
|--------|------|---------|
| `INVALID_ARGUMENT` | 参数无效 | 检查请求参数 |
| `ALREADY_EXISTS` | 邮箱已存在 | 使用其他邮箱 |
| `PERMISSION_DENIED` | 无管理员权限 | 需要管理员角色 |

## 请求示例与响应

### 注册失败示例

**请求**:
```bash
POST /api/v1/auth/register
{
  "email": "invalid-email",
  "password": "weak"
}
```

**响应** (400):
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "Invalid email format",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

### 登录失败示例

**请求**:
```bash
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "WrongPassword"
}
```

**响应** (400):
```json
{
  "code": "INVALID_ARGUMENT",
  "message": "Invalid email or password",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

### Token 过期示例

**请求**:
```bash
GET /api/v1/auth/me
Authorization: Bearer expired_token
```

**响应** (401):
```json
{
  "code": "UNAUTHENTICATED",
  "message": "Token has expired",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

### 权限不足示例

**请求**:
```bash
GET /api/v1/users/other-user-id
Authorization: Bearer valid_token
```

**响应** (403):
```json
{
  "code": "PERMISSION_DENIED",
  "message": "You do not have permission to access this resource",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

## 错误处理最佳实践

### 客户端处理

```javascript
try {
  const response = await fetch('/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  
  const data = await response.json();
  
  if (!response.ok) {
    switch (data.code) {
      case 'INVALID_ARGUMENT':
        // 显示表单验证错误
        showFormError(data.message);
        break;
      case 'UNAUTHENTICATED':
        // 重定向到登录页
        redirectToLogin();
        break;
      case 'PERMISSION_DENIED':
        // 显示权限错误
        showPermissionError();
        break;
      default:
        // 显示通用错误
        showGenericError(data.message);
    }
  }
} catch (error) {
  // 网络错误
  showNetworkError();
}
```

### 重试策略

对于某些错误,可以实现重试:

- `INTERNAL_ERROR`: 可重试 (指数退避)
- `UNAUTHENTICATED`: 尝试刷新 Token 后重试
- `INVALID_ARGUMENT`: 不重试 (需要用户修正)

## 调试提示

1. **使用 trace_id**: 每个错误响应都包含 `trace_id`,可用于日志追踪
2. **检查 Swagger 文档**: 查看端点的详细错误说明
3. **查看应用日志**: 使用 `trace_id` 在日志中搜索详细信息
4. **使用 Postman**: 导入 Collection 快速测试各个端点

## 相关资源

- [快速开始指南](./QUICK_START.md)
- [Swagger 文档](http://localhost:8080/swagger/index.html)
- [Postman Collection](../../backend/api/postman-collection.json)
