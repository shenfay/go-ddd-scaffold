package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础仓储操作
	Save(ctx context.Context, user *aggregate.User) error
	FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)
	Delete(ctx context.Context, id vo.UserID) error
	Exists(ctx context.Context, id vo.UserID) (bool, error)

	// 查询操作
	FindByUsername(ctx context.Context, username string) (*aggregate.User, error)
	FindByEmail(ctx context.Context, email string) (*aggregate.User, error)
	FindByStatus(ctx context.Context, status vo.UserStatus) ([]*aggregate.User, error)

	// 分页查询
	FindAll(ctx context.Context, pagination kernel.Pagination) (*kernel.PaginatedResult[*aggregate.User], error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status vo.UserStatus) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, users []*aggregate.User) error
	DeleteBatch(ctx context.Context, ids []vo.UserID) error

	// 乐观锁支持
	SaveWithVersion(ctx context.Context, user *aggregate.User, expectedVersion int) error

	// 事务支持 - 在事务中执行操作
	SaveInTransaction(ctx context.Context, user *aggregate.User, tx interface{}) error
}

// LoginStatsRepository 登录统计仓储接口
type LoginStatsRepository interface {
	// 基础仓储操作
	Save(ctx context.Context, stats *aggregate.LoginStats) error
	FindByUserID(ctx context.Context, userID vo.UserID) (*aggregate.LoginStats, error)
	Delete(ctx context.Context, userID vo.UserID) error
	Exists(ctx context.Context, userID vo.UserID) (bool, error)

	// 事务支持
	SaveInTransaction(ctx context.Context, stats *aggregate.LoginStats, tx interface{}) error
}
