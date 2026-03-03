import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

// 学习闭环的各个阶段
const STAGES = [
  { id: 'teach', name: '讲解', icon: '📖' },
  { id: 'practice', name: '练习', icon: '✏️' },
  { id: 'test', name: '测验', icon: '📝' },
  { id: 'diagnose', name: '诊断', icon: '🔍' },
  { id: 'achievement', name: '成就', icon: '🎉' },
];

// 模拟学习内容
const learningContent = {
  fractions: {
    domainName: '数与代数',
    trunkName: '分数',
    currentLevel: 3,
    nodes: [
      { id: 'frac-def', name: '分数的定义', type: 'C', level: 1, status: 'completed' },
      { id: 'frac-compare', name: '分数比较', type: 'C', level: 2, status: 'completed' },
      { id: 'frac-add', name: '同分母加减', type: 'S', level: 2, status: 'completed' },
      { id: 'frac-same-denom', name: '通分', type: 'S', level: 3, status: 'current' },
      { id: 'frac-diff-denom', name: '异分母加减', type: 'S', level: 3, status: 'locked' },
      { id: 'frac-problem', name: '分数应用题', type: 'P', level: 4, status: 'locked' },
    ]
  }
};

// 当前节点的练习题
const sampleQuestions = [
  {
    id: 1,
    question: '将 1/2 和 1/3 通分后，它们的分母是？',
    options: ['2', '3', '6', '12'],
    correct: 2,
    explanation: '2 和 3 的最小公倍数是 6，所以通分后分母是 6。',
  },
  {
    id: 2,
    question: '2/4 和 3/6 通分后相等吗？',
    options: ['相等', '不相等', '无法确定'],
    correct: 0,
    explanation: '2/4 = 1/2，3/6 = 1/2，所以它们相等！',
  },
];

const StageIndicator = ({ stages, currentStage }) => (
  <div style={{
    display: 'flex',
    justifyContent: 'space-between',
    padding: '16px',
    backgroundColor: 'white',
    borderRadius: '12px',
    marginBottom: '16px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.05)',
  }}>
    {stages.map((stage, index) => (
      <div key={stage.id} style={{ textAlign: 'center', flex: 1 }}>
        <div style={{
          width: '36px',
          height: '36px',
          borderRadius: '50%',
          backgroundColor: index <= currentStage ? '#3B82F6' : '#E5E7EB',
          color: index <= currentStage ? 'white' : '#9CA3AF',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          margin: '0 auto 4px',
          fontSize: '16px',
        }}>
          {index < currentStage ? '✓' : stage.icon}
        </div>
        <span style={{
          fontSize: '11px',
          color: index <= currentStage ? '#3B82F6' : '#9CA3AF',
        }}>
          {stage.name}
        </span>
      </div>
    ))}
  </div>
);

const TeachStage = ({ node, onNext }) => (
  <div style={{ padding: '20px', textAlign: 'center' }}>
    <div style={{
      width: '120px',
      height: '120px',
      borderRadius: '20px',
      backgroundColor: '#FEF3C7',
      margin: '0 auto 20px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontSize: '48px',
    }}>
      📐
    </div>
    <h2 style={{ margin: '0 0 12px', fontSize: '22px', color: '#1F2937' }}>{node.name}</h2>
    <p style={{ fontSize: '14px', color: '#6B7280', lineHeight: '1.6', marginBottom: '20px' }}>
      把一个整体平均分成若干份，表示其中一份或几份的数叫做分数。
      <br /><br />
      例如：把一个披萨平均分成 4 份，拿走 1 份就是 1/4。
    </p>
    <div style={{
      backgroundColor: '#ECFDF5',
      padding: '16px',
      borderRadius: '12px',
      textAlign: 'left',
    }}>
      <h4 style={{ margin: '0 0 8px', fontSize: '14px', color: '#059669' }}>💡 为什么要通分？</h4>
      <p style={{ margin: 0, fontSize: '13px', color: '#047857' }}>
        当我们需要比较两个分数大小，或者进行分数加减时，如果分母不同，
        就需要把它们变成相同的分母，这就是"通分"。
      </p>
    </div>
    <button
      onClick={onNext}
      style={{
        marginTop: '24px',
        padding: '14px 32px',
        backgroundColor: '#3B82F6',
        color: 'white',
        border: 'none',
        borderRadius: '12px',
        fontSize: '16px',
        fontWeight: '500',
        cursor: 'pointer',
      }}
    >
      开始练习 →
    </button>
  </div>
);

const PracticeStage = ({ questions, onComplete }) => {
  const [currentQ, setCurrentQ] = useState(0);
  const [selected, setSelected] = useState(null);
  const [showResult, setShowResult] = useState(false);

  const q = questions[currentQ];

  const handleSelect = (idx) => {
    if (showResult) return;
    setSelected(idx);
  };

  const handleSubmit = () => {
    if (selected === null) return;
    setShowResult(true);
  };

  const handleNext = () => {
    if (currentQ < questions.length - 1) {
      setCurrentQ(currentQ + 1);
      setSelected(null);
      setShowResult(false);
    } else {
      onComplete();
    }
  };

  const isCorrect = selected === q.correct;

  return (
    <div style={{ padding: '20px' }}>
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '20px',
      }}>
        <span style={{ fontSize: '14px', color: '#6B7280' }}>
          练习 {currentQ + 1}/{questions.length}
        </span>
        <span style={{ fontSize: '14px', color: '#3B82F6' }}>正确率: 0%</span>
      </div>

      <h3 style={{ fontSize: '18px', color: '#1F2937', marginBottom: '20px' }}>{q.question}</h3>

      <div style={{ marginBottom: '20px' }}>
        {q.options.map((opt, idx) => (
          <div
            key={idx}
            onClick={() => handleSelect(idx)}
            style={{
              padding: '16px',
              marginBottom: '10px',
              backgroundColor: showResult
                ? idx === q.correct ? '#D1FAE5' : idx === selected ? '#FEE2E2' : '#F9FAFB'
                : selected === idx ? '#DBEAFE' : '#F9FAFB',
              border: `2px solid ${selected === idx ? '#3B82F6' : '#E5E7EB'}`,
              borderRadius: '12px',
              cursor: showResult ? 'default' : 'pointer',
              transition: '0.2s',
            }}
          >
            <span style={{
              display: 'inline-block',
              width: '24px',
              height: '24px',
              borderRadius: '50%',
              backgroundColor: idx === q.correct ? '#10B981' : '#E5E7EB',
              color: 'white',
              textAlign: 'center',
              lineHeight: '24px',
              marginRight: '12px',
              fontSize: '12px',
            }}>
              {['A', 'B', 'C', 'D'][idx]}
            </span>
            {opt}
          </div>
        ))}
      </div>

      {showResult && (
        <div style={{
          backgroundColor: isCorrect ? '#D1FAE5' : '#FEE2E2',
          padding: '16px',
          borderRadius: '12px',
          marginBottom: '16px',
        }}>
          <p style={{
            margin: 0,
            color: isCorrect ? '#065F46' : '#991B1B',
            fontSize: '14px',
          }}>
            <strong>{isCorrect ? '✅ 回答正确！' : '❌ 回答错误'}</strong>
            <br />
            {q.explanation}
          </p>
        </div>
      )}

      {!showResult ? (
        <button
          onClick={handleSubmit}
          disabled={selected === null}
          style={{
            width: '100%',
            padding: '14px',
            backgroundColor: selected !== null ? '#3B82F6' : '#9CA3AF',
            color: 'white',
            border: 'none',
            borderRadius: '12px',
            fontSize: '16px',
            cursor: selected !== null ? 'pointer' : 'not-allowed',
          }}
        >
          提交答案
        </button>
      ) : (
        <button
          onClick={handleNext}
          style={{
            width: '100%',
            padding: '14px',
            backgroundColor: '#3B82F6',
            color: 'white',
            border: 'none',
            borderRadius: '12px',
            fontSize: '16px',
            cursor: 'pointer',
          }}
        >
          {currentQ < questions.length - 1 ? '下一题' : '完成练习 →'}
        </button>
      )}
    </div>
  );
};

const DiagnoseStage = ({ onFinish }) => (
  <div style={{ padding: '20px', textAlign: 'center' }}>
    <div style={{
      width: '80px',
      height: '80px',
      borderRadius: '50%',
      backgroundColor: '#D1FAE5',
      margin: '0 auto 20px',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontSize: '36px',
    }}>
      🔍
    </div>
    <h2 style={{ margin: '0 0 12px', fontSize: '20px', color: '#1F2937' }}>诊断结果</h2>
    
    <div style={{
      backgroundColor: 'white',
      borderRadius: '16px',
      padding: '20px',
      textAlign: 'left',
      marginBottom: '16px',
    }}>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '12px' }}>
        <span style={{ fontSize: '24px', marginRight: '12px' }}>📊</span>
        <div>
          <div style={{ fontWeight: '500', color: '#1F2937' }}>掌握度评估</div>
          <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#10B981' }}>85%</div>
        </div>
      </div>
      
      <div style={{ marginTop: '16px' }}>
        <div style={{ fontSize: '13px', color: '#6B7280', marginBottom: '8px' }}>能力 Breakdown</div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <AbilityBar label="概念[C]" value={90} color="#3B82F6" />
          <AbilityBar label="技能[S]" value={75} color="#10B981" />
          <AbilityBar label="思维[T]" value={80} color="#8B5CF6" />
        </div>
      </div>
    </div>

    <div style={{
      backgroundColor: '#FEF3C7',
      borderRadius: '12px',
      padding: '16px',
      textAlign: 'left',
      marginBottom: '20px',
    }}>
      <div style={{ fontWeight: '500', color: '#92400E', marginBottom: '8px' }}>
        💡 学习建议
      </div>
      <div style={{ fontSize: '13px', color: '#B45309' }}>
        你的通分概念掌握得很好！在异分母加减计算时还需加强练习。
        建议完成"倍数小齿轮"支线任务来强化底层技能。
      </div>
    </div>

    <button
      onClick={onFinish}
      style={{
        width: '100%',
        padding: '14px',
        backgroundColor: '#3B82F6',
        color: 'white',
        border: 'none',
        borderRadius: '12px',
        fontSize: '16px',
        cursor: 'pointer',
      }}
    >
      继续学习 →
    </button>
  </div>
);

const AbilityBar = ({ label, value, color }) => (
  <div style={{ flex: 1 }}>
    <div style={{ fontSize: '10px', color: '#6B7280', marginBottom: '4px' }}>{label}</div>
    <div style={{ height: '6px', backgroundColor: '#E5E7EB', borderRadius: '3px' }}>
      <div style={{
        width: value + '%',
        height: '100%',
        backgroundColor: color,
        borderRadius: '3px',
      }} />
    </div>
  </div>
);

const AchievementStage = ({ onFinish }) => (
  <div style={{ padding: '40px 20px', textAlign: 'center' }}>
    <div style={{
      fontSize: '64px',
      marginBottom: '20px',
      animation: 'bounce 1s infinite',
    }}>
      🏆
    </div>
    <h2 style={{ margin: '0 0 12px', fontSize: '22px', color: '#1F2937' }}>恭喜完成本节学习！</h2>
    
    <div style={{
      backgroundColor: '#FEF3C7',
      borderRadius: '16px',
      padding: '20px',
      marginBottom: '20px',
    }}>
      <div style={{ fontSize: '48px', marginBottom: '12px' }}>🎖️</div>
      <div style={{ fontWeight: 'bold', color: '#92400E', fontSize: '16px' }}>分数 Lv3 通关</div>
      <div style={{ fontSize: '13px', color: '#B45309', marginTop: '4px' }}>
        获得 50 经验值
      </div>
    </div>

    <button
      onClick={onFinish}
      style={{
        width: '100%',
        padding: '14px',
        backgroundColor: '#10B981',
        color: 'white',
        border: 'none',
        borderRadius: '12px',
        fontSize: '16px',
        cursor: 'pointer',
      }}
    >
      返回知识地图
    </button>
  </div>
);

const LearningPage = () => {
  const { domainId, trunkId } = useParams();
  const navigate = useNavigate();
  const [stage, setStage] = useState(0); // 0:讲解 1:练习 2:测验(跳过) 3:诊断 4:成就

  const content = learningContent[trunkId] || learningContent.fractions;
  const currentNode = content.nodes.find(n => n.status === 'current') || content.nodes[0];

  const handleNext = () => setStage(stage + 1);
  const handleComplete = () => setStage(3); // 跳到诊断
  const handleFinish = () => navigate('/knowledge-map');

  const renderStage = () => {
    switch (stage) {
      case 0:
        return <TeachStage node={currentNode} onNext={handleNext} />;
      case 1:
        return <PracticeStage questions={sampleQuestions} onComplete={handleComplete} />;
      case 3:
        return <DiagnoseStage onFinish={handleFinish} />;
      case 4:
        return <AchievementStage onFinish={handleFinish} />;
      default:
        return null;
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: '#F9FAFB',
    }}>
      {/* Header */}
      <div style={{
        position: 'sticky',
        top: 0,
        backgroundColor: 'white',
        padding: '12px 16px',
        borderBottom: '1px solid #E5E7EB',
        display: 'flex',
        alignItems: 'center',
        zIndex: 10,
      }}>
        <button
          onClick={() => navigate('/knowledge-map')}
          style={{
            background: 'none',
            border: 'none',
            fontSize: '20px',
            cursor: 'pointer',
            padding: '4px',
          }}
        >
          ←
        </button>
        <div style={{ marginLeft: '8px' }}>
          <div style={{ fontSize: '14px', color: '#6B7280' }}>{content.domainName}</div>
          <div style={{ fontSize: '16px', fontWeight: '500', color: '#1F2937' }}>{content.trunkName}</div>
        </div>
      </div>

      {/* Stage Indicator */}
      <div style={{ padding: '16px 16px 0' }}>
        <StageIndicator stages={STAGES} currentStage={stage} />
      </div>

      {/* Stage Content */}
      <div style={{ marginTop: '8px' }}>
        {renderStage()}
      </div>
    </div>
  );
};

export default LearningPage;
