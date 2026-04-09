# 配置管理指南

## 📁 配置文件结构

```
backend/configs/
├── server.yaml          # 服务器配置（端口、超时、模式）
├── database.yaml        # 数据库配置（连接、池配置）
├── redis.yaml           # Redis 配置（地址、密码、池）
├── auth.yaml            # 认证配置（JWT、设备、Token）
├── rate_limit.yaml      # 限流配置（通用、登录、注册）
├── app.yaml             # 应用配置（日志、CORS、邮件、Asynq、监控）
├── prometheus.yml       # Prometheus 监控配置
├── .env                 # 环境变量覆盖（Git 忽略）
└── .env.example         # 环境变量模板（提交 Git）
```

## 🔧 配置优先级

```
环境变量 > .env 文件 > YAML 配置文件 > 代码默认值
```

**说明：**
- **YAML 配置文件**：包含通用配置，提交到 Git
- **.env 文件**：包含环境特定配置（如密码），不提交到 Git
- **环境变量**：最高优先级，适合 Docker/K8s 部署

## 📝 YAML 配置文件

### server.yaml - 服务器配置
```yaml
server:
  port: 8080
  mode: debug  # debug, release, test
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 60s
```

### database.yaml - 数据库配置
```yaml
database:
  host: localhost
  port: 5432
  name: ddd_scaffold
  user: postgres
  password: ""  # 从环境变量读取
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m
```

### redis.yaml - Redis 配置
```yaml
redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10
```

### auth.yaml - 认证配置
```yaml
# JWT 配置
jwt:
  secret: your-secret-key
  access_expire: 30m
  refresh_expire: 168h  # 7 天
  issuer: go-ddd-scaffold

# 设备管理配置
device:
  max_devices_per_user: 5
  auto_revoke_oldest: false

# Token 过期配置
token:
  email_verification_expire: 24h
  password_reset_expire: 1h
```

### rate_limit.yaml - 限流配置
```yaml
rate_limit:
  enabled: true
  general:
    rate: 60        # 每分钟 60 次
    burst: 100      # 突发 100 次
  login:
    rate: 5         # 每分钟 5 次
    burst: 10
  register:
    rate: 10        # 每分钟 10 次
    burst: 20
```

### app.yaml - 应用配置
```yaml
# 日志配置
logger:
  level: debug
  format: console  # json, console
  output_path: backend/logs/app.log

# CORS 配置
cors:
  allowed_origins: ["http://localhost:3000"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["Authorization", "Content-Type"]
  allow_credentials: true
  max_age: 3600

# 邮件配置
email:
  from: noreply@example.com
  verification_url_template: "http://localhost:3000/verify-email?token=%s&user_id=%s"
  password_reset_url_template: "http://localhost:3000/reset-password?token=%s&user_id=%s"

# Asynq 消息队列配置
asynq:
  addr: localhost:6379
  concurrency: 10
  queues:
    critical: 6
    default: 3
    low: 1

# 监控指标配置
metrics:
  enabled: true
  http:
    enabled: true
  database:
    enabled: true
  redis:
    enabled: true
```

## 🔐 环境变量覆盖

### 命名规则

环境变量格式：`APP_{配置路径}`（`.` 替换为 `_`，全部大写）

**示例映射：**

| YAML 配置键 | 环境变量 |
|------------|---------|
| `server.port` | `APP_SERVER_PORT` |
| `database.host` | `APP_DATABASE_HOST` |
| `database.password` | `APP_DATABASE_PASSWORD` |
| `device.max_devices_per_user` | `APP_DEVICE_MAX_DEVICES_PER_USER` |
| `jwt.secret` | `APP_JWT_SECRET` |

### .env 文件示例

```bash
# 复制模板
cp backend/configs/.env.example backend/configs/.env

# 编辑 .env 文件
vim backend/configs/.env
```

```bash
# 服务器配置
APP_SERVER_PORT=8080
APP_SERVER_MODE=debug

# 数据库配置
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_NAME=ddd_scaffold
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=your-secure-password

# Redis 配置
APP_REDIS_ADDR=localhost:6379
APP_REDIS_PASSWORD=your-redis-password

# JWT 配置（生产环境必须修改）
APP_JWT_SECRET=your-super-secret-key-change-in-production

# 日志配置
APP_LOGGER_LEVEL=info

# 监控配置
APP_METRICS_ENABLED=true
```

## 🚀 不同环境配置

### 本地开发

使用默认的 YAML 配置 + `.env` 覆盖敏感信息：

```bash
cd backend
cp configs/.env.example configs/.env
# 编辑 .env，设置本地数据库密码等
```

### Docker 部署

通过环境变量注入配置：

```bash
docker run -d \
  -e APP_DATABASE_HOST=prod-db.example.com \
  -e APP_DATABASE_PASSWORD=secure-password \
  -e APP_JWT_SECRET=production-secret \
  -e APP_SERVER_MODE=release \
  your-app:latest
```

### Kubernetes 部署

使用 ConfigMap 和 Secret：

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_SERVER_MODE: "release"
  APP_LOGGER_LEVEL: "info"
  APP_DATABASE_HOST: "prod-db.example.com"
---
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  APP_DATABASE_PASSWORD: c2VjdXJlLXBhc3N3b3Jk  # base64 encoded
  APP_JWT_SECRET: cHJvZHVjdGlvbi1zZWNyZXQ=
```

## ⚙️ 配置加载逻辑

配置加载顺序（`internal/infra/config/config.go`）：

1. **加载 YAML 配置文件**（使用独立 viper 实例）
   - server.yaml
   - database.yaml
   - redis.yaml
   - auth.yaml
   - rate_limit.yaml
   - app.yaml

2. **加载 .env 文件**（覆盖 YAML 配置）

3. **绑定环境变量**（优先级最高）
   - 设置前缀：`APP`
   - 键名映射：`.` → `_`

4. **反序列化到 Config 结构体**

## 🔍 调试配置

查看实际加载的配置值：

```go
cfg, _ := config.Load("development")
log.Printf("Database: %+v", cfg.Database)
log.Printf("JWT Secret: %s", cfg.JWT.Secret)
log.Printf("Device Max: %d", cfg.Device.MaxDevicesPerUser)
```

## 📋 最佳实践

### ✅ 推荐做法

1. **YAML 文件提交 Git**
   - 包含所有配置键和默认值
   - 敏感信息留空或使用占位符

2. **.env 文件不提交 Git**
   - 已在 `.gitignore` 中
   - 只包含环境特定的值

3. **生产环境使用环境变量**
   - Docker/K8s 通过 Secret 注入
   - 避免在配置文件中硬编码密码

4. **按功能拆分配置文件**
   - 职责清晰，易于维护
   - 减少 Git 冲突

### ❌ 避免做法

1. **不要在 YAML 中硬编码密码**
   ```yaml
   # ❌ 错误
   database:
     password: my-secret-password
   
   # ✅ 正确
   database:
     password: ""  # 从环境变量读取
   ```

2. **不要修改 .env.example 为真实值**
   ```bash
   # ❌ 错误
   APP_DATABASE_PASSWORD=real-password
   
   # ✅ 正确
   APP_DATABASE_PASSWORD=your-password  # 占位符
   ```

3. **不要在代码中硬编码配置**
   ```go
   // ❌ 错误
   dbHost := "localhost"
   
   // ✅ 正确
   dbHost := cfg.Database.Host
   ```

## 🔗 相关文档

- [快速开始指南](../development/GETTING_STARTED.md)
- [监控配置](../operations/MONITORING_SETUP.md)
- [Docker 部署](../deployment/DOCKER_DEPLOYMENT.md)
