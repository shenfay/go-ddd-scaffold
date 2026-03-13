# Swagger API 文档使用指南

## 快速开始

### 1. 生成 Swagger 文档

```bash
# 使用 Makefile
make swagger-gen

# 或直接运行 swag 命令
swag init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal
```

### 2. 启动服务并访问 Swagger UI

```bash
# 启动 API 服务
make run

# 访问 Swagger UI
# 开发环境下，访问：http://localhost:8080/swagger/index.html
```

## 已实现的接口

### 认证模块 (/api/v1/auth)

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /auth/login | 用户登录 |
| POST | /auth/register | 用户注册 |
| POST | /auth/refresh | 刷新令牌 |
| POST | /auth/logout | 用户登出 |

### 用户管理模块 (/api/v1/users)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /users | 列出用户（分页） |
| POST | /users | 创建用户 |
| GET | /users/{id} | 获取用户详情 |
| PUT | /users/{id} | 更新用户信息 |
| POST | /users/{id}/activate | 激活用户 |
| POST | /users/{id}/deactivate | 禁用用户 |
| PUT | /users/{id}/password | 修改用户密码 |

## 添加新接口的 Swagger 注释

### 基本格式

在 Handler 方法上方添加注释：

```go
// @Summary 接口摘要
// @Description 详细描述
// @Tags 模块名称
// @Accept json
// @Produce json
// @Param id path string true "参数说明"
// @Success 200 {object} httpShared.APIResponse{data=YourResponseType}
// @Failure 400 {object} httpShared.APIResponse
// @Router /your/path [post]
func (h *Handler) YourMethod(c *gin.Context) {
    // ...
}
```

### 常用注释标签

- `@Summary`: 接口简要说明
- `@Description`: 详细描述
- `@Tags`: 分组标签
- `@Accept`: 请求 Content-Type
- `@Produce`: 响应 Content-Type
- `@Param`: 请求参数
  - 格式：`@Param 参数名 参数位置 参数类型 是否必需 "说明"`
  - 参数位置：`path`, `query`, `body`, `header`, `form`
- `@Success`: 成功响应
  - 格式：`@Success 状态码 {对象} 返回类型`
- `@Failure`: 失败响应
- `@Router`: 路由路径和方法
- `@Security`: 认证方式（如 `BearerAuth`）

### 定义响应模型

在 handler 文件末尾添加：

```go
// ==================== Swagger 模型定义 ====================

// YourResponse 你的响应结构
// @Description 响应数据描述
type YourResponse struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
}
```

## 安全认证

本项目使用 JWT Bearer Token 认证。在 Swagger UI 中测试时：

1. 点击右上角 "Authorize" 按钮
2. 在 Value 中输入：`Bearer your_token_here`
3. 点击 "Authorize"
4. 现在可以测试需要认证的接口

## 文件结构

```
backend/docs/swagger/
├── docs.go          # Go 代码形式的文档
├── swagger.json     # JSON 格式的 OpenAPI 规范
└── swagger.yaml     # YAML 格式的 OpenAPI 规范
```

## 注意事项

1. **仅开发环境可用**: Swagger UI 只在 `gin.DebugMode` 下启用
2. **及时更新注释**: 修改接口后，记得更新 Swagger 注释并重新生成
3. **生产环境禁用**: 生产环境建议关闭 Swagger UI，避免泄露 API 信息

## 故障排查

### Swagger UI 无法访问

1. 确认服务已启动且在开发模式下运行
2. 检查是否已运行 `make swagger-gen` 生成文档
3. 确认路由已正确注册：`http://localhost:8080/swagger/*any`

### 文档未更新

1. 清理旧的文档：`rm -rf docs/swagger/*`
2. 重新生成：`make swagger-gen`
3. 重启服务

### 类型找不到错误

确保所有响应类型都已定义，并且使用了正确的包路径引用。对于共享类型，在 `response.go` 中统一添加 Swagger 模型定义。

## 参考资源

- [Swag 官方文档](https://github.com/swaggo/swag)
- [Gin-Swagger 中间件](https://github.com/swaggo/gin-swagger)
- [OpenAPI 规范](https://swagger.io/specification/)
