package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	userAggregate "github.com/shenfay/go-ddd-scaffold/internal/domain/user/aggregate"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
)

// GetUserUseCase 获取用户用例
type GetUserUseCase struct {
	uow application.UnitOfWork
}

// NewGetUserUseCase 创建获取用户用例
func NewGetUserUseCase(uow application.UnitOfWork) *GetUserUseCase {
	return &GetUserUseCase{
		uow: uow,
	}
}

// Execute 执行获取用户用例
func (uc *GetUserUseCase) Execute(ctx context.Context, userID vo.UserID) (*userAggregate.User, error) {
	userRepo := uc.uow.UserRepository()
	return userRepo.FindByID(ctx, userID)
}
