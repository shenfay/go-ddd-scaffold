import React, { useState, useEffect } from 'react';
import KnowledgeGraph from '../../../../../interaction/components/graph/KnowledgeGraph';
import MixedKnowledgeView from '../../../../../interaction/components/MixedKnowledgeView';
import kgRepository from '../../../../../data/repositories/knowledgeRepository';

const KnowledgeGraphDemoPage = () => {
  const [activeTab, setActiveTab] = useState('mixed');
  const [selectedNode, setSelectedNode] = useState(null);
  const [demoData, setDemoData] = useState({
    nodes: [],
    edges: []
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  // 从后端 API 获取数据
  useEffect(() => {
    const fetchData = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // 使用 Repository 统一获取数据
        const graphData = await kgRepository.getGraphData();

        setDemoData({
          nodes: graphData.nodes,
          edges: graphData.edges
        });
        setIsLoading(false);
      } catch (err) {
        console.error('[KnowledgeGraphDemoPage] 获取数据失败:', err);
        setError(err.message);
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleNodeClick = (node) => {
    setSelectedNode(node);
  };

  const getNodeTypeLabel = (type) => {
    switch(type) {
      case 'C': return '概念';
      case 'S': return '支撑技能';
      case 'T': return '思维模式';
      case 'P': return '问题模型';
      default: return '未知类型';
    }
  };

  // 加载状态
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-lg text-gray-700">正在加载知识图谱数据...</p>
        </div>
      </div>
    );
  }

  // 错误状态
  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center">
        <div className="text-center bg-white rounded-xl shadow-lg p-8">
          <div className="text-6xl mb-4">❌</div>
          <h2 className="text-2xl font-bold text-gray-800 mb-2">数据加载失败</h2>
          <p className="text-gray-600 mb-4">{error}</p>
          <button 
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            onClick={() => window.location.reload()}
          >
            重新加载
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-6 overflow-y-auto">
      <div className="max-w-7xl mx-auto">
        {/* 页面标题 */}
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-800 mb-2">Go DDD Scaffold 知识图谱系统</h1>
          <p className="text-lg text-gray-600">探索数学知识的结构化呈现与智能学习路径</p>
        </div>

        {/* 控制面板 */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-6">
          <div className="flex flex-wrap gap-4 items-center justify-between">
            <div className="flex space-x-2">
              <button
                className={`px-6 py-3 rounded-lg font-medium transition-all ${
                  activeTab === 'mixed' 
                    ? 'bg-blue-600 text-white shadow-md' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('mixed')}
              >
                混合视图
              </button>
              <button
                className={`px-6 py-3 rounded-lg font-medium transition-all ${
                  activeTab === 'graph' 
                    ? 'bg-blue-600 text-white shadow-md' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('graph')}
              >
                关系图谱
              </button>
              <button
                className={`px-6 py-3 rounded-lg font-medium transition-all ${
                  activeTab === 'stats' 
                    ? 'bg-blue-600 text-white shadow-md' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('stats')}
              >
                数据统计
              </button>
            </div>
            
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-600">
                节点总数: {demoData.nodes.length}
              </div>
              <div className="text-sm text-gray-600">
                关系数: {demoData.edges.length}
              </div>
            </div>
          </div>
        </div>

        {/* 主要内容区域 */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* 侧边栏 - 节点信息 */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-xl shadow-lg p-6 h-full">
              <h3 className="text-xl font-bold text-gray-800 mb-4">节点详情</h3>
              
              {selectedNode ? (
                <div className="space-y-4">
                  <div className="p-4 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg">
                    <h4 className="font-bold text-lg text-gray-800 mb-2">{selectedNode.name}</h4>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-600">类型:</span>
                        <span className="font-medium">{getNodeTypeLabel(selectedNode.type)}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">等级:</span>
                        <span className="font-medium">Lv{selectedNode.level}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-600">ID:</span>
                        <span className="font-mono text-xs">{selectedNode.id}</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="p-4 bg-gray-50 rounded-lg">
                    <h5 className="font-semibold text-gray-700 mb-2">掌握度</h5>
                    <div className="w-full bg-gray-200 rounded-full h-3">
                      <div 
                        className="bg-gradient-to-r from-green-400 to-green-600 h-3 rounded-full transition-all duration-300"
                        style={{ width: `${Math.random() * 100}%` }}
                      ></div>
                    </div>
                    <div className="text-right text-xs text-gray-500 mt-1">
                      {Math.floor(Math.random() * 100)}%
                    </div>
                  </div>
                  
                  <div className="p-4 bg-gray-50 rounded-lg">
                    <h5 className="font-semibold text-gray-700 mb-2">前置知识</h5>
                    <ul className="text-sm text-gray-600 space-y-1">
                      {demoData.edges
                        .filter(edge => edge.target === selectedNode.id && edge.type === 'PREREQ')
                        .map((edge, index) => {
                          const sourceNode = demoData.nodes.find(n => n.id === edge.source);
                          return sourceNode ? (
                            <li key={index} className="flex items-center">
                              <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
                              {sourceNode.name}
                            </li>
                          ) : null;
                        })}
                    </ul>
                  </div>
                </div>
              ) : (
                <div className="text-center py-12 text-gray-500">
                  <div className="text-6xl mb-4">🔍</div>
                  <p>点击图谱中的节点查看详情</p>
                </div>
              )}
            </div>
          </div>

          {/* 主内容区域 */}
          <div className="lg:col-span-3">
            <div className="bg-white rounded-xl shadow-lg p-6 h-full min-h-[600px]">
              {activeTab === 'mixed' && (
                <div className="h-full">
                  <MixedKnowledgeView initialView="graph" />
                </div>
              )}
              
              {activeTab === 'graph' && (
                <div className="h-full">
                  <div className="mb-4">
                    <h3 className="text-xl font-bold text-gray-800">知识关系图谱</h3>
                    <p className="text-gray-600">可视化展示数学知识点之间的关联关系</p>
                  </div>
                  <div className="h-[500px]">
                    <KnowledgeGraph 
                      nodes={demoData.nodes} 
                      edges={demoData.edges} 
                      onNodeClick={handleNodeClick} 
                    />
                  </div>
                </div>
              )}
              
              {activeTab === 'stats' && (
                <div className="h-full">
                  <h3 className="text-xl font-bold text-gray-800 mb-6">数据统计</h3>
                  
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                    <div className="bg-gradient-to-r from-green-400 to-green-600 rounded-lg p-6 text-white">
                      <div className="text-3xl font-bold">{demoData.nodes.filter(n => n.type === 'C').length}</div>
                      <div className="text-sm opacity-90">概念节点</div>
                    </div>
                    
                    <div className="bg-gradient-to-r from-blue-400 to-blue-600 rounded-lg p-6 text-white">
                      <div className="text-3xl font-bold">{demoData.nodes.filter(n => n.type === 'S').length}</div>
                      <div className="text-sm opacity-90">技能节点</div>
                    </div>
                    
                    <div className="bg-gradient-to-r from-purple-400 to-purple-600 rounded-lg p-6 text-white">
                      <div className="text-3xl font-bold">{demoData.nodes.filter(n => n.type === 'T' || n.type === 'P').length}</div>
                      <div className="text-sm opacity-90">高级节点</div>
                    </div>
                  </div>
                  
                  <div className="space-y-4">
                    <h4 className="font-semibold text-gray-700">节点分布</h4>
                    {[1, 2, 3, 4, 5].map(level => (
                      <div key={level} className="flex items-center space-x-4">
                        <div className="w-16 text-sm font-medium text-gray-600">Lv{level}</div>
                        <div className="flex-1">
                          <div className="w-full bg-gray-200 rounded-full h-3">
                            <div 
                              className={`h-3 rounded-full transition-all duration-300 ${
                                level === 1 ? 'bg-green-500' :
                                level === 2 ? 'bg-blue-500' :
                                level === 3 ? 'bg-yellow-500' :
                                level === 4 ? 'bg-orange-500' : 'bg-red-500'
                              }`}
                              style={{ 
                                width: `${(demoData.nodes.filter(n => n.level === level).length / demoData.nodes.length) * 100}%` 
                              }}
                            ></div>
                          </div>
                        </div>
                        <div className="w-12 text-right text-sm text-gray-600">
                          {demoData.nodes.filter(n => n.level === level).length}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* 功能介绍 */}
        <div className="mt-8 bg-white rounded-xl shadow-lg p-6">
          <h3 className="text-xl font-bold text-gray-800 mb-4">系统特性</h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="text-center p-4">
              <div className="text-4xl mb-3">🔄</div>
              <h4 className="font-semibold text-gray-800 mb-2">动态关系图谱</h4>
              <p className="text-gray-600 text-sm">可视化展示知识点之间的先修、支撑和思维关联关系</p>
            </div>
            <div className="text-center p-4">
              <div className="text-4xl mb-3">🎯</div>
              <h4 className="font-semibold text-gray-800 mb-2">C/S/T/P分类</h4>
              <p className="text-gray-600 text-sm">概念、技能、思维、问题模型四种节点类型精准定位学习内容</p>
            </div>
            <div className="text-center p-4">
              <div className="text-4xl mb-3">📊</div>
              <h4 className="font-semibold text-gray-800 mb-2">智能分析</h4>
              <p className="text-gray-600 text-sm">基于掌握度和关系图谱的个性化学习路径推荐</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeGraphDemoPage;