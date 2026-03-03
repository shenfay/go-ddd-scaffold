# Frontend Scaffold 快速开始指南

## 5 分钟创建 React 五层架构项目

### 第一步：安装 Skill

```bash
npx skills install frontend-scaffold
```

### 第二步：创建新项目

```bash
frontend-scaffold create \
  --name mathfun-student \
  --template education \
  --threejs \
  --tailwind
```

这会生成一个完整的教育科技 React 应用，包含：
- ✅ 五层架构目录结构
- ✅ Vite + Tailwind CSS 配置
- ✅ Three.js 集成（用于 3D 场景）
- ✅ Zustand 状态管理
- ✅ 儿童友好的设计系统

### 第三步：安装依赖

```bash
cd mathfun-student
pnpm install
```

### 第四步：启动开发服务器

```bash
pnpm dev
```

浏览器会自动打开 http://localhost:3000

---

## 项目结构详解

生成的项目遵循严格的**五层架构**：

```
mathfun-student/
├── src/
│   ├── presentation/          # 表现层
│   │   ├── pages/            # 页面组件
│   │   │   └── HomePage.jsx  # 主页示例
│   │   ├── components/ui/    # UI 基础组件
│   │   ├── components/layouts/ # 布局组件
│   │   ├── routes/           # 路由配置
│   │   ├── styles/           # 全局样式
│   │   └── assets/           # 静态资源
│   │
│   ├── interaction/          # 交互层
│   │   ├── three/scenes/     # Three.js 场景
│   │   ├── three/models/     # 3D 模型
│   │   ├── components/       # 交互容器组件
│   │   └── hooks/            # 交互 Hooks
│   │
│   ├── business/             # 业务逻辑层
│   │   ├── store/            # Zustand Store
│   │   │   └── slices/       # 状态切片
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
│       ├── locales/zh-CN/    # 国际化
│       ├── types/            # 类型定义
│       └── hooks/            # 通用 Hooks
│
├── public/                    # 公共资源
├── package.json
├── vite.config.js
├── tailwind.config.js
└── tsconfig.json
```

---

## 核心文件说明

### vite.config.js

```javascript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@presentation': path.resolve(__dirname, './src/presentation'),
      '@interaction': path.resolve(__dirname, './src/interaction'),
      '@business': path.resolve(__dirname, './src/business'),
      '@data': path.resolve(__dirname, './src/data'),
      '@shared': path.resolve(__dirname, './src/shared'),
    },
  },
  server: {
    port: 3000,
    open: true,
  },
})
```

### tailwind.config.js

```javascript
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 儿童友好的配色
        primary: {
          400: '#60a5fa',  // 主要按钮颜色
          500: '#3b82f6',
          600: '#2563eb',
        },
        success: '#22c55e',
        warning: '#f59e0b',
        error: '#ef4444',
      },
      fontSize: {
        'base': '18px',  // 适合儿童阅读
        'lg': '20px',
        'xl': '24px',
      },
      borderRadius: {
        'DEFAULT': '0.75rem',
        'lg': '1rem',
        'xl': '2rem',
      },
    },
  },
}
```

### src/main.jsx

```jsx
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './presentation/styles/index.css'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
```

### src/App.jsx

```jsx
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { HomePage } from '@presentation/pages/HomePage'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
      </Routes>
    </Router>
  )
}

export default App
```

---

## 添加新页面

### 使用命令添加

```bash
frontend-scaffold add page --name LearningCenter --layout MainLayout
```

### 手动创建

在 `src/presentation/pages/` 下创建 `LearningCenter.jsx`：

```jsx
// LearningCenter 页面
// 职责：纯 UI 渲染

import { useLearningStore } from '@business/store'

export default function LearningCenter() {
  const { currentTask } = useLearningStore()
  
  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-purple-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8 text-primary-600">
          学习中心
        </h1>
        
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <p className="text-lg mb-4">
            当前任务：{currentTask?.title || '暂无任务'}
          </p>
          
          <button className="
            bg-primary-400 hover:bg-primary-500
            text-white font-bold text-xl
            px-8 py-4 rounded-2xl
            shadow-lg transform hover:scale-105
            transition-all duration-200
          ">
            开始学习
          </button>
        </div>
      </div>
    </div>
  )
}
```

然后在 `src/App.jsx` 中添加路由：

```jsx
import LearningCenter from '@presentation/pages/LearningCenter'

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/learning" element={<LearningCenter />} />
      </Routes>
    </Router>
  )
}
```

---

## 添加 3D 场景

### 创建 Three.js 场景

在 `src/interaction/three/scenes/` 下创建 `FractionWorld.jsx`：

```jsx
import { Canvas } from '@react-three/fiber'
import { OrbitControls, Sky } from '@react-three/drei'

export default function FractionWorld() {
  return (
    <Canvas camera={{ position: [0, 2, 5], fov: 60 }}>
      {/* 天空盒 */}
      <Sky sunPosition={[100, 20, 100]} />
      
      {/* 环境光 */}
      <ambientLight intensity={0.5} />
      
      {/* 方向光 */}
      <directionalLight position={[10, 10, 5]} intensity={1} />
      
      {/* 3D 对象 */}
      <mesh position={[0, 1, 0]}>
        <boxGeometry args={[1, 1, 1]} />
        <meshStandardMaterial color="hotpink" />
      </mesh>
      
      {/* 轨道控制器 */}
      <OrbitControls />
    </Canvas>
  )
}
```

### 在页面中使用

```jsx
import FractionWorld from '@interaction/three/scenes/FractionWorld'

function LearningPage() {
  return (
    <div className="h-screen">
      <h1 className="text-3xl font-bold mb-4">分数世界</h1>
      <div className="h-[600px] rounded-2xl overflow-hidden">
        <FractionWorld />
      </div>
    </div>
  )
}
```

---

## 添加状态管理

### 创建 Zustand Store

在 `src/business/store/slices/` 下创建 `learningProgress.js`：

```javascript
// learningProgress Slice
export const createLearningProgressSlice = (set, get) => ({
  // State
  progress: 0,
  achievements: [],
  currentTask: null,
  
  // Actions
  setProgress: (value) => set({ progress: value }),
  addAchievement: (achievement) => 
    set((state) => ({ 
      achievements: [...state.achievements, achievement] 
    })),
  setCurrentTask: (task) => set({ currentTask: task }),
  completeTask: () => {
    const state = get()
    set({
      progress: state.progress + 10,
      currentTask: null,
    })
  },
})
```

### 合并到主 Store

在 `src/business/store/index.js`：

```javascript
import { create } from 'zustand'
import { subscribeWithSelector } from 'zustand/middleware'
import { createDefaultSlice } from './slices/default.js'
import { createLearningProgressSlice } from './slices/learningProgress.js'

const useStore = create(
  subscribeWithSelector((...args) => ({
    ...createDefaultSlice(...args),
    ...createLearningProgressSlice(...args),
  }))
)

// 选择器
export const useLearningStore = () => useStore((state) => state)
export const selectProgress = (state) => state.progress

export default useStore
```

### 在组件中使用

```jsx
import { useLearningStore } from '@business/store'

function ProgressDisplay() {
  const { progress, achievements, completeTask } = useLearningStore()
  
  return (
    <div>
      <p>进度：{progress}%</p>
      <p>成就数：{achievements.length}</p>
      <button onClick={completeTask}>完成任务</button>
    </div>
  )
}
```

---

## 儿童友好设计实践

### 明亮温暖的色彩

```jsx
<button className="
  bg-blue-400 hover:bg-blue-500
  text-white font-bold text-lg
  px-8 py-4 rounded-2xl
  shadow-lg
">
  开始学习
</button>
```

### 大字体易读

```jsx
<h1 className="text-4xl font-bold mb-4">
  欢迎来到数学王国！
</h1>
<p className="text-lg leading-relaxed">
  让我们一起探索奇妙的数学世界吧~
</p>
```

### 圆角安全设计

```jsx
<div className="rounded-2xl shadow-xl p-8 bg-white">
  {/* 内容区域 */}
</div>
```

### 动画反馈

```jsx
import { motion } from 'framer-motion'

<motion.button
  whileHover={{ scale: 1.05 }}
  whileTap={{ scale: 0.95 }}
  className="bg-primary-400 text-white px-8 py-4 rounded-2xl"
>
  点击我
</motion.button>
```

---

## 下一步

### 1. 完善业务逻辑

在 `src/business/services/` 下添加业务服务：

```javascript
// LearningService.js
export class LearningService {
  async startTask(taskId) {
    // 开始任务的 бизнес逻辑
  }
  
  async completeTask(taskId, answers) {
    // 完成任务并评估
  }
}
```

### 2. 集成后端 API

在 `src/data/api/` 下创建 API 客户端：

```javascript
// apiClient.js
import axios from 'axios'

const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

export const learningAPI = {
  getTasks: () => apiClient.get('/tasks'),
  submitAnswer: (taskId, answer) => 
    apiClient.post(`/tasks/${taskId}/submit`, { answer }),
}
```

### 3. 添加更多功能

- 🎨 主题切换（白天/夜晚模式）
- 🌍 国际化支持
- 📱 响应式布局优化
- 🔊 音效和语音反馈
- 📊 数据可视化图表

---

## 获取帮助

- 📖 详细文档：查看 [REFERENCE.md](./REFERENCE.md)
- 💡 使用示例：查看 [EXAMPLES.md](./EXAMPLES.md)
- ❓ 遇到问题：咨询 DDD Architect Agent

祝你开发顺利！🚀
