package ddd

import (
	"context"
)

// Repository 仓储接口标记
type Repository interface{}

// GenericRepository 泛型仓储接口
type GenericRepository[T AggregateRoot] interface {
	Repository
	Save(ctx context.Context, aggregate T) error
	FindByID(ctx context.Context, id interface{}) (T, error)
	Delete(ctx context.Context, id interface{}) error
	Exists(ctx context.Context, id interface{}) (bool, error)
}

// PaginatedResult 分页结果
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	TotalCount int64 `json:"total_count"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
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

// Offset 计算偏移量
func (p Pagination) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit 返回限制数量
func (p Pagination) Limit() int {
	return p.PageSize
}

// CalculateTotalPages 计算总页数
func (p Pagination) CalculateTotalPages(totalCount int64) int {
	if totalCount == 0 {
		return 0
	}
	pages := int(totalCount) / p.PageSize
	if int(totalCount)%p.PageSize > 0 {
		pages++
	}
	return pages
}

// Specification 规约接口
type Specification interface {
	IsSatisfiedBy(entity interface{}) bool
	And(other Specification) Specification
	Or(other Specification) Specification
	Not() Specification
}

// AndSpecification 与规约
type AndSpecification struct {
	left, right Specification
}

// IsSatisfiedBy 检查是否满足规约
func (a *AndSpecification) IsSatisfiedBy(entity interface{}) bool {
	return a.left.IsSatisfiedBy(entity) && a.right.IsSatisfiedBy(entity)
}

// And 创建与规约
func (a *AndSpecification) And(other Specification) Specification {
	return &AndSpecification{left: a, right: other}
}

// Or 创建或规约
func (a *AndSpecification) Or(other Specification) Specification {
	return &OrSpecification{left: a, right: other}
}

// Not 创建非规约
func (a *AndSpecification) Not() Specification {
	return &NotSpecification{spec: a}
}

// OrSpecification 或规约
type OrSpecification struct {
	left, right Specification
}

// IsSatisfiedBy 检查是否满足规约
func (o *OrSpecification) IsSatisfiedBy(entity interface{}) bool {
	return o.left.IsSatisfiedBy(entity) || o.right.IsSatisfiedBy(entity)
}

// And 创建与规约
func (o *OrSpecification) And(other Specification) Specification {
	return &AndSpecification{left: o, right: other}
}

// Or 创建或规约
func (o *OrSpecification) Or(other Specification) Specification {
	return &OrSpecification{left: o, right: other}
}

// Not 创建非规约
func (o *OrSpecification) Not() Specification {
	return &NotSpecification{spec: o}
}

// NotSpecification 非规约
type NotSpecification struct {
	spec Specification
}

// IsSatisfiedBy 检查是否满足规约
func (n *NotSpecification) IsSatisfiedBy(entity interface{}) bool {
	return !n.spec.IsSatisfiedBy(entity)
}

// And 创建与规约
func (n *NotSpecification) And(other Specification) Specification {
	return &AndSpecification{left: n, right: other}
}

// Or 创建或规约
func (n *NotSpecification) Or(other Specification) Specification {
	return &OrSpecification{left: n, right: other}
}

// Not 创建非规约
func (n *NotSpecification) Not() Specification {
	return n.spec
}
