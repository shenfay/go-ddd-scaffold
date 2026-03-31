package aggregate

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/common"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// User 用户聚合根
// 使用组合模式替代继承，提高灵活性和可测试性
// 移除了登录统计相关字段（lastLoginAt, loginCount, lockedUntil, failedAttempts）
// 这些字段已迁移到独立的 LoginStats 聚合根，解决高频更新导致的乐观锁冲突
// 个人资料字段已封装到 UserProfile 值对象
type User struct {
	meta      *common.EntityMeta // 组合元数据
	username  *vo.UserName
	email     *vo.Email
	password  *vo.HashedPassword
	status    vo.UserStatus
	profile   *vo.UserProfile
	createdAt time.Time
	updatedAt time.Time
}

// NewUser 使用已哈希的密码创建新用户
func NewUser(username, email, hashedPassword string, idGenerator func() int64) (*User, error) {
	// 创建默认个人资料
	prof, err := vo.NewUserProfile("", "", "", vo.UserGenderUnknown, "", "")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &User{
		meta:      common.NewEntityMeta(nil, now),
		status:    vo.UserStatusActive,
		profile:   prof,
		createdAt: now,
		updatedAt: now,
	}

	// 使用 ID 生成器生成唯一 ID
	newUserID := idGenerator()
	userIDVO := vo.NewUserID(newUserID)
	user.meta.SetID(userIDVO)

	// 验证和设置用户名
	un, err := vo.NewUserName(username)
	if err != nil {
		return nil, err
	}
	user.username = un

	// 验证和设置邮箱
	em, err := vo.NewEmail(email)
	if err != nil {
		return nil, err
	}
	user.email = em

	// 设置已哈希的密码
	user.password = vo.NewHashedPassword(hashedPassword)

	// 产生 UserRegistered 领域事件
	registeredEvent := event.NewUserRegisteredEvent(
		userIDVO,
		username,
		email,
		user.status.String(),
		user.profile.DisplayName(),
		"", // registrationIP - 暂时为空，可通过后续重构传入
		0,  // tenantID - 暂时为 0，可通过后续重构传入
	)
	user.meta.ApplyEvent(registeredEvent)

	return user, nil
}

// ID 返回聚合根 ID
func (u *User) ID() interface{} {
	return u.meta.ID()
}

// Version 返回当前版本号
func (u *User) Version() int {
	return u.meta.Version()
}

// IncrementVersion 增加版本号
func (u *User) IncrementVersion() {
	u.meta.IncrementVersion()
	u.updatedAt = time.Now()
}

// ApplyEvent 应用领域事件
func (u *User) ApplyEvent(event common.DomainEvent) {
	u.meta.ApplyEvent(event)
}

// GetUncommittedEvents 获取未提交的事件
func (u *User) GetUncommittedEvents() []common.DomainEvent {
	return u.meta.GetUncommittedEvents()
}

// ClearUncommittedEvents 清除已提交的事件
func (u *User) ClearUncommittedEvents() {
	u.meta.ClearUncommittedEvents()
}

// Username 获取用户名
func (u *User) Username() *vo.UserName {
	return u.username
}

// Email 获取邮箱
func (u *User) Email() *vo.Email {
	return u.email
}

// Password 获取密码
func (u *User) Password() *vo.HashedPassword {
	return u.password
}

// Status 获取用户状态
func (u *User) Status() vo.UserStatus {
	return u.status
}

// Profile 获取个人资料
func (u *User) Profile() *vo.UserProfile {
	return u.profile
}

// DisplayName 获取显示名称
func (u *User) DisplayName() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.DisplayName()
}

// FirstName 获取名字
func (u *User) FirstName() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.FirstName()
}

// LastName 获取姓氏
func (u *User) LastName() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.LastName()
}

// Gender 获取性别
func (u *User) Gender() vo.UserGender {
	if u.profile == nil {
		return vo.UserGenderUnknown
	}
	return u.profile.Gender()
}

// PhoneNumber 获取电话号码
func (u *User) PhoneNumber() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.PhoneNumber()
}

// AvatarURL 获取头像 URL
func (u *User) AvatarURL() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.AvatarURL()
}

// UpdateProfile 更新个人资料
func (u *User) UpdateProfile(profile *vo.UserProfile) {
	u.profile = profile
	u.updatedAt = time.Now()
}

// SetDisplayName 设置显示名称
func (u *User) SetDisplayName(displayName string) error {
	newProfile, err := u.profile.UpdateDisplayName(displayName)
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// SetFirstName 设置名字
func (u *User) SetFirstName(firstName string) error {
	newProfile, err := u.profile.UpdateName(firstName, u.profile.LastName())
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// SetLastName 设置姓氏
func (u *User) SetLastName(lastName string) error {
	newProfile, err := u.profile.UpdateName(u.profile.FirstName(), lastName)
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// SetGender 设置性别
func (u *User) SetGender(gender vo.UserGender) error {
	newProfile, err := u.profile.UpdateGender(gender)
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// SetPhoneNumber 设置电话号码
func (u *User) SetPhoneNumber(phoneNumber string) error {
	newProfile, err := u.profile.UpdatePhoneNumber(phoneNumber)
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// SetAvatarURL 设置头像 URL
func (u *User) SetAvatarURL(avatarURL string) error {
	newProfile, err := u.profile.UpdateAvatarURL(avatarURL)
	if err != nil {
		return err
	}
	u.profile = newProfile
	u.updatedAt = time.Now()
	return nil
}

// Activate 激活用户
func (u *User) Activate() error {
	if u.status != vo.UserStatusPending {
		return common.NewBusinessError(CodeUserNotPending, "user is not in pending status")
	}

	u.status = vo.UserStatusActive
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// Deactivate 禁用用户
func (u *User) Deactivate(reason string) error {
	if u.status == vo.UserStatusInactive {
		return common.NewBusinessError(CodeUserAlreadyInactive, "user is already inactive")
	}

	u.status = vo.UserStatusInactive
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// Lock 锁定用户状态
func (u *User) Lock() error {
	if u.status == vo.UserStatusLocked {
		return common.NewBusinessError(CodeUserAlreadyLocked, "user is already locked")
	}

	u.status = vo.UserStatusLocked
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// Unlock 解锁用户状态
func (u *User) Unlock() error {
	if u.status != vo.UserStatusLocked {
		return common.NewBusinessError(CodeUserNotLocked, "user is not locked")
	}

	u.status = vo.UserStatusActive
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// ChangePassword 修改密码
func (u *User) ChangePassword(newPassword string, ipAddress string) error {
	// TODO: 这里应该验证新密码强度并加密
	u.password = vo.NewHashedPassword(newPassword)
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// UpdateEmail 更新邮箱
func (u *User) UpdateEmail(newEmail string) error {
	email, err := vo.NewEmail(newEmail)
	if err != nil {
		return err
	}

	u.email = email
	u.meta.SetUpdatedAt(time.Now())
	u.meta.IncrementVersion()

	return nil
}

// IsLocked 检查用户状态是否为锁定
func (u *User) IsLocked() bool {
	return u.status == vo.UserStatusLocked
}

// CanLogin 检查用户状态是否允许登录
// 注意：实际登录检查需要结合 LoginStats 的锁定状态
func (u *User) CanLogin() bool {
	return u.status == vo.UserStatusActive
}

// FullName 获取完整姓名
func (u *User) FullName() string {
	if u.profile == nil {
		return ""
	}
	return u.profile.FullName()
}

// GetFullName 获取完整姓名（带默认值）
func (u *User) GetFullName(defaultName string) string {
	fullName := u.FullName()
	if fullName == "" {
		return defaultName
	}
	return fullName
}

// CreatedAt 获取创建时间
func (u *User) CreatedAt() time.Time {
	return u.meta.CreatedAt()
}

// UpdatedAt 获取更新时间
func (u *User) UpdatedAt() time.Time {
	return u.meta.UpdatedAt()
}

// SetCreatedAt 设置创建时间 (用于 Builder)
func (u *User) SetCreatedAt(t time.Time) {
	u.meta.SetCreatedAt(t)
}

// SetUpdatedAt 设置更新时间 (用于 Builder)
func (u *User) SetUpdatedAt(t time.Time) {
	u.meta.SetUpdatedAt(t)
}
