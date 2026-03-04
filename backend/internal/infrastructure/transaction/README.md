# Transaction & UnitOfWork 使用指南

## 概述

本项目采用 **UnitOfWork + Repository** 模式来实现跨聚合根的事务一致性。

## 核心接口

### Transaction（事务）
```go
type Transaction interface {
    Commit() error           // 提交事务
    Rollback() error         // 回滚事务
    GetDB() *gorm.DB        // 获取底层 GORM 实例
}
```

### UnitOfWork（工作单元）
```go
type UnitOfWork interface {
    Begin(ctx context.Context) (Transaction, error)
}
```

## 使用示例

### 1. 在应用服务中使用 UnitOfWork

```go
func (s *TenantService) CreateTenant(
    ctx context.Context, 
    name string, 
    maxMembers int, 
    ownerID uuid.UUID,
) (*tenantEntity.Tenant, error) {
    // 1. 开启事务
    tx, err := s.uow.Begin(ctx)
    if err != nil {
        return nil, err
    }
    
    // 2. 使用 defer 确保异常时回滚
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        }
    }()
    
    // 3. 创建租户实体
    tenant := tenantEntity.NewTenant(name, maxMembers)
    
    // 4. 在事务中保存租户（需要仓储支持 WithTx 方法）
    if err := s.tenantRepo.CreateWithTx(tx.GetDB(), ctx, tenant); err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 5. 添加创建者为 owner（聚合根方法）
    member, err := tenant.AddMember(ownerID, sharedEntity.RoleOwner, nil)
    if err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 6. 在事务中保存成员关系
    if err := s.memberRepo.CreateWithTx(tx.GetDB(), ctx, member); err != nil {
        tx.Rollback()
        return nil, err
    }
    
    // 7. 提交事务
    if err := tx.Commit(); err != nil {
        return nil, err
    }
    
    return tenant, nil
}
```

### 2. 仓储层支持事务

需要在仓储接口中添加 `WithTx` 方法：

```go
// TenantRepository 租户仓储接口
type TenantRepository interface {
    // 基础方法
    Create(ctx context.Context, tenant *entity.Tenant) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
    Update(ctx context.Context, tenant *entity.Tenant) error
    Delete(ctx context.Context, id uuid.UUID) error
    
    // 事务支持方法
    CreateWithTx(tx *gorm.DB, ctx context.Context, tenant *entity.Tenant) error
    GetByIDWithTx(tx *gorm.DB, ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
    UpdateWithTx(tx *gorm.DB, ctx context.Context, tenant *entity.Tenant) error
}
```

### 3. 依赖注入配置（Wire）

```go
// InitializeUnitOfWork 初始化工作单元
func InitializeUnitOfWork(db *gorm.DB) transaction.UnitOfWork {
    return transaction.NewGormUnitOfWork(db)
}

// InitializeTenantService 初始化租户服务
func InitializeTenantService(
    uow transaction.UnitOfWork,
    tenantRepo repository.TenantRepository,
    memberRepo repository.TenantMemberRepository,
) *service.TenantService {
    return service.NewTenantService(uow, tenantRepo, memberRepo)
}
```

## 最佳实践

### ✅ DO - 推荐做法

1. **在应用服务层管理事务边界**
   - 领域层不应该知道事务的存在
   - 基础设施层只负责提供事务能力

2. **使用 defer+recover 保证回滚**
   ```go
   defer func() {
       if p := recover(); p != nil {
           tx.Rollback()
           panic(p)
       }
   }()
   ```

3. **保持事务尽可能短**
   - 避免在事务中进行网络调用、文件 IO 等耗时操作
   - 只做必要的数据库操作

4. **明确标注事务性方法**
   ```go
   // CreateTenant 创建租户（事务性方法）
   func (s *Service) CreateTenant(...) { ... }
   ```

### ❌ DON'T - 避免做法

1. **不要在领域层开启事务**
   ```go
   // 错误示例
   type Tenant struct {
       func AddMember(...) {
           tx := db.Begin() // ❌ 领域对象不应该直接访问数据库
       }
   }
   ```

2. **不要嵌套过深的事务**
   ```go
   // 尽量避免
   func MethodA() {
       tx.Begin()
       MethodB() // MethodB 内部又开启了新事务
   }
   ```

3. **不要忘记回滚**
   ```go
   // 错误示例
   tx, _ := uow.Begin()
   doSomething() // 如果这里 panic，事务永远不会回滚
   tx.Commit()
   ```

## 事务隔离级别

PostgreSQL 支持的隔离级别：

- **Read Committed**（默认）：适合大多数场景
- **Repeatable Read**：需要一致性读时使用
- **Serializable**：最强隔离，性能开销大

设置隔离级别：
```go
tx := db.Session(&gorm.Session{PrepareStmt: true}).Begin()
db.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ")
```

## 常见问题

### Q: 什么时候需要使用 UnitOfWork？

**A:** 当你的业务操作需要同时修改多个聚合根，并且这些修改要么全部成功、要么全部失败时。

例如：创建租户时，需要同时：
1. 保存租户信息
2. 添加创建者为 owner 角色

这两个操作必须原子性地完成。

### Q: UnitOfWork 和 @Transactional 注解有什么区别？

**A:** 
- Spring 的 `@Transactional` 是声明式事务管理（基于 AOP）
- UnitOfWork 是编程式事务管理（手动控制）

Go 语言更倾向于显式的编程式事务，因为：
1. 代码更清晰，事务边界一目了然
2. 不依赖反射和字节码增强
3. 更容易理解和调试

### Q: 如何处理分布式事务？

**A:** 本项目的 UnitOfWork 仅适用于单体应用的本地事务。如果需要分布式事务（微服务场景），建议使用：

1. **Saga 模式**：通过事件编排协调多个服务
2. **TCC 模式**：Try-Confirm-Cancel 三段式提交
3. **消息队列 + 最终一致性**：适合对实时性要求不高的场景

## 参考资料

- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [UnitOfWork Pattern](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [GORM Transactions](https://gorm.io/docs/transactions.html)
