import React, { useState, useRef, Suspense, useEffect } from 'react';
import { Canvas, useFrame, useThree } from '@react-three/fiber';
import { OrbitControls, Text, Sky, Environment, Stars, PerspectiveCamera, PointerLockControls } from '@react-three/drei';
import * as THREE from 'three';
import SoundManager from '../../../../shared/utils/SoundManager';

// --- 纹理生成工具 ---

// 生成草地纹理
const createGrassTexture = () => {
  const canvas = document.createElement('canvas');
  canvas.width = 256;
  canvas.height = 256;
  const ctx = canvas.getContext('2d');
  
  // 基础草地颜色
  ctx.fillStyle = '#5AD17A';
  ctx.fillRect(0, 0, 256, 256);
  
  // 添加草地纹理细节
  for (let i = 0; i < 1000; i++) {
    const x = Math.random() * 256;
    const y = Math.random() * 256;
    const size = Math.random() * 2;
    const opacity = Math.random() * 0.3;
    
    ctx.fillStyle = `rgba(76, 175, 80, ${opacity})`; // 深绿色
    ctx.fillRect(x, y, size, size * 2);
  }
  
  const texture = new THREE.CanvasTexture(canvas);
  texture.magFilter = THREE.NearestFilter;
  texture.minFilter = THREE.NearestFilter;
  return texture;
};

// 生成地面纹理
const createStreetTexture = () => {
  const canvas = document.createElement('canvas');
  canvas.width = 256;
  canvas.height = 256;
  const ctx = canvas.getContext('2d');
  
  // 基础街道颜色
  ctx.fillStyle = '#555555';
  ctx.fillRect(0, 0, 256, 256);
  
  // 添加混凝土纹理
  for (let i = 0; i < 2000; i++) {
    const x = Math.random() * 256;
    const y = Math.random() * 256;
    const brightness = Math.random() * 50;
    const opacity = Math.random() * 0.2;
    
    ctx.fillStyle = `rgba(${brightness}, ${brightness}, ${brightness}, ${opacity})`;
    ctx.fillRect(x, y, 1, 1);
  }
  
  const texture = new THREE.CanvasTexture(canvas);
  texture.magFilter = THREE.NearestFilter;
  texture.minFilter = THREE.NearestFilter;
  return texture;
};

// --- 背景元素 ---

// 目标位置标记（点击地面后出现的圆圈）
const TargetMarker = ({ position }) => (
  <mesh position={[position.x, 0.05, position.z]} rotation={[-Math.PI / 2, 0, 0]}>
    <ringGeometry args={[0.3, 0.4, 32]} />
    <meshBasicMaterial color="#FFEB3B" transparent opacity={0.6} />
  </mesh>
);

// 卡通山脉
const Mountains = () => (
  <group position={[0, 0, -25]}>
    <mesh position={[-15, 4, 0]} rotation={[0, Math.PI / 4, 0]}>
      <coneGeometry args={[10, 15, 4]} />
      <meshStandardMaterial color="#95A5A6" flatShading />
    </mesh>
    <mesh position={[15, 3, -2]} rotation={[0, -Math.PI / 4, 0]}>
      <coneGeometry args={[12, 12, 4]} />
      <meshStandardMaterial color="#7F8C8D" flatShading />
    </mesh>
  </group>
);

// 卡通灌木
const Bush = ({ position, color = "#2ECC71" }) => (
  <mesh position={position} castShadow>
    <sphereGeometry args={[0.6, 8, 8]} />
    <meshStandardMaterial color={color} flatShading />
  </mesh>
);

// 卡通建筑物（带窗户）
const Building = ({ position, width = 3, height = 4, depth = 2, color = "#E74C3C" }) => {
  return (
    <group position={position}>
      {/* 主体 */}
      <mesh castShadow>
        <boxGeometry args={[width, height, depth]} />
        <meshStandardMaterial color={color} flatShading />
      </mesh>
      {/* 屋顶 */}
      <mesh position={[0, height / 2, 0]} castShadow>
        <coneGeometry args={[width * 0.6, height * 0.4, 4]} />
        <meshStandardMaterial color={"#C0392B"} flatShading />
      </mesh>
      {/* 窗户 */}
      {Array.from({ length: 2 }).map((_, i) => (
        <mesh key={`window-${i}`} position={[-width / 2 - 0.01, height / 2, -0.5 + i * 1.2]} castShadow>
          <boxGeometry args={[0.4, 0.4, 0.1]} />
          <meshStandardMaterial color="#87CEEB" flatShading />
        </mesh>
      ))}
      {/* 门 */}
      <mesh position={[0, -height / 4, -depth / 2 - 0.01]} castShadow>
        <boxGeometry args={[0.6, 1.5, 0.1]} />
        <meshStandardMaterial color="#8B4513" flatShading />
      </mesh>
    </group>
  );
};

// 旋转的立方体（装饰物）
const RotatingCube = ({ position, size = 0.5, color = "#3498DB" }) => {
  const meshRef = useRef();
  
  useFrame(() => {
    if (meshRef.current) {
      meshRef.current.rotation.x += 0.01;
      meshRef.current.rotation.y += 0.01;
    }
  });
  
  return (
    <mesh ref={meshRef} position={position} castShadow>
      <boxGeometry args={[size, size, size]} />
      <meshStandardMaterial color={color} flatShading />
    </mesh>
  );
};

// 卡通树木
const Tree = ({ position, trunkColor = "#8B4513", leavesColor = "#27AE60" }) => (
  <group position={position}>
    {/* 树干 */}
    <mesh position={[0, 1, 0]} castShadow>
      <cylinderGeometry args={[0.3, 0.4, 2, 8]} />
      <meshStandardMaterial color={trunkColor} flatShading />
    </mesh>
    {/* 树冠 - 三个球叠罪 */}
    <mesh position={[0, 2.5, 0]} castShadow>
      <sphereGeometry args={[1.2, 8, 8]} />
      <meshStandardMaterial color={leavesColor} flatShading />
    </mesh>
    <mesh position={[-0.3, 3.2, 0]} castShadow>
      <sphereGeometry args={[1, 8, 8]} />
      <meshStandardMaterial color={leavesColor} flatShading />
    </mesh>
    <mesh position={[0.3, 3.2, 0]} castShadow>
      <sphereGeometry args={[1, 8, 8]} />
      <meshStandardMaterial color={leavesColor} flatShading />
    </mesh>
  </group>
);

// 围栅/桶
const Fence = ({ position, length = 5 }) => {
  const fenceColor = "#A0522D";
  return (
    <group position={position}>
      {Array.from({ length: Math.floor(length) }).map((_, i) => (
        <group key={i} position={[i * 1.2, 0, 0]}>
          {/* 朴子 */}
          <mesh position={[0, 0.5, 0]} castShadow>
            <boxGeometry args={[0.2, 1, 0.2]} />
            <meshStandardMaterial color={fenceColor} flatShading />
          </mesh>
          {/* 横檇 */}
          <mesh position={[0, 0.8, 0]} castShadow>
            <boxGeometry args={[1.2, 0.15, 0.15]} />
            <meshStandardMaterial color={fenceColor} flatShading />
          </mesh>
        </group>
      ))}
    </group>
  );
};

// 水体（湖或河）
const Water = ({ position, width = 10, depth = 5 }) => (
  <mesh position={[...position, -0.05]} rotation={[-Math.PI / 2, 0, 0]}>
    <planeGeometry args={[width, depth]} />
    <meshStandardMaterial color="#3498DB" transparent opacity={0.6} flatShading />
  </mesh>
);

// 农场区域
const Farm = ({ position }) => (
  <group position={position}>
    {/* 农场背景 */}
    <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0, 0]} receiveShadow>
      <planeGeometry args={[8, 8]} />
      <meshStandardMaterial color="#D2B48C" flatShading />
    </mesh>
    {/* 农作物纲 */}
    {Array.from({ length: 4 }).map((_, i) =>
      Array.from({ length: 4 }).map((_, j) => (
        <mesh key={`crop-${i}-${j}`} position={[-3 + i * 2, 0.05, -3 + j * 2]} castShadow>
          <boxGeometry args={[0.6, 0.3, 0.6]} />
          <meshStandardMaterial color="#7CB342" flatShading />
        </mesh>
      ))
    )}
    {/* 农场围栅 */}
    <Fence position={[-4, 0, -4]} length={8} />
  </group>
);

// 牞场区域
const Pasture = ({ position }) => (
  <group position={position}>
    {/* 牞场背景 - 深绿草地 */}
    <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0, 0]} receiveShadow>
      <planeGeometry args={[8, 8]} />
      <meshStandardMaterial color="#558B2F" flatShading />
    </mesh>
    {/* 羊或牡等动物 */}
    {Array.from({ length: 3 }).map((_, i) => (
      <mesh key={`animal-${i}`} position={[-2 + i * 3, 0.2, -1]} castShadow>
        <boxGeometry args={[0.6, 0.4, 1]} />
        <meshStandardMaterial color="#FFFFFF" flatShading />
      </mesh>
    ))}
    {/* 牞场围栅 */}
    <Fence position={[-4, 0, -4]} length={8} />
    {/* 树木装饰 */}
    <Tree position={[-3, 0, 3]} />
    <Tree position={[3, 0, 3]} />
  </group>
);

// 花朵装饰
const Flower = ({ position, color = "#FF69B4" }) => (
  <group position={position}>
    {/* 花茎 */}
    <mesh position={[0, 0.3, 0]} castShadow>
      <cylinderGeometry args={[0.08, 0.08, 0.6, 8]} />
      <meshStandardMaterial color="#228B22" flatShading />
    </mesh>
    {/* 花瓣 - 5个花瓣 */}
    {Array.from({ length: 5 }).map((_, i) => {
      const angle = (i / 5) * Math.PI * 2;
      return (
        <mesh key={`petal-${i}`} position={[Math.cos(angle) * 0.25, 0.65, Math.sin(angle) * 0.25]} castShadow>
          <sphereGeometry args={[0.15, 8, 8]} />
          <meshStandardMaterial color={color} flatShading />
        </mesh>
      );
    })}
    {/* 花心 */}
    <mesh position={[0, 0.65, 0]} castShadow>
      <sphereGeometry args={[0.1, 8, 8]} />
      <meshStandardMaterial color="#FFD700" flatShading />
    </mesh>
  </group>
);

// 石头/岩石
const Rock = ({ position, size = 0.5, color = "#8B7355" }) => (
  <mesh position={position} castShadow>
    <dodecahedronGeometry args={[size, 0]} />
    <meshStandardMaterial color={color} flatShading />
  </mesh>
);

// 小方石床（肤色水拣床）
const SmallPlanter = ({ position, color = "#8B6F47" }) => (
  <group position={position}>
    {/* 对称箱 */}
    <mesh castShadow>
      <boxGeometry args={[0.8, 0.4, 0.8]} />
      <meshStandardMaterial color={color} flatShading />
    </mesh>
    {/* 垨里的塌击 */}
    <mesh position={[0, 0.3, 0]} castShadow>
      <sphereGeometry args={[0.35, 8, 8]} />
      <meshStandardMaterial color="#7CB342" flatShading />
    </mesh>
  </group>
);

// --- 街道元素 ---

const Street = ({ onGroundClick }) => {
  const grassTexture = createGrassTexture();
  const streetTexture = createStreetTexture();
  
  return (
    <group onPointerDown={(e) => {
      e.stopPropagation();
      // 只有点击地面时才触发移动
      if (e.faceIndex !== undefined) onGroundClick(e.point);
    }}>
      {/* 草地背景 - 使用生成的纹理 */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, -0.01, 0]} receiveShadow>
        <planeGeometry args={[100, 100]} />
        <meshStandardMaterial map={grassTexture} /> 
      </mesh>
      {/* 主街道 - 使用生成的纹理 */}
      <mesh rotation={[-Math.PI / 2, 0, 0]} position={[0, 0, 0]} receiveShadow>
        <planeGeometry args={[10, 100]} />
        <meshStandardMaterial map={streetTexture} />
      </mesh>
      {/* 车道线 */}
      {Array.from({ length: 30 }).map((_, i) => (
        <mesh key={i} position={[0, 0.02, -28 + i * 2]}>
          <planeGeometry args={[0.4, 1]} />
          <meshStandardMaterial color="white" opacity={0.7} transparent />
        </mesh>
      ))}
    </group>
  );
};

// --- 角色系统 ---

const Player = ({ onMove, targetPos }) => {
  const rootRef = useRef();
  const visualsRef = useRef();
  const [keys, setKeys] = useState({});
  const speed = 0.15;
  const currentTarget = useRef(new THREE.Vector3(0, 0, 5));
  
  // 定义可行走的边界（基于场景地面范围）
  // 草地背景是 100x100 的平面，中心在 [0, 0, 0]
  // 实际边界：X轴 [-50, 50]，Z轴 [-50, 50]
  // 为了避免边界硬切，设置为 90% 范围内
  const boundaries = {
    minX: -48,  // 草地左边界
    maxX: 48,   // 草地右边界
    minZ: -48,  // 草地后边界
    maxZ: 48    // 草地前边界
  };

  // 边界检测函数
  const clampPosition = (pos) => {
    pos.x = Math.max(boundaries.minX, Math.min(boundaries.maxX, pos.x));
    pos.z = Math.max(boundaries.minZ, Math.min(boundaries.maxZ, pos.z));
    return pos;
  };

  // 碰撞检测：检查是否与障碍物撞撞
  const checkCollisions = (pos, nextPos) => {
    const minDistance = 0.4 + 0.5; // 玩家半径 + 消牛的边上分
    
    // 检查 NPC 碰撞（假设 NPC 半径为 0.45 + 0.35身体半径 ≈ 0.8）
    const npcPositions = [
      [-3, 0, -5],    // 小明
      [3, 0, -12],    // 王老师
      [6, 0, 2]       // 李老板
    ];
    
    for (let npcPos of npcPositions) {
      const npc = new THREE.Vector3(...npcPos);
      const dist = nextPos.distanceTo(npc);
      if (dist < minDistance + 0.8) {
        // 碰撞，不洈离开当前位置
        return pos;
      }
    }
    
    // 检查灌木碰撞（假设灌木半径为 0.6）
    const bushPositions = [
      [-6, 0.5, -10],
      [7, 0.5, -15]
    ];
    
    for (let bushPos of bushPositions) {
      const bush = new THREE.Vector3(...bushPos);
      const dist = nextPos.distanceTo(bush);
      if (dist < minDistance + 0.6) {
        // 碰撞，不洈离开当前位置
        return pos;
      }
    }
    
    // 检查建筑物碰撞（使用AABB - 轴对掐丘枓检测）
    const buildings = [
      { pos: [-15, 2, -10], width: 4, height: 4, depth: 3 },  // 红色建筑
      { pos: [10, 2, -15], width: 3.5, height: 4, depth: 3 }, // 海蜂色建筑
      { pos: [-8, 2, 5], width: 3, height: 4, depth: 2.5 },   // 不丈建筑
      { pos: [15, 2, 10], width: 3.5, height: 4, depth: 3 }   // 红赤色建筑
    ];
    
    for (let building of buildings) {
      // AABB 碰撞检测
      const halfWidth = building.width / 2 + 0.5;  // 消边距
      const halfDepth = building.depth / 2 + 0.5;
      
      if (Math.abs(nextPos.x - building.pos[0]) < halfWidth &&
          Math.abs(nextPos.z - building.pos[2]) < halfDepth) {
        // 碰撞，不洈离开当前位置
        return pos;
      }
    }
    
    // 检查树木碰撞（树半径约 1.5）
    const treePositions = [
      [-20, 0, 15],   // 树
      [18, 0, -20],   // 树
      [-35, 0, -5],   // 农场区树
      [35, 0, 5]      // 牞场区树
    ];
    
    for (let treePos of treePositions) {
      const tree = new THREE.Vector3(...treePos);
      const dist = nextPos.distanceTo(tree);
      if (dist < minDistance + 1.5) {
        // 碰撞，不洈离开当前位置
        return pos;
      }
    }
    
    // 检查农场区域碰撞（AABB）
    const farmArea = { pos: [-30, 0, -5], width: 10, depth: 10 };
    const farmHalfWidth = farmArea.width / 2 + 0.3;
    const farmHalfDepth = farmArea.depth / 2 + 0.3;
    
    if (Math.abs(nextPos.x - farmArea.pos[0]) < farmHalfWidth &&
        Math.abs(nextPos.z - farmArea.pos[2]) < farmHalfDepth) {
      // 碰撞，不洈离开当前位置
      return pos;
    }
    
    // 检查牞场区域碰撞（AABB）
    const pastureArea = { pos: [30, 0, 5], width: 10, depth: 10 };
    const pastureHalfWidth = pastureArea.width / 2 + 0.3;
    const pastureHalfDepth = pastureArea.depth / 2 + 0.3;
    
    if (Math.abs(nextPos.x - pastureArea.pos[0]) < pastureHalfWidth &&
        Math.abs(nextPos.z - pastureArea.pos[2]) < pastureHalfDepth) {
      // 碰撞，不洈离开当前位置
      return pos;
    }
    
    // 没有碰撞，返回新位置
    return nextPos;
  };

  // 可选：添加一个边界可视化辅助组件（开发时使用）
  const showBoundaryHelper = false; // 设置为 true 来显示边界框

  useEffect(() => {
    if (targetPos) currentTarget.current.copy(targetPos);
  }, [targetPos]);

  useEffect(() => {
    const handleDown = (e) => setKeys(prev => ({ ...prev, [e.code]: true }));
    const handleUp = (e) => setKeys(prev => ({ ...prev, [e.code]: false }));
    window.addEventListener('keydown', handleDown);
    window.addEventListener('keyup', handleUp);
    return () => {
      window.removeEventListener('keydown', handleDown);
      window.removeEventListener('keyup', handleUp);
    };
  }, []);

  useFrame((state) => {
    if (!rootRef.current || !visualsRef.current) return;

    const pos = rootRef.current.position;
    const nextPos = pos.clone(); // 洈计算下一个位置
    let isMoving = false;

    // 1. 优先处理键盘控制
    const keyboardMove = new THREE.Vector3(0, 0, 0);
    if (keys['KeyW'] || keys['ArrowUp']) keyboardMove.z -= 1;
    if (keys['KeyS'] || keys['ArrowDown']) keyboardMove.z += 1;
    if (keys['KeyA'] || keys['ArrowLeft']) keyboardMove.x -= 1;
    if (keys['KeyD'] || keys['ArrowRight']) keyboardMove.x += 1;

    if (keyboardMove.length() > 0) {
      keyboardMove.normalize().multiplyScalar(speed);
      nextPos.add(keyboardMove);
      // 应用边界检测和碎撞检测
      clampPosition(nextPos);
      const safePos = checkCollisions(pos, nextPos);
      pos.copy(safePos);
      visualsRef.current.rotation.y = Math.atan2(keyboardMove.x, keyboardMove.z);
      currentTarget.current.copy(pos); 
      isMoving = true;
    } 
    // 2. 鼠标点击控制
    else if (pos.distanceTo(currentTarget.current) > 0.2) {
      const direction = new THREE.Vector3().subVectors(currentTarget.current, pos).normalize();
      nextPos.add(direction.multiplyScalar(speed));
      // 应用边界检测和碎撞检测
      clampPosition(nextPos);
      const safePos = checkCollisions(pos, nextPos);
      pos.copy(safePos);
      visualsRef.current.rotation.y = Math.atan2(direction.x, direction.z);
      isMoving = true;
    }

    if (isMoving) {
      // 只晃动视觉部分，不晃动根节点
      visualsRef.current.rotation.z = Math.sin(state.clock.elapsedTime * 10) * 0.1;
      onMove(pos);
    } else {
      visualsRef.current.rotation.z = 0;
    }
  });

  return (
    <group ref={rootRef} position={[0, 0, 5]}>
      {/* 将所有视觉元素放入一个独立的组，单独进行晃动动画 */}
      <group ref={visualsRef}>
        {/* 头部 */}
        <mesh position={[0, 1.6, 0]} castShadow>
          <sphereGeometry args={[0.45, 16, 16]} />
          <meshStandardMaterial color="#FFE0BD" />
        </mesh>
        {/* 眼睛 */}
        <mesh position={[-0.15, 1.7, 0.35]}><sphereGeometry args={[0.06, 8, 8]} /><meshStandardMaterial color="black" /></mesh>
        <mesh position={[0.15, 1.7, 0.35]}><sphereGeometry args={[0.06, 8, 8]} /><meshStandardMaterial color="black" /></mesh>
        {/* 身体 - 蓝色卫衣 */}
        <mesh position={[0, 0.8, 0]} castShadow>
          <cylinderGeometry args={[0.35, 0.45, 1, 16]} />
          <meshStandardMaterial color="#3498DB" />
        </mesh>
        {/* 小背包 */}
        <mesh position={[0, 1, -0.4]}><boxGeometry args={[0.5, 0.6, 0.2]} /><meshStandardMaterial color="#E67E22" /></mesh>
      </group>
      {/* 边界可视化辅助（调试用） */}
      {showBoundaryHelper && (
        <group>
          <lineSegments>
            <bufferGeometry>
              <bufferAttribute
                attach="attributes-position"
                count={8}
                array={new Float32Array([
                  boundaries.minX, 0, boundaries.minZ,
                  boundaries.maxX, 0, boundaries.minZ,
                  boundaries.maxX, 0, boundaries.minZ,
                  boundaries.maxX, 0, boundaries.maxZ,
                  boundaries.maxX, 0, boundaries.maxZ,
                  boundaries.minX, 0, boundaries.maxZ,
                  boundaries.minX, 0, boundaries.maxZ,
                  boundaries.minX, 0, boundaries.minZ
                ])}
                itemSize={3}
              />
            </bufferGeometry>
            <lineBasicMaterial color="#FF0000" />
          </lineSegments>
        </group>
      )}
    </group>
  );
};

// 动态方向光（基于时间旋转）
const DynamicDirectionalLight = () => {
  const lightRef = useRef();
  
  useEffect(() => {
    // 改进阴影质量
    if (lightRef.current) {
      lightRef.current.shadow.mapSize.width = 2048;
      lightRef.current.shadow.mapSize.height = 2048;
      lightRef.current.shadow.camera.near = 0.5;
      lightRef.current.shadow.camera.far = 500;
      lightRef.current.shadow.camera.left = -100;
      lightRef.current.shadow.camera.right = 100;
      lightRef.current.shadow.camera.top = 100;
      lightRef.current.shadow.camera.bottom = -100;
      lightRef.current.shadow.bias = -0.0001; // 减少阻光
      lightRef.current.shadow.blurSamples = 8; // 浄布窗口不是常量了
    }
  }, []);
  
  useFrame((state) => {
    if (lightRef.current) {
      // 基于时间旋转方向光
      const angle = state.clock.elapsedTime * 0.3; // 每秒旋转速度
      lightRef.current.position.x = Math.cos(angle) * 30;
      lightRef.current.position.z = Math.sin(angle) * 25 - 15;
      lightRef.current.position.y = 30;
    }
  });
  
  return (
    <directionalLight 
      ref={lightRef} 
      position={[20, 30, 10]} 
      intensity={1.3} 
      castShadow
    />
  );
};

// 跟随相机
const FollowCamera = ({ targetPos }) => {
  const { camera } = useThree();
  const controlsRef = useRef();

  useFrame(() => {
    if (targetPos) {
      // 1. 计算相机目标位置（第三人称偏移）
      const desiredPos = new THREE.Vector3(targetPos.x, targetPos.y + 8, targetPos.z + 12);
      camera.position.lerp(desiredPos, 0.05); // 降低平滑系数使跟随更稳定

      // 2. 如果有 OrbitControls，更新它的 target，而不是直接 lookAt
      if (controlsRef.current) {
        controlsRef.current.target.lerp(new THREE.Vector3(targetPos.x, targetPos.y + 1, targetPos.z), 0.1);
        controlsRef.current.update();
      } else {
        camera.lookAt(targetPos.x, targetPos.y + 1, targetPos.z);
      }
    }
  });

  return <OrbitControls ref={controlsRef} makeDefault minDistance={10} maxDistance={40} />;
};

// --- 主场景 ---

const MathCityScene = () => {
  const [playerPos, setPlayerPos] = useState(new THREE.Vector3(0, 0, 5));
  const [targetPos, setTargetPos] = useState(new THREE.Vector3(0, 0, 5));
  const [showMarker, setShowMarker] = useState(false);
  const [dialogue, setDialogue] = useState(null);
  const [nearbyNPC, setNearbyNPC] = useState(null); // 跟踪附近的NPC
  const [mousePointer, setMousePointer] = useState('default'); // 鼠标光标状态
  const [dialogueCooldown, setDialogueCooldown] = useState(0); // 对话框冷却时间
  
  // NPC 列表及其信息
  const npcs = [
    { id: 1, type: 'child', name: '小明', position: [-3, 0, -5] },
    { id: 2, type: 'teacher', name: '王老师', position: [3, 0, -12] },
    { id: 3, type: 'shopkeeper', name: '李老板', position: [6, 0, 2] }
  ];
  
  // NPC接近距离阈值
  const NPC_INTERACTION_DISTANCE = 3;

  // 检测玩家是否接近任何NPC
  useEffect(() => {
    let closestNPC = null;
    let closestDistance = Infinity;
    
    npcs.forEach(npc => {
      const npcPos = new THREE.Vector3(...npc.position);
      const distance = playerPos.distanceTo(npcPos);
      
      if (distance < NPC_INTERACTION_DISTANCE && distance < closestDistance) {
        closestNPC = npc;
        closestDistance = distance;
      }
    });
    
    setNearbyNPC(closestNPC);
    setMousePointer(closestNPC ? 'pointer' : 'default');
  }, [playerPos]);
  
  // 应用鼠标光标样式
  useEffect(() => {
    document.body.style.cursor = mousePointer;
    return () => {
      document.body.style.cursor = 'default';
    };
  }, [mousePointer]);

  const handleGroundClick = (point) => {
    setTargetPos(new THREE.Vector3(point.x, 0, point.z));
    setShowMarker(true);
    // 1.5秒后隐藏标记
    setTimeout(() => setShowMarker(false), 1500);
  };

  const handleCharacterClick = (name, type) => {
    // 仅在冷却执完或厶有时才要求拨打，禁止在冷却中再次点击
    if (dialogueCooldown > 0) return;
      
    const messages = {
      child: "嘯！想跟我一起比赛算术吗？",
      teacher: "你好孩子，数学是通往未来的大门。",
      scientist: "观察这些建筑的比例，数学无处不在。",
      shopkeeper: "今天所有的圆柱形商品都打八折哦！"
    };
    setDialogue({ name, message: messages[type] || "你好呀！" });
  };
  
  // 接近NPC时自动显示对话框
  const handleNPCProximity = (npc) => {
    if (npc) {
      const messages = {
        child: "嘯！想跟我一起比赛算术吗？",
        teacher: "你好孩子，数学是通往未来的大门。",
        scientist: "观察这些建筑的比例，数学无处不在。",
        shopkeeper: "今天所有的圆柱形商品都打八折哦！"
      };
      setDialogue({ name: npc.name, message: messages[npc.type] || "你好呀！" });
    }
  };

  // 监听接近NPC状态的变化 - 只显示UI效果，不自动弹出对话
  useEffect(() => {
    // 此效果仅用于软件UI效果，实际对话弹出显示由handleCharacterClick控制
  }, [nearbyNPC]);

  // 当离开NPC范围时自动关闭对话框
  useEffect(() => {
    if (!nearbyNPC && dialogue) {
      setDialogue(null);
      // 离开范围时重置冷却时间
      setDialogueCooldown(0);
    }
  }, [nearbyNPC]);

  // 监听范围外NPC点击事件，触发小人移动
  useEffect(() => {
    const handleNPCClick = (e) => {
      const { position } = e.detail;
      // 点击NPC时，设置目标位置为NPC所在位置
      setTargetPos(new THREE.Vector3(position[0], 0, position[2]));
      setShowMarker(true);
      setTimeout(() => setShowMarker(false), 1500);
    };
    
    document.addEventListener('npcClicked', handleNPCClick);
    return () => document.removeEventListener('npcClicked', handleNPCClick);
  }, []);

  // 冷却时间倒计时
  useEffect(() => {
    if (dialogueCooldown > 0) {
      const timer = setTimeout(() => {
        setDialogueCooldown(prev => Math.max(0, prev - 0.1));
      }, 100);
      return () => clearTimeout(timer);
    }
  }, [dialogueCooldown]);

  // 关闭对话框并启动冷却
  const closeDialogue = () => {
    setDialogue(null);
    // 设置0.5秒冷却时间，防止连续点击
    setDialogueCooldown(0.5);
  };

  return (
    <div style={{ width: '100vw', height: '100vh', position: 'relative', overflow: 'hidden' }}>
      <div style={{ position: 'absolute', top: 20, left: 20, zIndex: 10, pointerEvents: 'none' }}>
        <h2 style={{ color: '#2C3E50', margin: 0, textShadow: '2px 2px white' }}>数学卡通小镇</h2>
        <p style={{ color: '#34495E' }}><b>点击地面</b> 或 使用 <b>WASD</b> 走动，寻找 NPC 对话</p>
      </div>

      <Canvas shadows>
        <Suspense fallback={null}>
          {/* 天空背景 - 使用 Drei 的 Sky 组件 */}
          <Sky distance={450000} sunPosition={[100, 20, 100]} intensity={0.8} turbidity={10} rayleigh={2} />
          
          {/* 星星效果 */}
          <Stars radius={100} depth={50} count={5000} factor={4} saturation={0} fade speed={1} />
          
          {/* 环境光和阴影 */}
          <ambientLight intensity={0.9} />
          {/* 动态方向光（因为会根据时间旋转） */}
          <DynamicDirectionalLight />
          
          {/* 补光 - 从另一个方向的间接光 */}
          <pointLight position={[-30, 10, -30]} intensity={0.4} />
          
          <Mountains />
          <Bush position={[-6, 0.5, -10]} color="#27AE60" />
          <Bush position={[7, 0.5, -15]} color="#2ECC71" />

          {/* 建筑物 */}
          <Building position={[-15, 2, -10]} width={4} height={4} depth={3} color="#E74C3C" />
          <Building position={[10, 2, -15]} width={3.5} height={4} depth={3} color="#F39C12" />
          <Building position={[-8, 2, 5]} width={3} height={4} depth={2.5} color="#9B59B6" />
          <Building position={[15, 2, 10]} width={3.5} height={4} depth={3} color="#E74C3C" />
          
          {/* 旋转的装饰物 */}
          <RotatingCube position={[-15, 5.5, -10]} size={0.6} color="#3498DB" />
          <RotatingCube position={[10, 5.5, -15]} size={0.6} color="#2ECC71" />
          <RotatingCube position={[-8, 5.5, 5]} size={0.6} color="#F39C12" />
          <RotatingCube position={[15, 5.5, 10]} size={0.6} color="#3498DB" />
          
          {/* 树木 */}
          <Tree position={[-20, 0, 15]} />
          <Tree position={[18, 0, -20]} />
          
          {/* 水体（湖/河） */}
          <Water position={[0, 0, 25]} width={20} depth={15} />
          
          {/* 农场区域 */}
          <Farm position={[-30, 0, -5]} />
          
          {/* 牞场区域 */}
          <Pasture position={[30, 0, 5]} />
          
          {/* 装饰花朵 */}
          <Flower position={[-25, 0, 10]} color="#FF1493" />
          <Flower position={[-35, 0, 15]} color="#FF69B4" />
          <Flower position={[25, 0, 10]} color="#FFB6C1" />
          <Flower position={[35, 0, 15]} color="#FF1493" />
          
          {/* 装饰石头 */}
          <Rock position={[-20, 0, 5]} size={0.6} color="#A0826D" />
          <Rock position={[20, 0, -5]} size={0.5} color="#8B7355" />
          <Rock position={[-35, 0, -15]} size={0.7} color="#C0B0A0" />
          <Rock position={[35, 0, -10]} size={0.6} color="#A0826D" />
          
          {/* 装饰花丛盆栽 */}
          <SmallPlanter position={[-15, 0, 20]} color="#CD853F" />
          <SmallPlanter position={[15, 0, 20]} color="#DAA520" />
          <SmallPlanter position={[-40, 0, 0]} color="#8B6F47" />
          <SmallPlanter position={[40, 0, 0]} color="#CD853F" />
          
          <Street onGroundClick={handleGroundClick} />
          {showMarker && <TargetMarker position={targetPos} />}
          
          <NPC type="child" name="小明" position={[-3, 0, -5]} onClick={handleCharacterClick} isNearby={nearbyNPC?.id === 1} playerPos={playerPos} />
          <NPC type="teacher" name="王老师" position={[3, 0, -12]} onClick={handleCharacterClick} isNearby={nearbyNPC?.id === 2} playerPos={playerPos} />
          <NPC type="shopkeeper" name="李老板" position={[6, 0, 2]} onClick={handleCharacterClick} isNearby={nearbyNPC?.id === 3} playerPos={playerPos} />

      <Player onMove={(pos) => setPlayerPos(pos.clone())} targetPos={targetPos} />
          <FollowCamera targetPos={playerPos} />
        </Suspense>
      </Canvas>

        {dialogue && (
        <div style={{
          position: 'absolute', bottom: 50, left: '50%', transform: 'translateX(-50%)',
          backgroundColor: 'white', padding: '20px', borderRadius: '15px',
          boxShadow: '0 5px 15px rgba(0,0,0,0.2)', border: '4px solid #3498DB', width: '80%', maxWidth: '600px'
        }}>
          <b style={{ color: '#3498DB', fontSize: '1.2em' }}>{dialogue.name}:</b>
          <p style={{ fontSize: '1.1em', margin: '10px 0' }}>{dialogue.message}</p>
          <button onClick={closeDialogue} style={{
            backgroundColor: '#3498DB', color: 'white', border: 'none', padding: '10px 20px', borderRadius: '8px', cursor: 'pointer'
          }}>确定</button>
        </div>
      )}  
      {nearbyNPC && !dialogue && (
        <div style={{
          position: 'absolute', bottom: 100, left: '50%', transform: 'translateX(-50%)',
          backgroundColor: '#FFD700', padding: '10px 20px', borderRadius: '8px',
          boxShadow: '0 3px 10px rgba(0,0,0,0.2)', color: '#333', fontSize: '1.1em',
          fontWeight: 'bold', animation: 'pulse 1s infinite'
        }}>
          💬 靠近 {nearbyNPC.name} - 点击交互!
          <style>{`
            @keyframes pulse {
              0%, 100% { opacity: 1; }
              50% { opacity: 0.7; }
            }
          `}</style>
        </div>
      )}
    </div>
  );
};

// 简单的 NPC 组件（基于之前的卡通形象）
const NPC = ({ type, position, name, onClick, isNearby, playerPos }) => {
  const [hovered, setHovered] = useState(false);
  
  // 计算与玩家的距离
  const npcPos = new THREE.Vector3(...position);
  const playerVector = new THREE.Vector3(...(playerPos ? [playerPos.x, playerPos.y, playerPos.z] : [0, 0, 0]));
  const distanceToPlayer = npcPos.distanceTo(playerVector);
  const isWithinInteractionRange = distanceToPlayer < 3;
  
  // 处理点击：在范围内时触发对话，范围外时产生移动效果
  const handleClick = (e) => {
    e.stopPropagation();
    if (isWithinInteractionRange) {
      // 在范围内：触发对话
      onClick(name, type);
    } else {
      // 范围外：产生移动效果（点击该位置移动小人过去）
      const moveTarget = new THREE.Vector3(...position);
      // 这里假设有一个全局的移动函数，或者通过其他方式传递
      // 为了简化，我们直接在这里处理
      document.dispatchEvent(new CustomEvent('npcClicked', { detail: { position } }));
    }
  };
    
  // 防止事件穿透到地面
  const handlePointerDown = (e) => {
    e.stopPropagation();
  };
  
  return (
    <group position={position} 
           onClick={handleClick}
           onPointerDown={handlePointerDown}
           onPointerOver={() => setHovered(true)}
           onPointerOut={() => setHovered(false)}>
      {/* 接近时的发光背景光圈 */}
      {isNearby && (
        <mesh position={[0, 0.8, 0]}>
          <cylinderGeometry args={[1.5, 1.5, 0.1, 32]} />
          <meshBasicMaterial color="#FFD700" transparent opacity={0.3} />
        </mesh>
      )}
      
      <mesh position={[0, 1.6, 0]} castShadow>
        <sphereGeometry args={[0.45, 16, 16]} />
        <meshStandardMaterial color={hovered || isNearby ? "#FFEB3B" : "#FFE0BD"} />
      </mesh>
      <mesh position={[0, 0.8, 0]} castShadow>
        <cylinderGeometry args={[0.3, 0.4, 1, 16]} />
        <meshStandardMaterial color={type === 'teacher' ? "#2ECC71" : "#E74C3C"} />
      </mesh>
      <Text position={[0, 2.4, 0]} fontSize={0.4} color={isNearby ? "#FFD700" : "#2C3E50"} anchorX="center">{name}</Text>
      {/* 范围外时显示点击可移动提示 */}
      {!isNearby && (hovered || isWithinInteractionRange === false) && (
        <Text position={[0, 3.2, 0]} fontSize={0.25} color="#888888" anchorX="center" anchorY="middle">
          点击移动
        </Text>
      )}
    </group>
  );
};

export default MathCityScene;