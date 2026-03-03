/**
 * UI 状态 Slice
 * 
 * 管理全局 UI 状态：模态框、侧边栏、通知等
 */

import { createSlice } from '@reduxjs/toolkit';

const uiSlice = createSlice({
  name: 'ui',
  initialState: {
    // 模态框状态
    modals: {
      isOpen: false,
      type: null,
      data: null
    },
    
    // 侧边栏状态
    sidebar: {
      isOpen: false
    },
    
    // 通知状态
    notification: {
      isOpen: false,
      type: 'info', // 'success', 'error', 'warning', 'info'
      message: '',
      duration: 3000
    },
    
    // 加载状态
    isLoading: false,
    
    // 主题
    theme: localStorage.getItem('theme') || 'light'
  },
  
  reducers: {
    // 模态框操作
    openModal: (state, action) => {
      state.modals = {
        isOpen: true,
        type: action.payload?.type || null,
        data: action.payload?.data || null
      };
    },
    closeModal: (state) => {
      state.modals = {
        isOpen: false,
        type: null,
        data: null
      };
    },
    
    // 侧边栏操作
    openSidebar: (state) => {
      state.sidebar.isOpen = true;
    },
    closeSidebar: (state) => {
      state.sidebar.isOpen = false;
    },
    toggleSidebar: (state) => {
      state.sidebar.isOpen = !state.sidebar.isOpen;
    },
    
    // 通知操作
    showNotification: (state, action) => {
      state.notification = {
        isOpen: true,
        type: action.payload?.type || 'info',
        message: action.payload?.message || '',
        duration: action.payload?.duration || 3000
      };
    },
    hideNotification: (state) => {
      state.notification = {
        ...state.notification,
        isOpen: false
      };
    },
    
    // 加载状态操作
    setLoading: (state, action) => {
      state.isLoading = action.payload;
    },
    
    // 主题操作
    setTheme: (state, action) => {
      const theme = action.payload;
      state.theme = theme;
      localStorage.setItem('theme', theme);
    },
    toggleTheme: (state) => {
      const newTheme = state.theme === 'light' ? 'dark' : 'light';
      state.theme = newTheme;
      localStorage.setItem('theme', newTheme);
    }
  }
});

export const {
  openModal,
  closeModal,
  openSidebar,
  closeSidebar,
  toggleSidebar,
  showNotification,
  hideNotification,
  setLoading,
  setTheme,
  toggleTheme
} = uiSlice.actions;

export default uiSlice.reducer;
