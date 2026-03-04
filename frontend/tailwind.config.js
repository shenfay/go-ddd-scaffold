/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  theme: {
    extend: {
      // 品牌色彩定义
      colors: {
        primary: {
          DEFAULT: '#2E7D32', // 森林绿 - 主色调
          light: '#4CAF50',
          dark: '#1B5E20',
        },
        secondary: {
          DEFAULT: '#F9F8F4', // 暖沙色 - 背景色
          light: '#FFFBF7',
          dark: '#ECEBE6',
        },
        danger: {
          DEFAULT: '#E53935', // 警告红
          light: '#EF5350',
          dark: '#C62828',
        },
        text: {
          primary: '#1A1A1A',
          secondary: '#666666',
          muted: '#999999',
        },
      },
      
      // 字体配置 - 苹方字体
      fontFamily: {
        sans: ['PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '-apple-system', 'sans-serif'],
      },
      
      // 字体大小规范
      fontSize: {
        'title': ['20pt', { lineHeight: '1.3', fontWeight: '600' }], // Semibold
        'body': ['15pt', { lineHeight: '1.6', fontWeight: '400' }],   // Regular
        'small': ['13pt', { lineHeight: '1.5', fontWeight: '400' }],
      },
      
      // 间距规范 (基于 16px 页面边距，12px 卡片间距)
      spacing: {
        'page': '16px',    // 页面边距
        'card': '12px',    // 卡片间距
        'section': '24px', // 区块间距
      },
      
      // 圆角规范 (4px 基础圆角)
      borderRadius: {
        'sm': '4px',       // 小圆角
        'md': '8px',       // 中等圆角
        'lg': '12px',      // 大圆角
        'xl': '16px',      // 超大圆角
      },
      
      // 线条粗细 (用于图标等)
      borderWidth: {
        '1.5': '1.5px',    // 图标线条粗细
      },
    },
  },
  plugins: [],
}