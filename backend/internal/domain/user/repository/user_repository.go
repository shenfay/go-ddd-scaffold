package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/model"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	ddd.Repository

	// 基础仓储操作
	Save(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id model.UserID) (*model.User, error)
	Delete(ctx context.Context, id model.UserID) error
	Exists(ctx context.Context, id model.UserID) (bool, error)

	// 查询操作
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByStatus(ctx context.Context, status model.UserStatus) ([]*model.User, error)

	// 分页查询
	FindAll(ctx context.Context, pagination ddd.Pagination) (*ddd.PaginatedResult[*model.User], error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status model.UserStatus) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, users []*model.User) error
	DeleteBatch(ctx context.Context, ids []model.UserID) error

	// 乐观锁支持
	SaveWithVersion(ctx context.Context, user *model.User, expectedVersion int) error
}
