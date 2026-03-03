import React, { useState, useRef, useMemo } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Text, Sphere, Box, Cone, Torus } from '@react-three/drei';
import * as THREE from 'three';

// 知识节点组件 - 概念节点 (C) - 绿色球体
const ConceptNode = ({ position, name = "概念", level = 1, onClick, isSelected }) => {
  const meshRef = useRef();
  const [hovered, setHovered] = useState(false);

  useFrame(() => {
    if (meshRef.current && (hovered || isSelected)) {
      meshRef.current.rotation.y += 0.01;
      meshRef.current.scale.setScalar(1.2);
    } else if (meshRef.current) {
      meshRef.current.scale.lerp(new THREE.Vector3(1, 1, 1), 0.1);
    }
  });

  // 根据等级设置颜色深浅
  const levelColors = [
    "#90EE90", // Lv1 - 浅绿色
    "#32CD32", // Lv2 - 酸橙绿
    "#228B22", // Lv3 - 森林绿
    "#006400", // Lv4 - 深绿色
    "#004d00", // Lv5 - 最深绿
  ];
  
  const color = levelColors[level - 1] || levelColors[0];

  return (
    <group position={position}>
      <Sphere
        ref={meshRef}
        args={[0.8, 32, 32]}
        onClick={onClick}
        onPointerOver={(e) => {
          e.stopPropagation();
          setHovered(true);
        }}
        onPointerOut={(e) => {
          e.stopPropagation();
          setHovered(false);
        }}
      >
        <meshStandardMaterial 
          color={isSelected ? "#FFD700" : hovered ? "#FFA500" : color}
          roughness={0.2}
          metalness={0.1}
          emissive={isSelected ? "#FFD700" : hovered ? "#FFA500" : color}
          emissiveIntensity={isSelected ? 0.3 : hovered ? 0.2 : 0.1}
        />
      </Sphere>
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

// 知识节点组件 - 支撑技能节点 (S) - 蓝色立方体
const SupportSkillNode = ({ position, name = "技能", level = 1, onClick, isSelected }) => {
  const meshRef = useRef();
  const [hovered, setHovered] = useState(false);

  useFrame(() => {
    if (meshRef.current && (hovered || isSelected)) {
      meshRef.current.rotation.y += 0.01;
      meshRef.current.scale.setScalar(1.2);
    } else if (meshRef.current) {
      meshRef.current.scale.lerp(new THREE.Vector3(1, 1, 1), 0.1);
    }
  });

  // 根据等级设置颜色深浅
  const levelColors = [
    "#87CEEB", // Lv1 - 天蓝色
    "#4682B4", // Lv2 - 钢蓝色
    "#4169E1", // Lv3 - 皇家蓝
    "#00008B", // Lv4 - 深蓝色
    "#00004d", // Lv5 - 最深蓝
  ];
  
  const color = levelColors[level - 1] || levelColors[0];

  return (
    <group position={position}>
      <Box
        ref={meshRef}
        args={[1.2, 1.2, 1.2]}
        onClick={onClick}
        onPointerOver={(e) => {
          e.stopPropagation();
          setHovered(true);
        }}
        onPointerOut={(e) => {
          e.stopPropagation();
          setHovered(false);
        }}
      >
        <meshStandardMaterial 
          color={isSelected ? "#FFD700" : hovered ? "#FFA500" : color}
          roughness={0.3}
          metalness={0.2}
          emissive={isSelected ? "#FFD700" : hovered ? "#FFA500" : color}
          emissiveIntensity={isSelected ? 0.3 : hovered ? 0.2 : 0.1}
        />
      </Box>
      <Text
        position={[0, 1.8, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -1.8, 0]}
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

// 知识节点组件 - 思维模式节点 (T) - 黄色圆锥
const ThinkingNode = ({ position, name = "思维", level = 1, onClick, isSelected }) => {
  const meshRef = useRef();
  const [hovered, setHovered] = useState(false);

  useFrame(() => {
    if (meshRef.current && (hovered || isSelected)) {
      meshRef.current.rotation.y += 0.01;
      meshRef.current.scale.setScalar(1.2);
    } else if (meshRef.current) {
      meshRef.current.scale.lerp(new THREE.Vector3(1, 1, 1), 0.1);
    }
  });

  // 根据等级设置颜色深浅
  const levelColors = [
    "#FFD700", // Lv1 - 金色
    "#FFA500", // Lv2 - 橙色
    "#FF8C00", // Lv3 - 深橙色
    "#FF4500", // Lv4 - 橙红色
    "#8B0000", // Lv5 - 暗红色
  ];
  
  const color = levelColors[level - 1] || levelColors[0];

  return (
    <group position={position}>
      <Cone
        ref={meshRef}
        args={[0.7, 1.5, 32]}
        onClick={onClick}
        onPointerOver={(e) => {
          e.stopPropagation();
          setHovered(true);
        }}
        onPointerOut={(e) => {
          e.stopPropagation();
          setHovered(false);
        }}
      >
        <meshStandardMaterial 
          color={isSelected ? "#FFFFFF" : hovered ? "#FFA500" : color}
          roughness={0.2}
          metalness={0.3}
          emissive={isSelected ? "#FFFFFF" : hovered ? "#FFA500" : color}
          emissiveIntensity={isSelected ? 0.3 : hovered ? 0.2 : 0.1}
        />
      </Cone>
      <Text
        position={[0, 2.2, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -2.2, 0]}
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

// 知识节点组件 - 问题模型节点 (P) - 紫色环面
const ProblemNode = ({ position, name = "问题", level = 1, onClick, isSelected }) => {
  const meshRef = useRef();
  const [hovered, setHovered] = useState(false);

  useFrame(() => {
    if (meshRef.current && (hovered || isSelected)) {
      meshRef.current.rotation.y += 0.01;
      meshRef.current.scale.setScalar(1.2);
    } else if (meshRef.current) {
      meshRef.current.scale.lerp(new THREE.Vector3(1, 1, 1), 0.1);
    }
  });

  // 根据等级设置颜色深浅
  const levelColors = [
    "#DDA0DD", // Lv1 - 梅花色
    "#9370DB", // Lv2 - 中紫色
    "#8A2BE2", // Lv3 - 蓝紫色
    "#4B0082", // Lv4 - 深紫罗兰色
    "#2F0F4F", // Lv5 - 最深紫色
  ];
  
  const color = levelColors[level - 1] || levelColors[0];

  return (
    <group position={position}>
      <Torus
        ref={meshRef}
        args={[1, 0.3, 16, 100]}
        onClick={onClick}
        onPointerOver={(e) => {
          e.stopPropagation();
          setHovered(true);
        }}
        onPointerOut={(e) => {
          e.stopPropagation();
          setHovered(false);
        }}
      >
        <meshStandardMaterial 
          color={isSelected ? "#FFFFFF" : hovered ? "#FFA500" : color}
          roughness={0.3}
          metalness={0.4}
          emissive={isSelected ? "#FFFFFF" : hovered ? "#FFA500" : color}
          emissiveIntensity={isSelected ? 0.3 : hovered ? 0.2 : 0.1}
        />
      </Torus>
      <Text
        position={[0, 2.2, 0]}
        fontSize={0.4}
        color="white"
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      <Text
        position={[0, -2.2, 0]}
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

// 连接线组件
const ConnectionLine = ({ start, end, type }) => {
  const points = useMemo(() => [
    new THREE.Vector3(...start),
    new THREE.Vector3(...end)
  ], [start, end]);

  return (
    <line>
      <bufferGeometry attach="geometry">
        <bufferAttribute
          attachObject={['attributes', 'position']}
          count={points.length}
          array={new Float32Array(points.flatMap(p => [p.x, p.y, p.z]))}
          itemSize={3}
        />
      </bufferGeometry>
      <lineBasicMaterial 
        attach="material" 
        color="#94a3b8" 
        linewidth={2}
        transparent
        opacity={0.6}
      />
    </line>
  );
};

const KnowledgeGraph3DCanvas = ({ nodes, edges, onNodeClick, selectedNode }) => {
  return (
    <Canvas 
      camera={{ position: [15, 15, 15], fov: 50 }} 
      shadows
      style={{ height: '100%', background: 'linear-gradient(to bottom, #1e3a8a, #3b82f6)' }}
    >
      {/* 环境光 */}
      <ambientLight intensity={0.4} />
      
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
      
      {/* 星空背景 */}
      <pointLight position={[-10, -10, -10]} intensity={0.2} color="#4f46e5" />
      
      {/* 连接线 */}
      {edges?.map((edge, index) => {
        const sourceNode = nodes?.find(n => n.id === edge.source);
        const targetNode = nodes?.find(n => n.id === edge.target);
        
        if (sourceNode && targetNode) {
          return (
            <ConnectionLine 
              key={index} 
              start={[sourceNode.x, sourceNode.y, 0]} 
              end={[targetNode.x, targetNode.y, 0]} 
              type={edge.type} 
            />
          );
        }
        return null;
      })}

      {/* 知识节点 */}
      {nodes?.map((node) => {
        const isSelected = selectedNode && selectedNode.id === node.id;
        const position = [node.x, node.y, 0];
        
        switch (node.type) {
          case 'C':
            return (
              <ConceptNode 
                key={node.id} 
                position={position} 
                name={node.name} 
                level={node.level} 
                onClick={() => onNodeClick(node)}
                isSelected={isSelected}
              />
            );
          case 'S':
            return (
              <SupportSkillNode 
                key={node.id} 
                position={position} 
                name={node.name} 
                level={node.level} 
                onClick={() => onNodeClick(node)}
                isSelected={isSelected}
              />
            );
          case 'T':
            return (
              <ThinkingNode 
                key={node.id} 
                position={position} 
                name={node.name} 
                level={node.level} 
                onClick={() => onNodeClick(node)}
                isSelected={isSelected}
              />
            );
          case 'P':
            return (
              <ProblemNode 
                key={node.id} 
                position={position} 
                name={node.name} 
                level={node.level} 
                onClick={() => onNodeClick(node)}
                isSelected={isSelected}
              />
            );
          default:
            return (
              <ConceptNode 
                key={node.id} 
                position={position} 
                name={node.name} 
                level={node.level} 
                onClick={() => onNodeClick(node)}
                isSelected={isSelected}
              />
            );
        }
      })}

      {/* 添加地面网格 */}
      <gridHelper args={[50, 50, '#4f46e5', '#312e81']} position={[0, -2, 0]} />
      
      {/* 相机控制器 */}
      <OrbitControls 
        enablePan={true} 
        enableZoom={true} 
        enableRotate={true}
        minDistance={5}
        maxDistance={50}
      />
    </Canvas>
  );
};

const KnowledgeGraph3D = ({ initialNodes, initialEdges, onNodeClick }) => {
  const [selectedNode, setSelectedNode] = useState(null);
  const [nodes, setNodes] = useState(initialNodes || [
    { id: '1', name: '整数', type: 'C', level: 1, x: 0, y: 0 },
    { id: '2', name: '分数', type: 'C', level: 2, x: 5, y: 3 },
    { id: '3', name: '小数', type: 'C', level: 3, x: 10, y: 0 },
    { id: '4', name: '加法', type: 'S', level: 1, x: -5, y: 3 },
    { id: '5', name: '减法', type: 'S', level: 1, x: -5, y: -3 },
    { id: '6', name: '乘法', type: 'S', level: 3, x: 15, y: 3 },
    { id: '7', name: '数形结合', type: 'T', level: 3, x: 0, y: 8 },
    { id: '8', name: '鸡兔同笼', type: 'P', level: 3, x: 10, y: 8 },
  ]);
  
  const [edges, setEdges] = useState(initialEdges || [
    { source: '1', target: '2', type: 'PREREQ' },
    { source: '2', target: '3', type: 'PREREQ' },
    { source: '1', target: '4', type: 'SUP_SKILL' },
    { source: '1', target: '5', type: 'SUP_SKILL' },
    { source: '2', target: '6', type: 'SUP_SKILL' },
    { source: '4', target: '7', type: 'THINK_PAT' },
    { source: '5', target: '7', type: 'THINK_PAT' },
    { source: '7', target: '8', type: 'PREREQ' },
  ]);

  const handleNodeClick = (node) => {
    setSelectedNode(node);
    if (onNodeClick) {
      onNodeClick(node);
    }
  };

  return (
    <div className="w-full h-full flex flex-col">
      <div className="flex-grow relative">
        <KnowledgeGraph3DCanvas 
          nodes={nodes} 
          edges={edges} 
          onNodeClick={handleNodeClick}
          selectedNode={selectedNode}
        />
      </div>
      
      {/* 信息面板 */}
      {selectedNode && (
        <div className="absolute top-4 left-4 bg-white bg-opacity-90 p-4 rounded-lg shadow-lg max-w-xs">
          <h3 className="font-bold text-lg text-gray-800">{selectedNode.name}</h3>
          <div className="mt-2 space-y-1 text-sm">
            <div><span className="font-medium">类型:</span> 
              {selectedNode.type === 'C' ? ' 概念' : 
               selectedNode.type === 'S' ? ' 支撑技能' : 
               selectedNode.type === 'T' ? ' 思维模式' : 
               selectedNode.type === 'P' ? ' 问题模型' : ' 未知'}
            </div>
            <div><span className="font-medium">等级:</span> Lv{selectedNode.level}</div>
            <div><span className="font-medium">ID:</span> {selectedNode.id}</div>
          </div>
        </div>
      )}
      
      {/* 图例 */}
      <div className="absolute top-4 right-4 bg-white bg-opacity-90 p-4 rounded-lg shadow-lg">
        <h3 className="font-bold text-gray-800 mb-2">图例</h3>
        <div className="space-y-2 text-xs">
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-green-400 mr-2"></div>
            <span>概念 (C)</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded bg-blue-400 mr-2"></div>
            <span>技能 (S)</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded-full bg-yellow-400 mr-2"></div>
            <span>思维 (T)</span>
          </div>
          <div className="flex items-center">
            <div className="w-3 h-3 rounded bg-purple-400 mr-2"></div>
            <span>问题 (P)</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeGraph3D;