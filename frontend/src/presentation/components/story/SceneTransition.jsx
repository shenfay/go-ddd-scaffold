/**
 * 场景切换组件
 * 
 * 处理不同场景之间的切换动画和过渡效果
 */

import React, { useState, useEffect } from 'react';

const SceneTransition = ({ currentScene, onSceneChange }) => {
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [previousScene, setPreviousScene] = useState(currentScene);

  // 场景配置
  const sceneConfig = {
    modern: {
      name: '现代教室',
      backgroundColor: '#dbeafe',
      description: '你在现代的数学教室里学习',
      icon: '📚'
    },
    ancient: {
      name: '古希腊集市',
      backgroundColor: '#fef3c7',
      description: '你穿越到了古希腊的集市',
      icon: '🏛️'
    },
    discovery: {
      name: '发现时刻',
      backgroundColor: '#dcfce7',
      description: '你正在见证伟大的发现',
      icon: '✨'
    }
  };

  // 处理场景变化
  useEffect(() => {
    if (currentScene !== previousScene) {
      setIsTransitioning(true);
      
      // 模拟过渡动画时间
      const timer = setTimeout(() => {
        setIsTransitioning(false);
        setPreviousScene(currentScene);
      }, 800);
      
      return () => clearTimeout(timer);
    }
  }, [currentScene, previousScene]);

  // 获取场景样式
  const getSceneStyle = (scene) => {
    const config = sceneConfig[scene] || sceneConfig.modern;
    return {
      backgroundColor: config.backgroundColor,
      transition: 'all 0.8s ease-in-out',
      opacity: isTransitioning && scene === currentScene ? 0 : 1
    };
  };

  // 场景切换按钮
  const renderSceneButtons = () => {
    const scenes = Object.keys(sceneConfig);
    
    return (
      <div style={{
        position: 'absolute',
        top: '20px',
        left: '20px',
        display: 'flex',
        gap: '10px',
        zIndex: 100
      }}>
        {scenes.map(scene => {
          const config = sceneConfig[scene];
          const isActive = currentScene === scene;
          
          return (
            <button
              key={scene}
              onClick={() => !isTransitioning && onSceneChange(scene)}
              disabled={isTransitioning || isActive}
              style={{
                padding: '8px 16px',
                backgroundColor: isActive ? '#3b82f6' : 'rgba(255, 255, 255, 0.8)',
                color: isActive ? 'white' : '#374151',
                border: `2px solid ${isActive ? '#3b82f6' : '#d1d5db'}`,
                borderRadius: '6px',
                cursor: isTransitioning || isActive ? 'not-allowed' : 'pointer',
                opacity: isTransitioning ? 0.5 : 1,
                transition: 'all 0.3s ease',
                display: 'flex',
                alignItems: 'center',
                gap: '4px',
                fontSize: '14px'
              }}
            >
              <span>{config.icon}</span>
              <span>{config.name}</span>
            </button>
          );
        })}
      </div>
    );
  };

  // 场景描述
  const renderSceneDescription = () => {
    const config = sceneConfig[currentScene] || sceneConfig.modern;
    
    return (
      <div style={{
        position: 'absolute',
        top: '80px',
        left: '20px',
        padding: '12px 16px',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderRadius: '8px',
        boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
        maxWidth: '300px',
        zIndex: 90
      }}>
        <h3 style={{
          margin: '0 0 8px 0',
          fontSize: '18px',
          fontWeight: '600',
          color: '#1f2937'
        }}>
          {config.icon} {config.name}
        </h3>
        <p style={{
          margin: 0,
          fontSize: '14px',
          color: '#6b7280'
        }}>
          {config.description}
        </p>
      </div>
    );
  };

  // 过渡动画效果
  const renderTransitionEffect = () => {
    if (!isTransitioning) return null;
    
    return (
      <div style={{
        position: 'absolute',
        top: 0,
        left: 0,
        width: '100%',
        height: '100%',
        backgroundColor: 'rgba(255, 255, 255, 0.8)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        zIndex: 1000,
        animation: 'pulse 0.8s ease-in-out'
      }}>
        <div style={{
          textAlign: 'center'
        }}>
          <div style={{
            fontSize: '48px',
            marginBottom: '16px',
            animation: 'spin 0.8s linear infinite'
          }}>
            🌀
          </div>
          <div style={{
            fontSize: '18px',
            fontWeight: '500',
            color: '#374151'
          }}>
            时空穿越中...
          </div>
        </div>
      </div>
    );
  };

  return (
    <div style={{
      position: 'relative',
      width: '100%',
      height: '100%',
      ...getSceneStyle(currentScene)
    }}>
      {/* 场景切换按钮 */}
      {renderSceneButtons()}
      
      {/* 场景描述 */}
      {renderSceneDescription()}
      
      {/* 过渡动画 */}
      {renderTransitionEffect()}
      
      {/* 添加CSS动画 */}
      <style>{`
        @keyframes pulse {
          0% { opacity: 0; }
          50% { opacity: 1; }
          100% { opacity: 0; }
        }
        
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  );
};

export default SceneTransition;