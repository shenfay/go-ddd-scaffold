# Wire Auto Generator Skill 使用指南

## 📦 安装说明

### 方式 1: 项目级 Skill（推荐）

Skill 已创建在项目目录：
```
.qoder/skills/wire-auto-generator/
```

Qoder 会自动识别项目级 Skill，无需额外配置。

### 方式 2: 用户级 Skill

如果希望全局使用此 Skill，请手动复制：

```bash
# 创建用户级 Skill 目录
mkdir -p ~/.qoder/skills/wire-auto-generator

# 复制 Skill 文件
cp -r .qoder/skills/wire-auto-generator/* ~/.qoder/skills/wire-auto-generator/

# 添加执行权限
chmod +x ~/.qoder/skills/wire-auto-generator/wire-generator.sh
```

---

## 🚀 快速开始

### 基本用法

在 Qoder 对话中输入：

```
/wire-auto-generator
```

或指定目录：

```
/wire-auto-generator ./backend/internal/infrastructure/wire
```

### 命令行工具

也可以直接使用脚本：

```bash
# 进入项目根目录
cd /Users/shenfay/Projects/ddd-scaffold

# 运行 Wire 生成
.qoder/skills/wire-auto-generator/wire-generator.sh

# 或使用参数
.qoder/skills/wire-auto-generator/wire-generator.sh \
  --output-dir ./backend/internal/infrastructure/wire \
  --verbose
```

---

## 📋 命令参数

### 选项参数

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--help` | `-h` | 显示帮助信息 | - |
| `--verbose` | `-v` | 详细输出模式 | false |
| `--dry-run` | `-d` | 仅检查，不生成代码 | false |
| `--auto-fix` | `-f` | 自动修复可修复的问题 | true |
| `--output-dir DIR` | `-o` | 指定 Wire 配置目录 | `./backend/internal/infrastructure/wire` |

### 使用示例

```bash
# 基本使用
.wire-generator.sh

# 详细模式
.wire-generator.sh -v

# 仅检查
.wire-generator.sh -d

# 指定目录
.wire-generator.sh -o ./backend/internal/infrastructure/wire

# 组合使用
.wire-generator.sh -v -f -o ./backend/internal/infrastructure/wire
```

---

## 🔍 工作流程

### Step 1: 环境检查

```
✅ 检查 Go 环境
✅ 检查 Wire 是否安装
✅ 检查项目结构
```

### Step 2: 配置分析

```
✅ 扫描 injector.go
✅ 分析 Provider Set
✅ 检查类型导出情况
✅ 检测循环依赖
```

### Step 3: 代码生成

```
🚀 开始生成 Wire 代码...
✅ wire: wrote wire_gen.go
```

### Step 4: 编译验证

```
✅ go build 编译通过
```

### Step 5: 问题修复（如有）

```
❌ 发现编译错误
💡 诊断错误类型
🔧 提供修复方案
```

---

## 🎯 典型场景

### 场景 1: 首次使用 Wire

**需求**: 项目中第一次配置 Wire

**步骤**:
1. 运行 `/wire-auto-generator`
2. Skill 会提示创建 injector.go
3. 参考提供的模板创建
4. 重新运行生成代码

**输出**:
```
⚠️  未找到 injector.go
💡 建议创建 injector.go 定义 Provider Set
✅ 模板已创建：./backend/internal/infrastructure/wire/injector.go.template

请参考模板创建 injector.go
```

### 场景 2: 添加新 Service

**需求**: 添加了新的 Application Service，需要更新 Wire 配置

**步骤**:
1. 在 injector.go 中添加新的 Provider
2. 运行 `/wire-auto-generator`
3. 自动生成新的 wire_gen.go

**示例**:
```go
// 在 injector.go 中添加
var UserServiceSet = wire.NewSet(
    service.NewUserService,
)

func InitializeUserService(...) service.UserService {
    wire.Build(
        RepositorySet,
        UserServiceSet,
    )
    return nil
}
```

**输出**:
```
✅ 发现 4 个 Provider Set
✅ 发现 5 个 Injector 函数
🚀 开始生成 Wire 代码...
✅ wire: wrote wire_gen.go
✅ go build 编译通过
```

### 场景 3: 修复编译错误

**需求**: Wire 生成后编译失败

**步骤**:
1. 运行 `/wire-auto-generator --auto-fix`
2. Skill 诊断错误并提供修复方案
3. 根据建议修复代码
4. 重新生成

**常见错误及修复**:

#### 错误 1: 类型未导出
```
❌ wire.Bind(new(Service), new(*impl))
💡 原因：impl 首字母小写
🔧 修复：移除 wire.Bind 或重命名为 Impl
```

#### 错误 2: 重复定义
```
❌ userEventBusAdapter redeclared in this block
💡 原因：在两个文件中定义
🔧 修复：删除 user.go 中的定义（第 63-77 行）
```

#### 错误 3: 导入缺失
```
❌ undefined: transaction.UnitOfWork
💡 原因：缺少 import
🔧 修复：添加 "go-ddd-scaffold/internal/infrastructure/transaction"
```

---

## 📚 最佳实践

### 1. Provider Set 组织

```go
// ✅ 推荐：按功能模块分组
var RepositorySet = wire.NewSet(
    repo.NewUserDAORepository,
    repo.NewTenantDAORepository,
    repo.NewTenantMemberDAORepository,
)

var TransactionSet = wire.NewSet(
    transaction.NewGormUnitOfWork,
)

var AuthServiceSet = wire.NewSet(
    auth.NewJWTService,
    InitializeCasbinService,
)

// 在 Injector 中使用
func InitializeUserService(...) service.UserService {
    wire.Build(
        RepositorySet,
        TransactionSet,
        AuthServiceSet,
        service.NewUserService,
    )
    return nil
}
```

### 2. 避免 wire.Bind

```go
// ❌ 不推荐：需要导出具体类型
wire.Bind(new(transaction.UnitOfWork), new(*transaction.GormUnitOfWork))

// ✅ 推荐：直接使用返回值
var TransactionSet = wire.NewSet(transaction.NewGormUnitOfWork)
```

### 3. 辅助类型定义

```go
// ✅ 推荐：定义在 injector.go 中
type userEventBusAdapter struct {
    bus *event.EventBus
}

func newUserEventBusAdapter(bus *event.EventBus) service.EventBus {
    return &userEventBusAdapter{bus: bus}
}

func (a *userEventBusAdapter) Publish(ctx context.Context, event service.DomainEvent) error {
    if event == nil {
        return nil
    }
    return a.bus.Publish(ctx, event)
}
```

### 4. 文件组织

```
backend/internal/infrastructure/wire/
├── injector.go          # 主要注入器配置（新建）
├── providers.go         # Provider 定义（已有）
├── providers_event.go   # 事件相关 Provider（已有）
├── providers_monitoring.go  # 监控相关 Provider（已有）
├── user.go              # User 模块专用配置（已有）
└── wire_gen.go          # 自动生成的代码（生成）
```

---

## 🐛 故障排除

### Q1: wire 命令找不到

**解决**:
```bash
# 安装 Wire
go install github.com/google/wire/cmd/wire@latest

# 或使用 go run（脚本会自动处理）
go run github.com/google/wire/cmd/wire@latest gen ./...
```

### Q2: 生成的代码不更新

**解决**:
```bash
# 清理缓存
rm backend/internal/infrastructure/wire/wire_gen.go

# 重新生成
wire gen ./internal/infrastructure/wire

# 或使用 Skill
/wire-auto-generator --auto-fix
```

### Q3: 大量编译错误

**解决**:
```bash
# 逐步排查
1. 检查 injector.go 语法是否正确
2. 确认所有导入路径正确
3. 验证类型是否导出（首字母大写）
4. 运行 go mod tidy 清理依赖
5. 使用 Skill 诊断：/wire-auto-generator -v
```

### Q4: 循环依赖错误

**解决**:
```bash
# 查看错误
import cycle not allowed

# 重构方案
1. 识别循环依赖的包
2. 引入中间接口层
3. 使用依赖倒置原则
4. 参考 DDD 架构分层
```

---

## 📊 输出示例

### 成功场景

```
🚀 Wire 自动化生成工具启动

💡 检查 Go 环境...
✅ Go 环境：go version go1.21.0 darwin/arm64

💡 检查 Wire 工具...
✅ Wire 已安装：wire version v0.6.0

💡 检查项目结构...
✅ 项目结构正常

💡 分析 Wire 配置：./backend/internal/infrastructure/wire
✅ 发现 3 个 Provider Set
✅ 发现 4 个 Injector 函数
✅ 检测到 wire.Bind 使用

🚀 开始生成 Wire 代码...
✅ wire: wrote wire_gen.go
✅ 生成 wire_gen.go (103 行)

💡 验证编译...
✅ 编译通过

🎉 Wire 代码生成完成！

下一步建议:
  1. 在 main.go 中使用生成的 Injector 函数
  2. 替换手动依赖初始化代码
  3. 运行集成测试验证
```

### 失败场景

```
🚀 Wire 自动化生成工具启动

💡 检查 Go 环境...
✅ Go 环境：go version go1.21.0 darwin/arm64

💡 检查 Wire 工具...
⚠️  Wire 未安装，将使用 go run 方式

💡 检查项目结构...
✅ 项目结构正常

💡 分析 Wire 配置：./backend/internal/infrastructure/wire
✅ 发现 2 个 Provider Set
✅ 发现 3 个 Injector 函数

🚀 开始生成 Wire 代码...
❌ Wire 生成失败

错误详情:
wire: /Users/shenfay/Projects/ddd-scaffold/backend/internal/infrastructure/wire/injector.go:39:46: undefined: auth.JwtService

💡 诊断编译错误...
❌ 发现未导出的类型
💡 解决方案：将类型名首字母大写或使用接口类型

尝试自动修复...
💡 建议：在 injector.go 中修改
   var AuthServiceSet = wire.NewSet(auth.NewJWTService)
   // 移除 wire.Bind

请修复上述问题后重新运行 /wire-auto-generator
```

---

## 🎓 学习资源

- [Wire 官方文档](https://github.com/google/wire)
- [DDD 架构重构计划](../../../docs/DDD_ARCHITECTURE_RESTRUCTURE_PLAN.md)
- [UnitOfWork 集成完成报告](../../../docs/UNIT_OF_WORK_INTEGRATION_COMPLETE.md)
- [Wire Provider 配置示例](../backend/internal/infrastructure/wire/providers.go)
- [injector.go 示例](../backend/internal/infrastructure/wire/injector.go)

---

## 🔄 版本历史

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
**许可证**: MIT
