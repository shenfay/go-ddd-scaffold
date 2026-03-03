# API端点生成器技术参考

## 目录结构规范

```
.qoder/skills/api-endpoint-generator/
├── SKILL.md              # 主技能定义文件（必需，大写）
├── REFERENCE.md          # 技术参考文档（大写）
├── EXAMPLES.md           # 使用示例文档（大写）
├── QUICKSTART.md         # 快速入门指南
├── config.yaml           # 配置文件
├── scripts/              # 辅助脚本目录
│   └── helper.sh         # 主要辅助脚本
└── templates/            # 模板文件目录
    ├── handler.tmpl      # HTTP处理器模板
    ├── dto.tmpl          # DTO模板
    ├── router.tmpl       # 路由注册模板
    └── swagger.tmpl      # Swagger注解模板
```

## 核心组件详解

### 1. 代码生成引擎

#### 模板系统架构
```go
type ApiGenerator struct {
    TemplateEngine  *template.Template
    DddModelReader  *DddModelReader
    ConfigManager   *ConfigManager
    OutputWriter    *OutputWriter
}

type GenerationContext struct {
    AggregateName   string
    ResourceName    string
    Fields          []FieldDefinition
    Operations      []OperationType
    ValidationRules map[string][]ValidationRule
    Config          *ApiConfig
}
```

#### 端点生成器
```go
func (g *ApiGenerator) GenerateEndpoints(ctx *GenerationContext) error {
    // 生成处理器
    if err := g.generateHandlers(ctx); err != nil {
        return err
    }
    
    // 生成DTO
    if err := g.generateDTOs(ctx); err != nil {
        return err
    }
    
    // 生成路由
    if err := g.generateRoutes(ctx); err != nil {
        return err
    }
    
    // 生成Swagger文档
    if err := g.generateSwaggerDocs(ctx); err != nil {
        return err
    }
    
    return nil
}
```

### 2. DDD模型解析器

#### 聚合解析
```go
type DddModelReader struct {
    projectPath string
}

func (r *DddModelReader) ReadAggregate(aggregateName string) (*AggregateInfo, error) {
    // 解析聚合根定义
    aggregateRoot := r.parseAggregateRoot(aggregateName)
    
    // 解析实体定义
    entities := r.parseEntities(aggregateName)
    
    // 解析值对象
    valueObjects := r.parseValueObjects(aggregateName)
    
    return &AggregateInfo{
        Name:         aggregateName,
        Root:         aggregateRoot,
        Entities:     entities,
        ValueObjects: valueObjects,
    }, nil
}
```

### 3. 验证规则引擎

#### 验证规则定义
```yaml
validation_rules:
  - name: "required"
    description: "字段必填"
    template: "if {{.FieldName}} == \"\" { return errors.New(\"{{.FieldName}} is required\") }"
  
  - name: "min_length"
    description: "最小长度验证"
    template: "if len({{.FieldName}}) < {{.MinLength}} { return errors.New(\"{{.FieldName}} too short\") }"
  
  - name: "max_length"
    description: "最大长度验证"
    template: "if len({{.FieldName}}) > {{.MaxLength}} { return errors.New(\"{{.FieldName}} too long\") }"
  
  - name: "email"
    description: "邮箱格式验证"
    template: "if !emailRegex.MatchString({{.FieldName}}) { return errors.New(\"invalid email format\") }"
```

### 4. 配置管理系统

#### 完整配置结构
```yaml
api_generation:
  # RESTful API配置
  rest:
    base_path: "/api/v1"
    version: "v1"
    id_parameter: ":id"
    
    # HTTP方法映射
    methods:
      create: "POST"
      read: "GET"
      update: "PUT"
      delete: "DELETE"
      list: "GET"
    
    # 响应配置
    responses:
      success_status: 200
      created_status: 201
      not_found_status: 404
      validation_error_status: 400
      internal_error_status: 500
  
  # 分页配置
  pagination:
    enabled: true
    default_limit: 20
    max_limit: 100
    page_parameter: "page"
    limit_parameter: "limit"
  
  # 排序配置
  sorting:
    enabled: true
    default_field: "created_at"
    default_direction: "desc"
    field_parameter: "sort"
    direction_parameter: "order"

# 集成配置
integrations:
  ddd_modeling_assistant:
    enabled: true
    skill_name: "ddd-modeling-assistant"
    auto_discover: true
  
  error_handler_builder:
    enabled: true
    skill_name: "error-handler-builder"
    apply_standard_errors: true
  
  api_doc_generator:
    enabled: true
    skill_name: "api-doc-generator"
    auto_invoke: true
    output_formats: ["json", "yaml"]

# 模板配置
templates:
  handler_template: "templates/handler.tmpl"
  dto_template: "templates/dto.tmpl"
  router_template: "templates/router.tmpl"
  swagger_template: "templates/swagger.tmpl"
```

## API参考

### 核心命令

#### /api-generate-endpoints
生成RESTful API端点
```
/api-generate-endpoints --aggregate {aggregate_name} [--output {path}] [--config {config_file}]
```

#### /api-list-aggregates
列出可用的DDD聚合
```
/api-list-aggregates [--filter {pattern}]
```

#### /api-validate-config
验证API生成配置
```
/api-validate-config [--config {config_file}]
```

#### /api-preview-endpoints
预览将要生成的端点
```
/api-preview-endpoints --aggregate {aggregate_name}
```

### 辅助命令

#### /api-update-docs
更新API文档
```
/api-update-docs [--aggregate {aggregate_name}] [--force]
```

#### /api-generate-tests
生成API测试代码
```
/api-generate-tests --aggregate {aggregate_name} [--test-type unit|integration]
```

## 生成代码结构

### 1. Handler结构
```go
// UserHandler 用户API处理器
type UserHandler struct {
    userService *UserService
    validator   *Validator
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.handleError(c, err)
        return
    }
    
    if err := h.validator.Validate(req); err != nil {
        h.handleValidationError(c, err)
        return
    }
    
    user, err := h.userService.CreateUser(req.ToDomain())
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, CreateUserResponse{User: user})
}
```

### 2. DTO结构
```go
// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=2,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
    User UserDTO `json:"user"`
}

// UserDTO 用户数据传输对象
type UserDTO struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
```

### 3. 路由注册
```go
// RegisterUserRoutes 注册用户路由
func RegisterUserRoutes(router *gin.RouterGroup, handler *UserHandler) {
    users := router.Group("/users")
    {
        users.POST("", handler.CreateUser)
        users.GET("", handler.ListUsers)
        users.GET("/:id", handler.GetUser)
        users.PUT("/:id", handler.UpdateUser)
        users.DELETE("/:id", handler.DeleteUser)
    }
}
```

## 错误处理机制

### 标准错误响应格式
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数验证失败",
    "details": [
      {
        "field": "email",
        "message": "邮箱格式不正确"
      }
    ]
  }
}
```

### 错误码定义
```go
const (
    // 4xx 客户端错误
    ErrCodeValidationError = "VALIDATION_ERROR"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeForbidden       = "FORBIDDEN"
    
    // 5xx 服务器错误
    ErrCodeInternalServerError = "INTERNAL_ERROR"
    ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
)
```

## 性能优化

### 代码生成优化
- 模板预编译缓存
- 并行文件生成
- 增量更新支持
- 生成进度可视化

### 运行时优化
- 连接池配置
- 查询优化建议
- 缓存策略推荐
- 监控指标集成

## 安全考虑

### 输入验证
- 参数绑定安全检查
- SQL注入防护
- XSS攻击防范
- 文件上传安全

### 访问控制
- JWT令牌验证
- RBAC权限控制
- 速率限制
- CORS配置

## 扩展机制

### 插件架构
支持自定义扩展：
- 自定义验证规则
- 额外的中间件
- 第三方认证集成
- 自定义响应格式

### 模板扩展
- 自定义代码模板
- 多语言支持
- 框架适配器
- 代码风格配置

---
*本文档遵循Qoder Skills技术规范，定期更新以反映最新功能和最佳实践*