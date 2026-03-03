# Tenant Builder 快速开始指南

## 10 分钟搭建多租户 SaaS 架构

### 第一步：安装 Skill

```bash
npx skills install tenant-builder
```

### 第二步：创建租户表

```bash
tenant-builder create \
  --module tenant \
  --description "create tenants table"
```

这会生成 Goose 迁移文件：
```
backend/migrations/sql/20260225170000_create_tenants_table.sql
```

### 第三步：创建家庭系统表

```bash
tenant-builder create \
  --module family \
  --description "create family system"
```

### 第四步：创建订阅管理表

```bash
tenant-builder create \
  --module subscription \
  --description "create subscription management"
```

### 第五步：执行数据库迁移

```bash
# 确保数据库已启动
docker-compose up -d postgres

# 执行所有迁移
tenant-builder migrate up --env dev
```

查看迁移状态：
```bash
tenant-builder migrate status --env dev
```

### 第六步：生成 GORM Model 和 DAO

```bash
tenant-builder generate \
  --tables tenants,families,family_members,subscriptions \
  --output ./internal/infrastructure/persistence/gorm
```

生成的文件：
```
internal/infrastructure/persistence/gorm/
├── model/
│   ├── tenant.go
│   ├── family.go
│   ├── family_member.go
│   └── subscription.go
└── dao/
    ├── tenant_repository.go
    ├── tenant_repository_impl.go
    ├── family_repository.go
    └── ...
```

### 第七步：初始化 Casbin 策略

```bash
tenant-builder init-casbin --default-role admin
```

这会生成 `backend/config/auth/rbac.conf` 文件，包含完整的 RBAC 策略。

### 第八步：生成租户中间件（可选）

```bash
# 手动创建中间件文件
# 参考 EXAMPLES.md 中的中间件实现代码
```

### 第九步：在应用中集成

#### 1. 注册中间件

```go
// internal/infrastructure/server/http_server.go
r.Use(middleware.TenantMiddleware())
```

#### 2. 注入 Repository

```go
// cmd/server/main.go
tenantRepo := dao.NewTenantRepository(db)
familyRepo := dao.NewFamilyRepository(db)

tenantService := service.NewTenantService(tenantRepo, familyRepo)
tenantHandler := http.NewTenantHandler(tenantService)
```

#### 3. 使用示例

```go
// 创建租户
func CreateTenant(c *gin.Context) {
    req := &CreateTenantRequest{}
    if err := c.ShouldBindJSON(req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    entity := &model.Tenant{
        Name:      req.Name,
        Subdomain: req.Subdomain,
        Plan:      "free",
    }
    
    err := tenantService.CreateTenant(c.Request.Context(), entity)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, entity)
}
```

## 完整工作流（一键完成）

```bash
tenant-builder full-workflow \
  --modules tenant,family,subscription \
  --env dev \
  --generate-dao \
  --init-casbin \
  --with-middleware
```

这会引导你完成：
1. ✅ 创建所有迁移文件
2. ✅ 执行数据库迁移
3. ✅ 生成 DAO 代码
4. ✅ 初始化 Casbin 策略
5. ✅ 生成中间件代码

## 生成的数据库表结构

### tenants 表

```sql
CREATE TABLE tenants (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    subdomain       VARCHAR(50) UNIQUE NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    plan            VARCHAR(20) DEFAULT 'free',
    expired_at      TIMESTAMP WITH TIME ZONE,
    max_users       INT DEFAULT 10,
    max_storage     BIGINT DEFAULT 1073741824,  -- 1GB
    current_users   INT DEFAULT 1,
    current_storage BIGINT DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);
```

### families 表

```sql
CREATE TABLE families (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id),
    name            VARCHAR(100) NOT NULL,
    creator_id      BIGINT NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);
```

### family_members 表

```sql
CREATE TABLE family_members (
    id              BIGSERIAL PRIMARY KEY,
    family_id       BIGINT NOT NULL REFERENCES families(id),
    user_id         BIGINT NOT NULL,
    role            VARCHAR(20) NOT NULL DEFAULT 'child',
    joined_at       TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    invited_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL,
    UNIQUE(family_id, user_id)
);
```

### subscriptions 表

```sql
CREATE TABLE subscriptions (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT UNIQUE NOT NULL REFERENCES tenants(id),
    plan_id         BIGINT REFERENCES subscription_plans(id),
    status          VARCHAR(20) DEFAULT 'active',
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    auto_renew      BOOLEAN DEFAULT false,
    last_payment_at TIMESTAMP WITH TIME ZONE,
    next_billing_at TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);
```

## Casbin 策略配置

生成的 `rbac.conf` 包含以下策略：

```conf
# 租户管理员
p, admin, tenants, read, allow
p, admin, tenants, update, allow
p, admin, families, manage, allow

# 家长角色
p, parent, children, read, allow
p, parent, children, manage, allow
p, parent, learning_progress, read, allow

# 孩子角色
p, child, learning_resources, read, allow
p, child, games, play, allow

# 教育工作者角色
p, educator, students, read, allow
p, educator, students, assign_tasks, allow
```

## 数据隔离机制

### 中间件自动注入 tenant_id

```go
// middleware/tenant.go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从子域名提取租户
        subdomain := extractSubdomain(c.Request.Host)
        tenant := getTenantBySubdomain(subdomain)
        
        if tenant == nil {
            c.JSON(404, gin.H{"error": "tenant not found"})
            c.Abort()
            return
        }
        
        // 注入上下文
        ctx := context.WithValue(c.Request.Context(), "tenant_id", tenant.ID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}
```

### Repository 自动添加租户过滤

```go
// dao/family_repository_impl.go
func (r *FamilyRepositoryImpl) FindAll(ctx context.Context, offset, limit int) ([]model.Family, int64, error) {
    // 从上下文获取 tenant_id
    tenantID := ctx.Value("tenant_id").(int64)
    
    var families []model.Family
    var total int64
    
    // 自动添加租户过滤条件
    tx := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
    
    tx.Model(&model.Family{}).Where("tenant_id = ?", tenantID).Count(&total)
    tx.Offset(offset).Limit(limit).Find(&families)
    
    return families, total, tx.Error
}
```

## 套餐计划配置

### 免费版（Free）
- 最多 3 个用户
- 100MB 存储空间
- 基础游戏
- 有限进度追踪

### 基础版（Basic）
- 最多 10 个用户
- 1GB 存储空间
- 完整进度追踪
- 家庭报告
- ¥29.9/月 或 ¥299/年

### 高级版（Premium）
- 最多 50 个用户
- 5GB 存储空间
- 全部游戏
- 高级分析
- 优先支持
- API 访问
- ¥59.9/月 或 ¥599/年

## 下一步

### 1. 添加邀请功能

参考 [EXAMPLES.md](./EXAMPLES.md) 中的邀请系统设计。

### 2. 实现订阅支付

集成微信支付或支付宝的订阅支付接口。

### 3. 完善权限控制

根据实际业务需求调整 Casbin 策略。

### 4. 添加用量统计

实现租户用量监控和超限处理逻辑。

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：咨询 DDD Architect Agent

祝你开发顺利！🚀
