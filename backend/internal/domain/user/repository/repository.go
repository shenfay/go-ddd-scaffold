package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// 类型别名，方便使用
type (
	UserID     = user.UserID
	UserStatus = user.UserStatus
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 基础仓储操作
	Save(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id UserID) (*user.User, error)
	Delete(ctx context.Context, id UserID) error
	Exists(ctx context.Context, id UserID) (bool, error)

	// 查询操作
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByStatus(ctx context.Context, status UserStatus) ([]*user.User, error)

	// 分页查询
	FindAll(ctx context.Context, pagination common.Pagination) (*common.PaginatedResult[*user.User], error)

	// 统计操作
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status UserStatus) (int64, error)

	// 批量操作
	SaveBatch(ctx context.Context, users []*user.User) error
	DeleteBatch(ctx context.Context, ids []UserID) error

	// 乐观锁支持
	SaveWithVersion(ctx context.Context, user *user.User, expectedVersion int) error

	// 事务支持 - 在事务中执行操作
	SaveInTransaction(ctx context.Context, user *user.User, tx interface{}) error
}

// LoginStatsRepository 登录统计仓储接口
type LoginStatsRepository interface {
	// 基础仓储操作
	Save(ctx context.Context, stats *user.LoginStats) error
	FindByUserID(ctx context.Context, userID UserID) (*user.LoginStats, error)
	Delete(ctx context.Context, userID UserID) error
	Exists(ctx context.Context, userID UserID) (bool, error)

	// 事务支持
	SaveInTransaction(ctx context.Context, stats *user.LoginStats, tx interface{}) error
}
