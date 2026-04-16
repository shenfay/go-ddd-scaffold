# API 快速开始指南

本指南将帮助您在 5 分钟内快速上手 DDD Scaffold API。

## 📋 前置条件

- 已安装 Go 1.21+
- 已安装 PostgreSQL 16+
- 已安装 Redis 7+
- 已安装 Postman 或类似 API 测试工具

## 🚀 快速启动

### 1. 启动依赖服务

```bash
# 使用 Docker Compose 启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis
```

### 2. 运行数据库迁移

```bash
cd backend
make migrate-up
```

### 3. 启动 API 服务

```bash
# 方式一: 使用 Make
make run

# 方式二: 直接运行
go run cmd/api/main.go
```

API 服务将在 `http://localhost:8080` 启动。

### 4. 查看 Swagger 文档

访问 `http://localhost:8080/swagger/index.html` 查看完整的 API 文档。

## 📝 快速测试

### 使用 curl

#### 1. 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

#### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

响应示例:
```json
{
  "code": "SUCCESS",
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 3600,
    "user": {
      "id": "01H...",
      "email": "test@example.com",
      "email_verified": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

#### 3. 获取当前用户信息

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 使用 Postman

1. 导入 Postman Collection:
   - 打开 Postman
   - 点击 Import
   - 选择 `backend/api/postman-collection.json`

2. 设置环境变量:
   - `base_url`: `http://localhost:8080/api/v1`
   - `access_token`: (登录后自动设置)
   - `refresh_token`: (登录后自动设置)

3. 按顺序执行:
   - Register → Login → Get Current User

## 🔐 认证说明

### Token 类型

| Token 类型 | 有效期 | 用途 |
|-----------|--------|------|
| Access Token | 1 小时 | API 请求认证 |
| Refresh Token | 7 天 | 刷新 Access Token |

### 认证流程

```
注册 → 登录 → 获取 Token → 使用 Token 访问 API → Token 过期前刷新
```

### 刷新 Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## 📚 API 端点

### 认证相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/login` | 用户登录 | ❌ |
| POST | `/auth/logout` | 用户登出 | ✅ |
| POST | `/auth/refresh` | 刷新 Token | ❌ |
| GET | `/auth/me` | 当前用户 | ✅ |
| POST | `/auth/verify-email` | 验证邮箱 | ❌ |
| POST | `/auth/resend-verification` | 重发验证邮件 | ✅ |
| GET | `/auth/devices` | 设备列表 | ✅ |
| DELETE | `/auth/devices/:id` | 撤销设备 | ✅ |

### 用户管理

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/users/:id` | 获取用户 | ✅ |
| POST | `/users` | 创建用户 | ✅ |

### 健康检查

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/health` | 健康检查 | ❌ |

## ⚠️ 错误处理

### 错误响应格式

```json
{
  "code": "INVALID_ARGUMENT",
  "message": "Invalid email format",
  "timestamp": "2024-01-01T00:00:00Z",
  "trace_id": "abc123"
}
```

### 常见错误码

| 错误码 | HTTP 状态 | 说明 |
|--------|----------|------|
| `INVALID_ARGUMENT` | 400 | 参数错误 |
| `UNAUTHORIZED` | 401 | 未认证 |
| `FORBIDDEN` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `CONFLICT` | 409 | 资源冲突 |
| `INTERNAL_ERROR` | 500 | 服务器错误 |

## 🛡️ 安全要求

### 密码规则

- 最少 8 个字符
- 最多 72 个字符 (bcrypt 限制)
- 必须包含大写字母
- 必须包含小写字母
- 必须包含数字
- 必须包含特殊字符

### 速率限制

- 登录: 5 次失败后锁定账户
- API 请求: 根据配置限制

## 📖 更多资源

- **Swagger 文档**: `http://localhost:8080/swagger/index.html`
- **Postman Collection**: `backend/api/postman-collection.json`
- **架构文档**: `docs/architecture/`
- **开发指南**: `docs/development/`

## ❓ 常见问题

### Q: 如何重置密码?

A: 使用密码重置流程:
1. 请求密码重置 Token
2. 检查邮箱获取 Token
3. 使用 Token 重置密码

### Q: Token 过期了怎么办?

A: 使用 Refresh Token 获取新的 Token 对,或重新登录。

### Q: 如何查看 API 请求日志?

A: 检查应用日志文件 `logs/app.log` 或使用监控面板。

### Q: 如何修改配置?

A: 编辑 `backend/configs/` 目录下的配置文件,或设置环境变量。

## 🎯 下一步

- [ ] 阅读 [架构文档](../architecture/DDD_ARCHITECTURE.md)
- [ ] 了解 [领域模型](../architecture/DOMAIN_MODEL.md)
- [ ] 查看 [数据库设计](../database/SCHEMA_DESIGN.md)
- [ ] 学习 [开发规范](../development/DEVELOPMENT_GUIDE.md)
