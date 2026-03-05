/**
 * 统一错误处理工具类
 * 
 * 功能：
 * 1. 统一的错误码映射
 * 2. 友好的错误提示
 * 3. 网络异常处理
 * 4. 重试机制
 */

// ============================================
// 错误码定义（与后端保持一致）
// ============================================

export const ERROR_CODES = {
  // 通用错误
  Success: 'Success',
  InvalidParameter: 'InvalidParameter',
  MissingParameter: 'MissingParameter',
  Unauthorized: 'Unauthorized',
  Forbidden: 'Forbidden',
  NotFound: 'NotFound',
  MethodNotAllowed: 'MethodNotAllowed',
  TooManyRequests: 'TooManyRequests',
  ValidationFailed: 'ValidationFailed',
  ResourceConflict: 'ResourceConflict',
  
  // 用户模块错误
  UserExists: 'User.Exists',
  UserNotFound: 'User.NotFound',
  InvalidPassword: 'User.InvalidPassword',
  UserUnauthorized: 'User.Unauthorized',
  TenantLimitExceed: 'User.TenantLimitExceed',
  InvalidEmail: 'User.InvalidEmail',
  InvalidRole: 'User.InvalidRole',
  
  // 系统错误
  InternalError: 'System.InternalError',
  DatabaseError: 'System.DatabaseError',
  CacheUnavailable: 'System.CacheUnavailable',
  ExternalServiceError: 'System.ExternalServiceError',
  Timeout: 'System.Timeout',
  
  // 网络错误
  NetworkTimeout: 'Network.Timeout',
  NetworkUnavailable: 'Network.Unavailable',
  ConnectionRefused: 'Network.ConnectionRefused',
  DNSResolutionFailed: 'Network.DNSResolutionFailed',
  SSLHandshakeFailed: 'Network.SSLHandshakeFailed',
  RequestEntityTooLarge: 'Network.RequestEntityTooLarge',
  ServiceUnavailable: 'Network.ServiceUnavailable',
};

// ============================================
// HTTP 状态码映射
// ============================================

const HTTP_STATUS_MAP = {
  // 通用错误映射
  [ERROR_CODES.Success]: 200,
  [ERROR_CODES.InvalidParameter]: 400,
  [ERROR_CODES.MissingParameter]: 400,
  [ERROR_CODES.Unauthorized]: 401,
  [ERROR_CODES.Forbidden]: 403,
  [ERROR_CODES.NotFound]: 404,
  [ERROR_CODES.MethodNotAllowed]: 405,
  [ERROR_CODES.TooManyRequests]: 429,
  [ERROR_CODES.ValidationFailed]: 400,
  [ERROR_CODES.ResourceConflict]: 409,
  
  // 用户模块映射
  [ERROR_CODES.UserExists]: 409,
  [ERROR_CODES.UserNotFound]: 404,
  [ERROR_CODES.InvalidPassword]: 401,
  [ERROR_CODES.UserUnauthorized]: 401,
  [ERROR_CODES.TenantLimitExceed]: 403,
  [ERROR_CODES.InvalidEmail]: 400,
  [ERROR_CODES.InvalidRole]: 400,
  
  // 系统错误映射
  [ERROR_CODES.InternalError]: 500,
  [ERROR_CODES.DatabaseError]: 503,
  [ERROR_CODES.CacheUnavailable]: 503,
  [ERROR_CODES.ExternalServiceError]: 502,
  [ERROR_CODES.Timeout]: 504,
  
  // 网络错误映射
  [ERROR_CODES.NetworkTimeout]: 504,
  [ERROR_CODES.NetworkUnavailable]: 503,
  [ERROR_CODES.ConnectionRefused]: 503,
  [ERROR_CODES.DNSResolutionFailed]: 503,
  [ERROR_CODES.SSLHandshakeFailed]: 502,
  [ERROR_CODES.RequestEntityTooLarge]: 413,
  [ERROR_CODES.ServiceUnavailable]: 503,
};

// ============================================
// 友好的错误提示消息
// ============================================

const FRIENDLY_MESSAGES = {
  // 通用错误
  [ERROR_CODES.Success]: '操作成功',
  [ERROR_CODES.InvalidParameter]: '请求参数错误，请检查输入',
  [ERROR_CODES.MissingParameter]: '缺少必要参数',
  [ERROR_CODES.Unauthorized]: '登录已过期，请重新登录',
  [ERROR_CODES.Forbidden]: '抱歉，您没有权限执行此操作',
  [ERROR_CODES.NotFound]: '请求的资源不存在',
  [ERROR_CODES.MethodNotAllowed]: '不支持的请求方法',
  [ERROR_CODES.TooManyRequests]: '请求过于频繁，请稍后再试',
  [ERROR_CODES.ValidationFailed]: '数据验证失败，请检查输入',
  [ERROR_CODES.ResourceConflict]: '资源冲突，请刷新后重试',
  
  // 用户模块
  [ERROR_CODES.UserExists]: '该用户已存在，请使用其他账号',
  [ERROR_CODES.UserNotFound]: '用户不存在',
  [ERROR_CODES.InvalidPassword]: '密码错误，请重试',
  [ERROR_CODES.UserUnauthorized]: '登录已过期或账户已被禁用',
  [ERROR_CODES.TenantLimitExceed]: '租户用户数已达上限',
  [ERROR_CODES.InvalidEmail]: '邮箱格式不正确',
  [ERROR_CODES.InvalidRole]: '无效的用户角色',
  
  // 系统错误
  [ERROR_CODES.InternalError]: '系统繁忙，请稍后重试',
  [ERROR_CODES.DatabaseError]: '数据库服务异常，请稍后重试',
  [ERROR_CODES.CacheUnavailable]: '缓存服务异常，请稍后重试',
  [ERROR_CODES.ExternalServiceError]: '外部服务调用失败',
  [ERROR_CODES.Timeout]: '请求超时，请稍后重试',
  
  // 网络错误
  [ERROR_CODES.NetworkTimeout]: '网络连接超时，请检查网络后重试',
  [ERROR_CODES.NetworkUnavailable]: '网络不可用，请检查网络连接',
  [ERROR_CODES.ConnectionRefused]: '连接被拒绝，请稍后重试',
  [ERROR_CODES.DNSResolutionFailed]: 'DNS 解析失败，请检查网络设置',
  [ERROR_CODES.SSLHandshakeFailed]: '安全连接建立失败',
  [ERROR_CODES.RequestEntityTooLarge]: '上传的文件过大',
  [ERROR_CODES.ServiceUnavailable]: '服务暂时不可用，请稍后重试',
};

// ============================================
// 错误处理工具类
// ============================================

class ErrorHandler {
  /**
   * 获取 HTTP 状态码
   * @param {string} errorCode - 错误码
   * @returns {number} HTTP 状态码
   */
  getHTTPStatus(errorCode) {
    return HTTP_STATUS_MAP[errorCode] || 500;
  }

  /**
   * 获取友好的错误提示消息
   * @param {string} errorCode - 错误码
   * @param {string} customMessage - 自定义消息（可选）
   * @returns {string} 友好的错误提示
   */
  getFriendlyMessage(errorCode, customMessage = '') {
    if (customMessage) {
      return customMessage;
    }
    
    return FRIENDLY_MESSAGES[errorCode] || '系统繁忙，请稍后重试';
  }

  /**
   * 判断是否是客户端错误（4xx）
   * @param {string} errorCode - 错误码
   * @returns {boolean}
   */
  isClientError(errorCode) {
    const status = this.getHTTPStatus(errorCode);
    return status >= 400 && status < 500;
  }

  /**
   * 判断是否是服务端错误（5xx）
   * @param {string} errorCode - 错误码
   * @returns {boolean}
   */
  isServerError(errorCode) {
    const status = this.getHTTPStatus(errorCode);
    return status >= 500;
  }

  /**
   * 判断是否是网络错误
   * @param {string} errorCode - 错误码
   * @returns {boolean}
   */
  isNetworkError(errorCode) {
    return errorCode.startsWith('Network.');
  }

  /**
   * 判断是否应该重试
   * @param {string} errorCode - 错误码
   * @returns {boolean}
   */
  shouldRetry(errorCode) {
    const retryCodes = [
      ERROR_CODES.NetworkTimeout,
      ERROR_CODES.NetworkUnavailable,
      ERROR_CODES.ConnectionRefused,
      ERROR_CODES.DNSResolutionFailed,
      ERROR_CODES.Timeout,
      ERROR_CODES.DatabaseError,
      ERROR_CODES.CacheUnavailable,
      ERROR_CODES.ExternalServiceError,
      ERROR_CODES.ServiceUnavailable,
    ];
    
    return retryCodes.includes(errorCode);
  }

  /**
   * 处理 HTTP 响应错误
   * @param {Response} response - Fetch API 响应对象
   * @returns {Promise<Object>} 错误信息对象
   */
  async handleHTTPResponse(response) {
    const statusCode = response.status;
    
    try {
      const data = await response.json();
      
      return {
        success: false,
        errorCode: data.code || this.mapStatusCodeToErrorCode(statusCode),
        httpStatus: statusCode,
        message: data.message || this.getFriendlyMessage(
          data.code || this.mapStatusCodeToErrorCode(statusCode)
        ),
        details: data.details,
      };
    } catch (e) {
      // 无法解析 JSON，使用默认错误
      return {
        success: false,
        errorCode: this.mapStatusCodeToErrorCode(statusCode),
        httpStatus: statusCode,
        message: this.getFriendlyMessage(
          this.mapStatusCodeToErrorCode(statusCode)
        ),
      };
    }
  }

  /**
   * 将 HTTP 状态码映射到错误码
   * @param {number} statusCode - HTTP 状态码
   * @returns {string} 错误码
   */
  mapStatusCodeToErrorCode(statusCode) {
    switch (statusCode) {
      case 400:
        return ERROR_CODES.InvalidParameter;
      case 401:
        return ERROR_CODES.Unauthorized;
      case 403:
        return ERROR_CODES.Forbidden;
      case 404:
        return ERROR_CODES.NotFound;
      case 405:
        return ERROR_CODES.MethodNotAllowed;
      case 409:
        return ERROR_CODES.ResourceConflict;
      case 413:
        return ERROR_CODES.RequestEntityTooLarge;
      case 429:
        return ERROR_CODES.TooManyRequests;
      case 500:
        return ERROR_CODES.InternalError;
      case 502:
        return ERROR_CODES.ExternalServiceError;
      case 503:
        return ERROR_CODES.ServiceUnavailable;
      case 504:
        return ERROR_CODES.Timeout;
      default:
        return ERROR_CODES.InternalError;
    }
  }

  /**
   * 显示错误提示（支持多种 UI 框架）
   * @param {string} message - 错误消息
   * @param {Object} options - 选项
   */
  showError(message, options = {}) {
    const { duration = 3000, type = 'error' } = options;
    
    // 如果使用了 Ant Design
    if (window.antd && window.antd.message) {
      window.antd.message.error(message, duration);
      return;
    }
    
    // 如果使用了 Element UI
    if (window.ElementUI && window.ElementUI.Message) {
      window.ElementUI.Message.error({ message, duration });
      return;
    }
    
    // 默认：使用浏览器 alert（不推荐，仅作为降级方案）
    console.error(`[${type}] ${message}`);
    // alert(message); // 仅在调试时启用
  }

  /**
   * 处理请求错误并显示提示
   * @param {Error} error - 错误对象
   * @param {Object} options - 选项
   */
  handleRequestError(error, options = {}) {
    const { showError = true, onRetry = null } = options;
    
    let errorCode = ERROR_CODES.InternalError;
    let message = '请求失败，请稍后重试';
    
    // 处理不同类型的错误
    if (error.name === 'TypeError') {
      // 网络错误
      if (error.message.includes('timeout')) {
        errorCode = ERROR_CODES.NetworkTimeout;
        message = '请求超时，请检查网络连接';
      } else if (error.message.includes('Failed to fetch')) {
        errorCode = ERROR_CODES.NetworkUnavailable;
        message = '网络不可用，请检查网络连接';
      } else {
        errorCode = ERROR_CODES.ConnectionRefused;
        message = '连接失败，请检查网络或服务状态';
      }
    } else if (error.response) {
      // HTTP 响应错误
      errorCode = this.mapStatusCodeToErrorCode(error.response.status);
      message = error.response.data?.message || this.getFriendlyMessage(errorCode);
    }
    
    // 显示错误提示
    if (showError) {
      this.showError(message);
    }
    
    // 如果需要重试
    if (onRetry && this.shouldRetry(errorCode)) {
      console.log('建议重试，错误码:', errorCode);
    }
    
    // 返回标准化错误信息
    return {
      success: false,
      errorCode,
      message,
      originalError: error,
    };
  }

  /**
   * 带重试的请求包装器
   * @param {Function} requestFn - 请求函数
   * @param {Object} options - 选项
   * @returns {Promise<any>}
   */
  async withRetry(requestFn, options = {}) {
    const { maxRetries = 3, delay = 1000, onError = null } = options;
    
    let lastError = null;
    
    for (let attempt = 0; attempt <= maxRetries; attempt++) {
      try {
        return await requestFn();
      } catch (error) {
        lastError = error;
        
        // 分析错误是否应该重试
        let shouldRetry = false;
        if (error.response) {
          const errorCode = this.mapStatusCodeToErrorCode(error.response.status);
          shouldRetry = this.shouldRetry(errorCode);
        } else if (error.name === 'TypeError') {
          // 网络错误通常可以重试
          shouldRetry = true;
        }
        
        // 不应该重试或已达到最大次数，直接抛出
        if (!shouldRetry || attempt === maxRetries) {
          break;
        }
        
        // 等待后重试
        console.log(`请求失败，${delay}ms 后重试 (${attempt + 1}/${maxRetries})`);
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
    
    // 所有重试都失败
    if (onError) {
      onError(lastError);
    }
    
    throw lastError;
  }

  /**
   * 批量处理多个请求的错误
   * @param {Array<Promise>} promises - Promise 数组
   * @param {Object} options - 选项
   * @returns {Promise<Array>} 结果数组
   */
  async allSettled(promises, options = {}) {
    const { stopOnError = false } = options;
    
    const results = await Promise.allSettled(promises);
    
    // 如果有错误且要求遇到错误就停止
    if (stopOnError) {
      const errorResult = results.find(r => r.status === 'rejected');
      if (errorResult) {
        throw errorResult.reason;
      }
    }
    
    return results;
  }
}

// 导出单例
export const errorHandler = new ErrorHandler();

// 导出工具函数
export const {
  getHTTPStatus,
  getFriendlyMessage,
  isClientError,
  isServerError,
  isNetworkError,
  shouldRetry,
  handleHTTPResponse,
  handleRequestError,
  showError,
  withRetry,
} = errorHandler;

// 默认导出
export default errorHandler;
