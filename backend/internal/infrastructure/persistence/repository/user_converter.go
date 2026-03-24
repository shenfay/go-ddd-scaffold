package repository

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// UserConverter 用户领域对象与数据模型转换器
type UserConverter struct{}

// NewUserConverter 创建用户转换器
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToDomain 将数据模型转换为领域对象
func (c *UserConverter) ToDomain(m *model.User) *aggregate.User {
	if m == nil {
		return nil
	}

	// 使用 Builder 模式构建领域对象
	builder := aggregate.NewUserBuilder()

	// 设置基本字段
	builder.WithID(m.ID)

	// 设置用户名和邮箱（直接使用字符串，Builder 内部处理）
	if m.Username != "" {
		builder.WithUsername(m.Username)
	}
	if m.Email != "" {
		builder.WithEmail(m.Email)
	}

	// 设置密码
	if m.PasswordHash != "" {
		builder.WithPasswordHash(m.PasswordHash)
	}

	// 设置状态
	builder.WithStatus(vo.UserStatus(m.Status))

	// 设置可选字段
	if m.DisplayName != nil {
		builder.WithDisplayName(*m.DisplayName)
	}
	if m.Gender != nil {
		builder.WithGender(vo.UserGender(*m.Gender))
	}
	if m.PhoneNumber != nil {
		builder.WithPhoneNumber(*m.PhoneNumber)
	}
	if m.AvatarURL != nil {
		builder.WithAvatarURL(*m.AvatarURL)
	}
	if m.Version != nil {
		builder.WithVersion(int(*m.Version))
	}

	// 设置时间戳
	if m.CreatedAt != nil && m.UpdatedAt != nil {
		builder.WithTimestamps(*m.CreatedAt, *m.UpdatedAt)
	}

	return builder.Build()
}

// FromDomain 将领域对象转换为数据模型
func (c *UserConverter) FromDomain(u *aggregate.User) *model.User {
	if u == nil {
		return nil
	}

	displayName := u.DisplayName()
	phoneNumber := u.PhoneNumber()
	avatarURL := u.AvatarURL()
	version := int(u.Version())

	return &model.User{
		ID:           u.ID().(vo.UserID).Int64(),
		Username:     u.Username().Value(),
		Email:        u.Email().Value(),
		PasswordHash: u.Password().Value(),
		Status:       int16(u.Status()),
		DisplayName:  util.StringPtrNilIfEmpty(displayName),
		Gender:       util.Int16PtrNilIfZero(int16(u.Gender())),
		PhoneNumber:  util.StringPtrNilIfEmpty(phoneNumber),
		AvatarURL:    util.StringPtrNilIfEmpty(avatarURL),
		Version:      util.Int32PtrNilIfZero(int32(version)),
		CreatedAt:    util.Time(u.CreatedAt()),
		UpdatedAt:    util.Time(u.UpdatedAt()),
	}
}
