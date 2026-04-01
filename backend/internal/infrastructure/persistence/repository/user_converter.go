package repository

import (
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/builder"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
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
func (c *UserConverter) ToDomain(m *model.User) *user.User {
	if m == nil {
		return nil
	}

	builder := builder.NewUserBuilder()

	// 必填字段
	builder.WithID(m.ID).
		WithUsername(m.Username).
		WithEmail(m.Email).
		WithPasswordHash(m.PasswordHash).
		WithStatus(vo.UserStatus(m.Status))

	// 可选字段 - 使用指针，有值才设置
	c.applyOptionalFields(builder, m)

	// 时间戳
	c.applyTimestamps(builder, m)

	user, err := builder.Build()
	if err != nil {
		// 数据库数据应该是有效的，如果出错说明数据不一致
		panic("failed to build user from database: " + err.Error())
	}

	return user
}

// applyOptionalFields 应用可选字段
func (c *UserConverter) applyOptionalFields(builder *builder.UserBuilder, m *model.User) {
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
}

// applyTimestamps 应用时间戳字段
func (c *UserConverter) applyTimestamps(builder *builder.UserBuilder, m *model.User) {
	if m.CreatedAt != nil && m.UpdatedAt != nil {
		builder.WithTimestamps(*m.CreatedAt, *m.UpdatedAt)
	}
}

// FromDomain 将领域对象转换为数据模型
func (c *UserConverter) FromDomain(u *user.User) *model.User {
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
