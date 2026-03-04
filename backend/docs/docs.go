// Package docs API 文档
//
// @title           Go DDD Scaffold API
// @version         1.0
// @description     Go DDD Scaffold 通用脚手架 API，基于领域驱动设计（DDD）架构
// @description     
// @description     ## 核心功能
// @description     - **用户认证**: JWT Token 认证、多租户支持
// @description     - **权限管理**: RBAC 基于角色的访问控制
// @description     - **Token 黑名单**: Redis 缓存、限流熔断保护
// @description     - **事件总线**: Redis Stream 持久化、重试机制
// @description     
// @description     ## 认证说明
// @description     所有受保护的接口需要在 Header 中添加：
// @description     ```
// @description     Authorization: Bearer {access_token}
// @description     X-Tenant-ID: {tenant_id}  // 多租户场景
// @description     ```
// @description     
// @termsOfService http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 在值中输入 "Bearer {token}"

// @securityDefinitions.apikey TenantAuth
// @in header
// @name X-Tenant-ID
// @description 租户 ID（多租户场景可选）

// @externalDocs.description  OpenAPI Specification
// @externalDocs.url          https://swagger.io/resources/open-api/

package docs

import "github.com/swaggo/swag"

// SwaggerInfo 存储 Swagger 信息
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "Go DDD Scaffold API",
	Description:      "Go DDD Scaffold 通用脚手架 API，基于领域驱动设计（DDD）架构",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

const docTemplate = `
openapi: 3.0.0
info:
  title: Go DDD Scaffold API
  description: |
    Go DDD Scaffold 通用脚手架 API，基于领域驱动设计（DDD）架构
    
    ## 核心功能
    - **用户认证**: JWT Token 认证、多租户支持
    - **权限管理**: RBAC 基于角色的访问控制
    - **Token 黑名单**: Redis 缓存、限流熔断保护
    - **事件总线**: Redis Stream 持久化、重试机制
    
    ## 认证说明
    所有受保护的接口需要在 Header 中添加：
    - Authorization: Bearer {access_token}
    - X-Tenant-ID: {tenant_id} （多租户场景）
  version: 1.0.0
servers:
  - url: http://localhost:8080/api
    description: Development server
  - url: https://api.example.com/api
    description: Production server
`
