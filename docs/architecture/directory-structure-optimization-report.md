# 工程目录结构优化报告 v2

> 更新时间：2026-03-20 11:15  
> 分析范围：`go-ddd-scaffold` 后端工程  
> 目标：基于 Clean Architecture + DDD 最佳实践，跟踪已完成的优化并给出后续建议

---

## 一、已完成优化分析

### ✅ 已完成项 1：`domain/user` 子包合并（🔴 高优先级 → 已完成）

**优化前：**
```
domain/user/
  aggregate/user.go
  valueobject/email.go
  valueobject/user_id.go
  ...
  event/user_registered.go
  ...
  repository/user_repository.go
  service/password_hasher.go
  types.go   ← re-export 补丁
```

**优化后：**
```
domain/user/
  user.go          ← 聚合根
  builder.go       ← 构建器
  vo.go            ← 所有值对象
  events.go        ← 领域事件
  repository.go    ← Repository 接口
  service.go       ← 领域服务接口
  event_handlers.go ← 事件处理器（见下文讨论）
```

**评价：** ⭐⭐⭐⭐⭐ 完美执行，删除了 `types.go` 补丁，import 路径统一为 `domain/user`。

---

### ✅ 已完成项 2：`audit/loginlog` 的 `event_handler` 迁移（🔴 高优先级 → 已完成）

**优化前：**
```
domain/audit/
  entity.go
  event_handler.go   ← 职责错位

domain/loginlog/
  entity.go
  event_handler.go
```

**优化后：**
```
domain/audit/
  entity.go          ← 只保留实体定义

domain/loginlog/
  entity.go

interfaces/event/
  subscriber.go
  audit_subscriber.go    ← 从 domain/audit 移过来 ✅
  loginlog_subscriber.go ← 从 domain/loginlog 移过来 ✅
```

**评价：** ⭐⭐⭐⭐⭐ 职责正确归位，`domain/audit` 和 `domain/loginlog` 现在只包含纯实体定义，副作用处理器移到了 `interfaces/event`。

---

### ✅ 已完成项 3：HTTP 接口层统一（🟡 中优先级 → 已完成）

**优化前：**
```
interfaces/http/auth/
  handler.go
  provider.go       ← 命名模糊
  request.go
  response.go
  （无 mapper.go）
```

**优化后：**
```
interfaces/http/auth/
  handler.go
  request.go
  response.go
  mapper.go         ← 新增 ✅
  routes.go         ← provider.go 改名 ✅

interfaces/http/user/
  handler.go
  request.go
  response.go
  mapper.go
  routes.go         ← provider.go 改名 ✅
```

**评价：** ⭐⭐⭐⭐⭐ 
- `provider.go` 改名为 `routes.go`，语义清晰
- `auth/mapper.go` 补齐，与 `user` 风格统一
- `routes.go` 采用 `Register(routerGroup *gin.RouterGroup)` 设计，支持灵活挂载

---

## 二、当前目录结构现状

```
internal/
├── domain/
│   ├── user/
│   │   ├── user.go           ✅ 聚合根
│   │   ├── builder.go        ✅ 构建器
│   │   ├── vo.go             ✅ 值对象
│   │   ├── events.go         ✅ 领域事件
│   │   ├── repository.go     ✅ 仓储接口
│   │   ├── service.go        ✅ 领域服务接口
│   │   └── event_handlers.go ⚠️ 见下文讨论
│   │
│   ├── tenant/               ✅ 扁平结构，与 user 一致
│   │   ├── entity.go
│   │   ├── vo.go
│   │   ├── events.go
│   │   ├── repository.go
│   │   └── service.go
│   │
│   ├── audit/
│   │   └── entity.go         ✅ 只保留实体
│   │
│   ├── loginlog/
│   │   └── entity.go         ✅ 只保留实体
│   │
│   └── shared/kernel/        ✅ 共享内核
│
├── application/              ✅ 应用层结构良好
│   ├── auth/
│   ├── user/
│   └── shared/dto/
│
├── infrastructure/
│   ├── persistence/          ✅ 仓储实现
│   ├── domain_event/         ⚠️ 待合并
│   ├── messaging/            ⚠️ 待合并
│   ├── auth/
│   ├── cache/
│   ├── config/
│   ├── logging/
│   ├── snowflake/
│   └── task_queue/
│
├── interfaces/
│   ├── http/                 ✅ HTTP 层统一完成
│   │   ├── auth/
│   │   │   ├── handler.go
│   │   │   ├── mapper.go
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── routes.go
│   │   ├── user/
│   │   │   ├── handler.go
│   │   │   ├── mapper.go
│   │   │   ├── request.go
│   │   │   ├── response.go
│   │   │   └── routes.go
│   │   ├── middleware/
│   │   ├── router.go
│   │   ├── response.go
│   │   └── types.go
│   │
│   └── event/                ✅ 事件订阅者
│       ├── subscriber.go
│       ├── audit_subscriber.go
│       └── loginlog_subscriber.go
│
├── bootstrap/
│   ├── bootstrap.go
│   ├── auth_domain.go
│   ├── user_domain.go
│   └── helpers/              ⚠️ 待清理
│
└── container/
    └── container.go
```

---

## 三、剩余待优化项

### 🟡 中优先级：`domain_event` + `messaging` 合并

**现状：**
```
infrastructure/
  domain_event/
    bus.go
    publisher.go
    repository.go
    store.go
  messaging/
    event_bus.go   ← 重复定义
```

**建议：** 合并为 `eventbus/` 包，消除职责重叠。

```
infrastructure/
  eventbus/
    bus.go        ← 事件总线实现
    publisher.go  ← 事件发布器
    store.go      ← 事件持久化
```

---

### 🟢 低优先级：`bootstrap/helpers` 清理

**现状：**
```
bootstrap/
  helpers/
    event_registrar.go
    http_builder.go
```

**建议：** 内联到 `bootstrap.go` 或 `bootstrap/http.go`，减少目录层级。

---

### 🟢 低优先级：`event_handlers.go` 归属讨论

**现状：**
```
domain/user/
  event_handlers.go   ← 这是什么？
```

**问题：** 文件名暗示这是事件处理器，但事件处理器应该是**响应事件**的副作用逻辑，放在 `domain` 层不合适。

**可能的情况：**
1. 如果是**聚合根内部的事件应用逻辑**（如 `ApplyEvent`），应该内联到 `user.go`
2. 如果是**响应其他领域事件的副作用**，应该移到 `interfaces/event/user_subscriber.go`

**建议：** 检查文件内容，确认其职责后决定是否迁移。

---

### 🟢 低优先级：`domain/audit` 和 `domain/loginlog` 是否必要

**现状：**
```
domain/audit/
  entity.go   ← 只有 AuditLog 实体定义

domain/loginlog/
  entity.go   ← 只有 LoginLog 实体定义
```

**问题：** `AuditLog` 和 `LoginLog` 是**写时序日志**，不是业务聚合，它们的实体定义可以放到：
- `infrastructure/persistence/model/`（如果只用于持久化）
- `application/audit/model/`（如果应用层需要）
- `interfaces/event/model/`（如果只有事件订阅者使用）

**建议：** 视后续需求决定是否迁移。当前保留在 `domain/` 也无大碍，只是略显冗余。

---

## 四、整体评价

| 维度 | v1 评价 | v2 评价 | 变化 |
|------|---------|---------|------|
| 领域层结构 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | +2 |
| 职责归位 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | +2 |
| HTTP 接口层 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | +2 |
| 基础设施层 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | - |
| 启动流程 | ⭐⭐⭐ | ⭐⭐⭐ | - |

**总结：** 核心的高优先级优化已全部完成，工程结构质量显著提升。剩余工作量为低优先级清理，可根据实际需求灵活安排。

---

## 五、后续建议优先级

| 优先级 | 任务 | 预估工时 | 备注 |
|--------|------|----------|------|
| 🟡 中 | `domain_event` + `messaging` 合并 | 1h | 消除重复 |
| 🟢 低 | `event_handlers.go` 归属确认 | 30min | 视内容决定 |
| 🟢 低 | `bootstrap/helpers` 清理 | 30min | 可选 |
| 🟢 低 | `audit/loginlog` 实体迁移 | 30min | 可选 |

---

*报告完*
