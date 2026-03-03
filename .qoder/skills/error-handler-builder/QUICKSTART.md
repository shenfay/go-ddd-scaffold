# 错误处理构建器快速入门

## 快速开始指南

### 1. 基本操作流程

#### 启动错误处理生成
```
# 生成基础错误处理代码
/error-build-handler

# 指定配置文件
/error-build-handler --config custom-error-config.yaml

# 指定输出目录
/error-build-handler --output ./backend/internal/pkg
```

#### 查看错误码列表
```
# 列出所有错误码
/error-list-codes

# 按分类筛选错误码
/error-list-codes --category BUSINESS

# 搜索特定错误码
/error-list-codes --filter "*USER*"
```

### 2. 常用命令速查

| 命令 | 用途 | 示例 |
|------|------|------|
| `/error-build-handler` | 生成错误处理代码 | `/error-build-handler --with-monitoring` |
| `/error-list-codes` | 查看错误码列表 | `/error-list-codes --category SYSTEM` |
| `/error-validate-config` | 验证配置文件 | `/error-validate-config --config errors.yaml` |
| `/error-generate-tests` | 生成测试代码 | `/error-generate-tests --coverage 90` |
| `/error-analyze-existing` | 分析现有错误使用 | `/error-analyze-existing --source ./backend` |
| `/error-generate-migration` | 生成迁移方案 | `/error-generate-migration --analysis analysis.json` |

### 3. 典型使用场景

#### 场景1：新项目错误处理初始化
```
# 1. 生成标准错误处理体系
/error-build-handler --with-monitoring --with-logging

# 2. 验证生成结果
/error-validate-config

# 3. 生成测试代码
/error-generate-tests --coverage-target 85
```

#### 场景2：现有系统错误处理重构
```
# 1. 分析现有错误使用情况
/error-analyze-existing-errors --source ./backend

# 2. 生成迁移方案
/error-generate-migration-plan --analysis-result analysis_report.json

# 3. 执行重构
/error-refactor-existing --migration-plan migration_plan.yaml
```

#### 场景3：团队错误规范统一
```
# 1. 使用团队标准配置
/error-build-handler --config team-error-standard.yaml --team-mode

# 2. 生成多语言支持
/error-build-handler --languages "zh,en,ja" --localization-file locales/

# 3. 集成监控告警
/error-build-handler --with-prometheus --with-alerting
```

## 高级功能

### 配置自定义
```
# 使用自定义错误分类
/error-build-handler --categories custom-categories.yaml

# 指定日志级别
/error-build-handler --log-level debug --include-stack-trace

# 启用特定功能
/error-build-handler --with-tracing --with-metrics --with-profiling
```

### 集成开发
```
# 与API端点生成器集成
/error-build-handler --integrate api-endpoint-generator

# 与DDD建模助手集成
/error-build-handler --from-ddd-model latest

# 生成完整的微服务错误处理
/error-build-handler --microservice-pattern --service-name user-service
```

### 性能优化
```
# 生成高性能错误处理
/error-build-handler --performance-optimized --pool-size 1000

# 启用异步处理
/error-build-handler --async-logging --buffer-size 10000

# 配置采样率
/error-build-handler --sampling-rate 0.1 --sample-errors
```

## 配置说明

技能使用 `.qoder/skills/error-handler-builder/config.yaml` 进行配置：

```yaml
error_handling:
  response_format:
    standard: true
    include_stack_trace: false
    include_error_id: true
    mask_sensitive_data: true
  
  error_categories:
    business_errors:
      prefix: "BUSINESS_"
      log_level: "warn"
    system_errors:
      prefix: "SYSTEM_"
      log_level: "error"
    validation_errors:
      prefix: "VALIDATION_"
      log_level: "info"

logging:
  enabled: true
  format: "json"
  levels: ["debug", "info", "warn", "error"]

integrations:
  api_endpoint_generator: "api-endpoint-generator"
  ddd_modeling_assistant: "ddd-modeling-assistant"
```

## 故障排除

### 常见问题

**错误码冲突**
```
/error-list-codes --conflicts-only
# 查看冲突的错误码
/error-resolve-conflicts --strategy prefix-module
# 解决命名冲突
```

**配置文件错误**
```
/error-validate-config --config my-errors.yaml
# 验证配置文件语法
```

**生成代码不符合预期**
```
/error-preview-generation
# 预览将要生成的内容
/error-build-handler --dry-run
# 试运行模式
```

**集成Skills调用失败**
```
/skills
# 检查依赖Skills是否可用
```

### 获取帮助
请参阅详细文档：
- SKILL.md - 主技能文档
- REFERENCE.md - 技术参考
- EXAMPLES.md - 使用示例

## 最佳实践

1. **建立错误码规范** - 制定团队统一的错误码命名和分类标准
2. **合理的日志级别** - 根据错误严重程度设置适当的日志级别
3. **用户友好的提示** - 错误消息要清晰易懂，避免技术术语
4. **安全敏感信息** - 不要在错误响应中泄露敏感信息
5. **监控告警配置** - 及时发现和处理系统异常

## 下一步

- 探索 [EXAMPLES.md](EXAMPLES.md) 了解详细的使用场景
- 查看 [REFERENCE.md](REFERENCE.md) 了解技术细节
- 检查 [SKILL.md](SKILL.md) 了解完整的技能文档

---
*如需支持，请联系 MathFun 开发团队*