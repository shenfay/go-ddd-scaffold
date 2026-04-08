package user

import "context"

// UserRepository 用户仓储接口
type UserRepository interface {
	// Create 创建用户
	Create(ctx context.Context, user *User) error

	// Save 保存用户（创建或更新）
	Save(ctx context.Context, user *User) error

	// FindByID 根据 ID 查找用户
	FindByID(ctx context.Context, id string) (*User, error)

	// FindByEmail 根据邮箱查找用户
	FindByEmail(ctx context.Context, email string) (*User, error)

	// ExistsByEmail 检查邮箱是否存在
	ExistsByEmail(ctx context.Context, email string) bool

	// Update 更新用户
	Update(ctx context.Context, user *User) error
}
