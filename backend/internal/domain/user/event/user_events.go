package event

import (
	"time"

	"go-ddd-scaffold/internal/domain/user/entity"

	"github.com/google/uuid"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	UserID    uuid.UUID       `json:"userId"`
	Email     string          `json:"email"`
	Role      entity.UserRole `json:"role"`
	TenantID  *uuid.UUID      `json:"tenantId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// UserLoggedInEvent 用户登录事件
type UserLoggedInEvent struct {
	UserID    uuid.UUID `json:"userId"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip,omitempty"`
	UserAgent string    `json:"userAgent,omitempty"`
}
