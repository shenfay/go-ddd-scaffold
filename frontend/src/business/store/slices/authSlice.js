/**
 * 认证状态 Slice
 * 
 * 管理用户认证状态：登录、登出、令牌管理
 */

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import userService from '../../../data/api/services/userService.js';

/**
 * 异步 Thunk: 用户登录
 */
export const loginUser = createAsyncThunk(
  'auth/loginUser',
  async ({ email, password }, { rejectWithValue }) => {
    try {
      const response = await userService.login(email, password);
      const token = response.data?.token;
      
      if (token) {
        localStorage.setItem('auth_token', token);
      }
      
      return response.data;
    } catch (error) {
      return rejectWithValue(error.message);
    }
  }
);

/**
 * 异步 Thunk: 用户登出
 */
export const logoutUser = createAsyncThunk(
  'auth/logoutUser',
  async (_, { rejectWithValue }) => {
    try {
      await userService.logout();
      localStorage.removeItem('auth_token');
      return null;
    } catch (error) {
      return rejectWithValue(error.message);
    }
  }
);

/**
 * 认证 Slice
 */
const authSlice = createSlice({
  name: 'auth',
  initialState: {
    token: localStorage.getItem('auth_token') || null,
    isLoading: false,
    isAuthenticated: !!localStorage.getItem('auth_token'),
    error: null
  },
  reducers: {
    // 同步 Action: 清除错误
    clearError: (state) => {
      state.error = null;
    },
    
    // 同步 Action: 设置令牌
    setToken: (state, action) => {
      state.token = action.payload;
      state.isAuthenticated = !!action.payload;
      if (action.payload) {
        localStorage.setItem('auth_token', action.payload);
      }
    },
    
    // 同步 Action: 清除令牌
    clearToken: (state) => {
      state.token = null;
      state.isAuthenticated = false;
      localStorage.removeItem('auth_token');
    }
  },
  extraReducers: (builder) => {
    // 处理登录
    builder
      .addCase(loginUser.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(loginUser.fulfilled, (state, action) => {
        state.isLoading = false;
        state.isAuthenticated = true;
        state.token = action.payload?.token;
      })
      .addCase(loginUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload;
        state.isAuthenticated = false;
      })
      
      // 处理登出
      .addCase(logoutUser.pending, (state) => {
        state.isLoading = true;
      })
      .addCase(logoutUser.fulfilled, (state) => {
        state.isLoading = false;
        state.token = null;
        state.isAuthenticated = false;
        state.error = null;
      })
      .addCase(logoutUser.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload;
      });
  }
});

export const { clearError, setToken, clearToken } = authSlice.actions;
export default authSlice.reducer;
