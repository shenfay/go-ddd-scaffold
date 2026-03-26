package builder

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserBuilder 用户构建器（用于从数据库重建聚合根）
// 职责：将数据模型转换为领域对象，重建聚合根状态
type UserBuilder struct {
	id           int64
	username     string
	email        string
	passwordHash string
	status       vo.UserStatus
	profile      *vo.UserProfile
	version      int
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUserBuilder 创建用户构建器
func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		status: vo.UserStatusActive,
	}
}

// WithID 设置用户 ID
func (b *UserBuilder) WithID(id int64) *UserBuilder {
	b.id = id
	return b
}

// WithUsername 设置用户名
func (b *UserBuilder) WithUsername(username string) *UserBuilder {
	b.username = username
	return b
}

// WithEmail 设置邮箱
func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.email = email
	return b
}

// WithPasswordHash 设置密码哈希值
func (b *UserBuilder) WithPasswordHash(hash string) *UserBuilder {
	b.passwordHash = hash
	return b
}

// WithStatus 设置用户状态
func (b *UserBuilder) WithStatus(status vo.UserStatus) *UserBuilder {
	b.status = status
	return b
}

// WithProfile 设置个人资料
func (b *UserBuilder) WithProfile(prof *vo.UserProfile) *UserBuilder {
	b.profile = prof
	return b
}

// WithDisplayName 设置显示名称
func (b *UserBuilder) WithDisplayName(name string) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile(name, "", "", vo.UserGenderUnknown, "", "")
	} else {
		b.profile, _ = b.profile.UpdateDisplayName(name)
	}
	return b
}

// WithFirstName 设置名字
func (b *UserBuilder) WithFirstName(name string) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile("", name, "", vo.UserGenderUnknown, "", "")
	} else {
		b.profile, _ = b.profile.UpdateName(name, b.profile.LastName())
	}
	return b
}

// WithLastName 设置姓氏
func (b *UserBuilder) WithLastName(name string) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile("", "", name, vo.UserGenderUnknown, "", "")
	} else {
		b.profile, _ = b.profile.UpdateName(b.profile.FirstName(), name)
	}
	return b
}

// WithGender 设置性别
func (b *UserBuilder) WithGender(gender vo.UserGender) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile("", "", "", gender, "", "")
	} else {
		b.profile, _ = b.profile.UpdateGender(gender)
	}
	return b
}

// WithPhoneNumber 设置电话号码
func (b *UserBuilder) WithPhoneNumber(phone string) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile("", "", "", vo.UserGenderUnknown, phone, "")
	} else {
		b.profile, _ = b.profile.UpdatePhoneNumber(phone)
	}
	return b
}

// WithAvatarURL 设置头像 URL
func (b *UserBuilder) WithAvatarURL(url string) *UserBuilder {
	if b.profile == nil {
		b.profile, _ = vo.NewUserProfile("", "", "", vo.UserGenderUnknown, "", url)
	} else {
		b.profile, _ = b.profile.UpdateAvatarURL(url)
	}
	return b
}

// WithVersion 设置版本号（用于乐观锁）
func (b *UserBuilder) WithVersion(version int) *UserBuilder {
	b.version = version
	return b
}

// WithTimestamps 设置时间戳
func (b *UserBuilder) WithTimestamps(createdAt, updatedAt time.Time) *UserBuilder {
	b.createdAt = createdAt
	b.updatedAt = updatedAt
	return b
}

// Build 构建最终的用户对象
func (b *UserBuilder) Build() (*aggregate.User, error) {
	// 创建用户（使用已验证的数据）
	user, err := aggregate.NewUser(b.username, b.email, b.passwordHash, func() int64 { return b.id })
	if err != nil {
		return nil, err
	}

	// 应用可选字段
	if b.profile != nil {
		user.UpdateProfile(b.profile)
	}

	// 应用版本号
	for i := 0; i < b.version; i++ {
		user.IncrementVersion()
	}

	// 应用时间戳
	user.SetCreatedAt(b.createdAt)
	user.SetUpdatedAt(b.updatedAt)

	return user, nil
}
