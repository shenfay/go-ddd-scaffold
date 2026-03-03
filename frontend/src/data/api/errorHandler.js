/**
 * 全局错误处理器
 * 
 * 统一处理所有应用中的错误：
 * - 未捕获的 Promise 拒绝
 * - 全局错误事件
 * - API 错误
 * - React 错误
 */

import logger from '../../shared/utils/logger.js';

class ErrorHandler {
  constructor() {
    this.errorHandlers = [];
    this.isInitialized = false;
  }

  /**
   * 初始化全局错误处理
   */
  initialize() {
    if (this.isInitialized) return;

    // 处理未捕获的错误
    window.addEventListener('error', (event) => {
      this.handleError(event.error, {
        type: 'uncaught_error',
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno
      });
    });

    // 处理未捕获的 Promise 拒绝
    window.addEventListener('unhandledrejection', (event) => {
      this.handleError(event.reason, {
        type: 'unhandled_promise_rejection'
      });
    });

    // 监听自定义认证过期事件
    window.addEventListener('auth:expired', (event) => {
      this.handleAuthExpired(event.detail);
    });

    logger.info('✅ 全局错误处理器已初始化');
    this.isInitialized = true;
  }

  /**
   * 注册错误处理回调
   */
  addErrorHandler(handler) {
    this.errorHandlers.push(handler);
  }

  /**
   * 移除错误处理回调
   */
  removeErrorHandler(handler) {
    this.errorHandlers = this.errorHandlers.filter(h => h !== handler);
  }

  /**
   * 处理错误
   */
  handleError(error, context = {}) {
    // 格式化错误
    const formattedError = this._formatError(error, context);

    // 记录错误
    logger.error(formattedError.message, formattedError);

    // 通知所有错误处理器
    this.errorHandlers.forEach(handler => {
      try {
        handler(formattedError);
      } catch (err) {
        logger.error('Error handler failed:', err);
      }
    });

    return formattedError;
  }

  /**
   * 处理认证过期
   */
  handleAuthExpired(detail = {}) {
    const error = {
      code: 'AUTH_EXPIRED',
      message: detail.message || '认证已过期，请重新登录',
      type: 'auth_error',
      severity: 'medium'
    };

    logger.warn('认证已过期', error);

    // 通知所有错误处理器
    this.errorHandlers.forEach(handler => {
      try {
        handler(error);
      } catch (err) {
        logger.error('Error handler failed:', err);
      }
    });
  }

  /**
   * 格式化错误对象
   * @private
   */
  _formatError(error, context = {}) {
    if (typeof error === 'string') {
      return {
        code: 'UNKNOWN_ERROR',
        message: error,
        type: context.type || 'error',
        severity: 'medium',
        ...context
      };
    }

    if (error instanceof Error) {
      return {
        code: error.code || 'ERROR',
        message: error.message,
        stack: error.stack,
        type: context.type || 'error',
        severity: this._getSeverity(error),
        ...context
      };
    }

    if (typeof error === 'object') {
      return {
        ...error,
        code: error.code || 'UNKNOWN_ERROR',
        message: error.message || 'Unknown error',
        type: context.type || 'error',
        severity: error.severity || 'medium',
        ...context
      };
    }

    return {
      code: 'UNKNOWN_ERROR',
      message: String(error),
      type: context.type || 'error',
      severity: 'medium',
      ...context
    };
  }

  /**
   * 判断错误严重程度
   * @private
   */
  _getSeverity(error) {
    if (error.message && error.message.includes('CRITICAL')) return 'critical';
    if (error.message && error.message.includes('AUTH')) return 'high';
    if (error.code && error.code.includes('TIMEOUT')) return 'medium';
    if (error.code && error.code.includes('NETWORK')) return 'medium';
    return 'low';
  }

  /**
   * 生成错误边界组件用的错误信息
   */
  generateErrorBoundaryMessage(error, errorInfo) {
    return {
      code: 'REACT_ERROR_BOUNDARY',
      message: error.message || 'An unexpected error occurred',
      type: 'react_error',
      severity: 'high',
      componentStack: errorInfo?.componentStack || '',
      stackTrace: error.stack || ''
    };
  }
}

// 创建全局错误处理器实例
const errorHandler = new ErrorHandler();

export default errorHandler;
export { ErrorHandler };
