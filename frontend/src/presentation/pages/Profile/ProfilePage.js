import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// 用户模拟数据
const userData = {
  name: '小明',
  age: 8,
  avatar: '👦',
  level: 5,
  exp: 1250,
  nextLevelExp: 2000,
  streak: 7,
  totalDays: 45,
  achievements: 12,
  learningStats: {
    totalTime: '12h 30m',
    thisWeek: '3h 45m',
    nodesCompleted: 28,
    questionsAnswered: 156,
    correctRate: 85,
  },
  settings: {
    sound: true,
    music: true,
    vibration: true,
    ageMode: '7-8岁',
  },
};

const StatItem = ({ icon, label, value }) => (
  <div style={{
    textAlign: 'center',
    padding: '12px',
    backgroundColor: '#F9FAFB',
    borderRadius: '12px',
  }}>
    <div style={{ fontSize: '20px', marginBottom: '4px' }}>{icon}</div>
    <div style={{ fontSize: '18px', fontWeight: 'bold', color: '#1F2937' }}>{value}</div>
    <div style={{ fontSize: '11px', color: '#6B7280' }}>{label}</div>
  </div>
);

const MenuItem = ({ icon, label, value, onClick, danger }) => (
  <div
    onClick={onClick}
    style={{
      display: 'flex',
      alignItems: 'center',
      padding: '16px',
      backgroundColor: 'white',
      borderBottom: '1px solid #F3F4F6',
      cursor: 'pointer',
    }}
  >
    <span style={{ fontSize: '20px', marginRight: '12px' }}>{icon}</span>
    <span style={{ flex: 1, color: danger ? '#EF4444' : '#1F2937' }}>{label}</span>
    {value && (
      <span style={{ color: '#9CA3AF', fontSize: '14px' }}>{value}</span>
    )}
    <span style={{ color: '#D1D5DB', marginLeft: '8px' }}>›</span>
  </div>
);

const ProfilePage = () => {
  const navigate = useNavigate();
  const [showAgeModal, setShowAgeModal] = useState(false);

  const expProgress = (userData.exp / userData.nextLevelExp) * 100;

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#F9FAFB',
    }}>
      {/* Profile Header */}
      <div style={{
        backgroundColor: '#3B82F6',
        padding: '40px 20px 80px',
        borderRadius: '0 0 30px 30px',
        color: 'white',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
          <div style={{
            width: '80px',
            height: '80px',
            borderRadius: '50%',
            backgroundColor: 'white',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '40px',
          }}>
            {userData.avatar}
          </div>
          <div>
            <h1 style={{ margin: 0, fontSize: '24px' }}>{userData.name}</h1>
            <div style={{ fontSize: '14px', opacity: 0.9 }}>Lv.{userData.level} · {userData.age}岁</div>
          </div>
        </div>

        {/* XP Progress */}
        <div style={{
          marginTop: '24px',
          backgroundColor: 'rgba(255,255,255,0.2)',
          borderRadius: '10px',
          padding: '12px',
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '6px', fontSize: '12px' }}>
            <span>经验值</span>
            <span>{userData.exp} / {userData.nextLevelExp}</span>
          </div>
          <div style={{
            height: '8px',
            backgroundColor: 'rgba(255,255,255,0.3)',
            borderRadius: '4px',
          }}>
            <div style={{
              width: expProgress + '%',
              height: '100%',
              backgroundColor: '#FFD700',
              borderRadius: '4px',
            }} />
          </div>
        </div>
      </div>

      {/* Stats Grid */}
      <div style={{
        marginTop: '-50px',
        padding: '0 16px',
      }}>
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(4, 1fr)',
          gap: '8px',
          marginBottom: '16px',
        }}>
          <StatItem icon="🔥" label="连续天数" value={userData.streak} />
          <StatItem icon="📅" label="学习天数" value={userData.totalDays} />
          <StatItem icon="🏅" label="成就" value={userData.achievements} />
          <StatItem icon="📈" label="正确率" value={userData.learningStats.correctRate + '%'} />
        </div>
      </div>

      {/* Learning Stats */}
      <div style={{
        margin: '16px',
        backgroundColor: 'white',
        borderRadius: '16px',
        padding: '16px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.04)',
      }}>
        <h3 style={{ margin: '0 0 12px', fontSize: '14px', color: '#6B7280' }}>学习数据</h3>
        <div style={{ display: 'flex', justifyContent: 'space-between' }}>
          <div style={{ fontSize: '13px', color: '#6B7280' }}>总学习时长</div>
          <div style={{ fontWeight: '500', color: '#1F2937' }}>{userData.learningStats.totalTime}</div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '8px' }}>
          <div style={{ fontSize: '13px', color: '#6B7280' }}>本周学习</div>
          <div style={{ fontWeight: '500', color: '#1F2937' }}>{userData.learningStats.thisWeek}</div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '8px' }}>
          <div style={{ fontSize: '13px', color: '#6B7280' }}>完成知识点</div>
          <div style={{ fontWeight: '500', color: '#1F2937' }}>{userData.learningStats.nodesCompleted}个</div>
        </div>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '8px' }}>
          <div style={{ fontSize: '13px', color: '#6B7280' }}>答题数量</div>
          <div style={{ fontWeight: '500', color: '#1F2937' }}>{userData.learningStats.questionsAnswered}题</div>
        </div>
      </div>

      {/* Settings */}
      <div style={{
        margin: '16px',
        backgroundColor: 'white',
        borderRadius: '16px',
        overflow: 'hidden',
        boxShadow: '0 2px 8px rgba(0,0,0,0.04)',
      }}>
        <h3 style={{ margin: '0', padding: '16px 16px 8px', fontSize: '14px', color: '#6B7280' }}>设置</h3>
        
        <MenuItem
          icon="🎂"
          label="年龄设置"
          value={userData.settings.ageMode}
          onClick={() => setShowAgeModal(true)}
        />
        <MenuItem
          icon="🔊"
          label="音效"
          value={userData.settings.sound ? '开' : '关'}
        />
        <MenuItem
          icon="🎵"
          label="背景音乐"
          value={userData.settings.music ? '开' : '关'}
        />
        <MenuItem
          icon="📳"
          label="振动反馈"
          value={userData.settings.vibration ? '开' : '关'}
        />
      </div>

      {/* Menu */}
      <div style={{
        margin: '16px',
        backgroundColor: 'white',
        borderRadius: '16px',
        overflow: 'hidden',
        boxShadow: '0 2px 8px rgba(0,0,0,0.04)',
      }}>
        <MenuItem
          icon="👨‍👩‍👧"
          label="切换到家长端"
          onClick={() => navigate('/parent')}
        />
        <MenuItem
          icon="❓"
          label="帮助与反馈"
        />
        <MenuItem
          icon="ℹ️"
          label="关于我们"
        />
      </div>

      <div style={{ height: '40px' }} />

      {/* Age Selection Modal */}
      {showAgeModal && (
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
        }} onClick={() => setShowAgeModal(false)}>
          <div style={{
            backgroundColor: 'white',
            borderRadius: '20px',
            padding: '24px',
            margin: '20px',
            width: '300px',
          }} onClick={e => e.stopPropagation()}>
            <h3 style={{ margin: '0 0 16px', textAlign: 'center', color: '#1F2937' }}>
              选择年龄
            </h3>
            {['3-4岁', '5-6岁', '7-8岁', '9-10岁', '11-12岁'].map(age => (
              <div
                key={age}
                style={{
                  padding: '14px',
                  marginBottom: '8px',
                  backgroundColor: userData.settings.ageMode === age ? '#DBEAFE' : '#F9FAFB',
                  borderRadius: '10px',
                  textAlign: 'center',
                  cursor: 'pointer',
                  color: userData.settings.ageMode === age ? '#1E40AF' : '#4B5563',
                }}
              >
                {age}
              </div>
            ))}
            <button
              onClick={() => setShowAgeModal(false)}
              style={{
                marginTop: '12px',
                width: '100%',
                padding: '12px',
                backgroundColor: '#F3F4F6',
                border: 'none',
                borderRadius: '10px',
                cursor: 'pointer',
              }}
            >
              取消
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProfilePage;
