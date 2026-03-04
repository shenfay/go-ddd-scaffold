# Swagger API文档使用指南

## 📖 概述

本项目使用 [Swag](https://github.com/swaggo/swag) 自动生成符合 OpenAPI/Swagger 2.0 规范的 API文档。文档包含完整的接口定义、参数说明、响应格式和认证方式。

---

## 🚀 快速开始

### 1. 安装 Swag 工具

```bash
# 安装 swag 命令行工具
go install github.com/swaggo/swag/cmd/swag@latest

# 验证安装
swag --version
```

### 2. 生成 API文档

```bash
cd backend

# 基本用法
swag init -g cmd/server/main.go -o ./docs

# 推荐配置（解析依赖和内部包）
swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal

# 生成特定格式
swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal --outputTypes json,yaml
```

### 3. 查看文档

启动服务后访问：

```bash
# 启动后端服务
go run cmd/server/main.go

# 浏览器访问
http://localhost:8080/swagger/index.html
```

或使用 Swagger UI 独立服务：

```bash
# 安装 swagger-ui 服务器
go install github.com/swaggo/http-swagger/example/docs-server@latest

# 启动文档服务器
docs-server --dir ./docs --port 8081

# 访问
http://localhost:8081
```

---

## 📝 Swagger注释规范

### Handler注释模板

```go
// @Summary      用户登录
// @Description  用户通过邮箱和密码进行身份验证，返回 JWT Token
// @Tags         用户认证
// @Accept       json
// @Produce      json
// @Router       /api/auth/login [post]
// @Param        request body dto.LoginRequest true "登录信息"
// @Success      200  {object}  response.Response{data=dto.LoginResponse}  "登录成功"
// @Failure      400  {object}  response.Response  "请求参数错误"
// @Failure      401  {object}  response.Response  "用户名或密码错误"
// @Failure      403  {object}  response.Response  "账户被禁用"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Security     BearerAuth
func (h *AuthHandler) Login(c *gin.Context) {
    // ...
}
```

### 常用注解说明

| 注解 | 说明 | 示例 |
|------|------|------|
| **@Summary** | 接口简短描述 | `@Summary 用户登录` |
| **@Description** | 接口详细描述 | `@Description 用户通过邮箱和密码进行身份验证` |
| **@Tags** | 接口分组标签 | `@Tags 用户认证` |
| **@Accept** | 请求格式 | `@Accept json` |
| **@Produce** | 响应格式 | `@Produce json` |
| **@Router** | 路由路径和方法 | `@Router /api/auth/login [post]` |
| **@Param** | 请求参数 | `@Param request body dto.LoginRequest true "登录信息"` |
| **@Success** | 成功响应 | `@Success 200 {object} dto.User "登录成功"` |
| **@Failure** | 失败响应 | `@Failure 400 {object} response.Response "参数错误"` |
| **@Security** | 安全认证 | `@Security BearerAuth` |

### 结构体字段注释

```go
type LoginRequest struct {
    Email    string `json:"email" example:"user@example.com" validate:"required,email"`     // 用户邮箱
    Password string `json:"password" example:"StrongPass123" validate:"required,min=8"`     // 密码（至少 8 位）
}
```

---

## 📂 生成的文件

执行 `swag init` 后会生成以下文件：

```
backend/docs/
├── docs.go          // Go代码，用于注册 Swagger 信息
├── swagger.json     // JSON 格式的 API文档
└── swagger.yaml     // YAML 格式的 API文档
```

### 文件大小参考

- **docs.go**: ~23 KB
- **swagger.json**: ~22 KB
- **swagger.yaml**: ~12 KB

---

## 🔧 Makefile 命令

项目提供了便捷的 Makefile 命令管理文档：

```bash
# 查看所有相关命令
make help

# 生成 Swagger 文档
make swagger

# 启动 Swagger UI 服务器
make swagger-serve

# 完整构建流程（包含文档生成）
make all
```

---

## 🎯 已生成的 API 列表

### 认证接口 (Auth)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | ❌ |
| POST | `/api/auth/login` | 用户登录 | ❌ |
| POST | `/api/auth/logout` | 用户登出 | ✅ |
| GET | `/api/auth/me` | 获取当前用户信息 | ✅ |

### 租户接口 (Tenant)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/tenants` | 创建租户 | ✅ |
| GET | `/api/tenants/{id}` | 获取租户信息 | ✅ |
| PUT | `/api/tenants/{id}` | 更新租户信息 | ✅ |

### 用户接口 (User)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/users` | 获取用户列表 | ✅ |
| GET | `/api/users/{id}` | 获取用户详情 | ✅ |
| PUT | `/api/users/{id}` | 更新用户信息 | ✅ |
| DELETE | `/api/users/{id}` | 删除用户 | ✅ |

---

## 🔐 安全认证配置

### Bearer Token 认证

在 `docs.go` 中已配置：

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在值中输入 "Bearer {token}"
```

### 多租户支持

```go
// @securityDefinitions.apikey TenantAuth
// @in header
// @name X-Tenant-ID
// @description 租户 ID（多租户场景可选）
```

### 使用示例

在 Swagger UI 中点击 "Authorize" 按钮，输入：

```
BearerAuth: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
TenantAuth: tenant-uuid-here
```

---

## 📊 Swagger.json 内容示例

```json
{
    "swagger": "2.0",
    "info": {
        "contact": {}
    },
    "paths": {
        "/api/auth/login": {
            "post": {
                "description": "用户登录接口",
                "consumes": ["application/json"],
                "produces": ["application/json"],
                "tags": ["auth"],
                "summary": "用户登录",
                "parameters": [
                    {
                        "description": "登录信息",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/go-ddd-scaffold_internal_application_user_dto.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "$ref": "#/definitions/go-ddd-scaffold_internal_application_user_dto.LoginResponse"
                        }
                    }
                }
            }
        }
    }
}
```

---

## 🛠️ 故障排查

### 常见问题

#### 1. Swag 命令未找到

**错误：**
```
zsh: command not found: swag
```

**解决方案：**
```bash
# 确保 GOPATH/bin 在 PATH 中
export PATH=$PATH:$(go env GOPATH)/bin

# 或者使用完整路径
$(go env GOPATH)/bin/swag init -g cmd/server/main.go -o ./docs
```

#### 2. 文档生成不完整

**错误：**
```
warning: failed to get package name in dir: ./
```

**解决方案：**
```bash
# 添加 --parseDependency 和 --parseInternal 参数
swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal
```

#### 3. 路由重复声明警告

**错误：**
```
warning: route POST /api/tenants is declared multiple times
```

**原因：** 同一个路由在多个地方声明

**解决方案：**
- 检查是否有重复的 `@Router` 注解
- 确保每个 Handler 的路由唯一

#### 4. 文档与代码不同步

**解决方案：**
```bash
# 每次修改 API 后重新生成文档
make swagger

# 或在 CI/CD 中自动执行
swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal
```

---

## 🔄 自动化集成

### Git Hooks 自动更新

创建 `.git/hooks/pre-commit`：

```bash
#!/bin/bash

echo "Generating Swagger documentation..."
cd backend
$(go env GOPATH)/bin/swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal

if [ $? -ne 0 ]; then
    echo "❌ Swagger generation failed!"
    exit 1
fi

echo "✅ Documentation updated!"
```

### GitHub Actions CI/CD

创建 `.github/workflows/swagger.yml`：

```yaml
name: Generate Swagger Docs

on:
  push:
    branches: [ main ]
    paths:
      - 'backend/internal/interfaces/**'
      - 'backend/cmd/**'

jobs:
  swagger:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install swag
        run: go install github.com/swaggo/swag/cmd/swag@latest
      
      - name: Generate Swagger docs
        run: |
          cd backend
          swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal
      
      - name: Commit changes
        uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: 'docs: auto-generate Swagger documentation'
          file_pattern: 'backend/docs/*'
```

---

## 📈 最佳实践

### 1. 注释完整性

- ✅ 每个公开 API 必须有 `@Summary` 和 `@Description`
- ✅ 所有参数必须标注是否必需 (`true/false`)
- ✅ 提供所有可能的响应码说明
- ✅ 添加 `@Tags` 进行分类

### 2. 示例数据

为所有 DTO 字段提供示例值：

```go
type CreateUserRequest struct {
    Email    string `json:"email" example:"user@example.com"`
    Password string `json:"password" example:"StrongPass123"`
    Nickname string `json:"nickname" example:"John Doe"`
}
```

### 3. 错误响应统一

所有错误响应使用统一格式：

```go
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
```

### 4. 定期审查

- 每月检查一次文档完整性
- 新增 API 时同步更新注释
- 删除废弃接口时清理文档

---

## 🎨 自定义配置

### 修改 Swagger 信息

编辑 `backend/docs/docs.go`：

```go
// @title           Go DDD Scaffold API
// @version         1.0
// @description     Go DDD Scaffold 通用脚手架 API，基于领域驱动设计（DDD）架构
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@example.com

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api
```

### 添加 Logo

在 Swagger UI 中添加项目 Logo：

```css
/* custom.css */
.swagger-ui .topbar .download-url-wrapper {
    display: none;
}

.swagger-ui .topbar .link {
    background-image: url('logo.png');
    background-size: contain;
}
```

---

## 📚 参考资源

### 官方文档

- [Swag GitHub](https://github.com/swaggo/swag)
- [Swagger 2.0 Specification](https://swagger.io/specification/v2/)
- [OpenAPI Specification](https://swagger.io/specification/)

### 注释语法参考

- [Declarative Comments Format](https://github.com/swaggo/swag#declarative-comments-format)

### 工具推荐

- **Swagger UI**: 可视化界面浏览和测试 API
- **Postman**: 导入 swagger.json 快速测试
- **Insomnia**: 轻量级 API 测试工具

---

## 🎉 总结

通过 Swag 自动生成 API文档，我们可以：

✅ **提高效率** - 无需手动维护文档  
✅ **保证准确** - 文档与代码实时同步  
✅ **便于协作** - 团队成员快速了解接口  
✅ **提升质量** - 标准化的接口描述  

**立即开始使用：**

```bash
make swagger
open http://localhost:8080/swagger/index.html
```
