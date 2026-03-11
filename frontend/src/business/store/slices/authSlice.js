/**
 * 认证状态 Slice
 * 
 * 管理用户认证状态：登录、登出、令牌管理、租户选择
 */

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import userService from '../../../data/api/services/userService.js';

/**
 * 异步 Thunk: 用户登录
 */
export const loginUser= createAsyncThunk(
  'auth/loginUser',
  async ({ email, password }, { rejectWithValue }) => {
  try {
     // 登录前强制清除旧 Token（确保不使用已失效的 Token）
  const oldToken = localStorage.getItem('auth_token');
    if (oldToken) {
    console.log('[Login] 发现旧 Token，强制清除:', oldToken.substring(0, 30) + '...');
    console.log('[Login] 旧 Token 完整内容:', oldToken);
    localStorage.removeItem('auth_token');
    }
    
   console.log('[Login] 开始调用登录 API...');
   const response = await userService.login(email, password);
      // 登录成功，保存新 token
      // 注意：response 已经是格式化后的数据，直接包含 accessToken
  console.log('[Login] 登录成功，获取到响应:', response);
  console.log('[Login] 新 Token 内容:', response.accessToken);
    
   const token = response.accessToken;
      
      if (token) {
     localStorage.setItem('auth_token', token);
     console.log('[Login] 新 Token 已保存到 localStorage');
     console.log('[Login] localStorage 当前内容:', localStorage.getItem('auth_token'));
      } else {
     console.error('[Login] 错误：响应中没有 accessToken!');
     console.error('[Login] 完整响应:', JSON.stringify(response, null, 2));
      }
      
      return response;
    } catch(error) {
  console.error('[Login] 登录失败:', error.message);
  console.error('[Login] 完整错误:', error);
      return rejectWithValue(error.message);
    }
  }
);

/**
 * 异步 Thunk: 用户登出
 */
export const logoutUser= createAsyncThunk(
  'auth/logoutUser',
  async (_, { rejectWithValue }) => {
  try {
     // 尝试调用后端登出接口（将 Token 加入黑名单）
    await userService.logout();
    } catch (error) {
     // 即使后端接口失败，也要清除本地存储
   console.warn('登出接口调用失败，但会清除本地 Token:', error.message);
    }
    
    // 强制清除所有认证和租户信息
  localStorage.removeItem('auth_token');
  localStorage.removeItem('current_tenant_id');
  localStorage.removeItem('user_tenants');
   
   return null;
  }
);

/**
 * 认证 Slice
 */
const authSlice = createSlice({
  name: 'auth',
  initialState: {
    token: null,
    currentTenantId: null,
    isLoading: false,
    isAuthenticated: false,
    error: null
  },
  reducers: {
    // 同步 Action: 从 localStorage 恢复认证状态
    restoreAuthState: (state) => {
      const token = localStorage.getItem('auth_token');
      const tenantId = localStorage.getItem('current_tenant_id');
      state.token = token;
      state.currentTenantId = tenantId;
      state.isAuthenticated = !!token;
    },
    
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
    
    // 同步 Action: 设置当前租户
    setCurrentTenant: (state, action) => {
      state.currentTenantId = action.payload;
      if (action.payload) {
        localStorage.setItem('current_tenant_id', action.payload);
      } else {
        localStorage.removeItem('current_tenant_id');
      }
    },
    
    // 同步 Action: 清除令牌和租户
    clearToken: (state) => {
      state.token = null;
      state.currentTenantId = null;
      state.isAuthenticated = false;
      localStorage.removeItem('auth_token');
      localStorage.removeItem('current_tenant_id');
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
        state.token = action.payload?.accessToken;
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

export const { clearError, setToken, clearToken, restoreAuthState, setCurrentTenant } = authSlice.actions;
export default authSlice.reducer;
