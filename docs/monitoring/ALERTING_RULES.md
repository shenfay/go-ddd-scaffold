# 告警规则配置

本文档定义了 Prometheus 告警规则,用于监控系统健康和业务指标。

## 告警规则文件

告警规则配置文件位于: `backend/configs/alerting_rules.yml`

## 告警分类

### 1. 服务健康告警

#### ServiceDown
- **条件**: 服务实例 down
- **持续时间**: 1 分钟
- **严重级别**: 🔴 Critical
- **通知**: PagerDuty, Slack
- **处理**: 自动重启,检查日志

```yaml
alert: ServiceDown
expr: up == 0
for: 1m
labels:
  severity: critical
annotations:
  summary: "Service {{ $labels.instance }} is down"
  description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 1 minute."
```

#### HighErrorRate
- **条件**: HTTP 5xx 错误率 > 5%
- **持续时间**: 5 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查应用日志,查看依赖服务状态

```yaml
alert: HighErrorRate
expr: |
  sum(rate(http_requests_total{status=~"5.."}[5m])) 
  / 
  sum(rate(http_requests_total[5m])) 
  > 0.05
for: 5m
labels:
  severity: warning
annotations:
  summary: "High HTTP error rate"
  description: "HTTP 5xx error rate is {{ $value | humanizePercentage }} (> 5%)"
```

#### HighLatency
- **条件**: P95 延迟 > 1s
- **持续时间**: 10 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查数据库慢查询,缓存命中率

```yaml
alert: HighLatency
expr: |
  histogram_quantile(0.95, 
    sum(rate(http_request_duration_seconds_bucket[10m])) 
    by (le, path)
  ) > 1
for: 10m
labels:
  severity: warning
annotations:
  summary: "High latency on {{ $labels.path }}"
  description: "P95 latency is {{ $value }}s (> 1s)"
```

---

### 2. 数据库告警

#### DatabaseConnectionPoolExhausted
- **条件**: 数据库连接使用率 > 90%
- **持续时间**: 5 分钟
- **严重级别**: 🔴 Critical
- **通知**: PagerDuty, Slack
- **处理**: 增加连接池大小,检查连接泄漏

```yaml
alert: DatabaseConnectionPoolExhausted
expr: |
  db_connections_open / db_connections_max 
  > 0.9
for: 5m
labels:
  severity: critical
annotations:
  summary: "Database connection pool nearly exhausted"
  description: "Connection pool usage is {{ $value | humanizePercentage }} (> 90%)"
```

#### SlowQueriesIncreasing
- **条件**: 慢查询速率 > 10/min
- **持续时间**: 10 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 优化查询,添加索引

```yaml
alert: SlowQueriesIncreasing
expr: |
  rate(db_slow_queries_total[10m]) 
  > 10/60
for: 10m
labels:
  severity: warning
annotations:
  summary: "High rate of slow queries"
  description: "Slow query rate is {{ $value | humanize }} per second"
```

---

### 3. Redis 告警

#### RedisMemoryHigh
- **条件**: Redis 内存使用 > 80%
- **持续时间**: 5 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 清理过期 key,增加内存

```yaml
alert: RedisMemoryHigh
expr: |
  redis_memory_used_bytes 
  > 800 * 1024 * 1024  # 800MB threshold
for: 5m
labels:
  severity: warning
annotations:
  summary: "Redis memory usage high"
  description: "Redis memory usage is {{ $value | humanize1024 }}B"
```

#### RedisHitRateLow
- **条件**: Redis 命中率 < 50%
- **持续时间**: 15 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查缓存策略,key 过期时间

```yaml
alert: RedisHitRateLow
expr: |
  redis_hit_rate 
  < 0.5
for: 15m
labels:
  severity: warning
annotations:
  summary: "Redis cache hit rate is low"
  description: "Cache hit rate is {{ $value | humanizePercentage }} (< 50%)"
```

---

### 4. 认证业务告警

#### LoginFailureRateHigh
- **条件**: 登录失败率 > 20%
- **持续时间**: 5 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查是否有暴力破解,账户锁定策略

```yaml
alert: LoginFailureRateHigh
expr: |
  rate(auth_failure_total{type="login"}[5m]) 
  / 
  (rate(auth_success_total{type="login"}[5m]) + rate(auth_failure_total{type="login"}[5m])) 
  > 0.2
for: 5m
labels:
  severity: warning
annotations:
  summary: "High login failure rate"
  description: "Login failure rate is {{ $value | humanizePercentage }} (> 20%)"
```

#### RegistrationDrop
- **条件**: 注册率下降 > 50% (对比 1 小时前)
- **持续时间**: 30 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查注册流程,邮箱服务状态

```yaml
alert: RegistrationDrop
expr: |
  rate(user_registrations_total[30m]) 
  < 
  rate(user_registrations_total[30m] offset 1h) * 0.5
for: 30m
labels:
  severity: warning
annotations:
  summary: "User registration rate dropped significantly"
  description: "Registration rate dropped by more than 50% compared to 1 hour ago"
```

#### PasswordResetFailureRateHigh
- **条件**: 密码重置失败率 > 30%
- **持续时间**: 10 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查 Token 生成逻辑,邮件服务

```yaml
alert: PasswordResetFailureRateHigh
expr: |
  rate(password_reset_failed_total[10m]) 
  / 
  (rate(password_reset_completed_total[10m]) + rate(password_reset_failed_total[10m])) 
  > 0.3
for: 10m
labels:
  severity: warning
annotations:
  summary: "High password reset failure rate"
  description: "Password reset failure rate is {{ $value | humanizePercentage }} (> 30%)"
```

---

### 5. 限流告警

#### RateLimitTriggered
- **条件**: 限流拒绝率 > 10/min
- **持续时间**: 5 分钟
- **严重级别**: 🟡 Warning
- **通知**: Slack
- **处理**: 检查是否有异常流量,调整限流阈值

```yaml
alert: RateLimitTriggered
expr: |
  sum(rate(ratelimit_rejected_total[5m])) 
  > 10/60
for: 5m
labels:
  severity: warning
annotations:
  summary: "High rate of rate-limited requests"
  description: "Rate limit rejection rate is {{ $value | humanize }} per second"
```

---

## 告警通知配置

### Prometheus Alertmanager 配置

```yaml
global:
  resolve_timeout: 5m
  slack_api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
  pagerduty_url: 'https://events.pagerduty.com/v2/enqueue'

route:
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'slack-notifications'
  routes:
    - match:
        severity: critical
      receiver: 'pagerduty-critical'
      repeat_interval: 1h
    - match:
        severity: warning
      receiver: 'slack-notifications'
      repeat_interval: 4h

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#alerts'
        send_resolved: true
        title: '{{ .GroupLabels.alertname }}'
        text: |
          *Alert:* {{ .CommonAnnotations.summary }}
          *Description:* {{ .CommonAnnotations.description }}
          *Severity:* {{ .CommonLabels.severity }}
          *Instances:* {{ range .Alerts }}{{ .Labels.instance }} {{ end }}

  - name: 'pagerduty-critical'
    pagerduty_configs:
      - service_key: 'YOUR_PAGERDUTY_SERVICE_KEY'
        severity: 'critical'
        description: '{{ .CommonAnnotations.summary }}'
```

---

## 告警响应流程

### Critical 告警

1. **立即响应** (5 分钟内)
   - 确认告警真实性
   - 检查服务状态
   - 查看错误日志

2. **诊断问题** (15 分钟内)
   - 定位根因
   - 评估影响范围
   - 制定修复方案

3. **解决问题** (30 分钟内)
   - 执行修复
   - 验证恢复
   - 记录事件

### Warning 告警

1. **及时响应** (30 分钟内)
   - 确认告警
   - 初步分析

2. **计划修复** (2 小时内)
   - 制定修复计划
   - 安排维护窗口

3. **执行修复** (24 小时内)
   - 实施修复
   - 监控效果

---

## 告警优化建议

### 减少误报

1. **调整阈值**: 根据历史数据调整
2. **增加持续时间**: 避免短暂波动触发
3. **使用率而非绝对值**: 更准确反映问题

### 减少告警疲劳

1. **告警分组**: 相关告警合并
2. **静默规则**: 维护期间静默
3. **升级策略**: 未处理告警升级

### 持续改进

1. **定期审查**: 每周审查告警规则
2. **效果评估**: 统计告警准确率
3. **规则优化**: 删除无效告警

---

## 监控仪表盘

- **业务指标**: [Grafana - Business Metrics](http://localhost:3000/d/ddd-scaffold-monitoring-business)
- **认证监控**: [Grafana - Authentication](http://localhost:3000/d/ddd-scaffold-monitoring-auth)
- **API 性能**: [Grafana - API Performance](http://localhost:3000/d/ddd-scaffold-monitoring-api)
- **数据库**: [Grafana - Database](http://localhost:3000/d/ddd-scaffold-monitoring-db)
- **Redis**: [Grafana - Redis](http://localhost:3000/d/ddd-scaffold-monitoring-redis)
- **限流**: [Grafana - Rate Limiting](http://localhost:3000/d/ddd-scaffold-monitoring-ratelimit)

---

## 相关资源

- [Prometheus 告警文档](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)
- [Alertmanager 配置](https://prometheus.io/docs/alerting/latest/configuration/)
- [告警最佳实践](https://prometheus.io/docs/practices/alerting/)
