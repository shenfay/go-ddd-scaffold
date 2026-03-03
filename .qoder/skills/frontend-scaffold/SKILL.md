---
name: frontend-scaffold
description: React 前端五层架构脚手架生成工具。基于 Vite + Tailwind CSS + Zustand，自动生成表现层、交互层、业务逻辑层、数据层、共享层的完整项目结构。适用于教育科技前端快速搭建。
version: "1.0.0"
author: MathFun Team
tags: [react, vite, tailwind, zustand, five-layer, scaffold, frontend, education]
---

# Frontend Scaffold - React 五层架构脚手架生成工具

## 功能概述

这是一个智能化的 React 前端脚手架生成工具，专为 MathFun 项目设计。它基于 **五层架构体系**（表现层、交互层、业务逻辑层、数据层、共享层），使用 **Vite** 构建工具、**Tailwind CSS** 样式框架和 **Zustand** 状态管理，提供完整的前端项目结构和最佳实践模板。

## 核心能力

### 1. 五层架构生成
- **表现层 (Presentation)** - 纯 UI 组件、页面、路由、主题
- **交互层 (Interaction)** - Three.js 场景、交互 Hooks、容器组件
- **业务逻辑层 (Business)** - Zustand Store、领域模型、业务服务
- **数据层 (Data)** - API 调用、Repository、数据映射
- **共享层 (Shared)** - 工具函数、常量、国际化、类型定义

### 2. 技术栈集成
- **Vite** - 快速冷启动、热模块替换
- **React 18** - 声明式 UI、函数组件、Hooks
- **Tailwind CSS** - 原子化 CSS、响应式设计
- **Zustand** - 轻量级全局状态管理
- **React Router** - 客户端路由
- **Framer Motion** - 动画效果
- **Three.js / R3F** - 3D 图形渲染

### 3. 儿童友好设计
- **明亮色彩** - 温暖、活泼的配色方案
- **大字体** - 适合儿童阅读的字号
- **圆角设计** - 安全、友好的视觉风格
- **动画反馈** - 即时、有趣的交互反馈
- **可爱图标** - 简洁、卡通化的图标系统

### 4. 多端统一
- **响应式布局** - 适配桌面、平板、手机
- **Taro 集成** - 小程序支持（可选）
- **PWA 支持** - 离线访问和安装体验
- **统一组件** - 一套代码多端运行

## 使用场景

### 适用情况
- 教育科技前端项目初始化
- 需要五层架构的 React 应用
- 包含 3D 交互的学习产品
- 多端统一的 Web 应用
- 儿童友好的界面设计

### 不适用情况
- 简单的后台管理系统
- 不需要 3D 功能的纯展示网站
- 企业官网等静态站点

## 基本使用

### 快速开始（5 分钟）

```bash
# 1. 生成新项目
frontend-scaffold create --name mathfun-student --template education

# 2. 进入目录
cd mathfun-student

# 3. 安装依赖
pnpm install

# 4. 启动开发服务器
pnpm dev
```

### 添加新页面

```bash
# 添加主页
frontend-scaffold add page --name Home --layout MainLayout

# 添加学习页面
frontend-scaffold add page --name Learning --layout LearningLayout

# 添加个人中心
frontend-scaffold add page --name Profile --layout AuthLayout
```

### 添加 3D 场景

```bash
# 添加 Three.js 场景
frontend-scaffold add scene --name FractionWorld --type interactive

# 添加 NPC 角色
frontend-scaffold add npc --name MathGuide --model humanoid
```

### 添加状态管理

```bash
# 添加 Zustand Store
frontend-scaffold add store --name learningProgress --slices progress,achievements

# 添加业务服务
frontend-scaffold add service --name LearningFlow --methods start,complete,evaluate
```

## 参数说明

| 参数 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|--------|------|
| `--name` | string | 是 | - | 项目名称 |
| `--template` | string | 否 | education | 模板类型 (education/dashboard/landing) |
| `--typescript` | flag | 否 | true | 使用 TypeScript |
| `--tailwind` | flag | 否 | true | 集成 Tailwind CSS |
| `--threejs` | flag | 否 | false | 集成 Three.js |
| `--pwa` | flag | 否 | false | 启用 PWA 支持 |

## 生成的代码结构

```
my-app/
├── src/
│   ├── presentation/          # 表现层
│   │   ├── pages/            # 页面组件
│   │   ├── components/       # UI 基础组件
│   │   ├── routes/           # 路由配置
│   │   └── styles/           # 全局样式
│   │
│   ├── interaction/          # 交互层
│   │   ├── three/            # Three.js 场景
│   │   ├── components/       # 交互容器组件
│   │   └── hooks/            # 交互 Hooks
│   │
│   ├── business/             # 业务逻辑层
│   │   ├── store/            # Zustand Store
│   │   ├── models/           # 领域模型
│   │   └── services/         # 业务服务
│   │
│   ├── data/                 # 数据层
│   │   ├── api/              # API 调用
│   │   ├── endpoints/        # 端点定义
│   │   ├── mappers/          # 数据映射
│   │   └── repositories/     # 数据仓库
│   │
│   └── shared/               # 共享层
│       ├── utils/            # 工具函数
│       ├── constants/        # 常量定义
│       ├── locales/          # 国际化
│       ├── types/            # 类型定义
│       └── hooks/            # 通用 Hooks
│
├── public/                    # 静态资源
├── package.json
├── vite.config.js
├── tailwind.config.js
└── tsconfig.json
```

## 最佳实践

### 组件命名规范

```jsx
// ✅ 正确：使用 PascalCase
import HomePage from './pages/HomePage';
import LearningCard from './components/LearningCard';

// ❌ 错误：避免 camelCase
import homePage from './pages/homePage';
```

### 五层职责划分

```jsx
// ✅ 表现层：只负责渲染
function LearningPage() {
  const { currentTask } = useLearningStore();
  return <div>{currentTask.title}</div>;
}

// ✅ 交互层：处理用户操作
function useLearningInteraction() {
  const completeTask = useLearningStore(state => state.completeTask);
  
  const handleComplete = useCallback(() => {
    // 播放音效、动画等交互反馈
    playSound('success');
    showAnimation('confetti');
    
    // 调用业务逻辑
    completeTask();
  }, []);
  
  return { handleComplete };
}

// ✅ 业务逻辑层：管理状态和规则
const useLearningStore = create((set) => ({
  currentTask: null,
  completeTask: () => {
    // 业务规则：完成任务后更新进度
    set((state) => ({
      progress: state.progress + 10
    }));
  }
}));
```

### Tailwind CSS 使用

```jsx
// ✅ 正确：儿童友好的设计
<button className="
  bg-blue-400 hover:bg-blue-500
  text-white font-bold text-lg
  px-8 py-4 rounded-2xl
  shadow-lg transform hover:scale-105
  transition-all duration-200
">
  开始学习
</button>

// ❌ 错误：成人化设计
<button className="bg-gray-500 text-sm px-4 py-2 rounded">
  Click Me
</button>
```

## 故障排除

### 常见问题

**Vite 启动失败**
- 检查 Node.js 版本（>= 16）
- 删除 node_modules 重新安装
- 检查 vite.config.js 配置

**Tailwind 样式不生效**
- 确认 tailwind.config.js 包含正确的内容路径
- 检查 CSS 文件是否引入 @tailwind 指令
- 清除浏览器缓存

**Zustand 状态不更新**
- 确保在组件外部使用 store
- 检查 selector 是否正确
- 避免直接修改状态，始终使用 setState

### 获取帮助
- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- 🚀 快速开始：查看 [QUICKSTART.md](./QUICKSTART.md)

## 版本历史

- v1.0.0 (2026-02-25): 初始版本发布
  - 完整的五层架构生成
  - Vite + Tailwind + Zustand 集成
  - 儿童友好设计模板
  - Three.js 场景支持
  - 多端统一配置

---
*本技能遵循 Qoder Skills 规范，专为 MathFun 项目优化设计*
