# DB Migrator 使用示例

## 示例 1: 用户管理系统数据库

### 场景
快速搭建用户管理系统的完整数据库结构，包含用户表、角色表、用户角色关联表。

### 执行命令

```bash
# 1. 创建用户表
db-migrator create \
  --table users \
  --description "create users table"

# 2. 创建角色表
db-migrator create \
  --table roles \
  --description "create roles table"

# 3. 创建用户角色关联表
db-migrator create \
  --table user_roles \
  --description "create user role mapping table"

# 4. 执行迁移
db-migrator migrate up --env dev

# 5. 生成 DAO 代码
db-migrator generate \
  --tables users,roles,user_roles \
  --output ./internal/infrastructure/persistence/gorm
```

### 生成的表结构

#### users 表
```sql
CREATE TABLE users (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(50) NOT NULL,
    email           VARCHAR(100) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    avatar_url      VARCHAR(500),
    status          VARCHAR(20) DEFAULT 'active',
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON COLUMN users.status IS '用户状态：active, inactive, banned';
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

#### roles 表
```sql
CREATE TABLE roles (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(50) UNIQUE NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.code IS '角色代码：admin, member, guest';
```

#### user_roles 表（多对多关联）
```sql
CREATE TABLE user_roles (
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id         BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    granted_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    granted_by      BIGINT REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);

COMMENT ON TABLE user_roles IS '用户角色关联表';
CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
```

### GORM Model 关联

```go
// model/user.go
type User struct {
    ID        int64
    Name      string
    Email     string
    Roles     []Role `gorm:"Many2Many:user_roles;" json:"roles"`
}

// model/role.go
type Role struct {
    ID    int64
    Code  string
    Name  string
    Users []User `gorm:"Many2Many:user_roles;" json:"users"`
}
```

---

## 示例 2: 电商订单系统

### 场景
构建完整的电商订单模块，包含商品、订单、订单明细、物流信息。

### 执行命令

```bash
# 批量创建迁移
db-migrator create --table products --description "create products table"
db-migrator create --table orders --description "create orders table"
db-migrator create --table order_items --description "create order items table"
db-migrator create --table shipments --description "create shipments table"

# 执行所有迁移
db-migrator migrate up --env dev

# 生成代码（包含测试）
db-migrator generate \
  --tables products,orders,order_items,shipments \
  --output ./internal/infrastructure/persistence/gorm \
  --with-tests
```

### 订单表设计

```sql
-- 订单表
CREATE TABLE orders (
    id              BIGSERIAL PRIMARY KEY,
    order_no        VARCHAR(50) UNIQUE NOT NULL,
    user_id         BIGINT NOT NULL,
    total_amount    DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    pay_amount      DECIMAL(10,2) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method  VARCHAR(20),
    payment_time    TIMESTAMP WITH TIME ZONE,
    shipping_address JSONB NOT NULL,
    remark          VARCHAR(500),
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE orders IS '订单表';
COMMENT ON COLUMN orders.order_no IS '订单编号';
COMMENT ON COLUMN orders.status IS '订单状态：pending, paid, shipped, completed, cancelled';
COMMENT ON COLUMN orders.shipping_address IS '收货地址（JSON 格式）';
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created ON orders(created_at);

-- 订单明细表
CREATE TABLE order_items (
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      BIGINT NOT NULL,
    product_name    VARCHAR(200) NOT NULL,
    product_image   VARCHAR(500),
    price           DECIMAL(10,2) NOT NULL,
    quantity        INT NOT NULL DEFAULT 1,
    subtotal        DECIMAL(10,2) NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE order_items IS '订单明细表';
CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_product ON order_items(product_id);
```

### 复杂查询方法

在 Repository 实现中添加自定义查询：

```go
// dao/order_repository_impl.go

// FindByUserID 根据用户 ID 查找订单
func (r *OrderRepositoryImpl) FindByUserID(ctx context.Context, userID int64, offset, limit int) ([]model.Order, int64, error) {
    var orders []model.Order
    var total int64
    
    // 统计总数
    if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 查询订单
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Preload("Items").
        Order("created_at DESC").
        Offset(offset).
        Limit(limit).
        Find(&orders).Error
    
    return orders, total, err
}

// GetOrderWithDetails 获取订单及详情
func (r *OrderRepositoryImpl) GetOrderWithDetails(ctx context.Context, orderID int64) (*model.Order, error) {
    var order model.Order
    err := r.db.WithContext(ctx).
        Preload("Items").
        First(&order, orderID).Error
    return &order, err
}
```

---

## 示例 3: 博客内容管理系统

### 场景
为博客系统创建文章、分类、标签、评论等表结构。

### 特殊需求
- 文章支持 Markdown 和富文本
- 标签多对多关联
- 评论嵌套结构（父子评论）
- 浏览量、点赞数统计

### 执行命令

```bash
db-migrator create --table articles --description "create articles table"
db-migrator create --table categories --description "create categories table"
db-migrator create --table tags --description "create tags table"
db-migrator create --table article_tags --description "create article tag mapping table"
db-migrator create --table comments --description "create comments table"

db-migrator migrate up --env dev
db-migrator generate --tables articles,categories,tags,comments --with-comments
```

### 文章表设计

```sql
CREATE TABLE articles (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) UNIQUE NOT NULL,
    summary         VARCHAR(500),
    content         TEXT NOT NULL,
    cover_image     VARCHAR(500),
    author_id       BIGINT NOT NULL,
    category_id     BIGINT,
    status          VARCHAR(20) DEFAULT 'draft',
    view_count      BIGINT DEFAULT 0,
    like_count      BIGINT DEFAULT 0,
    comment_count   BIGINT DEFAULT 0,
    published_at    TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE articles IS '文章表';
COMMENT ON COLUMN articles.slug IS '文章别名（URL 友好）';
COMMENT ON COLUMN articles.status IS '状态：draft, published, archived';
COMMENT ON COLUMN articles.view_count IS '浏览量';
COMMENT ON COLUMN articles.like_count IS '点赞数';
CREATE INDEX idx_articles_author ON articles(author_id);
CREATE INDEX idx_articles_category ON articles(category_id);
CREATE INDEX idx_articles_status ON articles(status);
CREATE INDEX idx_articles_published ON articles(published_at);
```

### 评论嵌套结构

```sql
CREATE TABLE comments (
    id              BIGSERIAL PRIMARY KEY,
    article_id      BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id         BIGINT,
    parent_id       BIGINT REFERENCES comments(id) ON DELETE CASCADE,
    content         TEXT NOT NULL,
    ip_address      VARCHAR(50),
    user_agent      VARCHAR(200),
    status          VARCHAR(20) DEFAULT 'pending',
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE comments IS '评论表';
COMMENT ON COLUMN comments.parent_id IS '父评论 ID（用于嵌套回复）';
CREATE INDEX idx_comments_article ON comments(article_id);
CREATE INDEX idx_comments_parent ON comments(parent_id);
```

### GORM Model 自关联

```go
// model/comment.go
type Comment struct {
    ID        int64
    ArticleID int64
    ParentID  *int64
    Content   string
    
    // 自关联：父评论
    Parent    *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
    
    // 自关联：子评论列表
    Replies   []Comment `gorm:"foreignKey:ParentID" json:"replies"`
    
    Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
}
```

---

## 示例 4: 多租户 SaaS 系统

### 场景
构建支持多租户的 SaaS 平台，每个租户有独立的数据空间。

### 执行命令

```bash
# 创建租户表
db-migrator create --table tenants --description "create tenants table"

# 创建租户配置表
db-migrator create --table tenant_settings --description "create tenant settings table"

# 执行迁移
db-migrator migrate up --env dev

# 生成代码
db-migrator generate --tables tenants,tenant_settings
```

### 租户隔离设计

```sql
-- 租户表
CREATE TABLE tenants (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    subdomain       VARCHAR(50) UNIQUE NOT NULL,
    status          VARCHAR(20) DEFAULT 'active',
    plan            VARCHAR(20) DEFAULT 'free',
    expired_at      TIMESTAMP WITH TIME ZONE,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP WITH TIME ZONE NULL
);

COMMENT ON TABLE tenants IS '租户表';
COMMENT ON COLUMN tenants.subdomain IS '子域名（用于租户识别）';
COMMENT ON COLUMN tenants.plan IS '套餐：free, basic, premium';
CREATE INDEX idx_tenants_subdomain ON tenants(subdomain);

-- 用户表（添加租户关联）
ALTER TABLE users ADD COLUMN tenant_id BIGINT REFERENCES tenants(id);
CREATE INDEX idx_users_tenant ON users(tenant_id);

COMMENT ON COLUMN users.tenant_id IS '所属租户 ID';
```

### 租户中间件集成

```go
// infrastructure/middleware/tenant_middleware.go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从子域名或 Header 中提取租户 ID
        tenantID := extractTenantID(c)
        
        if tenantID == "" {
            c.JSON(400, gin.H{"error": "tenant not found"})
            c.Abort()
            return
        }
        
        // 设置租户上下文
        ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
        c.Request = c.Request.WithContext(ctx)
        
        c.Next()
    }
}

// dao/user_repository_impl.go - 租户数据隔离
func (r *UserRepositoryImpl) FindAll(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
    tenantID := ctx.Value("tenant_id").(int64)
    
    var users []model.User
    var total int64
    
    // 只查询当前租户的用户
    tx := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID)
    
    tx.Model(&model.User{}).Where("tenant_id = ?", tenantID).Count(&total)
    tx.Offset(offset).Limit(limit).Find(&users)
    
    return users, total, tx.Error
}
```

---

这些示例展示了 db-migrator 在各种场景下的应用。你可以根据实际业务需求调整和扩展。
