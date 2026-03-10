// Package assembler User DTO 转换器
package assembler

import (
	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// ToDTO 将领域实体转换为 User DTO
func ToDTO(user *entity.User) *dto.User {
	if user == nil {
		return nil
	}

	return &dto.User{
		ID:        user.ID.String(),
		Email:     user.Email.String(),
		Nickname: user.Nickname.String(),
		Phone:     user.Phone,
		Bio:       user.Bio,
		Avatar:    user.Avatar,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToLoginResponseDTO 将用户实体和 Token 转换为登录响应 DTO
func ToLoginResponseDTO(user*entity.User, accessToken string) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	return &dto.LoginResponse{
		User:        ToDTO(user),
		AccessToken: accessToken,
	}
}

// FromRegisterRequest 从注册请求创建用户实体
func FromRegisterRequest(req *dto.RegisterRequest, hashedPassword entity.HashedPassword) (*entity.User, error) {
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
