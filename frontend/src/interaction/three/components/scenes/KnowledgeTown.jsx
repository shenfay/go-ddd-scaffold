import React from 'react';
import { useState, useEffect } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Text, Sky, Environment, Stars } from '@react-three/drei';
import SoundManager from '../../../../shared/utils/SoundManager';

// 知识节点组件 - C (概念节点) - 绿色球体
const ConceptNode = ({ position, name = "概念", level = 1, onClick }) => {
  const [hovered, setHovered] = React.useState(false);

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  // 根据等级设置颜色深浅
  const levelColors = [
    "#90EE90", // Lv1 - 浅绿色
    "#32CD32", // Lv2 - 酸橙绿
    "#228B22", // Lv3 - 森林绿
    "#006400", // Lv4 - 深绿色
  ];
  
  const color = levelColors[Math.min(level - 1, levelColors.length - 1)] || "#90EE90";

  return (
    <group position={position}>
      <mesh 
        onClick={onClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <sphereGeometry args={[0.8, 32, 32]} />
        <meshStandardMaterial 
          color={hovered ? "#FFD700" : color}
          roughness={0.2}
          metalness={0.1}
          emissive={hovered ? "#FFD700" : color}
          emissiveIntensity={hovered ? 0.3 : 0.1}
        />
      </mesh>
      <Text
        position={[0, 1.2, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -1.2, 0]}
        fontSize={0.3}
        color="yellow"
        anchorX="center"
        anchorY="middle"
      >
        Lv{level}
      </Text>
    </group>
  );
};

// 知识节点组件 - S (支撑技能) - 蓝色立方体
const SupportSkillNode = ({ position, name = "技能", level = 1, onClick }) => {
  const [hovered, setHovered] = React.useState(false);

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  // 根据等级设置颜色深浅
  const levelColors = [
    "#87CEEB", // Lv1 - 天蓝色
    "#4682B4", // Lv2 - 钢蓝色
    "#4169E1", // Lv3 - 皇家蓝
    "#00008B", // Lv4 - 深蓝色
  ];
  
  const color = levelColors[Math.min(level - 1, levelColors.length - 1)] || "#87CEEB";

  return (
    <group position={position}>
      <mesh 
        onClick={onClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <boxGeometry args={[1.2, 1.2, 1.2]} />
        <meshStandardMaterial 
          color={hovered ? "#FFD700" : color}
          roughness={0.3}
          metalness={0.2}
          emissive={hovered ? "#FFD700" : color}
          emissiveIntensity={hovered ? 0.3 : 0.1}
        />
      </mesh>
      <Text
        position={[0, 1.5, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -1.5, 0]}
        fontSize={0.3}
        color="yellow"
        anchorX="center"
        anchorY="middle"
      >
        Lv{level}
      </Text>
    </group>
  );
};

// 知识节点组件 - T (思维模式) - 黄色圆锥
const ThinkingNode = ({ position, name = "思维", level = 1, onClick }) => {
  const [hovered, setHovered] = React.useState(false);

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  // 根据等级设置颜色深浅
  const levelColors = [
    "#FFD700", // Lv1 - 金色
    "#FFA500", // Lv2 - 橙色
    "#FF8C00", // Lv3 - 深橙色
    "#FF4500", // Lv4 - 橙红色
  ];
  
  const color = levelColors[Math.min(level - 1, levelColors.length - 1)] || "#FFD700";

  return (
    <group position={position}>
      <mesh 
        onClick={onClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <coneGeometry args={[0.7, 1.5, 32]} />
        <meshStandardMaterial 
          color={hovered ? "#FFFFFF" : color}
          roughness={0.2}
          metalness={0.3}
          emissive={hovered ? "#FFFFFF" : color}
          emissiveIntensity={hovered ? 0.3 : 0.1}
        />
      </mesh>
      <Text
        position={[0, 2, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -2, 0]}
        fontSize={0.3}
        color="yellow"
        anchorX="center"
        anchorY="middle"
      >
        Lv{level}
      </Text>
    </group>
  );
};

// 知识节点组件 - P (问题模型) - 紫色环面
const ProblemNode = ({ position, name = "问题", level = 1, onClick }) => {
  const [hovered, setHovered] = React.useState(false);

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  // 根据等级设置颜色深浅
  const levelColors = [
    "#DDA0DD", // Lv1 - 梅花色
    "#9370DB", // Lv2 - 中紫色
    "#8A2BE2", // Lv3 - 蓝紫色
    "#4B0082", // Lv4 - 深紫罗兰色
  ];
  
  const color = levelColors[Math.min(level - 1, levelColors.length - 1)] || "#DDA0DD";

  return (
    <group position={position}>
      <mesh 
        onClick={onClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <torusGeometry args={[1, 0.3, 16, 100]} />
        <meshStandardMaterial 
          color={hovered ? "#FFFFFF" : color}
          roughness={0.3}
          metalness={0.4}
          emissive={hovered ? "#FFFFFF" : color}
          emissiveIntensity={hovered ? 0.3 : 0.1}
        />
      </mesh>
      <Text
        position={[0, 2, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -2, 0]}
        fontSize={0.3}
        color="yellow"
        anchorX="center"
        anchorY="middle"
      >
        Lv{level}
      </Text>
    </group>
  );
};

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
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchNodes = async () => {
      try {
        // 1. 首先获取所有的 Domains
        const domainsResponse = await fetch('http://localhost:8080/api/knowledge/domains');
        if (!domainsResponse.ok) throw new Error('Failed to fetch domains');
        const domains = await domainsResponse.json();

        if (domains.length === 0) {
          setLoading(false);
          return;
        }

        const firstDomain = domains[0];

        // 2. 获取该 Domain 下的所有 Trunks
        const trunksResponse = await fetch(`http://localhost:8080/api/knowledge/trunks/domain/${firstDomain.id}`);
        if (!trunksResponse.ok) throw new Error('Failed to fetch trunks');
        const trunks = await trunksResponse.json();

        if (trunks.length === 0) {
          setLoading(false);
          return;
        }

        const firstTrunk = trunks[0];

        // 3. 获取该 Trunk 下的所有 Nodes
        const nodesResponse = await fetch(`http://localhost:8080/api/knowledge/nodes/${firstTrunk.id}`);
        if (!nodesResponse.ok) throw new Error('Failed to fetch nodes');
        const nodesData = await nodesResponse.json();

        setNodes(nodesData);
        setLoading(false);
      } catch (err) {
        console.error(err);
        setError(err.message);
        setLoading(false);
      }
    };

    fetchNodes();
  }, []);

  // 计算节点位置
  const calculatePosition = (index, total) => {
    const radius = 10;
    const angle = (index / total) * 2 * Math.PI;
    return [
      Math.cos(angle) * radius,
      0.5, // y position
      Math.sin(angle) * radius
    ];
  };

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

        {/* 渲染从 API 获取的节点 */}
        {loading ? (
          <Text position={[0, 2, 0]} fontSize={0.5} color="white">加载中...</Text>
        ) : error ? (
          <Text position={[0, 2, 0]} fontSize={0.5} color="red">错误: {error}</Text>
        ) : (
          nodes.map((node, index) => {
            const position = calculatePosition(index, nodes.length);
            switch(node.type) {
              case 'C':
                return <ConceptNode key={node.id} position={position} name={node.name_child_key || "未知概念"} level={3} />;
              case 'S':
                return <SupportSkillNode key={node.id} position={position} name={node.name_child_key || "未知技能"} level={3} />;
              case 'T':
                return <ThinkingNode key={node.id} position={position} name={node.name_child_key || "未知思维"} level={3} />;
              case 'P':
                return <ProblemNode key={node.id} position={position} name={node.name_child_key || "未知问题"} level={3} />;
              default:
                return null;
            }
          })
        )}
        
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