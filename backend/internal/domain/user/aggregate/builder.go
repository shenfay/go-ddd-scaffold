package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
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
	b.user.SetID(vo.NewUserID(id))
	return b
}

// WithUsername 设置用户名（直接赋值，不验证，因为数据来自数据库）
func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	// 从数据库加载的数据已经验证过，可以直接构造
	un := &vo.UserName{}
	if validated, err := vo.NewUserName(username); err == nil {
		*un = *validated
	}
	b.user.username = un
	return b
}

// WithEmail 设置邮箱（直接赋值，不验证）
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	em := &vo.Email{}
	if validated, err := vo.NewEmail(email); err == nil {
		*em = *validated
	}
	b.user.email = em
	return b
}

// WithPasswordHash 设置密码哈希值
func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.user.password = vo.NewHashedPassword(hash)
	return b
}

// WithStatus 设置用户状态
func (b *UserBuilder) WithStatus(status vo.UserStatus) *UserBuilder {
	b.user.status = status
	return b
}

// WithGender 设置性别
func (b *UserBuilder) WithGender(gender vo.UserGender) *UserBuilder {
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
