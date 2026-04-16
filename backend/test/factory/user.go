package factory

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// UserFactory 用户测试数据工厂
type UserFactory struct{}

// NewUserFactory 创建用户工厂实例
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// CreateUser 创建测试用户(默认已验证)
func (f *UserFactory) CreateUser(opts ...UserOption) *user.User {
	// 默认值
	email := "test@example.com"
	password := "TestPassword123!"
	verified := true
	locked := false
	failedAttempts := 0

	// 应用选项
	for _, opt := range opts {
		opt(&email, &password, &verified, &locked, &failedAttempts)
	}

	u, _ := user.NewUser(email, password)
	u.ID = "user-test-" + email[:5]

	if verified {
		u.VerifyEmail()
	}

	if locked {
		for i := 0; i < failedAttempts; i++ {
			u.IncrementFailedAttempts(failedAttempts)
		}
	}

	return u
}

// CreateUsers 批量创建测试用户
func (f *UserFactory) CreateUsers(count int) []*user.User {
	users := make([]*user.User, count)
	for i := 0; i < count; i++ {
		users[i] = f.CreateUser(WithEmailPrefix("user"))
	}
	return users
}

// UserOption 用户选项函数
type UserOption func(email, password *string, verified, locked *bool, failedAttempts *int)

// WithEmail 设置邮箱
func WithEmail(email string) UserOption {
	return func(e, _ *string, _, _ *bool, _ *int) {
		*e = email
	}
}

// WithEmailPrefix 使用邮箱前缀(自动生成完整邮箱)
func WithEmailPrefix(prefix string) UserOption {
	return func(e, _ *string, _, _ *bool, _ *int) {
		*e = prefix + "@example.com"
	}
}

// WithPassword 设置密码
func WithPassword(password string) UserOption {
	return func(_, p *string, _, _ *bool, _ *int) {
		*p = password
	}
}

// WithUnverified 设置未验证邮箱
func WithUnverified() UserOption {
	return func(_, _ *string, v *bool, _ *bool, _ *int) {
		*v = false
	}
}

// WithLocked 设置锁定账户
func WithLocked(failedAttempts int) UserOption {
	return func(_, _ *string, _ *bool, l *bool, fa *int) {
		*l = true
		*fa = failedAttempts
	}
}

// WithFailedAttempts 设置失败次数
func WithFailedAttempts(count int) UserOption {
	return func(_, _ *string, _ *bool, _ *bool, fa *int) {
		*fa = count
	}
}

// TokenPair 测试用 Token 对
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

// CreateTokenPair 创建测试 Token 对
func (f *UserFactory) CreateTokenPair(userID, email string) *TokenPair {
	return &TokenPair{
		AccessToken:  "access-token-" + userID,
		RefreshToken: "refresh-token-" + userID,
		ExpiresIn:    3600,
	}
}

// CreateExpiredTokenPair 创建已过期的 Token 对
func (f *UserFactory) CreateExpiredTokenPair(userID, email string) *TokenPair {
	return &TokenPair{
		AccessToken:  "expired-access-token-" + userID,
		RefreshToken: "expired-refresh-token-" + userID,
		ExpiresIn:    0,
	}
}

// DeviceInfo 测试用设备信息
type DeviceInfo struct {
	IP         string
	UserAgent  string
	DeviceType string
}

// CreateDeviceInfo 创建设备信息
func (f *UserFactory) CreateDeviceInfo(opts ...DeviceOption) *DeviceInfo {
	info := &DeviceInfo{
		IP:         "192.168.1.1",
		UserAgent:  "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		DeviceType: "web",
	}

	for _, opt := range opts {
		opt(info)
	}

	return info
}

// DeviceOption 设备选项函数
type DeviceOption func(*DeviceInfo)

// WithIP 设置 IP
func WithIP(ip string) DeviceOption {
	return func(d *DeviceInfo) {
		d.IP = ip
	}
}

// WithUserAgent 设置 UserAgent
func WithUserAgent(ua string) DeviceOption {
	return func(d *DeviceInfo) {
		d.UserAgent = ua
	}
}

// WithDeviceType 设置设备类型
func WithDeviceType(deviceType string) DeviceOption {
	return func(d *DeviceInfo) {
		d.DeviceType = deviceType
	}
}

// CreateTimestamp 创建测试时间戳
func (f *UserFactory) CreateTimestamp(offset time.Duration) time.Time {
	return time.Now().Add(offset)
}

// CreateFutureTimestamp 创建未来时间戳
func (f *UserFactory) CreateFutureTimestamp(hours int) time.Time {
	return time.Now().Add(time.Duration(hours) * time.Hour)
}

// CreatePastTimestamp 创建过去时间戳
func (f *UserFactory) CreatePastTimestamp(hours int) time.Time {
	return time.Now().Add(-time.Duration(hours) * time.Hour)
}
