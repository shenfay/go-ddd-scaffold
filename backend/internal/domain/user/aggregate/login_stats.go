package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// LoginStats 用户登录统计聚合根
// 从User聚合根中拆分出来，解决高频更新导致的乐观锁冲突问题
type LoginStats struct {
	kernel.BaseEntity

	userID         vo.UserID
	lastLoginAt    *time.Time
	loginCount     int
	failedAttempts int
	lockedUntil    *time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

// NewLoginStats 创建新的登录统计
func NewLoginStats(userID vo.UserID) *LoginStats {
	now := time.Now()
	return &LoginStats{
		userID:     userID,
		loginCount: 0,
		createdAt:  now,
		updatedAt:  now,
	}
}

// UserID 获取用户ID
func (ls *LoginStats) UserID() vo.UserID {
	return ls.userID
}

// LastLoginAt 获取最后登录时间
func (ls *LoginStats) LastLoginAt() *time.Time {
	return ls.lastLoginAt
}

// LoginCount 获取登录次数
func (ls *LoginStats) LoginCount() int {
	return ls.loginCount
}

// FailedAttempts 获取失败尝试次数
func (ls *LoginStats) FailedAttempts() int {
	return ls.failedAttempts
}

// LockedUntil 获取锁定截止时间
func (ls *LoginStats) LockedUntil() *time.Time {
	return ls.lockedUntil
}

// RecordLogin 记录成功登录
func (ls *LoginStats) RecordLogin() {
	now := time.Now()
	ls.lastLoginAt = &now
	ls.loginCount++
	ls.failedAttempts = 0
	ls.updatedAt = now
	ls.IncrementVersion()
}

// RecordFailedLogin 记录失败登录
func (ls *LoginStats) RecordFailedLogin() {
	ls.failedAttempts++
	ls.updatedAt = time.Now()
	ls.IncrementVersion()
}

// Lock 锁定账户
func (ls *LoginStats) Lock(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	ls.lockedUntil = &lockUntil
	ls.updatedAt = time.Now()
	ls.IncrementVersion()
}

// Unlock 解锁账户
func (ls *LoginStats) Unlock() {
	ls.lockedUntil = nil
	ls.failedAttempts = 0
	ls.updatedAt = time.Now()
	ls.IncrementVersion()
}

// IsLocked 检查是否被锁定
func (ls *LoginStats) IsLocked() bool {
	if ls.lockedUntil == nil {
		return false
	}
	if time.Now().After(*ls.lockedUntil) {
		// 锁定时间已过，自动解锁
		ls.Unlock()
		return false
	}
	return true
}

// ResetFailedAttempts 重置失败尝试次数
func (ls *LoginStats) ResetFailedAttempts() {
	ls.failedAttempts = 0
	ls.updatedAt = time.Now()
	ls.IncrementVersion()
}

// CanLogin 检查是否可以登录
func (ls *LoginStats) CanLogin() bool {
	return !ls.IsLocked()
}

// CreatedAt 获取创建时间
func (ls *LoginStats) CreatedAt() time.Time {
	return ls.createdAt
}

// UpdatedAt 获取更新时间
func (ls *LoginStats) UpdatedAt() time.Time {
	return ls.updatedAt
}
