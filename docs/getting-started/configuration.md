# 配置说明

## 📁 配置文件结构

```yaml
config/
└── config.yaml          # 主配置文件
```

---

## 🔧 完整配置示例

```yaml
# 应用配置
app:
  name: "ddd-scaffold"
  port: 8080
  env: "development"  # development, staging, production

# 数据库配置
database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "ddd_scaffold"
  sslmode: "disable"
  
  # 连接池配置
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: "1h"
  
# Redis 配置
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 100
  
# JWT 配置
jwt:
  secret_key: "your-secret-key-change-in-production"
  expire_in: "24h"
  
# 日志配置
log:
  level: "info"  # debug, info, warn, error
  format: "json" # json, console
  
# LLM 配置（可选）
llm:
  provider: "openai"
  api_key: "sk-xxx"
  base_url: "https://api.openai.com/v1"
  model: "gpt-4"
```

---

## 🌍 环境变量覆盖

生产环境建议使用环境变量：

```bash
# .env 文件
DB_HOST=production-db.example.com
DB_PASSWORD=super-secret-password
JWT_SECRET_KEY=production-secret-key
REDIS_HOST=production-redis.example.com
```

在代码中读取：

```go
// 自动从环境变量读取
dbHost := os.Getenv("DB_HOST")
if dbHost == "" {
    dbHost = config.Database.Host
}
```

---

## 📊 配置项说明

### App 配置

| 字段 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| name | 应用名称 | - | "ddd-scaffold" |
| port | HTTP 端口 | 8080 | 8080 |
| env | 运行环境 | development | development/staging/production |

### Database 配置

| 字段 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| host | 数据库主机 | localhost | localhost |
| port | 数据库端口 | 5432 | 5432 |
| user | 用户名 | postgres | postgres |
| password | 密码 | - | postgres |
| dbname | 数据库名 | ddd_scaffold | ddd_scaffold |
| sslmode | SSL 模式 | disable | disable/require |
| max_idle_conns | 最大空闲连接 | 10 | 10 |
| max_open_conns | 最大打开连接 | 100 | 100 |
| conn_max_lifetime | 连接生命周期 | 1h | 1h |

### Redis 配置

| 字段 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| host | Redis 主机 | localhost | localhost |
| port | Redis 端口 | 6379 | 6379 |
| password | 密码 | - | redis-password |
| db | 数据库编号 | 0 | 0 |
| pool_size | 连接池大小 | 100 | 100 |

### JWT 配置

| 字段 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| secret_key | 密钥 | - | your-secret-key |
| expire_in | 过期时间 | 24h | 24h |

---

## 🔐 安全建议

### 生产环境配置

```yaml
# ✅ 推荐
jwt:
  secret_key: "${JWT_SECRET_KEY}"  # 使用环境变量
  expire_in: "1h"  # 更短的过期时间

database:
  sslmode: "require"  # 启用 SSL
  password: "${DB_PASSWORD}"  # 环境变量
```

### ❌ 避免的做法

```yaml
# ❌ 不要在配置文件中写死敏感信息
jwt:
  secret_key: "my-super-secret-key"  # 禁止！
  
database:
  password: "admin123"  # 禁止！
```

---

## 🧪 多环境配置

### development.yaml

```yaml
app:
  env: "development"
  port: 8080

database:
  host: "localhost"
  sslmode: "disable"

log:
  level: "debug"
```

### production.yaml

```yaml
app:
  env: "production"
  port: 80

database:
  host: "${DB_HOST}"
  sslmode: "require"

log:
  level: "warn"
```

使用不同配置：

```bash
# 开发环境
CONFIG_PATH=config/development.yaml make run

# 生产环境
CONFIG_PATH=config/production.yaml make run
```

---

## 📝 最佳实践

1. **敏感信息使用环境变量**
   ```bash
   export JWT_SECRET_KEY=$(openssl rand -base64 32)
   ```

2. **不同环境使用不同配置文件**
   ```
   config/
   ├── development.yaml
   ├── staging.yaml
   └── production.yaml
   ```

3. **配置验证**
   ```go
   func ValidateConfig(cfg *Config) error {
       if cfg.JWT.SecretKey == "" {
           return errors.New("JWT secret key is required")
       }
       return nil
   }
   ```

4. **配置热重载**（可选）
   ```go
   // 使用 Viper 监听配置变化
   viper.WatchConfig()
   viper.OnConfigChange(func(e fsnotify.Event) {
       log.Println("Config file changed:", e.Name)
   })
   ```

---

**相关文档**:
- [安装指南](installation.md)
- [部署指南](../deployment/docker-deployment.md)
