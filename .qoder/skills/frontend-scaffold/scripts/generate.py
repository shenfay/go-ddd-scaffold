#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Frontend Scaffold - React 五层架构脚手架生成工具
基于 Vite + Tailwind CSS + Zustand
"""

import os
import sys
import argparse
import json
from pathlib import Path
from datetime import datetime


def parse_args():
    """解析命令行参数"""
    parser = argparse.ArgumentParser(
        description='React 五层架构脚手架生成工具',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 创建新项目
  %(prog)s create --name mathfun-student --template education
  
  # 添加页面
  %(prog)s add page --name Home --layout MainLayout
  
  # 添加 3D 场景
  %(prog)s add scene --name FractionWorld --type interactive
  
  # 添加状态管理
  %(prog)s add store --name learningProgress --slices progress,achievements
        """
    )
    
    subparsers = parser.add_subparsers(dest='command', help='可用命令')
    
    # create 子命令
    create_parser = subparsers.add_parser('create', help='创建新项目')
    create_parser.add_argument('--name', type=str, required=True, help='项目名称')
    create_parser.add_argument('--template', type=str, default='education', 
                              choices=['education', 'dashboard', 'landing'],
                              help='模板类型')
    create_parser.add_argument('--typescript', action='store_true', default=True, help='使用 TypeScript')
    create_parser.add_argument('--tailwind', action='store_true', default=True, help='集成 Tailwind CSS')
    create_parser.add_argument('--threejs', action='store_true', default=False, help='集成 Three.js')
    create_parser.add_argument('--pwa', action='store_true', default=False, help='启用 PWA 支持')
    create_parser.add_argument('--output', type=str, default='./', help='输出目录')
    
    # add 子命令
    add_parser = subparsers.add_parser('add', help='添加新元素')
    add_subparsers = add_parser.add_subparsers(dest='action', help='添加类型')
    
    # add page
    page_parser = add_subparsers.add_parser('page', help='添加页面')
    page_parser.add_argument('--name', type=str, required=True, help='页面名称')
    page_parser.add_argument('--layout', type=str, default='MainLayout', help='布局组件')
    page_parser.add_argument('--path', type=str, help='路由路径')
    
    # add component
    comp_parser = add_subparsers.add_parser('component', help='添加组件')
    comp_parser.add_argument('--name', type=str, required=True, help='组件名称')
    comp_parser.add_argument('--layer', type=str, default='presentation',
                            choices=['presentation', 'interaction', 'business'],
                            help='所属层级')
    
    # add scene
    scene_parser = add_subparsers.add_parser('scene', help='添加 3D 场景')
    scene_parser.add_argument('--name', type=str, required=True, help='场景名称')
    scene_parser.add_argument('--type', type=str, default='interactive',
                             choices=['interactive', 'display', 'game'],
                             help='场景类型')
    
    # add npc
    npc_parser = add_subparsers.add_parser('npc', help='添加 NPC 角色')
    npc_parser.add_argument('--name', type=str, required=True, help='NPC 名称')
    npc_parser.add_argument('--model', type=str, default='humanoid', help='模型类型')
    
    # add store
    store_parser = add_subparsers.add_parser('store', help='添加 Zustand Store')
    store_parser.add_argument('--name', type=str, required=True, help='Store 名称')
    store_parser.add_argument('--slices', type=str, help='状态切片，逗号分隔')
    
    # add service
    service_parser = add_subparsers.add_parser('service', help='添加业务服务')
    service_parser.add_argument('--name', type=str, required=True, help='服务名称')
    service_parser.add_argument('--methods', type=str, help='方法列表，逗号分隔')
    
    return parser.parse_args()


def create_project_structure(name, template, output_dir):
    """创建项目基础结构"""
    base_path = Path(output_dir) / name
    
    # 五层架构目录
    layers = {
        'presentation': ['pages', 'components/ui', 'components/layouts', 'routes', 'styles', 'assets'],
        'interaction': ['three/scenes', 'three/models', 'components', 'hooks'],
        'business': ['store/slices', 'models', 'services'],
        'data': ['api', 'endpoints', 'mappers', 'repositories'],
        'shared': ['utils', 'constants', 'locales/zh-CN', 'types', 'hooks', 'themes']
    }
    
    print(f"\n🚀 创建项目：{name}\n")
    
    for layer, subdirs in layers.items():
        for subdir in subdirs:
            dir_path = base_path / 'src' / layer / subdir
            dir_path.mkdir(parents=True, exist_ok=True)
            print(f"✓ 创建目录：src/{layer}/{subdir}")
    
    # 创建其他必要目录
    (base_path / 'public').mkdir(parents=True, exist_ok=True)
    (base_path / 'public' / 'images').mkdir(parents=True, exist_ok=True)
    (base_path / 'public' / 'sounds').mkdir(parents=True, exist_ok=True)
    (base_path / 'public' / 'models3d').mkdir(parents=True, exist_ok=True)
    
    print("✓ 创建公共资源目录\n")
    
    return base_path


def generate_package_json(name, template):
    """生成 package.json"""
    package = {
        "name": name,
        "private": True,
        "version": "0.1.0",
        "type": "module",
        "scripts": {
            "dev": "vite",
            "build": "vite build",
            "preview": "vite preview",
            "lint": "eslint . --ext js,jsx --report-unused-disable-directives --max-warnings 0"
        },
        "dependencies": {
            "react": "^18.2.0",
            "react-dom": "^18.2.0",
            "react-router-dom": "^6.15.0",
            "zustand": "^4.4.0",
            "framer-motion": "^10.16.0",
            "axios": "^1.5.0",
            "katex": "^0.16.0"
        },
        "devDependencies": {
            "@types/react": "^18.2.0",
            "@types/react-dom": "^18.2.0",
            "@vitejs/plugin-react": "^4.0.0",
            "autoprefixer": "^10.4.15",
            "eslint": "^8.45.0",
            "eslint-plugin-react": "^7.32.2",
            "eslint-plugin-react-hooks": "^4.6.0",
            "eslint-plugin-react-refresh": "^0.4.3",
            "postcss": "^8.4.29",
            "tailwindcss": "^3.3.0",
            "vite": "^4.4.0"
        }
    }
    
    if template == 'education':
        package["dependencies"]["@react-three/fiber"] = "^8.15.0"
        package["dependencies"]["@react-three/drei"] = "^9.80.0"
        package["dependencies"]["three"] = "^0.155.0"
    
    return package


def generate_vite_config(name):
    """生成 vite.config.js"""
    config = f"""import {{ defineConfig }} from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({{
  plugins: [react()],
  resolve: {{
    alias: {{
      '@': path.resolve(__dirname, './src'),
      '@presentation': path.resolve(__dirname, './src/presentation'),
      '@interaction': path.resolve(__dirname, './src/interaction'),
      '@business': path.resolve(__dirname, './src/business'),
      '@data': path.resolve(__dirname, './src/data'),
      '@shared': path.resolve(__dirname, './src/shared'),
    }},
  }},
  server: {{
    port: 3000,
    open: true,
  }},
  build: {{
    outDir: 'build',
    sourcemap: true,
  }},
}})
"""
    return config


def generate_tailwind_config():
    """生成 tailwind.config.js"""
    config = """/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        // 儿童友好的配色方案
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',  // 主要按钮颜色
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
        },
        success: '#22c55e',  // 成功反馈
        warning: '#f59e0b',  // 警告提示
        error: '#ef4444',    // 错误提示
      },
      fontSize: {
        // 适合儿童阅读的字号
        'base': '18px',
        'lg': '20px',
        'xl': '24px',
        '2xl': '32px',
        '3xl': '40px',
      },
      borderRadius: {
        // 圆角设计，更安全友好
        'none': '0',
        'sm': '0.5rem',
        DEFAULT: '0.75rem',
        'md': '1rem',
        'lg': '1.5rem',
        'xl': '2rem',
        'full': '9999px',
      },
      animation: {
        'bounce-slow': 'bounce 2s infinite',
        'pulse-fast': 'pulse 1s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
    },
  },
  plugins: [],
}
"""
    return config


def generate_index_html(name):
    """生成 index.html"""
    html = f"""<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{name}</title>
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
"""
    return html


def generate_main_jsx():
    """生成 main.jsx"""
    code = """import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.jsx'
import './presentation/styles/index.css'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
"""
    return code


def generate_app_jsx():
    """生成 App.jsx"""
    code = """import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { HomePage } from '@presentation/pages/HomePage'
// 导入更多页面组件

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        {/* 添加更多路由 */}
      </Routes>
    </Router>
  )
}

export default App
"""
    return code


def generate_index_css():
    """生成 index.css"""
    code = """@tailwind base;
@tailwind components;
@tailwind utilities;

/* 全局样式重置和基础样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f0f9ff;
  color: #1e293b;
}

/* 儿童友好的滚动条样式 */
::-webkit-scrollbar {
  width: 12px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
"""
    return code


def create_sample_page(layer, name, layout):
    """创建示例页面"""
    page_code = f"""// {name} 页面组件
// 职责：纯 UI 渲染，不包含业务逻辑

import {{ useLearningStore }} from '@business/store'
import {{ useInteraction }} from '@interaction/hooks/useInteraction'

export default function {name}() {{
  // 从业务层获取数据
  const {{ currentTask }} = useLearningStore()
  
  // 使用交互层处理用户操作
  const {{ handleClick }} = useInteraction()
  
  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-purple-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8 text-primary-600">
          {name}
        </h1>
        
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <p className="text-lg mb-4">
            当前任务：{{currentTask?.title || '暂无任务'}}
          </p>
          
          <button
            onClick={{handleClick}}
            className="
              bg-primary-400 hover:bg-primary-500
              text-white font-bold text-xl
              px-8 py-4 rounded-2xl
              shadow-lg transform hover:scale-105
              transition-all duration-200
            "
          >
            开始学习
          </button>
        </div>
      </div>
    </div>
  )
}}
"""
    
    return page_code


def create_store_template(name, slices):
    """创建 Zustand Store 模板"""
    slice_list = [s.strip() for s in slices.split(',')] if slices else ['default']
    
    slices_code = ""
    for slice_name in slice_list:
        slice_name_camel = slice_name
        slice_name_pascal = slice_name.title().replace('_', '')
        
        slices_code += f"""
// {slice_name_pascal} Slice
export const create{slice_name_pascal}Slice = (set, get) => ({{
  // State
  {slice_name}: null,
  
  // Actions
  set{slice_name_pascal}: (data) => set((state) => ({{ ...state, {slice_name}: data }})),
  reset{slice_name_pascal}: () => set((state) => ({{ ...state, {slice_name}: null }})),
}}))
"""
    
    store_code = f"""// {name} Store
// 使用 Zustand 进行全局状态管理
// 遵循五层架构：业务逻辑层

import {{ create }} from 'zustand'
import {{ subscribeWithSelector }} from 'zustand/middleware'
"""
    
    for slice_name in slice_list:
        slice_name_pascal = slice_name.title().replace('_', '')
        store_code += f"import {{ create{slice_name_pascal}Slice }} from './slices/{slice_name}.js'\n"
    
    store_code += f"""

const {name}Store = create(
  subscribeWithSelector((...args) => ({{
    // 合并所有切片
    ...createDefaultSlice(...args),
"""
    
    for slice_name in slice_list:
        slice_name_pascal = slice_name.title().replace('_', '')
        store_code += f"    ...create{slice_name_pascal}Slice(...args),\n"
    
    store_code += """  }))
)

// 选择器
export const use{name_title}Store = () => {name}Store((state) => state)
export const select{name_title} = (state) => state.{first_slice}

export default {name}Store
""".format(name=name, name_title=name.title(), first_slice=slice_list[0])
    
    return store_code


def main():
    """主函数"""
    args = parse_args()
    
    if args.command == 'create':
        # 创建项目结构
        base_path = create_project_structure(args.name, args.template, args.output)
        
        # 生成配置文件
        package_json = generate_package_json(args.name, args.template)
        with open(base_path / 'package.json', 'w', encoding='utf-8') as f:
            json.dump(package_json, f, indent=2, ensure_ascii=False)
        print("✓ 生成 package.json")
        
        vite_config = generate_vite_config(args.name)
        with open(base_path / 'vite.config.js', 'w', encoding='utf-8') as f:
            f.write(vite_config)
        print("✓ 生成 vite.config.js")
        
        tailwind_config = generate_tailwind_config()
        with open(base_path / 'tailwind.config.js', 'w', encoding='utf-8') as f:
            f.write(tailwind_config)
        print("✓ 生成 tailwind.config.js")
        
        # 生成入口文件
        index_html = generate_index_html(args.name)
        with open(base_path / 'index.html', 'w', encoding='utf-8') as f:
            f.write(index_html)
        print("✓ 生成 index.html")
        
        main_jsx = generate_main_jsx()
        with open(base_path / 'src' / 'main.jsx', 'w', encoding='utf-8') as f:
            f.write(main_jsx)
        print("✓ 生成 main.jsx")
        
        app_jsx = generate_app_jsx()
        with open(base_path / 'src' / 'App.jsx', 'w', encoding='utf-8') as f:
            f.write(app_jsx)
        print("✓ 生成 App.jsx")
        
        index_css = generate_index_css()
        with open(base_path / 'src' / 'presentation' / 'styles' / 'index.css', 'w', encoding='utf-8') as f:
            f.write(index_css)
        print("✓ 生成 index.css")
        
        # 创建示例页面
        sample_page = create_sample_page('presentation', 'HomePage', 'MainLayout')
        with open(base_path / 'src' / 'presentation' / 'pages' / 'HomePage.jsx', 'w', encoding='utf-8') as f:
            f.write(sample_page)
        print("✓ 生成 HomePage.jsx")
        
        # 创建 README
        readme = f"""# {args.name}

Generated by Frontend Scaffold - React 五层架构脚手架

## 快速开始

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

## 项目结构

本项目遵循 **五层架构**：

- `src/presentation` - 表现层（UI 组件、页面）
- `src/interaction` - 交互层（Three.js、交互 Hooks）
- `src/business` - 业务逻辑层（Store、Models、Services）
- `src/data` - 数据层（API、Repositories）
- `src/shared` - 共享层（Utils、Constants、Types）

## 技术栈

- React 18
- Vite
- Tailwind CSS
- Zustand
- React Router
- Framer Motion
- Three.js (可选)

祝你开发顺利！🚀
"""
        with open(base_path / 'README.md', 'w', encoding='utf-8') as f:
            f.write(readme)
        print("✓ 生成 README.md")
        
        print(f"\n✅ 项目创建完成！\n")
        print(f"下一步:")
        print(f"  cd {args.name}")
        print(f"  pnpm install")
        print(f"  pnpm dev\n")
    
    elif args.command == 'add':
        if args.action == 'page':
            print(f"📄 添加页面：{args.name}")
            # TODO: 实现页面生成逻辑
        elif args.action == 'component':
            print(f"🧩 添加组件：{args.name}")
            # TODO: 实现组件生成逻辑
        elif args.action == 'scene':
            print(f"🎨 添加 3D 场景：{args.name}")
            # TODO: 实现场景生成逻辑
        elif args.action == 'store':
            print(f"🗄️ 添加 Store: {args.name}")
            # TODO: 实现 Store 生成逻辑
        else:
            print("❌ 未知的添加类型")
            sys.exit(1)
    
    else:
        print("❌ 错误：未知命令。使用 --help 查看可用命令")
        sys.exit(1)


if __name__ == '__main__':
    main()
