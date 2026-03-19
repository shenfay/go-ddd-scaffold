package dto

import "github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"

// Pagination 分页参数（应用层统一使用）
type Pagination struct {
	Page     int `json:"page" form:"page" validate:"min=1"`
	PageSize int `json:"page_size" form:"page_size" validate:"min=1,max=100"`
}

// NewPagination 创建分页参数
func NewPagination(page, pageSize int) Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// ToDDDPagination 转换为 DDD 分页参数
func (p Pagination) ToDDDPagination() kernel.Pagination {
	return kernel.NewPagination(p.Page, p.PageSize)
}

// Offset 计算偏移量
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 返回限制数量
func (p Pagination) Limit() int {
	return p.PageSize
}

// PaginatedResult 分页结果（应用层统一使用）
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// NewPaginatedResult 创建分页结果
func NewPaginatedResult[T any](items []T, totalCount int64, page, pageSize int) *PaginatedResult[T] {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}
	return &PaginatedResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// FromDDDPaginatedResult 从 DDD 分页结果转换
func FromDDDPaginatedResult[T any](result *kernel.PaginatedResult[T]) *PaginatedResult[T] {
	return &PaginatedResult[T]{
		Items:      result.Items,
		TotalCount: result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPage,
	}
}
