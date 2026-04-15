package authentication

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// PasswordResetToken 密码重置令牌聚合根
//
// 职责:
// - 封装密码重置令牌的创建和验证逻辑
// - 确保令牌的安全性和有效期管理
// - 维护令牌的使用状态
type PasswordResetToken struct {
	ID        string    // 令牌唯一标识(ULID)
	UserID    string    // 关联用户ID
	Token     string    // 加密安全的随机令牌
	ExpiresAt time.Time // 过期时间
	Used      bool      // 是否已使用
	CreatedAt time.Time // 创建时间
}

// NewPasswordResetToken 创建密码重置令牌
//
// 业务流程:
// 1. 生成加密安全的随机令牌(32字节)
// 2. 设置1小时有效期
// 3. 初始化令牌状态
//
// 参数:
//   - userID: 请求密码重置的用户ID
//
// 返回:
//   - *PasswordResetToken: 创建的令牌实例
//   - error: 生成失败时返回错误
func NewPasswordResetToken(userID string) (*PasswordResetToken, error) {
	token, err := GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	return &PasswordResetToken{
		ID:        GenerateID(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
		CreatedAt: time.Now(),
	}, nil
}

// IsExpired 检查令牌是否过期
//
// 返回:
//   - bool: 已过期返回 true,否则返回 false
func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// MarkAsUsed 标记令牌已使用
//
// 注意: 此操作不可逆,用于防止令牌重放攻击
func (t *PasswordResetToken) MarkAsUsed() {
	t.Used = true
}

// GenerateSecureToken 生成加密安全的随机令牌(包级别共享函数)
//
// 使用 crypto/rand 生成 32 字节的随机数,并转换为十六进制字符串
// 返回 64 字符的令牌,具有足够的熵值防止暴力破解
//
// 返回:
//   - string: 十六进制编码的令牌
//   - error: 随机数生成失败时返回错误
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateID 生成 ULID 标识(包级别共享函数)
//
// TODO: 使用 pkg/utils 中的 ULID 生成函数
// 临时实现:使用时间戳+随机数
func GenerateID() string {
	return time.Now().Format("20060102150405.000000000")
}
