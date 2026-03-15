package user

import (
	"strconv"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// Mapper DTO 转换器
type Mapper struct{}

// NewMapper 创建转换器
func NewMapper() *Mapper {
	return &Mapper{}
}

// parseUserID 解析用户ID
func (m *Mapper) parseUserID(id string) (user.UserID, error) {
	value, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return user.UserID{}, err
	}
	return user.NewUserID(value), nil
}

// parseGender 解析性别
func (m *Mapper) parseGender(gender string) user.UserGender {
	switch gender {
	case "male":
		return user.UserGenderMale
	case "female":
		return user.UserGenderFemale
	case "other":
		return user.UserGenderOther
	default:
		return user.UserGenderUnknown
	}
}

// parseStatus 解析状态
func (m *Mapper) parseStatus(status string) user.UserStatus {
	switch status {
	case "active":
		return user.UserStatusActive
	case "inactive":
		return user.UserStatusInactive
	case "pending":
		return user.UserStatusPending
	case "locked":
		return user.UserStatusLocked
	default:
		return user.UserStatusPending
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
func (m *Mapper) ParseUserID(id string) (user.UserID, error) {
	return m.parseUserID(id)
}
