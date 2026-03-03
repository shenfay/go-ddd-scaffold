# DDD模块开发助手技术参考

## 目录结构规范

```
.qoder/skills/dddd-modeling-assistant/
├── SKILL.md              # 主技能定义文件（必需，大写）
├── REFERENCE.md          # 技术参考文档（大写）
├── EXAMPLES.md           # 使用示例文档（大写）
├── QUICKSTART.md         # 快速入门指南
├── config.yaml           # 配置文件
├── scripts/              # 辅助脚本目录
│   └── helper.sh         # 主要辅助脚本
└── templates/            # 模板文件目录
    ├── domain.tmpl       # 领域层模板
    ├── application.tmpl  # 应用层模板
    └── infrastructure.tmpl  # 基础设施层模板
```

## 核心组件详解

### 1. 对话引擎架构

#### 多轮对话状态管理
```go
type ConversationState struct {
    Stage           ConversationStage  // 当前对话阶段
    DomainModel     *DomainModel       // 领域模型缓存
    UserInputs      []UserInput        // 用户输入历史
    GeneratedCode   *GeneratedCode     // 已生成代码信息
    ValidationIssues []ValidationIssue // 验证问题列表
}

type ConversationStage string
const (
    StageBusinessUnderstanding  ConversationStage = "business_understanding"
    StageDomainConcepts         ConversationStage = "domain_concepts"
    StageAggregationDesign      ConversationStage = "aggregation_design"
    StageBusinessRules          ConversationStage = "business_rules"
    StageTechnicalRequirements  ConversationStage = "technical_requirements"
    StageCodeGeneration         ConversationStage = "code_generation"
)
```

#### 领域概念识别引擎
```go
type DomainAnalyzer struct {
    EntityRecognizer     *EntityRecognizer
    ValueObjectDetector  *ValueObjectDetector
    EventIdentifier      *EventIdentifier
    ServiceLocator       *ServiceLocator
}

func (da *DomainAnalyzer) Analyze(input string) *DomainAnalysis {
    // 实体识别
    entities := da.EntityRecognizer.Recognize(input)
    
    // 值对象检测
    valueObjects := da.ValueObjectDetector.Detect(input)
    
    // 事件识别
    events := da.EventIdentifier.Identify(input)
    
    return &DomainAnalysis{
        Entities:     entities,
        ValueObjects: valueObjects,
        Events:       events,
        Services:     da.ServiceLocator.Locate(input),
    }
}
```

### 2. 代码生成引擎

#### 模板系统架构
```go
type CodeGenerator struct {
    TemplateEngine  *template.Template
    OutputFormatter *OutputFormatter
    Validator       *CodeValidator
}

type GenerationContext struct {
    ModuleName      string
    DomainModel     *DomainModel
    Architecture    string  // "clean" | "classic" | "hexagonal"
    OutputDir       string
    Config          *GenerationConfig
}
```

#### 目录结构生成器
```go
func (cg *CodeGenerator) GenerateDirectoryStructure(ctx *GenerationContext) error {
    dirs := []string{
        fmt.Sprintf("internal/domain/%s", ctx.ModuleName),
        fmt.Sprintf("internal/domain/%s/aggregate", ctx.ModuleName),
        fmt.Sprintf("internal/domain/%s/entity", ctx.ModuleName),
        fmt.Sprintf("internal/domain/%s/valueobject", ctx.ModuleName),
        fmt.Sprintf("internal/domain/%s/event", ctx.ModuleName),
        fmt.Sprintf("internal/domain/%s/service", ctx.ModuleName),
        fmt.Sprintf("internal/application/%s/service", ctx.ModuleName),
        fmt.Sprintf("internal/infrastructure/persistence/gorm/%s", ctx.ModuleName),
        fmt.Sprintf("internal/interfaces/http/%s", ctx.ModuleName),
    }
    
    for _, dir := range dirs {
        if err := os.MkdirAll(filepath.Join(ctx.OutputDir, dir), 0755); err != nil {
            return err
        }
    }
    return nil
}
```

### 3. 领域建模规则引擎

#### 聚合边界验证规则
```yaml
aggregation_rules:
  - name: "业务一致性原则"
    description: "聚合边界内的对象应保持业务一致性"
    validation: |
      // 检查聚合根是否包含相关业务概念
      len(aggregate.Entities) > 1 &&
      aggregate.HasSingleRoot() &&
      aggregate.BoundaryEnclosesRelatedConcepts()
  
  - name: "事务边界原则"
    description: "聚合是数据一致性的事务边界"
    validation: |
      // 检查聚合变更的原子性要求
      aggregate.RequiresAtomicChanges() &&
      aggregate.MinimizesCrossBoundaryUpdates()
  
  - name: "领域事件通信原则"
    description: "跨聚合通信通过领域事件实现"
    validation: |
      // 检查跨聚合引用是否通过事件
      !aggregate.DirectlyReferencesOtherAggregates() ||
      aggregate.UsesDomainEventsForCommunication()
```

#### 值对象识别规则
```yaml
value_object_rules:
    - name: "业务含义原则"
      description: "值对象应具有明确的业务含义"
      pattern: "\\b(金额|数量|状态|比率|时间范围)\\b"
      
    - name: "不可变性原则"
      description: "值对象应该是不可变的"
      validation: |
        // 检查是否只有getter方法
        valueObject.HasOnlyGetters() &&
        !valueObject.HasSetters()
      
    - name: "相等性原则"
      description: "值对象通过值相等性判断"
      validation: |
        // 检查是否实现了相等性比较
        valueObject.ImplementsEquals() &&
        valueObject.ImplementsHashCode()
```

### 4. 配置管理系统

#### 完整配置结构
```yaml
ddd:
  architecture: "clean"
  module_naming:
    convention: "kebab-case"  # 模块命名规范
    prefix: ""               # 模块前缀
    suffix: ""               # 模块后缀
  
  directory_structure:
    domain: "internal/domain/{module}"
    application: "internal/application/{module}"
    infrastructure: "internal/infrastructure/persistence/gorm/{module}"
    interfaces: "internal/interfaces/http/{module}"
  
  code_generation:
    templates:
      domain: "templates/domain.tmpl"
      application: "templates/application.tmpl"
      infrastructure: "templates/infrastructure.tmpl"
    
    options:
      create_unit_tests: true
      create_integration_tests: false
      generate_swagger_docs: true
      create_database_migrations: true

# 集成配置
integrations:
  database_migrator:
    skill_name: "database-migrator"
    auto_invoke: true
    migration_template: "ddd_module"
  
  api_doc_generator:
    skill_name: "api-doc-generator"
    auto_invoke: true
    scan_paths: ["internal/interfaces/http"]
```

## API参考

### 核心命令

#### /ddd-develop-module
启动完整的DDD模块开发流程
```
/ddd-develop-module [--module-name {name}] [--output-dir {path}]
```

#### /ddd-analyze-business
单独进行业务需求分析
```
/ddd-analyze-business "{business_description}"
```

#### /ddd-design-domain
进行领域建模设计
```
/ddd-design-domain --entities "{entity_list}" --concepts "{concept_list}"
```

#### /ddd-generate-code
基于已完成的分析生成代码
```
/ddd-generate-code [--from-session {session_id}] [--force]
```

### 辅助命令

#### /ddd-export-model
导出领域模型设计文档
```
/ddd-export-model [--format json|yaml|markdown] [--output {file}]
```

#### /ddd-validate-design
验证领域设计质量
```
/ddd-validate-design [--strict] [--report {file}]
```

#### /ddd-list-sessions
列出当前会话历史
```
/ddd-list-sessions [--active-only]
```

## 领域建模最佳实践

### 1. 实体识别模式
```go
// 实体识别关键词模式
var entityPatterns = []string{
    `\b(用户|订单|商品|学习会话|知识点)\b`,
    `\b(管理|操作|处理)\b.*\b(对象|实体)\b`,
    `\b(具有|包含|拥有).*\b(唯一标识|ID)\b`,
}

// 实体特征检查清单
type EntityChecklist struct {
    HasIdentity           bool  // 是否有唯一标识
    HasLifecycle          bool  // 是否有生命周期
    MutableState          bool  // 状态是否可变
    BusinessBehavior      bool  // 是否有业务行为
    ReferenceByIdentity   bool  // 是否通过标识引用
}
```

### 2. 值对象识别模式
```go
// 值对象识别关键词模式
var valueObjectPatterns = []string{
    `\b(金额|数量|状态|比率|时间范围|地址)\b`,
    `\b(描述|计算|度量)\b.*\b(结果|值)\b`,
    `\b(不可变|值相等|无标识)\b`,
}

// 值对象特征检查清单
type ValueObjectChecklist struct {
    Immutable             bool  // 是否不可变
    ValueEquality         bool  // 是否基于值相等
    NoIdentity            bool  // 是否无唯一标识
    SelfValidating        bool  // 是否自我验证
    Behaviorless          bool  // 是否无业务行为
}
```

### 3. 聚合设计原则
```go
// 聚合设计检查清单
type AggregateDesignChecklist struct {
    SingleRoot            bool  // 是否有单一聚合根
    BusinessConsistency   bool  // 是否保证业务一致性
    BoundaryClarity       bool  // 边界是否清晰
    EventCommunication    bool  // 跨聚合是否通过事件
    SizeAppropriate       bool  // 聚合大小是否合适
}
```

## 错误处理机制

### 错误分类体系
1. **对话理解错误** - 无法正确解析用户输入
2. **领域建模错误** - 领域设计违反DDD原则
3. **代码生成错误** - 模板渲染或文件写入失败
4. **集成调用错误** - 调用其他Skills失败
5. **配置错误** - 配置文件缺失或格式错误

### 错误处理策略
- 提供清晰的错误信息和修复建议
- 支持对话流程回退和重试
- 记录详细的错误日志
- 区分致命错误和可恢复错误

## 性能优化

### 对话状态优化
- 会话状态压缩存储
- LRU缓存常用分析结果
- 异步处理耗时操作
- 进度状态实时反馈

### 代码生成优化
- 模板预编译
- 并行文件生成
- 增量更新支持
- 生成进度可视化

## 安全考虑

### 访问控制
- 会话隔离机制
- 敏感信息过滤
- 输出目录权限控制
- 配置文件加密

### 数据保护
- 用户输入脱敏处理
- 生成代码安全扫描
- 依赖包安全检查
- 环境变量安全处理

## 扩展机制

### 插件架构
支持自定义扩展：
- 自定义分析规则
- 额外的代码模板
- 第三方工具集成
- 自定义验证规则

### 集成能力
- 与其他Skills无缝集成
- CI/CD流水线集成
- 版本控制系统集成
- 项目管理工具集成

---
*本文档遵循Qoder Skills技术规范，定期更新以反映最新功能和最佳实践*