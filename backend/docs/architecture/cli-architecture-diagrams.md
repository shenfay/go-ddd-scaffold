# CLI 工具架构图

## 整体架构

```mermaid
graph TB
    User[用户] --> Main[main.go 入口层]
    Main --> Command[command 层<br/>命令解析]
    Command --> Generators[generators 层<br/>代码生成]
    Generators --> FileSystem[文件系统]
    
    subgraph "入口层"
        Main
    end
    
    subgraph "命令解析层"
        Command
        Root[root.go<br/>统一注册]
        Init[init.go]
        Generate[generate.go]
        Migrate[migrate_docs.go]
    end
    
    subgraph "代码生成层"
        Generators
        Types[types.go<br/>选项定义]
        DAO[dao_generator.go]
        Entity[entity_generator.go]
        Repository[repository_generator.go]
        Service[service_generator.go]
    end
    
    Main --> Root
    Root --> Init
    Root --> Generate
    Root --> Migrate
    
    Generate --> DAO
    Generate --> Entity
    Generate --> Repository
    Generate --> Service
```

## 数据流

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as Cobra 解析器
    participant CMD as Command 层
    participant GEN as Generators 层
    participant FS as 文件系统
    
    U->>C: go-ddd-scaffold generate dao user -f "..."
    C->>CMD: 解析参数和标志
    CMD->>CMD: 验证输入
    CMD->>GEN: NewDAOGenerator(opts)
    GEN->>GEN: parseFields()
    GEN->>GEN: generateInterface()
    GEN->>GEN: generateImplementation()
    GEN->>GEN: generateModel()
    GEN->>FS: 写入文件
    FS-->>GEN: 成功
    GEN-->>CMD: nil (成功)
    CMD-->>U: ✓ Generated successfully
```

## 命令层次结构

```mermaid
graph LR
    Root[go-ddd-scaffold] --> Init[init<br/>初始化项目]
    Root --> Generate[generate<br/>代码生成]
    Root --> Migrate[migrate<br/>数据库迁移]
    Root --> Docs[docs<br/>文档生成]
    Root --> Clean[clean<br/>清理文件]
    Root --> Version[version<br/>版本信息]
    
    Generate --> Entity[entity<br/>领域实体]
    Generate --> DAO[dao<br/>数据访问层]
    Generate --> Repo[repository<br/>仓储层]
    Generate --> Service[service<br/>应用服务]
    Generate --> Handler[handler<br/>CQRS 处理器]
    Generate --> DTO[dto<br/>数据传输对象]
    
    Migrate --> Up[up<br/>运行迁移]
    Migrate --> Down[down<br/>回滚迁移]
    Migrate --> Create[create<br/>创建迁移]
```

## 类型依赖关系

```mermaid
graph TB
    subgraph "Options 类型"
        DAOOpts[DAOOptions]
        EntityOpts[EntityOptions]
        RepoOpts[RepositoryOptions]
        ServiceOpts[ServiceOptions]
    end
    
    subgraph "Generators"
        DAOG gen[DAOGenerator]
        EntityG gen[EntityGenerator]
        RepoG gen[RepositoryGenerator]
        ServiceG gen[ServiceGenerator]
    end
    
    subgraph "Commands"
        DAOCmd[generateDAOCmd]
        EntityCmd[generateEntityCmd]
        RepoCmd[generateRepositoryCmd]
        ServiceCmd[generateServiceCmd]
    end
    
    DAOOpts --> DAOG gen
    EntityOpts --> EntityG en
    RepoOpts --> RepoG en
    ServiceOpts --> ServiceG en
    
    DAOCmd --> DAOG gen
    EntityCmd --> EntityG en
    RepoCmd --> RepoG en
    ServiceCmd --> ServiceG en
```

## 文件组织结构

```
cmd/cli/
├── main.go (35 lines)
│   └── 职责：启动程序，不包含业务逻辑
│
├── internal/
│   │
│   ├── command/ (命令定义层)
│   │   ├── root.go (~50 lines)
│   │   │   └── RegisterAll(): 统一注册所有命令
│   │   │
│   │   ├── init.go (~45 lines)
│   │   │   └── initCmd(): 项目初始化命令
│   │   │
│   │   ├── generate.go (~160 lines)
│   │   │   ├── generateCmd(): 生成根命令
│   │   │   ├── generateEntityCmd()
│   │   │   ├── generateDAOCmd()
│   │   │   ├── generateRepositoryCmd()
│   │   │   ├── generateServiceCmd()
│   │   │   ├── generateHandlerCmd()
│   │   │   └── generateDTOCmd()
│   │   │
│   │   ├── migrate_docs.go (~65 lines)
│   │   │   ├── migrateCmd()
│   │   │   └── docsCmd()
│   │   │
│   │   └── config.go (TODO)
│   │       └── configCmd(): 配置管理
│   │
│   └── generators/ (代码生成层)
│       ├── types.go (~75 lines)
│       │   └── 所有生成器的 Options 定义
│       │
│       ├── dao_generator.go (~580 lines)
│       │   ├── DAOGenerator 结构体
│       │   ├── NewDAOGenerator()
│       │   ├── Generate()
│       │   ├── parseFields()
│       │   ├── generateInterface()
│       │   ├── generateImplementation()
│       │   └── generateModel()
│       │
│       ├── entity_generator.go (TODO)
│       ├── repository_generator.go (TODO)
│       ├── service_generator.go (TODO)
│       ├── handler_generator.go (TODO)
│       ├── dto_generator.go (TODO)
│       │
│       ├── init_generator.go (~25 lines)
│       │   └── InitGenerator 存根
│       │
│       └── stubs.go (~30 lines)
│           └── 其他生成器存根
│
└── templates/ (可选)
    ├── dao/
    ├── entity/
    └── ...
```

## 命令执行流程示例

### 生成 DAO 的完整流程

```bash
go-ddd-scaffold generate dao user -f "username:string,email:string" -t users
```

**步骤分解**:

```
1. main.go
   └─ newRootCmd() → 创建根命令

2. command/root.go
   └─ RegisterAll() → 注册所有子命令
      └─ generateCmd() → 添加 generate 命令
         └─ generateDAOCmd() → 添加 dao 子命令

3. Cobra 框架
   └─ 解析命令行参数
      ├─ args[0] = "user"
      ├─ flags["fields"] = "username:string,email:string"
      └─ flags["table-name"] = "users"

4. command/generate.go
   └─ generateDAOCmd().RunE()
      ├─ opts.Name = "user"
      ├─ opts.Fields = "username:string,email:string"
      ├─ opts.TableName = "users"
      └─ NewDAOGenerator(opts)

5. generators/dao_generator.go
   └─ DAOGenerator.Generate()
      ├─ parseFields() → []Field
      │   ├─ Field{Name:"Username", Type:"string", GoType:"string"}
      │   └─ Field{Name:"Email", Type:"string", GoType:"string"}
      │
      ├─ generateInterface(fields)
      │   └─ 渲染模板 → user_dao.go
      │
      ├─ generateImplementation(fields)
      │   └─ 渲染模板 → user_dao_impl.go
      │
      └─ generateModel(fields)
          └─ 渲染模板 → user_model.go

6. 文件系统
   └─ 写入 3 个文件
      ├─ internal/infrastructure/dao/user_dao.go
      ├─ internal/infrastructure/dao/user_dao_impl.go
      └─ internal/infrastructure/dao/user_model.go
```

## 设计模式应用

### 1. Command Pattern（命令模式）

```mermaid
classDiagram
    class Command {
        <<interface>>
        Execute() error
    }
    
    class InitCommand {
        Execute() error
    }
    
    class GenerateCommand {
        Execute() error
    }
    
    class MigrateCommand {
        Execute() error
    }
    
    Command <|-- InitCommand
    Command <|-- GenerateCommand
    Command <|-- MigrateCommand
```

### 2. Strategy Pattern（策略模式）

```mermaid
classDiagram
    class Generator {
        <<interface>>
        Generate() error
    }
    
    class DAOGenerator {
        Generate() error
    }
    
    class EntityGenerator {
        Generate() error
    }
    
    class RepositoryGenerator {
        Generate() error
    }
    
    Generator <|-- DAOGenerator
    Generator <|-- EntityGenerator
    Generator <|-- RepositoryGenerator
```

### 3. Factory Pattern（工厂模式）

```mermaid
classDiagram
    class DAOGenerator {
        +NewDAOGenerator(opts) DAOGenerator
        +Generate() error
    }
    
    class EntityGenerator {
        +NewEntityGenerator(opts) EntityGenerator
        +Generate() error
    }
    
    class RepositoryGenerator {
        +NewRepositoryGenerator(opts) RepositoryGenerator
        +Generate() error
    }
```

## 接口设计

### Generator 接口

```go
type Generator interface {
    Generate() error
}
```

所有生成器都实现这个接口：
- `DAOGenerator`
- `EntityGenerator`
- `RepositoryGenerator`
- `ServiceGenerator`
- `HandlerGenerator`
- `DTOGenerator`
- `InitGenerator`

### 统一的调用方式

```go
// Command 层调用
generator := generators.NewDAOGenerator(opts)
err := generator.Generate()  // 所有生成器使用相同的接口

if err != nil {
    return err
}
fmt.Println("✓ Generation completed")
```

## 配置系统架构

```mermaid
graph TB
    subgraph "配置来源"
        CLI[命令行参数]
        ENV[环境变量]
        File[配置文件 YAML]
        Default[默认值]
    end
    
    subgraph "配置处理"
        Merge[合并配置]
        Validate[验证配置]
        Apply[应用配置]
    end
    
    subgraph "配置使用"
        Cmd[Command 层]
        Gen[Generators 层]
    end
    
    CLI --> Merge
    ENV --> Merge
    File --> Merge
    Default --> Merge
    
    Merge --> Validate
    Validate --> Apply
    Apply --> Cmd
    Apply --> Gen
```

## 测试架构

```mermaid
graph TB
    subgraph "单元测试"
        UT1[测试 DAOGenerator]
        UT2[测试 EntityGenerator]
        UT3[测试 parseFields]
    end
    
    subgraph "集成测试"
        IT1[测试 generate dao 命令]
        IT2[测试 init 命令]
        IT3[测试完整流程]
    end
    
    subgraph "E2E 测试"
        E2E1[测试真实项目生成]
        E2E2[测试批量生成]
    end
    
    UT1 --> generators 层
    UT2 --> generators 层
    UT3 --> generators 层
    
    IT1 --> command 层
    IT2 --> command 层
    IT3 --> command 层 + generators 层
    
    E2E1 --> 完整 CLI
    E2E2 --> 完整 CLI
```

## 性能优化策略

### 1. 并行生成

```mermaid
graph LR
    Start[开始生成] --> Parallel{并行执行}
    Parallel --> Gen1[生成 Interface]
    Parallel --> Gen2[生成 Implementation]
    Parallel --> Gen3[生成 Model]
    Gen1 --> Wait[等待完成]
    Gen2 --> Wait
    Gen3 --> Wait
    Wait --> End[生成完成]
```

### 2. 模板缓存

```mermaid
graph TB
    Request[请求模板] --> Cache{缓存命中？}
    Cache -->|是 | Return[返回缓存模板]
    Cache -->|否 | Parse[解析模板]
    Parse --> Store[存入缓存]
    Store --> Return
```

## 扩展点设计

```mermaid
graph TB
    subgraph "内置功能"
        Builtin[内置命令和生成器]
    end
    
    subgraph "扩展机制"
        Plugin[插件系统]
        Template[自定义模板]
        Hook[钩子函数]
    end
    
    subgraph "第三方扩展"
        Custom[自定义生成器]
        External[外部模板]
        Script[脚本钩子]
    end
    
    Builtin --> Plugin
    Plugin --> Custom
    Builtin --> Template
    Template --> External
    Builtin --> Hook
    Hook --> Script
```

## 错误处理流程

```mermaid
sequenceDiagram
    participant U as 用户
    participant C as Command 层
    participant G as Generators 层
    participant E as 错误处理器
    
    U->>C: 执行命令
    C->>G: 调用 Generate()
    
    alt 成功
        G-->>C: nil
        C-->>U: ✓ 成功消息
    else 失败
        G->>E: 返回错误
        E->>E: 格式化错误
        E-->>C: 友好错误消息
        C-->>U: ✗ 错误提示
    end
```

## 总结

通过以上多维度架构图，我们可以清晰地看到：

1. **分层清晰**: 入口层 → 命令层 → 生成层
2. **职责单一**: 每层只负责一件事
3. **易于扩展**: 遵循开闭原则
4. **可测试性强**: 各层独立，便于测试
5. **用户体验好**: 清晰的命令结构

这是一个**企业级、专业化、可扩展**的 CLI 工具架构设计。
