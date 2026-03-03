import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App'; // 注意：这里会查找 ./App.jsx

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);