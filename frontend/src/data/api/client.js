/**
 * HTTP 客户端配置
 * 
 * 统一管理所有 API 请求，包括：
 * - 基础 URL 配置
 * - 请求拦截器
 * - 响应拦截器
 * - 超时设置
 * - 通用错误处理
 * - 重试机制
 */

import { errorHandler } from '../../shared/utils/errorHandler';

class HttpClient {
  constructor() {
    this.baseURL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:3000/api';
    this.timeout = 30000; // 30 seconds
    this.interceptors = {
      request: [],
      response: [],
      error: []
    };
    this.retryConfig = {
      maxRetries: 3,
      retryDelay: 1000, // 1 second
      shouldRetry: true
    };
  }

  /**
   * 设置基础 URL
   * @param {string} url
   */
  setBaseURL(url) {
    this.baseURL = url;
  }

  /**
   * 获取基础 URL
   */
  getBaseURL() {
    return this.baseURL;
  }

  /**
   * 添加请求拦截器
   * @param {Function} handler
   */
  addRequestInterceptor(handler) {
    this.interceptors.request.push(handler);
  }

  /**
   * 添加响应拦截器
   * @param {Function} handler
   */
  addResponseInterceptor(handler) {
    this.interceptors.response.push(handler);
  }

  /**
   * 添加错误拦截器
   * @param {Function} handler
   */
  addErrorInterceptor(handler) {
    this.interceptors.error.push(handler);
  }

  /**
   * 执行请求拦截器
   * @private
   */
  _executeRequestInterceptors(config) {
    let finalConfig = { ...config };
    this.interceptors.request.forEach(handler => {
      finalConfig = handler(finalConfig) || finalConfig;
    });
    return finalConfig;
  }

  /**
   * 执行响应拦截器
   * @private
   */
  _executeResponseInterceptors(response) {
    let finalResponse = { ...response };
    this.interceptors.response.forEach(handler => {
      finalResponse = handler(finalResponse) || finalResponse;
    });
    return finalResponse;
  }

  /**
   * 执行错误拦截器
   * @private
   */
  _executeErrorInterceptors(error) {
    let finalError = error;
    this.interceptors.error.forEach(handler => {
      finalError = handler(finalError) || finalError;
    });
    return finalError;
  }

  /**
   * 通用请求方法（带重试和错误处理）
   * @param {string} method - HTTP 方法 (GET, POST, PUT, DELETE, etc.)
   * @param {string} path - API 路径
   * @param {object} options - 请求选项 (data, params, headers, etc.)
   * @param {boolean} enableRetry - 是否启用重试（默认 true）
   */
  async request(method, path, options = {}, enableRetry = true) {
    const url = `${this.baseURL}${path}`;

    // 构建请求配置
    let config = {
      method,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      timeout: options.timeout || this.timeout
    };

    // 如果有数据，添加到请求体
    if (options.data) {
      config.body = JSON.stringify(options.data);
    }

    // 处理查询参数
    let finalURL = url;
    if (options.params) {
      const queryString = new URLSearchParams(options.params).toString();
      finalURL = `${url}?${queryString}`;
    }

    // 执行请求拦截器
    config = this._executeRequestInterceptors(config);

    try {
      // 发送请求
      const response = await fetch(finalURL, config);

      // 解析响应
      const contentType = response.headers.get('content-type');
      let data;

      if (contentType && contentType.includes('application/json')) {
        data = await response.json();
      } else {
        data = await response.text();
      }

      // 构建响应对象
      const responseObj = {
        status: response.status,
        statusText: response.statusText,
        headers: response.headers,
        data,
        ok: response.ok
      };

      // 检查 HTTP 状态码
      if (!response.ok) {
        // 使用统一的错误处理器
        const errorInfo = await errorHandler.handleHTTPResponse(response);
        const error = new Error(errorInfo.message);
        error.response = responseObj;
        error.status = response.status;
        // 优先使用后端返回的错误码，如果没有则使用映射的错误码
        error.errorCode = data?.code || errorInfo.errorCode;
        
        // 判断是否应该重试
        if (enableRetry && this.retryConfig.shouldRetry && errorHandler.shouldRetry(error.errorCode)) {
          console.log(`请求失败，准备重试... 错误码：${error.errorCode}`);
          // 这里可以添加重试逻辑，或者抛出错误由上层处理
        }
        
        throw error;
      }

      // 执行响应拦截器
      return this._executeResponseInterceptors(responseObj);
    } catch (error) {
      // 使用统一的错误处理器
      const errorInfo = errorHandler.handleRequestError(error, { showError: false });
      
      // 执行错误拦截器
      const finalError = this._executeErrorInterceptors(error);
      
      // 只有在拦截器没有处理过的情况下才显示错误（避免重复提示）
      if (!finalError._handled) {
        errorHandler.showError(errorInfo.message);
      }
      
      throw finalError;
    }
  }

  /**
   * GET 请求
   */
  get(path, options = {}) {
    return this.request('GET', path, options);
  }

  /**
   * POST 请求
   */
  post(path, data, options = {}) {
    return this.request('POST', path, { ...options, data });
  }

  /**
   * PUT 请求
   */
  put(path, data, options = {}) {
    return this.request('PUT', path, { ...options, data });
  }

  /**
   * PATCH 请求
   */
  patch(path, data, options = {}) {
    return this.request('PATCH', path, { ...options, data });
  }

  /**
   * DELETE 请求
   */
  delete(path, options = {}) {
    return this.request('DELETE', path, options);
  }

  /**
   * 配置重试策略
   * @param {Object} config - 重试配置
   */
  configureRetry(config = {}) {
    if (config.maxRetries !== undefined) {
      this.retryConfig.maxRetries = config.maxRetries;
    }
    if (config.retryDelay !== undefined) {
      this.retryConfig.retryDelay = config.retryDelay;
    }
    if (config.shouldRetry !== undefined) {
      this.retryConfig.shouldRetry = config.shouldRetry;
    }
  }
}

// 创建全局 HTTP 客户端实例
const httpClient = new HttpClient();

export default httpClient;
export { HttpClient };
