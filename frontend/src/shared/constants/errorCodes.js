/**
 * 前端错误码常量定义
 * 与后端 errors.go 保持一致
 */

// 通用错误码
export const ERROR_CODES = {
  // 成功
  SUCCESS: 'Success',
  
  // 通用错误
  INVALID_PARAMETER: 'InvalidParameter',
  MISSING_PARAMETER: 'MissingParameter',
  UNAUTHORIZED: 'Unauthorized',
  FORBIDDEN: 'Forbidden',
  NOT_FOUND: 'NotFound',
  METHOD_NOT_ALLOWED: 'MethodNotAllowed',
  TOO_MANY_REQUESTS: 'TooManyRequests',
  VALIDATION_FAILED: 'ValidationFailed',
  RESOURCE_CONFLICT: 'ResourceConflict',
  UNSUPPORTED_MEDIA_TYPE: 'UnsupportedMediaType',
  
  // 知识图谱 Domain 错误
  KG_DOMAIN_NOT_FOUND: 'KG.Domain.NotFound',
  KG_DOMAIN_ALREADY_EXISTS: 'KG.Domain.AlreadyExists',
  KG_DOMAIN_INVALID_DATA: 'KG.Domain.InvalidData',
  
  // 知识图谱 Trunk 错误
  KG_TRUNK_NOT_FOUND: 'KG.Trunk.NotFound',
  KG_TRUNK_ALREADY_EXISTS: 'KG.Trunk.AlreadyExists',
  KG_TRUNK_NOT_IN_DOMAIN: 'KG.Trunk.NotInDomain',
  KG_TRUNK_INVALID_DATA: 'KG.Trunk.InvalidData',
  
  // 知识图谱 Node 错误
  KG_NODE_NOT_FOUND: 'KG.Node.NotFound',
  KG_NODE_ALREADY_EXISTS: 'KG.Node.AlreadyExists',
  KG_NODE_INVALID_TYPE: 'KG.Node.InvalidType',
  KG_NODE_NOT_IN_TRUNK: 'KG.Node.NotInTrunk',
  KG_NODE_INVALID_DATA: 'KG.Node.InvalidData',
  
  // 知识图谱 Relationship 错误
  KG_RELATIONSHIP_NOT_FOUND: 'KG.Relationship.NotFound',
  KG_RELATIONSHIP_ALREADY_EXISTS: 'KG.Relationship.AlreadyExists',
  KG_RELATIONSHIP_INVALID_TYPE: 'KG.Relationship.InvalidType',
  KG_RELATIONSHIP_CYCLE_DETECTED: 'KG.Relationship.CycleDetected',
  KG_RELATIONSHIP_INVALID_DATA: 'KG.Relationship.InvalidData',
  
  // 系统错误
  SYSTEM_INTERNAL_ERROR: 'System.InternalError',
  SYSTEM_DATABASE_ERROR: 'System.DatabaseError',
  SYSTEM_CACHE_UNAVAILABLE: 'System.CacheUnavailable',
  SYSTEM_EXTERNAL_SERVICE_ERROR: 'System.ExternalServiceError',
  SYSTEM_TIMEOUT: 'System.Timeout'
};

// 错误分类
export const ERROR_CATEGORIES = {
  COMMON: 'Common',
  KG: 'KG',
  SYSTEM: 'System',
  AUTH: 'Auth',
  USER: 'User'
};

// HTTP 状态码映射
export const HTTP_STATUS_MAP = {
  [ERROR_CODES.INVALID_PARAMETER]: 400,
  [ERROR_CODES.MISSING_PARAMETER]: 400,
  [ERROR_CODES.VALIDATION_FAILED]: 400,
  [ERROR_CODES.KG_DOMAIN_INVALID_DATA]: 400,
  [ERROR_CODES.KG_TRUNK_INVALID_DATA]: 400,
  [ERROR_CODES.KG_NODE_INVALID_DATA]: 400,
  [ERROR_CODES.KG_NODE_INVALID_TYPE]: 400,
  [ERROR_CODES.KG_NODE_NOT_IN_TRUNK]: 400,
  [ERROR_CODES.KG_RELATIONSHIP_INVALID_DATA]: 400,
  [ERROR_CODES.KG_RELATIONSHIP_INVALID_TYPE]: 400,
  [ERROR_CODES.KG_RELATIONSHIP_CYCLE_DETECTED]: 400,
  [ERROR_CODES.UNAUTHORIZED]: 401,
  [ERROR_CODES.FORBIDDEN]: 403,
  [ERROR_CODES.NOT_FOUND]: 404,
  [ERROR_CODES.KG_DOMAIN_NOT_FOUND]: 404,
  [ERROR_CODES.KG_TRUNK_NOT_FOUND]: 404,
  [ERROR_CODES.KG_NODE_NOT_FOUND]: 404,
  [ERROR_CODES.KG_RELATIONSHIP_NOT_FOUND]: 404,
  [ERROR_CODES.METHOD_NOT_ALLOWED]: 405,
  [ERROR_CODES.RESOURCE_CONFLICT]: 409,
  [ERROR_CODES.KG_DOMAIN_ALREADY_EXISTS]: 409,
  [ERROR_CODES.KG_TRUNK_ALREADY_EXISTS]: 409,
  [ERROR_CODES.KG_NODE_ALREADY_EXISTS]: 409,
  [ERROR_CODES.KG_RELATIONSHIP_ALREADY_EXISTS]: 409,
  [ERROR_CODES.UNSUPPORTED_MEDIA_TYPE]: 415,
  [ERROR_CODES.TOO_MANY_REQUESTS]: 429,
  [ERROR_CODES.SYSTEM_INTERNAL_ERROR]: 500,
  [ERROR_CODES.SYSTEM_DATABASE_ERROR]: 500,
  [ERROR_CODES.SYSTEM_EXTERNAL_SERVICE_ERROR]: 502,
  [ERROR_CODES.SYSTEM_CACHE_UNAVAILABLE]: 503,
  [ERROR_CODES.SYSTEM_TIMEOUT]: 504
};

// 错误消息映射（用于前端显示）
export const ERROR_MESSAGES = {
  [ERROR_CODES.SUCCESS]: '操作成功',
  [ERROR_CODES.INVALID_PARAMETER]: '无效请求，请检查输入参数',
  [ERROR_CODES.MISSING_PARAMETER]: '缺少必要参数',
  [ERROR_CODES.UNAUTHORIZED]: '未授权，请先登录',
  [ERROR_CODES.FORBIDDEN]: '禁止访问',
  [ERROR_CODES.NOT_FOUND]: '请求的资源不存在',
  [ERROR_CODES.METHOD_NOT_ALLOWED]: '不支持的请求方法',
  [ERROR_CODES.TOO_MANY_REQUESTS]: '请求过于频繁，请稍后再试',
  [ERROR_CODES.VALIDATION_FAILED]: '参数校验失败',
  [ERROR_CODES.RESOURCE_CONFLICT]: '资源冲突',
  [ERROR_CODES.UNSUPPORTED_MEDIA_TYPE]: '不支持的媒体类型',
  
  // KG Domain
  [ERROR_CODES.KG_DOMAIN_NOT_FOUND]: '知识领域不存在',
  [ERROR_CODES.KG_DOMAIN_ALREADY_EXISTS]: '知识领域已存在',
  [ERROR_CODES.KG_DOMAIN_INVALID_DATA]: '知识领域数据无效',
  
  // KG Trunk
  [ERROR_CODES.KG_TRUNK_NOT_FOUND]: '知识主线不存在',
  [ERROR_CODES.KG_TRUNK_ALREADY_EXISTS]: '知识主线已存在',
  [ERROR_CODES.KG_TRUNK_NOT_IN_DOMAIN]: '知识主线不属于指定的领域',
  [ERROR_CODES.KG_TRUNK_INVALID_DATA]: '知识主线数据无效',
  
  // KG Node
  [ERROR_CODES.KG_NODE_NOT_FOUND]: '知识节点不存在',
  [ERROR_CODES.KG_NODE_ALREADY_EXISTS]: '知识节点已存在',
  [ERROR_CODES.KG_NODE_INVALID_TYPE]: '无效的节点类型，必须是 C/S/T/P 之一',
  [ERROR_CODES.KG_NODE_NOT_IN_TRUNK]: '知识节点不属于指定的主线',
  [ERROR_CODES.KG_NODE_INVALID_DATA]: '知识节点数据无效',
  
  // KG Relationship
  [ERROR_CODES.KG_RELATIONSHIP_NOT_FOUND]: '知识关系不存在',
  [ERROR_CODES.KG_RELATIONSHIP_ALREADY_EXISTS]: '知识关系已存在',
  [ERROR_CODES.KG_RELATIONSHIP_INVALID_TYPE]: '无效的关系类型，必须是 PREREQ/SUP_SKILL/THINK_PAT 之一',
  [ERROR_CODES.KG_RELATIONSHIP_CYCLE_DETECTED]: '检测到循环引用，无法建立关系',
  [ERROR_CODES.KG_RELATIONSHIP_INVALID_DATA]: '知识关系数据无效',
  
  // System
  [ERROR_CODES.SYSTEM_INTERNAL_ERROR]: '系统内部错误，请稍后重试',
  [ERROR_CODES.SYSTEM_DATABASE_ERROR]: '数据库操作失败',
  [ERROR_CODES.SYSTEM_CACHE_UNAVAILABLE]: '缓存服务不可用',
  [ERROR_CODES.SYSTEM_EXTERNAL_SERVICE_ERROR]: '外部服务调用失败',
  [ERROR_CODES.SYSTEM_TIMEOUT]: '请求超时，请稍后重试'
};

// 错误类型判断函数
export const isErrorType = {
  // 是否为找不到资源错误
  isNotFound: (errorCode) => {
    return errorCode?.endsWith('.NotFound') || errorCode === ERROR_CODES.NOT_FOUND;
  },
  
  // 是否为参数校验错误
  isValidation: (errorCode) => {
    return [
      ERROR_CODES.VALIDATION_FAILED,
      ERROR_CODES.INVALID_PARAMETER,
      ERROR_CODES.MISSING_PARAMETER
    ].includes(errorCode);
  },
  
  // 是否为资源冲突错误
  isConflict: (errorCode) => {
    return errorCode?.endsWith('.AlreadyExists') || errorCode === ERROR_CODES.RESOURCE_CONFLICT;
  },
  
  // 是否为服务器错误
  isServer: (errorCode) => {
    return errorCode?.startsWith('System.') || [
      ERROR_CODES.SYSTEM_INTERNAL_ERROR,
      ERROR_CODES.SYSTEM_DATABASE_ERROR
    ].includes(errorCode);
  },
  
  // 是否为知识图谱相关错误
  isKG: (errorCode) => {
    return errorCode?.startsWith('KG.');
  },
  
  // 是否为认证相关错误
  isAuth: (errorCode) => {
    return [
      ERROR_CODES.UNAUTHORIZED,
      ERROR_CODES.FORBIDDEN
    ].includes(errorCode);
  }
};

// 获取错误对应的 HTTP 状态码
export function getHttpStatus(errorCode) {
  return HTTP_STATUS_MAP[errorCode] || 400;
}

// 获取错误消息
export function getErrorMessage(errorCode, defaultMessage = '未知错误') {
  return ERROR_MESSAGES[errorCode] || defaultMessage;
}

// 获取错误分类
export function getErrorCategory(errorCode) {
  if (!errorCode) return ERROR_CATEGORIES.COMMON;
  
  const parts = errorCode.split('.');
  if (parts.length > 0) {
    switch (parts[0]) {
      case 'KG':
        return ERROR_CATEGORIES.KG;
      case 'System':
        return ERROR_CATEGORIES.SYSTEM;
      case 'Auth':
        return ERROR_CATEGORIES.AUTH;
      case 'User':
        return ERROR_CATEGORIES.USER;
      default:
        return ERROR_CATEGORIES.COMMON;
    }
  }
  
  return ERROR_CATEGORIES.COMMON;
}