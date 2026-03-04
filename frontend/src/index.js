import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import App from './App.jsx';
import store from './business/store';
import { initializeInterceptors } from './data/api/interceptors';

// 初始化 API 拦截器
initializeInterceptors();

// 开发环境下暴露 store 到全局，便于调试
if (process.env.NODE_ENV === 'development') {
  window.store = store;
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <Provider store={store}>
    <App />
  </Provider>
);
