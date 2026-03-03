#!/bin/bash
# API Generator 使用示例脚本

set -e

echo "======================================"
echo "API Generator 使用示例"
echo "======================================"
echo ""

# 示例 1: 基础用法
echo "📦 示例 1: 快速生成用户管理 API"
echo "--------------------------------------"
cat << 'EOF'
命令:
/api-generator --aggregate UserAggregate --with-validation --auth jwt

生成的端点:
POST   /api/v1/users          # 创建用户
GET    /api/v1/users/:id      # 获取用户详情
GET    /api/v1/users          # 用户列表
PUT    /api/v1/users/:id      # 更新用户
DELETE /api/v1/users/:id      # 删除用户

特点:
- JWT 认证保护所有端点
- 请求参数自动验证
- Swagger 文档注解
- 统一响应格式
EOF
echo ""
echo ""

# 示例 2: 电商系统
echo "📦 示例 2: 完整电商系统 API"
echo "--------------------------------------"
cat << 'EOF'
命令:
/api-generator \\
  --aggregates Product,Order,Inventory \\
  --with-validation \\
  --auth jwt \\
  --with-tests

生成的模块:
✓ 商品管理 (Product)
  - CRUD 操作
  - 批量导入导出
  - 库存状态展示

✓ 订单管理 (Order)
  - 创建订单
  - 订单状态流转
  - 订单查询（支持多条件筛选）

✓ 库存管理 (Inventory)
  - 库存查询
  - 库存预占/释放
  - 库存预警

特性:
- 事务支持
- 并发控制
- 乐观锁机制
EOF
echo ""
echo ""

# 示例 3: 博客系统
echo "📦 示例 3: 博客系统 API"
echo "--------------------------------------"
cat << 'EOF'
命令:
/api-generator \\
  --aggregates Post,Comment,Tag \\
  --prefix /api/v1/blog \\
  --with-validation \\
  --auth jwt

特殊功能:
- Markdown 内容支持
- 标签多对多关系
- 评论嵌套结构
- 浏览量统计
- 草稿箱功能

自定义端点:
POST   /api/v1/blog/posts/:id/publish    # 发布文章
POST   /api/v1/blog/posts/:id/archive    # 归档文章
GET    /api/v1/blog/posts/:id/statistics # 查看统计数据
EOF
echo ""
echo ""

# 示例 4: RBAC 权限控制
echo "📦 示例 4: 基于角色的访问控制"
echo "--------------------------------------"
cat << 'EOF'
命令:
/api-generator --aggregate Document --auth casbin

权限模型:
角色：admin, member, guest

权限配置:
admin:   全部权限 (CRUD)
member:  读取 + 创建
guest:   仅读取

实现方式:
- Casbin RBAC 引擎
- 租户级权限隔离
- 动态策略加载
- 权限缓存优化
EOF
echo ""
echo ""

# 示例 5: 文件上传
echo "📦 示例 5: 文件上传 API"
echo "--------------------------------------"
cat << 'EOF'
命令:
# 先生成基础 API
/api-generator --aggregate File --with-validation --auth jwt

# 手动扩展上传方法
编辑：internal/interfaces/http/file/file_handler.go

支持的验证:
✓ 文件类型检查
✓ 文件大小限制
✓ 病毒扫描集成
✓ 图片尺寸验证

存储选项:
- 本地文件系统
- MinIO 对象存储
- AWS S3
- 阿里云 OSS
EOF
echo ""
echo ""

# 示例 6: 版本化 API
echo "📦 示例 6: API 版本管理"
echo "--------------------------------------"
cat << 'EOF'
命令:
# v1 版本
/api-generator --aggregate User --prefix /api/v1 --output ./api-v1

# v2 版本（新特性）
/api-generator --aggregate User --prefix /api/v2 --output ./api-v2

版本差异:
v1: 基础用户管理
v2: 新增手机号、时区、头像等字段

向后兼容:
- 同时支持 v1 和 v2
- v1 不推荐但继续工作
- 迁移指南文档
EOF
echo ""
echo ""

# 后续步骤
echo "🚀 生成完成后的操作"
echo "--------------------------------------"
cat << 'EOF'
# 1. 注册路由
编辑 cmd/server/main.go，注册生成的路由

# 2. 配置 JWT
export JWT_SECRET="your-secret-key"

# 3. 生成 API 文档
swag init -g cmd/server/main.go

# 4. 启动应用
make run

# 5. 查看文档
浏览器访问：http://localhost:8080/swagger/index.html

# 6. 测试 API
curl -X POST http://localhost:8080/api/v1/users \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN" \\
  -d '{"name":"John","email":"john@example.com"}'
EOF
echo ""
echo ""

echo "======================================"
echo "💡 提示"
echo "======================================"
cat << 'EOF'
- 首次使用建议从单个聚合开始
- 启用 --with-validation 确保数据安全
- 生产环境必须启用 --auth jwt
- 查看 EXAMPLES.md 了解更多场景
- 参考 REFERENCE.md 获取详细配置
- 咨询 API Design Agent 优化设计
EOF
echo ""
