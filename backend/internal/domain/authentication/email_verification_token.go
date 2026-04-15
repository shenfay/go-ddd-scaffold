package authentication

import (
	"time"
)

// EmailVerificationToken 邮箱验证令牌聚合根
//
// 职责:
// - 封装邮箱验证令牌的创建和验证逻辑
// - 确保令牌的安全性和有效期管理
// - 维护令牌的使用状态
type EmailVerificationToken struct {
	ID        string    // 令牌唯一标识(ULID)
	UserID    string    // 关联用户ID
	Token     string    // 加密安全的随机令牌
	ExpiresAt time.Time // 过期时间
	Used      bool      // 是否已使用
	CreatedAt time.Time // 创建时间
}

// NewEmailVerificationToken 创建邮箱验证令牌
//
// 业务流程:
// 1. 生成加密安全的随机令牌(32字节)
// 2. 设置24小时有效期
// 3. 初始化令牌状态
//
// 参数:
//   - userID: 需要验证邮箱的用户ID
//
// 返回:
//   - *EmailVerificationToken: 创建的令牌实例
//   - error: 生成失败时返回错误
func NewEmailVerificationToken(userID string) (*EmailVerificationToken, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	return &EmailVerificationToken{
		ID:        GenerateID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}, nil
}

// IsExpired 检查令牌是否过期
//
// 返回:
//   - bool: 已过期返回 true,否则返回 false
func (t *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// MarkAsUsed 标记令牌已使用
//
// 注意: 此操作不可逆,用于防止令牌重放攻击
func (t *EmailVerificationToken) MarkAsUsed() {
	t.Used = true
}
