package vo

import (
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
)

// ============================================================================
// UserProfile - 用户个人资料值对象
// ============================================================================

// UserProfile 用户个人资料值对象
type UserProfile struct {
	displayName string
	firstName   string
	lastName    string
	gender      UserGender
	phoneNumber string
	avatarURL   string
}

// NewUserProfile 创建用户个人资料
func NewUserProfile(displayName, firstName, lastName string, gender UserGender, phoneNumber, avatarURL string) (*UserProfile, error) {
	profile := &UserProfile{
		displayName: strings.TrimSpace(displayName),
		firstName:   strings.TrimSpace(firstName),
		lastName:    strings.TrimSpace(lastName),
		gender:      gender,
		phoneNumber: strings.TrimSpace(phoneNumber),
		avatarURL:   strings.TrimSpace(avatarURL),
	}

	if err := profile.Validate(); err != nil {
		return nil, err
	}

	return profile, nil
}

// Validate 验证个人资料
func (p *UserProfile) Validate() error {
	// 显示名称长度限制
	if len(p.displayName) > 100 {
		return &kernel.ValidationError{
			Field:   "display_name",
			Message: "display name cannot exceed 100 characters",
		}
	}

	// 名字长度限制
	if len(p.firstName) > 50 {
		return &kernel.ValidationError{
			Field:   "first_name",
			Message: "first name cannot exceed 50 characters",
		}
	}

	// 姓氏长度限制
	if len(p.lastName) > 50 {
		return &kernel.ValidationError{
			Field:   "last_name",
			Message: "last name cannot exceed 50 characters",
		}
	}

	// 手机号格式验证（简单验证）
	if p.phoneNumber != "" && len(p.phoneNumber) > 20 {
		return &kernel.ValidationError{
			Field:   "phone_number",
			Message: "phone number format is invalid",
		}
	}

	// 头像URL长度限制
	if len(p.avatarURL) > 500 {
		return &kernel.ValidationError{
			Field:   "avatar_url",
			Message: "avatar URL cannot exceed 500 characters",
		}
	}

	return nil
}

// DisplayName 返回显示名称
func (p *UserProfile) DisplayName() string {
	return p.displayName
}

// FirstName 返回名字
func (p *UserProfile) FirstName() string {
	return p.firstName
}

// LastName 返回姓氏
func (p *UserProfile) LastName() string {
	return p.lastName
}

// FullName 返回全名
func (p *UserProfile) FullName() string {
	if p.firstName == "" && p.lastName == "" {
		return p.displayName
	}
	return strings.TrimSpace(p.firstName + " " + p.lastName)
}

// Gender 返回性别
func (p *UserProfile) Gender() UserGender {
	return p.gender
}

// PhoneNumber 返回手机号
func (p *UserProfile) PhoneNumber() string {
	return p.phoneNumber
}

// AvatarURL 返回头像URL
func (p *UserProfile) AvatarURL() string {
	return p.avatarURL
}

// Equals 比较两个个人资料是否相等
func (p *UserProfile) Equals(other *UserProfile) bool {
	if other == nil {
		return false
	}
	return p.displayName == other.displayName &&
		p.firstName == other.firstName &&
		p.lastName == other.lastName &&
		p.gender == other.gender &&
		p.phoneNumber == other.phoneNumber &&
		p.avatarURL == other.avatarURL
}

// UpdateDisplayName 更新显示名称
func (p *UserProfile) UpdateDisplayName(displayName string) (*UserProfile, error) {
	return NewUserProfile(displayName, p.firstName, p.lastName, p.gender, p.phoneNumber, p.avatarURL)
}

// UpdateName 更新姓名
func (p *UserProfile) UpdateName(firstName, lastName string) (*UserProfile, error) {
	return NewUserProfile(p.displayName, firstName, lastName, p.gender, p.phoneNumber, p.avatarURL)
}

// UpdateGender 更新性别
func (p *UserProfile) UpdateGender(gender UserGender) (*UserProfile, error) {
	return NewUserProfile(p.displayName, p.firstName, p.lastName, gender, p.phoneNumber, p.avatarURL)
}

// UpdatePhoneNumber 更新手机号
func (p *UserProfile) UpdatePhoneNumber(phoneNumber string) (*UserProfile, error) {
	return NewUserProfile(p.displayName, p.firstName, p.lastName, p.gender, phoneNumber, p.avatarURL)
}

// UpdateAvatarURL 更新头像URL
func (p *UserProfile) UpdateAvatarURL(avatarURL string) (*UserProfile, error) {
	return NewUserProfile(p.displayName, p.firstName, p.lastName, p.gender, p.phoneNumber, avatarURL)
}
