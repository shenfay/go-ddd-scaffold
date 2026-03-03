import React, { useState } from 'react';
import KnowledgeGraph from './graph/KnowledgeGraph';
import KnowledgeGraph3D from './KnowledgeGraph3D';
import Scene_KnowledgeVillage from './KnowledgeVillage';
import KnowledgeTownScene from '../three/components/scenes/KnowledgeTown';

const MixedKnowledgeView = ({ initialView = '3d', domainId = 'math' }) => {
  const [currentView, setCurrentView] = useState(initialView); // '3d' or 'graph'
  const [selectedLevel, setSelectedLevel] = useState('village'); // village, town, city

  // Sample data for the knowledge graph
  const knowledgeData = {
    nodes: [
      { id: '1', name: '整数', type: 'C', level: 1, x: 100, y: 200 },
      { id: '2', name: '分数', type: 'C', level: 2, x: 300, y: 150 },
      { id: '3', name: '小数', type: 'C', level: 3, x: 500, y: 200 },
      { id: '4', name: '加法', type: 'S', level: 1, x: 200, y: 300 },
      { id: '5', name: '减法', type: 'S', level: 1, x: 400, y: 300 },
      { id: '6', name: '乘法', type: 'S', level: 3, x: 600, y: 300 },
      { id: '7', name: '数形结合', type: 'T', level: 3, x: 300, y: 50 },
      { id: '8', name: '鸡兔同笼', type: 'P', level: 3, x: 500, y: 50 },
    ],
    edges: [
      { source: '1', target: '2', type: 'PREREQ' },
      { source: '2', target: '3', type: 'PREREQ' },
      { source: '1', target: '4', type: 'SUP_SKILL' },
      { source: '1', target: '5', type: 'SUP_SKILL' },
      { source: '2', target: '6', type: 'SUP_SKILL' },
      { source: '4', target: '7', type: 'THINK_PAT' },
      { source: '5', target: '7', type: 'THINK_PAT' },
      { source: '7', target: '8', type: 'PREREQ' },
    ]
  };

  const handleNodeClick = (node) => {
    console.log('Node clicked:', node);
    // Update level based on node level
    if (node.level <= 2) {
      setSelectedLevel('village');
    } else if (node.level <= 4) {
      setSelectedLevel('town');
    } else {
      setSelectedLevel('city');
    }
  };

  const render3DScene = () => {
    if (selectedLevel === 'village') {
      return <Scene_KnowledgeVillage />;
    } else if (selectedLevel === 'town') {
      return <KnowledgeTownScene />;
    } else {
      // For city level, we use the 3D knowledge graph
      return (
        <div className="w-full h-full">
          <KnowledgeGraph3D
            initialNodes={knowledgeData.nodes}
            initialEdges={knowledgeData.edges}
            onNodeClick={handleNodeClick}
          />
        </div>
      );
    }
  };

  return (
    <div className="w-full h-screen flex flex-col">
      {/* Control Panel */}
      <div className="bg-gray-800 text-white p-4 flex justify-between items-center">
        <div className="flex space-x-4">
          <button
            className={`px-4 py-2 rounded-lg ${currentView === '3d' ? 'bg-blue-600' : 'bg-gray-700 hover:bg-gray-600'}`}
            onClick={() => setCurrentView('3d')}
          >
            3D 视图
          </button>
          <button
            className={`px-4 py-2 rounded-lg ${currentView === 'graph' ? 'bg-blue-600' : 'bg-gray-700 hover:bg-gray-600'}`}
            onClick={() => setCurrentView('graph')}
          >
            图表视图
          </button>
        </div>
        
        <div className="flex items-center space-x-2">
          <span>当前等级:</span>
          <select 
            value={selectedLevel} 
            onChange={(e) => setSelectedLevel(e.target.value)}
            className="bg-gray-700 text-white px-2 py-1 rounded"
          >
            <option value="village">村庄 (Lv1-Lv2)</option>
            <option value="town">小镇 (Lv3-Lv4)</option>
            <option value="city">城市 (Lv5)</option>
          </select>
        </div>
      </div>

      {/* Content Area */}
      <div className="flex-grow">
        {currentView === '3d' ? (
          <div className="w-full h-full">
            {render3DScene()}
          </div>
        ) : (
          <div className="w-full h-full p-4 bg-white">
            <div className="mb-4">
              <h2 className="text-2xl font-bold text-gray-800">知识图谱关系图</h2>
              <p className="text-gray-600">可视化展示数学知识节点之间的关系</p>
            </div>
            <div className="w-full h-[calc(100%-80px)]">
              <KnowledgeGraph 
                nodes={knowledgeData.nodes} 
                edges={knowledgeData.edges} 
                onNodeClick={handleNodeClick} 
              />
            </div>
          </div>
        )}
      </div>

      {/* Status Bar */}
      <div className="bg-gray-100 border-t p-2 text-sm text-gray-600 flex justify-between">
        <div>当前视图: {currentView === '3d' ? '3D 场景' : '知识图谱图表'}</div>
        <div>当前场景: {selectedLevel === 'village' ? '知识村庄' : selectedLevel === 'town' ? '知识小镇' : '知识城市'}</div>
      </div>
    </div>
  );
};

export default MixedKnowledgeView;