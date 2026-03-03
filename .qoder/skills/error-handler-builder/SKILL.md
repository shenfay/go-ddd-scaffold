---
name: error-handler-builder
description: 统一错误处理构建器。自动生成标准化的错误处理代码，包含错误分类、统一响应格式和日志记录机制。
version: "1.0.0"
author: MathFun Team
tags: [error-handling, exception, logging, middleware, automation]
---

# 错误处理构建器

## 功能概述

这是一个智能化的错误处理构建工具，专为MathFun项目设计。它能够自动生成标准化的错误处理代码，包括错误分类、统一响应格式、日志记录和监控集成，确保整个系统的错误处理一致性。

## 核心能力

### 1. 智能错误分类
- 自动识别业务错误、系统错误和第三方错误
- 支持自定义错误类型和错误码
- 生成错误层级结构和继承关系
- 提供错误转换和映射机制

### 2. 统一响应格式
- 标准化的JSON错误响应结构
- 支持多语言错误消息
- 可配置的错误详情级别
- 前端友好的错误提示生成

### 3. 完整日志体系
- 结构化错误日志记录
- 支持不同级别的日志输出
- 集成链路追踪标识
- 自动关联请求上下文信息

### 4. 监控告警集成
- 错误统计和指标收集
- 告警规则自动生成
- 监控面板配置模板
- 性能影响分析

## 使用场景

### 适用情况
- 新项目统一错误处理规范
- 现有系统错误处理重构
- 微服务间错误处理标准化
- 团队错误处理最佳实践推广

### 不适用情况
- 已有成熟的错误处理体系
- 非Web应用场景
- 对错误处理有特殊要求的系统
- 需要高度定制化错误逻辑的场景

## 基本使用

### 快速开始
```
# 生成基础错误处理代码
/error-build-handler

# 基于项目配置生成
/error-build-handler --config error-config.yaml

# 指定输出目录
/error-build-handler --output ./custom/error/path
```

### 高级用法
```
# 生成带监控的错误处理
/error-build-handler --with-monitoring --with-alerting

# 生成多语言支持
/error-build-handler --languages "zh,en,ja"

# 生成完整的测试套件
/error-build-handler --with-tests --coverage-target 90
```

## 配置说明

技能使用 `.qoder/skills/error-handler-builder/config.yaml` 进行配置：

```yaml
error_handling:
  # 响应格式配置
  response_format:
    standard: true
    include_stack_trace: false
    include_error_id: true
    mask_sensitive_data: true
  
  # 错误分类配置
  error_categories:
    business_errors:
      prefix: "BUSINESS_"
      log_level: "warn"
    system_errors:
      prefix: "SYSTEM_"
      log_level: "error"
    third_party_errors:
      prefix: "EXTERNAL_"
      log_level: "error"
  
  # 日志配置
  logging:
    enabled: true
    format: "json"
    levels: ["debug", "info", "warn", "error"]
    include_context: true

# 集成配置
integrations:
  api_endpoint_generator: "api-endpoint-generator"
  ddd_modeling_assistant: "ddd-modeling-assistant"
  monitoring_system: "prometheus"
```

## 最佳实践

1. **统一错误码规范** - 建立团队统一的错误码命名和分类标准
2. **合理的日志级别** - 根据错误严重程度设置适当的日志级别
3. **用户友好的提示** - 错误消息要清晰易懂，避免技术术语
4. **安全敏感信息** - 不要在错误响应中泄露敏感信息
5. **监控告警配置** - 及时发现和处理系统异常

## 故障排除

### 常见问题

**错误码冲突**
- 检查错误码命名规范
- 确认不同模块间错误码唯一性
- 使用前缀避免命名冲突

**日志格式不一致**
- 统一日志配置模板
- 规范化日志字段命名
- 建立日志级别标准

**监控告警误报**
- 调整告警阈值设置
- 优化错误分类逻辑
- 完善告警抑制规则

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

---
*本技能遵循Qoder Skills规范，专为MathFun项目错误处理优化设计*