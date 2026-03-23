package application

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	persistenceRepo "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/transaction"
	"gorm.io/gorm"
)

// UnitOfWork 工作单元接口，用于管理事务和仓储
type UnitOfWork interface {
	// Transaction 在事务中执行函数
	Transaction(ctx context.Context, fn func(context.Context) error) error
	// UserRepository 返回用户仓储
	UserRepository() repository.UserRepository
	// LoginStatsRepository 返回登录统计仓储
	LoginStatsRepository() repository.LoginStatsRepository
}

// unitOfWork 工作单元实现
type unitOfWork struct {
	db             *gorm.DB
	query          *dao.Query
	userRepo       repository.UserRepository
	loginStatsRepo repository.LoginStatsRepository
}

// NewUnitOfWork 创建工作单元实例
func NewUnitOfWork(db *gorm.DB, query *dao.Query) UnitOfWork {
	return &unitOfWork{
		db:             db,
		query:          query,
		userRepo:       persistenceRepo.NewUserRepository(query),
		loginStatsRepo: persistenceRepo.NewLoginStatsRepository(query),
	}
}

// Transaction 在事务中执行函数
func (u *unitOfWork) Transaction(ctx context.Context, fn func(context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务注入上下文
		txCtx := transaction.WithTransaction(ctx, tx)
		return fn(txCtx)
	})
}

// UserRepository 返回用户仓储
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return u.userRepo
}

// LoginStatsRepository 返回登录统计仓储
func (u *unitOfWork) LoginStatsRepository() repository.LoginStatsRepository {
	return u.loginStatsRepo
}
