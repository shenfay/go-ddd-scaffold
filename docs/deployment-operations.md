# Go DDD Scaffold 部署运维文档

## 文档概述

本文档详细描述了 go-ddd-scaffold 项目的部署方案、运维配置、监控告警以及故障排查指南。

## 部署架构设计

### 单体应用部署架构
```
┌─────────────────────────────────────────────────────┐
│                    Load Balancer                     │
│                   (Nginx/Haproxy)                    │
├─────────────────────────────────────────────────────┤
│                 Application Servers                  │
│              [Go App Instance 1..N]                  │
├─────────────┬───────────────────┬───────────────────┤
│   Cache     │    Database       │   Monitoring      │
│  (Redis)    │  (PostgreSQL)     │   (Prometheus)    │
└─────────────┴───────────────────┴───────────────────┘
```

### 环境配置管理

#### 多环境配置文件
```yaml
# configs/development.yaml
server:
  port: 8080
  mode: debug
  read_timeout: 30s
  write_timeout: 30s

database:
  host: localhost
  port: 5432
  name: scaffold_dev
  user: postgres
  password: dev_password
  max_idle_conns: 5
  max_open_conns: 20

redis:
  addr: localhost:6379
  password: ""
  db: 0

logging:
  level: debug
  format: console
  file: "./logs/app.log"
```

```yaml
# configs/production.yaml
server:
  port: 80
  mode: release
  read_timeout: 60s
  write_timeout: 60s

database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  name: ${DB_NAME}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  max_idle_conns: 20
  max_open_conns: 100
  conn_max_lifetime: 1h

redis:
  addr: ${REDIS_ADDR}
  password: ${REDIS_PASSWORD}
  db: 0

logging:
  level: info
  format: json
  file: "/var/log/scaffold/app.log"
```

## 部署方案

### 1. 本地开发环境部署

#### Docker Compose部署（推荐）
```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENV_MODE=development
      - DB_HOST=postgres
      - REDIS_ADDR=redis:6379
    depends_on:
      - postgres
      - redis
    volumes:
      - ./logs:/app/logs

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=scaffold_dev
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
    depends_on:
      - app

volumes:
  postgres_data:
```

#### 启动命令
```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看应用日志
docker-compose logs -f app

# 运行数据库迁移
docker-compose exec app ./migrate.sh up
```

### 2. 生产环境部署

#### Kubernetes部署方案
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: scaffold-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: scaffold
  template:
    metadata:
      labels:
        app: scaffold
    spec:
      containers:
      - name: app
        image: your-registry/scaffold-app:v1.0.0
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: scaffold-config
        - secretRef:
            name: scaffold-secrets
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 30
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: scaffold-service
spec:
  selector:
    app: scaffold
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

#### Helm Chart配置
```yaml
# charts/scaffold/values.yaml
replicaCount: 3

image:
  repository: your-registry/scaffold-app
  tag: v1.0.0
  pullPolicy: IfNotPresent

env:
  ENV_MODE: production
  LOG_LEVEL: info

resources:
  limits:
    cpu: 500m
    memory: 512Mi
  requests:
    cpu: 250m
    memory: 256Mi

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  hosts:
    - host: api.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
```

### 3. 传统服务器部署

#### systemd服务配置
```ini
# /etc/systemd/system/scaffold.service
[Unit]
Description=Go DDD Scaffold Application
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=scaffold
Group=scaffold
WorkingDirectory=/opt/scaffold
ExecStart=/opt/scaffold/bin/server
Restart=always
RestartSec=10
Environment=ENV_MODE=production
EnvironmentFile=/etc/scaffold/environment

[Install]
WantedBy=multi-user.target
```

#### 部署脚本
```bash
#!/bin/bash
# deploy.sh

set -e

APP_NAME="scaffold"
VERSION="v1.0.0"
DEPLOY_DIR="/opt/${APP_NAME}"
BACKUP_DIR="/opt/${APP_NAME}_backup_$(date +%Y%m%d_%H%M%S)"

echo "Starting deployment of ${APP_NAME} ${VERSION}"

# 备份当前版本
if [ -d "${DEPLOY_DIR}" ]; then
    echo "Backing up current version..."
    cp -r "${DEPLOY_DIR}" "${BACKUP_DIR}"
fi

# 创建部署目录
mkdir -p "${DEPLOY_DIR}"/{bin,config,logs}

# 下载新版本
echo "Downloading new version..."
curl -L "https://github.com/your-org/${APP_NAME}/releases/download/${VERSION}/${APP_NAME}.tar.gz" \
    | tar -xz -C "${DEPLOY_DIR}"

# 复制配置文件
cp -n ./config/* "${DEPLOY_DIR}/config/"

# 设置权限
chown -R scaffold:scaffold "${DEPLOY_DIR}"
chmod +x "${DEPLOY_DIR}/bin/server"

# 重启服务
echo "Restarting service..."
systemctl daemon-reload
systemctl restart "${APP_NAME}.service"

# 等待服务启动
sleep 10

# 检查服务状态
if systemctl is-active --quiet "${APP_NAME}.service"; then
    echo "Deployment successful!"
else
    echo "Deployment failed, rolling back..."
    rm -rf "${DEPLOY_DIR}"
    mv "${BACKUP_DIR}" "${DEPLOY_DIR}"
    systemctl restart "${APP_NAME}.service"
    exit 1
fi
```

## 监控告警配置

### Prometheus监控指标
```go
// metrics/metrics.go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP请求指标
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    // 业务指标
    ActiveUsers = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_users",
            Help: "Number of active users",
        },
    )

    DatabaseConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "database_connections",
            Help: "Current database connections",
        },
    )

    // 系统指标
    Goroutines = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "goroutines",
            Help: "Number of goroutines",
        },
    )
)
```

### Grafana仪表板配置
```json
{
  "dashboard": {
    "title": "Scaffold Application Metrics",
    "panels": [
      {
        "title": "HTTP Requests Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}} {{status}}"
          }
        ]
      },
      {
        "title": "Request Latency",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      },
      {
        "title": "Database Connections",
        "type": "gauge",
        "targets": [
          {
            "expr": "database_connections",
            "legendFormat": "Current Connections"
          }
        ]
      }
    ]
  }
}
```

### 告警规则配置
```yaml
# prometheus/rules/alerts.yml
groups:
- name: scaffold-alerts
  rules:
  - alert: HighErrorRate
    expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is above 5% for the last 5 minutes"

  - alert: HighLatency
    expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High request latency"
      description: "95th percentile latency is above 2 seconds"

  - alert: DatabaseDown
    expr: database_connections == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Database connection lost"
      description: "No database connections available"
```

## 日志管理

### 结构化日志配置
```go
// logging/logger.go
package logging

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger(config LoggerConfig) (*zap.Logger, error) {
    var encoderCfg zapcore.EncoderConfig
    if config.Format == "json" {
        encoderCfg = zap.NewProductionEncoderConfig()
    } else {
        encoderCfg = zap.NewDevelopmentEncoderConfig()
    }

    encoderCfg.TimeKey = "timestamp"
    encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

   	var level zapcore.Level
    if err := level.UnmarshalText([]byte(config.Level)); err != nil {
        return nil, err
    }

    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderCfg),
        zapcore.AddSync(&lumberjack.Logger{
            Filename:   config.File,
            MaxSize:    100, // MB
            MaxBackups: 3,
            MaxAge:     28, // days
            Compress:   true,
        }),
        level,
    )

    return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)), nil
}
```

### 日志轮转配置
```yaml
logging:
  level: info
  format: json
  file: "/var/log/scaffold/app.log"
  maxSize: 100
  maxBackups: 5
  maxAge: 30
  compress: true
```

### ELK日志收集
```yaml
# filebeat/filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/scaffold/*.log
  json.keys_under_root: true
  json.add_error_key: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "scaffold-%{+yyyy.MM.dd}"

setup.kibana:
  host: "kibana:5601"

setup.template.name: "scaffold"
setup.template.pattern: "scaffold-*"
```

## 健康检查与探针

### 应用健康检查接口
```go
// health/health.go
package health

import (
    "context"
    "net/http"
    "time"
)

type Checker interface {
    Check(ctx context.Context) error
}

type HealthHandler struct {
    checkers map[string]Checker
    timeout  time.Duration
}

func NewHealthHandler(timeout time.Duration) *HealthHandler {
    return &HealthHandler{
        checkers: make(map[string]Checker),
        timeout:  timeout,
    }
}

func (h *HealthHandler) AddChecker(name string, checker Checker) {
    h.checkers[name] = checker
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
    defer cancel()

    results := make(map[string]string)
    healthy := true

    for name, checker := range h.checkers {
        if err := checker.Check(ctx); err != nil {
            results[name] = err.Error()
            healthy = false
        } else {
            results[name] = "OK"
        }
    }

    status := http.StatusOK
    if !healthy {
        status = http.StatusServiceUnavailable
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  getStatus(healthy),
        "checks":  results,
        "uptime":  time.Since(startTime).String(),
        "version": version,
    })
}

func getStatus(healthy bool) string {
    if healthy {
        return "healthy"
    }
    return "unhealthy"
}
```

### 健康检查配置
```go
// main.go
func setupHealthChecks() *health.HealthHandler {
    handler := health.NewHealthHandler(5 * time.Second)
    
    // 数据库健康检查
    handler.AddChecker("database", &health.DatabaseChecker{
        DB: db,
    })
    
    // Redis健康检查
    handler.AddChecker("redis", &health.RedisChecker{
        Client: redisClient,
    })
    
    // 业务健康检查
    handler.AddChecker("business", &health.BusinessChecker{})
    
    return handler
}
```

## 故障排查指南

### 常见问题诊断

#### 1. 启动失败
```bash
# 检查服务状态
systemctl status scaffold.service

# 查看详细日志
journalctl -u scaffold.service -f

# 检查端口占用
netstat -tlnp | grep :8080

# 验证配置文件
./server --config=config.yaml --validate
```

#### 2. 数据库连接问题
```bash
# 测试数据库连接
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME

# 检查连接池状态
curl http://localhost:8080/debug/vars | jq '.database'

# 查看慢查询日志
tail -f /var/log/postgresql/postgresql-15-main.log
```

#### 3. 性能问题排查
```bash
# CPU和内存使用情况
top -p $(pgrep server)

# Go程序性能分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存分配分析
go tool pprof http://localhost:8080/debug/pprof/heap

# 协程分析
curl http://localhost:8080/debug/pprof/goroutine?debug=2
```

#### 4. 网络问题
```bash
# 检查网络连通性
ping $DB_HOST
telnet $DB_HOST $DB_PORT

# DNS解析检查
nslookup $DB_HOST

# 防火墙规则检查
iptables -L -n | grep 8080
```

### 应急处理流程

#### 服务恢复步骤
```bash
#!/bin/bash
# emergency-recovery.sh

SERVICE_NAME="scaffold"
BACKUP_DIR="/opt/${SERVICE_NAME}_backup"

echo "Starting emergency recovery procedure..."

# 1. 停止当前服务
systemctl stop ${SERVICE_NAME}.service

# 2. 恢复备份
if [ -d "${BACKUP_DIR}" ]; then
    echo "Restoring from backup..."
    rm -rf /opt/${SERVICE_NAME}
    cp -r ${BACKUP_DIR} /opt/${SERVICE_NAME}
fi

# 3. 重启服务
systemctl start ${SERVICE_NAME}.service

# 4. 验证服务状态
sleep 10
if systemctl is-active --quiet ${SERVICE_NAME}.service; then
    echo "Service recovered successfully"
else
    echo "Recovery failed, manual intervention required"
    exit 1
fi

# 5. 通知相关人员
send_alert "Service ${SERVICE_NAME} has been recovered"
```

#### 数据库恢复流程
```bash
#!/bin/bash
# database-recovery.sh

DB_NAME="scaffold_db"
BACKUP_FILE="/backup/$(date +%Y%m%d)_${DB_NAME}.sql"

echo "Starting database recovery..."

# 1. 停止应用服务
systemctl stop scaffold.service

# 2. 删除当前数据库
psql -c "DROP DATABASE IF EXISTS ${DB_NAME};"

# 3. 创建新数据库
psql -c "CREATE DATABASE ${DB_NAME};"

# 4. 恢复数据
psql ${DB_NAME} < ${BACKUP_FILE}

# 5. 运行迁移
./migrate.sh up

# 6. 重启服务
systemctl start scaffold.service

echo "Database recovery completed"
```

## 备份策略

### 自动备份脚本
```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup"
DATE=$(date +%Y%m%d_%H%M%S)
DB_NAME="scaffold_db"

# 创建备份目录
mkdir -p ${BACKUP_DIR}/${DATE}

# 数据库备份
pg_dump ${DB_NAME} > ${BACKUP_DIR}/${DATE}/database.sql

# 应用配置备份
tar -czf ${BACKUP_DIR}/${DATE}/config.tar.gz /opt/scaffold/config/

# 日志备份
tar -czf ${BACKUP_DIR}/${DATE}/logs.tar.gz /opt/scaffold/logs/

# 清理旧备份（保留7天）
find ${BACKUP_DIR} -type d -mtime +7 -exec rm -rf {} \;

echo "Backup completed: ${BACKUP_DIR}/${DATE}"
```

### 备份验证
```bash
#!/bin/bash
# verify-backup.sh

BACKUP_DIR="/backup"
LATEST_BACKUP=$(ls -t ${BACKUP_DIR} | head -1)

echo "Verifying backup: ${LATEST_BACKUP}"

# 验证数据库备份
if [ -f "${BACKUP_DIR}/${LATEST_BACKUP}/database.sql" ]; then
    echo "✓ Database backup exists"
else
    echo "✗ Database backup missing"
    exit 1
fi

# 验证配置备份
if [ -f "${BACKUP_DIR}/${LATEST_BACKUP}/config.tar.gz" ]; then
    echo "✓ Config backup exists"
else
    echo "✗ Config backup missing"
    exit 1
fi

echo "Backup verification completed successfully"
```

这个部署运维文档为项目提供了完整的生产环境部署和运维管理指南。