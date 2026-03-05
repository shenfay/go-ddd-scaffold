/**
 * 响应拦截器
 * 
 * 处理所有 API 响应的通用逻辑：
 * - 处理成功响应（统一错误码检查）
 * - 处理错误响应（友好的错误提示）
 * - 记录响应日志
 * - 处理认证过期
 * - 支持重试机制
 */

import httpClient from '../client.js';
import { errorHandler } from '../../../shared/utils/errorHandler';

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
 * 提取响应数据，统一响应格式，检查业务错误码
 */
export function successResponseInterceptor(response) {
  // 如果响应数据具有标准格式 { code, data, message }
  if (response.data && typeof response.data === 'object' && 'code' in response.data) {
    const backendCode = response.data.code;
    
    // 检查后端返回的业务状态码
    // 成功的标准：code === "Success" 或 code === 0 或 code === "0"
    const isSuccess = backendCode === 'Success' || backendCode === 0 || backendCode === '0';
    
    if (!isSuccess) {
      // 业务错误，使用统一的错误处理器
      const message = response.data.message || '操作失败';
      
      // 显示友好的错误提示
      errorHandler.showError(message);
      
      // 抛出错误，让上层知道请求失败
      const error = new Error(message);
      error.errorCode = backendCode;
      error.response = response;
      error.details = response.data.details;
      throw error;
    }
    
    // 成功响应，提取 data 字段
    return {
      ...response,
      data: response.data.data || response.data
    };
  }

  return response;
}

/**
 * 错误处理拦截器
 * 处理所有错误情况，提供友好的错误提示
 */
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

    // 优先使用后端返回的错误码（如果还没有设置）
    const errorCode = error.errorCode || data?.code || `HTTP_${status}`;
    
    // 处理 401 未授权 - 需要区分是 Token 过期还是业务错误
    if (status === 401) {
      // 检查是否是业务错误（如密码错误、用户不存在等）
      const businessErrorCodes = [
        'User.InvalidPassword',
        'User.Unauthorized',
        'User.NotFound',
      ];
      
      if (businessErrorCodes.includes(errorCode)) {
        // 业务错误，显示具体错误消息
        const message = data?.message || errorHandler.getFriendlyMessage(errorCode);
        errorHandler.showError(message);
        error.message = message;
        error.errorCode = errorCode;
        throw error;
      } else {
        // Token 过期或无效，清除认证信息并提示重新登录
        handleUnauthorized(error);
        error.message = '登录已过期，请重新登录';
        error.errorCode = errorCode;
        throw error;
      }
    }

    // 处理 403 禁止访问
    if (status === 403) {
      errorHandler.showError('抱歉，您没有权限执行此操作');
      error.message = '没有权限执行此操作';
      error.errorCode = errorCode;
      error._handled = true;
      throw error;
    }

    // 处理 404 资源不存在
    if (status === 404) {
      errorHandler.showError('请求的资源不存在');
      error.message = '资源不存在';
      error.errorCode = errorCode;
      error._handled = true;
      throw error;
    }

    // 处理 429 请求过于频繁
    if (status === 429) {
      errorHandler.showError('请求过于频繁，请稍后再试');
      error.message = '请求过于频繁';
      error.errorCode = errorCode;
      error._handled = true;
      throw error;
    }

    // 处理 500+ 服务器错误
    if (status >= 500) {
      errorHandler.showError('系统繁忙，请稍后重试');
      error.message = '系统内部错误';
      error.errorCode = errorCode;
      error._handled = true;
      throw error;
    }

    // 使用后端返回的具体错误信息
    if (data?.code) {
      const message = data.message || errorHandler.getFriendlyMessage(errorCode);
      errorHandler.showError(message);
      error.message = message;
      error.errorCode = data.code;
      error._handled = true;
      throw error;
    }

    // 其他错误
    const message = data?.message || `请求失败：HTTP ${status}`;
    errorHandler.showError(message);
    error.message = message;
    error.errorCode = errorCode;
    error._handled = true;
    throw error;
  }

  // 处理网络错误（没有收到响应）
  if (error.message === 'Failed to fetch' || !error.response) {
    errorHandler.showError('网络连接失败，请检查您的网络设置');
    error.message = '网络连接失败';
    error.errorCode = 'NETWORK_ERROR';
    error._handled = true;
    throw error;
  }

  // 处理超时错误
  if (error.name === 'AbortError' || error.message.includes('timeout')) {
    errorHandler.showError('请求超时，请稍后重试');
    error.message = '请求超时';
    error.errorCode = 'TIMEOUT';
    error._handled = true;
    throw error;
  }

  // 未知错误
  console.error('[Unknown Error]', error);
  errorHandler.showError('发生未知错误，请稍后重试');
  error.message = '发生未知错误';
  error._handled = true;
  throw error;
}

/**
 * 处理未授权情况
 * @private
 * @param {Error} error - 错误对象
 */
function handleUnauthorized(error) {
  // 清除本地存储的认证信息
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_info');
  localStorage.removeItem('current_user_id');

  // 触发全局事件，通知应用需要重新登录
  window.dispatchEvent(new CustomEvent('auth:expired', {
    detail: { message: '登录已过期，请重新登录' }
  }));

  // 可选：自动重定向到登录页面
  // window.location.href = '/login';
  
  console.warn('Token 已过期或无效:', error);
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
