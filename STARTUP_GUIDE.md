# 服务启动指南

## 快速启动

### 方式一：分别启动（推荐）

#### 1. 启动后端服务

```bash
cd backend
go run cmd/server/main.go
```

**访问地址**:
- API 服务：http://localhost:8080
- Prometheus 指标：http://localhost:8080/metrics
- 健康检查：http://localhost:8080/health

#### 2. 启动前端服务

```bash
cd frontend
npm start
```

**访问地址**:
- 前端应用：http://localhost:3000

---

### 方式二：使用 Docker Compose（生产环境）

```bash
cd backend/deployments/docker
docker-compose up -d
```

**访问地址**:
- 前端应用：http://localhost:3000
- 后端 API: http://localhost:8080
- Grafana: http://localhost:3001
- Prometheus: http://localhost:9090

---

## 服务状态检查

### 后端健康检查

```bash
# 简单健康检查
curl http://localhost:8080/health

# 详细健康检查
curl http://localhost:8080/health/detail

# 查看 Prometheus 指标
curl http://localhost:8080/metrics
```

### 前端服务检查

打开浏览器访问：http://localhost:3000

---

## API 测试示例

### 1. 用户注册

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "nickname": "Test User",
    "tenantId": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

### 2. 用户登录

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 3. 获取用户信息

```bash
curl -X GET http://localhost:8080/api/users/info \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 4. 创建租户

```bash
curl -X POST http://localhost:8080/api/tenants \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Tenant",
    "description": "我的租户描述",
    "maxMembers": 100
  }'
```

---

## 环境变量配置

### 后端环境变量

创建 `backend/.env` 文件：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go_ddd_scaffold

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 配置
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRE_HOURS=72

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug  # debug, release, test
```

### 前端环境变量

创建 `frontend/.env.development` 文件：

```bash
# API 基础 URL
REACT_APP_API_BASE_URL=http://localhost:8080/api

# 调试模式
REACT_APP_DEBUG=true

# 应用版本
REACT_APP_VERSION=0.1.0
```

---

## 常见问题排查

### 后端启动失败

#### 问题：数据库连接失败

```bash
# 检查 MySQL 是否运行
mysql.server status

# 或查看 Docker 容器
docker ps | grep mysql
```

**解决方案**:
1. 确保 MySQL 服务正在运行
2. 检查数据库配置是否正确
3. 确认数据库已创建

#### 问题：Redis 连接失败

```bash
# 检查 Redis 是否运行
redis-cli ping
```

**解决方案**:
1. 启动 Redis 服务
2. 检查 Redis 配置

### 前端启动失败

#### 问题：端口被占用

```bash
# 查看端口占用
lsof -i :3000

# 杀死进程
kill -9 <PID>
```

#### 问题：API 代理失败

检查 `frontend/package.json` 中的 proxy 配置：
```json
{
  "proxy": "http://localhost:8080"
}
```

---

## 开发工具推荐

### API调试工具

- **Postman**: API 接口调试
- **Insomnia**: 轻量级 API 工具
- **curl**: 命令行调试

### 浏览器开发工具

- Chrome DevTools
- React Developer Tools
- Redux DevTools

### 监控工具

- **Grafana**: 可视化监控仪表盘
- **Prometheus**: 指标采集和告警
- **AlertManager**: 告警通知管理

---

## 性能优化建议

### 后端优化

1. 使用 GIN_MODE=release 生产模式
2. 启用数据库连接池
3. 配置 Redis 缓存
4. 启用 Prometheus 监控

### 前端优化

1. 生产环境使用 npm run build
2. 启用 CDN 加速
3. 开启 Gzip 压缩
4. 使用浏览器缓存

---

## 安全注意事项

### 生产环境配置

1. **修改默认密钥**
   - JWT_SECRET
   - 数据库密码
   - Redis 密码

2. **启用 HTTPS**
   - SSL 证书配置
   - 强制 HTTPS 跳转

3. **CORS 配置**
   - 限制允许的域名
   - 配置跨域凭证

4. **速率限制**
   - API 请求限流
   - 防止暴力破解

---

## 日志查看

### 后端日志

```bash
# 实时查看日志
tail-f backend/logs/app.log

# 查看错误日志
tail-f backend/logs/error.log
```

### 前端日志

在浏览器控制台中查看：
- Console 面板
- Network 面板（API 请求）

---

## 下一步

服务启动成功后，可以：

1. ✅ 访问 http://localhost:3000 查看前端界面
2. ✅ 测试登录/注册功能
3. ✅ 创建和管理租户
4. ✅ 查看个人资料
5. ✅ 使用 Postman 测试 API 接口
6. ✅ 查看监控仪表盘（Grafana）

---

**祝开发愉快！** 🚀
