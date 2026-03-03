/**
 * 请求拦截器
 * 
 * 处理所有 API 请求的通用逻辑：
 * - 添加认证令牌
 * - 添加通用请求头
 * - 记录请求日志
 */

import httpClient from '../client.js';

/**
 * 认证拦截器
 * 为每个请求添加认证令牌（如果存在）
 */
export function authInterceptor(config) {
  const token = localStorage.getItem('auth_token');
  
  if (token) {
    config.headers = config.headers || {};
    config.headers['Authorization'] = `Bearer ${token}`;
  }

  return config;
}

/**
 * 通用请求头拦截器
 * 添加通用的请求头信息
 */
export function commonHeaderInterceptor(config) {
  config.headers = config.headers || {};
  
  // 添加客户端版本
  config.headers['X-Client-Version'] = process.env.REACT_APP_VERSION || '1.0.0';
  
  // 添加客户端类型
  config.headers['X-Client-Type'] = 'web';
  
  // 添加请求 ID（用于追踪）
  config.headers['X-Request-ID'] = generateRequestId();
  
  // 添加时间戳
  config.headers['X-Timestamp'] = new Date().toISOString();

  return config;
}

/**
 * 请求日志拦截器
 * 记录所有 API 请求
 */
export function requestLoggingInterceptor(config) {
  const startTime = Date.now();
  
  // 在 config 中存储开始时间，供响应拦截器使用
  config._startTime = startTime;
  
  console.log('[API Request]', {
    method: config.method,
    url: config.url,
    headers: config.headers,
    body: config.body ? JSON.parse(config.body) : null
  });

  return config;
}

/**
 * 生成唯一的请求 ID
 * @private
 */
function generateRequestId() {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

/**
 * 初始化请求拦截器
 * 在应用启动时调用此函数
 */
export function initializeRequestInterceptors() {
  httpClient.addRequestInterceptor(authInterceptor);
  httpClient.addRequestInterceptor(commonHeaderInterceptor);
  httpClient.addRequestInterceptor(requestLoggingInterceptor);
}

export default {
  authInterceptor,
  commonHeaderInterceptor,
  requestLoggingInterceptor,
  initializeRequestInterceptors
};
