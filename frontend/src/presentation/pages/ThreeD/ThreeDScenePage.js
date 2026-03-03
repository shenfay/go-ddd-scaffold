import React, { Suspense, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

// 3D 组件
import KnowledgeTown from '../../../interaction/three/components/scenes/KnowledgeTown';
import KnowledgeVillage from '../../../interaction/three/components/scenes/KnowledgeVillage';

const sceneConfigs = {
  town: {
    name: '知识小镇',
    description: '探索充满数学魅力的奇妙小镇',
    bgGradient: 'linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%)',
  },
  village: {
    name: '知识村庄',
    description: '发现数学世界的宁静村落',
    bgGradient: 'linear-gradient(135deg, #134e5e 0%, #71b280 100%)',
  },
};

const ThreeDScenePage = () => {
  const { sceneType } = useParams();
  const navigate = useNavigate();
  const config = sceneConfigs[sceneType] || sceneConfigs.town;
  const [error, setError] = useState(null);
  
  const handleBack = () => {
    navigate('/knowledge-map');
  };

  const renderScene = () => {
    switch (sceneType) {
      case 'village':
        return <KnowledgeVillage />;
      case 'town':
      default:
        return <KnowledgeTown />;
    }
  };

  return (
    <div style={{
      width: '100vw',
      height: '100vh',
      position: 'fixed',
      top: 0,
      left: 0,
      background: config.bgGradient,
      overflow: 'hidden',
    }}>
      {/* 返回按钮 */}
      <button
        onClick={handleBack}
        style={{
          position: 'fixed',
          top: '20px',
          left: '20px',
          zIndex: 100,
          padding: '12px 20px',
          backgroundColor: 'rgba(255,255,255,0.15)',
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(255,255,255,0.2)',
          borderRadius: '12px',
          color: 'white',
          fontSize: '14px',
          cursor: 'pointer',
          display: 'flex',
          alignItems: 'center',
          gap: '8px',
          transition: '0.2s',
        }}
      >
        ← 返回知识地图
      </button>

      {/* 场景标题 */}
      <div style={{
        position: 'fixed',
        top: '20px',
        right: '20px',
        zIndex: 100,
        padding: '12px 20px',
        backgroundColor: 'rgba(255,255,255,0.15)',
        backdropFilter: 'blur(10px)',
        border: '1px solid rgba(255,255,255,0.2)',
        borderRadius: '12px',
        color: 'white',
        textAlign: 'right',
      }}>
        <div style={{ fontSize: '18px', fontWeight: 'bold' }}>{config.name}</div>
        <div style={{ fontSize: '12px', opacity: 0.8 }}>{config.description}</div>
      </div>

      {/* 底部提示 */}
      <div style={{
        position: 'fixed',
        bottom: '20px',
        left: '50%',
        transform: 'translateX(-50%)',
        zIndex: 100,
        padding: '10px 20px',
        backgroundColor: 'rgba(0,0,0,0.4)',
        backdropFilter: 'blur(10px)',
        borderRadius: '20px',
        color: 'white',
        fontSize: '12px',
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
      }}>
        🖱️ 拖拽旋转 · 滚轮缩放 · 点击节点学习
      </div>

      {/* 3D 场景 */}
      <Suspense fallback={
        <div style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100%',
          color: 'white',
          fontSize: '18px',
        }}>
          加载 3D 场景中...
        </div>
      }>
        <div style={{ width: '100%', height: '100%' }}>
          {renderScene()}
        </div>
      </Suspense>
    </div>
  );
};

export default ThreeDScenePage;
