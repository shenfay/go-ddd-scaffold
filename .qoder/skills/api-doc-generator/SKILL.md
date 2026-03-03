---
name: api-doc-generator
description: 自动生成API文档的智能工具。支持Swagger/OpenAPI规范，可从Go代码注释提取API信息，生成完整的API文档。适用于RESTful API项目文档维护和同步。
version: "1.0.0"
author: MathFun Team
tags: [api, documentation, swagger, openapi, go, swag]
---

# API文档生成器

## 功能概述

这是一个智能化的API文档生成工具，专门为MathFun项目设计。它能够自动扫描Go代码中的注释，提取API接口信息，并生成符合Swagger/OpenAPI 2.0/3.0规范的完整API文档。

## 核心能力

### 1. 智能注释解析
- 自动识别Go代码中的Swagger注释
- 支持标准的swag注解语法
- 提取路由、参数、响应等API元信息
- 支持结构体字段注释解析

### 2. 多格式文档生成
- Swagger 2.0 JSON/YAML格式
- OpenAPI 3.0 JSON/YAML格式
- 通过swag命令自动生成标准文档

### 3. 代码同步机制
- 自动检测代码变更
- 增量更新文档内容
- 保持文档与代码一致性
- 支持版本控制集成

### 4. 自定义模板支持
- 可配置的文档模板
- 支持品牌化定制
- 多语言文档支持
- 扩展字段定义

## 使用场景

### 适用情况
- RESTful API项目文档维护
- 微服务接口文档生成
- API设计规范检查
- 开发团队文档标准化
- 第三方开发者文档提供

### 不适用情况
- 非HTTP协议接口
- 二进制协议通信
- 简单的内部工具接口
- 临时调试接口

## 快速开始

### 基本使用流程

1. **扫描并生成文档**
   ```
   /api-generate-docs
   ```

2. **指定输出格式**
   ```
   /api-generate-docs --format html
   ```

3. **增量更新文档**
   ```
   /api-update-docs
   ```

4. **验证文档质量**
   ```
   /api-validate-docs
   ```

### 高级用法
```
# 生成特定包的文档
/api-generate-docs --package internal/interfaces/http

# 自定义输出路径
/api-generate-docs --output backend/docs/swagger.json

# 包含示例数据
/api-generate-docs --include-examples
```

## 配置说明

技能使用 `.qoder/skills/api-doc-generator/config.yaml` 进行配置：

```yaml
# API文档生成配置
api:
  # 扫描路径配置
  scan_paths:
    - backend/internal/interfaces/http
    - backend/cmd/server
    - backend/pkg/api
  
  # 输出配置
  output:
    formats: [json, yaml]  # 支持的输出格式
    directory: backend/docs      # 输出目录
    filename: swagger            # 基础文件名
  
  # Swagger配置
  swagger:
    version: "2.0"
    info:
      title: "MathFun API Documentation"
      description: "数学教育平台API接口文档"
      version: "1.0.0"
      contact:
        name: "MathFun开发团队"
        email: "dev@mathfun.com"
    
    # 服务器配置
    host: "api.mathfun.com"
    base_path: "/v1"
    schemes: [https]
  
  # 注释解析配置
  parsing:
    include_examples: true       # 是否包含示例
    parse_structs: true          # 是否解析结构体
    required_tags: [summary]     # 必需的注释标签

# 代码质量检查
validation:
  enabled: true
  rules:
    - require_summary          # 必需summary描述
    - require_description      # 必需详细描述
    - validate_response_types  # 验证响应类型
    - check_duplicate_routes   # 检查重复路由

# 自定义字段
custom_fields:
  - x-author: "开发者"
  - x-created-at: "创建时间"
  - x-business-domain: "业务领域"
```

## 最佳实践

### 注释规范
```go
// @Summary 用户登录
// @Description 用户通过邮箱和密码进行身份验证
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} LoginResponse "登录成功"
// @Failure 400 {object} ErrorResponse "请求参数错误"
// @Failure 401 {object} ErrorResponse "认证失败"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    // 实现逻辑
}
```

### 结构体注释
```go
type User struct {
    ID       uint64 `json:"id" example:"1"`                    // 用户ID
    Username string `json:"username" example:"john_doe"`       // 用户名
    Email    string `json:"email" example:"john@example.com"`  // 邮箱地址
    Created  string `json:"created_at" example:"2026-01-01T00:00:00Z"` // 创建时间
}
```

### 团队协作
- 统一注释风格和格式
- 定期同步API文档
- 代码审查时检查注释完整性
- 建立文档更新流程

## 故障排除

### 常见问题

**文档生成为空**
- 检查代码中是否有有效的Swagger注释
- 验证扫描路径配置是否正确
- 确认Go代码能否正常编译

**注释解析错误**
- 检查注释语法是否符合swag规范
- 验证标签拼写和格式
- 确认结构体字段映射正确

**输出格式问题**
- 检查配置文件中的格式设置
- 验证输出目录权限
- 确认依赖工具版本兼容性

### 获取帮助
- 查看详细文档：REFERENCE.md
- 参考使用示例：EXAMPLES.md
- 快速入门指南：QUICKSTART.md

## 版本历史

- v1.0.0 (2026-01-26): 初始版本发布
  - 基础API文档生成功能
  - 多格式输出支持
  - 注释解析和验证
  - 配置文件管理

---
*本技能遵循Qoder Skills规范，专为MathFun项目优化设计*