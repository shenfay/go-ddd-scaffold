package user

import (
	"time"

	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// ConvertUserToDTO 将领域对象转换为 DTO
func ConvertUserToDTO(user *user.User) *UserDTO {
	return &UserDTO{
		ID:          user.ID().(vo.UserID).Int64(),
		Username:    user.Username().Value(),
		Email:       user.Email().Value(),
		DisplayName: user.DisplayName(),
		Status:      user.Status().String(),
		CreatedAt:   user.CreatedAt(),
	}
}
