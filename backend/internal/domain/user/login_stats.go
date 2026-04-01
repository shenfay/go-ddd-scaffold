package user

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
)

// LoginStats 用户登录统计聚合根
// 从 User 聚合根中拆分出来，解决高频更新导致的乐观锁冲突问题
// 使用组合模式替代继承，提高灵活性
type LoginStats struct {
	meta           *common.EntityMeta // 组合元数据
	userID         UserID
	lastLoginAt    *time.Time
	loginCount     int
	failedAttempts int
	lockedUntil    *time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

// NewLoginStats 创建新的登录统计
func NewLoginStats(userID UserID) *LoginStats {
	now := time.Now()
	return &LoginStats{
		meta:       common.NewEntityMeta(nil, now),
		userID:     userID,
		loginCount: 0,
		createdAt:  now,
		updatedAt:  now,
	}
}

// ID 返回聚合根 ID
func (ls *LoginStats) ID() interface{} {
	return ls.meta.ID()
}

// Version 返回当前版本号
func (ls *LoginStats) Version() int {
	return ls.meta.Version()
}

// IncrementVersion 增加版本号
func (ls *LoginStats) IncrementVersion() {
	ls.meta.IncrementVersion()
	ls.updatedAt = time.Now()
}

// ApplyEvent 应用领域事件
func (ls *LoginStats) ApplyEvent(event common.DomainEvent) {
	ls.meta.ApplyEvent(event)
}

// GetUncommittedEvents 获取未提交的事件
func (ls *LoginStats) GetUncommittedEvents() []common.DomainEvent {
	return ls.meta.GetUncommittedEvents()
}

// ClearUncommittedEvents 清除已提交的事件
func (ls *LoginStats) ClearUncommittedEvents() {
	ls.meta.ClearUncommittedEvents()
}

// UserID 获取用户 ID
func (ls *LoginStats) UserID() UserID {
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
	ls.meta.SetUpdatedAt(now)
	ls.meta.IncrementVersion()
}

// RecordFailedLogin 记录失败登录
func (ls *LoginStats) RecordFailedLogin() {
	ls.failedAttempts++
	ls.meta.SetUpdatedAt(time.Now())
	ls.meta.IncrementVersion()
}

// Lock 锁定账户
func (ls *LoginStats) Lock(duration time.Duration) {
	lockUntil := time.Now().Add(duration)
	ls.lockedUntil = &lockUntil
	ls.meta.SetUpdatedAt(time.Now())
	ls.meta.IncrementVersion()
}

// Unlock 解锁账户
func (ls *LoginStats) Unlock() {
	ls.lockedUntil = nil
	ls.failedAttempts = 0
	ls.meta.SetUpdatedAt(time.Now())
	ls.meta.IncrementVersion()
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
	ls.meta.SetUpdatedAt(time.Now())
	ls.meta.IncrementVersion()
}

// CanLogin 检查是否可以登录
func (ls *LoginStats) CanLogin() bool {
	return !ls.IsLocked()
}

// CreatedAt 获取创建时间
func (ls *LoginStats) CreatedAt() time.Time {
	return ls.meta.CreatedAt()
}

// UpdatedAt 获取更新时间
func (ls *LoginStats) UpdatedAt() time.Time {
	return ls.meta.UpdatedAt()
}
