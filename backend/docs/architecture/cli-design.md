# CLI 工具架构设计文档

## 概述

本文档详细描述了 `go-ddd-scaffold` CLI 工具的架构设计，包括目录结构、核心组件、设计原则和扩展机制。

## 设计目标

1. **模块化**: 各组件职责清晰，低耦合高内聚
2. **可扩展**: 易于添加新的生成器和命令
3. **可测试**: 命令逻辑与生成逻辑分离，便于单元测试
4. **用户友好**: 清晰的命令结构和帮助信息
5. **配置灵活**: 支持配置文件、环境变量和命令行参数

## 目录结构

```
backend/cmd/cli/
├── main.go                          # 程序入口，只负责启动和版本信息
├── internal/                        # 内部包（不对外暴露）
│   ├── command/                     # 命令定义层
│   │   ├── root.go                  # 根命令和命令注册
│   │   ├── init.go                  # init 命令
│   │   ├── generate.go              # generate 命令及其子命令
│   │   ├── migrate_docs.go          # migrate 和 docs 命令
│   │   └── config.go                # 配置管理命令
│   └── generators/                  # 代码生成器层
│       ├── types.go                 # 生成器选项和类型定义
│       ├── dao_generator.go         # DAO 生成器实现
│       ├── entity_generator.go      # 实体生成器
│       ├── repository_generator.go  # Repository 生成器
│       ├── service_generator.go     # Service 生成器
│       ├── handler_generator.go     # Handler 生成器
│       ├── dto_generator.go         # DTO 生成器
│       ├── init_generator.go        # 项目初始化生成器
│       └── stubs.go                 # 其他生成器存根
└── templates/                       # 模板文件（可选，也可以嵌入到二进制中）
    ├── dao/
    ├── entity/
    └── ...
```

## 核心组件

### 1. main.go - 程序入口

**职责**:
- 创建根命令
- 设置版本信息
- 执行命令

**设计原则**:
- 不包含任何业务逻辑
- 不直接注册子命令（通过 `command.RegisterAll` 统一注册）
- 保持简洁，便于测试和维护

```go
func main() {
    rootCmd := newRootCmd()
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### 2. command 层 - 命令定义

**职责**:
- 定义 CLI 命令结构
- 解析命令行参数
- 验证用户输入
- 调用相应的生成器

**设计模式**: Command Pattern

每个命令都是一个独立的函数，返回 `*cobra.Command`，例如：
- `initCmd()` - 项目初始化
- `generateCmd()` - 代码生成根命令
- `generateDAOCmd()` - DAO 生成

**示例**:
```go
func generateDAOCmd() *cobra.Command {
    var opts generators.DAOOptions
    
    cmd := &cobra.Command{
        Use:   "dao [name]",
        Short: "Generate DAO layer",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.Name = args[0]
            generator := generators.NewDAOGenerator(opts)
            return generator.Generate()
        },
    }
    
    // 定义标志
    cmd.Flags().StringVarP(&opts.Fields, "fields", "f", "", "Field definitions")
    cmd.Flags().StringVarP(&opts.TableName, "table-name", "t", "", "Table name")
    
    return cmd
}
```

### 3. generators 层 - 代码生成器

**职责**:
- 实现具体的代码生成逻辑
- 处理模板渲染
- 文件写入

**设计模式**: Strategy Pattern

每个生成器实现统一的接口：
```go
type Generator interface {
    Generate() error
}
```

**生成器结构**:
```go
type DAOGenerator struct {
    opts DAOOptions
}

func (g *DAOGenerator) Generate() error {
    // 1. 创建目录
    // 2. 解析字段
    // 3. 生成接口文件
    // 4. 生成实现文件
    // 5. 生成模型文件
}
```

### 4. types.go - 类型定义

定义所有生成器的选项结构，用于在 command 层和 generators 层之间传递参数：

```go
type DAOOptions struct {
    Name          string
    Fields        string
    TableName     string
    OutputDir     string
    WithInterface bool
    WithCondition bool
}
```

## 数据流

```
用户输入 → Cobra 解析 → Command 层 → Generators 层 → 文件系统
           ↓              ↓            ↓
        参数验证       选项封装      模板渲染
```

### 示例流程：生成 DAO

1. **用户输入**: `go-ddd-scaffold generate dao user -f "username:string,email:string"`
2. **Cobra 解析**: 解析命令和标志
3. **Command 层**: 
   - 验证参数
   - 封装 `DAOOptions`
   - 创建 `DAOGenerator`
   - 调用 `Generate()`
4. **Generators 层**:
   - 解析字段字符串为 `[]Field`
   - 创建输出目录
   - 渲染模板
   - 写入文件
5. **输出**: 生成 3 个文件（interface, implementation, model）

## 设计原则

### 1. 关注点分离 (Separation of Concerns)

- **main.go**: 只负责启动
- **command 层**: 只负责命令解析和参数验证
- **generators 层**: 只负责代码生成

### 2. 依赖倒置 (Dependency Inversion)

command 层依赖 generators 层的抽象（接口），而不是具体实现：

```go
// Command 层
generator := generators.NewDAOGenerator(opts)
generator.Generate()  // 依赖接口，不关心实现细节
```

### 3. 单一职责 (Single Responsibility)

每个生成器只负责一种类型的代码生成：
- `DAOGenerator` → DAO 层代码
- `EntityGenerator` → 领域实体代码
- `ServiceGenerator` → 应用服务代码

### 4. 约定优于配置 (Convention over Configuration)

默认值遵循约定：
- 表名默认为实体名的复数形式 (`user` → `users`)
- 输出目录使用标准路径 (`internal/infrastructure/dao`)
- 字段命名自动转换（驼峰 ↔ 蛇形）

## 扩展机制

### 添加新的生成器

1. **在 `generators/types.go` 中添加选项**:
```go
type ProjectorOptions struct {
    Name       string
    Domain     string
    OutputDir  string
}
```

2. **创建生成器实现**:
```go
// generators/projector_generator.go
type ProjectorGenerator struct {
    opts ProjectorOptions
}

func NewProjectorGenerator(opts ProjectorOptions) *ProjectorGenerator {
    return &ProjectorGenerator{opts: opts}
}

func (g *ProjectorGenerator) Generate() error {
    // 实现生成逻辑
}
```

3. **在 `command/generate.go` 中添加命令**:
```go
func generateProjectorCmd() *cobra.Command {
    var opts generators.ProjectorOptions
    
    cmd := &cobra.Command{
        Use:   "projector [name]",
        Short: "Generate projector for read model",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.Name = args[0]
            generator := generators.NewProjectorGenerator(opts)
            return generator.Generate()
        },
    }
    
    // 定义标志...
    return cmd
}
```

4. **注册命令**:
```go
// command/generate.go
func generateCmd() *cobra.Command {
    cmd := &cobra.Command{Use: "generate"}
    
    cmd.AddCommand(generateEntityCmd())
    cmd.AddCommand(generateDAOCmd())
    // ... 其他命令
    cmd.AddCommand(generateProjectorCmd())  // 添加新命令
    
    return cmd
}
```

## 配置管理

### 配置优先级

```
命令行参数 > 环境变量 > 配置文件 > 默认值
```

### 配置文件格式

支持 YAML 格式的配置文件 (`~/.go-ddd-scaffold.yaml`):

```yaml
defaults:
  author: "Your Name"
  email: "your.email@example.com"
  license: "MIT"
  
generator:
  dao:
    output_dir: "internal/infrastructure/dao"
    with_interface: true
  entity:
    with_vo: false
    with_aggregate: true
```

## 测试策略

### 单元测试

针对 generators 层编写单元测试：

```go
func TestDAOGenerator_Generate(t *testing.T) {
    opts := DAOOptions{
        Name:      "user",
        Fields:    "username:string,email:string",
        TableName: "users",
    }
    
    generator := NewDAOGenerator(opts)
    err := generator.Generate()
    
    if err != nil {
        t.Fatalf("Generate failed: %v", err)
    }
    
    // 验证生成的文件
}
```

### 集成测试

测试完整的命令执行流程：

```go
func TestGenerateDAOCommand(t *testing.T) {
    cmd := generateDAOCmd()
    cmd.SetArgs([]string{"user", "-f", "username:string", "-t", "users"})
    
    err := cmd.Execute()
    if err != nil {
        t.Fatalf("Command execution failed: %v", err)
    }
}
```

## 最佳实践

### 1. 命令命名

- 使用动词开头：`init`, `generate`, `migrate`, `clean`
- 使用别名提高易用性：`gen` 作为 `generate` 的别名
- 保持命令简短且有描述性

### 2. 标志设计

- 常用标志使用短格式：`-f`, `-t`, `-o`
- 提供合理的默认值
- 使用 `BoolVarP` 而非 `BoolVar` 以支持短标志

### 3. 错误处理

- 使用 `RunE` 而非 `Run` 以支持错误返回
- 设置 `SilenceErrors: true` 避免重复输出
- 提供友好的错误消息

### 4. 输出格式

- 成功消息使用绿色 ✓ 前缀
- 警告消息使用黄色 ⚠ 前缀
- 错误消息使用红色 ✗ 前缀

## 性能优化

### 1. 并行生成

对于不相关的文件生成，可以并行执行：

```go
func (g *DAOGenerator) Generate() error {
    errChan := make(chan error, 3)
    
    go func() {
        errChan <- g.generateInterface(fields)
    }()
    
    go func() {
        errChan <- g.generateImplementation(fields)
    }()
    
    go func() {
        errChan <- g.generateModel(fields)
    }()
    
    // 等待所有 goroutine 完成
    for i := 0; i < 3; i++ {
        if err := <-errChan; err != nil {
            return err
        }
    }
    
    return nil
}
```

### 2. 模板缓存

预加载和缓存常用模板，避免重复解析：

```go
var templateCache = sync.Map{}

func loadTemplate(name, content string) (*template.Template, error) {
    if cached, ok := templateCache.Load(name); ok {
        return cached.(*template.Template), nil
    }
    
    tmpl, err := template.New(name).Parse(content)
    if err != nil {
        return nil, err
    }
    
    templateCache.Store(name, tmpl)
    return tmpl, nil
}
```

## 总结

当前的 CLI 工具设计遵循了现代 CLI 工具的最佳实践：

✅ **清晰的层次结构**: main → command → generators  
✅ **职责分离**: 每个层有明确的职责边界  
✅ **易于扩展**: 添加新功能不需要修改现有代码  
✅ **用户友好**: 清晰的命令结构和帮助信息  
✅ **可测试性强**: 各层独立，便于单元测试  

这种设计使得 CLI 工具可以随着项目发展而不断扩展，同时保持良好的可维护性。
