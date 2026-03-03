import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// 成就数据
const achievements = [
  {
    id: 1,
    name: '分数 Lv3 通关',
    description: '完成分数第三章学习',
    icon: '🎖️',
    category: 'knowledge',
    unlocked: true,
    unlockedAt: '2026-02-15',
    color: '#3B82F6',
  },
  {
    id: 2,
    name: '探索小达人',
    description: '探索了5个不同的数学世界',
    icon: '🧭',
    category: 'exploration',
    unlocked: true,
    unlockedAt: '2026-02-10',
    color: '#10B981',
  },
  {
    id: 3,
    name: '数学小天才',
    description: '连续答对10道题',
    icon: '🧠',
    category: 'challenge',
    unlocked: true,
    unlockedAt: '2026-02-08',
    color: '#F59E0B',
  },
  {
    id: 4,
    name: '倍数大师',
    description: '完成倍数与因数主线',
    icon: '✖️',
    category: 'knowledge',
    unlocked: false,
    progress: 75,
    color: '#8B5CF6',
  },
  {
    id: 5,
    name: '几何探索者',
    description: '探索所有几何图形',
    icon: '📐',
    category: 'exploration',
    unlocked: false,
    progress: 40,
    color: '#EC4899',
  },
  {
    id: 6,
    name: '坚持不懈',
    description: '连续学习7天',
    icon: '📅',
    category: 'challenge',
    unlocked: true,
    unlockedAt: '2026-02-18',
    color: '#14B8A6',
  },
  {
    id: 7,
    name: '全对超人',
    description: '一次测验全部正确',
    icon: '💯',
    category: 'challenge',
    unlocked: true,
    unlockedAt: '2026-02-12',
    color: '#EF4444',
  },
  {
    id: 8,
    name: '故事探险家',
    description: '完成所有数学历史故事',
    icon: '📚',
    category: 'exploration',
    unlocked: false,
    progress: 60,
    color: '#6366F1',
  },
];

const categories = [
  { id: 'all', name: '全部', icon: '🏆' },
  { id: 'knowledge', name: '知识', icon: '📖' },
  { id: 'exploration', name: '探索', icon: '🧭' },
  { id: 'challenge', name: '挑战', icon: '⚡' },
];

const StatsCard = () => (
  <div style={{
    backgroundColor: 'white',
    borderRadius: '16px',
    padding: '20px',
    marginBottom: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
  }}>
    <div style={{ display: 'flex', justifyContent: 'space-around', textAlign: 'center' }}>
      <div>
        <div style={{ fontSize: '28px', fontWeight: 'bold', color: '#3B82F6' }}>1,250</div>
        <div style={{ fontSize: '12px', color: '#6B7280' }}>总经验值</div>
      </div>
      <div>
        <div style={{ fontSize: '28px', fontWeight: 'bold', color: '#10B981' }}>12</div>
        <div style={{ fontSize: '12px', color: '#6B7280' }}>已解锁</div>
      </div>
      <div>
        <div style={{ fontSize: '28px', fontWeight: 'bold', color: '#F59E0B' }}>18</div>
        <div style={{ fontSize: '12px', color: '#6B7280' }}>进行中</div>
      </div>
    </div>
  </div>
);

const AchievementCard = ({ achievement, onClick }) => (
  <div
    onClick={() => onClick(achievement)}
    style={{
      backgroundColor: achievement.unlocked ? 'white' : '#F9FAFB',
      borderRadius: '12px',
      padding: '16px',
      marginBottom: '10px',
      boxShadow: achievement.unlocked ? '0 2px 8px rgba(0,0,0,0.06)' : 'none',
      border: achievement.unlocked ? '1px solid #E5E7EB' : '1px dashed #D1D5DB',
      opacity: achievement.unlocked ? 1 : 0.6,
      cursor: 'pointer',
      transition: '0.2s',
    }}
  >
    <div style={{ display: 'flex', alignItems: 'center' }}>
      <div style={{
        width: '50px',
        height: '50px',
        borderRadius: '12px',
        backgroundColor: achievement.unlocked ? achievement.color + '20' : '#E5E7EB',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: '24px',
        marginRight: '12px',
      }}>
        {achievement.unlocked ? achievement.icon : '🔒'}
      </div>
      <div style={{ flex: 1 }}>
        <div style={{
          fontWeight: '500',
          color: achievement.unlocked ? '#1F2937' : '#9CA3AF',
        }}>
          {achievement.name}
        </div>
        <div style={{ fontSize: '12px', color: '#6B7280' }}>{achievement.description}</div>
        {achievement.unlocked ? (
          <div style={{ fontSize: '11px', color: '#9CA3AF', marginTop: '4px' }}>
            {achievement.unlockedAt} 解锁
          </div>
        ) : (
          <div style={{ marginTop: '8px' }}>
            <div style={{
              height: '4px',
              backgroundColor: '#E5E7EB',
              borderRadius: '2px',
            }}>
              <div style={{
                width: achievement.progress + '%',
                height: '100%',
                backgroundColor: achievement.color,
                borderRadius: '2px',
              }} />
            </div>
            <div style={{ fontSize: '10px', color: '#9CA3AF', marginTop: '2px' }}>
              {achievement.progress}% 进度
            </div>
          </div>
        )}
      </div>
    </div>
  </div>
);

const AchievementsPage = () => {
  const navigate = useNavigate();
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [selectedAchievement, setSelectedAchievement] = useState(null);

  const filteredAchievements = selectedCategory === 'all'
    ? achievements
    : achievements.filter(a => a.category === selectedCategory);

  const unlockedCount = achievements.filter(a => a.unlocked).length;

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#F9FAFB',
      padding: '16px',
    }}>
      <header style={{ marginBottom: '20px', paddingTop: '40px' }}>
        <h1 style={{ margin: 0, fontSize: '24px', color: '#1F2937' }}>🏆 成就中心</h1>
        <p style={{ margin: '4px 0 0', fontSize: '14px', color: '#6B7280' }}>
          已解锁 {unlockedCount}/{achievements.length} 个成就
        </p>
      </header>

      <StatsCard />

      {/* Category Tabs */}
      <div style={{
        display: 'flex',
        gap: '8px',
        marginBottom: '16px',
        overflowX: 'auto',
        paddingBottom: '4px',
      }}>
        {categories.map(cat => (
          <button
            key={cat.id}
            onClick={() => setSelectedCategory(cat.id)}
            style={{
              padding: '8px 16px',
              borderRadius: '20px',
              border: 'none',
              backgroundColor: selectedCategory === cat.id ? '#3B82F6' : 'white',
              color: selectedCategory === cat.id ? 'white' : '#6B7280',
              fontSize: '13px',
              whiteSpace: 'nowrap',
              cursor: 'pointer',
              boxShadow: selectedCategory === cat.id ? 'none' : '0 1px 2px rgba(0,0,0,0.05)',
            }}
          >
            {cat.icon} {cat.name}
          </button>
        ))}
      </div>

      {/* Achievement List */}
      <div>
        {filteredAchievements.map(achievement => (
          <AchievementCard
            key={achievement.id}
            achievement={achievement}
            onClick={setSelectedAchievement}
          />
        ))}
      </div>

      {/* Achievement Detail Modal */}
      {selectedAchievement && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.5)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          zIndex: 100,
        }} onClick={() => setSelectedAchievement(null)}>
          <div style={{
            backgroundColor: 'white',
            borderRadius: '20px',
            padding: '24px',
            margin: '20px',
            maxWidth: '320px',
            textAlign: 'center',
          }} onClick={e => e.stopPropagation()}>
            <div style={{
              width: '80px',
              height: '80px',
              borderRadius: '20px',
              backgroundColor: selectedAchievement.unlocked 
                ? selectedAchievement.color + '20' 
                : '#F3F4F6',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: '40px',
              margin: '0 auto 16px',
            }}>
              {selectedAchievement.unlocked ? selectedAchievement.icon : '🔒'}
            </div>
            <h3 style={{ margin: '0 0 8px', fontSize: '18px', color: '#1F2937' }}>
              {selectedAchievement.name}
            </h3>
            <p style={{ margin: '0 0 16px', fontSize: '14px', color: '#6B7280' }}>
              {selectedAchievement.description}
            </p>
            {selectedAchievement.unlocked ? (
              <div style={{
                backgroundColor: '#D1FAE5',
                padding: '12px',
                borderRadius: '10px',
                fontSize: '13px',
                color: '#065F46',
              }}>
                ✅ 已于 {selectedAchievement.unlockedAt} 解锁
              </div>
            ) : (
              <div>
                <div style={{
                  height: '8px',
                  backgroundColor: '#E5E7EB',
                  borderRadius: '4px',
                  marginBottom: '8px',
                }}>
                  <div style={{
                    width: selectedAchievement.progress + '%',
                    height: '100%',
                    backgroundColor: selectedAchievement.color,
                    borderRadius: '4px',
                  }} />
                </div>
                <div style={{ fontSize: '13px', color: '#6B7280' }}>
                  继续加油！还差 {100 - selectedAchievement.progress}% 达成
                </div>
              </div>
            )}
            <button
              onClick={() => setSelectedAchievement(null)}
              style={{
                marginTop: '20px',
                padding: '10px 24px',
                backgroundColor: '#F3F4F6',
                border: 'none',
                borderRadius: '8px',
                fontSize: '14px',
                color: '#4B5563',
                cursor: 'pointer',
              }}
            >
              关闭
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default AchievementsPage;
