import React from 'react';
import { useNavigate } from 'react-router-dom';

const ExamplesPage = () => {
  const navigate = useNavigate();

  const examples = [
    { name: '知识村庄场景', path: '/examples/knowledge-village' },
    { name: '知识小镇场景', path: '/examples/knowledge-town' },
    { name: '知识图谱演示', path: '/examples/knowledge-graph' },
    { name: '交互式知识图谱', path: '/examples/knowledge-graph-interactive' },
    { name: '世界演进动画', path: '/examples/evolution-animation' },
    { name: 'NPC 互动', path: '/examples/npc-interaction' },
    { name: '宠物互动', path: '/examples/pet-interaction' },
    { name: '数学城市街道', path: '/examples/math-city' },
  ];

  const handleSelect = (path) => {
    navigate(path);
  };

  return (
    <div className="p-8 bg-gray-100 min-h-screen">
      <h1 className="text-2xl font-bold mb-6 text-center">选择一个示例运行:</h1>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 max-w-4xl mx-auto mb-8">
        {examples.map((example, index) => (
          <button 
            key={index}
            className="example-selector-button bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded m-1"
            onClick={() => handleSelect(example.path)}
          >
            {example.name}
          </button>
        ))}
      </div>
    </div>
  );
};

export default ExamplesPage;