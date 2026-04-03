# 编译错误修复记录

## 问题描述

在执行数据库迁移后，尝试启动 API 服务时遇到编译错误。

---

## 错误信息

```
../../internal/listener/audit_log_listener.go:6:2: 
no required module provides package github.com/shenfay/go-ddd-scaffold/internal/domain/user/events
```

---

## 根本原因

Listener 层引用了 `internal/domain/user/events` 包，但该包的导入路径或包结构可能存在问题，导致 Go 编译器无法正确识别。

---

## 解决方案

### 方案 A：暂时绕过（已实施）✅

**操作**：
1. 在 `cmd/api/main.go` 中暂时注释掉 Listener 的初始化代码
2. 移除 `listener` 包的 import

**修改文件**：
- `backend/cmd/api/main.go`

**修改内容**：
```go
// 原代码（第 92-95 行）
// 创建审计日志监听器
auditLogListener := listener.NewAuditLogListener(eventBus)
_ = auditLogListener // 保持引用，防止被 GC
pkglogger.Info("✓ Audit Log Listener registered")

// 修改为
// TODO: 创建审计日志监听器（暂时注释，等待编译错误修复）
// auditLogListener := listener.NewAuditLogListener(eventBus)
// _ = auditLogListener // 保持引用，防止被 GC
// pkglogger.Info("✓ Audit Log Listener registered")
```

**优点**：
- ✅ 可立即让 API 编译通过
- ✅ 不影响核心功能（EventBus 仍可正常工作）
- ✅ EventBus 已集成到 AuthService

**缺点**：
- ⚠️ Listener 功能暂时不可用
- ⚠️ 审计日志不会自动记录

---

### 方案 B：彻底修复（待实施）

**需要检查的点**：

1. **检查 events.go 的 package 声明**
```bash
head -1 backend/internal/domain/user/events.go
# 应该输出：package events
```

2. **检查导入路径**
```go
// audit_log_listener.go 中的导入
import (
    "github.com/shenfay/go-ddd-scaffold/internal/domain/user/events"
)
```

3. **验证目录结构**
```bash
find backend/internal/domain -name "*.go" | head -10
```

4. **尝试重新生成 ULID 函数**

可能是 `generateUserID` 函数未定义导致的连锁错误。

---

## 当前状态

### ✅ 已解决
- API 服务可以编译
- EventBus 已正确初始化
- AuthService 已集成 EventBus

### ⚠️ 待解决
- Listener 层的编译错误
- 审计日志自动记录功能

---

## 下一步行动

### 立即可做（推荐）

1. **测试现有功能**
```bash
# 启动 API 服务
cd backend/cmd/api && go run main.go

# 在另一个终端测试登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"test@example.com","password":"Test123456"}'
```

2. **验证 EventBus 工作**
   - 检查日志中是否有 "✓ Event Bus initialized"
   - 登录成功后手动检查 audit_logs 表

---

### 短期修复（今天）

1. **排查 Listener 编译错误**
```bash
cd backend
go build ./internal/domain/user/...
go build ./internal/listener/...
```

2. **恢复 Listener 功能**
   - 修复编译错误后
   - 取消注释 cmd/api/main.go 中的 Listener 代码
   - 重新测试完整流程

---

## 经验总结

### 教训
1. 新增目录结构时应立即验证编译
2. Listener 层应该在最后一步集成
3. 应该先保证最小可用版本（MVP）

### 改进
1. 下次重构时采用渐进式策略：
   - 第一步：实现核心功能（已完成）
   - 第二步：验证编译通过（进行中）
   - 第三步：逐步添加增强功能（Listener、Worker）

---

**修复时间**: 2026-04-03  
**状态**: ⚠️ 部分完成（API 可编译，Listener 待修复）  
**负责人**: 后端团队
