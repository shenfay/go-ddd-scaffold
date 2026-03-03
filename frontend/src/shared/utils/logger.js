/**
 * 日志管理系统
 * 
 * 统一管理所有日志记录，支持：
 * - 不同日志级别（DEBUG, INFO, WARN, ERROR）
 * - 日志格式化
 * - 日志收集
 * - 环境感知（开发环境显示更多细节）
 */

const LogLevel = {
  DEBUG: 0,
  INFO: 1,
  WARN: 2,
  ERROR: 3
};

class Logger {
  constructor() {
    this.logLevel = process.env.NODE_ENV === 'production' ? LogLevel.INFO : LogLevel.DEBUG;
    this.logs = [];
    this.maxLogs = 1000;
    this.enableConsole = true;
    this.enableRemoteLogging = false;
    this.remoteEndpoint = null;
  }

  /**
   * 设置日志级别
   */
  setLogLevel(level) {
    this.logLevel = level;
  }

  /**
   * 启用/禁用远程日志
   */
  enableRemote(endpoint) {
    this.enableRemoteLogging = true;
    this.remoteEndpoint = endpoint;
  }

  /**
   * 禁用远程日志
   */
  disableRemote() {
    this.enableRemoteLogging = false;
  }

  /**
   * 格式化日志信息
   * @private
   */
  _formatLog(level, message, data) {
    const timestamp = new Date().toISOString();
    return {
      timestamp,
      level: Object.keys(LogLevel).find(key => LogLevel[key] === level),
      message,
      data,
      url: typeof window !== 'undefined' ? window.location.href : 'N/A'
    };
  }

  /**
   * 记录调试日志
   */
  debug(message, data = null) {
    if (this.logLevel <= LogLevel.DEBUG) {
      const log = this._formatLog(LogLevel.DEBUG, message, data);
      this._store(log);
      
      if (this.enableConsole) {
        console.log(`%c[DEBUG] ${message}`, 'color: #0066cc', data);
      }
    }
  }

  /**
   * 记录信息日志
   */
  info(message, data = null) {
    if (this.logLevel <= LogLevel.INFO) {
      const log = this._formatLog(LogLevel.INFO, message, data);
      this._store(log);
      
      if (this.enableConsole) {
        console.log(`%c[INFO] ${message}`, 'color: #00cc00', data);
      }
    }
  }

  /**
   * 记录警告日志
   */
  warn(message, data = null) {
    if (this.logLevel <= LogLevel.WARN) {
      const log = this._formatLog(LogLevel.WARN, message, data);
      this._store(log);
      
      if (this.enableConsole) {
        console.warn(`%c[WARN] ${message}`, 'color: #ffaa00', data);
      }
    }
  }

  /**
   * 记录错误日志
   */
  error(message, error = null) {
    if (this.logLevel <= LogLevel.ERROR) {
      const data = error instanceof Error
        ? {
            message: error.message,
            stack: error.stack,
            code: error.code
          }
        : error;

      const log = this._formatLog(LogLevel.ERROR, message, data);
      this._store(log);
      
      if (this.enableConsole) {
        console.error(`%c[ERROR] ${message}`, 'color: #ff0000', data);
      }

      // 发送到远程
      if (this.enableRemoteLogging) {
        this._sendToRemote(log);
      }
    }
  }

  /**
   * 存储日志
   * @private
   */
  _store(log) {
    this.logs.push(log);
    
    // 限制日志数量
    if (this.logs.length > this.maxLogs) {
      this.logs.shift();
    }
  }

  /**
   * 发送日志到远程服务
   * @private
   */
  async _sendToRemote(log) {
    if (!this.remoteEndpoint) return;

    try {
      await fetch(this.remoteEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(log)
      });
    } catch (err) {
      // 远程日志失败不应该阻塞应用
      console.warn('Failed to send remote log:', err);
    }
  }

  /**
   * 获取所有日志
   */
  getLogs() {
    return [...this.logs];
  }

  /**
   * 获取特定级别的日志
   */
  getLogsByLevel(level) {
    return this.logs.filter(log => LogLevel[log.level] === level);
  }

  /**
   * 清空日志
   */
  clearLogs() {
    this.logs = [];
  }

  /**
   * 导出日志为 JSON
   */
  exportLogs() {
    return JSON.stringify(this.logs, null, 2);
  }

  /**
   * 导出日志为 CSV
   */
  exportLogsAsCSV() {
    const headers = ['Timestamp', 'Level', 'Message', 'Data', 'URL'];
    const rows = this.logs.map(log => [
      log.timestamp,
      log.level,
      log.message,
      typeof log.data === 'object' ? JSON.stringify(log.data) : log.data,
      log.url
    ]);

    const csvContent = [
      headers.join(','),
      ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
    ].join('\n');

    return csvContent;
  }

  /**
   * 下载日志文件
   */
  downloadLogs(filename = 'logs.json', format = 'json') {
    const content = format === 'csv' ? this.exportLogsAsCSV() : this.exportLogs();
    const element = document.createElement('a');
    const file = new Blob([content], { type: 'text/plain' });

    element.href = URL.createObjectURL(file);
    element.download = filename;
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
  }

  /**
   * 绩效计时
   */
  time(label) {
    console.time(label);
  }

  /**
   * 结束计时
   */
  timeEnd(label) {
    console.timeEnd(label);
  }
}

// 创建全局 logger 实例
const logger = new Logger();

export default logger;
export { Logger, LogLevel };
