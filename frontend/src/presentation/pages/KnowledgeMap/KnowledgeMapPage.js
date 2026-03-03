import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// 模拟数据：数学领域
const domains = [
  {
    id: 'numbers',
    name: '数与代数',
    icon: '🔢',
    emoji: '🔢',
    color: '#3B82F6',
    bgGradient: 'linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)',
    description: '数字、运算、代数式',
    decoration: '＋−×÷',
    trunks: [
      { id: 'fractions', name: '分数', level: 3.2, maxLevel: 5, color: '#60A5FA' },
      { id: 'multiples', name: '倍数与因数', level: 2.8, maxLevel: 5, color: '#93C5FD' },
      { id: 'equations', name: '方程', level: 1.5, maxLevel: 5, color: '#DBEAFE' },
    ]
  },
  {
    id: 'geometry',
    name: '图形与几何',
    icon: '📐',
    emoji: '📐',
    color: '#10B981',
    bgGradient: 'linear-gradient(135deg, #10B981 0%, #34D399 100%)',
    description: '图形、测量、空间',
    decoration: '△□◇○',
    trunks: [
      { id: 'triangles', name: '三角形', level: 2.5, maxLevel: 5, color: '#34D399' },
      { id: 'areas', name: '面积', level: 2.0, maxLevel: 5, color: '#6EE7B7' },
      { id: 'angles', name: '角度', level: 1.0, maxLevel: 5, color: '#A7F3D0' },
    ]
  },
  {
    id: 'statistics',
    name: '统计与概率',
    icon: '📊',
    emoji: '📊',
    color: '#F59E0B',
    bgGradient: 'linear-gradient(135deg, #F59E0B 0%, #FBBF24 100%)',
    description: '数据收集、统计图表',
    decoration: '📈📉📊',
    trunks: [
      { id: 'averages', name: '平均数', level: 1.5, maxLevel: 5, color: '#FBBF24' },
      { id: 'charts', name: '统计图表', level: 1.0, maxLevel: 5, color: '#FCD34D' },
    ]
  },
  {
    id: 'application',
    name: '综合与实践',
    icon: '🧩',
    emoji: '🧩',
    color: '#8B5CF6',
    bgGradient: 'linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%)',
    description: '问题解决、跨学科',
    decoration: '❓💡✨',
    trunks: [
      { id: 'modeling', name: '数学建模', level: 1.0, maxLevel: 5, color: '#A78BFA' },
      { id: 'puzzles', name: '数学游戏', level: 2.0, maxLevel: 5, color: '#C4B5FD' },
    ]
  }
];

const LevelIndicator = ({ level, maxLevel, compact }) => {
  const stars = Math.floor(level);
  const hasHalfStar = level % 1 >= 0.5;
  
  if (compact) {
    // 紧凑模式：小圆点
    return (
      <div style={{ display: 'flex', alignItems: 'center', gap: '3px' }}>
        {[...Array(maxLevel)].map((_, i) => (
          <span
            key={i}
            style={{
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              backgroundColor: i < level ? '#FFD700' : '#E5E7EB',
              border: '1px solid #D1D5DB',
            }}
          />
        ))}
        <span style={{ marginLeft: '6px', fontSize: '11px', color: '#6B7280' }}>
          Lv{level.toFixed(1)}
        </span>
      </div>
    );
  }
  
  // 普通模式
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '4px' }}>
      {[...Array(maxLevel)].map((_, i) => (
        <span
          key={i}
          style={{
            width: '12px',
            height: '12px',
            borderRadius: '50%',
            backgroundColor: i < level ? '#FFD700' : '#E5E7EB',
            border: '1px solid #D1D5DB',
          }}
        />
      ))}
      <span style={{ marginLeft: '8px', fontSize: '12px', color: '#6B7280' }}>
        Lv{level.toFixed(1)}
      </span>
    </div>
  );
};

const DomainCard = ({ domain, isExpanded, onClick, onSelectTrunk, listMode }) => {
  const [expanded, setExpanded] = useState(false);
  
  const isExpandedFinal = listMode ? true : (isExpanded || expanded);

  // 计算平均等级
  const avgLevel = domain.trunks.reduce((sum, t) => sum + t.level, 0) / domain.trunks.length;

  if (listMode) {
    // 列表模式样式
    return (
      <div style={{
        backgroundColor: 'white',
        borderRadius: '16px',
        padding: '16px',
        marginBottom: '12px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
        border: '1px solid #F3F4F6',
      }}>
        <div
          onClick={() => setExpanded(!expanded)}
          style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}
        >
          <div style={{
            width: '48px', height: '48px', borderRadius: '12px',
            backgroundColor: domain.color + '20', display: 'flex',
            alignItems: 'center', justifyContent: 'center', fontSize: '24px', marginRight: '12px'
          }}>
            {domain.icon}
          </div>
          <div style={{ flex: 1 }}>
            <h3 style={{ margin: 0, fontSize: '16px', color: '#1F2937' }}>{domain.name}</h3>
            <p style={{ margin: '4px 0 0', fontSize: '12px', color: '#6B7280' }}>{domain.description}</p>
          </div>
          <span style={{ fontSize: '20px', transform: expanded ? 'rotate(180deg)' : 'rotate(0)', transition: '0.2s' }}>▼</span>
        </div>

        {expanded && (
          <div style={{ marginTop: '16px', paddingTop: '16px', borderTop: '1px solid #F3F4F6' }}>
            {domain.trunks.map(trunk => (
              <div
                key={trunk.id}
                onClick={() => onSelectTrunk(domain.id, trunk)}
                style={{
                  display: 'flex', alignItems: 'center', padding: '12px',
                  marginBottom: '8px', backgroundColor: trunk.color + '15',
                  borderRadius: '10px', cursor: 'pointer'
                }}
              >
                <div style={{ flex: 1 }}>
                  <span style={{ fontSize: '14px', fontWeight: '500', color: '#374151' }}>{trunk.name}</span>
                  <LevelIndicator level={trunk.level} maxLevel={trunk.maxLevel} />
                </div>
                <button style={{
                  padding: '6px 12px', backgroundColor: '#3B82F6', color: 'white',
                  border: 'none', borderRadius: '6px', fontSize: '12px', cursor: 'pointer'
                }}>
                  开始学习 →
                </button>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  // 探索模式（卡片样式）
  return (
    <div
      onClick={onClick}
      style={{
        background: domain.bgGradient,
        borderRadius: '20px',
        padding: '20px',
        minHeight: '180px',
        cursor: 'pointer',
        transition: 'all 0.3s ease',
        transform: isExpanded ? 'scale(1.02)' : 'scale(1)',
        boxShadow: isExpanded 
          ? `0 20px 40px ${domain.color}40` 
          : '0 4px 12px rgba(0,0,0,0.1)',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      {/* 装饰元素 */}
      <div style={{
        position: 'absolute',
        top: '-10px',
        right: '-10px',
        fontSize: '80px',
        opacity: 0.1,
        color: 'white',
        fontWeight: 'bold',
      }}>
        {domain.emoji}
      </div>
      <div style={{
        position: 'absolute',
        bottom: '-20px',
        left: '-20px',
        fontSize: '100px',
        opacity: 0.05,
        color: 'white',
      }}>
        {domain.decoration}
      </div>

      {/* 内容 */}
      <div style={{ position: 'relative', zIndex: 1 }}>
        <div style={{
          fontSize: '40px',
          marginBottom: '12px',
          filter: 'drop-shadow(0 2px 4px rgba(0,0,0,0.2))',
        }}>
          {domain.emoji}
        </div>
        <h3 style={{
          margin: 0,
          fontSize: '18px',
          color: 'white',
          fontWeight: 'bold',
          textShadow: '0 1px 2px rgba(0,0,0,0.2)',
        }}>
          {domain.name}
        </h3>
        <p style={{
          margin: '4px 0 0',
          fontSize: '12px',
          color: 'rgba(255,255,255,0.8)',
        }}>
          {domain.description}
        </p>

        {/* 进度条 */}
        <div style={{
          marginTop: '16px',
          backgroundColor: 'rgba(255,255,255,0.2)',
          borderRadius: '8px',
          padding: '8px',
        }}>
          <div style={{
            display: 'flex',
            justifyContent: 'space-between',
            marginBottom: '4px',
          }}>
            <span style={{ fontSize: '11px', color: 'rgba(255,255,255,0.9)' }}>
              平均等级
            </span>
            <span style={{ fontSize: '11px', color: 'white', fontWeight: 'bold' }}>
              Lv{avgLevel.toFixed(1)}
            </span>
          </div>
          <div style={{
            height: '6px',
            backgroundColor: 'rgba(255,255,255,0.3)',
            borderRadius: '3px',
          }}>
            <div style={{
              width: (avgLevel / 5 * 100) + '%',
              height: '100%',
              backgroundColor: '#FFD700',
              borderRadius: '3px',
            }} />
          </div>
        </div>

        {/* 技能线数量 */}
        <div style={{
          marginTop: '12px',
          display: 'flex',
          gap: '6px',
          flexWrap: 'wrap',
        }}>
          {domain.trunks.slice(0, 3).map(trunk => (
            <span key={trunk.id} style={{
              padding: '4px 8px',
              backgroundColor: 'rgba(255,255,255,0.2)',
              borderRadius: '12px',
              fontSize: '10px',
              color: 'white',
            }}>
              {trunk.name}
            </span>
          ))}
          {domain.trunks.length > 3 && (
            <span style={{
              padding: '4px 8px',
              backgroundColor: 'rgba(255,255,255,0.2)',
              borderRadius: '12px',
              fontSize: '10px',
              color: 'white',
            }}>
              +{domain.trunks.length - 3}
            </span>
          )}
        </div>
      </div>

      {/* 展开时显示技能线列表 */}
      {isExpanded && (
        <div style={{
          marginTop: '16px',
          paddingTop: '16px',
          borderTop: '1px solid rgba(255,255,255,0.3)',
        }}>
          {domain.trunks.map(trunk => (
            <div
              key={trunk.id}
              onClick={(e) => { e.stopPropagation(); onSelectTrunk(domain.id, trunk); }}
              style={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                padding: '10px 12px',
                marginBottom: '6px',
                backgroundColor: 'rgba(255,255,255,0.9)',
                borderRadius: '10px',
                cursor: 'pointer',
                transition: '0.2s',
              }}
            >
              <span style={{ fontSize: '13px', color: '#374151', fontWeight: '500' }}>
                {trunk.name}
              </span>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <LevelIndicator level={trunk.level} maxLevel={trunk.maxLevel} compact />
                <span style={{
                  fontSize: '10px',
                  color: '#3B82F6',
                  backgroundColor: '#DBEAFE',
                  padding: '2px 6px',
                  borderRadius: '8px',
                }}>
                  开始
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

const KnowledgeMapPage = () => {
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState('cards'); // 'list' or 'cards'
  const [selectedDomain, setSelectedDomain] = useState(null);

  const handleSelectTrunk = (domainId, trunk) => {
    navigate(`/learning/${domainId}/${trunk.id}`);
  };

  const handleDomainClick = (domain) => {
    if (viewMode === 'cards') {
      setSelectedDomain(selectedDomain?.id === domain.id ? null : domain);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#F0F9FF',
      backgroundImage: 'radial-gradient(circle at 20% 80%, rgba(59, 130, 246, 0.1) 0%, transparent 50%), radial-gradient(circle at 80% 20%, rgba(16, 185, 129, 0.1) 0%, transparent 50%)',
    }}>
      {/* Header */}
      <header style={{
        padding: '16px',
        paddingTop: '48px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}>
        <div>
          <h1 style={{ margin: 0, fontSize: '28px', color: '#1F2937', fontWeight: 'bold' }}>
            🗺️ 知识地图
          </h1>
          <p style={{ margin: '4px 0 0', fontSize: '14px', color: '#6B7280' }}>
            选择一个数学世界开始探险
          </p>
        </div>
        
        {/* Mode Switcher */}
        <div style={{
          display: 'flex',
          backgroundColor: 'white',
          borderRadius: '12px',
          padding: '4px',
          boxShadow: '0 2px 4px rgba(0,0,0,0.05)',
        }}>
          <button
            onClick={() => { setViewMode('cards'); setSelectedDomain(null); }}
            style={{
              padding: '8px 12px',
              borderRadius: '8px',
              border: 'none',
              backgroundColor: viewMode === 'cards' ? '#3B82F6' : 'transparent',
              color: viewMode === 'cards' ? 'white' : '#6B7280',
              fontSize: '12px',
              cursor: 'pointer',
              transition: '0.2s',
            }}
          >
            🎴 探索模式
          </button>
          <button
            onClick={() => { setViewMode('list'); setSelectedDomain(null); }}
            style={{
              padding: '8px 12px',
              borderRadius: '8px',
              border: 'none',
              backgroundColor: viewMode === 'list' ? '#3B82F6' : 'transparent',
              color: viewMode === 'list' ? 'white' : '#6B7280',
              fontSize: '12px',
              cursor: 'pointer',
              transition: '0.2s',
            }}
          >
            📋 列表模式
          </button>
        </div>
      </header>

      <div style={{ padding: '16px', paddingTop: '0px' }}>
        {/* Cards Mode */}
        {viewMode === 'cards' && (
          <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(2, 1fr)',
            gap: '16px',
          }}>
            {domains.map(domain => (
              <DomainCard
                key={domain.id}
                domain={domain}
                isExpanded={selectedDomain?.id === domain.id}
                onClick={() => handleDomainClick(domain)}
                onSelectTrunk={handleSelectTrunk}
              />
            ))}
          </div>
        )}

        {/* List Mode */}
        {viewMode === 'list' && (
          <div>
            {domains.map(domain => (
              <DomainCard
                key={domain.id}
                domain={domain}
                isExpanded={true}
                onClick={() => {}}
                onSelectTrunk={handleSelectTrunk}
                listMode={true}
              />
            ))}
          </div>
        )}

        {/* 3D 场景入口 */}
        <div style={{ marginTop: '24px' }}>
          <div style={{
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            borderRadius: '20px',
            padding: '24px',
            color: 'white',
            position: 'relative',
            overflow: 'hidden',
          }}>
            {/* 装饰 */}
            <div style={{
              position: 'absolute',
              top: '-30px',
              right: '-30px',
              fontSize: '120px',
              opacity: 0.1,
            }}>🌟</div>
            
            <div style={{ position: 'relative', zIndex: 1 }}>
              <h3 style={{ margin: '0 0 8px', fontSize: '18px', fontWeight: 'bold' }}>
                🌟 进入 3D 数学世界
              </h3>
              <p style={{ margin: '0 0 16px', fontSize: '13px', opacity: 0.9 }}>
                沉浸式探索数学知识，发现隐藏的奥秘
              </p>
              
              <div style={{ display: 'flex', gap: '12px' }}>
                <button
                  onClick={() => navigate('/3d/town')}
                  style={{
                    flex: 1,
                    padding: '14px 16px',
                    backgroundColor: 'rgba(255,255,255,0.2)',
                    backdropFilter: 'blur(10px)',
                    border: '1px solid rgba(255,255,255,0.3)',
                    borderRadius: '14px',
                    color: 'white',
                    cursor: 'pointer',
                    transition: '0.2s',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '10px',
                  }}
                >
                  <span style={{ fontSize: '28px' }}>🏰</span>
                  <div style={{ textAlign: 'left' }}>
                    <div style={{ fontWeight: 'bold', fontSize: '14px' }}>知识小镇</div>
                    <div style={{ fontSize: '11px', opacity: 0.8 }}>进阶探索</div>
                  </div>
                </button>
                
                <button
                  onClick={() => navigate('/3d/village')}
                  style={{
                    flex: 1,
                    padding: '14px 16px',
                    backgroundColor: 'rgba(255,255,255,0.2)',
                    backdropFilter: 'blur(10px)',
                    border: '1px solid rgba(255,255,255,0.3)',
                    borderRadius: '14px',
                    color: 'white',
                    cursor: 'pointer',
                    transition: '0.2s',
                    display: 'flex',
                    alignItems: 'center',
                    gap: '10px',
                  }}
                >
                  <span style={{ fontSize: '28px' }}>🏡</span>
                  <div style={{ textAlign: 'left' }}>
                    <div style={{ fontWeight: 'bold', fontSize: '14px' }}>知识村庄</div>
                    <div style={{ fontSize: '11px', opacity: 0.8 }}>入门探索</div>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeMapPage;
