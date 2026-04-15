package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"gorm.io/gorm"
)

// EmailVerificationTokenPO 邮箱验证令牌持久化对象
type EmailVerificationTokenPO struct {
	ID        string    `gorm:"column:id;type:varchar(26);primaryKey"`
	UserID    string    `gorm:"column:user_id;type:varchar(26);not null;index"`
	Token     string    `gorm:"column:token;type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	Used      bool      `gorm:"column:used;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// TableName 指定表名
func (EmailVerificationTokenPO) TableName() string {
	return "email_verification_tokens"
}

// EmailVerificationTokenRepositoryImpl 邮箱验证令牌仓储实现
type EmailVerificationTokenRepositoryImpl struct {
	db *gorm.DB
}

// NewEmailVerificationTokenRepository 创建邮箱验证令牌仓储实例
func NewEmailVerificationTokenRepository(db *gorm.DB) *EmailVerificationTokenRepositoryImpl {
	return &EmailVerificationTokenRepositoryImpl{db: db}
}

// Create 创建邮箱验证令牌
func (r *EmailVerificationTokenRepositoryImpl) Create(ctx context.Context, token *authentication.EmailVerificationToken) error {
	po := &EmailVerificationTokenPO{
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
func (r *EmailVerificationTokenRepositoryImpl) FindByToken(ctx context.Context, token string) (*authentication.EmailVerificationToken, error) {
	var po EmailVerificationTokenPO
	if err := r.db.WithContext(ctx).
		Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).
		First(&po).Error; err != nil {
		return nil, err
	}
	return toDomainEmailToken(&po), nil
}

// FindByUserID 查找用户未使用的验证令牌
func (r *EmailVerificationTokenRepositoryImpl) FindByUserID(ctx context.Context, userID string) (*authentication.EmailVerificationToken, error) {
	var po EmailVerificationTokenPO
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND used = ? AND expires_at > ?", userID, false, time.Now()).
		Order("created_at DESC").
		First(&po).Error; err != nil {
		return nil, err
	}
	return toDomainEmailToken(&po), nil
}

// MarkAsUsed 标记令牌已使用
func (r *EmailVerificationTokenRepositoryImpl) MarkAsUsed(ctx context.Context, tokenID string) error {
	return r.db.WithContext(ctx).
		Model(&EmailVerificationTokenPO{}).
		Where("id = ?", tokenID).
		Update("used", true).Error
}

// DeleteExpired 删除所有过期令牌
func (r *EmailVerificationTokenRepositoryImpl) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&EmailVerificationTokenPO{}).Error
}

// toDomainEmailToken 持久化对象转领域对象
func toDomainEmailToken(po *EmailVerificationTokenPO) *authentication.EmailVerificationToken {
	return &authentication.EmailVerificationToken{
		ID:        po.ID,
		UserID:    po.UserID,
		Token:     po.Token,
		ExpiresAt: po.ExpiresAt,
		Used:      po.Used,
		CreatedAt: po.CreatedAt,
	}
}
