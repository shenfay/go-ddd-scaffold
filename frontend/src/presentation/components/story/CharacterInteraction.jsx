/**
 * NPC交互组件
 * 
 * 处理与毕达哥拉斯等NPC的交互和对话
 */

import React, { useState, useEffect } from 'react';

const CharacterInteraction = ({ scene, onTaskStart }) => {
  const [currentDialog, setCurrentDialog] = useState(0);
  const [isTalking, setIsTalking] = useState(false);
  const [showOptions, setShowOptions] = useState(false);

  // NPC配置
  const npcConfig = {
    pythagoras: {
      name: '毕达哥拉斯',
      avatar: '🧙‍♂️',
      color: '#8b5cf6',
      dialogs: {
        modern: [
          {
            text: '你好！我是古希腊的数学家毕达哥拉斯。',
            action: 'greeting'
          },
          {
            text: '我想带你回到我的时代，一起发现著名的毕达哥拉斯定理！',
            action: 'invitation'
          },
          {
            text: '准备好了吗？让我们开始这段奇妙的数学之旅！',
            action: 'start',
            options: [
              { text: '准备好了！', action: 'accept' },
              { text: '再等一下', action: 'wait' }
            ]
          }
        ],
        ancient: [
          {
            text: '欢迎来到古希腊！这里是公元前6世纪的萨摩斯岛。',
            action: 'welcome'
          },
          {
            text: '你看，我正在研究这些直角三角形...',
            action: 'introduction'
          },
          {
            text: '我发现了一个神奇的规律，你想和我一起探索吗？',
            action: 'discovery',
            options: [
              { text: '当然想！', action: 'explore' },
              { text: '这是什么规律？', action: 'question' }
            ]
          }
        ],
        discovery: [
          {
            text: '看这里！我发现直角三角形两条直角边的平方和...',
            action: 'revelation'
          },
          {
            text: '等于斜边的平方！这就是著名的毕达哥拉斯定理！',
            action: 'theorem'
          },
          {
            text: '让我们用实际测量来验证这个发现吧！',
            action: 'verify',
            options: [
              { text: '开始测量', action: 'measure' },
              { text: '先理解一下', action: 'understand' }
            ]
          }
        ]
      }
    }
  };

  const currentNPC = npcConfig.pythagoras;
  const dialogs = currentNPC.dialogs[scene] || [];
  const currentDialogData = dialogs[currentDialog];

  // 处理对话进度
  useEffect(() => {
    if (dialogs.length > 0 && currentDialog < dialogs.length) {
      setIsTalking(true);
      
      // 模拟说话动画
      const timer = setTimeout(() => {
        setIsTalking(false);
        if (currentDialogData.options) {
          setShowOptions(true);
        }
      }, 2000);
      
      return () => clearTimeout(timer);
    }
  }, [currentDialog, dialogs, currentDialogData]);

  // 处理选项选择
  const handleOptionSelect = (option) => {
    setShowOptions(false);
    
    // 根据选项执行不同动作
    switch (option.action) {
      case 'accept':
      case 'explore':
      case 'measure':
        // 触发任务开始
        onTaskStart('pythagorean_measurement');
        // 移动到下一个对话
        if (currentDialog < dialogs.length - 1) {
          setTimeout(() => setCurrentDialog(prev => prev + 1), 500);
        }
        break;
        
      case 'wait':
      case 'understand':
        // 等待后继续对话
        setTimeout(() => setCurrentDialog(prev => prev + 1), 1000);
        break;
        
      case 'question':
        // 特殊处理：提供更多解释
        setTimeout(() => setCurrentDialog(prev => prev + 1), 500);
        break;
        
      default:
        if (currentDialog < dialogs.length - 1) {
          setCurrentDialog(prev => prev + 1);
        }
        break;
    }
  };

  // 重置对话
  const resetDialog = () => {
    setCurrentDialog(0);
    setShowOptions(false);
  };

  // 渲染NPC形象
  const renderNPC = () => {
    return (
      <div style={{
        position: 'absolute',
        bottom: '100px',
        right: '100px',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        zIndex: 50
      }}>
        {/* NPC头像 */}
        <div style={{
          fontSize: '80px',
          marginBottom: '10px',
          animation: isTalking ? 'bounce 0.5s infinite' : 'none'
        }}>
          {currentNPC.avatar}
        </div>
        
        {/* NPC名字标签 */}
        <div style={{
          padding: '6px 12px',
          backgroundColor: currentNPC.color,
          color: 'white',
          borderRadius: '20px',
          fontSize: '14px',
          fontWeight: '500',
          boxShadow: '0 2px 4px rgba(0, 0, 0, 0.2)'
        }}>
          {currentNPC.name}
        </div>
      </div>
    );
  };

  // 渲染对话框
  const renderDialog = () => {
    if (!currentDialogData) return null;
    
    return (
      <div style={{
        position: 'absolute',
        bottom: '250px',
        right: '200px',
        maxWidth: '350px',
        zIndex: 60
      }}>
        <div style={{
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '20px',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
          border: `3px solid ${currentNPC.color}`,
          position: 'relative'
        }}>
          {/* 对话文本 */}
          <div style={{
            fontSize: '16px',
            lineHeight: '1.5',
            color: '#374151',
            marginBottom: showOptions ? '15px' : '0'
          }}>
            {currentDialogData.text}
            {isTalking && (
              <span style={{ 
                display: 'inline-block',
                marginLeft: '5px',
                animation: 'blink 1s infinite'
              }}>
                ▋
              </span>
            )}
          </div>
          
          {/* 选项按钮 */}
          {showOptions && currentDialogData.options && (
            <div style={{
              display: 'flex',
              flexDirection: 'column',
              gap: '8px'
            }}>
              {currentDialogData.options.map((option, index) => (
                <button
                  key={index}
                  onClick={() => handleOptionSelect(option)}
                  style={{
                    padding: '10px 16px',
                    backgroundColor: 'rgba(139, 92, 246, 0.1)',
                    border: `2px solid ${currentNPC.color}`,
                    borderRadius: '8px',
                    color: currentNPC.color,
                    fontSize: '14px',
                    fontWeight: '500',
                    cursor: 'pointer',
                    transition: 'all 0.2s ease'
                  }}
                  onMouseEnter={(e) => {
                    e.target.style.backgroundColor = currentNPC.color;
                    e.target.style.color = 'white';
                  }}
                  onMouseLeave={(e) => {
                    e.target.style.backgroundColor = 'rgba(139, 92, 246, 0.1)';
                    e.target.style.color = currentNPC.color;
                  }}
                >
                  {option.text}
                </button>
              ))}
            </div>
          )}
          
          {/* 对话框尾巴 */}
          <div style={{
            position: 'absolute',
            bottom: '-12px',
            right: '30px',
            width: '0',
            height: '0',
            borderLeft: '12px solid transparent',
            borderRight: '12px solid transparent',
            borderTop: `12px solid white`
          }}></div>
        </div>
      </div>
    );
  };

  // 渲染进度指示器
  const renderProgress = () => {
    return (
      <div style={{
        position: 'absolute',
        top: '20px',
        right: '20px',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        zIndex: 40
      }}>
        <span style={{ fontSize: '14px', color: '#6b7280' }}>对话进度:</span>
        <div style={{
          display: 'flex',
          gap: '4px'
        }}>
          {dialogs.map((_, index) => (
            <div
              key={index}
              style={{
                width: '8px',
                height: '8px',
                borderRadius: '50%',
                backgroundColor: index <= currentDialog ? currentNPC.color : '#d1d5db'
              }}
            ></div>
          ))}
        </div>
      </div>
    );
  };

  // 如果没有对话数据，显示简单的NPC
  if (dialogs.length === 0) {
    return (
      <div style={{
        position: 'absolute',
        bottom: '100px',
        right: '100px',
        fontSize: '60px',
        zIndex: 50
      }}>
        {currentNPC.avatar}
      </div>
    );
  }

  return (
    <div>
      {renderNPC()}
      {renderDialog()}
      {renderProgress()}
      
      {/* CSS动画 */}
      <style>{`
        @keyframes bounce {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-10px); }
        }
        
        @keyframes blink {
          0%, 100% { opacity: 1; }
          50% { opacity: 0; }
        }
      `}</style>
    </div>
  );
};

export default CharacterInteraction;