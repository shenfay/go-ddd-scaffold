---
name: api-generator
description: 智能 API 端点生成器。根据 DDD 模型自动生成 RESTful API 端点，包含完整的 CRUD 操作、参数验证、JWT 认证和 Swagger 文档注解。适用于快速构建标准化的后端 API。
version: "2.0.0"
author: MathFun Team
tags: [api, rest, crud, swagger, validation, automation, jwt, authentication]
---

# API Generator - API 端点生成器

## 功能概述

这是一个智能化的 API 端点生成工具，基于 MathFun 项目的最佳实践设计。它能够根据 DDD 领域模型自动生成完整的 RESTful API 端点，包括标准的 CRUD 操作、参数验证逻辑、JWT 认证中间件和 Swagger 文档注解。

## 核心能力

### 1. 智能端点生成
- **基于 DDD 聚合** - 自动生成标准 CRUD 端点
- **批量操作支持** - 支持批量创建、更新、删除
- **复合查询生成** - 支持多条件查询、分页、排序
- **Swagger 文档** - 自动生成完整的 API 文档注解
- **参数验证** - 内置请求参数验证和错误处理
- **JWT 认证** - 自动集成 JWT Token 认证中间件

### 2. 代码质量保障
- **RESTful 设计** - 遵循 REST 架构规范
- **统一响应格式** - 标准化的响应结构
- **DTO 自动生成** - 输入输出数据传输对象
- **单元测试模板** - 包含完整的测试用例
- **中间件集成** - 认证、授权、限流中间件

### 3. 灵活配置
- **自定义路径** - 可配置端点路径和前缀
- **HTTP 方法控制** - 指定支持的 HTTP 方法
- **字段级验证** - 灵活的验证规则配置
- **响应格式定制** - 自定义状态码和响应体
- **分页排序** - 内置分页和排序参数

## 使用场景

### 适用情况
- DDD 聚合对应的 API 端点生成
- 标准化 REST API 快速开发
- 需要 JWT 认证的项目
- 团队 API 规范统一
- 快速原型验证

### 不适用情况
- 非 RESTful 风格的 API（如 GraphQL）
- 需要复杂业务逻辑的端点
- 高度定制化的接口需求
- WebSocket 端点生成

## 基本使用

### 快速开始
```bash
# 基于 DDD 聚合生成 API 端点
/api-generator --aggregate UserAggregate

# 指定输出目录
/api-generator --aggregate Product --output ./custom/path

# 自定义端点配置
/api-generator --config custom-api-config.yaml
```

### 高级用法
```bash
# 批量生成多个聚合的 API
/api-generator --aggregates "User,Order,Product"

# 生成带验证的 API 端点
/api-generator --aggregate User --with-validation

# 生成完整测试套件
/api-generator --aggregate Order --with-tests

# 启用 JWT 认证
/api-generator --aggregate User --auth jwt
```

### 参数说明
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--aggregate` | string | 是 | - | 聚合根名称 |
| `--aggregates` | array | 否 | - | 多个聚合根列表 |
| `--output` | string | 否 | ./generated | 输出目录 |
| `--with-validation` | flag | 否 | false | 包含参数验证 |
| `--with-tests` | flag | 否 | false | 包含测试文件 |
| `--auth` | string | 否 | none | 认证方式 (none/jwt/casbin) |
| `--prefix` | string | 否 | /api/v1 | API 路径前缀 |

## 生成的端点示例

### 输入：UserAggregate

```bash
/api-generator --aggregate UserAggregate --with-validation --auth jwt
```

### 输出：标准 CRUD 端点

```go
// internal/interfaces/http/user/user_handler.go

// @Summary 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body dto.CreateUserRequest true "用户信息"
// @Success 201 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(err))
        return
    }
    
    // JWT 认证获取当前用户 ID
    userID := c.GetString("user_id")
    
    // 调用应用服务
    user, err := h.userService.CreateUser(c.Request.Context(), req, userID)
    if err != nil {
        c.JSON(response.HTTPStatus(err), response.Error(err))
        return
    }
    
    c.JSON(http.StatusCreated, response.Success(user))
}

// @Summary 获取用户详情
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(response.HTTPStatus(err), response.Error(err))
        return
    }
    
    c.JSON(http.StatusOK, response.Success(user))
}

// @Summary 用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=dto.ListUsersResponse}
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
    
    users, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
    if err != nil {
        c.JSON(response.HTTPStatus(err), response.Error(err))
        return
    }
    
    c.JSON(http.StatusOK, response.Success(users))
}

// @Summary 更新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Param user body dto.UpdateUserRequest true "用户信息"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
    id := c.Param("id")
    
    var req dto.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(err))
        return
    }
    
    user, err := h.userService.UpdateUser(c.Request.Context(), id, req)
    if err != nil {
        c.JSON(response.HTTPStatus(err), response.Error(err))
        return
    }
    
    c.JSON(http.StatusOK, response.Success(user))
}

// @Summary 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "用户 ID"
// @Success 204
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
    id := c.Param("id")
    
    err := h.userService.DeleteUser(c.Request.Context(), id)
    if err != nil {
        c.JSON(response.HTTPStatus(err), response.Error(err))
        return
    }
    
    c.Status(http.StatusNoContent)
}
```

## 配置说明

技能使用 `.qoder/skills/api-generator/config.yaml` 进行配置：

```yaml
api_generation:
  # RESTful 配置
  rest:
    base_path: "/api/v1"
    id_param: ":id"
    pagination:
      default_limit: 20
      max_limit: 100
      
  # 验证配置
  validation:
    enabled: true
    rules:
      - "required"
      - "min_length"
      - "max_length"
      - "email"
      - "phone"
      - "uuid"
  
  # 认证配置
  authentication:
    type: "jwt"  # none | jwt | casbin
    jwt:
      header: "Authorization"
      token_prefix: "Bearer"
      claims_key: "user_id"
  
  # 文档配置
  documentation:
    swagger_version: "2.0"
    generate_examples: true
    include_security: true
    produce_types:
      - "application/json"
    consume_types:
      - "application/json"

# 集成配置
integrations:
  ddd_modeling_assistant: "ddd-modeling-assistant"
  error_handler_builder: "error-handler-builder"
  api_doc_generator: "api-doc-generator"
```

## 最佳实践

1. **先进行 DDD 建模** - 使用 ddd-scaffold 完成领域建模
2. **明确 API 需求** - 确定需要暴露的业务能力
3. **合理配置验证** - 根据业务场景设置验证规则
4. **及时更新文档** - 保持 API 文档与代码同步
5. **统一错误处理** - 使用统一的错误响应格式

## 故障排除

### 常见问题

**端点生成失败**
- 检查 DDD 聚合定义是否完整
- 确认配置文件格式正确
- 验证输出目录权限

**验证规则不生效**
- 检查验证规则语法
- 确认字段映射正确
- 验证依赖包是否安装

**JWT 认证不工作**
- 确认 JWT 中间件已注册
- 检查 Token 格式是否正确
- 验证 Claims 解析逻辑

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

## 版本历史

- v2.0.0 (2026-02-25): 重大更新
  - 新增 JWT 认证支持
  - 增强参数验证功能
  - 优化 Swagger 文档生成
  - 统一错误处理机制
  
- v1.0.0 (2026-01-15): 初始版本
  - 基础 CRUD 端点生成
  - Swagger 文档注解
  - 参数验证支持

---
*本技能遵循 Qoder Skills 规范，专为 MathFun 项目 API 开发优化设计*
