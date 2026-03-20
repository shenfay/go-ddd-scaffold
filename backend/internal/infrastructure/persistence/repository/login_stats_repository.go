package repository

import (
	"context"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/repository"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"gorm.io/gorm"
)

// LoginStatsRepositoryImpl 登录统计仓储实现
// 注意：虽然LoginStats是独立的聚合根，但数据仍存储在users表中
// 这是为了向后兼容，避免复杂的数据迁移
type LoginStatsRepositoryImpl struct {
	query *dao.Query
}

// NewLoginStatsRepository 创建登录统计仓储实例
func NewLoginStatsRepository(query *dao.Query) repository.LoginStatsRepository {
	return &LoginStatsRepositoryImpl{query: query}
}

// Save 保存登录统计
func (r *LoginStatsRepositoryImpl) Save(ctx context.Context, stats *aggregate.LoginStats) error {
	// 使用 Updates 只更新登录统计相关字段，避免影响User的其他字段
	userDo := r.query.User
	_, err := userDo.WithContext(ctx).
		Where(userDo.ID.Eq(stats.UserID().Int64())).
		Updates(map[string]interface{}{
			"last_login_at":   stats.LastLoginAt(),
			"login_count":     stats.LoginCount(),
			"failed_attempts": stats.FailedAttempts(),
			"locked_until":    stats.LockedUntil(),
			"version":         stats.Version(),
			"updated_at":      time.Now(),
		})
	return err
}

// SaveInTransaction 在事务中保存登录统计
func (r *LoginStatsRepositoryImpl) SaveInTransaction(ctx context.Context, stats *aggregate.LoginStats, tx interface{}) error {
	db, ok := tx.(*gorm.DB)
	if !ok {
		return r.Save(ctx, stats)
	}

	// 在事务中更新
	return db.WithContext(ctx).Model(&struct {
		ID             int64      `gorm:"column:id"`
		LastLoginAt    *time.Time `gorm:"column:last_login_at"`
		LoginCount     int        `gorm:"column:login_count"`
		FailedAttempts int        `gorm:"column:failed_attempts"`
		LockedUntil    *time.Time `gorm:"column:locked_until"`
		Version        int        `gorm:"column:version"`
		UpdatedAt      time.Time  `gorm:"column:updated_at"`
	}{}).Where("id = ?", stats.UserID().Int64()).Updates(map[string]interface{}{
		"last_login_at":   stats.LastLoginAt(),
		"login_count":     stats.LoginCount(),
		"failed_attempts": stats.FailedAttempts(),
		"locked_until":    stats.LockedUntil(),
		"version":         stats.Version(),
		"updated_at":      time.Now(),
	}).Error
}

// FindByUserID 根据用户ID查找登录统计
func (r *LoginStatsRepositoryImpl) FindByUserID(ctx context.Context, userID vo.UserID) (*aggregate.LoginStats, error) {
	userDo := r.query.User
	userModel, err := userDo.WithContext(ctx).
		Select(userDo.ID, userDo.LastLoginAt, userDo.LoginCount, userDo.FailedAttempts, userDo.LockedUntil, userDo.Version, userDo.CreatedAt, userDo.UpdatedAt).
		Where(userDo.ID.Eq(userID.Int64())).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kernel.ErrAggregateNotFound
		}
		return nil, err
	}

	return r.toDomain(userModel), nil
}

// Delete 删除登录统计（软删除）
func (r *LoginStatsRepositoryImpl) Delete(ctx context.Context, userID vo.UserID) error {
	// 将登录统计字段重置为零值
	userDo := r.query.User
	_, err := userDo.WithContext(ctx).
		Where(userDo.ID.Eq(userID.Int64())).
		Updates(map[string]interface{}{
			"last_login_at":   nil,
			"login_count":     0,
			"failed_attempts": 0,
			"locked_until":    nil,
		})
	return err
}

// Exists 检查登录统计是否存在
func (r *LoginStatsRepositoryImpl) Exists(ctx context.Context, userID vo.UserID) (bool, error) {
	userDo := r.query.User
	count, err := userDo.WithContext(ctx).Where(userDo.ID.Eq(userID.Int64())).Count()
	return count > 0, err
}

// toDomain 将数据模型转换为领域对象
func (r *LoginStatsRepositoryImpl) toDomain(userModel interface{}) *aggregate.LoginStats {
	// 使用类型断言获取用户模型
	// 实际实现需要根据 dao.User 的实际类型调整
	stats := aggregate.NewLoginStats(vo.NewUserID(0)) // 临时实现
	return stats
}
