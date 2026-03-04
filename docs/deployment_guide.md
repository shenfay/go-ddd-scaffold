# DDD-Scaffold 部署指南

## 目录

- [环境要求](#环境要求)
- [本地开发部署](#本地开发部署)
- [Docker 容器化部署](#docker-容器化部署)
- [Kubernetes 集群部署](#kubernetes-集群部署)
- [生产环境配置](#生产环境配置)
- [监控与日志](#监控与日志)
- [故障排查](#故障排查)

---

## 环境要求

### 基础环境

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| **Go** | >= 1.21 | 推荐使用最新稳定版 |
| **PostgreSQL** | >= 13 | 生产推荐 15+ |
| **Redis** | >= 6.0 | 用于缓存、黑名单、EventBus |
| **Node.js** | >= 18 | 前端构建（可选） |

### 硬件资源（生产环境）

| 资源 | 最小配置 | 推荐配置 |
|------|---------|---------|
| **CPU** | 2 Core | 4-8 Core |
| **内存** | 2 GB | 4-8 GB |
| **磁盘** | 20 GB | 100 GB SSD |
| **带宽** | 10 Mbps | 100 Mbps |

---

## 本地开发部署

### 1. 克隆项目

```bash
git clone https://github.com/your-org/ddd-scaffold.git
cd ddd-scaffold
```

### 2. 安装依赖

#### 后端依赖

```bash
cd backend
go mod download
```

#### 前端依赖（可选）

```bash
cd frontend
pnpm install  # 或 npm install
```

### 3. 启动基础设施服务

#### 使用 Docker Compose（推荐）

```bash
cd deployments/docker
docker-compose up -d postgres redis
```

#### 手动安装

**macOS (Homebrew):**
```bash
brew install postgresql@15 redis
brew services start postgresql@15
brew services start redis
```

**Linux (Ubuntu):**
```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib redis-server
sudo systemctl start postgresql
sudo systemctl start redis-server
```

### 4. 初始化数据库

```bash
# 创建数据库
createdb go_ddd_scaffold
createdb go_ddd_scaffold_test

# 运行迁移
cd backend
go run cmd/migrate/main.go up
```

### 5. 配置文件

编辑 `backend/config/config.yaml`：

```yaml
app:
  name: "Go DDD Scaffold"
  port: 8080
  env: "development"

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  dbname: "go_ddd_scaffold"
  sslmode: "disable"
  
  # 连接池配置
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 1h
  conn_max_idle_time: 5m

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 100
  min_idle_conns: 10
  
  # EventBus 配置
  event_bus:
    stream_key: "domain_events"
    max_retries: 3
    retry_base_delay: 1s
    poll_interval: 5s
    batch_size: 100

jwt:
  secret_key: "mathfun-jwt-secret-key-2026-change-in-production"
  expire_in: 24h

log:
  level: "debug"  # development: debug, production: info
```

### 6. 启动服务

#### 后端服务

```bash
cd backend
go run cmd/server/main.go
```

访问：http://localhost:8080

#### Swagger API 文档

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init -g cmd/server/main.go -o ./docs

# 启动服务后访问
http://localhost:8080/swagger/index.html
```

#### 前端服务（可选）

```bash
cd frontend
pnpm dev
```

访问：http://localhost:3000

---

## Docker 容器化部署

### 1. 构建镜像

#### 后端镜像

```bash
cd backend
docker build -t ddd-scaffold-backend:latest .
```

#### 前端镜像（可选）

```bash
cd frontend
docker build -t ddd-scaffold-frontend:latest .
```

### 2. 使用 Docker Compose 部署

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  # PostgreSQL 数据库
  postgres:
    image: postgres:15-alpine
    container_name: ddd-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_ddd_scaffold
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ../migrations/sql:/docker-entrypoint-initdb.d
    networks:
      - ddd-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: ddd-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - ddd-network
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # 后端服务
  backend:
    image: ddd-scaffold-backend:latest
    container_name: ddd-backend
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    volumes:
      - ./config:/app/config
    networks:
      - ddd-network
    restart: unless-stopped

  # 前端服务（可选）
  frontend:
    image: ddd-scaffold-frontend:latest
    container_name: ddd-frontend
    depends_on:
      - backend
    ports:
      - "80:80"
    networks:
      - ddd-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  ddd-network:
    driver: bridge
```

### 3. 启动服务

```bash
docker-compose up -d
```

### 4. 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f backend
docker-compose logs -f postgres
```

### 5. 停止服务

```bash
docker-compose down

# 删除数据卷（谨慎使用）
docker-compose down -v
```

---

## Kubernetes 集群部署

### 1. 准备 K8s 资源配置

详见 `deployments/k8s/` 目录：

```yaml
# deployments/k8s/backend-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: ddd-scaffold-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: app-config
              key: db.host
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### 2. 部署到 K8s

```bash
cd deployments/k8s

# 创建命名空间
kubectl create namespace ddd-scaffold

# 应用配置
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f postgres-statefulset.yaml
kubectl apply -f redis-statefulset.yaml
kubectl apply -f backend-deployment.yaml
kubectl apply -f backend-service.yaml
kubectl apply -f ingress.yaml
```

### 3. 自动扩缩容

```bash
# 创建 HPA
kubectl autoscale deployment backend --cpu-percent=80 --min=3 --max=10

# 查看 HPA 状态
kubectl get hpa
```

---

## 生产环境配置

### 1. 环境变量配置

创建 `.env.production`：

```bash
# 应用配置
APP_ENV=production
APP_PORT=8080

# 数据库配置（使用环境变量覆盖）
DB_HOST=prod-db.example.com
DB_PORT=5432
DB_USER=prod_user
DB_PASSWORD=${DB_PASSWORD_SECRET}
DB_NAME=go_ddd_scaffold_prod
DB_SSLMODE=require

# Redis 配置
REDIS_HOST=prod-redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=${REDIS_PASSWORD_SECRET}
REDIS_DB=0

# JWT 配置（必须修改为强密钥）
JWT_SECRET_KEY=${JWT_SECRET_KEY_GENERATED}
JWT_EXPIRE_IN=24h

# 日志配置
LOG_LEVEL=info
LOG_FORMAT=json
```

### 2. 安全加固

#### 生成 JWT 密钥

```bash
# 生成随机密钥
openssl rand -base64 64
```

#### 配置 SSL/TLS

```yaml
# config.yaml
database:
  sslmode: "require"
  
# 启用 HTTPS
server:
  https:
    enabled: true
    cert_file: /etc/ssl/certs/server.crt
    key_file: /etc/ssl/private/server.key
```

### 3. 性能优化

#### 数据库连接池调优

```yaml
database:
  max_idle_conns: 20          # 增加空闲连接
  max_open_conns: 200         # 根据并发调整
  conn_max_lifetime: 30m      # 延长连接生命周期
  conn_max_idle_time: 10m     # 减少频繁重建
```

#### Redis 连接池调优

```yaml
redis:
  pool_size: 200              # 增加连接池大小
  min_idle_conns: 20          # 保持最小空闲连接
```

### 4. 备份策略

#### 数据库备份

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backups/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/go_ddd_scaffold_$DATE.sql.gz"

# 全量备份
pg_dump -h $DB_HOST -U postgres go_ddd_scaffold | gzip > $BACKUP_FILE

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE"
```

#### Redis 持久化

```yaml
# Redis 配置
appendonly yes                    # 启用 AOF
appendfsync everysec              # 每秒同步
save 900 1                        # RDB 快照
save 300 10
save 60 10000
```

---

## 监控与日志

### 1. Prometheus 监控

#### 部署 Prometheus

```yaml
# prometheus-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: prometheus
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: prometheus
        image: prom/prometheus:v2.45.0
        ports:
        - containerPort: 9090
        volumeMounts:
        - name: prometheus-config
          mountPath: /etc/prometheus
      volumes:
      - name: prometheus-config
        configMap:
          data:
            prometheus.yml: |
              global:
                scrape_interval: 15s
              scrape_configs:
              - job_name: 'backend'
                static_configs:
                - targets: ['backend:8080']
                metrics_path: '/metrics'
```

#### 关键监控指标

```prometheus
# QPS
rate(http_requests_total[1m])

# 延迟分布
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# 错误率
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m])

# Redis 连接数
redis_connected_clients

# 熔断器状态
circuit_breaker_state{resource="redis"}
```

### 2. Grafana 可视化

导入 Dashboard ID:
- **Backend Performance**: 自定义
- **Redis Overview**: 763
- **PostgreSQL Overview**: 9628

### 3. 日志聚合

#### ELK Stack 配置

```yaml
# Filebeat 配置
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/app/*.log
  json.keys_under_root: true
  
output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "ddd-scaffold-%{+yyyy.MM.dd}"
```

---

## 故障排查

### 常见问题

#### 1. 数据库连接失败

**现象：**
```
failed to connect to database: dial tcp: lookup postgres: no such host
```

**解决方案：**
```bash
# 检查 DNS 解析
kubectl get svc postgres

# 检查网络策略
kubectl describe networkpolicy

# 验证连接
psql -h postgres -U postgres -c "SELECT 1"
```

#### 2. Redis 连接超时

**现象：**
```
redis: connection pool: failed to dial after 5 attempts
```

**解决方案：**
```bash
# 检查 Redis 健康状态
redis-cli ping

# 查看 Redis 日志
kubectl logs redis-0

# 增加连接超时
redis.options.read_timeout = 5s
```

#### 3. JWT Token 验证失败

**现象：**
```
token signature is invalid
```

**解决方案：**
```bash
# 检查 JWT 密钥一致性
kubectl get secret app-secret -o jsonpath='{.data.JWT_SECRET}' | base64 -d

# 确保所有实例使用相同密钥
kubectl rollout restart deployment/backend
```

#### 4. 熔断器频繁跳闸

**现象：**
```
circuit breaker state: OPEN
```

**解决方案：**
```bash
# 检查 Redis 性能
redis-cli slowlog get 10

# 调整熔断器配置
config.MaxFailures = 10
config.ResetTimeout = 60 * time.Second

# 增加 Redis 资源
kubectl set resources redis --limits=memory=1Gi,cpu=500m
```

### 性能问题诊断

#### 慢查询分析

```sql
-- PostgreSQL 慢查询
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

#### Redis 慢日志

```bash
redis-cli slowlog get 10
redis-cli config set slowlog-log-slower-than 10000  # 10ms
```

#### 应用性能分析

```bash
# 启用 pprof
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# 查看 goroutine
go tool pprof http://localhost:8080/debug/pprof/goroutine
```

---

## 总结

### 部署检查清单

- [ ] 数据库迁移完成
- [ ] Redis 连接正常
- [ ] JWT 密钥已更新
- [ ] 环境变量配置正确
- [ ] 健康检查通过
- [ ] 监控指标正常
- [ ] 日志输出正常
- [ ] 备份策略生效

### 性能基准

| 指标 | 目标值 | 实测值 |
|------|--------|--------|
| **P99 延迟** | < 100ms | ~50ms |
| **QPS** | > 1000 | ~2000 |
| **错误率** | < 0.1% | ~0.01% |
| **可用性** | > 99.9% | 99.95% |

### 下一步优化

1. **水平扩展** - 增加后端实例数
2. **读写分离** - PostgreSQL 主从复制
3. **Redis Cluster** - 分布式缓存
4. **CDN 加速** - 静态资源分发
5. **多区域部署** - 提高全球访问速度
