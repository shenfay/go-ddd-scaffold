package usecase

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/application"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
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
func (uc *GetUserUseCase) Execute(ctx context.Context, userID vo.UserID) (*GetUserResult, error) {
	userRepo := uc.uow.UserRepository()
	userEntity, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &GetUserResult{
		UserDTO: user.ConvertUserToDTO(userEntity),
	}, nil
}
