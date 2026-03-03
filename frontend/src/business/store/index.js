/**
 * Redux Store 配置
 * 
 * 使用 Redux Toolkit 创建和配置全局 Redux store
 */

import { configureStore } from '@reduxjs/toolkit';

// 导入所有切片
import authSlice from './slices/authSlice.js';
import uiSlice from './slices/uiSlice.js';
import userSlice from './slices/userSlice.js';
import learningSlice from './slices/learningSlice.js';

/**
 * 创建 Redux store
 */
const store = configureStore({
  reducer: {
    auth: authSlice,
    ui: uiSlice,
    user: userSlice,
    learning: learningSlice
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // 忽略某些不可序列化的值
        ignoredActions: [],
        ignoredPaths: []
      }
    }).concat([
      // 可以添加自定义中间件
    ]),
  devTools: process.env.NODE_ENV !== 'production'
});

export default store;
