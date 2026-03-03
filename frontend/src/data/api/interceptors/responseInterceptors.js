/**
 * 响应拦截器
 * 
 * 处理所有 API 响应的通用逻辑：
 * - 处理成功响应
 * - 处理错误响应
 * - 记录响应日志
 * - 处理认证过期
 */

import httpClient from '../client.js';

/**
 * 响应日志拦截器
 * 记录所有 API 响应
 */
export function responseLoggingInterceptor(response) {
  const duration = Date.now() - (response._startTime || 0);
  
  console.log('[API Response]', {
    status: response.status,
    statusText: response.statusText,
    duration: `${duration}ms`,
    data: response.data
  });

  return response;
}

/**
 * 成功响应处理拦截器
 * 提取响应数据，统一响应格式
 */
export function successResponseInterceptor(response) {
  // 如果响应数据具有标准格式 { code, data, message }，则直接返回
  if (response.data && typeof response.data === 'object' && 'code' in response.data) {
    if (response.data.code === 0 || response.data.code === 200) {
      return {
        ...response,
        data: response.data.data || response.data
      };
    }
  }

  return response;
}

/**
 * 错误处理拦截器
 * 处理所有错误情况
 */
import { ERROR_CODES, isErrorType, getErrorMessage } from '../../../shared/constants/errorCodes';

export function errorInterceptor(error) {
  // 处理响应错误
  if (error.response) {
    const status = error.response.status;
    const data = error.response.data;
    
    console.error('[API Error]', {
      status,
      code: data?.code,
      message: data?.message || error.message,
      data
    });

    // 优先使用后端返回的错误码
    const errorCode = data?.code || `HTTP_${status}`;
    
    // 处理 401 未授权 - 认证过期
    if (status === 401 || isErrorType.isAuth(errorCode)) {
      handleUnauthorized();
      const err = new Error(getErrorMessage(ERROR_CODES.UNAUTHORIZED));
      err.code = ERROR_CODES.UNAUTHORIZED;
      err.status = status;
      throw err;
    }

    // 处理 403 禁止
    if (status === 403) {
      const err = new Error(getErrorMessage(ERROR_CODES.FORBIDDEN));
      err.code = ERROR_CODES.FORBIDDEN;
      err.status = status;
      throw err;
    }

    // 处理 404 未找到
    if (status === 404 || isErrorType.isNotFound(errorCode)) {
      const err = new Error(getErrorMessage(ERROR_CODES.NOT_FOUND));
      err.code = ERROR_CODES.NOT_FOUND;
      err.status = status;
      throw err;
    }

    // 处理 429 请求过于频繁
    if (status === 429) {
      const err = new Error(getErrorMessage(ERROR_CODES.TOO_MANY_REQUESTS));
      err.code = ERROR_CODES.TOO_MANY_REQUESTS;
      err.status = status;
      throw err;
    }

    // 处理 500+ 服务器错误
    if (status >= 500 || isErrorType.isServer(errorCode)) {
      const err = new Error(getErrorMessage(ERROR_CODES.SYSTEM_INTERNAL_ERROR));
      err.code = ERROR_CODES.SYSTEM_INTERNAL_ERROR;
      err.status = status;
      throw err;
    }

    // 使用后端返回的具体错误信息
    if (data?.code) {
      const err = new Error(data.message || getErrorMessage(data.code));
      err.code = data.code;
      err.status = status;
      err.details = data.error?.details;
      throw err;
    }

    // 其他错误
    const err = new Error(data?.message || `请求失败: HTTP ${status}`);
    err.code = `HTTP_${status}`;
    err.status = status;
    throw err;
  }

  // 处理网络错误
  if (error.message === 'Failed to fetch' || !error.response) {
    const err = new Error('网络连接失败，请检查您的网络设置');
    err.code = 'NETWORK_ERROR';
    throw err;
  }

  // 处理超时错误
  if (error.name === 'AbortError') {
    const err = new Error('请求超时，请稍后重试');
    err.code = 'TIMEOUT';
    throw err;
  }

  throw error;
}

/**
 * 处理未授权情况
 * @private
 */
function handleUnauthorized() {
  // 清除本地存储的认证信息
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_info');

  // 触发全局事件，通知应用需要重新登录
  window.dispatchEvent(new CustomEvent('auth:expired', {
    detail: { message: '认证已过期，请重新登录' }
  }));

  // 可选：自动重定向到登录页面
  // window.location.href = '/login';
}

/**
 * 初始化响应拦截器
 * 在应用启动时调用此函数
 */
export function initializeResponseInterceptors() {
  httpClient.addResponseInterceptor(responseLoggingInterceptor);
  httpClient.addResponseInterceptor(successResponseInterceptor);
  httpClient.addErrorInterceptor(errorInterceptor);
}

export default {
  responseLoggingInterceptor,
  successResponseInterceptor,
  errorInterceptor,
  initializeResponseInterceptors
};
