package user

import (
	"strconv"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/vo"
)

// Mapper DTO 转换器
type Mapper struct{}

// NewMapper 创建转换器
func NewMapper() *Mapper {
	return &Mapper{}
}

// parseUserID 解析用户ID
func (m *Mapper) parseUserID(id string) (vo.UserID, error) {
	value, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return vo.UserID{}, err
	}
	return vo.NewUserID(value), nil
}

// parseGender 解析性别
func (m *Mapper) parseGender(gender string) vo.UserGender {
	switch gender {
	case "male":
		return vo.UserGenderMale
	case "female":
		return vo.UserGenderFemale
	case "other":
		return vo.UserGenderOther
	default:
		return vo.UserGenderUnknown
	}
}

// parseStatus 解析状态
func (m *Mapper) parseStatus(status string) vo.UserStatus {
	switch status {
	case "active":
		return vo.UserStatusActive
	case "inactive":
		return vo.UserStatusInactive
	case "pending":
		return vo.UserStatusPending
	case "locked":
		return vo.UserStatusLocked
	default:
		return vo.UserStatusPending
	}
}

// ToUpdateUserCommand 转换为更新用户命令
func (m *Mapper) ToUpdateUserCommand(req *UpdateUserRequest, userID string) (*userApp.UpdateUserProfileCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	cmd := &userApp.UpdateUserProfileCommand{
		UserID:      uid,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
	}

	if req.Gender != nil {
		gender := m.parseGender(*req.Gender)
		cmd.Gender = &gender
	}

	return cmd, nil
}

// ToChangePasswordCommand 转换为修改密码命令
func (m *Mapper) ToChangePasswordCommand(req *ChangePasswordRequest, userID string) (*userApp.ChangePasswordCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	return &userApp.ChangePasswordCommand{
		UserID:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}, nil
}

// ParseUserID 解析用户 ID（公开方法）
func (m *Mapper) ParseUserID(id string) (vo.UserID, error) {
	return m.parseUserID(id)
}

// ToRegisterUserCommand 转换为注册用户命令
func (m *Mapper) ToRegisterUserCommand(req *RegisterUserRequest) *userApp.RegisterUserCommand {
	return &userApp.RegisterUserCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

// ToAuthenticateUserCommand 转换为认证用户命令
func (m *Mapper) ToAuthenticateUserCommand(req *AuthenticateUserRequest, ipAddress, userAgent string) *userApp.AuthenticateUserCommand {
	return &userApp.AuthenticateUserCommand{
		Username:  req.Username,
		Password:  req.Password,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}
}
