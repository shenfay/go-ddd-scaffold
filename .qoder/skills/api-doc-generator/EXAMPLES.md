# API文档生成器使用示例

## 基础使用示例

### 示例 1：用户认证API文档

**场景**：为用户登录注册功能生成API文档

```
/api-generate-docs --package backend/internal/interfaces/http/auth
```

**控制器代码示例**：
```go
package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// AuthHandler 认证相关接口
type AuthHandler struct{}

// @Summary 用户登录
// @Description 用户通过邮箱和密码进行身份验证，成功后返回JWT令牌
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录凭证"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "用户名或密码错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_REQUEST",
            Message: "请求参数格式错误",
        })
        return
    }
    
    // 验证用户凭证
    user, token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, ErrorResponse{
            Code:    "AUTH_FAILED",
            Message: "用户名或密码错误",
        })
        return
    }
    
    c.JSON(http.StatusOK, LoginResponse{
        AccessToken: token,
        TokenType:   "Bearer",
        ExpiresIn:   3600,
        User: UserInfo{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
        },
    })
}

// @Summary 用户注册
// @Description 新用户注册账户
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param register body RegisterRequest true "注册信息"
// @Success 201 {object} RegisterResponse "注册成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 409 {object} ErrorResponse "邮箱已被注册"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
    // 实现注册逻辑
}
```

**数据模型**：
```go
// LoginRequest 登录请求
type LoginRequest struct {
    // @Description 用户邮箱地址
    Email string `json:"email" example:"user@example.com" validate:"required,email"`
    
    // @Description 用户密码
    Password string `json:"password" example:"password123" validate:"required,min=6"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    // @Description JWT访问令牌
    AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
    
    // @Description 令牌类型
    TokenType string `json:"token_type" example:"Bearer"`
    
    // @Description 令牌有效期（秒）
    ExpiresIn int `json:"expires_in" example:"3600"`
    
    // @Description 用户基本信息
    User UserInfo `json:"user"`
}

// UserInfo 用户基本信息
type UserInfo struct {
    // @Description 用户ID
    ID uint64 `json:"id" example:"12345"`
    
    // @Description 用户名
    Username string `json:"username" example:"john_doe"`
    
    // @Description 邮箱地址
    Email string `json:"email" example:"user@example.com"`
}
```

### 示例 2：用户管理API文档

**场景**：为用户管理功能生成完整API文档

```
/api-generate-docs --format html --include-examples
```

**完整控制器示例**：
```go
package user

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
)

// UserHandler 用户管理接口
type UserHandler struct{}

// @Summary 获取用户列表
// @Description 分页获取用户信息列表，支持状态筛选和关键词搜索
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1) minimum(1)
// @Param size query int false "每页数量" default(10) minimum(1) maximum(100)
// @Param status query string false "用户状态" Enums(active,inactive,suspended)
// @Param search query string false "搜索关键词"
// @Success 200 {object} UserListResponse "成功返回用户列表"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
    // 实现分页查询逻辑
}

// @Summary 获取用户详情
// @Description 根据用户ID获取详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} UserDetailResponse "成功返回用户详情"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{
            Code:    "INVALID_ID",
            Message: "无效的用户ID",
        })
        return
    }
    
    // 获取用户详情逻辑
}

// @Summary 更新用户信息
// @Description 更新指定用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body UpdateUserRequest true "更新的用户信息"
// @Success 200 {object} UserDetailResponse "更新成功"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    // 实现更新逻辑
}

// @Summary 删除用户
// @Description 软删除指定用户（标记为已删除）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 204 "删除成功"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    // 实现删除逻辑
}
```

## 高级使用示例

### 示例 3：知识图谱API文档

**场景**：为MathFun的知识图谱相关API生成专业文档

```
/api-generate-docs --package backend/internal/interfaces/http/knowledge --format yaml
```

**复杂API示例**：
```go
package knowledge

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// KnowledgeHandler 知识图谱接口处理器
type KnowledgeHandler struct{}

// @Summary 获取知识点详情
// @Description 根据知识点ID获取详细信息和关联关系
// @Tags 知识图谱
// @Accept json
// @Produce json
// @Param id path string true "知识点ID" format(uuid)
// @Param expand query string false "展开关联数据" Enums(prerequisites,related,applications)
// @Success 200 {object} KnowledgeNodeDetail "成功返回知识点详情"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "知识点不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /knowledge/nodes/{id} [get]
func (h *KnowledgeHandler) GetNodeDetail(c *gin.Context) {
    // 实现获取节点详情逻辑
}

// @Summary 查询学习路径
// @Description 根据起始和目标知识点计算最优学习路径
// @Tags 知识图谱
// @Accept json
// @Produce json
// @Param path body LearningPathRequest true "路径查询参数"
// @Success 200 {object} LearningPathResponse "成功返回学习路径"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 422 {object} ErrorResponse "无法找到有效路径"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /knowledge/paths [post]
func (h *KnowledgeHandler) CalculateLearningPath(c *gin.Context) {
    // 实现路径计算逻辑
}

// @Summary 批量导入知识点
// @Description 批量导入知识点数据和关系
// @Tags 知识图谱
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "知识点数据文件" format(csv)
// @Param format formData string false "文件格式" Enums(csv,json) default(csv)
// @Success 202 {object} BatchImportResponse "开始处理导入任务"
// @Failure 400 {object} ErrorResponse "文件格式错误"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /knowledge/import [post]
func (h *KnowledgeHandler) BatchImportNodes(c *gin.Context) {
    // 实现批量导入逻辑
}
```

**复杂数据模型**：
```go
// KnowledgeNodeDetail 知识点详情
type KnowledgeNodeDetail struct {
    // @Description 知识点唯一标识
    ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000" format:"uuid"`
    
    // @Description 知识点名称
    Name string `json:"name" example:"分数的基本概念"`
    
    // @Description 知识点描述
    Description string `json:"description" example:"理解分数的含义和基本运算规则"`
    
    // @Description 知识点类型
    // @Enum concept,skill,problem,application
    NodeType string `json:"node_type" example:"concept"`
    
    // @Description 能力等级要求
    CompetencyLevel int `json:"competency_level" example:"2" minimum:"1" maximum:"5"`
    
    // @Description 前置知识点列表
    Prerequisites []string `json:"prerequisites,omitempty" example:"['basic_arithmetic']"`
    
    // @Description 相关知识点列表
    RelatedNodes []string `json:"related_nodes,omitempty"`
    
    // @Description 应用场景示例
    Applications []ApplicationExample `json:"applications,omitempty"`
    
    // @Description 创建时间
    CreatedAt string `json:"created_at" example:"2026-01-01T00:00:00Z" format:"date-time"`
    
    // @Description 更新时间
    UpdatedAt string `json:"updated_at" example:"2026-01-01T00:00:00Z" format:"date-time"`
}

// ApplicationExample 应用场景示例
type ApplicationExample struct {
    // @Description 场景标题
    Title string `json:"title" example:"购物中的分数应用"`
    
    // @Description 场景描述
    Description string `json:"description" example:"在超市购物时如何计算折扣价格"`
    
    // @Description 难度等级
    Difficulty int `json:"difficulty" example:"2" minimum:"1" maximum:"3"`
    
    // @Description 示例数据
    ExampleData map[string]interface{} `json:"example_data,omitempty"`
}
```

### 示例 4：学习进度API文档

**场景**：为学习进度跟踪功能生成API文档

```
/api-generate-docs --validate --include-examples
```

**带验证的API示例**：
```go
package learning

import (
    "net/http"
    "time"
    "github.com/gin-gonic/gin"
)

// LearningHandler 学习相关接口处理器
type LearningHandler struct{}

// @Summary 开始学习会话
// @Description 为指定知识点创建新的学习会话
// @Tags 学习管理
// @Accept json
// @Produce json
// @Param session body StartSessionRequest true "会话启动参数"
// @Success 201 {object} SessionResponse "会话创建成功"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "知识点不存在"
// @Failure 409 {object} ErrorResponse "已有活跃会话"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /learning/sessions [post]
func (h *LearningHandler) StartSession(c *gin.Context) {
    // 实现会话启动逻辑
}

// @Summary 提交练习答案
// @Description 提交练习题的答案并获得即时反馈
// @Tags 学习管理
// @Accept json
// @Produce json
// @Param submission body ExerciseSubmission true "练习提交"
// @Success 200 {object} ExerciseFeedback "答题反馈"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "练习不存在"
// @Failure 423 {object} ErrorResponse "会话已结束"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /learning/exercises/submit [post]
func (h *LearningHandler) SubmitExercise(c *gin.Context) {
    // 实现答案提交逻辑
}

// @Summary 结束学习会话
// @Description 结束当前学习会话并保存进度
// @Tags 学习管理
// @Accept json
// @Produce json
// @Param id path string true "会话ID" format(uuid)
// @Param end body EndSessionRequest false "会话结束信息"
// @Success 200 {object} SessionSummary "会话总结"
// @Failure 400 {object} ErrorResponse "参数错误"
// @Failure 404 {object} ErrorResponse "会话不存在"
// @Failure 500 {object} ErrorResponse "服务器内部错误"
// @Security ApiKeyAuth
// @Router /learning/sessions/{id}/end [post]
func (h *LearningHandler) EndSession(c *gin.Context) {
    // 实现会话结束逻辑
}
```

## 团队协作示例

### 示例 5：多人协作开发API文档

**场景**：团队成员协同开发API文档

```
# 团队成员A：生成基础文档
/api-generate-docs --package backend/internal/interfaces/http/user

# 团队成员B：添加示例数据
/api-generate-docs --include-examples --update

# 团队负责人：验证文档质量
/api-validate-docs --strict --report backend/docs/validation-report.md

# 全团队：预览最终效果
/api-preview-docs --port 8080
```

### 示例 6：CI/CD集成

**场景**：在持续集成流程中自动生成API文档

```yaml
# .github/workflows/api-docs.yml
name: API Documentation Generator

on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/internal/interfaces/http/**'
      - '.qoder/skills/api-doc-generator/**'

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Generate API Documentation
        run: |
          cd .qoder/skills/api-doc-generator
          /api-generate-docs --format json --output backend/docs/swagger.json
          /api-validate-docs --strict
      
      - name: Deploy Documentation
        run: |
          # 部署到文档服务器
          # 或推送到GitHub Pages
```

## 故障排除示例

### 示例 7：处理常见问题

**问题1：文档生成为空**
```
# 检查注释完整性
/api-check-comments --verbose

# 手动修复缺失注释
# 然后重新生成
/api-generate-docs --force
```

**问题2：注释语法错误**
```go
// 错误示例：
// @Param user_id path int true "用户ID"  // 缺少类型说明

// 正确示例：
// @Param user_id path int true "用户ID" minimum(1)
```

**问题3：响应模型未定义**
```go
// 确保响应结构体有正确的注释
type UserResponse struct {
    // @Description 用户信息
    User User `json:"user"`
    
    // @Description 操作时间戳
    Timestamp string `json:"timestamp" format:"date-time"`
}
```

### 示例 8：性能优化场景

**场景**：大型项目中提高文档生成效率

```
# 增量更新模式
/api-update-docs --incremental

# 只扫描变更的文件
/api-generate-docs --changed-only

# 并行处理多个包
/api-generate-docs --parallel --workers 4
```

这些示例展示了API文档生成器在各种实际场景中的灵活应用，从基础的REST API到复杂的教育领域专用接口，都能提供专业的文档生成支持。