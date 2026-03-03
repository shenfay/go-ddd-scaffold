/**
 * 剧情体验测试页面
 * 
 * 用于测试和演示剧情体验功能
 */

import React from 'react';
import { Provider } from 'react-redux';
import { configureStore } from '@reduxjs/toolkit';
import { StoryEngine, storyReducer } from './index.js';

// 创建测试用的Redux store
const store = configureStore({
  reducer: {
    story: storyReducer
  }
});

const StoryTestPage = () => {
  return (
    <Provider store={store}>
      <div style={{
        minHeight: '100vh',
        backgroundColor: '#f9fafb',
        padding: '20px'
      }}>
        {/* 页面标题 */}
        <div style={{
          textAlign: 'center',
          marginBottom: '40px'
        }}>
          <h1 style={{
            fontSize: '32px',
            fontWeight: '700',
            color: '#1f2937',
            marginBottom: '12px'
          }}>
            🎭 数学剧情体验测试
          </h1>
          <p style={{
            fontSize: '18px',
            color: '#6b7280',
            maxWidth: '600px',
            margin: '0 auto'
          }}>
            点击右下角的"时空穿越"按钮开始体验毕达哥拉斯定理的发现之旅
          </p>
        </div>
        
        {/* 测试内容区域 */}
        <div style={{
          maxWidth: '800px',
          margin: '0 auto',
          backgroundColor: 'white',
          borderRadius: '16px',
          padding: '30px',
          boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
        }}>
          <h2 style={{
            fontSize: '24px',
            fontWeight: '600',
            color: '#1f2937',
            marginBottom: '20px'
          }}>
            📚 学习内容展示
          </h2>
          
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
            gap: '20px',
            marginBottom: '30px'
          }}>
            <div style={{
              padding: '20px',
              backgroundColor: '#eff6ff',
              borderRadius: '12px',
              border: '2px solid #3b82f6'
            }}>
              <h3 style={{ 
                margin: '0 0 10px 0', 
                color: '#1e40af',
                fontSize: '18px'
              }}>
                🔷 直角三角形
              </h3>
              <p style={{ 
                margin: 0, 
                color: '#374151',
                fontSize: '14px'
              }}>
                具有一个90度角的三角形，两条直角边与斜边满足特殊关系
              </p>
            </div>
            
            <div style={{
              padding: '20px',
              backgroundColor: '#f0fdf4',
              borderRadius: '12px',
              border: '2px solid #10b981'
            }}>
              <h3 style={{ 
                margin: '0 0 10px 0', 
                color: '#065f46',
                fontSize: '18px'
              }}>
                📐 毕达哥拉斯定理
              </h3>
              <p style={{ 
                margin: 0, 
                color: '#374151',
                fontSize: '14px'
              }}>
                直角三角形中，两直角边的平方和等于斜边的平方
              </p>
            </div>
            
            <div style={{
              padding: '20px',
              backgroundColor: '#fffbeb',
              borderRadius: '12px',
              border: '2px solid #f59e0b'
            }}>
              <h3 style={{ 
                margin: '0 0 10px 0', 
                color: '#92400e',
                fontSize: '18px'
              }}>
                🎯 学习目标
              </h3>
              <p style={{ 
                margin: 0, 
                color: '#374151',
                fontSize: '14px'
              }}>
                通过互动体验理解并验证毕达哥拉斯定理的数学原理
              </p>
            </div>
          </div>
          
          <div style={{
            backgroundColor: '#f3f4f6',
            padding: '20px',
            borderRadius: '12px',
            border: '1px solid #e5e7eb'
          }}>
            <h3 style={{
              margin: '0 0 15px 0',
              color: '#374151',
              fontSize: '18px'
            }}>
              📋 功能说明
            </h3>
            <ul style={{
              margin: 0,
              padding: '0 0 0 20px',
              color: '#6b7280',
              fontSize: '14px'
            }}>
              <li style={{ marginBottom: '8px' }}>点击"时空穿越"按钮开始剧情体验</li>
              <li style={{ marginBottom: '8px' }}>与古希腊数学家毕达哥拉斯进行互动对话</li>
              <li style={{ marginBottom: '8px' }}>通过测量不同三角形验证定理</li>
              <li style={{ marginBottom: '8px' }}>完成任务后获得成就奖励</li>
              <li>支持场景切换和进度跟踪</li>
            </ul>
          </div>
        </div>
        
        {/* 剧情引擎 */}
        <StoryEngine>
          {/* 这里可以放置主应用内容 */}
        </StoryEngine>
      </div>
    </Provider>
  );
};

export default StoryTestPage;