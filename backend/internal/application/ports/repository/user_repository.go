package repository

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// UserRepository 用户仓储端口
type UserRepository interface {
	// Save 保存用户
	Save(ctx context.Context, user *aggregate.User) error

	// FindByID 根据 ID 查找用户
	FindByID(ctx context.Context, id vo.UserID) (*aggregate.User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email vo.Email) (*aggregate.User, error)

	// FindByUsername 根据用户名查找用户
	FindByUsername(ctx context.Context, username vo.UserName) (*aggregate.User, error)

	// Update 更新用户
	Update(ctx context.Context, user *aggregate.User) error

	// Delete 删除用户
	Delete(ctx context.Context, id vo.UserID) error
}
