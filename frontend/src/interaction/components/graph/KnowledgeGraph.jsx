import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';

const KnowledgeGraph = ({ nodes, edges, onNodeClick }) => {
  const svgRef = useRef();

  // Default sample data if no data provided
  const defaultNodes = [
    { id: '1', name: '整数', type: 'C', level: 1, x: 100, y: 200 },
    { id: '2', name: '分数', type: 'C', level: 2, x: 300, y: 150 },
    { id: '3', name: '小数', type: 'C', level: 3, x: 500, y: 200 },
    { id: '4', name: '加法', type: 'S', level: 1, x: 200, y: 300 },
    { id: '5', name: '减法', type: 'S', level: 1, x: 400, y: 300 },
    { id: '6', name: '乘法', type: 'S', level: 3, x: 600, y: 300 },
    { id: '7', name: '数形结合', type: 'T', level: 3, x: 300, y: 50 },
    { id: '8', name: '鸡兔同笼', type: 'P', level: 3, x: 500, y: 50 },
  ];

  const defaultEdges = [
    { source: '1', target: '2', type: 'PREREQ' },
    { source: '2', target: '3', type: 'PREREQ' },
    { source: '1', target: '4', type: 'SUP_SKILL' },
    { source: '1', target: '5', type: 'SUP_SKILL' },
    { source: '2', target: '6', type: 'SUP_SKILL' },
    { source: '4', target: '7', type: 'THINK_PAT' },
    { source: '5', target: '7', type: 'THINK_PAT' },
    { source: '7', target: '8', type: 'PREREQ' },
  ];

  const graphNodes = nodes || defaultNodes;
  const graphEdges = edges || defaultEdges;

  useEffect(() => {
    if (!svgRef.current) return;

    // Clear previous SVG content
    d3.select(svgRef.current).selectAll("*").remove();

    const svg = d3.select(svgRef.current);
    const width = svg.node().clientWidth;
    const height = svg.node().clientHeight;

    // Create simulation
    const simulation = d3.forceSimulation(graphNodes)
      .force('link', d3.forceLink(graphEdges).id(d => d.id).distance(100))
      .force('charge', d3.forceManyBody().strength(-300))
      .force('center', d3.forceCenter(width / 2, height / 2))
      .force('collision', d3.forceCollide().radius(40));

    // Draw links
    const link = svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(graphEdges)
      .enter()
      .append('line')
      .attr('stroke', '#94a3b8')
      .attr('stroke-width', 1.5)
      .attr('stroke-dasharray', d => {
        if (d.type === 'PREREQ') return '5,5';
        if (d.type === 'SUP_SKILL') return '10,5';
        return '1,0';
      });

    // Draw nodes
    const node = svg.append('g')
      .attr('class', 'nodes')
      .selectAll('g')
      .data(graphNodes)
      .enter()
      .append('g')
      .call(d3.drag()
        .on('start', (event, d) => {
          if (!event.active) simulation.alphaTarget(0.3).restart();
          d.fx = d.x;
          d.fy = d.y;
        })
        .on('drag', (event, d) => {
          d.fx = event.x;
          d.fy = event.y;
        })
        .on('end', (event, d) => {
          if (!event.active) simulation.alphaTarget(0);
          d.fx = null;
          d.fy = null;
        }));

    // Add circles to nodes
    node.append('circle')
      .attr('r', d => {
        // Different sizes based on level
        return 15 + (d.level - 1) * 5;
      })
      .attr('fill', d => {
        // Different colors based on type
        if (d.type === 'C') { // Concept
          const levelColors = ['#90EE90', '#32CD32', '#228B22', '#006400', '#004d00'];
          return levelColors[d.level - 1] || levelColors[0];
        } else if (d.type === 'S') { // Support Skill
          const levelColors = ['#87CEEB', '#4682B4', '#4169E1', '#00008B', '#00004d'];
          return levelColors[d.level - 1] || levelColors[0];
        } else if (d.type === 'T') { // Thinking Pattern
          const levelColors = ['#FFD700', '#FFA500', '#FF8C00', '#FF4500', '#8B0000'];
          return levelColors[d.level - 1] || levelColors[0];
        } else if (d.type === 'P') { // Problem Model
          const levelColors = ['#DDA0DD', '#9370DB', '#8A2BE2', '#4B0082', '#320064'];
          return levelColors[d.level - 1] || levelColors[0];
        }
        return '#ccc';
      })
      .attr('stroke', '#fff')
      .attr('stroke-width', 2)
      .on('click', (event, d) => {
        if (onNodeClick) {
          onNodeClick(d);
        }
      })
      .on('mouseover', function(event, d) {
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', d => 20 + (d.level - 1) * 5)
          .attr('stroke-width', 3);
          
        // Show tooltip
        const tooltip = d3.select('.tooltip');
        if (!tooltip.empty()) {
          tooltip
            .html(`
              <div class="bg-gray-800 text-white p-2 rounded text-sm">
                <div class="font-bold">${d.name}</div>
                <div>Type: ${d.type === 'C' ? '概念' : d.type === 'S' ? '支撑技能' : d.type === 'T' ? '思维模式' : '问题模型'}</div>
                <div>Level: ${d.level}</div>
              </div>
            `)
            .style('visibility', 'visible')
            .style('left', (event.pageX + 10) + 'px')
            .style('top', (event.pageY - 10) + 'px');
        }
      })
      .on('mouseout', function(event, d) {
        d3.select(this)
          .transition()
          .duration(200)
          .attr('r', d => 15 + (d.level - 1) * 5)
          .attr('stroke-width', 2);
              
        // Hide tooltip
        const tooltip = d3.select('.tooltip');
        if (!tooltip.empty()) {
          tooltip.style('visibility', 'hidden');
        }
      });

    // Add labels to nodes
    node.append('text')
      .attr('text-anchor', 'middle')
      .attr('dy', 4)
      .attr('fill', '#fff')
      .attr('font-size', '12px')
      .attr('font-weight', 'bold')
      .text(d => d.name.length > 6 ? d.name.substring(0, 6) + '..' : d.name);

    // Update positions on each tick
    simulation.on('tick', () => {
      link
        .attr('x1', d => d.source.x)
        .attr('y1', d => d.source.y)
        .attr('x2', d => d.target.x)
        .attr('y2', d => d.target.y);

      node.attr('transform', d => `translate(${d.x},${d.y})`);
    });

    // Cleanup function
    return () => {
      simulation.stop();
    };
  }, [graphNodes, graphEdges, onNodeClick]);

  return (
    <div className="relative w-full h-full">
      <svg 
        ref={svgRef} 
        className="w-full h-full bg-gray-50 rounded-lg border border-gray-200"
        style={{ minHeight: '500px' }}
      />
      <div className="tooltip absolute pointer-events-none z-10" style={{ visibility: 'hidden' }} />
      
      {/* Legend */}
      <div className="absolute top-4 right-4 bg-white p-4 rounded-lg shadow-lg border border-gray-200">
        <h3 className="font-bold mb-2">图例</h3>
        <div className="space-y-2">
          <div className="flex items-center">
            <div className="w-4 h-4 rounded-full bg-green-400 mr-2"></div>
            <span className="text-xs">概念 (C)</span>
          </div>
          <div className="flex items-center">
            <div className="w-4 h-4 rounded bg-blue-400 mr-2"></div>
            <span className="text-xs">支撑技能 (S)</span>
          </div>
          <div className="flex items-center">
            <div className="w-4 h-4 rounded-full bg-yellow-400 mr-2"></div>
            <span className="text-xs">思维模式 (T)</span>
          </div>
          <div className="flex items-center">
            <div className="w-4 h-4 rounded bg-purple-400 mr-2"></div>
            <span className="text-xs">问题模型 (P)</span>
          </div>
          <div className="pt-2 border-t border-gray-200">
            <div className="text-xs text-gray-600">虚线: 先修关系</div>
            <div className="text-xs text-gray-600">点划线: 支持关系</div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default KnowledgeGraph;