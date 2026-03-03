# API文档生成器技术参考

## 目录结构规范

```
.qoder/skills/api-doc-generator/
├── SKILL.md              # 主技能定义文件（必需，大写）
├── REFERENCE.md          # 技术参考文档（大写）
├── EXAMPLES.md           # 使用示例文档（大写）
├── QUICKSTART.md         # 快速入门指南
├── config.yaml           # 配置文件
├── scripts/              # 辅助脚本目录
│   └── helper.sh         # 主要辅助脚本
└── templates/            # 模板文件目录
    ├── swagger.tmpl      # Swagger模板
    └── html.tmpl         # HTML文档模板
```

## 核心组件详解

### 1. 注释解析引擎

#### 支持的注释标签
```go
// @Summary 简要描述
// @Description 详细描述
// @Tags 标签分类
// @Accept 接受的MIME类型
// @Produce 返回的MIME类型
// @Param 参数定义
// @Success 成功响应
// @Failure 错误响应
// @Router 路由定义
// @Security 安全方案
// @Deprecated 标记为废弃
```

#### 参数定义语法
```go
// @Param {name} {in} {type} {required} "{description}" {schema}
// 示例：
// @Param id path int true "用户ID" 
// @Param page query int false "页码" default(1)
// @Param user body User true "用户信息"
```

#### 响应定义语法
```go
// @Success {code} {type} {model} "{description}"
// @Failure {code} {type} {model} "{description}"
// 示例：
// @Success 200 {object} User "用户信息"
// @Failure 404 {object} ErrorResponse "用户不存在"
```

### 2. 配置管理系统

#### 完整配置结构
```yaml
api:
  scan_paths:              # 扫描路径列表
    - backend/internal/interfaces/http
    - backend/cmd/server
  
  output:                  # 输出配置
    formats: [json, yaml, html]
    directory: backend/docs
    filename: swagger
  
  swagger:                 # Swagger配置
    version: "2.0"
    info:
      title: "API文档"
      description: "接口文档描述"
      version: "1.0.0"
      contact:
        name: "开发团队"
        email: "dev@example.com"
    host: "api.example.com"
    base_path: "/v1"
    schemes: [https]
  
  parsing:                 # 解析配置
    include_examples: true
    parse_structs: true
    required_tags: [summary]

validation:                # 验证规则
  enabled: true
  rules:
    - require_summary
    - require_description
    - validate_response_types

custom_fields:             # 自定义字段
  - x-business-owner: "业务负责人"
  - x-api-version: "API版本"
```

### 3. 模板系统

#### Swagger模板变量
```
{{.Title}}              # 文档标题
{{.Version}}            # API版本
{{.Description}}        # 文档描述
{{.Host}}               # 主机地址
{{.BasePath}}           # 基础路径
{{.Paths}}              # 路径定义
{{.Definitions}}        # 数据模型定义
```

#### HTML模板结构
```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Info.Title}}</title>
    <meta charset="utf-8">
</head>
<body>
    <header>
        <h1>{{.Info.Title}}</h1>
        <p>{{.Info.Description}}</p>
    </header>
    
    <nav>
        <!-- 目录导航 -->
    </nav>
    
    <main>
        {{range $path, $methods := .Paths}}
            <section class="endpoint">
                <h2>{{$path}}</h2>
                {{range $method, $operation := $methods}}
                    <div class="operation {{$method}}">
                        <!-- 操作详情 -->
                    </div>
                {{end}}
            </section>
        {{end}}
    </main>
</body>
</html>
```

## API参考

### 核心命令

#### /api-generate-docs
生成完整的API文档
```
/api-generate-docs [--format {json|yaml|html}] [--package {path}] [--output {path}]
```

#### /api-update-docs
增量更新API文档
```
/api-update-docs [--force] [--clean]
```

#### /api-validate-docs
验证API文档质量和完整性
```
/api-validate-docs [--strict] [--report {path}]
```

#### /api-preview-docs
预览生成的API文档
```
/api-preview-docs [--port 8080]
```

### 辅助命令

#### /api-list-endpoints
列出所有API端点
```
/api-list-endpoints [--filter {tag}]
```

#### /api-check-comments
检查代码注释完整性
```
/api-check-comments [--fix] [--verbose]
```

#### /api-compare-versions
比较不同版本的API差异
```
/api-compare-versions {version1} {version2}
```

## 注释最佳实践

### 1. 控制器方法注释
```go
// UserHandler 用户相关接口处理器
type UserHandler struct{}

// @Summary 获取用户列表
// @Description 分页获取用户信息列表，支持多种筛选条件
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) minimum(1)
// @Param size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param status query string false "用户状态" Enums(active,inactive,suspended)
// @Success 200 {object} UserListResponse "成功返回用户列表"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
    // 实现逻辑
}
```

### 2. 数据模型注释
```go
// User 用户信息
type User struct {
    // @Description 用户唯一标识
    ID uint64 `json:"id" example:"12345"`
    
    // @Description 用户名，用于登录
    Username string `json:"username" example:"john_doe" validate:"required,min=3,max=20"`
    
    // @Description 邮箱地址
    Email string `json:"email" example:"john@example.com" validate:"required,email"`
    
    // @Description 用户状态
    // @Enum active,inactive,suspended
    Status string `json:"status" example:"active" validate:"required"`
    
    // @Description 创建时间
    CreatedAt time.Time `json:"created_at" example:"2026-01-01T00:00:00Z"`
    
    // @Description 更新时间
    UpdatedAt time.Time `json:"updated_at" example:"2026-01-01T00:00:00Z"`
}
```

### 3. 请求/响应结构体
```go
// LoginRequest 登录请求
type LoginRequest struct {
    // @Description 用户邮箱
    Email string `json:"email" example:"user@example.com" validate:"required,email"`
    
    // @Description 用户密码
    Password string `json:"password" example:"password123" validate:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    // @Description 访问令牌
    AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    
    // @Description 令牌类型
    TokenType string `json:"token_type" example:"Bearer"`
    
    // @Description 过期时间（秒）
    ExpiresIn int `json:"expires_in" example:"3600"`
    
    // @Description 用户基本信息
    User UserInfo `json:"user"`
}
```

## 错误处理机制

### 错误分类
1. **配置错误** - 配置文件缺失或格式错误
2. **解析错误** - 注释语法错误或不完整
3. **生成错误** - 模板渲染或文件写入失败
4. **验证错误** - 文档质量检查失败

### 错误处理策略
- 提供详细的错误位置和原因
- 支持自动修复常见问题
- 生成错误报告和建议
- 区分警告和致命错误

## 性能优化

### 扫描优化
- 增量扫描只处理变更文件
- 并行处理多个包
- 缓存已解析的注释
- 跳过无关的文件类型

### 生成优化
- 模板预编译
- 批量文件写入
- 内存使用优化
- 进度状态显示

## 安全考虑

### 访问控制
- 配置文件权限管理
- 敏感信息过滤
- 输出目录权限控制

### 数据保护
- 避免暴露敏感端点
- 过滤内部API文档
- 环境变量安全处理

## 扩展机制

### 插件架构
支持自定义扩展：
- 自定义注释解析器
- 额外的输出格式
- 自定义验证规则
- 第三方集成

### 集成能力
- CI/CD流水线集成
- Git钩子集成
- 监控系统集成
- 文档门户集成

---
*本文档遵循Qoder Skills技术规范，定期更新以反映最新功能和最佳实践*