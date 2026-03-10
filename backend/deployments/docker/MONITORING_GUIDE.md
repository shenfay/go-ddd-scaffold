# Grafana 仪表盘配置

## 应用监控仪表盘 (Dashboard ID: 1)

```json
{
  "dashboard": {
    "title": "DDD Scaffold Application Metrics",
    "panels": [
      {
        "id": 1,
        "title": "HTTP Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{method}} {{endpoint}}"
          }
        ]
      },
      {
        "id": 2,
        "title": "Request Latency (p50, p90, p99)",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p50"
          },
          {
            "expr": "histogram_quantile(0.90, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p90"
          },
          {
            "expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "p99"
          }
        ]
      },
      {
        "id": 3,
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(business_errors_total[5m])",
            "legendFormat": "{{error_type}}"
          }
        ]
      },
      {
        "id": 4,
        "title": "Database Operations",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(db_operations_total[5m])",
            "legendFormat": "{{operation}} - {{status}}"
          }
        ]
      },
      {
        "id": 5,
        "title": "Cache Hit Ratio",
        "type": "gauge",
        "targets": [
          {
            "expr": "avg(cache_hit_ratio)",
            "legendFormat": "Hit Ratio"
          }
        ]
      }
    ]
  }
}
```

## 告警规则配置

### 1. 高错误率告警

```yaml
groups:
  - name: application_alerts
    interval: 30s
   rules:
      # HTTP 5xx 错误率过高
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
         severity: critical
        annotations:
         summary: "高错误率检测"
          description: "接口 {{ $labels.endpoint }} 的 5xx 错误率超过 5%（当前值：{{ $value }}）"
      
      # 业务错误数激增
      - alert: BusinessErrorsSpike
        expr: rate(business_errors_total[5m]) > 10
        for: 3m
        labels:
         severity: warning
        annotations:
         summary: "业务错误激增"
          description: "业务错误数达到 {{ $value }}/秒"
      
      # 响应延迟过高
      - alert: HighLatency
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
         severity: warning
        annotations:
         summary: "响应延迟过高"
          description: "P95 延迟超过 1 秒（当前值：{{ $value }}秒）"
      
      # 数据库操作失败率高
      - alert: DatabaseFailureRate
        expr: rate(db_operations_total{status="failed"}[5m]) / rate(db_operations_total[5m]) > 0.1
        for: 5m
        labels:
         severity: critical
        annotations:
         summary: "数据库失败率过高"
          description: "数据库操作失败率超过 10%（当前值：{{ $value | humanizePercentage }}）"
      
      # 缓存命中率过低
      - alert: LowCacheHitRatio
        expr: avg(cache_hit_ratio) < 0.5
        for: 10m
        labels:
         severity: warning
        annotations:
         summary: "缓存命中率过低"
          description: "缓存命中率低于 50%（当前值：{{ $value | humanizePercentage }}）"
      
      # 服务不可用
      - alert: ServiceDown
        expr: up == 0
        for: 1m
        labels:
         severity: critical
        annotations:
         summary: "服务不可用"
          description: "实例 {{ $labels.instance }} 已下线"
```

### 2. 资源使用告警

```yaml
  - name: resource_alerts
    interval: 30s
   rules:
      # CPU 使用率过高
      - alert: HighCPUUsage
        expr: process_cpu_seconds_total > 0.8
        for: 5m
        labels:
         severity: warning
        annotations:
         summary: "CPU 使用率过高"
          description: "CPU 使用率超过 80%（当前值：{{ $value | humanizePercentage }}）"
      
      # 内存使用过多
      - alert: HighMemoryUsage
        expr: process_resident_memory_bytes > 1073741824  # 1GB
        for: 5m
        labels:
         severity: warning
        annotations:
         summary: "内存使用过多"
          description: "内存使用超过 1GB（当前值：{{ $value | humanizeBytes }}）"
      
      # 磁盘空间不足
      - alert: LowDiskSpace
        expr: node_filesystem_avail_bytes / node_filesystem_size_bytes < 0.1
        for: 10m
        labels:
         severity: critical
        annotations:
         summary: "磁盘空间不足"
          description: "磁盘可用空间低于 10%（当前值：{{ $value | humanizePercentage }}）"
```

## Prometheus 配置文件

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

rule_files:
  - "alerts.yml"

scrape_configs:
  # Prometheus 自监控
  - job_name: 'prometheus'
   static_configs:
      - targets: ['localhost:9090']
  
  # 应用指标采集
  - job_name: 'application'
   static_configs:
      - targets: ['backend:8080']
   metrics_path: '/metrics'
  
  # MySQL 监控
  - job_name: 'mysql'
   static_configs:
      - targets: ['mysqld-exporter:9104']
  
  # Node Exporter 监控
  - job_name: 'node'
   static_configs:
      - targets: ['node-exporter:9100']
```

## AlertManager 通知配置

```yaml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alertmanager@example.com'
  smtp_auth_username: 'alertmanager@example.com'
  smtp_auth_password: 'your-password'

route:
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'default-receiver'
  
  routes:
    - match:
       severity: critical
      receiver: 'critical-alerts'
    - match:
       severity: warning
      receiver: 'warning-alerts'

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'dev-team@example.com'
       send_resolved: true
  
  - name: 'critical-alerts'
    email_configs:
      - to: 'oncall@example.com'
       send_resolved: true
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
        channel: '#alerts-critical'
       send_resolved: true
  
  - name: 'warning-alerts'
    email_configs:
      - to: 'dev-team@example.com'
       send_resolved: true

inhibit_rules:
  - source_match:
     severity: 'critical'
   target_match:
     severity: 'warning'
    equal: ['alertname', 'instance']
```

## Docker Compose 监控栈

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:latest
   container_name: prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./alerts.yml:/etc/prometheus/alerts.yml
      - prometheus_data:/prometheus
   command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/usr/share/prometheus/console_libraries'
      - '--web.console.templates=/usr/share/prometheus/consoles'
    ports:
      - "9090:9090"
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
   container_name: grafana
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=grafana_password
      - GF_USERS_ALLOW_SIGN_UP=false
    ports:
      - "3000:3000"
    networks:
      - monitoring
    depends_on:
      - prometheus

  alertmanager:
    image: prom/alertmanager:latest
   container_name: alertmanager
    volumes:
      - ./alertmanager.yml:/etc/alertmanager/alertmanager.yml
      - alertmanager_data:/alertmanager
   command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
    ports:
      - "9093:9093"
    networks:
      - monitoring

  node-exporter:
    image: prom/node-exporter:latest
   container_name: node-exporter
   command:
      - '--path.rootfs=/host'
    volumes:
      - '/:/host:ro,rslave'
    ports:
      - "9100:9100"
    networks:
      - monitoring

volumes:
  prometheus_data: {}
  grafana_data: {}
  alertmanager_data: {}

networks:
  monitoring:
    driver: bridge
```

## 使用说明

### 1. 启动监控栈

```bash
cd deployments/docker
docker-compose -f docker-compose.monitoring.yml up -d
```

### 2. 访问服务

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/grafana_password)
- **AlertManager**: http://localhost:9093

### 3. 导入仪表盘

在 Grafana 中：
1. 登录 Grafana
2. Dashboard -> Import
3. 上传 dashboard.json 或输入 ID

### 4. 配置告警通知

1. 编辑 `alertmanager.yml`
2. 配置 SMTP 或 Slack Webhook
3. 重启 AlertManager

### 5. 查看指标

应用暴露的指标端点：`http://localhost:8080/metrics`

示例指标：
```
http_requests_total{method="POST",endpoint="/api/auth/login",status="200"} 1234
http_request_duration_seconds{method="POST",endpoint="/api/auth/login"} 0.015
business_errors_total{error_type="client_error",endpoint="/api/users/:id"} 5
cache_hit_ratio{cache_type="user"} 0.85
```

---

## 最佳实践

1. **指标命名规范**: 使用 `_total`, `_seconds`, `_bytes` 等标准后缀
2. **标签基数控制**: 避免高基数标签（如 user_id）
3. **告警分级**: critical/warning/info三级分类
4. **告警抑制**: 防止告警风暴
5. **定期演练**: 验证告警有效性
6. **文档更新**: 保持 Runbook 最新
