import React from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Text, Sky, Environment, Stars } from '@react-three/drei';
import SoundManager from '../../../../shared/utils/SoundManager';

// 小木屋 (来自上一个示例)
const Cabin = ({ position }) => {
  return (
    <group position={position}>
      <mesh position={[0, 1, 0]}>
        <boxGeometry args={[2, 2, 2]} />
        <meshStandardMaterial color="#8B4513" roughness={0.7} metalness={0.3} />
      </mesh>
      <mesh position={[0, 2.2, 0]} rotation={[0, 0, Math.PI / 4]}>
        <coneGeometry args={[1.8, 1.5, 4]} />
        <meshStandardMaterial color="#A52A2A" roughness={0.5} metalness={0.2} />
      </mesh>
      <mesh position={[0, 0, 1.01]}>
        <planeGeometry args={[0.8, 1.2]} />
        <meshStandardMaterial color="#654321" roughness={0.8} metalness={0.1} />
      </mesh>
    </group>
  );
};

// 砖瓦房组件
const BrickHouse = ({ position, color = "#CD853F" }) => {
  return (
    <group position={position}>
      {/* 房身 */}
      <mesh position={[0, 1.5, 0]}>
        <boxGeometry args={[2.5, 3, 2.5]} />
        <meshStandardMaterial color={color} roughness={0.6} metalness={0.4} />
      </mesh>
      {/* 屋顶 */}
      <mesh position={[0, 3.25, 0]} rotation={[0, 0, Math.PI / 4]}>
        <coneGeometry args={[2, 1.5, 4]} />
        <meshStandardMaterial color="#800000" roughness={0.5} metalness={0.2} />
      </mesh>
      {/* 门 */}
      <mesh position={[0, 0.75, 1.26]}>
        <planeGeometry args={[0.8, 1.5]} />
        <meshStandardMaterial color="#654321" roughness={0.8} metalness={0.1} />
      </mesh>
      {/* 窗户 */}
      <mesh position={[-0.8, 1.5, 1.26]}>
        <planeGeometry args={[0.5, 0.5]} />
        <meshStandardMaterial 
          color="#87CEEB" 
          transparent 
          opacity={0.7} 
          emissive="#87CEEB"
          emissiveIntensity={0.2}
        />
      </mesh>
      <mesh position={[0.8, 1.5, 1.26]}>
        <planeGeometry args={[0.5, 0.5]} />
        <meshStandardMaterial 
          color="#87CEEB" 
          transparent 
          opacity={0.7} 
          emissive="#87CEEB"
          emissiveIntensity={0.2}
        />
      </mesh>
    </group>
  );
};

// 图书馆组件
const Library = ({ position }) => {
  return (
    <group position={position}>
      {/* 主体 */}
      <mesh position={[0, 2, 0]}>
        <boxGeometry args={[4, 4, 3]} />
        <meshStandardMaterial color="#D2B48C" roughness={0.7} metalness={0.3} />
      </mesh>
      {/* 屋顶 */}
      <mesh position={[0, 4.2, 0]}>
        <boxGeometry args={[4.2, 0.4, 3.2]} />
        <meshStandardMaterial color="#A0522D" roughness={0.6} metalness={0.4} />
      </mesh>
      {/* 门廊 */}
      <mesh position={[0, 1, 1.6]}>
        <boxGeometry args={[2, 2, 0.2]} />
        <meshStandardMaterial color="#D2B48C" roughness={0.7} metalness={0.3} />
      </mesh>
      {/* 门廊柱子 */}
      <mesh position={[-0.8, 0, 1.6]}>
        <cylinderGeometry args={[0.1, 0.1, 2, 8]} />
        <meshStandardMaterial color="#A0522D" roughness={0.8} metalness={0.2} />
      </mesh>
      <mesh position={[0.8, 0, 1.6]}>
        <cylinderGeometry args={[0.1, 0.1, 2, 8]} />
        <meshStandardMaterial color="#A0522D" roughness={0.8} metalness={0.2} />
      </mesh>
      {/* 门 */}
      <mesh position={[0, 1, 1.61]}>
        <planeGeometry args={[1.5, 1.8]} />
        <meshStandardMaterial color="#654321" roughness={0.8} metalness={0.1} />
      </mesh>
      {/* 书本装饰 (简化) */}
      <mesh position={[-1, 2.5, 0]}>
        <boxGeometry args={[0.2, 0.8, 0.4]} />
        <meshStandardMaterial color="#8B4513" roughness={0.9} metalness={0.1} />
      </mesh>
      <mesh position={[1, 2.5, 0]}>
        <boxGeometry args={[0.2, 0.8, 0.4]} />
        <meshStandardMaterial color="#8B4513" roughness={0.9} metalness={0.1} />
      </mesh>
    </group>
  );
};

// 守护者/宠物组件 (带有多样化外观)
const Guardian = ({ position, type = 'basic' }) => {
  const [hovered, setHovered] = React.useState(false);
  const [clickCount, setClickCount] = React.useState(0);

  const handleClick = () => {
    SoundManager.playSound('town_pet_interaction');
    setClickCount(prev => prev + 1);
    console.log("Pet clicked!", clickCount + 1);
  };

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  // 根据类型设置不同外观
  let geometry, color, scale = 1;
  switch(type) {
    case 'cat':
      geometry = <sphereGeometry args={[0.4, 16, 16]} />;
      color = hovered ? "#FF69B4" : "#FFB6C1"; // 悬停时粉色
      scale = 0.8;
      break;
    case 'dog':
      geometry = <sphereGeometry args={[0.5, 16, 16]} />;
      color = hovered ? "#8B4513" : "#D2691E"; // 悬停时深棕色
      scale = 1;
      break;
    case 'robot':
      geometry = <boxGeometry args={[0.6, 0.6, 0.6]} />;
      color = hovered ? "#708090" : "#A9A9A9"; // 悬停时石板灰
      scale = 1.1;
      break;
    default:
      geometry = <sphereGeometry args={[0.5, 16, 16]} />;
      color = hovered ? "#4169E1" : "#00BFFF"; // 悬停时皇家蓝
      scale = 1;
  }

  return (
    <group position={position} scale={[scale, scale, scale]}>
      <mesh 
        onClick={handleClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        {geometry}
        <meshStandardMaterial 
          color={color}
          roughness={0.3}
          metalness={0.7}
          emissive={hovered ? color : "#000000"}
          emissiveIntensity={hovered ? 0.2 : 0}
        />
        <Text
          position={[0, 1, 0]}
          fontSize={0.4}
          color="black"
          anchorX="center"
          anchorY="middle"
        >
          {type === 'cat' ? '🐱' : type === 'dog' ? '🐶' : type === 'robot' ? '🤖' : '🌟'}
        </Text>
      </mesh>
      
      {/* 添加光环效果 */}
      {hovered && (
        <mesh>
          <torusGeometry args={[0.7, 0.05, 16, 100]} />
          <meshBasicMaterial color="#FFD700" transparent opacity={0.6} />
        </mesh>
      )}
    </group>
  );
};

const KnowledgeTownScene = () => {
  return (
    <div style={{ width: '100vw', height: '100vh' }}>
      <Canvas 
        camera={{ position: [15, 15, 15], fov: 50 }} 
        shadows
      >
        {/* 环境光 */}
        <ambientLight intensity={0.5} />
        
        {/* 主光源 */}
        <directionalLight 
          position={[10, 20, 10]} 
          intensity={1} 
          castShadow
          shadow-mapSize-width={2048}
          shadow-mapSize-height={2048}
        />
        
        {/* 点光源 */}
        <pointLight position={[10, 10, 10]} intensity={0.5} />
        
        {/* 天空环境 */}
        <Sky sunPosition={[100, 10, 100]} turbidity={10} rayleigh={3} mieCoefficient={0.005} mieDirectionalG={0.8} />
        
        {/* 星空背景 */}
        <Stars radius={100} depth={50} count={5000} factor={4} saturation={0} fade />
        
        {/* 背景色 */}
        <color attach="background" args={["#87CEEB"]} />
        
        {/* 地面 */}
        <mesh 
          rotation={[-Math.PI / 2, 0, 0]} 
          receiveShadow
        >
          <planeGeometry args={[50, 50]} />
          <meshStandardMaterial 
            color="#98FB98" 
            roughness={0.8} 
            metalness={0.2} 
          />
        </mesh>

        {/* 添加场景元素 */}
        <Cabin position={[-8, 0, -8]} />
        <BrickHouse position={[0, 0, -8]} color="#DEB887" />
        <BrickHouse position={[8, 0, -8]} color="#F4A460" />
        <Library position={[0, 0, 8]} />
        
        {/* 多样化的守护者/宠物 */}
        <Guardian position={[-3, 0.5, 0]} type="cat" />
        <Guardian position={[3, 0.5, 0]} type="dog" />
        <Guardian position={[0, 0.5, -3]} type="robot" />
        <Guardian position={[0, 0.5, 3]} type="basic" />
        
        {/* 添加路灯 */}
        <group position={[6, 0, 6]}>
          <mesh position={[0, 3, 0]}>
            <cylinderGeometry args={[0.1, 0.1, 6, 8]} />
            <meshStandardMaterial color="#C0C0C0" />
          </mesh>
          <mesh position={[0, 6.2, 0]}>
            <sphereGeometry args={[0.5, 16, 16]} />
            <meshStandardMaterial color="#FFFF00" emissive="#FFFF00" emissiveIntensity={0.5} />
          </mesh>
        </group>
        
        <group position={[-6, 0, -6]}>
          <mesh position={[0, 3, 0]}>
            <cylinderGeometry args={[0.1, 0.1, 6, 8]} />
            <meshStandardMaterial color="#C0C0C0" />
          </mesh>
          <mesh position={[0, 6.2, 0]}>
            <sphereGeometry args={[0.5, 16, 16]} />
            <meshStandardMaterial color="#FFFF00" emissive="#FFFF00" emissiveIntensity={0.5} />
          </mesh>
        </group>
        
        {/* 添加花坛 */}
        <group position={[-10, 0, 0]}>
          <mesh position={[0, 0.5, 0]}>
            <cylinderGeometry args={[1.5, 1.5, 1, 32]} />
            <meshStandardMaterial color="#8B4513" />
          </mesh>
          <mesh position={[0, 1, 0]}>
            <sphereGeometry args={[1, 16, 16]} />
            <meshStandardMaterial color="#FF69B4" />
          </mesh>
        </group>
        
        {/* 添加相机控制器 */}
        <OrbitControls 
          enablePan={true} 
          enableZoom={true} 
          enableRotate={true}
          minDistance={5}
          maxDistance={50}
        />
      </Canvas>
    </div>
  );
};

export default KnowledgeTownScene;