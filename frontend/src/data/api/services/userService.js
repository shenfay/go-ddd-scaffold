/**
 * 用户服务
 * 封装所有用户相关的 API 调用
 */

import httpClient from '../client.js';
import { getEndpoint } from '../../endpoints/endpoints.js';

class UserService {
  /**
   * 用户登录
   */
  async login(email, password) {
    const path = getEndpoint('user.login');
    return httpClient.post(path, { email, password });
  }

  /**
   * 用户注册
   */
  async register(userData) {
    const path = getEndpoint('user.register');
    return httpClient.post(path, userData);
  }

  /**
   * 获取用户信息
   */
  async getProfile() {
    const path = getEndpoint('user.profile');
    return httpClient.get(path);
  }

  /**
   * 更新用户信息
   */
  async updateProfile(userData) {
    const path = getEndpoint('user.updateProfile');
    return httpClient.put(path, userData);
  }

  /**
   * 修改密码
   */
  async changePassword(oldPassword, newPassword) {
    const path = getEndpoint('user.changePassword');
    return httpClient.post(path, { oldPassword, newPassword });
  }

  /**
   * 用户登出
   */
  async logout() {
    const path = getEndpoint('user.logout');
    return httpClient.post(path);
  }

  /**
   * 获取用户基本信息
   */
  async getInfo() {
    const path = getEndpoint('user.getInfo');
    return httpClient.get(path);
  }

  /**
   * 获取用户的租户列表
   */
  async getUserTenants() {
    const path = getEndpoint('tenant.userTenants');
    return httpClient.get(path);
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
