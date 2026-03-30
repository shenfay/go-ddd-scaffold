# Asynq 监控与告警部署指南

## 📋 目录

1. [Prometheus 配置](#prometheus-配置)
2. [Grafana 仪表盘导入](#grafana-仪表盘导入)
3. [AlertManager 告警配置](#alertmanager-告警配置)
4. [验证与测试](#验证与测试)

---

## 🔧 Prometheus 配置

### 1. 添加 scrape 配置

在 `prometheus.yml` 中添加 Worker 的监控端点：

```yaml
scrape_configs:
  - job_name: 'worker'
    static_configs:
      - targets: ['localhost:9091']  # Worker 暴露 metrics 的端口
    scrape_interval: 10s
    scrape_timeout: 5s
    metrics_path: /metrics
```

### 2. 添加告警规则

将告警规则文件添加到 Prometheus 配置：

```yaml
rule_files:
  - "asynq-alerts.yml"
```

或者直接复制配置文件：

```bash
cp deployments/prometheus/asynq-alerts.yml /etc/prometheus/rules/asynq-alerts.yml
systemctl reload prometheus
```

### 3. 重启 Prometheus

```bash
systemctl daemon-reload
systemctl restart prometheus
```

---

## 📊 Grafana 仪表盘导入

### 方式一：导入 JSON 文件（推荐）

1. **打开 Grafana**
   - 访问 http://localhost:3000
   - 登录（默认 admin/admin）

2. **导入仪表盘**
   - 点击左侧菜单 "+" → "Import"
   - 点击 "Upload dashboard JSON file"
   - 选择 `deployments/grafana/asynq-dashboard.json`
   - 点击 "Import"

3. **配置数据源**
   - 选择 Prometheus 数据源
   - 点击 "Import" 完成

### 方式二：使用 Dashboard ID

如果 Grafana 支持在线导入：

1. 访问 https://grafana.com/grafana/dashboards/
2. 搜索 "Asynq Task Queue Monitoring"
3. 输入 Dashboard ID
4. 点击 "Load"

### 仪表盘说明

| 面板 | 说明 | 告警阈值 |
|------|------|---------|
| 任务处理速率 | 实时显示各类型任务的处理速度 | - |
| 任务处理延迟 | P50/P95/P99 延迟监控 | P99 > 60s |
| 队列积压情况 | 各队列待处理任务数 | > 1000 |
| 任务失败率 | 任务失败占比 | > 5% |
| 总队列积压 | 所有队列积压总数 | > 100 |
| 总体失败率 | 所有任务的平均失败率 | > 10% |
| Worker 并发度 | 当前 Worker 并发数 | < 5 |
| 积压队列数 | 积压超过阈值的队列数 | > 10 |

---

## 🚨 AlertManager 告警配置

### 1. 配置通知渠道

编辑 `alertmanager.yml`：

```yaml
global:
  smtp_smarthost: 'smtp.example.com:587'
  smtp_from: 'alertmanager@example.com'
  smtp_auth_username: 'alertmanager@example.com'
  smtp_auth_password: 'password'

route:
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  receiver: 'default-receiver'
  
  routes:
    - match:
        severity: critical
      receiver: 'critical-receiver'
    - match:
        severity: warning
      receiver: 'warning-receiver'

receivers:
  - name: 'default-receiver'
    email_configs:
      - to: 'team@example.com'
        send_resolved: true
  
  - name: 'critical-receiver'
    email_configs:
      - to: 'oncall@example.com'
        send_resolved: true
    webhook_configs:
      - url: 'http://localhost:5001/webhook'  # PagerDuty/OpsGenie
  
  - name: 'warning-receiver'
    email_configs:
      - to: 'dev-team@example.com'
        send_resolved: true

inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'queue']
```

### 2. 告警级别说明

| 级别 | 触发条件 | 响应时间 | 通知渠道 |
|------|---------|---------|---------|
| **Critical** | 失败率>5%、关键队列积压>100、Worker 宕机 | 立即响应 | 邮件 + 电话 |
| **Warning** | 队列积压>1000、延迟>60s、并发度<5 | 30 分钟内 | 邮件 |

### 3. 告警处理流程

```
告警触发 → AlertManager 路由 → 通知接收者 → 确认告警 → 问题修复 → 告警恢复
```

---

## ✅ 验证与测试

### 1. 验证 Prometheus 配置

```bash
# 检查配置文件语法
promtool check config /etc/prometheus/prometheus.yml

# 检查告警规则
promtool check rules /etc/prometheus/rules/asynq-alerts.yml

# 重新加载配置
curl -X POST http://localhost:9090/-/reload
```

### 2. 验证 Grafana 仪表盘

访问导入的仪表盘，确认：
- ✅ 所有面板都有数据显示
- ✅ 时间范围选择器正常工作
- ✅ 自动刷新设置为 10s
- ✅ 变量替换正确（如数据源选择）

### 3. 测试告警

#### 方法一：手动触发告警

```bash
# 模拟高失败率（停止 Worker）
docker stop worker-container

# 等待 5 分钟，应该收到 AsynqWorkerDown 告警
```

#### 方法二：使用 Prometheus UI

1. 访问 http://localhost:9090
2. 进入 "Alerts" 页面
3. 查看告警状态：
   - Pending（等待中）
   - Firing（触发中）
   - Inactive（未激活）

#### 方法三：发送测试告警

```bash
# 发送测试告警到 AlertManager
curl -X POST http://localhost:9093/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d '[{
    "labels": {
      "alertname": "TestAlert",
      "severity": "warning"
    },
    "annotations": {
      "summary": "Test alert",
      "description": "This is a test alert"
    }
  }]'
```

### 4. 监控指标验证

```bash
# 检查指标是否暴露
curl http://localhost:9091/metrics | grep asynq

# 应该看到类似输出：
# asynq_tasks_processed_total{task_type="job:daily_report",status="success"} 100
# asynq_queue_size{queue="jobs_default"} 5
# asynq_worker_concurrency 20
```

---

## 🔍 故障排查

### 问题 1：看不到指标数据

**原因：**
- Worker 未暴露 metrics 端点
- Prometheus 未正确抓取

**解决方案：**
```bash
# 检查 Worker 是否正常启动
docker logs worker-container | grep "metrics"

# 检查 Prometheus targets 状态
curl http://localhost:9090/api/v1/targets

# 手动测试 metrics 端点
curl http://localhost:9091/metrics
```

### 问题 2：告警未触发

**原因：**
- 告警规则语法错误
- 评估间隔设置过长
- 表达式不匹配实际数据

**解决方案：**
```bash
# 检查告警规则状态
curl http://localhost:9090/api/v1/rules

# 测试告警表达式
# 在 Prometheus UI 的 Graph 页面输入表达式
sum(rate(asynq_tasks_processed_total{status="failed"}[5m])) 
/ 
sum(rate(asynq_tasks_processed_total[5m]))
```

### 问题 3：仪表盘无数据

**原因：**
- 数据源配置错误
- 指标名称不匹配
- 时间范围设置不当

**解决方案：**
1. 在 Grafana 中检查数据源连接
2. 使用 Explore 功能测试查询
3. 调整时间范围为最近 1 小时

---

## 📝 维护建议

### 日常监控

1. **每日检查**
   - 查看失败率趋势
   - 检查队列积压峰值
   - 确认 Worker 健康状态

2. **每周回顾**
   - 分析告警频率
   - 优化阈值设置
   - 清理过期告警

3. **每月优化**
   - 审查告警规则有效性
   - 更新仪表盘布局
   - 性能基准对比

### 阈值调优

根据实际业务情况调整阈值：

```yaml
# 失败率阈值（默认 5%）
- expr: ... > 0.05  # 调整为 0.03 (3%) 或 0.10 (10%)

# 队列积压阈值（默认 1000）
- expr: asynq_queue_size > 1000  # 根据处理能力调整

# 延迟阈值（默认 60 秒）
- expr: ... > 60  # 根据 SLA 调整
```

---

## 🎯 最佳实践

1. **分级告警**
   - Critical：需要立即响应
   - Warning：工作时间处理
   - Info：仅记录不通知

2. **避免告警疲劳**
   - 设置合理的 `for` 持续时间
   - 使用 `group_wait` 聚合告警
   - 配置 `repeat_interval` 避免重复通知

3. **自动化响应**
   - 集成 PagerDuty/OpsGenie
   - 配置自动扩容策略
   - 实现告警自动恢复检测

4. **文档化**
   - 记录每个告警的处理流程
   - 建立 On-call 轮值制度
   - 定期演练故障场景

---

## 📚 参考资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [Asynq 监控最佳实践](https://github.com/hibiken/asynq/wiki/Monitoring)
- [AlertManager 配置指南](https://prometheus.io/docs/alerting/latest/configuration/)
