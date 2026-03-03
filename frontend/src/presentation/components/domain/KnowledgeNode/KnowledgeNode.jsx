import React from 'react';

const KnowledgeNode = ({ 
  title, 
  description, 
  level, 
  type, // C: Concept, S: Support Skill, T: Thinking Pattern, P: Problem Model
  completed, 
  masteryScore,
  onClick 
}) => {
  // Define colors for different node types
  const typeStyles = {
    'C': { // Concept
      bgColor: 'bg-green-50',
      borderColor: 'border-green-200',
      textColor: 'text-green-800',
      icon: '🧠',
      title: '概念'
    },
    'S': { // Support Skill
      bgColor: 'bg-blue-50',
      borderColor: 'border-blue-200',
      textColor: 'text-blue-800',
      icon: '⚙️',
      title: '支撑技能'
    },
    'T': { // Thinking Pattern
      bgColor: 'bg-yellow-50',
      borderColor: 'border-yellow-200',
      textColor: 'text-yellow-800',
      icon: '💡',
      title: '思维模式'
    },
    'P': { // Problem Model
      bgColor: 'bg-purple-50',
      borderColor: 'border-purple-200',
      textColor: 'text-purple-800',
      icon: '🧩',
      title: '问题模型'
    }
  };

  // Define colors for different levels
  const levelColors = {
    'Lv1': 'bg-green-100 text-green-800',
    'Lv2': 'bg-blue-100 text-blue-800',
    'Lv3': 'bg-yellow-100 text-yellow-800',
    'Lv4': 'bg-orange-100 text-orange-800',
    'Lv5': 'bg-red-100 text-red-800',
  };

  // Get styles based on type
  const nodeStyle = typeStyles[type] || typeStyles['C'];

  // Calculate mastery percentage
  const masteryPercentage = masteryScore !== undefined ? Math.round(masteryScore * 100) : null;

  return (
    <div 
      className={`p-4 rounded-lg border-2 cursor-pointer transition-all hover:shadow-md ${
        completed ? 'border-green-500 bg-green-50' : `${nodeStyle.borderColor} ${nodeStyle.bgColor}`
      }`}
      onClick={onClick}
    >
      <div className="flex justify-between items-start">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-lg">{nodeStyle.icon}</span>
            <h3 className="font-semibold text-lg">{title}</h3>
          </div>
          <span className="text-xs text-gray-500">{nodeStyle.title}</span>
        </div>
        <div className="flex flex-col items-end">
          <span className={`text-xs px-2 py-1 rounded-full ${levelColors[level] || 'bg-gray-100 text-gray-800'}`}>
            {level}
          </span>
          {masteryPercentage !== null && (
            <div className="mt-1 text-xs text-gray-600">
              掌握度: {masteryPercentage}%
            </div>
          )}
        </div>
      </div>
      <p className="text-gray-600 mt-2 text-sm">{description}</p>
      {completed && (
        <div className="mt-2 text-green-600 text-sm flex items-center">
          <span className="mr-1">✓</span> 已完成
        </div>
      )}
      
      {/* Mastery indicator bar */}
      {masteryPercentage !== null && (
        <div className="mt-2 w-full bg-gray-200 rounded-full h-2">
          <div 
            className={`h-2 rounded-full ${
              masteryPercentage >= 80 ? 'bg-green-500' : 
              masteryPercentage >= 60 ? 'bg-yellow-500' : 
              masteryPercentage >= 40 ? 'bg-orange-500' : 'bg-red-500'
            }`}
            style={{ width: `${masteryPercentage}%` }}
          ></div>
        </div>
      )}
    </div>
  );
};

export default KnowledgeNode;