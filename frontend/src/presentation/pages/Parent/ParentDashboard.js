import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

// 模拟孩子学习数据
const childData = {
  name: '小明',
  grade: '三年级',
  lastLogin: '2026-02-21',
  activeLines: [
    { id: 'fractions', name: '分数技能线', level: 3.2, maxLevel: 5, color: '#3B82F6' },
    { id: 'multiples', name: '倍数与因数', level: 2.8, maxLevel: 5, color: '#10B981' },
  ],
  recentAchievements: [
    { name: '分数 Lv3 通关', date: '2026-02-15', icon: '🎖️' },
    { name: '坚持不懈', date: '2026-02-18', icon: '📅' },
  ],
  weeklyData: [
    { day: '周一', minutes: 25 },
    { day: '周二', minutes: 40 },
    { day: '周三', minutes: 15 },
    { day: '周四', minutes: 35 },
    { day: '周五', minutes: 50 },
    { day: '周六', minutes: 60 },
    { day: '周日', minutes: 30 },
  ],
  weakPoints: [
    {
      node: '[S] 最小公倍数求法',
      mastery: 50,
      relatedLine: '分数技能线',
      suggestion: '系统已推送"倍数小齿轮工坊"支线任务',
    },
    {
      node: '[P] 异分母分数加减',
      mastery: 40,
      relatedLine: '分数技能线',
      suggestion: '需要加强通分技能的练习',
    },
  ],
};

const MaxMinutes = 70;

const ChildInfoBar = () => (
  <div style={{
    backgroundColor: 'white',
    borderRadius: '16px',
    padding: '16px',
    marginBottom: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
    display: 'flex',
    justifyContent: 'space-between',
    alignItems: 'center',
  }}>
    <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
      <div style={{
        width: '48px',
        height: '48px',
        borderRadius: '50%',
        backgroundColor: '#FEF3C7',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: '24px',
      }}>
        👦
      </div>
      <div>
        <div style={{ fontWeight: 'bold', color: '#1F2937' }}>{childData.name}</div>
        <div style={{ fontSize: '13px', color: '#6B7280' }}>{childData.grade} · 上次登录: {childData.lastLogin}</div>
      </div>
    </div>
    <button
      onClick={() => {}}
      style={{
        padding: '8px 12px',
        backgroundColor: '#F3F4F6',
        border: 'none',
        borderRadius: '8px',
        fontSize: '12px',
        color: '#6B7280',
        cursor: 'pointer',
      }}
    >
      切换孩子
    </button>
  </div>
);

const ProgressCard = ({ title, line, level, maxLevel, color }) => (
  <div style={{
    backgroundColor: 'white',
    borderRadius: '12px',
    padding: '14px',
    marginBottom: '10px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
  }}>
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
      <span style={{ fontSize: '13px', color: '#6B7280' }}>{title}</span>
      <span style={{ fontSize: '14px', fontWeight: 'bold', color }}>Lv{level} / {maxLevel}</span>
    </div>
    <div style={{ height: '8px', backgroundColor: '#E5E7EB', borderRadius: '4px' }}>
      <div style={{
        width: (level / maxLevel * 100) + '%',
        height: '100%',
        backgroundColor: color,
        borderRadius: '4px',
      }} />
    </div>
  </div>
);

const WeeklyChart = () => (
  <div style={{
    backgroundColor: 'white',
    borderRadius: '16px',
    padding: '16px',
    marginBottom: '16px',
    boxShadow: '0 2px 8px rgba(0,0,0,0.06)',
  }}>
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
      <h3 style={{ margin: 0, fontSize: '15px', color: '#1F2937' }}>本周学习时长</h3>
      <span style={{ fontSize: '13px', color: '#6B7280' }}>255分钟</span>
    </div>
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', height: '100px' }}>
      {childData.weeklyData.map((d, i) => (
        <div key={i} style={{ textAlign: 'center', flex: 1 }}>
          <div style={{
            height: (d.minutes / MaxMinutes * 80) + 'px',
            backgroundColor: '#3B82F6',
            borderRadius: '4px 4px 0 0',
            margin: '0 4px',
            minHeight: '4px',
          }} />
          <div style={{ fontSize: '10px', color: '#9CA3AF', marginTop: '4px' }}>{d.day}</div>
        </div>
      ))}
    </div>
  </div>
);

const WeakPointCard = ({ point }) => (
  <div style={{
    backgroundColor: '#FEF3C7',
    borderRadius: '12px',
    padding: '14px',
    marginBottom: '10px',
    borderLeft: '4px solid #F59E0B',
  }}>
    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '6px' }}>
      <span style={{ fontWeight: '500', color: '#92400E', fontSize: '14px' }}>{point.node}</span>
      <span style={{ fontSize: '12px', color: '#B45309', backgroundColor: '#FDE68A', padding: '2px 8px', borderRadius: '10px' }}>
        掌握度 {point.mastery}%
      </span>
    </div>
    <div style={{ fontSize: '12px', color: '#B45309' }}>关联: {point.relatedLine}</div>
    <div style={{ fontSize: '12px', color: '#059669', marginTop: '6px', backgroundColor: '#D1FAE5', padding: '6px 8px', borderRadius: '6px' }}>
      💡 {point.suggestion}
    </div>
  </div>
);

const AchievementItem = ({ achievement }) => (
  <div style={{
    display: 'flex',
    alignItems: 'center',
    padding: '10px',
    backgroundColor: '#F9FAFB',
    borderRadius: '10px',
    marginBottom: '8px',
  }}>
    <span style={{ fontSize: '24px', marginRight: '10px' }}>{achievement.icon}</span>
    <div style={{ flex: 1 }}>
      <div style={{ fontSize: '13px', color: '#1F2937' }}>{achievement.name}</div>
      <div style={{ fontSize: '11px', color: '#9CA3AF' }}>{achievement.date}</div>
    </div>
  </div>
);

const ParentDashboard = () => {
  const navigate = useNavigate();
  const [activeTab, setActiveTab] = useState('overview');

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#F9FAFB',
      padding: '16px',
    }}>
      {/* Header */}
      <header style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: '16px',
        paddingTop: '40px',
      }}>
        <div>
          <h1 style={{ margin: 0, fontSize: '22px', color: '#1F2937' }}>👨‍👩‍👧 家长端</h1>
          <p style={{ margin: '4px 0 0', fontSize: '13px', color: '#6B7280' }}>孩子的学习情况</p>
        </div>
        <button
          onClick={() => navigate('/profile')}
          style={{
            padding: '8px 12px',
            backgroundColor: 'white',
            border: '1px solid #E5E7EB',
            borderRadius: '8px',
            fontSize: '12px',
            color: '#6B7280',
            cursor: 'pointer',
          }}
        >
          切换到孩子端
        </button>
      </header>

      {/* Tab Switcher */}
      <div style={{
        display: 'flex',
        gap: '8px',
        marginBottom: '16px',
      }}>
        {[
          { id: 'overview', label: '总览' },
          { id: 'skills', label: '技能线' },
          { id: 'reports', label: '报告' },
        ].map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            style={{
              flex: 1,
              padding: '10px',
              borderRadius: '10px',
              border: 'none',
              backgroundColor: activeTab === tab.id ? '#3B82F6' : 'white',
              color: activeTab === tab.id ? 'white' : '#6B7280',
              fontSize: '14px',
              cursor: 'pointer',
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {activeTab === 'overview' && (
        <>
          <ChildInfoBar />

          {/* Active Progress */}
          <div style={{ marginBottom: '16px' }}>
            <h3 style={{ margin: '0 0 8px', fontSize: '14px', color: '#6B7280' }}>最近在学习</h3>
            {childData.activeLines.map(line => (
              <ProgressCard
                key={line.id}
                title={line.name}
                level={line.level}
                maxLevel={line.maxLevel}
                color={line.color}
              />
            ))}
          </div>

          <WeeklyChart />

          {/* Recent Achievements */}
          <div style={{ marginBottom: '16px' }}>
            <h3 style={{ margin: '0 0 8px', fontSize: '14px', color: '#6B7280' }}>近期成就</h3>
            {childData.recentAchievements.map((a, i) => (
              <AchievementItem key={i} achievement={a} />
            ))}
          </div>

          {/* Weak Points Warning */}
          <div>
            <h3 style={{ margin: '0 0 8px', fontSize: '14px', color: '#6B7280' }}>需要关注</h3>
            {childData.weakPoints.map((point, i) => (
              <WeakPointCard key={i} point={point} />
            ))}
          </div>
        </>
      )}

      {activeTab === 'skills' && (
        <div style={{ padding: '20px', textAlign: 'center', color: '#6B7280' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>📊</div>
          <p>技能线详情页面</p>
          <p style={{ fontSize: '13px' }}>点击技能线卡片查看详细进度和节点掌握情况</p>
        </div>
      )}

      {activeTab === 'reports' && (
        <div style={{ padding: '20px', textAlign: 'center', color: '#6B7280' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>📋</div>
          <p>专题报告页面</p>
          <p style={{ fontSize: '13px' }}>查看诊断报告和成长报告</p>
          <button style={{
            marginTop: '16px',
            padding: '10px 20px',
            backgroundColor: '#3B82F6',
            color: 'white',
            border: 'none',
            borderRadius: '8px',
            cursor: 'pointer',
          }}>
            生成新报告
          </button>
        </div>
      )}

      <div style={{ height: '40px' }} />
    </div>
  );
};

export default ParentDashboard;
