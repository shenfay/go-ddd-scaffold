# Go DDD Scaffold - 通用的 DDD 项目脚手架模板

<div align="center">

[![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=flat&logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/Framework-React-20232A?style=flat&logo=react&logoColor=61DAFB)](https://reactjs.org/)
[![DDD](https://img.shields.io/badge/Architecture-DDD-8A2BE2?style=flat)](https://dddcommunity.org/)
[![Clean Architecture](https://img.shields.io/badge/Pattern-Clean_Architecture-4EC04E?style=flat)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

🎯 **基于领域驱动设计和洁净架构的通用项目模板**

</div>

## 🌟 项目简介

Go DDD Scaffold 是一个通用的 DDD 项目脚手架模板，提供完整的 Clean Architecture 实现。它帮助开发者快速启动标准化的 DDD 项目，内置多租户、RBAC 权限、用户管理等通用模块。

### 核心特性
- 🏗️ **标准 DDD 架构** - 四层架构（Domain/Application/Infrastructure/Interface）
- 🔧 **通用技术栈** - Go + Gin + GORM + PostgreSQL / React + TypeScript
- 👥 **多租户支持** - 内置 SaaS 多租户架构和 RBAC 权限控制
- 🎨 **示例领域** - User/Tenant 完整实现展示 DDD 最佳实践
- 🚀 **自动化工具** - Qoder Skills 智能代码生成

## 🏗️ 技术架构

### 后端技术栈
- **语言**: Go
- **架构**: 领域驱动设计 (DDD) + 洁净架构 (Clean Architecture)
- **Web 框架**: Gin + GORM
- **数据库**: PostgreSQL (支持迁移管理)
- **认证**: JWT + Casbin RBAC 权限控制
- **API**: Swagger 文档自动生成
- **实时通信**: WebSocket 连接池管理

### 前端技术栈
- **框架**: React 18 + TypeScript
- **架构**: 五层架构 (表现层、交互层、业务层、数据层、共享层)
- **样式**: Tailwind CSS
- **状态管理**: Zustand
- **3D 渲染**: Three.js + React Three Fiber

## 🚀 快速开始

### 后端启动

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 配置环境
cp config/config.yaml.example config/config.yaml

# 启动服务
go run cmd/server/main.go
```

### 前端启动

```bash
# 进入前端目录
cd frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev
```

## 🛠️ 核心功能模块

### 1. 用户与租户管理
- 🏠 **多租户架构** - 支持家庭组织单元
- 👨‍👩‍👧‍👦 **家庭成员管理** - 家长、学生角色权限
- 🔐 **RBAC 权限控制** - 基于 Casbin 的细粒度权限

### 2. 知识图谱系统
- 📚 **C/S/T/P 节点体系** - 概念、支撑、思维、问题四类节点
- 🧭 **学习路径规划** - 基于知识依赖的个性化路径
- 📊 **学习诊断** - 智能分析学习薄弱环节

### 3. 智能学习引擎
- 🎯 **6 阶段学习闭环** - 讲解-练习-测验-诊断-支线-成就
- 🎮 **游戏化元素** - 任务、徽章、排行榜激励机制
- 🤖 **AI 辅导** - 智能答疑与学习建议

### 4. 实时互动
- ⚡ **WebSocket 通信** - 学习状态实时同步
- 👥 **多人协作** - 在线答题竞赛、小组学习
- 🔔 **即时通知** - 学习进度、成就提醒

## 🎯 学习内容体系

### 知识领域世界观
- 🌳 **数学知识树** - 从数与代数到图形与几何
- 🧩 **Lv1-Lv5 能力等级** - 感知级到创新级进阶
- 🎭 **情境化学习** - 《小小超市大冒险》等主题场景

### 教学闭环设计
1. **讲解** - 动画演示、概念阐释
2. **练习** - 分层递进、即时反馈
3. **测验** - 定期评估、能力检测
4. **诊断** - 弱点分析、个性建议
5. **支线** - 补强训练、拓展提升
6. **成就** - 奖励激励、持续动力

## 🤖 智能化特性

### LLM 驱动的 NPC
- 🎭 **个性化角色** - 不同性格、特长的虚拟伙伴
- 💬 **自然对话** - 上下文感知的智能问答
- 🧠 **学习记忆** - 记录学习轨迹与偏好
- 🎯 **自适应教学** - 根据能力调整难度

### 智能推荐系统
- 🎯 **内容推荐** - 基于学习进度的个性化内容
- ⚖️ **难度平衡** - 保持挑战与能力匹配
- 📈 **进度预测** - 预估学习时间与效果

## 📊 技术亮点

### 自动化工具链
- 🔧 **Qoder Skills** - 自动化代码生成工具集
- 🏗️ **ddd-scaffold** - DDD 项目脚手架生成
- 🗃️ **db-migrator** - 数据库迁移与 DAO 生成
- 👥 **tenant-builder** - 多租户架构快速搭建
- 🔌 **websocket-integration** - WebSocket 快速集成
- 📄 **api-generator** - RESTful API 端点自动生成
- 📖 **api-doc-generator** - API 文档自动生成
- 🛡️ **error-handler-builder** - 统一错误处理构建

### 性能优化
- ⚡ **连接池管理** - 高效 WebSocket 连接复用
- 🚀 **缓存策略** - 多层级缓存提升响应速度
- 📊 **监控指标** - Prometheus 指标收集

## 📁 项目结构

```
go-ddd-scaffold/
├── backend/              # Go 后端服务
│   ├── cmd/             # 应用入口
│   ├── internal/        # 内部实现
│   │   ├── domain/      # 领域层 (DDD)
│   │   ├── application/ # 应用层
│   │   ├── interfaces/  # 接口层
│   │   └── infrastructure/ # 基础设施层
│   └── migrations/      # 数据库迁移
├── frontend/             # React 前端应用
│   ├── presentation/    # 表现层
│   ├── interaction/     # 交互层  
│   ├── business/        # 业务层
│   ├── data/           # 数据层
│   └── shared/         # 共享层
├── docs/                # 项目文档
└── .qoder/             # Qoder 工具配置
    ├── agents/         # AI 代理
    └── skills/         # 自动化技能
```

## 🚀 部署

### Docker 部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

## 🤝 贡献

我们欢迎各种形式的贡献：

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送分支
5. 创建 Pull Request

### 开发规范

- 使用 DDD 和洁净架构模式
- 遵循 Go 代码规范
- 编写单元测试
- 更新相关文档

## 📄 许可证

MIT License - 详见 [LICENSE](./LICENSE) 文件

## 📞 支持

如有任何问题，请通过以下方式联系我们：

- 🐛 [Issues](https://github.com/shenfay/math-fun/issues) - 报告 Bug 或提出功能请求
- 💬 [Discussions](https://github.com/shenfay/math-fun/discussions) - 讨论想法和建议

---

<div align="center">

**让每个孩子都能享受数学的乐趣！** 🌟

[⭐ Star 项目](https://github.com/shenfay/math-fun) [🐛 报告问题](https://github.com/shenfay/math-fun/issues) [🤝 参与贡献](https://github.com/shenfay/math-fun/pulls)

</div>