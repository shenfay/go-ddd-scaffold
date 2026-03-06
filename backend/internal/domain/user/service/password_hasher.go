package service

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher 密码哈希器接口
type PasswordHasher interface {
	// Hash 对明文密码进行哈希
	Hash(plain string) (string, error)
	// Verify 验证密码是否匹配
	Verify(hash, plain string) bool
}

// BcryptPasswordHasher 基于 bcrypt 的密码哈希器实现
type BcryptPasswordHasher struct {
	cost int // 加密成本因子
}

// NewBcryptPasswordHasher 创建 bcrypt 密码哈希器
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	return &BcryptPasswordHasher{cost: cost}
}

// NewDefaultBcryptPasswordHasher 创建默认配置的 bcrypt 密码哈希器（cost=12）
// 用于 Wire 依赖注入
// 生产环境推荐 cost=12，开发环境可调整为 10 以提升性能
func NewDefaultBcryptPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{cost: 12} // 生产环境成本因子
}

// Hash 对明文密码进行哈希
func (h *BcryptPasswordHasher) Hash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), h.cost)
	return string(bytes), err
}

// Verify 验证密码是否匹配
func (h *BcryptPasswordHasher) Verify(hash, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
