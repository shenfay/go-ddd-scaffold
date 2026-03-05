# 聚合根边界重新设计分析

## 当前问题分析

### 1. 聚合根职责不清
**现状问题**：
- `User`实体承担了过多聚合职责，但实际上应该是一个独立的聚合根
- `Tenant`聚合根没有很好地管理其成员关系
- 跨聚合的操作缺少事务一致性保证

### 2. 领域事件缺失
**现状问题**：
- 缺少聚合内操作的领域事件发布机制
- 没有建立事件溯源的能力
- 异步处理链路不完整

### 3. 边界控制不当
**现状问题**：
- 租户成员关系应该由`Tenant`聚合根统一管理
- 用户信息更新等操作应该在各自聚合内完成
- 缺少聚合间的引用规范

## 改进设计方案

### 方案一：双聚合根模式（推荐）

```
User Aggregate (用户聚合根)
├── User (聚合根实体)
├── UserProfile (值对象)
├── UserSettings (值对象)
└── 领域事件：UserRegistered, UserProfileUpdated, UserDeactivated

Tenant Aggregate (租户聚合根)  
├── Tenant (聚合根实体)
├── TenantMember (子实体)
├── TenantSettings (值对象)
└── 领域事件：TenantCreated, MemberAdded, MemberRemoved, TenantExpired
```

### 方案二：聚合根职责划分

**User聚合根职责**：
- 管理用户基本信息
- 处理用户认证相关逻辑
- 维护用户状态一致性

**Tenant聚合根职责**：
- 管理租户基本信息和配置
- 统一管理所有成员关系
- 控制租户生命周期
- 执行成员增删改操作

### 方案三：聚合间交互规范

**允许的操作**：
- 通过领域服务协调跨聚合操作
- 使用最终一致性保证数据同步
- 通过事件驱动实现松耦合

**禁止的操作**：
- 直接从一个聚合修改另一个聚合的状态
- 跨聚合的强一致性事务
- 聚合间的直接引用

## 实施步骤

### 第一步：明确聚合边界
1. 定义User聚合的完整边界
2. 定义Tenant聚合的完整边界
3. 建立聚合间交互规范

### 第二步：引入领域事件
1. 为每个聚合根定义核心领域事件
2. 建立事件发布/订阅机制
3. 实现事件存储和重放能力

### 第三步：优化仓储设计
1. 仓储只为单个聚合根服务
2. 跨聚合查询通过领域服务协调
3. 引入读模型优化查询性能

## 设计原则

### 1. 单一职责原则
每个聚合根只负责一个业务领域的完整性

### 2. 高内聚低耦合
聚合内实体高度内聚，聚合间松耦合

### 3. 最终一致性
接受最终一致性，避免分布式事务

### 4. 显式建模
通过显式的聚合根和领域事件表达业务意图

## 领域事件清单

### User聚合事件
- `UserRegistered`: 用户注册完成
- `UserProfileUpdated`: 用户资料更新
- `UserPasswordChanged`: 用户密码修改
- `UserDeactivated`: 用户停用
- `UserDeleted`: 用户删除

### Tenant聚合事件
- `TenantCreated`: 租户创建
- `MemberAdded`: 成员添加
- `MemberRemoved`: 成员移除
- `MemberRoleChanged`: 成员角色变更
- `TenantExpired`: 租户过期
- `TenantSettingsUpdated`: 租户设置更新

## 技术实现要点

### 1. 聚合根基类设计
```go
type AggregateRoot interface {
    ID() uuid.UUID
    Version() int
    IncrementVersion()
    ApplyEvent(event DomainEvent)
    GetUncommittedEvents() []DomainEvent
    ClearUncommittedEvents()
}
```

### 2. 仓储接口约束
```go
type Repository interface {
    Save(aggregate AggregateRoot) error
    GetByID(id uuid.UUID) (AggregateRoot, error)
    // 不允许跨聚合查询
}
```

### 3. 领域服务协调
```go
type TenantUserService struct {
    tenantRepo TenantRepository
    userRepo   UserRepository
    eventBus   EventBus
}
```

## 风险与挑战

### 1. 学习成本
团队需要理解DDD聚合根概念

### 2. 性能考虑
最终一致性可能影响用户体验

### 3. 调试复杂度
分布式事件增加了调试难度

### 4. 数据一致性
需要仔细设计补偿机制