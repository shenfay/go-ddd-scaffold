/**
 * 交互任务组件
 * 
 * 处理核心的交互式学习任务，如三角形测量
 */

import React, { useState, useEffect, useRef } from 'react';

const InteractiveTask = ({ currentTask, onComplete }) => {
  const [taskState, setTaskState] = useState({
    isActive: false,
    step: 0,
    measurements: [],
    currentTriangle: 0,
    isMeasuring: false,
    showResult: false
  });
  
  const canvasRef = useRef(null);

  // 任务配置
  const taskConfig = {
    pythagorean_measurement: {
      name: '毕达哥拉斯定理验证',
      description: '通过测量不同的直角三角形来发现边长之间的关系',
      triangles: [
        { a: 3, b: 4, c: 5, color: '#3b82f6' },
        { a: 5, b: 12, c: 13, color: '#10b981' },
        { a: 8, b: 15, c: 17, color: '#f59e0b' }
      ],
      steps: [
        '观察直角三角形',
        '测量两条直角边',
        '测量斜边',
        '计算并验证关系'
      ]
    }
  };

  const currentTaskConfig = taskConfig[currentTask] || taskConfig.pythagorean_measurement;
  const currentTriangle = currentTaskConfig.triangles[taskState.currentTriangle];

  // 初始化任务
  useEffect(() => {
    if (currentTask) {
      setTaskState(prev => ({
        ...prev,
        isActive: true,
        step: 0,
        measurements: [],
        currentTriangle: 0
      }));
    } else {
      setTaskState(prev => ({ ...prev, isActive: false }));
    }
  }, [currentTask]);

  // 绘制三角形
  useEffect(() => {
    if (!canvasRef.current || !taskState.isActive) return;
    
    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    const width = canvas.width;
    const height = canvas.height;
    
    // 清空画布
    ctx.clearRect(0, 0, width, height);
    
    if (!currentTriangle) return;
    
    // 计算三角形坐标（简化绘制）
    const scale = 20;
    const offsetX = width / 2 - (currentTriangle.a * scale) / 2;
    const offsetY = height / 2 - (currentTriangle.b * scale) / 2;
    
    ctx.strokeStyle = currentTriangle.color;
    ctx.lineWidth = 3;
    ctx.fillStyle = currentTriangle.color + '20';
    
    // 绘制三角形
    ctx.beginPath();
    ctx.moveTo(offsetX, offsetY); // 直角顶点
    ctx.lineTo(offsetX + currentTriangle.a * scale, offsetY); // 底边
    ctx.lineTo(offsetX, offsetY + currentTriangle.b * scale); // 高边
    ctx.closePath();
    ctx.fill();
    ctx.stroke();
    
    // 绘制边长标签
    ctx.fillStyle = '#1f2937';
    ctx.font = '14px Arial';
    ctx.fillText(`a = ${currentTriangle.a}`, offsetX + (currentTriangle.a * scale) / 2, offsetY - 10);
    ctx.fillText(`b = ${currentTriangle.b}`, offsetX - 30, offsetY + (currentTriangle.b * scale) / 2);
    ctx.fillText(`c = ${currentTriangle.c}`, offsetX + (currentTriangle.a * scale) / 2 + 10, offsetY + (currentTriangle.b * scale) / 2 + 20);
    
    // 绘制直角标记
    ctx.strokeStyle = '#ef4444';
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.moveTo(offsetX + 15, offsetY);
    ctx.lineTo(offsetX + 15, offsetY + 15);
    ctx.lineTo(offsetX, offsetY + 15);
    ctx.stroke();
  }, [taskState, currentTriangle]);

  // 处理测量操作
  const handleMeasure = (side) => {
    if (taskState.isMeasuring) return;
    
    setTaskState(prev => ({
      ...prev,
      isMeasuring: true
    }));
    
    // 模拟测量过程
    setTimeout(() => {
      const measurement = {
        triangle: taskState.currentTriangle,
        side: side,
        value: currentTriangle[side],
        timestamp: Date.now()
      };
      
      setTaskState(prev => ({
        ...prev,
        measurements: [...prev.measurements, measurement],
        isMeasuring: false,
        step: Math.min(prev.step + 1, 3)
      }));
    }, 1000);
  };

  // 处理下一步
  const handleNext = () => {
    if (taskState.step < 3) {
      setTaskState(prev => ({
        ...prev,
        step: prev.step + 1
      }));
    } else {
      // 完成当前三角形
      if (taskState.currentTriangle < currentTaskConfig.triangles.length - 1) {
        setTaskState(prev => ({
          ...prev,
          currentTriangle: prev.currentTriangle + 1,
          step: 0,
          measurements: []
        }));
      } else {
        // 所有三角形完成，显示结果
        setTaskState(prev => ({
          ...prev,
          showResult: true
        }));
      }
    }
  };

  // 验证毕达哥拉斯定理
  const verifyPythagoreanTheorem = () => {
    const a = currentTriangle.a;
    const b = currentTriangle.b;
    const c = currentTriangle.c;
    
    const leftSide = a * a + b * b;
    const rightSide = c * c;
    
    return {
      a: a,
      b: b,
      c: c,
      leftSide: leftSide,
      rightSide: rightSide,
      isCorrect: Math.abs(leftSide - rightSide) < 0.001,
      formula: `${a}² + ${b}² = ${leftSide} ${leftSide === rightSide ? '=' : '≠'} ${c}² = ${rightSide}`
    };
  };

  // 完成任务
  const handleComplete = () => {
    const result = {
      taskId: currentTask,
      triangles: currentTaskConfig.triangles.map((triangle, index) => ({
        ...triangle,
        verification: verifyPythagoreanTheorem()
      })),
      completedAt: Date.now()
    };
    
    onComplete(currentTask, result);
    setTaskState({
      isActive: false,
      step: 0,
      measurements: [],
      currentTriangle: 0,
      isMeasuring: false,
      showResult: false
    });
  };

  // 渲染测量界面
  const renderMeasurementInterface = () => {
    const steps = [
      { title: '观察', icon: '👀', description: '仔细观察直角三角形的结构' },
      { title: '测量', icon: '📏', description: '测量两条直角边的长度' },
      { title: '测量', icon: '📐', description: '测量斜边的长度' },
      { title: '验证', icon: '✅', description: '验证边长之间的数学关系' }
    ];
    
    return (
      <div style={{
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        width: '80%',
        maxWidth: '800px',
        backgroundColor: 'white',
        borderRadius: '16px',
        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.1)',
        padding: '30px',
        zIndex: 100
      }}>
        {/* 任务标题 */}
        <div style={{
          textAlign: 'center',
          marginBottom: '25px'
        }}>
          <h2 style={{
            margin: '0 0 10px 0',
            fontSize: '24px',
            fontWeight: '600',
            color: '#1f2937'
          }}>
            {currentTaskConfig.name}
          </h2>
          <p style={{
            margin: 0,
            fontSize: '16px',
            color: '#6b7280'
          }}>
            {currentTaskConfig.description}
          </p>
        </div>
        
        {/* 进度条 */}
        <div style={{
          marginBottom: '25px'
        }}>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            marginBottom: '10px'
          }}>
            {steps.map((step, index) => (
              <div
                key={index}
                style={{
                  textAlign: 'center',
                  flex: 1
                }}
              >
                <div style={{
                  width: '40px',
                  height: '40px',
                  borderRadius: '50%',
                  backgroundColor: index <= taskState.step ? '#3b82f6' : '#d1d5db',
                  color: 'white',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '20px',
                  margin: '0 auto 8px',
                  transition: 'all 0.3s ease'
                }}>
                  {index < taskState.step ? '✓' : step.icon}
                </div>
                <div style={{
                  fontSize: '12px',
                  color: index <= taskState.step ? '#3b82f6' : '#9ca3af',
                  fontWeight: index === taskState.step ? '500' : 'normal'
                }}>
                  {step.title}
                </div>
              </div>
            ))}
          </div>
          <div style={{
            fontSize: '14px',
            color: '#6b7280',
            textAlign: 'center',
            marginTop: '10px'
          }}>
            {steps[taskState.step].description}
          </div>
        </div>
        
        {/* 三角形画布 */}
        <div style={{
          display: 'flex',
          justifyContent: 'center',
          marginBottom: '25px'
        }}>
          <canvas
            ref={canvasRef}
            width={400}
            height={300}
            style={{
              border: '2px solid #e5e7eb',
              borderRadius: '8px',
              cursor: taskState.isMeasuring ? 'wait' : 'pointer'
            }}
          />
        </div>
        
        {/* 操作按钮 */}
        <div style={{
          display: 'flex',
          justifyContent: 'center',
          gap: '15px',
          marginBottom: '20px'
        }}>
          {taskState.step === 1 && (
            <button
              onClick={() => handleMeasure('a')}
              disabled={taskState.isMeasuring}
              style={{
                padding: '12px 24px',
                backgroundColor: '#3b82f6',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '16px',
                cursor: taskState.isMeasuring ? 'not-allowed' : 'pointer',
                opacity: taskState.isMeasuring ? 0.6 : 1,
                transition: 'all 0.2s ease'
              }}
            >
              {taskState.isMeasuring ? '测量中...' : `测量边 a (${currentTriangle.a})`}
            </button>
          )}
          
          {taskState.step === 2 && (
            <button
              onClick={() => handleMeasure('b')}
              disabled={taskState.isMeasuring}
              style={{
                padding: '12px 24px',
                backgroundColor: '#10b981',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '16px',
                cursor: taskState.isMeasuring ? 'not-allowed' : 'pointer',
                opacity: taskState.isMeasuring ? 0.6 : 1
              }}
            >
              {taskState.isMeasuring ? '测量中...' : `测量边 b (${currentTriangle.b})`}
            </button>
          )}
          
          {taskState.step === 3 && (
            <button
              onClick={() => handleMeasure('c')}
              disabled={taskState.isMeasuring}
              style={{
                padding: '12px 24px',
                backgroundColor: '#f59e0b',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '16px',
                cursor: taskState.isMeasuring ? 'not-allowed' : 'pointer',
                opacity: taskState.isMeasuring ? 0.6 : 1
              }}
            >
              {taskState.isMeasuring ? '测量中...' : `测量边 c (${currentTriangle.c})`}
            </button>
          )}
          
          {taskState.step > 0 && taskState.step < 3 && (
            <button
              onClick={handleNext}
              style={{
                padding: '12px 24px',
                backgroundColor: '#8b5cf6',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '16px',
                cursor: 'pointer'
              }}
            >
              下一步
            </button>
          )}
        </div>
        
        {/* 测量结果显示 */}
        {taskState.measurements.length > 0 && (
          <div style={{
            backgroundColor: '#f9fafb',
            borderRadius: '8px',
            padding: '15px',
            marginBottom: '20px'
          }}>
            <h4 style={{ margin: '0 0 10px 0', color: '#374151' }}>测量结果:</h4>
            <div style={{ display: 'flex', gap: '15px', flexWrap: 'wrap' }}>
              {taskState.measurements.map((measurement, index) => (
                <div
                  key={index}
                  style={{
                    backgroundColor: 'white',
                    padding: '8px 12px',
                    borderRadius: '6px',
                    border: `2px solid ${currentTriangle.color}`,
                    fontSize: '14px'
                  }}
                >
                  边{measurement.side} = {measurement.value}
                </div>
              ))}
            </div>
          </div>
        )}
        
        {/* 完成按钮 */}
        {taskState.showResult && (
          <div style={{ textAlign: 'center' }}>
            <div style={{
              backgroundColor: '#dcfce7',
              border: '2px solid #22c55e',
              borderRadius: '8px',
              padding: '20px',
              marginBottom: '20px'
            }}>
              <div style={{ fontSize: '24px', marginBottom: '10px' }}>🎉</div>
              <h3 style={{ margin: '0 0 10px 0', color: '#166534' }}>发现成功！</h3>
              <p style={{ margin: 0, color: '#166534' }}>
                你验证了毕达哥拉斯定理：a² + b² = c²
              </p>
            </div>
            <button
              onClick={handleComplete}
              style={{
                padding: '14px 28px',
                backgroundColor: '#22c55e',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                fontSize: '18px',
                fontWeight: '500',
                cursor: 'pointer',
                boxShadow: '0 4px 6px rgba(34, 197, 94, 0.25)'
              }}
            >
              完成任务 🎓
            </button>
          </div>
        )}
      </div>
    );
  };

  if (!taskState.isActive) {
    return null;
  }

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      width: '100vw',
      height: '100vh',
      backgroundColor: 'rgba(0, 0, 0, 0.5)',
      zIndex: 3000,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center'
    }}>
      {renderMeasurementInterface()}
    </div>
  );
};

export default InteractiveTask;