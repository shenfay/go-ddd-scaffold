// Package mapper 实体映射器
//
// 负责领域对象与持久化模型之间的双向转换
package mapper

import (
	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
	"github.com/google/uuid"
)

// EntityMapper 实体映射器接口
type EntityMapper interface {
	// ToEntity 将持久化模型转换为领域实体
	ToEntity(model *model.User) (*entity.User, error)

	// ToModel 将领域实体转换为持久化模型
	ToModel(entity *entity.User) (*model.User, error)
}

// entityMapperImpl EntityMapper 实现
type entityMapperImpl struct{}

// NewEntityMapper 创建 EntityMapper 实例
func NewEntityMapper() EntityMapper {
	return &entityMapperImpl{}
}

// ToEntity 将持久化模型转换为领域实体
func (m *entityMapperImpl) ToEntity(userModel *model.User) (*entity.User, error) {
	if userModel == nil {
		return nil, nil
	}

	user := &entity.User{}

	// 设置 ID
	if userModel.ID != nil {
		id, err := uuid.Parse(*userModel.ID)
		if err != nil {
			return nil, err
		}
		user.ID = id
	}

	// 设置时间
	if userModel.CreatedAt != nil {
		user.CreatedAt = *userModel.CreatedAt
	}
	if userModel.UpdatedAt != nil {
		user.UpdatedAt = *userModel.UpdatedAt
	}

	// 设置邮箱（使用值对象）
	if userModel.Email != "" {
		email, err := valueobject.NewEmail(userModel.Email)
		if err != nil {
			return nil, err
		}
		user.Email = email
	}

	// 设置密码（使用值对象）
	if userModel.Password != "" {
		hashedPassword, err := entity.NewHashedPassword(userModel.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	// 设置昵称（使用值对象）
	if userModel.Nickname != "" {
		nickname, err := valueobject.NewNickname(userModel.Nickname)
		if err != nil {
			return nil, err
		}
		user.Nickname = nickname
	}

	// 设置可选字段
	user.Avatar = userModel.Avatar
	user.Phone = userModel.Phone
	user.Bio = userModel.Bio

	// 设置状态
	if userModel.Status != nil {
		status, err := entity.ParseUserStatus(*userModel.Status)
		if err != nil {
			return nil, err
		}
		user.Status = status
	}

	return user, nil
}

// ToModel 将领域实体转换为持久化模型
func (m *entityMapperImpl) ToModel(user *entity.User) (*model.User, error) {
	if user == nil {
		return nil, nil
	}

	userModel := &model.User{}

	// 设置 ID
	idStr := user.ID.String()
	userModel.ID = &idStr

	// 设置时间
	userModel.CreatedAt = &user.CreatedAt
	userModel.UpdatedAt = &user.UpdatedAt

	// 设置基础字段（值对象）
	userModel.Email = user.Email.String()
	userModel.Password = user.Password.String()
	userModel.Nickname = user.Nickname.String()

	// 设置可选字段
	userModel.Avatar = user.Avatar
	userModel.Phone = user.Phone
	userModel.Bio = user.Bio

	// 设置状态
	statusStr := user.Status.String()
	userModel.Status = &statusStr

	return userModel, nil
}
