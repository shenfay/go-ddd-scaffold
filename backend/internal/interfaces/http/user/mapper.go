package user

import (
	"strconv"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/commands"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/queries"
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
func (m *Mapper) ToUpdateUserCommand(req *UpdateUserRequest, userID string) (*commands.UpdateUserCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	cmd := &commands.UpdateUserCommand{
		UserID:      uid,
		DisplayName: req.DisplayName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		AvatarURL:   req.AvatarURL,
	}

	if req.Gender != nil {
		gender := m.parseGender(*req.Gender)
		cmd.Gender = &gender
	}

	return cmd, nil
}

// ToChangePasswordCommand 转换为修改密码命令
func (m *Mapper) ToChangePasswordCommand(req *ChangePasswordRequest, userID string) (*commands.ChangePasswordCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	return &commands.ChangePasswordCommand{
		UserID:      uid,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}, nil
}

// ToDeactivateUserCommand 转换为禁用用户命令
func (m *Mapper) ToDeactivateUserCommand(req *DeactivateUserRequest, userID string) (*commands.DeactivateUserCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	return &commands.DeactivateUserCommand{
		UserID: uid,
		Reason: req.Reason,
	}, nil
}

// ToGetUserQuery 转换为获取用户查询
func (m *Mapper) ToGetUserQuery(userID string) (*queries.GetUserQuery, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	return &queries.GetUserQuery{
		UserID: uid,
	}, nil
}

// ToListUsersQuery 转换为列出用户查询
func (m *Mapper) ToListUsersQuery(req *ListUsersRequest) *queries.ListUsersQuery {
	query := &queries.ListUsersQuery{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	if req.Status != "" {
		status := m.parseStatus(req.Status)
		query.Status = &status
	}

	return query
}

// ToActivateUserCommand 转换为激活用户命令
func (m *Mapper) ToActivateUserCommand(userID string) (*commands.ActivateUserCommand, error) {
	uid, err := m.parseUserID(userID)
	if err != nil {
		return nil, err
	}

	return &commands.ActivateUserCommand{
		UserID: uid,
	}, nil
}
