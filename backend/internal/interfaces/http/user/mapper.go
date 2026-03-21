package user

import (
	"strconv"

	userApp "github.com/shenfay/go-ddd-scaffold/internal/application/user"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// ToUpdateProfileRequest 转换为更新用户资料请求
func (m *Mapper) ToUpdateProfileRequest(req *UpdateUserRequest) *userApp.UpdateProfileRequest {
	return &userApp.UpdateProfileRequest{
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Gender:      m.parseGenderPtr(req.Gender),
		PhoneNumber: req.PhoneNumber,
	}
}

// parseGenderPtr 解析性别指针
func (m *Mapper) parseGenderPtr(genderStr *string) *vo.UserGender {
	if genderStr == nil {
		return nil
	}
	gender := m.parseGender(*genderStr)
	return &gender
}

// ToChangePasswordRequest 转换为修改密码请求
func (m *Mapper) ToChangePasswordRequest(req *ChangePasswordRequest) *userApp.ChangePasswordRequest {
	return &userApp.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}

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

// ParseUserID 解析用户 ID（公开方法）
func (m *Mapper) ParseUserID(id string) (vo.UserID, error) {
	return m.parseUserID(id)
}
