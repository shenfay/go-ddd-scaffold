import React, { useState, useEffect } from 'react';

const SimpleStoryExperience = () => {
  const [isStoryActive, setIsStoryActive] = useState(false);
  const [currentScene, setCurrentScene] = useState('modern');
  const [progress, setProgress] = useState(0);
  const [showPythagoras, setShowPythagoras] = useState(false);
  const [showTriangle, setShowTriangle] = useState(false);
  const [showAchievement, setShowAchievement] = useState(false);

  const scenes = {
    modern: {
      title: '现代教室',
      description: '你在数学课上学习直角三角形...',
      bgColor: '#dbeafe',
      icon: '📚',
      character: null
    },
    ancient: {
      title: '古希腊集市',
      description: '你穿越到了古希腊，遇到了数学家毕达哥拉斯...',
      bgColor: '#fef3c7',
      icon: '🏛️',
      character: '🧙‍♂️ 毕达哥拉斯'
    },
    discovery: {
      title: '发现时刻',
      description: '通过测量验证，你发现了毕达哥拉斯定理！',
      bgColor: '#dcfce7',
      icon: '🔍',
      character: '✨ 数学之光'
    }
  };

  useEffect(() => {
    if (currentScene === 'ancient') {
      const timer = setTimeout(() => {
        setShowPythagoras(true);
      }, 1000);
      return () => clearTimeout(timer);
    } else {
      setShowPythagoras(false);
    }
  }, [currentScene]);

  useEffect(() => {
    if (currentScene === 'discovery') {
      const timer = setTimeout(() => {
        setShowTriangle(true);
      }, 1000);
      return () => clearTimeout(timer);
    } else {
      setShowTriangle(false);
    }
  }, [currentScene]);

  const handleStartStory = () => {
    setIsStoryActive(true);
    setProgress(0);
    setCurrentScene('modern');
  };

  const handleSceneChange = (scene) => {
    setCurrentScene(scene);
    const progressMap = {
      modern: 0,
      ancient: 33,
      discovery: 66
    };
    setProgress(progressMap[scene] || 0);
  };

  const handleCompleteStory = () => {
    setProgress(100);
    setShowAchievement(true);
    setTimeout(() => {
      setShowAchievement(false);
      alert('🎉 恭喜！你完成了毕达哥拉斯定理的发现之旅！\n\n获得了"数学探索者"成就！');
    }, 2000);
  };

  if (!isStoryActive) {
    return (
      <div style={{
        position: 'fixed',
        bottom: '20px',
        right: '20px',
        zIndex: 1000,
        animation: 'float 2s ease-in-out infinite'
      }}>
        <button
          onClick={handleStartStory}
          style={{
            padding: '16px 24px',
            backgroundColor: '#4f46e5',
            color: 'white',
            border: 'none',
            borderRadius: '12px',
            fontSize: '18px',
            fontWeight: '600',
            cursor: 'pointer',
            boxShadow: '0 8px 15px rgba(79, 70, 229, 0.3)',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            transform: 'scale(1)',
            transition: 'all 0.3s ease'
          }}
          onMouseOver={(e) => {
            e.target.style.transform = 'scale(1.05) translateY(-2px)';
            e.target.style.boxShadow = '0 12px 20px rgba(79, 70, 229, 0.4)';
          }}
          onMouseOut={(e) => {
            e.target.style.transform = 'scale(1)';
            e.target.style.boxShadow = '0 8px 15px rgba(79, 70, 229, 0.3)';
          }}
        >
          <span>🚀</span>
          时空穿越
        </button>
      </div>
    );
  }

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100vw',
      height: '100vh',
      zIndex: 2000,
      backgroundColor: scenes[currentScene].bgColor,
      display: 'flex',
      flexDirection: 'column',
      transition: 'background-color 0.5s ease'
    }}>
      {/* 顶部进度条 */}
      <div style={{
        padding: '20px',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        borderBottom: '1px solid rgba(0, 0, 0, 0.1)'
      }}>
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          marginBottom: '10px'
        }}>
          <h2 style={{ margin: 0, color: '#1f2937', display: 'flex', alignItems: 'center', gap: '10px' }}>
            <span>{scenes[currentScene].icon}</span>
            {scenes[currentScene].title}
          </h2>
          <div style={{
            padding: '6px 12px',
            backgroundColor: '#e5e7eb',
            borderRadius: '20px',
            fontSize: '14px',
            fontWeight: '500'
          }}>
            进度: {progress}%
          </div>
        </div>
        <div style={{
          width: '100%',
          height: '8px',
          backgroundColor: '#e5e7eb',
          borderRadius: '4px',
          overflow: 'hidden'
        }}>
          <div style={{
            width: `${progress}%`,
            height: '100%',
            backgroundColor: '#4f46e5',
            transition: 'width 0.8s cubic-bezier(0.25, 0.46, 0.45, 0.94)'
          }} />
        </div>
      </div>

      {/* 主要内容区域 */}
      <div style={{
        flex: 1,
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'center',
        alignItems: 'center',
        padding: '40px',
        textAlign: 'center',
        position: 'relative'
      }}>
        {/* 角色展示 */}
        {scenes[currentScene].character && (
          <div 
            style={{
              position: 'absolute',
              top: '20px',
              right: '40px',
              opacity: showPythagoras ? 1 : 0,
              transform: showPythagoras ? 'translateY(0)' : 'translateY(-20px)',
              transition: 'all 0.5s ease',
              fontSize: '48px',
              animation: showPythagoras ? 'bounce 1s ease infinite alternate' : 'none'
            }}
          >
            {scenes[currentScene].character}
          </div>
        )}

        <div 
          style={{
            maxWidth: '600px',
            backgroundColor: 'white',
            borderRadius: '16px',
            padding: '40px',
            boxShadow: '0 10px 25px rgba(0, 0, 0, 0.1)',
            marginBottom: '30px',
            transform: 'scale(1)',
            transition: 'transform 0.3s ease',
            animation: 'fadeInUp 0.6s ease'
          }}
          onMouseEnter={(e) => e.target.style.transform = 'scale(1.02)'}
          onMouseLeave={(e) => e.target.style.transform = 'scale(1)'}
        >
          <h3 style={{
            color: '#1f2937',
            fontSize: '24px',
            marginBottom: '20px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            gap: '10px'
          }}>
            {scenes[currentScene].icon} {scenes[currentScene].title}
          </h3>
          <p style={{
            color: '#6b7280',
            fontSize: '18px',
            lineHeight: '1.6',
            marginBottom: '30px',
            animation: 'fadeIn 0.8s ease'
          }}>
            {scenes[currentScene].description}
          </p>
          
          {currentScene === 'discovery' && (
            <div 
              style={{
                padding: '20px',
                backgroundColor: '#f0fdf4',
                borderRadius: '12px',
                border: '2px solid #10b981',
                opacity: showTriangle ? 1 : 0,
                transform: showTriangle ? 'scale(1)' : 'scale(0.8)',
                transition: 'all 0.5s ease 0.3s'
              }}
            >
              <h4 style={{ 
                margin: '0 0 15px 0', 
                color: '#065f46',
                fontSize: '20px'
              }}>
                📐 毕达哥拉斯定理
              </h4>
              <p style={{ 
                margin: 0, 
                color: '#374151',
                fontSize: '16px'
              }}>
                在直角三角形中，两直角边的平方和等于斜边的平方：<br/>
                <strong style={{
                  fontSize: '20px',
                  color: '#4f46e5',
                  display: 'block',
                  margin: '10px 0',
                  animation: 'pulse 2s infinite'
                }}>a² + b² = c²</strong>
              </p>
              
              {/* 三角形示意图 */}
              <div style={{
                marginTop: '15px',
                display: 'flex',
                justifyContent: 'center'
              }}>
                <svg width="150" height="150" viewBox="0 0 150 150">
                  <polygon 
                    points="30,120 120,120 120,30" 
                    fill="none" 
                    stroke="#4f46e5" 
                    strokeWidth="3"
                    style={{
                      opacity: showTriangle ? 1 : 0,
                      transform: showTriangle ? 'scale(1)' : 'scale(0.5)',
                      transition: 'all 0.8s ease 0.5s'
                    }}
                  />
                  <text x="75" y="140" textAnchor="middle" fill="#6b7280">a</text>
                  <text x="125" y="75" textAnchor="middle" fill="#6b7280">b</text>
                  <text x="25" y="80" textAnchor="middle" fill="#6b7280">c</text>
                </svg>
              </div>
            </div>
          )}
        </div>

        {/* 场景切换按钮 */}
        <div style={{
          display: 'flex',
          gap: '15px',
          marginBottom: '30px',
          animation: 'fadeIn 1s ease'
        }}>
          {Object.keys(scenes).map((sceneKey) => (
            sceneKey !== currentScene && (
              <button
                key={sceneKey}
                onClick={() => handleSceneChange(sceneKey)}
                style={{
                  padding: '12px 24px',
                  backgroundColor: '#6366f1',
                  color: 'white',
                  border: 'none',
                  borderRadius: '8px',
                  fontSize: '16px',
                  cursor: 'pointer',
                  transition: 'all 0.3s ease',
                  boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
                  transform: 'translateY(0)'
                }}
                onMouseOver={(e) => {
                  e.target.style.backgroundColor = '#4f46e5';
                  e.target.style.transform = 'translateY(-3px)';
                  e.target.style.boxShadow = '0 6px 12px rgba(0, 0, 0, 0.15)';
                }}
                onMouseOut={(e) => {
                  e.target.style.backgroundColor = '#6366f1';
                  e.target.style.transform = 'translateY(0)';
                  e.target.style.boxShadow = '0 4px 6px rgba(0, 0, 0, 0.1)';
                }}
              >
                前往 {scenes[sceneKey].title}
              </button>
            )
          ))}
        </div>

        {/* 完成按钮 */}
        {currentScene === 'discovery' && progress < 100 && (
          <button
            onClick={handleCompleteStory}
            style={{
              padding: '16px 32px',
              backgroundColor: '#10b981',
              color: 'white',
              border: 'none',
              borderRadius: '12px',
              fontSize: '18px',
              fontWeight: '600',
              cursor: 'pointer',
              boxShadow: '0 8px 15px rgba(16, 185, 129, 0.3)',
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              transform: 'scale(1)',
              transition: 'all 0.3s ease'
            }}
            onMouseOver={(e) => {
              e.target.style.transform = 'scale(1.05) translateY(-2px)';
              e.target.style.boxShadow = '0 12px 20px rgba(16, 185, 129, 0.4)';
            }}
            onMouseOut={(e) => {
              e.target.style.transform = 'scale(1)';
              e.target.style.boxShadow = '0 8px 15px rgba(16, 185, 129, 0.3)';
            }}
          >
            <span>🎓</span>
            完成探索
          </button>
        )}

        {/* 成就弹窗 */}
        {showAchievement && (
          <div style={{
            position: 'fixed',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            backgroundColor: 'white',
            padding: '40px',
            borderRadius: '20px',
            boxShadow: '0 20px 40px rgba(0, 0, 0, 0.3)',
            zIndex: 3000,
            textAlign: 'center',
            animation: 'achievementPop 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275)'
          }}>
            <div style={{ fontSize: '60px', marginBottom: '20px' }}>🎉</div>
            <h3 style={{ color: '#1f2937', margin: '10px 0' }}>成就解锁！</h3>
            <p style={{ color: '#6b7280', margin: '10px 0' }}>数学探索者</p>
            <p style={{ color: '#9ca3af', fontSize: '14px' }}>成功完成毕达哥拉斯定理发现之旅</p>
          </div>
        )}

        {/* 退出按钮 */}
        <button
          onClick={() => setIsStoryActive(false)}
          style={{
            marginTop: '20px',
            padding: '10px 20px',
            backgroundColor: '#ef4444',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            fontSize: '16px',
            cursor: 'pointer',
            transition: 'all 0.2s ease'
          }}
          onMouseOver={(e) => e.target.style.backgroundColor = '#dc2626'}
          onMouseOut={(e) => e.target.style.backgroundColor = '#ef4444'}
        >
          退出剧情
        </button>
      </div>

      {/* 添加CSS动画 */}
      <style>{`
        @keyframes float {
          0% { transform: translateY(0px); }
          50% { transform: translateY(-10px); }
          100% { transform: translateY(0px); }
        }
        
        @keyframes bounce {
          0% { transform: translateY(0); }
          100% { transform: translateY(-10px); }
        }
        
        @keyframes fadeInUp {
          from {
            opacity: 0;
            transform: translate3d(0, 30px, 0);
          }
          to {
            opacity: 1;
            transform: translate3d(0, 0, 0);
          }
        }
        
        @keyframes fadeIn {
          from { opacity: 0; }
          to { opacity: 1; }
        }
        
        @keyframes pulse {
          0% { transform: scale(1); }
          50% { transform: scale(1.05); }
          100% { transform: scale(1); }
        }
        
        @keyframes achievementPop {
          0% { 
            opacity: 0;
            transform: translate(-50%, -50%) scale(0.5);
          }
          70% { 
            transform: translate(-50%, -50%) scale(1.1);
          }
          100% { 
            opacity: 1;
            transform: translate(-50%, -50%) scale(1);
          }
        }
      `}</style>
    </div>
  );
};

export default SimpleStoryExperience;