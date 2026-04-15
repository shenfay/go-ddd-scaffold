package authentication

import "context"

// PasswordResetTokenRepository 密码重置令牌仓储接口
//
// 定义了密码重置令牌的持久化操作规范
// 由基础设施层实现,遵循依赖倒置原则
type PasswordResetTokenRepository interface {
	// Create 创建密码重置令牌
	Create(ctx context.Context, token *PasswordResetToken) error

	// FindByToken 根据令牌字符串查找
	// 自动过滤已使用和已过期的令牌
	FindByToken(ctx context.Context, token string) (*PasswordResetToken, error)

	// MarkAsUsed 标记令牌已使用
	MarkAsUsed(ctx context.Context, tokenID string) error

	// DeleteExpired 删除所有过期令牌(清理任务)
	DeleteExpired(ctx context.Context) error
}

// EmailVerificationTokenRepository 邮箱验证令牌仓储接口
//
// 定义了邮箱验证令牌的持久化操作规范
// 由基础设施层实现,遵循依赖倒置原则
type EmailVerificationTokenRepository interface {
	// Create 创建邮箱验证令牌
	Create(ctx context.Context, token *EmailVerificationToken) error

	// FindByToken 根据令牌字符串查找
	// 自动过滤已使用和已过期的令牌
	FindByToken(ctx context.Context, token string) (*EmailVerificationToken, error)

	// FindByUserID 查找用户未使用的验证令牌
	FindByUserID(ctx context.Context, userID string) (*EmailVerificationToken, error)

	// MarkAsUsed 标记令牌已使用
	MarkAsUsed(ctx context.Context, tokenID string) error

	// DeleteExpired 删除所有过期令牌(清理任务)
	DeleteExpired(ctx context.Context) error
}
