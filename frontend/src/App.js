import React from 'react';
import EnhancedStoryExperience from './presentation/components/story/EnhancedStoryExperience';

function App() {
  return (
    <div className="App">
      <div style={{
        minHeight: '100vh',
        backgroundColor: '#f0f9ff',
        padding: '40px',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center'
      }}>
        <h1 style={{
          color: '#1f2937',
          fontSize: '2.5rem',
          marginBottom: '20px',
          textAlign: 'center'
        }}>
          🎭 数学剧情体验系统
        </h1>
        <p style={{
          color: '#6b7280',
          fontSize: '1.2rem',
          marginBottom: '30px',
          textAlign: 'center',
          maxWidth: '600px'
        }}>
          欢迎来到毕达哥拉斯定理的发现之旅！
          点击右下角的"时空穿越"按钮开始您的数学探索之旅。
        </p>
        <div style={{
          marginTop: '40px',
          padding: '20px',
          backgroundColor: 'white',
          borderRadius: '12px',
          boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
          maxWidth: '500px'
        }}>
          <h3 style={{
            color: '#1f2937',
            marginBottom: '15px',
            textAlign: 'center'
          }}>
            📚 学习内容
          </h3>
          <ul style={{
            listStyle: 'none',
            padding: 0,
            margin: 0
          }}>
            <li style={{
              padding: '10px',
              marginBottom: '10px',
              backgroundColor: '#eff6ff',
              borderRadius: '8px',
              color: '#1e40af'
            }}>
              🔷 直角三角形的性质
            </li>
            <li style={{
              padding: '10px',
              marginBottom: '10px',
              backgroundColor: '#f0fdf4',
              borderRadius: '8px',
              color: '#065f46'
            }}>
              📐 毕达哥拉斯定理
            </li>
            <li style={{
              padding: '10px',
              backgroundColor: '#fffbeb',
              borderRadius: '8px',
              color: '#92400e'
            }}>
              🎯 数学发现的乐趣
            </li>
          </ul>
        </div>
      </div>
      <EnhancedStoryExperience />
    </div>
  );
}

export default App;