package handler

// APIResponse 通用 API 响应结构
// @Description API 统一响应格式（用于 Swagger 文档）
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

// PageMeta 分页元信息
// @Description 分页查询的元数据（用于 Swagger 文档）
type PageMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
