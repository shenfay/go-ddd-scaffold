# API Generator 使用示例

## 示例 1: 快速生成用户管理 API

### 场景
已有 UserAggregate，需要快速生成完整的 CRUD API 端点。

### 执行命令

```bash
/api-generator \
  --aggregate UserAggregate \
  --with-validation \
  --auth jwt \
  --output ./user-api
```

### 生成的文件

```
user-api/
├── internal/
│   ├── interfaces/http/user/
│   │   ├── user_handler.go    # HTTP Handler (CRUD)
│   │   └── user_router.go     # 路由配置
│   └── application/user/dto/
│       └── user_dto.go        # DTO 定义
└── docs/
    └── swagger.json           # API 文档
```

### 生成的 Handler 代码（部分）

```go
// @Summary Create User
// @Tags UserManagement
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dto.CreateUserRequest true "User information"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}
	
	userID := c.GetString("user_id")
	
	result, err := h.service.CreateUser(c.Request.Context(), req, userID)
	if err != nil {
		c.JSON(response.HTTPStatus(err), response.Error(err))
		return
	}
	
	c.JSON(http.StatusCreated, response.Success(result))
}
```

---

## 示例 2: 电商系统完整 API

### 场景
构建完整的电商平台 API，包含商品、订单、库存三个核心模块。

### 执行命令

```bash
/api-generator \
  --aggregates Product,Order,Inventory \
  --with-validation \
  --auth jwt \
  --with-tests \
  --output ./ecommerce-api
```

### 生成的 API 端点

#### 商品管理
- `POST /api/v1/products` - 创建商品
- `GET /api/v1/products/:id` - 获取商品详情
- `PUT /api/v1/products/:id` - 更新商品
- `DELETE /api/v1/products/:id` - 删除商品
- `GET /api/v1/products` - 商品列表（支持筛选、分页、排序）

#### 订单管理
- `POST /api/v1/orders` - 创建订单
- `GET /api/v1/orders/:id` - 订单详情
- `PUT /api/v1/orders/:id/status` - 更新订单状态
- `GET /api/v1/orders` - 订单列表
- `DELETE /api/v1/orders/:id` - 取消订单

#### 库存管理
- `GET /api/v1/inventory/:productId` - 查询库存
- `PUT /api/v1/inventory/:productId` - 更新库存
- `POST /api/v1/inventory/reserve` - 预占库存
- `POST /api/v1/inventory/release` - 释放库存

### 批量操作示例

```go
// POST /api/v1/products/batch
type BatchCreateProductsRequest struct {
	Items []CreateProductRequest `json:"items" binding:"required,dive"`
}

func (h *ProductHandler) BatchCreateProducts(c *gin.Context) {
	var req BatchCreateProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}
	
	// 验证数量限制
	if len(req.Items) > 100 {
		c.JSON(http.StatusBadRequest, response.Error(errors.New("max 100 items per batch")))
		return
	}
	
	results, err := h.service.BatchCreateProducts(c.Request.Context(), req.Items)
	if err != nil {
		c.JSON(response.HTTPStatus(err), response.Error(err))
		return
	}
	
	c.JSON(http.StatusCreated, response.Success(results))
}
```

---

## 示例 3: 博客系统 API

### 场景
为博客系统生成 API，包含文章、评论、标签管理。

### 特殊需求
- 文章需要支持 Markdown 格式
- 评论需要嵌套结构（父子评论）
- 标签支持多对多关系

### 执行命令

```bash
/api-generator \
  --aggregates Post,Comment,Tag \
  --with-validation \
  --auth jwt \
  --prefix /api/v1/blog
```

### 自定义 DTO

```go
// Post DTO
type CreatePostRequest struct {
	Title     string   `json:"title" binding:"required,min=5,max=200"`
	Content   string   `json:"content" binding:"required,min=100"`
	Summary   string   `json:"summary" binding:"omitempty,max=500"`
	TagIDs    []string `json:"tag_ids" binding:"required"`
	CategoryID string  `json:"category_id" binding:"required,uuid"`
	Status    string   `json:"status" binding:"required,oneof=draft published archived"`
}

type PostResponse struct {
	ID         string         `json:"id"`
	Title      string         `json:"title"`
	Content    string         `json:"content"`
	Summary    string         `json:"summary"`
	Tags       []TagResponse  `json:"tags"`
	Category   CategoryResponse `json:"category"`
	AuthorID   string         `json:"author_id"`
	ViewCount  int64          `json:"view_count"`
	Status     string         `json:"status"`
	PublishedAt *time.Time    `json:"published_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
```

### 嵌套评论支持

```go
// Comment DTO with nested structure
type CommentResponse struct {
	ID        string           `json:"id"`
	Content   string           `json:"content"`
	AuthorID  string           `json:"author_id"`
	PostID    string           `json:"post_id"`
	ParentID  string           `json:"parent_id,omitempty"`
	Replies   []CommentResponse `json:"replies,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
}
```

---

## 示例 4: 带权限控制的 API

### 场景
需要基于角色的访问控制（RBAC），不同角色有不同的 API 访问权限。

### 执行命令

```bash
/api-generator \
  --aggregate Document \
  --with-validation \
  --auth casbin
```

### Casbin 权限配置

```ini
# config/auth/rbac_with_domains.conf
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && r.obj == p.obj && r.act == p.act
```

### 权限策略示例

```csv
# config/auth/policy.csv
p, admin, tenant1, documents, create
p, admin, tenant1, documents, read
p, admin, tenant1, documents, update
p, admin, tenant1, documents, delete

p, member, tenant1, documents, read
p, member, tenant1, documents, create

p, guest, tenant1, documents, read
```

### Handler 中的权限检查

```go
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	
	// 从 JWT 中提取用户信息和租户
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")
	role := c.GetString("role")
	
	// Casbin 权限检查
	enforcer := h.enforcer
	allowed, err := enforcer.Enforce(role, tenantID, "documents", "delete")
	if err != nil || !allowed {
		c.JSON(http.StatusForbidden, response.Error(errors.New("permission denied")))
		return
	}
	
	err = h.service.DeleteDocument(c.Request.Context(), id)
	if err != nil {
		c.JSON(response.HTTPStatus(err), response.Error(err))
		return
	}
	
	c.Status(http.StatusNoContent)
}
```

---

## 示例 5: 文件上传 API

### 场景
需要生成支持文件上传的 API 端点。

### 手动扩展 Handler

在生成的 Handler 基础上添加上传方法：

```go
// @Summary Upload File
// @Tags FileManagement
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Success 200 {object} response.Response{data=FileUploadResponse}
// @Router /api/v1/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}
	
	// 验证文件类型
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"application/pdf": true,
	}
	
	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, response.Error(errors.New("unsupported file type")))
		return
	}
	
	// 验证文件大小（最大 10MB）
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, response.Error(errors.New("file too large")))
		return
	}
	
	// 生成唯一文件名
	filename := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(file.Filename))
	dst := filepath.Join("./uploads", filename)
	
	// 保存文件
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(err))
		return
	}
	
	// 返回文件信息
	response := &FileUploadResponse{
		FileID:   uuid.New().String(),
		Filename: filename,
		Size:     file.Size,
		URL:      fmt.Sprintf("/uploads/%s", filename),
	}
	
	c.JSON(http.StatusOK, response.Success(response))
}
```

---

## 示例 6: WebSocket + RESTful 混合 API

### 场景
需要实时通知功能，当数据变更时通过 WebSocket 推送。

### 前置条件

已安装 `websocket-integration` Skill

### 执行命令

```bash
# 生成 RESTful API
/api-generator --aggregate Notification --with-validation --auth jwt

# 集成 WebSocket
/websocket-integration --target-dir ./my-app
```

### 结合使用

```go
// RESTful API - 创建通知
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	var req dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err))
		return
	}
	
	notification, err := h.service.CreateNotification(c.Request.Context(), req)
	if err != nil {
		c.JSON(response.HTTPStatus(err), response.Error(err))
		return
	}
	
	// 通过 WebSocket 实时推送
	wsManager := GetWebSocketManager()
	wsManager.NotifyUser(notification.UserID, map[string]interface{}{
		"type": "new_notification",
		"data": notification,
	})
	
	c.JSON(http.StatusCreated, response.Success(notification))
}
```

---

## 示例 7: 版本化 API

### 场景
需要同时支持多个 API 版本（向后兼容）。

### 执行命令

```bash
# 生成 v1 版本
/api-generator \
  --aggregate User \
  --prefix /api/v1 \
  --output ./api-v1

# 生成 v2 版本（新特性）
/api-generator \
  --aggregate User \
  --prefix /api/v2 \
  --output ./api-v2
```

### 版本差异处理

```go
// v1 Handler
type UserHandlerV1 struct {
	service *service.UserService
}

// v2 Handler (支持更多字段)
type UserHandlerV2 struct {
	service *service.UserService
}

// v2 新增字段
type CreateUserRequestV2 struct {
	CreateUserRequestV1
	PhoneNumber string `json:"phone_number" binding:"omitempty,e164"`
	Timezone    string `json:"timezone" binding:"omitempty"`
}
```

### 路由注册

```go
v1 := r.Group("/api/v1")
{
	http.RegisterUserRoutesV1(v1, userHandlerV1)
}

v2 := r.Group("/api/v2")
{
	http.RegisterUserRoutesV2(v2, userHandlerV2)
}
```

---

这些示例展示了 API Generator 的各种使用场景。你可以根据实际需求组合使用不同的选项和配置。
