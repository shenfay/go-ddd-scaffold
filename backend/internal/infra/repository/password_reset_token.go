package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"gorm.io/gorm"
)

// PasswordResetTokenPO 密码重置令牌持久化对象
type PasswordResetTokenPO struct {
	ID        string    `gorm:"column:id;type:varchar(26);primaryKey"`
	UserID    string    `gorm:"column:user_id;type:varchar(26);not null;index"`
	Token     string    `gorm:"column:token;type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	Used      bool      `gorm:"column:used;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName 指定表名
func (PasswordResetTokenPO) TableName() string {
	return "password_reset_tokens"
}

// PasswordResetTokenRepositoryImpl 密码重置令牌仓储实现
type PasswordResetTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository 创建密码重置令牌仓储实例
func NewPasswordResetTokenRepository(db *gorm.DB) *PasswordResetTokenRepositoryImpl {
	return &PasswordResetTokenRepositoryImpl{db: db}
}

// Create 创建密码重置令牌
func (r *PasswordResetTokenRepositoryImpl) Create(ctx context.Context, token *authentication.PasswordResetToken) error {
	po := &PasswordResetTokenPO{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		Used:      token.Used,
		CreatedAt: token.CreatedAt,
	}
	return r.db.WithContext(ctx).Create(&po).Error
}

// FindByToken 根据令牌字符串查找
// 自动过滤已使用和已过期的令牌
func (r *PasswordResetTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*authentication.PasswordResetToken, error) {
	var po PasswordResetTokenPO
	if err := r.db.WithContext(ctx).
		Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).
		First(&po).Error; err != nil {
		return nil, err
	}
	return toDomainToken(&po), nil
}

// MarkAsUsed 标记令牌已使用
func (r *PasswordResetTokenRepositoryImpl) MarkAsUsed(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).
		Model(&PasswordResetTokenPO{}).
		Where("id = ?", tokenID).
		Update("used", true).Error
}

// DeleteExpired 删除所有过期令牌
func (r *PasswordResetTokenRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&PasswordResetTokenPO{}).Error
}

// toDomainToken 持久化对象转领域对象
func toDomainToken(po *PasswordResetTokenPO) *authentication.PasswordResetToken {
	return &authentication.PasswordResetToken{
		ID:        po.ID,
		UserID:    po.UserID,
		Token:     po.Token,
		ExpiresAt: po.ExpiresAt,
		Used:      po.Used,
		CreatedAt: po.CreatedAt,
	}
}
