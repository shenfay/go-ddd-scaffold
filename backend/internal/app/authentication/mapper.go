package authentication

import (
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

// ToAuthResponse 将ServiceAuthResponse转换为HTTP响应
func ToAuthResponse(resp *ServiceAuthResponse) map[string]interface{} {
	return map[string]interface{}{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    resp.ExpiresIn,
		"user": map[string]interface{}{
			"id":             resp.User.ID,
			"email":          resp.User.Email,
			"email_verified": resp.User.EmailVerified,
			"last_login_at":  resp.User.LastLoginAt,
		},
	}
}

// UserToDTO 将领域实体转换为DTO
func UserToDTO(u *user.User) map[string]interface{} {
	return map[string]interface{}{
		"id":             u.ID,
		"email":          u.Email,
		"email_verified": u.EmailVerified,
		"locked":         u.Locked,
		"last_login_at":  u.LastLoginAt,
		"created_at":     u.CreatedAt,
		"updated_at":     u.UpdatedAt,
	}
}

// FormatExpiresIn 格式化过期时间（秒）
func FormatExpiresIn(expire time.Duration) int {
	return int(expire.Seconds())
}
