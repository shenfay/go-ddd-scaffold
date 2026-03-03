# DDD模块开发助手快速入门

## 快速开始指南

### 1. 基本操作流程

#### 启动DDD开发流程
```
# 启动完整的DDD模块开发
/ddd-develop-module

# 指定模块名称
/ddd-develop-module --module-name my_module

# 指定输出目录
/ddd-develop-module --output-dir ./custom/path
```

#### 阶段性操作命令
```
# 单独进行业务分析
/ddd-analyze-business "业务需求描述"

# 进行领域设计
/ddd-design-domain --entities "实体列表" --concepts "概念列表"

# 基于分析生成代码
/ddd-generate-code

# 验证设计质量
/ddd-validate-design
```

### 2. 常用命令速查

| 命令 | 用途 | 示例 |
|------|------|------|
| `/ddd-develop-module` | 启动完整开发流程 | `/ddd-develop-module --module-name user` |
| `/ddd-analyze-business` | 业务需求分析 | `/ddd-analyze-business "用户管理系统"` |
| `/ddd-design-domain` | 领域建模设计 | `/ddd-design-domain --entities "用户,订单"` |
| `/ddd-generate-code` | 生成代码 | `/ddd-generate-code --force` |
| `/ddd-validate-design` | 验证设计质量 | `/ddd-validate-design --strict` |
| `/ddd-export-model` | 导出领域模型 | `/ddd-export-model --format json` |
| `/ddd-list-sessions` | 查看会话历史 | `/ddd-list-sessions` |

### 3. 典型使用场景

#### 场景1：新功能模块开发
```
# 1. 启动开发流程
/ddd-develop-module --module-name learning_analytics

# 2. 跟随对话指引完成各阶段
[多轮对话交互...]

# 3. 确认后自动生成完整代码
Skills: "代码生成完成！"
```

#### 场景2：现有系统重构
```
# 1. 分析现有业务
/ddd-analyze-business "将CRUD系统重构为DDD架构"

# 2. 设计新的领域模型
/ddd-design-domain --existing-structure "分析现有代码结构"

# 3. 增量生成代码
/ddd-generate-code --incremental --target existing_module
```

#### 场景3：团队协作开发
```
# 团队成员A：负责业务分析
/ddd-analyze-business "复杂的电商订单流程"

# 团队成员B：负责技术设计
/ddd-design-domain --technical-requirements "高并发、分布式部署"

# 团队负责人：整合并生成
/ddd-generate-code --combine-sessions "analysis_session,design_session"
```

## 高级功能

### 配置自定义
```
# 使用自定义配置文件
/ddd-develop-module --config custom-ddd-config.yaml

# 指定架构风格
/ddd-develop-module --architecture hexagonal

# 启用特定生成选项
/ddd-generate-code --with-tests --with-docs --with-migrations
```

### 集成开发
```
# 与database-migrator集成
/ddd-develop-module --integrate database-migrator

# 与api-doc-generator集成
/ddd-develop-module --integrate api-doc-generator

# 批量处理多个模块
/ddd-develop-module --batch-mode --modules "user,order,payment"
```

### 学习模式
```
# 启用学习模式
/ddd-develop-module --learning-mode

# 查看详细解释
/ddd-explain-concept "什么是聚合根？"

# 获取最佳实践建议
/ddd-best-practices --topic "聚合设计"
```

## 配置说明

技能使用 `.qoder/skills/dddd-modeling-assistant/config.yaml` 进行配置：

```yaml
ddd:
  architecture: "clean"
  generation:
    create_tests: true
    create_migrations: true
    create_api_docs: true
  module_naming: "kebab-case"

integrations:
  database_migrator: "database-migrator"
  api_doc_generator: "api-doc-generator"
```

## 故障排除

### 常见问题

**对话理解偏差**
```
/ddd-restart-analysis
# 重新开始业务分析阶段
```

**领域设计争议**
```
/ddd-validate-design --detailed
# 获取详细的设计验证报告
```

**代码生成失败**
```
/ddd-generate-code --debug --verbose
# 启用调试模式查看详细错误信息
```

**集成Skills调用失败**
```
# 检查依赖Skills是否可用
/skills
# 确认database-migrator和api-doc-generator状态
```

### 获取帮助
请参阅详细文档：
- SKILL.md - 主技能文档
- REFERENCE.md - 技术参考
- EXAMPLES.md - 使用示例

## 最佳实践

1. **充分沟通业务需求** - 提供具体、详细的业务场景描述
2. **逐步确认设计** - 每个阶段都要仔细确认设计决策
3. **重视领域词汇** - 建立统一的业务术语理解
4. **验证生成代码** - 生成后仔细检查代码质量
5. **持续改进模型** - 根据实际使用反馈优化领域模型

## 下一步

- 探索 [EXAMPLES.md](EXAMPLES.md) 了解详细的使用场景
- 查看 [REFERENCE.md](REFERENCE.md) 了解技术细节
- 检查 [SKILL.md](SKILL.md) 了解完整的技能文档

---
*如需支持，请联系 MathFun 开发团队*