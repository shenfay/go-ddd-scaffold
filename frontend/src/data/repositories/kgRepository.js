/**
 * 知识图谱数据仓储
 *
 * 统一管理知识图谱相关数据的获取
 * 职责：
 * 1. 封装 HTTP 请求逻辑
 * 2. 调用 Mapper 进行数据转换
 * 3. 提供简洁的数据访问接口
 */

import httpClient from '../api/client';
import { getEndpoint } from '../endpoints/endpoints';
import {
  mapDomain,
  mapTrunk,
  mapNode,
  mapRelationship,
  mapArray,
  mapToGraphNode,
  mapToGraphEdge,
  calculateNodePosition
} from '../mappers/knowledgeMapper';

/**
 * KgRepository - 知识图谱数据仓储
 */
class KgRepository {
  constructor() {
    // 配置基础 URL（后端运行在 8080 端口）
    httpClient.setBaseURL('http://localhost:8080/api');
  }

  /**
   * 获取所有领域
   * @returns {Promise<Array>} 领域列表
   */
  async getDomains() {
    try {
      const endpoint = getEndpoint('kg.getDomains');
      console.log('[KgRepository] 请求端点:', endpoint);
      console.log('[KgRepository] 基础URL:', httpClient.getBaseURL());

      const response = await httpClient.get(endpoint);
      console.log('[KgRepository] 响应数据:', response.data);

      return mapArray(response.data, mapDomain);
    } catch (error) {
      console.error('[KgRepository] 获取领域列表失败:', error);
      const detailedError = error.message || '网络连接失败，请检查您的网络设置';
      throw new Error(`获取领域数据失败: ${detailedError}`);
    }
  }

  /**
   * 根据 ID 获取领域
   * @param {string} id - 领域 ID
   * @returns {Promise<object|null>} 领域数据
   */
  async getDomainById(id) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getDomain', { id }));
      return mapDomain(response.data);
    } catch (error) {
      console.error('[KgRepository] 获取领域详情失败:', error);
      throw new Error(`获取领域详情失败: ${error.message}`);
    }
  }

  /**
   * 获取指定领域下的所有主线
   * @param {string} domainId - 领域 ID
   * @returns {Promise<Array>} 主线列表
   */
  async getTrunksByDomain(domainId) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getTrunksByDomain', { domainId }));
      return mapArray(response.data, mapTrunk);
    } catch (error) {
      console.error('[KgRepository] 获取主线列表失败:', error);
      throw new Error(`获取主线数据失败: ${error.message}`);
    }
  }

  /**
   * 根据 ID 获取主线
   * @param {string} id - 主线 ID
   * @returns {Promise<object|null>} 主线数据
   */
  async getTrunkById(id) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getTrunk', { id }));
      return mapTrunk(response.data);
    } catch (error) {
      console.error('[KgRepository] 获取主线详情失败:', error);
      throw new Error(`获取主线详情失败: ${error.message}`);
    }
  }

  /**
   * 获取指定主线下的所有节点
   * @param {string} trunkId - 主线 ID
   * @returns {Promise<Array>} 节点列表
   */
  async getNodesByTrunk(trunkId) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getNodesByTrunk', { trunkId }));
      return mapArray(response.data, mapNode);
    } catch (error) {
      console.error('[KgRepository] 获取节点列表失败:', error);
      throw new Error(`获取节点数据失败: ${error.message}`);
    }
  }

  /**
   * 根据 ID 获取节点
   * @param {string} id - 节点 ID
   * @returns {Promise<object|null>} 节点数据
   */
  async getNodeById(id) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getNode', { id }));
      return mapNode(response.data);
    } catch (error) {
      console.error('[KgRepository] 获取节点详情失败:', error);
      throw new Error(`获取节点详情失败: ${error.message}`);
    }
  }

  /**
   * 获取节点的所有关系
   * @param {string} nodeId - 节点 ID
   * @returns {Promise<Array>} 关系列表
   */
  async getNodeRelationships(nodeId) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getNodeRelationships', { id: nodeId }));
      return mapArray(response.data, mapRelationship);
    } catch (error) {
      console.error('[KgRepository] 获取节点关系失败:', error);
      // 关系获取失败不应阻止主流程
      return [];
    }
  }

  /**
   * 获取节点的前置知识
   * @param {string} nodeId - 节点 ID
   * @returns {Promise<Array>} 前置节点列表
   */
  async getNodePrerequisites(nodeId) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getNodePrerequisites', { id: nodeId }));
      return mapArray(response.data, mapNode);
    } catch (error) {
      console.error('[KgRepository] 获取前置知识失败:', error);
      return [];
    }
  }

  /**
   * 获取依赖该节点的后续节点
   * @param {string} nodeId - 节点 ID
   * @returns {Promise<Array>} 后续节点列表
   */
  async getNodeDependents(nodeId) {
    try {
      const response = await httpClient.get(getEndpoint('kg.getNodeDependents', { id: nodeId }));
      return mapArray(response.data, mapNode);
    } catch (error) {
      console.error('[KgRepository] 获取依赖节点失败:', error);
      return [];
    }
  }

  /**
   * 获取完整的图谱数据（用于展示）
   * @param {string} domainId - 可选，指定领域 ID，默认使用第一个领域
   * @returns {Promise<object>} 图谱数据 { nodes, edges }
   */
  async getGraphData(domainId = null) {
    try {
      // 1. 获取领域列表
      const domains = await this.getDomains();
      if (domains.length === 0) {
        return { nodes: [], edges: [] };
      }

      const targetDomainId = domainId || domains[0].id;
      const domain = domains.find(d => d.id === targetDomainId) || domains[0];

      // 2. 获取该领域下的所有主线
      const trunks = await this.getTrunksByDomain(targetDomainId);
      if (trunks.length === 0) {
        return { nodes: [], edges: [], domain };
      }

      // 3. 收集所有节点和关系
      const graphNodes = [];
      const graphEdges = [];
      let nodeIndex = 0;

      for (const trunk of trunks) {
        const nodes = await this.getNodesByTrunk(trunk.id);

        for (const node of nodes) {
          // 转换为图谱节点并计算位置
          const graphNode = mapToGraphNode(node, nodeIndex, 100); // 预估总数为 100
          if (graphNode) {
            graphNodes.push(graphNode);
            nodeIndex++;

            // 获取节点关系
            const relationships = await this.getNodeRelationships(node.id);
            for (const rel of relationships) {
              const edge = {
                source: rel.fromNodeId,
                target: rel.toNodeId,
                type: rel.relationshipType
              };
              graphEdges.push(edge);
            }
          }
        }
      }

      // 重新计算所有节点位置
      graphNodes.forEach((node, index) => {
        const position = calculateNodePosition(index, graphNodes.length, {
          centerX: 400,
          centerY: 300,
          radiusX: 280,
          radiusY: 220
        });
        node.x = position.x;
        node.y = position.y;
      });

      return {
        nodes: graphNodes,
        edges: graphEdges,
        domain,
        trunks
      };
    } catch (error) {
      console.error('[KgRepository] 获取图谱数据失败:', error);
      throw new Error(`获取图谱数据失败: ${error.message}`);
    }
  }

  /**
   * 简化版：获取图谱节点数据（不包含关系）
   * @param {string} domainId - 可选，指定领域 ID
   * @returns {Promise<Array>} 节点列表
   */
  async getGraphNodes(domainId = null) {
    try {
      const domains = await this.getDomains();
      if (domains.length === 0) return [];

      const targetDomainId = domainId || domains[0].id;
      const trunks = await this.getTrunksByDomain(targetDomainId);

      const allNodes = [];
      let index = 0;

      for (const trunk of trunks) {
        const nodes = await this.getNodesByTrunk(trunk.id);
        for (const node of nodes) {
          const graphNode = mapToGraphNode(node, index, 100);
          if (graphNode) {
            allNodes.push(graphNode);
            index++;
          }
        }
      }

      // 重新计算位置
      allNodes.forEach((node, i) => {
        const pos = calculateNodePosition(i, allNodes.length);
        node.x = pos.x;
        node.y = pos.y;
      });

      return allNodes;
    } catch (error) {
      console.error('[KgRepository] 获取图谱节点失败:', error);
      throw new Error(`获取图谱节点失败: ${error.message}`);
    }
  }

  /**
   * 获取指定主线的图谱节点
   * @param {string} trunkId - 主线 ID
   * @returns {Promise<Array>} 节点列表
   */
  async getTrunkGraphNodes(trunkId) {
    try {
      const nodes = await this.getNodesByTrunk(trunkId);

      return nodes.map((node, index) => {
        const graphNode = mapToGraphNode(node, index, nodes.length);
        const position = calculateNodePosition(index, nodes.length);
        graphNode.x = position.x;
        graphNode.y = position.y;
        return graphNode;
      });
    } catch (error) {
      console.error('[KgRepository] 获取主线节点失败:', error);
      throw new Error(`获取主线节点失败: ${error.message}`);
    }
  }
}

// 创建全局单例
const kgRepository = new KgRepository();

export default kgRepository;
export { KgRepository };
