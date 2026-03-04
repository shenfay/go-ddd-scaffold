# 前端设计系统

## 品牌定义

**极简 · 可依赖 · 自然**

---

## 🎨 Color 色彩

### 主色调
```css
--primary: #2E7D32;      /* 森林绿 - 象征自然、成长 */
--primary-light: #4CAF50;
--primary-dark: #1B5E20;
```

### 辅助色
```css
--secondary: #F9F8F4;    /* 暖沙色 - 温暖、可靠背景 */
--secondary-light: #FFFBF7;
--secondary-dark: #ECEBE6;
```

### 功能色
```css
--danger: #E53935;       /* 警告红 - 错误、警示 */
--danger-light: #EF5350;
--danger-dark: #C62828;
```

### 文本色
```css
--text-primary: #1A1A1A;   /* 主文本 */
--text-secondary: #666666; /* 次要文本 */
--text-muted: #999999;     /* 弱化文本 */
```

### 使用示例
```jsx
// Tailwind CSS 类名
<div className="bg-primary text-white">主色调背景</div>
<div className="text-danger">错误文本</div>
<div className="bg-secondary">暖沙色背景</div>
```

---

## 🔤 Type 字体

### 字体栈
```css
font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', -apple-system, sans-serif;
```

### 字号规范
| 类型 | 字号 | 字重 | 行高 | 使用场景 |
|------|------|------|------|----------|
| 标题 | 20pt | Semibold (600) | 1.3 | 页面标题、卡片标题 |
| 正文 | 15pt | Regular (400) | 1.6 | 常规文本内容 |
| 小字 | 13pt | Regular (400) | 1.5 | 注释、说明文本 |

### 标题层级
```css
h1 { font-size: 32pt; }  /* 页面主标题 */
h2 { font-size: 28pt; }  /* 区块标题 */
h3 { font-size: 24pt; }  /* 子区块标题 */
h4 { font-size: 22pt; }  /* 组标题 */
h5 { font-size: 20pt; }  /* 小标题 */
h6 { font-size: 18pt; }  /* 微小标题 */
```

### 使用示例
```jsx
<h5 className="text-title">标题（20pt Semibold）</h5>
<p className="text-body">正文（15pt Regular）</p>
<span className="text-small">小字（13pt Regular）</span>
```

---

## 🔷 Icons 图标

### 设计规范
- **线条粗细**: 1.5pt
- **圆角半径**: 4px
- **风格**: 等比线性缩放，简洁流畅

### 使用示例
```jsx
<svg width="24" height="24" viewBox="0 0 24 24" fill="none" 
     xmlns="http://www.w3.org/2000/svg"
     className="stroke-current">
  <path d="M12 2L2 7L12 12L22 7L12 2Z" 
        strokeWidth="1.5" 
        strokeLinecap="round" 
        strokeLinejoin="round"/>
</svg>
```

### 图标尺寸
| 尺寸 | 大小 | 使用场景 |
|------|------|----------|
| xs | 16x16 | 按钮内图标、标签图标 |
| sm | 20x20 | 列表项图标 |
| md | 24x24 | 导航图标、功能图标 |
| lg | 32x32 | 特性图标、空状态图标 |
| xl | 48x48 | 欢迎页插图图标 |

---

## 📐 Grid 栅格

### 间距系统
```css
--spacing-page: 16px;     /* 页面边距 */
--spacing-card: 12px;     /* 卡片间距 */
--spacing-section: 24px;  /* 区块间距 */
```

### 页面布局
```
┌─────────────────────────────────┐
│ 16px                            │
│ ┌───────────────────────────┐   │
│ │                           │   │
│ │   内容区域                 │ 16px
│ │                           │   │
│ └───────────────────────────┘   │
│ 16px                            │
└─────────────────────────────────┘
```

### 卡片布局
```
┌─────────────┐ 12px ┌─────────────┐
│             │      │             │
│   卡片 1     │      │   卡片 2     │
│             │      │             │
└─────────────┘      └─────────────┘
       12px
```

### 使用示例
```jsx
// 页面容器
<div className="p-page">页面边距 16px</div>

// 卡片网格
<div className="grid grid-cols-3 gap-card">
  <Card />  {/* 卡片间距 12px */}
  <Card />
  <Card />
</div>

// 区块间距
<section className="mb-section">区块间距 24px</section>
```

---

## 🔘 Components 基础组件

### 按钮
```jsx
// 主按钮
<button className="btn-primary">
  主要操作
</button>

// 次级按钮
<button className="btn-secondary">
  取消
</button>

// 禁用状态
<button className="btn-primary" disabled>
  提交中...
</button>
```

### 输入框
```jsx
<input 
  type="text" 
  className="input"
  placeholder="请输入..."
/>

// 聚焦状态自动应用绿色边框和阴影
```

### 卡片
```jsx
<div className="card">
  <h5>卡片标题</h5>
  <p>卡片内容</p>
</div>
```

### 消息提示
```jsx
<div className="error-message">
  错误信息
</div>

<div className="success-message">
  成功信息
</div>
```

---

## 🎯 设计原则应用

### 极简 (Minimalist)
- ✅ 留白充足，避免拥挤
- ✅ 单一主色调，不过度装饰
- ✅ 清晰的视觉层次
- ✅ 去除不必要的边框和分割线

### 可依赖 (Trustworthy)
- ✅ 一致的交互反馈
- ✅ 清晰的错误提示
- ✅ 稳定的布局结构
- ✅ 明确的视觉引导

### 自然 (Natural)
- ✅ 森林绿主色调
- ✅ 暖沙色温暖背景
- ✅ 圆润的圆角处理
- ✅ 流畅的过渡动画

---

## 📱 响应式断点

```css
/* 手机 */
@media (max-width: 640px) {
  --spacing-page: 12px;
}

/* 平板 */
@media (min-width: 641px) and (max-width: 1024px) {
  /* 单列布局 */
}

/* 桌面 */
@media (min-width: 1025px) {
  /* 多列布局 */
}
```

---

## 🎨 快速参考

### Tailwind CSS 类名速查
```bash
# 颜色
bg-primary      # 森林绿背景
text-primary    # 主文本
border-danger   # 红色边框

# 间距
p-page          # 16px 内边距
m-card          # 12px 外边距
gap-section     # 24px 间距

# 圆角
rounded-sm      # 4px
rounded-md      # 8px
rounded-lg      # 12px

# 字体
text-title      # 20pt Semibold
text-body       # 15pt Regular
font-sans       # 苹方字体
```

---

*最后更新：2026-03-04*
*版本：v1.0.0*
