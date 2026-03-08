# wire-auto-generator Skill

Wire 依赖注入自动化生成工具 - 一键生成 Wire 代码并诊断常见问题

## 🎯 功能特性

- ✅ 自动检测 Wire 配置完整性
- ✅ 一键生成 wire_gen.go 代码
- ✅ 智能诊断编译错误
- ✅ 提供最佳实践建议

## 🚀 使用方式

```bash
/qoder wire-auto-generator [目录路径]
```

## 📝 典型场景

### 场景 1: 首次使用 Wire

**用户需求**: "我需要配置 Wire 依赖注入"

**Skill 执行**:
1. 检查是否已有 injector.go
2. 提供模板示例
3. 运行 wire gen 生成代码
4. 验证编译通过

### 场景 2: 添加新的 Provider

**用户需求**: "我添加了新的 Service，需要更新 Wire 配置"

**Skill 执行**:
1. 检查 injector.go 中的 ProviderSet
2. 提示添加新的 Provider
3. 重新生成 wire_gen.go
4. 验证依赖注入正确性

### 场景 3: 修复编译错误

**用户需求**: "Wire 生成失败，有编译错误"

**Skill 执行**:
1. 分析错误信息
2. 识别常见错误类型
3. 提供修复方案
4. 自动或指导修复

## 🔍 常见错误诊断

### 错误 1: 类型未导出

```
❌ wire.Bind(new(Service), new(*impl))
原因：impl 首字母小写，未导出
修复：改为 Impl 或使用接口类型
```

### 错误 2: 重复定义

```
❌ userEventBusAdapter redeclared in this block
原因：在多个文件中定义了相同的辅助类型
修复：只在一个文件中保留定义
```

### 错误 3: 导入缺失

```
❌ undefined: transaction.UnitOfWork
原因：缺少 import 语句
修复：添加 "go-ddd-scaffold/internal/infrastructure/transaction"
```

### 错误 4: 循环依赖

```
❌ import cycle not allowed
原因：包之间存在循环导入
修复：重构代码结构，引入中间接口层
```

## 📚 最佳实践

### 1. Provider Set 组织

```go
// ✅ 推荐：按功能模块分组
var RepositorySet = wire.NewSet(
    repo.NewUserDAORepository,
    repo.NewTenantDAORepository,
)

var TransactionSet = wire.NewSet(
    transaction.NewGormUnitOfWork,
)

var AuthServiceSet = wire.NewSet(
    auth.NewJWTService,
    InitializeCasbinService,
)
```

### 2. Injector 函数命名

```go
// ✅ 推荐：清晰表达用途
func InitializeUserCommandService(...) service.UserCommandService
func InitializeTenantService(...) service.TenantService

// ❌ 不推荐：过于笼统
func InitializeService(...) interface{}
```

### 3. 依赖参数顺序

```go
// ✅ 推荐：从具体到抽象
func InitializeUserService(
    db *gorm.DB,                    // 基础设施
    logger *zap.Logger,             // 基础设施
    repo repository.UserRepository, // 仓储
    uow transaction.UnitOfWork,     // 事务
) service.UserService
```

### 4. 辅助类型定义

```go
// ✅ 推荐：定义在 injector.go 中
type userEventBusAdapter struct {
    bus *event.EventBus
}

func newUserEventBusAdapter(bus *event.EventBus) service.EventBus {
    return &userEventBusAdapter{bus: bus}
}
```

## 🛠️ 执行流程

### Step 1: 前置检查

```bash
✅ 检查 Go 环境
✅ 检查 Wire 是否安装
✅ 检查项目结构
```

### Step 2: 配置分析

```bash
✅ 扫描 injector.go
✅ 分析 Provider Set
✅ 检查类型导出情况
✅ 检测循环依赖
```

### Step 3: 代码生成

```bash
✅ 运行 wire gen
✅ 验证生成的 wire_gen.go
✅ 检查语法正确性
```

### Step 4: 编译验证

```bash
✅ go build 编译测试
✅ 检查导入完整性
✅ 验证依赖注入
```

### Step 5: 问题修复（如有）

```bash
✅ 诊断错误类型
✅ 提供修复建议
✅ 自动修复（如可能）
✅ 手动修复指导
```

## 📊 输出示例

### 成功场景

```
✅ Wire 配置检查通过
✅ 发现 3 个 Provider Set
✅ 发现 4 个 Injector 函数

🚀 开始生成 Wire 代码...
✅ wire: wrote wire_gen.go

🔍 编译验证...
✅ go build 编译通过

📦 生成完成！

新增文件:
  backend/internal/infrastructure/wire/wire_gen.go

修改文件:
  backend/internal/infrastructure/wire/injector.go

下一步建议:
1. 在 main.go 中使用 InitializeUserCommandService
2. 替换手动依赖初始化代码
3. 运行集成测试验证
```

### 失败场景

```
❌ Wire 配置检查发现问题

问题 1: 类型未导出
📍 位置：injector.go:34
❌ wire.Bind(new(transaction.UnitOfWork), new(*gormUnitOfWork))
💡 原因：gormUnitOfWork 首字母小写，未导出
🔧 修复方案:
   方案 A: 移除 wire.Bind (推荐)
      var TransactionSet = wire.NewSet(transaction.NewGormUnitOfWork)
   
   方案 B: 重命名为导出类型
      type GormUnitOfWork struct { ... }

问题 2: 重复定义
📍 位置：user.go:63, injector.go:108
❌ userEventBusAdapter redeclared in this block
💡 原因：在两个文件中定义了相同的类型
🔧 修复方案:
   删除 user.go 中的 userEventBusAdapter 定义 (第 63-77 行)

请修复上述问题后重新运行 /wire-auto-generator
```

## 🎓 学习资源

- [Wire 官方文档](https://github.com/google/wire)
- [DDD 架构重构计划](../../../docs/DDD_ARCHITECTURE_RESTRUCTURE_PLAN.md)
- [UnitOfWork 集成完成报告](../../../docs/UNIT_OF_WORK_INTEGRATION_COMPLETE.md)
- [Wire Provider 配置示例](../backend/internal/infrastructure/wire/providers.go)

## ⚙️ 配置选项

### 可选参数

```bash
--dry-run          # 仅检查，不生成代码
--verbose          # 详细输出模式
--auto-fix         # 自动修复可修复的问题
--output-dir DIR   # 指定输出目录（默认：当前目录）
```

### 环境变量

```bash
WIRE_VERBOSE=true  # 启用详细日志
WIRE_AUTO_FIX=true # 启用自动修复
```

## 🐛 故障排除

### Q1: wire 命令找不到

```bash
# 安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 或使用 go run
go run github.com/google/wire/cmd/wire@latest gen ./...
```

### Q2: 生成的代码不更新

```bash
# 清理缓存
rm backend/internal/infrastructure/wire/wire_gen.go

# 重新生成
wire gen ./internal/infrastructure/wire
```

### Q3: 编译出现大量错误

```bash
# 逐步排查
1. 检查 injector.go 语法
2. 确认所有导入路径正确
3. 验证类型是否导出
4. 运行 go mod tidy 清理依赖
```

## 📝 更新日志

### v1.0.0 (2026-03-08)
- ✅ 初始版本发布
- ✅ 支持 Wire 代码生成
- ✅ 支持常见错误诊断
- ✅ 支持最佳实践建议
- ✅ 集成 UnitOfWork 事务管理

---

**Skill 作者**: DDD Scaffold Team  
**最后更新**: 2026-03-08  
**适用版本**: Wire v0.6+
