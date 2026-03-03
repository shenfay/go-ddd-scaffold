/**
 * 用户信息状态 Slice
 * 
 * 管理用户信息：个人资料、偏好设置、学习进度等
 */

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import userService from '../../../data/api/services/userService.js';

/**
 * 异步 Thunk: 获取用户信息
 */
export const fetchUserProfile = createAsyncThunk(
  'user/fetchProfile',
  async (_, { rejectWithValue }) => {
    try {
      const response = await userService.getProfile();
      return response.data;
    } catch (error) {
      return rejectWithValue(error.message);
    }
  }
);

/**
 * 异步 Thunk: 更新用户信息
 */
export const updateUserProfile = createAsyncThunk(
  'user/updateProfile',
  async (userData, { rejectWithValue }) => {
    try {
      const response = await userService.updateProfile(userData);
      return response.data;
    } catch (error) {
      return rejectWithValue(error.message);
    }
  }
);

const userSlice = createSlice({
  name: 'user',
  initialState: {
    profile: null,
    isLoading: false,
    error: null,
    preferences: {
      language: localStorage.getItem('language') || 'zh-CN',
      soundEnabled: localStorage.getItem('soundEnabled') !== 'false',
      notificationsEnabled: localStorage.getItem('notificationsEnabled') !== 'false'
    }
  },
  
  reducers: {
    // 设置语言偏好
    setLanguage: (state, action) => {
      state.preferences.language = action.payload;
      localStorage.setItem('language', action.payload);
    },
    
    // 切换声音开关
    toggleSound: (state) => {
      state.preferences.soundEnabled = !state.preferences.soundEnabled;
      localStorage.setItem('soundEnabled', state.preferences.soundEnabled);
    },
    
    // 切换通知开关
    toggleNotifications: (state) => {
      state.preferences.notificationsEnabled = !state.preferences.notificationsEnabled;
      localStorage.setItem('notificationsEnabled', state.preferences.notificationsEnabled);
    },
    
    // 清除错误
    clearError: (state) => {
      state.error = null;
    }
  },
  
  extraReducers: (builder) => {
    // 处理获取用户信息
    builder
      .addCase(fetchUserProfile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.profile = action.payload;
      })
      .addCase(fetchUserProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload;
      })
      
      // 处理更新用户信息
      .addCase(updateUserProfile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(updateUserProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.profile = action.payload;
      })
      .addCase(updateUserProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload;
      });
  }
});

export const { setLanguage, toggleSound, toggleNotifications, clearError } = userSlice.actions;
export default userSlice.reducer;
