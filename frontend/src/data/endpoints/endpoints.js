/**
 * API 端点配置
 * 
 * 集中管理所有 API 端点，便于维护和修改
 * 按照业务模块进行分类组织
 */

const API_ENDPOINTS = {
  // 用户相关端点
  user: {
    login: '/user/login',
    logout: '/user/logout',
    register: '/user/register',
    profile: '/user/profile',
    updateProfile: '/user/profile',
    getInfo: '/user/info',
    changePassword: '/user/change-password'
  },

  // 学习相关端点
  learning: {
    getLessons: '/learning/lessons',
    getLesson: '/learning/lessons/:id',
    getLearningProgress: '/learning/progress',
    submitLesson: '/learning/lessons/:id/submit',
    getKnowledgeNode: '/learning/knowledge-nodes/:id',
    queryKnowledgeNodes: '/learning/knowledge-nodes'
  },

  // 知识图谱相关端点
  knowledge: {
    getGraph: '/knowledge/graph',
    getNode: '/knowledge/nodes/:id',
    getConnections: '/knowledge/nodes/:id/connections',
    searchNodes: '/knowledge/search'
  },

  // 知识图谱 KG 相关端点
  kg: {
    // 领域 (Domain)
    getDomains: '/knowledge/domains',
    getDomain: '/knowledge/domains/:id',

    // 主线 (Trunk)
    getTrunksByDomain: '/knowledge/trunks/domain/:domainId',
    getTrunk: '/knowledge/trunks/:id',

    // 节点 (Node)
    getNodesByTrunk: '/knowledge/nodes/:trunkId',
    getNode: '/knowledge/node/:id',

    // 节点关系
    getNodeRelationships: '/knowledge/node/:id/relationships',
    getNodePrerequisites: '/knowledge/node/:id/prerequisites',
    getNodeDependents: '/knowledge/node/:id/dependents'
  },

  // 任务和成就相关端点
  tasks: {
    getTasks: '/tasks',
    getTask: '/tasks/:id',
    completeTask: '/tasks/:id/complete',
    getAchievements: '/achievements',
    getAchievement: '/achievements/:id'
  },

  // 素材和资源相关端点
  resources: {
    getAssets: '/resources/assets',
    getAsset: '/resources/assets/:id',
    uploadAsset: '/resources/assets/upload',
    getModels: '/resources/models',
    getModel: '/resources/models/:id',
    getSounds: '/resources/sounds',
    getSound: '/resources/sounds/:id'
  },

  // 父母端相关端点
  parent: {
    getDashboard: '/parent/dashboard',
    getChildProgress: '/parent/children/:childId/progress',
    getChildStats: '/parent/children/:childId/stats',
    getNotifications: '/parent/notifications',
    markNotificationRead: '/parent/notifications/:id/read'
  },

  // 游戏相关端点
  game: {
    getScene: '/game/scenes/:id',
    saveGameState: '/game/state',
    getGameState: '/game/state',
    getLeaderboard: '/game/leaderboard',
    submitScore: '/game/scores'
  },

  // 系统相关端点
  system: {
    getConfig: '/system/config',
    getStatus: '/system/status',
    getVersion: '/system/version'
  }
};

/**
 * 替换路径参数
 * @param {string} path - API 路径，可能包含 :id 等参数
 * @param {object} params - 参数对象
 * @returns {string} 替换后的路径
 */
export function replacePathParams(path, params = {}) {
  let resultPath = path;
  Object.entries(params).forEach(([key, value]) => {
    resultPath = resultPath.replace(`:${key}`, value);
  });
  return resultPath;
}

/**
 * 获取完整的 API 路径
 * @param {string} endpoint - 端点名称
 * @param {object} params - 路径参数
 * @returns {string} 完整的 API 路径
 */
export function getEndpoint(endpoint, params = {}) {
  const pathKeys = endpoint.split('.');
  let path = API_ENDPOINTS;

  for (const key of pathKeys) {
    if (path[key]) {
      path = path[key];
    } else {
      console.warn(`Endpoint not found: ${endpoint}`);
      return null;
    }
  }

  if (typeof path === 'string') {
    return replacePathParams(path, params);
  }

  return path;
}

export default API_ENDPOINTS;
