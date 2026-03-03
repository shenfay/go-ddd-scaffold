/**
 * React 错误边界组件
 * 
 * 捕获并处理子组件中的 React 错误
 */

import React from 'react';
import errorHandler from '../../data/api/errorHandler.js';
import logger from '../../shared/utils/logger.js';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null
    };
  }

  static getDerivedStateFromError(error) {
    return { hasError: true };
  }

  componentDidCatch(error, errorInfo) {
    // 记录错误
    this.setState({
      error,
      errorInfo
    });

    // 通知全局错误处理器
    const formattedError = errorHandler.generateErrorBoundaryMessage(error, errorInfo);
    errorHandler.handleError(formattedError);
  }

  handleReset = () => {
    this.setState({
      hasError: false,
      error: null,
      errorInfo: null
    });
  };

  render() {
    if (this.state.hasError) {
      return (
        <div
          style={{
            padding: '20px',
            margin: '20px',
            border: '1px solid #ff0000',
            borderRadius: '4px',
            backgroundColor: '#ffe6e6',
            color: '#cc0000',
            fontFamily: 'system-ui, -apple-system, sans-serif'
          }}
        >
          <h2>⚠️ 出现了一个错误</h2>
          <p>{this.state.error && this.state.error.message}</p>
          
          {process.env.NODE_ENV !== 'production' && (
            <details
              style={{
                whiteSpace: 'pre-wrap',
                marginTop: '10px',
                padding: '10px',
                backgroundColor: '#fff',
                borderRadius: '4px',
                color: '#000',
                fontSize: '12px',
                maxHeight: '300px',
                overflow: 'auto'
              }}
            >
              <summary style={{ cursor: 'pointer', marginBottom: '10px' }}>
                详细错误信息
              </summary>
              <p>{this.state.errorInfo && this.state.errorInfo.componentStack}</p>
            </details>
          )}

          <button
            onClick={this.handleReset}
            style={{
              marginTop: '10px',
              padding: '8px 16px',
              backgroundColor: '#ff0000',
              color: '#fff',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            重新尝试
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
