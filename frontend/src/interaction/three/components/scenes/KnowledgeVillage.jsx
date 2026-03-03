import React from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Text, Sky, Environment } from '@react-three/drei';
import SoundManager from '../../../../shared/utils/SoundManager';

// 小木屋组件
const Cabin = ({ position, roofColor = "#A52A2A", wallColor = "#8B4513" }) => {
  return (
    <group position={position}>
      {/* 房身 */}
      <mesh position={[0, 1, 0]}>
        <boxGeometry args={[2, 2, 2]} />
        <meshStandardMaterial 
          color={wallColor} 
          roughness={0.7} 
          metalness={0.3} 
        />
      </mesh>
      {/* 屋顶 */}
      <mesh position={[0, 2.2, 0]} rotation={[0, 0, Math.PI / 4]}>
        <coneGeometry args={[1.8, 1.5, 4]} />
        <meshStandardMaterial 
          color={roofColor} 
          roughness={0.5} 
          metalness={0.2} 
        />
      </mesh>
      {/* 门 */}
      <mesh position={[0, 0, 1.01]}>
        <planeGeometry args={[0.8, 1.2]} />
        <meshStandardMaterial 
          color="#654321" 
          roughness={0.8} 
          metalness={0.1} 
        />
      </mesh>
      {/* 窗户 */}
      <mesh position={[-0.7, 0.5, 1.01]}>
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

// 菜园组件
const Garden = ({ position }) => {
  return (
    <group position={position}>
      {/* 围栏 */}
      <mesh position={[0, 0.5, 0]}>
        <boxGeometry args={[4, 1, 0.1]} />
        <meshStandardMaterial color="#8FBC8F" roughness={0.6} metalness={0.2} />
      </mesh>
      {/* 菜地 */}
      <mesh position={[0, 0, 0]} rotation={[-Math.PI / 2, 0, 0]}>
        <planeGeometry args={[3.8, 3.8]} />
        <meshStandardMaterial color="#228B22" side={2} roughness={0.9} metalness={0.1} />
      </mesh>
      {/* 几棵小树 */}
      <group position={[-1, 0, 0]}>
        <mesh position={[0, 1.5, 0]}>
          <sphereGeometry args={[0.5, 8, 8]} />
          <meshStandardMaterial color="#2E8B57" roughness={0.4} metalness={0.3} />
        </mesh>
        <mesh position={[0, 0.5, 0]}>
          <cylinderGeometry args={[0.1, 0.1, 1, 8]} />
          <meshStandardMaterial color="#8B4513" roughness={0.7} metalness={0.4} />
        </mesh>
      </group>
      <group position={[1, 0, 0]}>
        <mesh position={[0, 1.5, 0]}>
          <sphereGeometry args={[0.5, 8, 8]} />
          <meshStandardMaterial color="#2E8B57" roughness={0.4} metalness={0.3} />
        </mesh>
        <mesh position={[0, 0.5, 0]}>
          <cylinderGeometry args={[0.1, 0.1, 1, 8]} />
          <meshStandardMaterial color="#8B4513" roughness={0.7} metalness={0.4} />
        </mesh>
      </group>
    </group>
  );
};

// 守护者/宠物组件 (带有更多细节)
const Guardian = ({ position }) => {
  const [hovered, setHovered] = React.useState(false);

  const handleClick = () => {
    SoundManager.playSound('pet_interaction');
    console.log("Pet clicked!");
  };

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  return (
    <group>
      <mesh 
        position={position} 
        onClick={handleClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <sphereGeometry args={[0.5, 32, 32]} />
        <meshStandardMaterial 
          color={hovered ? "#FF6347" : "#FFA500"} // 悬停时变为番茄红
          roughness={0.2}
          metalness={0.1}
          emissive={hovered ? "#FF6347" : "#FFA500"}
          emissiveIntensity={hovered ? 0.2 : 0.1}
        />
        <Text
          position={[0, 1, 0]}
          fontSize={0.5}
          color="black"
          anchorX="center"
          anchorY="middle"
        >
          Pet
        </Text>
      </mesh>
      
      {/* 添加光环效果 */}
      {hovered && (
        <mesh position={position}>
          <torusGeometry args={[0.7, 0.05, 16, 100]} />
          <meshBasicMaterial color="#FFD700" transparent opacity={0.6} />
        </mesh>
      )}
    </group>
  );
};

// Updated component to fit the new routing structure
const Scene_KnowledgeVillage = () => {
  return (
    <div className="w-full h-screen bg-gray-800">
      <div className="container mx-auto px-4 py-8">
        <h2 className="text-2xl font-bold text-white mb-4">知识村庄场景 (Knowledge Village Scene)</h2>
        <p className="text-gray-300 mb-6">探索 Go DDD Scaffold 世界的初始阶段。</p>
        <div className="bg-gray-900 rounded-lg overflow-hidden shadow-xl">
          <Canvas 
            camera={{ position: [10, 10, 10], fov: 50 }} 
            className="w-full h-96"
            shadows
          >
            {/* 环境光 */}
            <ambientLight intensity={0.4} />
            
            {/* 主光源 */}
            <directionalLight 
              position={[10, 20, 15]} 
              intensity={1} 
              castShadow
              shadow-mapSize-width={1024}
              shadow-mapSize-height={1024}
            />
            
            {/* 点光源 */}
            <pointLight position={[10, 10, 10]} intensity={0.5} />
            
            {/* 天空环境 */}
            <Sky sunPosition={[100, 10, 100]} />
            
            {/* 环境贴图 */}
            <color attach="background" args={["#87CEEB"]} />
            
            {/* 地面 */}
            <mesh 
              rotation={[-Math.PI / 2, 0, 0]} 
              receiveShadow
            >
              <planeGeometry args={[50, 50]} />
              <meshStandardMaterial 
                color="#90EE90" 
                roughness={0.9} 
                metalness={0.1} 
              />
            </mesh>

            {/* 添加场景元素 */}
            <Cabin position={[-5, 0, -5]} roofColor="#8B0000" wallColor="#A0522D" />
            <Cabin position={[5, 0, -5]} roofColor="#2F4F4F" wallColor="#CD853F" />
            <Garden position={[0, 0, 5]} />
            <Guardian position={[0, 0.5, 0]} />
            
            {/* 添加更多环境元素 */}
            {/* 池塘 */}
            <mesh position={[-8, 0.01, 5]}>
              <circleGeometry args={[2, 32]} />
              <meshStandardMaterial 
                color="#4682B4" 
                roughness={0.1} 
                metalness={0.9} 
                side={2} // 双面渲染
              />
            </mesh>
            
            {/* 花园 */}
            <group position={[7, 0, 6]}>
              <mesh position={[-0.5, 0.2, 0]}>
                <sphereGeometry args={[0.2, 16, 16]} />
                <meshStandardMaterial color="#FF69B4" />
              </mesh>
              <mesh position={[0.5, 0.3, 0]}>
                <sphereGeometry args={[0.15, 16, 16]} />
                <meshStandardMaterial color="#FF1493" />
              </mesh>
              <mesh position={[0, 0.4, -0.5]}>
                <sphereGeometry args={[0.18, 16, 16]} />
                <meshStandardMaterial color="#FF6347" />
              </mesh>
            </group>
            
            {/* 添加相机控制器 */}
            <OrbitControls 
              enablePan={true} 
              enableZoom={true} 
              enableRotate={true} 
              minDistance={5}
              maxDistance={30}
            />
          </Canvas>
        </div>
      </div>
    </div>
  );
};

export default Scene_KnowledgeVillage;