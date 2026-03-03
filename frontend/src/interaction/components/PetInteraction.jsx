import React, { useState, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';
import SoundManager from '../../../../shared/utils/SoundManager';

const Pet = ({ type = 'basic', position = [0, 0, 0] }) => {
  const meshRef = useRef();
  const [isHappy, setIsHappy] = useState(false);
  const [color, setColor] = useState('#FFA500'); // Orange
  const [size, setSize] = useState(1);
  const [rotationSpeed, setRotationSpeed] = useState(0);
  const [clickCount, setClickCount] = useState(0);

  const handleClick = () => {
    SoundManager.playSound('pet_clicked');
    setIsHappy(true);
    setSize(1.5);
    setRotationSpeed(0.05);
    
    // 随机改变颜色
    const colors = ['#FFA500', '#ADFF2F', '#FF69B4', '#4169E1', '#32CD32'];
    setColor(colors[Math.floor(Math.random() * colors.length)]);
    
    setClickCount(prev => prev + 1);
    
    // 重置动画
    setTimeout(() => {
      setIsHappy(false);
      setSize(1);
      setRotationSpeed(0);
      setColor('#FFA500');
    }, 1000);
  };

  useFrame(() => {
    if (meshRef.current) {
      // 旋转效果
      meshRef.current.rotation.y += rotationSpeed;
      
      // 呼吸效果
      if (isHappy) {
        meshRef.current.scale.x = size + Math.sin(Date.now() / 200) * 0.1;
        meshRef.current.scale.y = size + Math.sin(Date.now() / 200) * 0.1;
        meshRef.current.scale.z = size + Math.sin(Date.now() / 200) * 0.1;
      } else {
        meshRef.current.scale.setScalar(size);
      }
    }
  });

  // 根据类型设置不同几何形状
  let geometry;
  switch(type) {
    case 'cat':
      geometry = <sphereGeometry args={[0.5, 16, 16]} />;
      break;
    case 'dog':
      geometry = <boxGeometry args={[0.6, 0.5, 0.4]} />;
      break;
    case 'bird':
      geometry = <coneGeometry args={[0.3, 0.8, 4]} />;
      break;
    case 'robot':
      geometry = <boxGeometry args={[0.5, 0.7, 0.5]} />;
      break;
    default:
      geometry = <sphereGeometry args={[0.5, 16, 16]} />;
  }

  return (
    <mesh 
      ref={meshRef} 
      position={position} 
      onClick={handleClick}
      onPointerOver={() => document.body.style.cursor = 'pointer'}
      onPointerOut={() => document.body.style.cursor = 'default'}
    >
      {geometry}
      <meshStandardMaterial 
        color={color} 
        roughness={0.2}
        metalness={0.8}
        emissive={isHappy ? color : "#000000"}
        emissiveIntensity={isHappy ? 0.2 : 0}
      />
    </mesh>
  );
};

const PetInteraction = () => {
  const [activePets, setActivePets] = useState([
    { id: 1, type: 'cat', position: [-2, 0.5, 0] },
    { id: 2, type: 'dog', position: [0, 0.5, 0] },
    { id: 3, type: 'bird', position: [2, 0.5, 0] },
    { id: 4, type: 'robot', position: [0, 0.5, -2] },
  ]);

  const [clickStats, setClickStats] = useState({});

  const handlePetClick = (id) => {
    setClickStats(prev => ({
      ...prev,
      [id]: (prev[id] || 0) + 1
    }));
  };

  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <div style={{ position: 'absolute', top: 10, left: 10, zIndex: 100, color: 'white', backgroundColor: 'rgba(0,0,0,0.5)', padding: '10px', borderRadius: '5px' }}>
        <h3>宠物互动乐园</h3>
        <p>点击宠物与它们互动！</p>
        <p>点击统计: {Object.values(clickStats).reduce((a, b) => a + b, 0)} 次</p>
      </div>
      
      <Canvas 
        camera={{ position: [5, 5, 5], fov: 50 }} 
        shadows
      >
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} castShadow />
        <spotLight position={[5, 10, 5]} angle={0.15} penumbra={1} intensity={0.5} castShadow />
        
        {/* 地面 */}
        <mesh 
          rotation={[-Math.PI / 2, 0, 0]} 
          receiveShadow
        >
          <planeGeometry args={[20, 20]} />
          <meshStandardMaterial 
            color="#90EE90" 
            roughness={0.9} 
            metalness={0.1} 
          />
        </mesh>
        
        {/* 装饰元素 */}
        {/* 树 */}
        <group position={[4, 0, 4]}>
          <mesh position={[0, 3, 0]}>
            <sphereGeometry args={[1.5, 16, 16]} />
            <meshStandardMaterial color="#2E8B57" />
          </mesh>
          <mesh position={[0, 1, 0]}>
            <cylinderGeometry args={[0.3, 0.3, 2, 8]} />
            <meshStandardMaterial color="#8B4513" />
          </mesh>
        </group>
        
        <group position={[-4, 0, -4]}>
          <mesh position={[0, 2.5, 0]}>
            <sphereGeometry args={[1.2, 16, 16]} />
            <meshStandardMaterial color="#3CB371" />
          </mesh>
          <mesh position={[0, 1, 0]}>
            <cylinderGeometry args={[0.25, 0.25, 2, 8]} />
            <meshStandardMaterial color="#8B4513" />
          </mesh>
        </group>
        
        {/* 池塘 */}
        <mesh position={[-3, 0.01, 3]}>
          <circleGeometry args={[1.5, 32]} />
          <meshStandardMaterial 
            color="#4682B4" 
            roughness={0.1} 
            metalness={0.9} 
            side={2} 
          />
        </mesh>
        
        {/* 花丛 */}
        <group position={[3, 0.2, -3]}>
          <mesh position={[-0.3, 0, 0]}>
            <sphereGeometry args={[0.15, 8, 8]} />
            <meshStandardMaterial color="#FF69B4" />
          </mesh>
          <mesh position={[0.3, 0.1, 0]}>
            <sphereGeometry args={[0.12, 8, 8]} />
            <meshStandardMaterial color="#FF1493" />
          </mesh>
          <mesh position={[0, 0.2, -0.3]}>
            <sphereGeometry args={[0.14, 8, 8]} />
            <meshStandardMaterial color="#FF6347" />
          </mesh>
        </group>
        
        {/* 宠物们 */}
        {activePets.map(pet => (
          <Pet 
            key={pet.id} 
            type={pet.type} 
            position={pet.position}
            onClick={() => handlePetClick(pet.id)}
          />
        ))}
        
        {/* 添加相机控制器 */}
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

export default PetInteraction;