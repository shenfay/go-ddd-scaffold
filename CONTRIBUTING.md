# 项目贡献指南

## 开发流程

### 1. 分支管理

- **main**: 生产环境分支，只有稳定的代码才能合并
- **develop**: 开发主分支，日常开发在此分支
- **feature/***: 功能分支，从 develop 创建，完成后合并回 develop
- **hotfix/***: 热修复分支，从 main 创建，修复后同时合并到 main 和 develop

### 2. 提交规范

遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

**Type 类型**:
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是 bug 修复）
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建过程或辅助工具变动

**示例**:
```
feat(user): 添加用户注册功能

- 实现用户注册接口
- 添加邮箱验证
- 完善错误处理

Closes #123
```

### 3. CI/CD 流程

#### 本地开发

```bash
# 后端
cd backend
go mod download
go build ./...
go test -v ./...

# 前端
cd frontend
npm install
npm run build
npm test
```

#### 自动化流程

所有 PR 会自动触发以下检查：

1. **代码编译**: 确保代码能够正常编译
2. **单元测试**: 所有测试必须通过
3. **代码覆盖率**: 后端覆盖率 > 80%
4. **性能基准**: 与当前版本对比，性能下降不超过 5%
5. **Swagger 生成**: 确保 API文档能正常生成

#### 部署流程

- **Staging 环境**: 合并到 main 分支自动部署
- **Production 环境**: 手动打 tag 触发部署

```bash
git tag v1.0.0
git push origin v1.0.0
```

### 4. 代码审查清单

#### 后端 (Go)

- [ ] 代码编译通过
- [ ] 所有测试通过（包括基准测试）
- [ ] 新增代码有对应的单元测试
- [ ] 错误处理完整（使用 AppError）
- [ ] 日志记录合理（使用 zap）
- [ ] Swagger注释完整
- [ ] 无内存泄漏风险
- [ ] 遵循 DDD 分层架构

#### 前端 (React)

- [ ] 代码编译通过
- [ ] 所有测试通过
- [ ] 组件有 PropTypes 或 TypeScript 类型定义
- [ ] 无控制台警告
- [ ] 响应式布局适配
- [ ] 可访问性（a11y）检查通过

### 5. 数据库迁移

使用 [Goose](https://github.com/pressly/goose) 管理数据库迁移：

```bash
# 创建新迁移
goose -dir migrations/sql create add_user_avatar_field sql

# 执行迁移
goose -dir migrations/sql up

# 查看状态
goose -dir migrations/sql status
```

**迁移文件命名规范**:
```
YYYYMMDDHHMMSS_description.sql
```

### 6. 环境变量配置

#### 后端环境变量

```bash
# 数据库
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=go_ddd_scaffold

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=72

# 服务器配置
SERVER_PORT=8080
SERVER_MODE=debug  # debug, release, test
```

#### 前端环境变量

```bash
# API 地址
REACT_APP_API_URL=http://localhost:8080/api

# 应用配置
REACT_APP_ENV=development
REACT_APP_VERSION=1.0.0
```

### 7. 故障排查

#### 常见问题

**CI/CD 失败**:
1. 检查 GitHub Actions logs
2. 本地复现问题
3. 修复后重新推送

**测试失败**:
1. 查看完整的错误信息
2. 本地运行相同测试
3. 检查测试数据依赖

**部署失败**:
1. 检查 Docker 镜像构建日志
2. 验证 Kubernetes 配置
3. 查看应用日志

### 8. 性能优化

#### 后端性能

- 使用基准测试识别瓶颈
- 优化数据库查询（添加索引）
- 合理使用缓存（Redis）
- 避免不必要的内存分配

#### 前端性能

- 使用 React.memo 优化渲染
- 代码分割（Code Splitting）
- 图片懒加载
- 减少不必要的重渲染

### 9. 安全注意事项

- 敏感信息使用环境变量
- 密码加密存储（bcrypt）
- JWT token 定期轮换
- SQL 注入防护（使用 GORM）
- XSS 防护（前端输入过滤）

### 10. 监控与告警

#### 后端监控

- Prometheus 指标采集
- Grafana 仪表盘
- 错误率监控
- 性能指标监控

#### 前端监控

- 页面加载时间
- JavaScript 错误率
- 用户行为分析

---

## 快速开始

### 第一次贡献

1. Fork 项目
2. Clone 到本地
3. 创建功能分支
4. 开发并测试
5. 提交代码
6. 创建 Pull Request

### 环境搭建

```bash
# 克隆项目
git clone https://github.com/your-username/ddd-scaffold.git
cd ddd-scaffold

# 后端
cd backend
go mod download
cp config/config.yaml.example config/config.yaml
# 编辑配置文件
go run cmd/server/main.go

# 前端
cd ../frontend
npm install
npm start
```

---

感谢你的贡献！🎉
