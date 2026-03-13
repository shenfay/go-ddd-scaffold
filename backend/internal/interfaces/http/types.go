package http

// ==================== Swagger API 通用类型定义 ====================
// 本文件包含所有用于 Swagger 文档的通用类型定义
// 模块特有类型请参考各模块的 response.go
// ============================================================

// APIResponse 通用 API 响应结构
// @Description API 统一响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// PageMeta 分页元信息
// @Description 分页查询的元数据
type PageMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
