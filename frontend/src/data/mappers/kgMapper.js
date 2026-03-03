/**
 * 知识图谱数据映射器
 *
 * 统一处理后端 API 数据与前端模型之间的转换
 * 职责：
 * 1. 后端字段名 -> 前端字段名映射
 * 2. 数据类型转换
 * 3. 默认值处理
 */

// 后端 Node 类型映射
const NODE_TYPE_MAP = {
  'C': 'C', // Concept - 概念
  'S': 'S', // Support Skill - 支撑技能
  'T': 'T', // Thinking Pattern - 思维模式
  'P': 'P'  // Problem Model - 问题模型
};

// 关系类型映射
const RELATIONSHIP_TYPE_MAP = {
  'PREREQ': 'PREREQ',      // 先修关系
  'SUP_SKILL': 'SUP_SKILL', // 支撑技能关系
  'THINK_PAT': 'THINK_PAT'  // 思维模式关系
};

/**
 * 转换 Domain 数据
 * @param {object} apiData - 后端返回的 Domain 数据
 * @returns {object} 转换后的 Domain 数据
 */
export function mapDomain(apiData) {
  if (!apiData) return null;

  return {
    id: apiData.id || '',
    nameKey: apiData.nameKey || apiData.name_key || '',
    descriptionKey: apiData.descriptionKey || apiData.description_key || '',
    academicSource: apiData.academicSource || apiData.academic_source || '',
    isActive: apiData.isActive ?? apiData.is_active ?? true,
    createdAt: apiData.createdAt || apiData.created_at || new Date().toISOString(),
    updatedAt: apiData.updatedAt || apiData.updated_at || new Date().toISOString()
  };
}

/**
 * 转换 Trunk 数据
 * @param {object} apiData - 后端返回的 Trunk 数据
 * @returns {object} 转换后的 Trunk 数据
 */
export function mapTrunk(apiData) {
  if (!apiData) return null;

  return {
    id: apiData.id || '',
    domainId: apiData.domainId || apiData.domain_id || '',
    nameKey: apiData.nameKey || apiData.name_key || '',
    descriptionKey: apiData.descriptionKey || apiData.description_key || '',
    academicSource: apiData.academicSource || apiData.academic_source || '',
    isActive: apiData.isActive ?? apiData.is_active ?? true,
    createdAt: apiData.createdAt || apiData.created_at || new Date().toISOString(),
    updatedAt: apiData.updatedAt || apiData.updated_at || new Date().toISOString()
  };
}

/**
 * 转换 Node 数据
 * @param {object} apiData - 后端返回的 Node 数据
 * @returns {object} 转换后的 Node 数据
 */
export function mapNode(apiData) {
  if (!apiData) return null;

  // 从 competencyLevelId 提取等级数字
  let level = 1;
  if (apiData.competencyLevelId || apiData.competency_level_id) {
    const levelId = apiData.competencyLevelId || apiData.competency_level_id;
    // 格式如 "level-1", "level-2" 等
    const match = levelId.match(/level[-_]?(\d+)/i);
    if (match) {
      level = parseInt(match[1], 10) || 1;
    }
  }

  return {
    id: apiData.id || '',
    trunkId: apiData.trunkId || apiData.trunk_id || '',
    type: NODE_TYPE_MAP[apiData.type] || apiData.type || 'C',
    competencyLevelId: apiData.competencyLevelId || apiData.competency_level_id || '',
    nameChildKey: apiData.nameChildKey || apiData.name_child_key || '',
    nameParentKey: apiData.nameParentKey || apiData.name_parent_key || '',
    descriptionKey: apiData.descriptionKey || apiData.description_key || '',
    academicConceptId: apiData.acpetencyConceptId || apiData.academic_concept_id || '',
    isActive: apiData.isActive ?? apiData.is_active ?? true,
    level, // 计算出的等级
    createdAt: apiData.createdAt || apiData.created_at || new Date().toISOString(),
    updatedAt: apiData.updatedAt || apiData.updated_at || new Date().toISOString()
  };
}

/**
 * 转换 NodeRelationship 数据
 * @param {object} apiData - 后端返回的关系数据
 * @returns {object} 转换后的关系数据
 */
export function mapRelationship(apiData) {
  if (!apiData) return null;

  return {
    id: apiData.id || '',
    fromNodeId: apiData.fromNodeId || apiData.from_node_id || '',
    toNodeId: apiData.toNodeId || apiData.to_node_id || '',
    relationshipType: RELATIONSHIP_TYPE_MAP[apiData.relationshipType] ||
                      apiData.relationship_type ||
                      apiData.type ||
                      'PREREQ',
    weight: apiData.weight ?? apiData.weight ?? 1.0,
    createdAt: apiData.createdAt || apiData.created_at || new Date().toISOString(),
    updatedAt: apiData.updatedAt || apiData.updated_at || new Date().toISOString()
  };
}

/**
 * 批量转换数组数据
 * @param {Array} apiArray - 后端返回的数组数据
 * @param {Function} mapperFn - 转换函数
 * @returns {Array} 转换后的数组
 */
export function mapArray(apiArray, mapperFn) {
  if (!Array.isArray(apiArray)) return [];
  return apiArray.map(mapperFn).filter(item => item !== null);
}

/**
 * 将后端节点转换为图谱展示节点
 * @param {object} apiData - 后端返回的 Node 数据
 * @param {number} index - 节点索引（用于位置计算）
 * @param {number} total - 节点总数（用于位置计算）
 * @returns {object} 图谱展示节点
 */
export function mapToGraphNode(apiData, index = 0, total = 1) {
  const node = mapNode(apiData);
  if (!node) return null;

  // 计算展示位置
  const position = calculateNodePosition(index, total);

  return {
    id: node.id,
    name: node.nameChildKey || '未知节点',
    type: node.type,
    level: node.level,
    x: position.x,
    y: position.y,
    // 保留原始数据
    _raw: node
  };
}

/**
 * 计算节点在图谱中的位置
 * @param {number} index - 节点索引
 * @param {number} total - 节点总数
 * @param {object} options - 配置选项
 * @returns {object} 位置坐标 {x, y}
 */
export function calculateNodePosition(index, total, options = {}) {
  const {
    centerX = 400,
    centerY = 300,
    radiusX = 250,
    radiusY = 200,
    startAngle = 0
  } = options;

  if (total === 0) {
    return { x: centerX, y: centerY };
  }

  if (total === 1) {
    return { x: centerX, y: centerY };
  }

  const angle = startAngle + (index / total) * 2 * Math.PI;

  return {
    x: centerX + Math.cos(angle) * radiusX,
    y: centerY + Math.sin(angle) * radiusY
  };
}

/**
 * 将后端边数据转换为图谱边数据
 * @param {object} apiData - 后端返回的关系数据
 * @returns {object} 图谱边数据
 */
export function mapToGraphEdge(apiData) {
  const rel = mapRelationship(apiData);
  if (!rel) return null;

  return {
    source: rel.fromNodeId,
    target: rel.toNodeId,
    type: rel.relationshipType,
    weight: rel.weight
  };
}

/**
 * 获取节点类型的展示标签
 * @param {string} type - 节点类型
 * @returns {string} 展示标签
 */
export function getNodeTypeLabel(type) {
  const labels = {
    'C': '概念',
    'S': '支撑技能',
    'T': '思维模式',
    'P': '问题模型'
  };
  return labels[type] || '未知类型';
}

/**
 * 获取节点类型的展示颜色
 * @param {string} type - 节点类型
 * @returns {string} CSS 颜色类名
 */
export function getNodeTypeColor(type) {
  const colors = {
    'C': 'bg-green-100 text-green-800 border-green-300',
    'S': 'bg-blue-100 text-blue-800 border-blue-300',
    'T': 'bg-yellow-100 text-yellow-800 border-yellow-300',
    'P': 'bg-purple-100 text-purple-800 border-purple-300'
  };
  return colors[type] || 'bg-gray-100 text-gray-800 border-gray-300';
}

export default {
  mapDomain,
  mapTrunk,
  mapNode,
  mapRelationship,
  mapArray,
  mapToGraphNode,
  mapToGraphEdge,
  calculateNodePosition,
  getNodeTypeLabel,
  getNodeTypeColor
};
