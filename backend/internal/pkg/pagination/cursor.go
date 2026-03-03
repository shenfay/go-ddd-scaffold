package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// PageRequest 分页请求参数（支持page+pageSize和limit+offset两种方式）
type PageRequest struct {
	Page     int    `form:"page" json:"page"`         // 页码，从1开始
	PageSize int    `form:"pageSize" json:"pageSize"` // 每页数量
	Limit    int    `form:"limit" json:"limit"`       // 每页数量（兼容）
	Offset   int    `form:"offset" json:"offset"`     // 偏移量（兼容）
	Sort     string `form:"sort" json:"sort"`         // 排序字段

}

// Cursor 游标分页游标
type Cursor struct {
	Offset    int64     `json:"offset"`    // 偏移量
	Timestamp time.Time `json:"timestamp"` // 时间戳
	ID        string    `json:"id"`        // 最后一条记录ID
}

// PageResponse 游标分页响应
type PageResponse struct {
	Items      interface{} `json:"items"`       // 数据列表
	NextCursor string      `json:"next_cursor"` // 下一页游标
	HasMore    bool        `json:"has_more"`    // 是否有更多数据
	Total      int64       `json:"total"`       // 总数
}

// DefaultLimit 默认每页数量
const DefaultLimit = 20
const MaxLimit = 100

// ParseCursor 解析游标
func ParseCursor(cursor string) (*Cursor, error) {
	if cursor == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format")
	}

	var c Cursor
	if err := json.Unmarshal(decoded, &c); err != nil {
		return nil, fmt.Errorf("failed to parse cursor")
	}

	return &c, nil
}

// EncodeCursor 编码游标
func EncodeCursor(cursor *Cursor) string {
	if cursor == nil {
		return ""
	}

	data, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(data)
}

// GetLimit 获取分页大小
func (r *PageRequest) GetLimit() int {
	// 优先使用PageSize
	if r.PageSize > 0 {
		if r.PageSize > MaxLimit {
			return MaxLimit
		}
		return r.PageSize
	}
	// 兼容Limit
	if r.Limit > 0 {
		if r.Limit > MaxLimit {
			return MaxLimit
		}
		return r.Limit
	}
	return DefaultLimit
}

// GetPage 获取页码
func (r *PageRequest) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

// GetOffset 计算偏移量
func (r *PageRequest) GetOffset() int {
	// 如果直接提供了Offset则使用
	if r.Offset > 0 {
		return r.Offset
	}
	// 否则根据page计算
	return (r.GetPage() - 1) * r.GetLimit()
}

// GetSort 获取排序参数
// 返回格式：字段前加 "-" 表示降序，否则升序
// 多字段用逗号分隔，如 "-priority,name"
func (r *PageRequest) GetSort() string {
	return r.Sort
}

// ParseSort 解析排序参数
// 返回：(字段名列表, 是否降序列表, 错误)
func (r *PageRequest) ParseSort() ([]string, []bool, error) {
	if r.Sort == "" {
		return nil, nil, nil
	}

	fields := strings.Split(r.Sort, ",")
	fieldNames := make([]string, 0, len(fields))
	descending := make([]bool, 0, len(fields))

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		// 判断是否降序
		if strings.HasPrefix(field, "-") {
			fieldNames = append(fieldNames, strings.TrimPrefix(field, "-"))
			descending = append(descending, true)
		} else {
			fieldNames = append(fieldNames, field)
			descending = append(descending, false)
		}
	}

	if len(fieldNames) == 0 {
		return nil, nil, nil
	}

	return fieldNames, descending, nil
}

// NewPageRequest 创建分页请求
func NewPageRequest(cursor string, limit int, sort string) *PageRequest {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return &PageRequest{
		Limit: limit,
		Sort:  sort,
	}
}

// BuildNextCursor 构建下一页游标
func BuildNextCursor(items interface{}, limit int, total int64) *PageResponse {
	resp := &PageResponse{
		Items:   items,
		HasMore: false,
		Total:   total,
	}

	if items == nil {
		return resp
	}

	switch v := items.(type) {
	case []interface{}:
		if len(v) == limit {
			resp.HasMore = true
			nextCursor := &Cursor{
				Offset: int64(len(v)),
				ID:     fmt.Sprintf("page_%d", len(v)),
			}
			resp.NextCursor = EncodeCursor(nextCursor)
		}
	}

	return resp
}

// OffsetBasedPageRequest 基于偏移量的分页请求
type OffsetBasedPageRequest struct {
	Page     int `form:"page" json:"page"`           // 页码，从1开始
	PageSize int `form:"page_size" json:"page_size"` // 每页数量
}

// DefaultPageSize 默认每页数量
const DefaultPageSize = 20
const MaxPageSize = 100

// NewOffsetBasedPageRequest 创建基于偏移量的分页请求
func NewOffsetBasedPageRequest(page, pageSize int) *OffsetBasedPageRequest {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return &OffsetBasedPageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

// GetOffset 获取偏移量
func (r *OffsetBasedPageRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// GetLimit 获取限制数量
func (r *OffsetBasedPageRequest) GetLimit() int {
	return r.PageSize
}

// OffsetBasedPageResponse 基于偏移量的分页响应
type OffsetBasedPageResponse struct {
	Items     interface{} `json:"items"`     // 数据列表
	Page      int         `json:"page"`      // 当前页码
	PageSize  int         `json:"pageSize"`  // 每页数量
	Total     int64       `json:"total"`     // 总数
	TotalPage int         `json:"totalPage"` // 总页数
	HasNext   bool        `json:"hasNext"`   // 是否有下一页
	HasPrev   bool        `json:"hasPrev"`   // 是否有上一页
}

// BuildOffsetBasedResponse 构建基于偏移量的分页响应
func BuildOffsetBasedResponse(items interface{}, page, pageSize int, total int64) *OffsetBasedPageResponse {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}

	return &OffsetBasedPageResponse{
		Items:     items,
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPage,
		HasNext:   page < totalPage,
		HasPrev:   page > 1,
	}
}
