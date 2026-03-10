/**
 * API 响应格式化器
 * 
 * 负责处理后端返回的统一响应格式，提取 data 字段
 * 支持 Success/Created/NoContent 等不同状态码
 */

/**
 * 格式化成功响应
 * @param {Object} response - 原始响应对象
 * @returns {any} 返回 data 字段内容
 */
export function formatSuccessResponse(response) {
  if (!response || !response.data) {
   return null;
  }

  const { code, message, data, requestId, timestamp } = response.data;
  
  // 根据 code 判断响应类型
  switch (code) {
    case'Success':
    case 'Created':
      // 有 data 字段，直接返回
     return data;
    
    case 'NoContent':
      // 204 No Content，返回 null
     return null;
    
    default:
      // 其他情况，尝试返回 data
     return data || response.data;
  }
}

/**
 * 格式化错误响应
 * @param {Object} error - 错误对象
 * @returns {Object} 格式化的错误信息
 */
export function formatErrorResponse(error) {
  if (!error || !error.response) {
   return {
     code: 'NetworkError',
     message: error.message || '网络错误',
      details: null
    };
  }

  const responseData = error.response.data;
  
  // 后端返回的统一错误格式
  if (responseData && typeof responseData === 'object') {
   return {
     code: responseData.code || 'UnknownError',
     message: responseData.message || '未知错误',
      details: responseData.error?.details || null,
     requestId: responseData.requestId,
      timestamp: responseData.timestamp
    };
  }

  // 非 JSON 响应
  return {
   code: 'HTTPError',
   message: error.response.statusText || 'HTTP 错误',
    details: null,
   status: error.response.status
  };
}

/**
 * 判断响应是否为成功状态
 * @param {Object} response - 响应对象
 * @returns {boolean}
 */
export function isSuccessResponse(response) {
  if (!response || !response.data) {
   return false;
  }

  const { code } = response.data;
  return ['Success', 'Created', 'NoContent'].includes(code);
}

/**
 * 提取分页数据
 * @param {Object} response - 响应对象
 * @returns {Object} 分页数据
 */
export function extractPageData(response) {
  const data = formatSuccessResponse(response);
  
  if (!data || !data.items) {
   return {
      items: [],
      total: 0,
      page: 1,
      pageSize: 10,
      totalPages: 0,
      hasNext: false,
      hasPrev: false
    };
  }

  return {
    items: data.items || [],
    total: data.total || 0,
    page: data.page || 1,
    pageSize: data.pageSize || 10,
    totalPages: Math.ceil(data.total/ data.pageSize),
    hasNext: data.page < Math.ceil(data.total/ data.pageSize),
    hasPrev: data.page > 1
  };
}
