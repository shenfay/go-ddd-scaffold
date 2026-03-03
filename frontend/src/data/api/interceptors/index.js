/**
 * 拦截器索引
 * 统一导出所有拦截器和初始化函数
 */

export * from './requestInterceptors.js';
export * from './responseInterceptors.js';

import { initializeRequestInterceptors } from './requestInterceptors.js';
import { initializeResponseInterceptors } from './responseInterceptors.js';

/**
 * 初始化所有拦截器
 * 在应用启动时调用此函数
 */
export function initializeInterceptors() {
  initializeRequestInterceptors();
  initializeResponseInterceptors();
  console.log('✅ API 拦截器已初始化');
}

export default {
  initializeInterceptors
};
