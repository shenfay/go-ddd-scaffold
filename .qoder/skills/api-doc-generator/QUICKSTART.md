# API文档生成器快速入门

## 快速开始指南

### 1. 基本操作流程

#### 生成API文档
```
# 生成默认格式的API文档
/api-generate-docs

# 指定输出格式
/api-generate-docs --format html

# 指定扫描包路径
/api-generate-docs --package backend/internal/interfaces/http
```

#### 验证和更新文档
```
# 验证文档质量
/api-validate-docs

# 增量更新文档
/api-update-docs

# 预览生成的文档
/api-preview-docs --port 8080
```

### 2. 常用命令速查

| 命令 | 用途 | 示例 |
|------|------|------|
| `/api-generate-docs` | 生成API文档 | `/api-generate-docs --format json,yaml` |
| `/api-update-docs` | 更新现有文档 | `/api-update-docs --force` |
| `/api-validate-docs` | 验证文档质量 | `/api-validate-docs --strict` |
| `/api-preview-docs` | 预览文档 | `/api-preview-docs --port 8080` |

### 3. 典型使用场景

#### 场景1：新项目初始化
```
# 1. 初始化API文档
/api-generate-docs --format html --include-examples

# 2. 验证文档质量
/api-validate-docs --report validation-report.md

# 3. 启动预览服务
/api-preview-docs --port 8080
```

#### 场景2：日常开发更新
```
# 代码变更后更新文档
/api-update-docs

# 检查注释完整性
/api-check-comments --verbose
```

#### 场景3：团队协作开发
```
# 团队成员A：生成基础文档
/api-generate-docs --package user

# 团队成员B：添加详细描述
/api-update-docs --include-descriptions

# 团队负责人：质量检查
/api-validate-docs --strict
```

## 高级功能

### 自定义配置
```
# 使用自定义配置文件
/api-generate-docs --config custom-config.yaml

# 指定输出目录
/api-generate-docs --output backend/docs/v2
```

### 批量处理
```
# 处理多个包
/api-generate-docs --packages "user,auth,knowledge"

# 并行处理
/api-generate-docs --parallel --workers 4
```

### 集成开发
```
# 与现有构建流程集成
/api-generate-docs --silent --exit-code

# 生成差异报告
/api-compare-versions v1.0.0 v1.1.0
```

## 配置说明

技能使用 `.qoder/skills/api-doc-generator/config.yaml` 进行配置：

```yaml
api:
  scan_paths:
    - backend/internal/interfaces/http
  output:
    formats: [json, yaml, html]
    directory: backend/docs
  swagger:
    version: "2.0"
    info:
      title: "MathFun API"
      version: "1.0.0"
```

## 故障排除

### 常见问题

**文档生成为空**
```
/api-check-comments --verbose
# 检查代码中是否有有效的Swagger注释
```

**注释解析错误**
```
/api-validate-docs --detailed
# 查看详细的错误信息和位置
```

**输出文件权限问题**
```
# 确保输出目录存在且有写权限
mkdir -p backend/docs
chmod 755 backend/docs
```

### 获取帮助
请参阅详细文档：
- SKILL.md - 主技能文档
- REFERENCE.md - 技术参考
- EXAMPLES.md - 使用示例

## 最佳实践

1. **保持注释更新** - 代码变更时同步更新API注释
2. **定期验证质量** - 建立文档质量检查流程
3. **团队统一规范** - 制定团队注释编写标准
4. **版本化管理** - 重要变更时生成版本快照
5. **自动化集成** - 将文档生成集成到CI/CD流程

## 下一步

- 探索 [EXAMPLES.md](EXAMPLES.md) 了解详细的使用场景
- 查看 [REFERENCE.md](REFERENCE.md) 了解技术细节
- 检查 [SKILL.md](SKILL.md) 了解完整的技能文档

---
*如需支持，请联系 MathFun 开发团队*