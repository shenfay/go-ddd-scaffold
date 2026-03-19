package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserBuilder 用户构建器（用于从数据库重建聚合根）
type UserBuilder struct {
	user *User
}

// NewUserBuilder 创建用户构建器
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &User{},
	}
}

// WithID 设置用户 ID
func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.user.SetID(valueobject.NewUserID(id))
	return b
}

// WithUsername 设置用户名（直接赋值，不验证，因为数据来自数据库）
func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	// 从数据库加载的数据已经验证过，可以直接构造
	un := &valueobject.UserName{}
	setUserNameValue(un, username)
	b.user.username = un
	return b
}

// setUserNameValue 内部辅助函数，用于设置用户名的值
func setUserNameValue(un *valueobject.UserName, value string) {
	// 使用反射或 unsafe 包来设置私有字段会比较复杂
	// 这里我们采用一个简单的方法：重新验证（性能损失可接受）
	if validated, err := valueobject.NewUserName(value); err == nil {
		*un = *validated
	}
}

// WithEmail 设置邮箱（直接赋值，不验证）
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	em := &valueobject.Email{}
	if validated, err := valueobject.NewEmail(email); err == nil {
		*em = *validated
	}
	b.user.email = em
	return b
}

// WithPasswordHash 设置密码哈希值
func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.user.password = valueobject.NewHashedPassword(hash)
	return b
}

// WithStatus 设置用户状态
func (b *UserBuilder) WithStatus(status valueobject.UserStatus) *UserBuilder {
	b.user.status = status
	return b
}

// WithGender 设置性别
func (b *UserBuilder) WithGender(gender valueobject.UserGender) *UserBuilder {
	b.user.gender = gender
	return b
}

// WithDisplayName 设置显示名称
func (b *UserBuilder) WithDisplayName(name string) *UserBuilder {
	b.user.displayName = name
	return b
}

// WithFirstName 设置名字
func (b *UserBuilder) WithFirstName(name string) *UserBuilder {
	b.user.firstName = name
	return b
}

// WithLastName 设置姓氏
func (b *UserBuilder) WithLastName(name string) *UserBuilder {
	b.user.lastName = name
	return b
}

// WithPhoneNumber 设置电话号码
func (b *UserBuilder) WithPhoneNumber(phone string) *UserBuilder {
	b.user.phoneNumber = phone
	return b
}

// WithAvatarURL 设置头像 URL
func (b *UserBuilder) WithAvatarURL(url string) *UserBuilder {
	b.user.avatarURL = url
	return b
}

// WithLastLoginAt 设置最后登录时间
func (b *UserBuilder) WithLastLoginAt(t *time.Time) *UserBuilder {
	b.user.lastLoginAt = t
	return b
}

// WithLoginCount 设置登录次数
func (b *UserBuilder) WithLoginCount(count int) *UserBuilder {
	b.user.loginCount = count
	return b
}

// WithFailedAttempts 设置失败尝试次数
func (b *UserBuilder) WithFailedAttempts(count int) *UserBuilder {
	b.user.failedAttempts = count
	return b
}

// WithLockedUntil 设置锁定截止时间
func (b *UserBuilder) WithLockedUntil(t *time.Time) *UserBuilder {
	b.user.lockedUntil = t
	return b
}

// WithVersion 设置版本号（用于乐观锁）
func (b *UserBuilder) WithVersion(version int) *UserBuilder {
	b.user.SetVersion(version)
	return b
}

// WithTimestamps 设置时间戳
func (b *UserBuilder) WithTimestamps(createdAt, updatedAt time.Time) *UserBuilder {
	b.user.createdAt = createdAt
	b.user.updatedAt = updatedAt
	return b
}

// Build 构建最终的用户对象
func (b *UserBuilder) Build() *User {
	return b.user
}
