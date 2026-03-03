/**
 * 成就解锁组件
 * 
 * 处理成就解锁的动画效果和展示
 */

import React, { useState, useEffect } from 'react';

const AchievementUnlock = ({ recentUnlocks, onAchievementShown }) => {
  const [currentAchievement, setCurrentAchievement] = useState(null);
  const [showAnimation, setShowAnimation] = useState(false);
  const [animationStage, setAnimationStage] = useState('hidden'); // hidden, appear, show, disappear

  // 处理新成就解锁
  useEffect(() => {
    if (recentUnlocks && recentUnlocks.length > 0) {
      const latestAchievement = recentUnlocks[recentUnlocks.length - 1];
      setCurrentAchievement(latestAchievement);
      setShowAnimation(true);
      setAnimationStage('appear');
      
      // 启动动画序列
      const timers = [
        setTimeout(() => setAnimationStage('show'), 300),      // 出现动画
        setTimeout(() => setAnimationStage('disappear'), 3000), // 显示3秒后消失
        setTimeout(() => {
          setShowAnimation(false);
          setAnimationStage('hidden');
          onAchievementShown();
        }, 3500) // 完全消失后清理
      ];
      
      return () => timers.forEach(timer => clearTimeout(timer));
    }
  }, [recentUnlocks, onAchievementShown]);

  // 如果没有成就要显示，返回null
  if (!showAnimation || !currentAchievement) {
    return null;
  }

  // 成就类型配置
  const achievementTypes = {
    pythagorean_discovery: {
      name: '毕达哥拉斯发现者',
      description: '成功体验了毕达哥拉斯定理的发现过程',
      icon: '🎓',
      color: '#3b82f6',
      rarity: 'legendary',
      particleEffect: 'sparkle'
    },
    explorer: {
      name: '数学探索家',
      description: '完成第一个数学发现之旅',
      icon: '🧭',
      color: '#10b981',
      rarity: 'rare',
      particleEffect: 'stars'
    },
    master: {
      name: '定理大师',
      description: '掌握多个重要数学定理',
      icon: '👑',
      color: '#f59e0b',
      rarity: 'epic',
      particleEffect: 'confetti'
    }
  };

  const achievementConfig = achievementTypes[currentAchievement.id] || achievementTypes.pythagorean_discovery;

  // 渲染粒子效果
  const renderParticleEffect = () => {
    const particles = [];
    const particleCount = 20;
    
    for (let i = 0; i < particleCount; i++) {
      particles.push(
        <div
          key={i}
          style={{
            position: 'absolute',
            width: '8px',
            height: '8px',
            backgroundColor: achievementConfig.color,
            borderRadius: '50%',
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animation: `float-${i} 2s ease-in-out infinite`,
            opacity: animationStage === 'show' ? 1 : 0,
            transition: 'opacity 0.3s ease'
          }}
        ></div>
      );
    }
    
    return particles;
  };

  // 渲染成就卡片
  const renderAchievementCard = () => {
    const baseStyle = {
      position: 'fixed',
      top: '50%',
      left: '50%',
      transform: 'translate(-50%, -50%)',
      zIndex: 5000,
      width: '400px',
      maxWidth: '90vw',
      textAlign: 'center'
    };
    
    // 根据动画阶段调整样式
    let cardStyle = { ...baseStyle };
    
    switch (animationStage) {
      case 'appear':
        cardStyle = {
          ...baseStyle,
          transform: 'translate(-50%, -50%) scale(0.8)',
          opacity: 0
        };
        break;
      case 'show':
        cardStyle = {
          ...baseStyle,
          transform: 'translate(-50%, -50%) scale(1)',
          opacity: 1
        };
        break;
      case 'disappear':
        cardStyle = {
          ...baseStyle,
          transform: 'translate(-50%, -50%) scale(1.1)',
          opacity: 0
        };
        break;
      default:
        cardStyle = {
          ...baseStyle,
          opacity: 0
        };
    }
    
    return (
      <div
        style={{
          ...cardStyle,
          transition: 'all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)'
        }}
      >
        {/* 粒子效果容器 */}
        <div style={{
          position: 'absolute',
          top: '-20px',
          left: '-20px',
          right: '-20px',
          bottom: '-20px',
          pointerEvents: 'none'
        }}>
          {renderParticleEffect()}
        </div>
        
        {/* 成就卡片 */}
        <div style={{
          backgroundColor: 'white',
          borderRadius: '20px',
          padding: '30px',
          boxShadow: '0 20px 40px rgba(0, 0, 0, 0.2)',
          border: `4px solid ${achievementConfig.color}`,
          position: 'relative',
          overflow: 'hidden'
        }}>
          {/* 背景渐变效果 */}
          <div style={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            background: `linear-gradient(135deg, ${achievementConfig.color}10 0%, ${achievementConfig.color}05 100%)`,
            zIndex: -1
          }}></div>
          
          {/* 解锁标题 */}
          <div style={{
            marginBottom: '20px'
          }}>
            <div style={{
              fontSize: '20px',
              fontWeight: '600',
              color: '#6b7280',
              marginBottom: '5px'
            }}>
              🎉 成就解锁！
            </div>
          </div>
          
          {/* 成就图标 */}
          <div style={{
            fontSize: '60px',
            marginBottom: '15px',
            animation: animationStage === 'show' ? 'pulse 1s ease-in-out' : 'none'
          }}>
            {achievementConfig.icon}
          </div>
          
          {/* 成就名称 */}
          <h2 style={{
            margin: '0 0 10px 0',
            fontSize: '24px',
            fontWeight: '700',
            color: '#1f2937'
          }}>
            {currentAchievement.name}
          </h2>
          
          {/* 成就描述 */}
          <p style={{
            margin: '0 0 20px 0',
            fontSize: '16px',
            color: '#6b7280',
            lineHeight: '1.5'
          }}>
            {currentAchievement.description}
          </p>
          
          {/* 稀有度标签 */}
          <div style={{
            display: 'inline-block',
            padding: '6px 16px',
            backgroundColor: achievementConfig.color,
            color: 'white',
            borderRadius: '20px',
            fontSize: '14px',
            fontWeight: '500'
          }}>
            {achievementConfig.rarity === 'legendary' && '🏆 传奇'}
            {achievementConfig.rarity === 'epic' && '💎 史诗'}
            {achievementConfig.rarity === 'rare' && '⭐ 稀有'}
          </div>
          
          {/* 时间戳 */}
          <div style={{
            marginTop: '15px',
            fontSize: '12px',
            color: '#9ca3af'
          }}>
            解锁于 {new Date(currentAchievement.unlockedAt).toLocaleString('zh-CN')}
          </div>
        </div>
      </div>
    );
  };

  // 渲染屏幕闪光效果
  const renderScreenFlash = () => {
    if (animationStage !== 'appear') return null;
    
    return (
      <div style={{
        position: 'fixed',
        top: 0,
        left: 0,
        width: '100vw',
        height: '100vh',
        backgroundColor: achievementConfig.color,
        opacity: 0.3,
        zIndex: 4999,
        animation: 'flash 0.5s ease-out'
      }}></div>
    );
  };

  return (
    <div>
      {renderScreenFlash()}
      {renderAchievementCard()}
      
      {/* CSS动画定义 */}
      <style>{`
        @keyframes pulse {
          0% { transform: scale(1); }
          50% { transform: scale(1.1); }
          100% { transform: scale(1); }
        }
        
        @keyframes flash {
          0% { opacity: 0; }
          50% { opacity: 0.4; }
          100% { opacity: 0; }
        }
        
        ${Array.from({ length: 20 }, (_, i) => `
          @keyframes float-${i} {
            0% { 
              transform: translate(0, 0) rotate(0deg);
              opacity: 0;
            }
            10% { opacity: 1; }
            90% { opacity: 1; }
            100% { 
              transform: translate(${(Math.random() - 0.5) * 100}px, ${(Math.random() - 0.5) * 100}px) rotate(360deg);
              opacity: 0;
            }
          }
        `).join('')}
      `}</style>
    </div>
  );
};

export default AchievementUnlock;