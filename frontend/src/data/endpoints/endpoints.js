/**
 * API 端点配置
 * 
 * 集中管理所有 API 端点，便于维护和修改
 * 按照业务模块进行分类组织
 */

const API_ENDPOINTS = {
  // 用户相关端点
  user: {
    login: '/users/login',
    logout: '/users/logout',
    register: '/users/register',
    getInfo: '/users/info',
    profile: '/users/profile',
    updateProfile: '/users/profile',
    changePassword: '/users/change-password'
  },
  
  // 租户相关端点
  tenant: {
    userTenants: '/tenants/my-tenants',  // 获取用户的租户列表
    select: '/tenants/select',  // 选择当前租户
    list: '/tenants',  // 租户列表
    create: '/tenants',  // 创建租户
    detail: '/tenants/:id'  // 租户详情
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
