---
name: ddd-modeling-assistant
description: 智能DDD建模工具。通过对话式引导完成领域建模分析，自动生成符合Clean Architecture规范的DDD模块代码。
version: "1.0.0"
author: MathFun Team
tags: [ddd, domain-driven-design, module-development, clean-architecture, code-generation]
---

# DDD建模助手

## 功能概述

这是一个智能化的DDD模块开发工具，采用对话式引导方式，帮助开发者完整地创建符合领域驱动设计规范的业务模块。从需求理解、领域建模到代码实现，提供一站式开发体验。

## 核心能力

### 1. 智能需求分析
- 多轮对话式业务需求梳理
- 自动识别领域概念和业务实体
- 分析业务规则和约束条件
- 确定技术集成要求

### 2. 领域建模指导
- 聚合边界设计建议
- 实体、值对象识别
- 领域事件模式推荐
- 领域服务定位指导

### 3. 代码自动生成
- 标准DDD目录结构创建
- 领域层代码生成（聚合、实体、值对象、事件）
- 应用层服务实现
- 基础设施层集成（DAO、数据模型）
- 接口层代码生成（HTTP handlers）

### 4. 配套资源生成
- 数据库迁移文件
- API文档（Swagger格式）
- 单元测试代码模板
- 集成测试场景

## 使用场景

### 适用情况
- 新业务模块开发
- 现有模块重构
- DDD实践学习
- 团队DDD标准化
- 快速原型验证

### 不适用情况
- 简单的CRUD操作
- 已有成熟框架的模块
- 非业务逻辑密集型功能
- 纯技术组件开发

## 对话式开发流程

### 第一阶段：业务理解
```
Skills: "请描述您要开发的业务模块"
用户: "学习进度跟踪功能"

Skills: "请提供更多细节：
1. 目标用户群体是谁？
2. 需要跟踪哪些学习活动？
3. 期望达成什么业务目标？"
```

### 第二阶段：领域概念识别
```
Skills: "基于您的描述，我识别出以下概念：
🎯 核心实体：学习会话、知识点、用户
📊 值对象：掌握程度、学习时长、进度状态
⚡ 领域事件：学习开始、练习完成

这些概念是否准确？还需要补充吗？"
```

### 第三阶段：聚合边界设计
```
Skills: "建议的聚合设计：
📦 学习进度聚合：会话 + 统计数据 + 掌握度计算
👤 用户聚合：基本信息 + 学习偏好
📖 知识点聚合：知识结构 + 前置关系

这样的边界划分是否合理？"
```

### 第四阶段：业务规则确认
```
Skills: "识别的关键业务规则：
❌ 同一用户对同一知识点只能有一个活跃会话
📈 掌握度 = 完成练习数 / 总练习数 × 100%
⏰ 学习时长计算规则

这些规则是否完整准确？"
```

### 第五阶段：技术要求确认
```
Skills: "技术集成要求：
🔗 依赖用户服务获取用户信息
📡 通过事件通知学习报告服务
📚 生成REST API接口
💾 使用GORM持久化

这些技术要求是否符合预期？"
```

### 第六阶段：代码生成执行
```
Skills: "确认无误后，开始生成代码..."
[执行完整的代码生成流程]
```

## 快速开始

### 基本使用
```
# 启动DDD模块开发流程
/ddd-develop-module

# 重新开始某个阶段
/ddd-restart-analysis
/ddd-restart-design
/ddd-restart-generation
```

### 高级用法
```
# 基于现有分析生成代码
/ddd-generate-from-analysis

# 导出领域模型设计
/ddd-export-domain-model

# 验证生成代码质量
/ddd-validate-generated-code
```

## 配置说明

技能使用 `.qoder/skills/dddd-modeling-assistant/config.yaml` 进行配置：

```yaml
# DDD开发配置
ddd:
  # 架构风格
  architecture: "clean"  # clean | classic | hexagonal
  
  # 目录结构模板
  templates:
    domain: "internal/domain/{module}"
    application: "internal/application/{module}"
    infrastructure: "internal/infrastructure/persistence/gorm/{module}"
    interfaces: "internal/interfaces/http/{module}"
  
  # 代码生成选项
  generation:
    create_tests: true           # 生成测试代码
    create_migrations: true      # 生成迁移文件
    create_api_docs: true        # 生成API文档
    include_examples: true       # 包含示例代码

# 领域建模规则
domain_modeling:
  # 聚合设计原则
  aggregation_rules:
    - "一个聚合包含相关的业务概念"
    - "聚合边界内保证数据一致性"
    - "跨聚合通过领域事件通信"
  
  # 值对象识别规则
  value_object_patterns:
    - "具有业务含义的属性组合"
    - "不可变性要求"
    - "值相等性判断"

# 质量检查配置
quality_checks:
  enabled: true
  rules:
    - "聚合根必须有明确的业务标识"
    - "领域事件命名遵循过去时态"
    - "值对象实现相等性比较"
    - "实体具有唯一标识"

# 集成配置
integrations:
  database_migrator: "database-migrator"    # 数据库迁移技能
  api_doc_generator: "api-doc-generator"    # API文档生成技能
```

## 最佳实践

### 领域建模原则
1. **业务导向** - 从真实业务场景出发设计领域模型
2. **一致性边界** - 明确聚合的业务一致性范围
3. **显式建模** - 让隐式的业务规则显性化
4. **演进式设计** - 允许模型随业务理解深化而演进

### 对话技巧
- 提供具体而非抽象的业务场景
- 描述真实的用户行为和业务流程
- 明确业务规则的触发条件和约束
- 说明与其他系统的集成需求

### 团队协作
- 领域专家参与建模过程
- 定期回顾和优化领域模型
- 建立领域词汇表统一认知
- 文档化重要的设计决策

## 故障排除

### 常见问题

**领域概念识别不准确**
- 提供更详细的业务场景描述
- 举例说明典型的操作流程
- 明确业务规则的边界条件

**聚合边界设计争议**
- 回到业务一致性需求
- 考虑数据变更的原子性要求
- 分析跨边界通信的成本

**代码生成不符合预期**
- 检查配置文件设置
- 确认目录结构权限
- 验证依赖技能是否可用

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

## 版本历史

- v1.0.0 (2026-01-26): 初始版本发布
  - 完整的对话式DDD开发流程
  - 多阶段领域建模指导
  - 自动代码生成能力
  - 与现有技能集成支持

---
*本技能遵循Qoder Skills规范，专为MathFun项目DDD实践优化设计*