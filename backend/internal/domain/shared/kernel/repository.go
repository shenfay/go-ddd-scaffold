package kernel

import (
	"context"
)

// Repository 仓储接口
type Repository interface {
	Save(ctx context.Context, aggregate AggregateRoot) error
	FindByID(ctx context.Context, id interface{}) (AggregateRoot, error)
	Delete(ctx context.Context, id interface{}) error
}

// GenericRepository 泛型仓储接口
type GenericRepository[T AggregateRoot] interface {
	Save(ctx context.Context, entity T) error
	FindByID(ctx context.Context, id interface{}) (T, error)
	Delete(ctx context.Context, id interface{}) error
}

// Pagination 分页参数
type Pagination struct {
	Page     int
	PageSize int
}

// NewPagination 创建分页参数
func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Offset 计算偏移量
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 返回每页大小
func (p Pagination) Limit() int {
	return p.PageSize
}

// PaginatedResult 分页结果
type PaginatedResult[T any] struct {
	Items     []T   `json:"items"`
	Total     int64 `json:"total"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalPage int   `json:"total_page"`
}

// NewPaginatedResult 创建分页结果
func NewPaginatedResult[T any](items []T, total int64, page, pageSize int) *PaginatedResult[T] {
	totalPage := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPage++
	}
	return &PaginatedResult[T]{
		Items:     items,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}
}
