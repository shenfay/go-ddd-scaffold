import React, { useRef, useState, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Text } from '@react-three/drei';
import * as THREE from 'three';
import SoundManager from '../../../../shared/utils/SoundManager';

// 一个会演变的几何体
function EvolvingShape({ isActive, evolutionStage }) {
  const meshRef = useRef();
  const [shapeType, setShapeType] = useState(0); // 0=sphere, 1=box, 2=cone, 3=dodecahedron
  
  useFrame((state, delta) => {
    if (isActive && meshRef.current) {
      // 旋转动画
      meshRef.current.rotation.x += delta * 0.5;
      meshRef.current.rotation.y += delta * 0.7;
      
      // 颜色变化
      const material = meshRef.current.material;
      const hue = (state.clock.elapsedTime * 0.1) % 1;
      material.color.setHSL(hue, 1, 0.5);
      
      // 规模演变
      const scale = 1 + Math.sin(state.clock.elapsedTime * 2) * 0.1;
      meshRef.current.scale.setScalar(scale);
    }
  });

  // 根据进化阶段选择形状
  const geometry = useMemo(() => {
    switch(evolutionStage) {
      case 1:
        return new THREE.SphereGeometry(1, 32, 32);
      case 2:
        return new THREE.BoxGeometry(1.5, 1.5, 1.5);
      case 3:
        return new THREE.ConeGeometry(1, 1.5, 8);
      case 4:
        return new THREE.DodecahedronGeometry(1, 0);
      default:
        return new THREE.TorusGeometry(1, 0.4, 16, 100);
    }
  }, [evolutionStage]);

  return (
    <mesh ref={meshRef}>
      <primitive object={geometry} attach="geometry" />
      <meshStandardMaterial 
        color="orange" 
        wireframe={isActive} 
        roughness={0.1}
        metalness={0.9}
        emissive={isActive ? "#FFA500" : "#000000"}
        emissiveIntensity={isActive ? 0.5 : 0}
      />
    </mesh>
  );
}

const EvolutionAnimation = () => {
  const [evolving, setEvolving] = useState(false);
  const [evolutionStage, setEvolutionStage] = useState(0);
  const [animationSpeed, setAnimationSpeed] = useState(1);
  const [showInstructions, setShowInstructions] = useState(true);

  const triggerEvolution = () => {
    SoundManager.playSound('evolution_start');
    setEvolving(true);
    
    // 循环演变阶段
    let stage = 1;
    const interval = setInterval(() => {
      setEvolutionStage(stage);
      stage++;
      if (stage > 4) {
        clearInterval(interval);
        setTimeout(() => {
          setEvolving(false);
          setEvolutionStage(0);
        }, 1000);
      }
    }, 1000);
  };

  const increaseSpeed = () => {
    setAnimationSpeed(prev => Math.min(prev + 0.2, 3));
    SoundManager.playSound('speed_increase');
  };

  const decreaseSpeed = () => {
    setAnimationSpeed(prev => Math.max(prev - 0.2, 0.5));
    SoundManager.playSound('speed_decrease');
  };

  const resetAnimation = () => {
    setEvolving(false);
    setEvolutionStage(0);
    setAnimationSpeed(1);
    SoundManager.playSound('reset');
  };

  return (
    <div style={{ width: '100vw', height: '100vh', position: 'relative' }}>
      {showInstructions && (
        <div 
          style={{ 
            position: 'absolute', 
            top: '10px', 
            left: '10px', 
            zIndex: 100, 
            backgroundColor: 'rgba(0,0,0,0.7)', 
            color: 'white', 
            padding: '15px', 
            borderRadius: '10px',
            maxWidth: '300px'
          }}
        >
          <h3>演进动画演示</h3>
          <p>观察几何体的演变过程！</p>
          <button 
            onClick={() => setShowInstructions(false)} 
            style={{ 
              marginTop: '10px', 
              padding: '5px 10px', 
              backgroundColor: '#4682B4', 
              color: 'white', 
              border: 'none', 
              borderRadius: '4px',
              cursor: 'pointer'
            }}
          >
            关闭说明
          </button>
        </div>
      )}
      
      <div style={{ position: 'absolute', top: '10px', right: '10px', zIndex: 100, color: 'white', backgroundColor: 'rgba(0,0,0,0.5)', padding: '10px', borderRadius: '5px' }}>
        <p>演变阶段: {evolutionStage}/4</p>
        <p>动画速度: {animationSpeed.toFixed(1)}x</p>
      </div>
      
      <div style={{ position: 'absolute', bottom: '20px', left: '50%', transform: 'translateX(-50%)', zIndex: 100 }}>
        <button 
          onClick={triggerEvolution} 
          disabled={evolving}
          style={{ 
            marginRight: '10px',
            padding: '10px 20px', 
            backgroundColor: evolving ? '#cccccc' : '#4682B4', 
            color: 'white', 
            border: 'none', 
            borderRadius: '4px', 
            cursor: evolving ? 'not-allowed' : 'pointer'
          }}
        >
          {evolving ? '演变中...' : '触发演进动画'}
        </button>
        <button 
          onClick={increaseSpeed} 
          style={{ 
            marginRight: '10px',
            padding: '10px 20px', 
            backgroundColor: '#32CD32', 
            color: 'white', 
            border: 'none', 
            borderRadius: '4px', 
            cursor: 'pointer'
          }}
        >
          加速
        </button>
        <button 
          onClick={decreaseSpeed} 
          style={{ 
            marginRight: '10px',
            padding: '10px 20px', 
            backgroundColor: '#FF6347', 
            color: 'white', 
            border: 'none', 
            borderRadius: '4px', 
            cursor: 'pointer'
          }}
        >
          减速
        </button>
        <button 
          onClick={resetAnimation} 
          style={{ 
            padding: '10px 20px', 
            backgroundColor: '#666', 
            color: 'white', 
            border: 'none', 
            borderRadius: '4px', 
            cursor: 'pointer'
          }}
        >
          重置
        </button>
      </div>
      
      <Canvas 
        camera={{ position: [5, 5, 5], fov: 50 }} 
        shadows
      >
        <ambientLight intensity={0.4} />
        <pointLight position={[10, 10, 10]} intensity={1} castShadow />
        <spotLight position={[5, 10, 5]} angle={0.15} penumbra={1} intensity={0.5} castShadow />
        
        {/* 添加背景网格 */}
        <gridHelper args={[20, 20, '#444444', '#222222']} />
        
        {/* 添加彩色灯光效果 */}
        <pointLight position={[5, 10, 5]} color="#FF6347" intensity={0.5} />
        <pointLight position={[-5, 10, -5]} color="#4169E1" intensity={0.5} />
        <pointLight position={[0, 10, -7]} color="#32CD32" intensity={0.5} />
        
        {/* 演变中的形状 */}
        <EvolvingShape isActive={evolving} evolutionStage={evolutionStage} />
        
        {/* 添加标签 */}
        <Text
          position={[0, -3, 0]}
          fontSize={0.5}
          color="white"
          anchorX="center"
          anchorY="middle"
        >
          {evolutionStage === 0 ? '准备演变...' : 
           evolutionStage === 1 ? '球形阶段' : 
           evolutionStage === 2 ? '立方体阶段' : 
           evolutionStage === 3 ? '锥形阶段' : 
           '多面体阶段'}
        </Text>
        
        {/* 添加轨道控制器 */}
        <OrbitControls 
          enablePan={true} 
          enableZoom={true} 
          enableRotate={true}
          minDistance={3}
          maxDistance={15}
        />
      </Canvas>
    </div>
  );
};

export default EvolutionAnimation;