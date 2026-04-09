# Grafana 仪表盘管理

## 📊 仪表盘概览

项目采用模块化仪表盘设计，按功能域拆分为 6 个独立仪表盘：

| 仪表盘 | 文件 | 面板数 | 说明 |
|--------|------|--------|------|
| 📊 API Performance | `api-performance.json` | 14 | HTTP 请求性能监控 |
| 🔐 Authentication | `authentication.json` | 4 | 认证相关指标 |
| 🗄️ Database | `database.json` | 5 | 数据库查询性能 |
| ⚡ Redis | `redis.json` | 3 | Redis 命令监控 |
| 🚦 Rate Limiting | `rate-limiting.json` | 1 | 限流拒绝统计 |
| 📈 Business Metrics | `business-metrics.json` | 1 | 业务指标 |

## 📁 文件位置

```
grafana/
├── dashboards/                    # 子仪表盘（生产使用）
│   ├── api-performance.json
│   ├── authentication.json
│   ├── database.json
│   ├── redis.json
│   ├── rate-limiting.json
│   └── business-metrics.json
└── grafana-dashboard.json         # 完整仪表盘（备份参考）
```

## 🚀 快速导入

### 方式 1：批量导入（推荐）

```bash
# 导入所有仪表盘
./scripts/import-all-dashboards.sh --api-key YOUR_API_KEY

# 使用环境变量
export GRAFANA_API_KEY=YOUR_API_KEY
./scripts/import-all-dashboards.sh
```

### 方式 2：导入单个仪表盘

```bash
./scripts/import-all-dashboards.sh \
  --api-key YOUR_API_KEY \
  --file grafana/dashboards/api-performance.json
```

### 方式 3：Dry Run 预览

```bash
./scripts/import-all-dashboards.sh --dry-run
```

## ⚙️ 配置选项

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--url URL` | Grafana URL | `http://localhost:3000` |
| `--api-key KEY` | API Key（Admin 权限） | 环境变量 `GRAFANA_API_KEY` |
| `--dir DIR` | 仪表盘目录 | `grafana/dashboards` |
| `--file FILE` | 单个仪表盘文件 | - |
| `--dry-run` | 预览模式，不实际导入 | - |

## 🔑 创建 API Key

1. 登录 Grafana
2. 进入 **Configuration > API Keys**
3. 点击 **Add API Key**
4. 设置：
   - Name: `import-script`
   - Role: **Admin**（必须）
   - Time to live: 留空（永久）或自定义
5. 复制生成的 Key

## 🎨 布局设计

每个仪表盘采用针对性布局优化：

### API Performance
```
行 1: [QPS][Error][P50][P95][P99][预留]  (统计卡片)
行 2: [QPS Trend ───────────────────]  (全宽)
行 3: [HTTP Status] [QPS by Status]
行 4-6: [各类时序图 ────────────────]  (全宽)
```

### Authentication
```
行 1: [Success Rate][Success][Failures]
行 2: [Failures by Reason ──────────]  (全宽)
```

### Database
```
行 1: [Conn Pool][Open][Max]
行 2: [DB Queries by Table ─────────]  (全宽)
行 3: [Query Latency][Slow Queries]
```

### Redis
```
行 1: [Commands][Latency]
行 2: [Redis Errors ──────────────]  (全宽)
```

## 🔄 更新仪表盘

### 修改面板
1. 在 Grafana 中编辑仪表盘
2. 导出 JSON：Dashboard Settings > JSON Model
3. 保存到 `grafana/dashboards/` 对应文件

### 批量导出（从 Grafana）
```bash
# 使用 Grafana API 导出所有仪表盘
curl -H "Authorization: Bearer YOUR_KEY" \
  http://localhost:3000/api/dashboards/uids | jq '.[].uid' | while read uid; do
    curl -H "Authorization: Bearer YOUR_KEY" \
      "http://localhost:3000/api/dashboards/uid/$uid" > "grafana/dashboards/$uid.json"
  done
```

## 🛠️ 故障排查

### 导入失败：401 Unauthorized
- **原因**：API Key 无效或过期
- **解决**：重新创建 Admin 权限的 API Key

### 导入失败：403 Forbidden
- **原因**：API Key 权限不足
- **解决**：确保 Role 设置为 Admin

### JSON 格式错误
```bash
# 验证 JSON 格式
python3 -c "import json; json.load(open('grafana/dashboards/xxx.json'))"
```

### 面板显示 "No Data"
- 确认服务正常运行且有请求
- 检查 Prometheus 数据源配置
- 验证指标名称是否匹配

## 📝 最佳实践

1. **版本控制**：仪表盘 JSON 文件纳入 Git 管理
2. **环境隔离**：不同环境使用不同 API Key
3. **定期备份**：修改前导出 JSON 备份
4. **测试验证**：使用 `--dry-run` 预览变更
5. **文档同步**：修改后更新本文档

## 🔗 相关文档

- [监控指标配置](./monitoring-metrics.md)
- [部署指南](./deployment-guide.md)
- [故障排查手册](./troubleshooting.md)
