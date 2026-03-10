/**
 * 用户服务
 * 封装所有用户相关的 API 调用
 */

import httpClient from '../client.js';
import { getEndpoint } from '../../endpoints/endpoints.js';
import { formatSuccessResponse, formatErrorResponse } from '../responseFormatter.js';

class UserService {
  /**
   * 用户登录
   */
  async login(email, password) {
   const path = getEndpoint('auth.login');
   const response = await httpClient.post(path, { email, password });
   return formatSuccessResponse(response);
  }

  /**
   * 用户注册
   */
  async register(userData) {
   const path = getEndpoint('auth.register');
   const response = await httpClient.post(path, userData);
   return formatSuccessResponse(response);
  }

  /**
   * 获取用户信息
   */
  async getProfile() {
   const path = getEndpoint('user.profile');
   const response = await httpClient.get(path);
   return formatSuccessResponse(response);
  }

  /**
   * 更新用户信息
   */
  async updateProfile(userData) {
   const path= getEndpoint('user.updateProfile');
   const response = await httpClient.put(path, userData);
   return formatSuccessResponse(response);
  }

  /**
   * 修改密码
   */
  async changePassword(oldPassword, newPassword) {
   const path = getEndpoint('user.changePassword');
   const response = await httpClient.post(path, { oldPassword, newPassword });
   return formatSuccessResponse(response);
  }

  /**
   * 用户登出（调用 /api/auth/logout，带 Token 黑名单机制）
   */
  async logout() {
   const path = getEndpoint('auth.logout');
   const response = await httpClient.post(path);
   return formatSuccessResponse(response);
  }

  /**
   * 获取用户基本信息
   */
  async getInfo() {
   const path= getEndpoint('user.getInfo');
   const response = await httpClient.get(path);
   return formatSuccessResponse(response);
  }

  /**
   * 获取指定用户详情
   */
  async getUser(userId) {
   const path = getEndpoint('user.getUser', { id: userId });
   const response = await httpClient.get(path);
   return formatSuccessResponse(response);
  }

  /**
   * 更新指定用户信息
   */
  async updateUser(userId, userData) {
   const path = getEndpoint('user.updateUser', { id: userId });
   const response = await httpClient.put(path, userData);
   return formatSuccessResponse(response);
  }

  /**
   * 获取用户的租户列表
   */
  async getUserTenants() {
   const path = getEndpoint('tenant.userTenants');
   const response = await httpClient.get(path);
   return formatSuccessResponse(response);
  }

  /**
   * 创建租户
   */
  async createTenant(tenantData) {
   const path = getEndpoint('tenant.create');
   const response = await httpClient.post(path, tenantData);
   return formatSuccessResponse(response);
  }

  /**
   * 选择当前使用的租户
   * @param {string} tenantId - 租户 ID
   */
  selectTenant(tenantId) {
    localStorage.setItem('current_tenant_id', tenantId);
    // 触发事件通知其他组件
    window.dispatchEvent(new CustomEvent('tenantChanged', { 
      detail: { tenantId } 
    }));
  }

  /**
   * 获取当前选择的租户 ID
   */
  getCurrentTenantId() {
    return localStorage.getItem('current_tenant_id');
  }
}

// 创建全局用户服务实例
const userService = new UserService();

export default userService;
export { UserService };
