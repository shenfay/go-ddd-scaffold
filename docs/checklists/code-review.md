# 代码审查清单

## 📋 审查流程

### 提交前自查（作者）
- [ ] 代码符合规范
- [ ] 测试通过
- [ ] 文档已更新

### 审查要点（Reviewer）
- [ ] 功能正确性
- [ ] 代码质量
- [ ] 规范遵守

---

## ✅ DDD 规范检查

### Domain Layer

- [ ] **纯业务逻辑**：是否只有领域逻辑，无基础设施依赖？
- [ ] **实体设计**：Entity 是否有明确的业务方法（而非 getter/setter）？
- [ ] **值对象**：ValueObject 是否不可变且有验证？
- [ ] **聚合根**：Aggregate Root 边界是否清晰？
- [ ] **领域事件**：重要的状态变更是否发布事件？

```go
// ✅ 正确
func (u *User) UpdateEmail(email Email) error {
    u.Email = email
    u.addEvent(UserEmailChangedEvent{...})
    return nil
}

// ❌ 错误
type User struct {
    Email string  // 直接使用字符串
}
func (u *User) SetEmail(email string) {  // Setter 模式
    u.Email = email
}
```

---

### Application Layer

- [ ] **职责单一**：是否只负责编排，不含业务逻辑？
- [ ] **DTO 使用**：是否返回 DTO 而非 Entity？
- [ ] **事务管理**：是否正确管理事务边界？
- [ ] **错误处理**：是否使用 AppError？

```go
// ✅ 正确：应用服务只协调
func (s *UserService) ChangeEmail(ctx context.Context, userID uuid.UUID, email string) error {
    user, _ := s.userRepo.GetByID(ctx, userID)
    user.UpdateEmail(valueobject.NewEmail(email))  // 业务逻辑在 Domain
    s.userRepo.Update(ctx, user)
    s.eventBus.Publish(user.Events())
    return nil
}

// ❌ 错误：包含业务判断
if email == oldEmail {
    return nil  // 业务逻辑不应该在这里
}
```

---

### Infrastructure Layer

- [ ] **接口实现**：是否实现了 Repository 接口？
- [ ] **Model 转换**：是否在内部完成 Model ↔ Entity 转换？
- [ ] **不返回 Model**：是否直接返回 Model 给应用层？

```go
// ✅ 正确
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    var model model.User
    r.db.First(&model, id)
    return r.toEntity(&model), nil  // 转换为 Entity
}

// ❌ 错误
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    var model model.User
    r.db.First(&model, id)
    return &model, nil  // 直接返回 Model
}
```

---

### Interfaces Layer

- [ ] **协议转换**：是否只做参数绑定和响应格式化？
- [ ] **无业务逻辑**：Handler 中是否有业务判断？
- [ ] **参数验证**：输入是否验证？
- [ ] **错误处理**：是否统一错误响应格式？

```go
// ✅ 正确：Handler 只转换
func (h *UserHandler) UpdateEmail(c *gin.Context) {
    var req UpdateEmailRequest
    c.ShouldBindJSON(&req)  // 参数绑定
    
    err := h.userService.ChangeEmail(...)  // 调用应用服务
    
    c.JSON(200, response.Success(user))  // 格式化响应
}

// ❌ 错误：包含业务逻辑
if req.Email == "" {
    // 应该在 validator 中验证
}
user, _ := userRepo.GetByID(...)  // 不应该直接访问仓储
```

---

## 🔒 安全检查

### 输入验证

- [ ] 所有用户输入是否验证？
- [ ] SQL 注入防护（参数化查询）？
- [ ] XSS 防护（输出转义）？
- [ ] CSRF 防护（Token 验证）？

```go
// ✅ 正确：使用 validator
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// ❌ 错误：未经验证直接使用
var email string
c.BindJSON(&email)
userRepo.Create(ctx, email)
```

---

### 敏感信息

- [ ] 密码是否加密存储？
- [ ] Token 是否安全存储？
- [ ] 日志中是否脱敏？
- [ ] 配置文件是否不包含敏感信息？

```go
// ✅ 正确：密码加密
hashedPassword, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

// ❌ 错误：明文密码
type User struct {
    Password string  // 明文存储
}

// ❌ 错误：日志记录敏感信息
log.Printf("用户登录，密码：%s", password)
```

---

## 🧪 测试检查

### 单元测试

- [ ] 核心业务逻辑是否有单元测试？
- [ ] 测试覆盖率是否达标（≥80%）？
- [ ] 测试是否独立（不依赖外部资源）？
- [ ] 测试命名是否清晰？

```go
// ✅ 正确：表驱动测试
func TestUser_UpdateEmail(t *testing.T) {
    tests := []struct {
        name      string
        newEmail  string
        wantError bool
    }{
        {"有效邮箱", "new@example.com", false},
        {"无效邮箱", "invalid", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

---

### 集成测试

- [ ] 关键业务流程是否有集成测试？
- [ ] 是否使用测试数据库？
- [ ] 测试数据是否清理？

---

## 📊 性能检查

### 数据库

- [ ] 是否有 N+1 查询问题？
- [ ] 索引是否合理？
- [ ] 是否使用预加载（Preload）？
- [ ] 连接池配置是否优化？

```go
// ✅ 正确：避免 N+1
db.Preload("Orders").Find(&users)

// ❌ 错误：N+1 查询
for _, user := range users {
    db.Where("user_id = ?", user.ID).Find(&orders)
}
```

---

### 缓存

- [ ] 热点数据是否使用缓存？
- [ ] 缓存策略是否合理（TTL、淘汰）？
- [ ] 缓存穿透/击穿/雪崩防护？

---

## 📝 代码质量

### 命名规范

- [ ] 包名小写、单数形式？
- [ ] 导出元素 PascalCase？
- [ ] 变量 camelCase、语义清晰？
- [ ] 常量全大写下划线？

```go
// ✅ 正确
package user
type UserRepository interface {}
var userID uuid.UUID
const StatusActive = "active"

// ❌ 错误
package Users      // 应该用 user
type IUser interface {}  // 不应该用 I 前缀
var uid uuid.UUID  // 应该用 userID
```

---

### 注释规范

- [ ] 所有导出元素是否有注释？
- [ ] 注释是否完整句子？
- [ ] 是否解释"为什么"而非"是什么"？

```go
// ✅ 正确
// UpdateEmail 更新用户邮箱
//
// 验证新邮箱有效性，并在变更后发布 UserEmailChangedEvent 事件。
// 如果新旧邮箱相同，返回 nil。
func (u *User) UpdateEmail(email Email) error {
    // 使用 bcrypt 加密，成本因子为 10
    // 这是在安全性和性能之间的平衡选择
}

// ❌ 错误
// 更新邮箱
func (u *User) UpdateEmail(email Email) error {
    // 加密密码
}
```

---

### 错误处理

- [ ] 是否统一使用 AppError？
- [ ] 错误码是否规范？
- [ ] 错误消息是否明确？
- [ ] 是否包装底层错误？

```go
// ✅ 正确
if err != nil {
    return errPkg.Wrap(err, "GET_USER_FAILED", "从数据库获取用户失败")
}

// ❌ 错误
if err != nil {
    return fmt.Errorf("failed: %v", err)  // 裸 error
}
```

---

## 🔄 架构检查

### 分层依赖

- [ ] 是否跨层调用？（禁止）
- [ ] 依赖方向是否正确？
- [ ] 循环依赖？（禁止）

```
✅ 正确的依赖方向:
Interfaces → Application → Domain ← Infrastructure

❌ 错误的依赖:
Domain → Infrastructure  // 禁止！
Interfaces → Domain      // 禁止跨层！
```

---

### 模块化

- [ ] 是否按领域划分目录？
- [ ] 模块间耦合是否过高？
- [ ] 是否有重复代码？

---

## 📋 审查结论

### 通过标准

- [ ] 所有必选项满足
- [ ] 无严重问题
- [ ] 一般问题 ≤ 5 个

### 不通过

出现以下情况直接拒绝：

- ❌ 安全问题
- ❌ 严重性能问题
- ❌ 违反核心规范
- ❌ 测试覆盖率 < 60%

---

## 🎯 常见问题 FAQ

### Q: 紧急修复是否可以跳过审查？
A: 不可以。紧急修复也需要审查，可以走快速通道但必须审查。

### Q: 审查意见有分歧怎么办？
A: 遵循规范文档，如规范未定义，由 Tech Lead 决定。

### Q: 如何跟踪审查意见？
A: 使用 GitHub/GitLab 的 Review 功能，所有意见必须有回应。

---

**版本**: v1.0  
**生效日期**: 2026-03-06  
**下次回顾**: 2026-03-20
