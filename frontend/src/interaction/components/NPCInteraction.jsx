import React, { useState } from 'react';
import { Canvas, useThree } from '@react-three/fiber';
import { Text, OrbitControls } from '@react-three/drei';
import * as THREE from 'three';
import SoundManager from '../../shared/utils/SoundManager';

// NPC 组件
const NPC = ({ name, characterType, onClick, position = [0, 0, 0] }) => {
  const [hovered, setHovered] = useState(false);
  const [clickEffect, setClickEffect] = useState(false);

  // 根据角色类型设置不同颜色
  const getColorByType = (type) => {
    switch(type) {
      case 'farmer':
        return '#4682B4'; // 钢蓝色
      case 'merchant':
        return '#8B4513'; // 棕色
      case 'teacher':
        return '#32CD32'; // 酸橙绿
      case 'explorer':
        return '#FF6347'; // 番茄红
      default:
        return '#4682B4';
    }
  };

  const handleClick = (e) => {
    e.stopPropagation(); // 防止事件冒泡
    setClickEffect(true);
    SoundManager.playSound('npc_click');
    onClick(name, characterType);
    setTimeout(() => setClickEffect(false), 500);
  };

  const handlePointerOver = () => {
    setHovered(true);
    document.body.style.cursor = 'pointer';
  };

  const handlePointerOut = () => {
    setHovered(false);
    document.body.style.cursor = 'default';
  };

  const color = getColorByType(characterType);

  return (
    <group position={position}>
      {/* 点击效果 */}
      {clickEffect && (
        <mesh position={[0, 2.5, 0]}>
          <sphereGeometry args={[0.3, 16, 16]} />
          <meshBasicMaterial color="#FFFF00" transparent opacity={0.7} />
        </mesh>
      )}
      
      {/* NPC 身体 */}
      <mesh 
        position={[0, 1, 0]} 
        onClick={handleClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <cylinderGeometry args={[0.4, 0.4, 2, 16]} />
        <meshStandardMaterial 
          color={hovered ? '#FFD700' : color} 
          roughness={0.2} 
          metalness={0.1} 
        />
      </mesh>
      
      {/* NPC 头部 */}
      <mesh 
        position={[0, 2.2, 0]} 
        onClick={handleClick}
        onPointerOver={handlePointerOver}
        onPointerOut={handlePointerOut}
      >
        <sphereGeometry args={[0.5, 16, 16]} />
        <meshStandardMaterial 
          color={hovered ? '#FFD700' : '#FFDAB9'} 
          roughness={0.3} 
          metalness={0.2} 
        />
      </mesh>
      
      {/* NPC 名字标签 */}
      <Text
        position={[0, 3, 0]}
        fontSize={0.4}
        color={hovered ? '#FFD700' : 'white'}
        anchorX="center"
        anchorY="middle"
      >
        {name}
      </Text>
      
      {/* 角色类型标识 */}
      <Text
        position={[0, 2.7, 0]}
        fontSize={0.2}
        color="#FFA500"
        anchorX="center"
        anchorY="middle"
      >
        {characterType === 'farmer' ? '🌾' : 
         characterType === 'merchant' ? '💰' : 
         characterType === 'teacher' ? '📚' : '🎒'}
      </Text>
    </group>
  );
};

// 对话框组件 (2D UI)
const DialogueBox = ({ message, npcName, onClose, options = [] }) => {
  if (!message) return null;

  return (
    <div
      style={{
        position: 'absolute',
        bottom: '20%',
        left: '50%',
        transform: 'translateX(-50%)',
        width: '60%',
        maxWidth: '500px',
        backgroundColor: 'rgba(0, 0, 0, 0.8)',
        color: 'white',
        padding: '20px',
        borderRadius: '10px',
        border: '2px solid #ccc',
        zIndex: 100,
        fontFamily: 'Arial, sans-serif',
      }}
    >
      <div style={{ marginBottom: '10px' }}>
        <strong style={{ color: '#FFD700' }}>{npcName}:</strong> {message}
      </div>
      
      {options.length > 0 && (
        <div style={{ marginTop: '10px' }}>
          {options.map((option, idx) => (
            <button
              key={idx}
              onClick={() => option.action && option.action()}
              style={{
                display: 'block',
                width: '100%',
                padding: '8px',
                margin: '5px 0',
                backgroundColor: '#4682B4',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer',
              }}
            >
              {option.text}
            </button>
          ))}
        </div>
      )}
      
      <button 
        onClick={onClose} 
        style={{ 
          marginTop: '10px', 
          padding: '5px 10px',
          backgroundColor: '#666',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          cursor: 'pointer',
        }}
      >
        关闭
      </button>
    </div>
  );
};

const NPCInteraction = () => {
  const [dialogue, setDialogue] = useState(null);
  const [activeNPC, setActiveNPC] = useState(null);

  const handleNPCClick = (name, characterType) => {
    setActiveNPC({ name, characterType });
    
    // 根据不同角色类型提供不同的对话内容
    let message = '';
    let options = [];
    
    switch(characterType) {
      case 'farmer':
        message = "你好小朋友！我家的菜园子需要帮忙规划一下面积，你能帮我吗？";
        options = [
          { text: "好的，我来帮你！", action: () => {
              SoundManager.playSound('success_sound');
              setDialogue({ message: "太好了！我们来规划一下，这个菜园是长方形的，长10米宽6米，你知道面积是多少吗？", npcName: name });
            }},
          { text: "我不太会，可以教我吗？", action: () => {
              setDialogue({ message: "当然可以！长方形面积等于长乘以宽。来，我们一起计算吧！", npcName: name });
            }}
        ];
        break;
      case 'merchant':
        message = "哎呀，我这批货的折扣算错了，能帮帮我吗？";
        options = [
          { text: "折扣怎么算呢？", action: () => {
              setDialogue({ message: "折扣是原价减去打折后的价格。比如原价100元打8折，就是100×0.8=80元，折扣是20元。", npcName: name });
            }},
          { text: "让我试试计算", action: () => {
              SoundManager.playSound('math_problem');
              setDialogue({ message: "好的！假设一件商品原价是120元，打7.5折，你能算出折扣金额吗？", npcName: name });
            }}
        ];
        break;
      case 'teacher':
        message = "欢迎来到数学课堂！今天我们学习分数的概念。";
        options = [
          { text: "什么是分数？", action: () => {
              setDialogue({ message: "分数表示一个整体被平均分成若干份，取其中的几份。比如1/2表示把一个整体平均分成2份，取其中1份。", npcName: name });
            }},
          { text: "分数怎么比较大小？", action: () => {
              setDialogue({ message: "当分母相同时，分子越大分数越大；当分子相同时，分母越小分数越大。比如1/3 > 1/4。", npcName: name });
            }}
        ];
        break;
      case 'explorer':
        message = "你好！我在寻找隐藏的数学宝藏，需要你的帮助。";
        options = [
          { text: "数学宝藏？听起来很有趣！", action: () => {
              setDialogue({ message: "是的！宝藏被藏在坐标(3, 4)处，你需要通过解数学谜题才能到达那里。准备好了吗？", npcName: name });
            }},
          { text: "坐标是什么？", action: () => {
              setDialogue({ message: "坐标是用来确定平面上点的位置的。比如(3, 4)表示从原点向右3格，向上4格。", npcName: name });
            }}
        ];
        break;
      default:
        message = "你好！很高兴见到你。";
    }
    
    setDialogue({ message, npcName: name, options });
  };

  const closeDialogue = () => {
    setDialogue(null);
    setActiveNPC(null);
  };

  return (
    <div style={{ width: '100vw', height: '100vh', position: 'relative' }}>
      <Canvas camera={{ position: [8, 8, 8], fov: 50 }}>
        <ambientLight intensity={0.5} />
        <pointLight position={[10, 10, 10]} intensity={1} />
        <spotLight position={[10, 10, 0]} angle={0.15} penumbra={1} intensity={1} castShadow />
        
        {/* 地面 */}
        <mesh rotation={[-Math.PI / 2, 0, 0]} receiveShadow>
          <planeGeometry args={[20, 20]} />
          <meshStandardMaterial color="#90EE90" roughness={0.8} metalness={0.2} />
        </mesh>
        
        {/* 添加一些环境装饰 */}
        {/* 树木 */}
        <group position={[-5, 0, -5]}>
          <mesh position={[0, 3, 0]}>
            <sphereGeometry args={[1.5, 16, 16]} />
            <meshStandardMaterial color="#2E8B57" />
          </mesh>
          <mesh position={[0, 1, 0]}>
            <cylinderGeometry args={[0.3, 0.3, 2, 8]} />
            <meshStandardMaterial color="#8B4513" />
          </mesh>
        </group>
        
        <group position={[6, 0, -4]}>
          <mesh position={[0, 2.5, 0]}>
            <sphereGeometry args={[1.2, 16, 16]} />
            <meshStandardMaterial color="#3CB371" />
          </mesh>
          <mesh position={[0, 1, 0]}>
            <cylinderGeometry args={[0.25, 0.25, 2, 8]} />
            <meshStandardMaterial color="#8B4513" />
          </mesh>
        </group>
        
        {/* 房屋 */}
        <group position={[0, 0, 5]}>
          <mesh position={[0, 1.5, 0]}>
            <boxGeometry args={[4, 3, 4]} />
            <meshStandardMaterial color="#D2B48C" />
          </mesh>
          <mesh position={[0, 3.2, 0]} rotation={[0, 0, Math.PI / 4]}>
            <coneGeometry args={[3, 2, 4]} />
            <meshStandardMaterial color="#A52A2A" />
          </mesh>
        </group>
        
        {/* NPCs */}
        <NPC 
          name="农夫老王" 
          characterType="farmer" 
          onClick={handleNPCClick} 
          position={[-3, 0, 0]} 
        />
        <NPC 
          name="商人李姐" 
          characterType="merchant" 
          onClick={handleNPCClick} 
          position={[3, 0, 0]} 
        />
        <NPC 
          name="张老师" 
          characterType="teacher" 
          onClick={handleNPCClick} 
          position={[0, 0, -3]} 
        />
        <NPC 
          name="探险家小刘" 
          characterType="explorer" 
          onClick={handleNPCClick} 
          position={[0, 0, 3]} 
        />
        
        {/* 添加相机控制器 */}
        <OrbitControls enablePan={true} enableZoom={true} enableRotate={true} />
      </Canvas>
      
      <DialogueBox 
        message={dialogue?.message} 
        npcName={dialogue?.npcName} 
        options={dialogue?.options}
        onClose={closeDialogue} 
      />
    </div>
  );
};

export default NPCInteraction;