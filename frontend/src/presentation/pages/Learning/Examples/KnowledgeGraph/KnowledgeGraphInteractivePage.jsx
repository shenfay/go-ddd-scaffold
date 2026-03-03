import React, { useState, useEffect } from 'react';
import KnowledgeGraph from '../../../../../interaction/components/graph/KnowledgeGraph';
import KnowledgeGraph3D from '../../../../../interaction/components/KnowledgeGraph3D';
import MixedKnowledgeView from '../../../../../interaction/components/MixedKnowledgeView';
import kgRepository from '../../../../../data/repositories/knowledgeRepository';

const KnowledgeGraphInteractivePage = () => {
  const [activeTab, setActiveTab] = useState('demo');
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
        console.error('[KnowledgeGraphInteractivePage] 获取数据失败:', err);
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

  const getNodeTypeColor = (type) => {
    switch(type) {
      case 'C': return 'bg-green-100 text-green-800';
      case 'S': return 'bg-blue-100 text-blue-800';
      case 'T': return 'bg-yellow-100 text-yellow-800';
      case 'P': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

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
          <h1 className="text-5xl font-bold text-gray-800 mb-4 flex items-center justify-center">
            <span className="mr-4">🧠</span>
            MathFun 知识图谱系统
            <span className="ml-4">🎯</span>
          </h1>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            探索数学知识的结构化呈现与智能学习路径，通过C/S/T/P四类节点构建完整的学习生态系统
          </p>
        </div>

        {/* 控制面板 */}
        <div className="bg-white rounded-2xl shadow-xl p-6 mb-8">
          <div className="flex flex-wrap gap-4 items-center justify-between">
            <div className="flex space-x-2 flex-wrap">
              <button
                className={`px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 ${
                  activeTab === 'demo' 
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('demo')}
              >
                交互演示
              </button>
              <button
                className={`px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 ${
                  activeTab === '3d' 
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('3d')}
              >
                3D 视图
              </button>
              <button
                className={`px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 ${
                  activeTab === 'graph' 
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('graph')}
              >
                关系图谱
              </button>
              <button
                className={`px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 ${
                  activeTab === 'mixed' 
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('mixed')}
              >
                混合视图
              </button>
              <button
                className={`px-6 py-3 rounded-xl font-medium transition-all transform hover:scale-105 ${
                  activeTab === 'insights' 
                    ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
                onClick={() => setActiveTab('insights')}
              >
                数据洞察
              </button>
            </div>
            
            <div className="flex items-center space-x-6">
              <div className="text-sm text-gray-600 bg-blue-50 px-4 py-2 rounded-full">
                📊 节点总数: {demoData.nodes.length}
              </div>
              <div className="text-sm text-gray-600 bg-purple-50 px-4 py-2 rounded-full">
                🔗 关系数: {demoData.edges.length}
              </div>
            </div>
          </div>
        </div>

        {/* 主要内容区域 */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* 侧边栏 - 节点信息 */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-2xl shadow-xl p-6 h-full sticky top-6">
              <h3 className="text-2xl font-bold text-gray-800 mb-6 flex items-center">
                <span className="mr-3">🔍</span>
                节点详情
              </h3>
              
              {selectedNode ? (
                <div className="space-y-6">
                  <div className="p-5 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl border border-blue-100">
                    <div className="flex items-center justify-between mb-4">
                      <h4 className="font-bold text-xl text-gray-800">{selectedNode.name}</h4>
                      <span className={`px-3 py-1 rounded-full text-xs font-semibold ${getNodeTypeColor(selectedNode.type)}`}>
                        {getNodeTypeLabel(selectedNode.type)}
                      </span>
                    </div>
                    <div className="space-y-3 text-sm">
                      <div className="flex justify-between items-center">
                        <span className="text-gray-600">能力等级:</span>
                        <span className="font-bold text-blue-600">Lv{selectedNode.level}</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-gray-600">节点ID:</span>
                        <span className="font-mono text-xs bg-gray-100 px-2 py-1 rounded">{selectedNode.id}</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="p-5 bg-gradient-to-r from-green-50 to-emerald-50 rounded-xl border border-green-100">
                    <h5 className="font-semibold text-gray-700 mb-3 flex items-center">
                      <span className="mr-2">📈</span>
                      掌握度
                    </h5>
                    <div className="w-full bg-gray-200 rounded-full h-4 mb-2">
                      <div 
                        className="bg-gradient-to-r from-green-500 to-emerald-600 h-4 rounded-full transition-all duration-500 ease-out"
                        style={{ width: `${Math.floor(Math.random() * 100)}%` }}
                      ></div>
                    </div>
                    <div className="text-right text-sm text-gray-500">
                      {Math.floor(Math.random() * 100)}%
                    </div>
                  </div>
                  
                  <div className="p-5 bg-gradient-to-r from-purple-50 to-fuchsia-50 rounded-xl border border-purple-100">
                    <h5 className="font-semibold text-gray-700 mb-3 flex items-center">
                      <span className="mr-2">🔗</span>
                      相关节点
                    </h5>
                    <div className="space-y-2">
                      <div className="text-sm">
                        <span className="text-gray-600 block mb-1">前置知识:</span>
                        <ul className="space-y-1">
                          {demoData.edges
                            .filter(edge => edge.target === selectedNode.id && edge.type === 'PREREQ')
                            .slice(0, 3)
                            .map((edge, index) => {
                              const sourceNode = demoData.nodes.find(n => n.id === edge.source);
                              return sourceNode ? (
                                <li key={index} className="flex items-center bg-white p-2 rounded-lg">
                                  <span className="w-2 h-2 bg-blue-500 rounded-full mr-2"></span>
                                  <span className="text-sm">{sourceNode.name}</span>
                                </li>
                              ) : null;
                            })}
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="text-center py-16 text-gray-500">
                  <div className="text-8xl mb-6">🎯</div>
                  <p className="text-lg font-medium mb-2">请选择一个节点</p>
                  <p className="text-sm">点击图谱中的任意节点查看详情</p>
                </div>
              )}
            </div>
          </div>

          {/* 主内容区域 */}
          <div className="lg:col-span-3">
            <div className="bg-white rounded-2xl shadow-xl p-6 h-full min-h-[700px]">
              {activeTab === 'demo' && (
                <div className="h-full">
                  <div className="mb-6">
                    <h3 className="text-2xl font-bold text-gray-800 mb-2">交互式知识图谱演示</h3>
                    <p className="text-gray-600">点击节点查看详情，拖拽节点调整位置，滚动缩放</p>
                  </div>
                  <div className="h-[600px] bg-gradient-to-br from-blue-50 to-indigo-100 rounded-xl p-4">
                    <KnowledgeGraph 
                      nodes={demoData.nodes} 
                      edges={demoData.edges} 
                      onNodeClick={handleNodeClick} 
                    />
                  </div>
                </div>
              )}
              
              {activeTab === '3d' && (
                <div className="h-full">
                  <div className="mb-6">
                    <h3 className="text-2xl font-bold text-gray-800 mb-2">3D 知识图谱可视化</h3>
                    <p className="text-gray-600">立体展示知识节点的空间关系和层次结构</p>
                  </div>
                  <div className="h-[600px] bg-gradient-to-br from-slate-900 to-blue-900 rounded-xl overflow-hidden">
                    <KnowledgeGraph3D 
                      initialNodes={demoData.nodes}
                      initialEdges={demoData.edges}
                      onNodeClick={handleNodeClick}
                    />
                  </div>
                </div>
              )}
              
              {activeTab === 'graph' && (
                <div className="h-full">
                  <div className="mb-6">
                    <h3 className="text-2xl font-bold text-gray-800 mb-2">知识关系图谱</h3>
                    <p className="text-gray-600">可视化展示数学知识点之间的关联关系</p>
                  </div>
                  <div className="h-[600px] bg-gradient-to-br from-blue-50 to-indigo-100 rounded-xl p-4">
                    <KnowledgeGraph 
                      nodes={demoData.nodes} 
                      edges={demoData.edges} 
                      onNodeClick={handleNodeClick} 
                    />
                  </div>
                </div>
              )}
              
              {activeTab === 'mixed' && (
                <div className="h-full">
                  <div className="mb-6">
                    <h3 className="text-2xl font-bold text-gray-800 mb-2">混合知识视图</h3>
                    <p className="text-gray-600">综合展示知识图谱的多种表现形式</p>
                  </div>
                  <div className="h-[600px]">
                    <MixedKnowledgeView initialView="graph" />
                  </div>
                </div>
              )}
              
              {activeTab === 'insights' && (
                <div className="h-full">
                  <h3 className="text-2xl font-bold text-gray-800 mb-6">数据洞察分析</h3>
                  
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                    <div className="bg-gradient-to-br from-green-400 to-green-600 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-transform">
                      <div className="text-4xl font-bold mb-2">{demoData.nodes.filter(n => n.type === 'C').length}</div>
                      <div className="text-sm opacity-90">概念节点</div>
                      <div className="mt-2 text-xs opacity-75">基础定义与性质</div>
                    </div>
                    
                    <div className="bg-gradient-to-br from-blue-400 to-blue-600 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-transform">
                      <div className="text-4xl font-bold mb-2">{demoData.nodes.filter(n => n.type === 'S').length}</div>
                      <div className="text-sm opacity-90">技能节点</div>
                      <div className="mt-2 text-xs opacity-75">程序性知识</div>
                    </div>
                    
                    <div className="bg-gradient-to-br from-purple-400 to-purple-600 rounded-2xl p-6 text-white shadow-lg transform hover:scale-105 transition-transform">
                      <div className="text-4xl font-bold mb-2">{demoData.nodes.filter(n => n.type === 'T' || n.type === 'P').length}</div>
                      <div className="text-sm opacity-90">高级节点</div>
                      <div className="mt-2 text-xs opacity-75">思维与问题模型</div>
                    </div>
                  </div>
                  
                  <div className="space-y-6">
                    <h4 className="font-semibold text-gray-700 text-lg">节点等级分布</h4>
                    {[1, 2, 3, 4, 5].map(level => {
                      const count = demoData.nodes.filter(n => n.level === level).length;
                      const percentage = (count / demoData.nodes.length) * 100;
                      
                      return (
                        <div key={level} className="flex items-center space-x-4">
                          <div className="w-16 text-sm font-medium text-gray-600">Lv{level}</div>
                          <div className="flex-1">
                            <div className="w-full bg-gray-200 rounded-full h-4">
                              <div 
                                className={`h-4 rounded-full transition-all duration-500 ease-out ${
                                  level === 1 ? 'bg-green-500' :
                                  level === 2 ? 'bg-blue-500' :
                                  level === 3 ? 'bg-yellow-500' :
                                  level === 4 ? 'bg-orange-500' : 'bg-red-500'
                                }`}
                                style={{ width: `${percentage}%` }}
                              ></div>
                            </div>
                          </div>
                          <div className="w-16 text-right text-sm text-gray-600 font-medium">
                            {count} ({Math.round(percentage)}%)
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* 功能介绍 */}
        <div className="mt-12 bg-white rounded-2xl shadow-xl p-8">
          <h3 className="text-2xl font-bold text-gray-800 mb-8 text-center">系统核心特性</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="text-center p-6 hover:bg-blue-50 rounded-xl transition-colors">
              <div className="text-5xl mb-4">🔄</div>
              <h4 className="font-bold text-lg text-gray-800 mb-3">动态关系图谱</h4>
              <p className="text-gray-600 leading-relaxed">可视化展示知识点之间的先修、支撑和思维关联关系，帮助理解知识结构</p>
            </div>
            <div className="text-center p-6 hover:bg-green-50 rounded-xl transition-colors">
              <div className="text-5xl mb-4">🎯</div>
              <h4 className="font-bold text-lg text-gray-800 mb-3">C/S/T/P分类</h4>
              <p className="text-gray-600 leading-relaxed">概念、技能、思维、问题模型四种节点类型精准定位学习内容</p>
            </div>
            <div className="text-center p-6 hover:bg-purple-50 rounded-xl transition-colors">
              <div className="text-5xl mb-4">📊</div>
              <h4 className="font-bold text-lg text-gray-800 mb-3">智能分析</h4>
              <p className="text-gray-600 leading-relaxed">基于掌握度和关系图谱的个性化学习路径推荐</p>
            </div>
            <div className="text-center p-6 hover:bg-orange-50 rounded-xl transition-colors">
              <div className="text-5xl mb-4">🎮</div>
              <h4 className="font-bold text-lg text-gray-800 mb-3">交互体验</h4>
              <p className="text-gray-600 leading-relaxed">3D可视化和2D图表相结合，提供丰富的交互体验</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeGraphInteractivePage;