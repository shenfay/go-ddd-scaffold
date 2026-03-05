// Package converter 用户模块统一转换器
//
// 负责领域对象与 DTO 之间的双向转换
package converter

import (
	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserConverter 用户转换器
type UserConverter struct{}

// NewUserConverter 创建用户转换器实例
func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToUserDTO 将领域实体转换为 User DTO
func (c *UserConverter) ToUserDTO(user *entity.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		ID:        user.ID.String(),
		Email:     user.Email.String(),
		Nickname:  user.Nickname.String(),
		Phone:     user.Phone,
		Bio:       user.Bio,
		Avatar:    user.Avatar,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToLoginResponseDTO 将用户实体和 Token 转换为登录响应 DTO
func (c *UserConverter) ToLoginResponseDTO(user *entity.User, accessToken string) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	return &dto.LoginResponse{
		User:        c.ToUserDTO(user),
		AccessToken: accessToken,
	}
}

// FromRegisterRequest 从注册请求创建用户实体
func (c *UserConverter) FromRegisterRequest(req *dto.RegisterRequest, hashedPassword entity.HashedPassword) (*entity.User, error) {
	if req == nil {
		return nil, nil
	}

	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	nickname, err := valueobject.NewNickname(req.Nickname)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    email,
		Password: hashedPassword,
		Nickname: nickname,
		Status:   entity.StatusActive,
	}

	return user, nil
}

// UpdateFromProfileRequest 从个人资料更新请求更新用户实体
func (c *UserConverter) UpdateFromProfileRequest(user *entity.User, req *dto.UpdateProfileRequest) error {
	if user == nil || req == nil {
		return nil
	}

	if req.Nickname != nil {
		nickname, err := valueobject.NewNickname(*req.Nickname)
		if err != nil {
			return err
		}
		user.Nickname = nickname
	}

	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	return nil
}
