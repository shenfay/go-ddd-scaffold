package service

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher 密码哈希接口
type PasswordHasher interface {
	// Hash 哈希密码
	Hash(password string) (string, error)
	// Verify 验证密码
	Verify(password, hash string) bool
}

// BcryptPasswordHasher 基于 bcrypt 的密码哈希实现
type BcryptPasswordHasher struct {
	cost int // bcrypt 成本因子
}

// NewBcryptPasswordHasher 创建 bcrypt 密码哈希器
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &BcryptPasswordHasher{cost: cost}
}

// Hash 哈希密码
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Verify 验证密码
func (h *BcryptPasswordHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
