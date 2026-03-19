package repository

import (
	"context"
	"errors"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"gorm.io/gorm"
)

// DAOQuery 定义 GORM Gen DAO 的通用接口
type DAOQuery interface {
	WithContext(ctx context.Context) interface{}
}

// DAOFinder 定义查找接口
type DAOFinder[M any] interface {
	Where(query interface{}, args ...interface{}) DAOFinder[M]
	First() (M, error)
	Count() (int64, error)
}

// DomainConverter 领域对象与数据模型转换接口
type DomainConverter[T any, M any] interface {
	// ToDomain 将数据模型转换为领域对象
	ToDomain(model M) T
	// FromDomain 将领域对象转换为数据模型
	FromDomain(domain T) M
}

// BaseRepository 泛型仓储基类
// T: 领域对象类型
// M: 数据模型类型
type BaseRepository[T any, M any] struct {
	db        *gorm.DB
	converter DomainConverter[T, M]
}

// NewBaseRepository 创建基础仓储
func NewBaseRepository[T any, M any](db *gorm.DB, converter DomainConverter[T, M]) *BaseRepository[T, M] {
	return &BaseRepository[T, M]{
		db:        db,
		converter: converter,
	}
}

// DB 获取 GORM DB 实例（供子类使用）
func (r *BaseRepository[T, M]) DB() *gorm.DB {
	return r.db
}

// Converter 获取转换器（供子类使用）
func (r *BaseRepository[T, M]) Converter() DomainConverter[T, M] {
	return r.converter
}

// HandleNotFound 统一处理记录不存在错误
func (r *BaseRepository[T, M]) HandleNotFound(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return kernel.ErrAggregateNotFound
	}
	return err
}

// CheckRowsAffected 检查影响行数
func (r *BaseRepository[T, M]) CheckRowsAffected(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return kernel.ErrAggregateNotFound
	}
	return nil
}
