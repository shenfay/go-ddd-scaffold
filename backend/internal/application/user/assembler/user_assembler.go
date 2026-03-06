// Package assembler 用户模块 Assembler 层
//
// 负责领域对象与 DTO 之间的双向转换
package assembler

import (
	"github.com/google/uuid"

	"go-ddd-scaffold/internal/application/user/dto"
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserAssembler 用户 Assembler 接口
type UserAssembler interface {
	// ToDTO 将领域实体转换为 User DTO
	ToDTO(user *entity.User) *dto.User

	// ToLoginResponseDTO 将用户实体和 Token 转换为登录响应 DTO
	ToLoginResponseDTO(user *entity.User, accessToken string) *dto.LoginResponse

	// FromRegisterRequest 从注册请求创建用户实体
	FromRegisterRequest(req *dto.RegisterRequest, hashedPassword entity.HashedPassword) (*entity.User, error)

	// UpdateFromProfileRequest 从个人资料更新请求更新用户实体
	UpdateFromProfileRequest(user *entity.User, req *dto.UpdateProfileRequest) error

	// ToUUIDPtr 将 string 转换为 *uuid.UUID
	ToUUIDPtr(s string) (*uuid.UUID, error)

	// ToStringPtr 将 uuid.UUID 转换为 *string
	ToStringPtr(id uuid.UUID) *string
}

// userAssemblerImpl UserAssembler 实现
type userAssemblerImpl struct{}

// NewUserAssembler 创建 UserAssembler 实例
func NewUserAssembler() UserAssembler {
	return &userAssemblerImpl{}
}

// ToDTO 将领域实体转换为 User DTO
func (a *userAssemblerImpl) ToDTO(user *entity.User) *dto.User {
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
func (a *userAssemblerImpl) ToLoginResponseDTO(user *entity.User, accessToken string) *dto.LoginResponse {
	if user == nil {
		return nil
	}

	return &dto.LoginResponse{
		User:        a.ToDTO(user),
		AccessToken: accessToken,
	}
}

// FromRegisterRequest 从注册请求创建用户实体
func (a *userAssemblerImpl) FromRegisterRequest(req *dto.RegisterRequest, hashedPassword entity.HashedPassword) (*entity.User, error) {
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
		ID:         uuid.New(), // ✅ 生成新 UUID
		Email:      email,
		Password:   hashedPassword,
		Nickname:   nickname,
		Status:     entity.StatusActive,
	}

	return user, nil
}

// UpdateFromProfileRequest 从个人资料更新请求更新用户实体
func (a *userAssemblerImpl) UpdateFromProfileRequest(user *entity.User, req *dto.UpdateProfileRequest) error {
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

// ToUUIDPtr 将 string 转换为 *uuid.UUID
func (a *userAssemblerImpl) ToUUIDPtr(s string) (*uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// ToStringPtr 将 uuid.UUID 转换为 *string
func (a *userAssemblerImpl) ToStringPtr(id uuid.UUID) *string {
	idStr := id.String()
	return &idStr
}
