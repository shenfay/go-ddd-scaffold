package user

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user"
	domainUser "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	httpiface "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// UserHandler 用户领域 HTTP 处理器
type UserHandler struct {
	userService user.UserService
	handler     *httpiface.Handler
	isMock      bool // Mock 模式标志
}

// NewUserHandler 创建用户领域 HTTP 处理器
func NewUserHandler(
	userService user.UserService,
	handler *httpiface.Handler,
) *UserHandler {
	return &UserHandler{
		userService: userService,
		handler:     handler,
		isMock:      userService == nil,
	}
}

// NewMockUserHandler 创建 Mock 用户处理器（用于演示）
func NewMockUserHandler(handler *httpiface.Handler) *UserHandler {
	return &UserHandler{
		userService: nil,
		handler:     handler,
		isMock:      true,
	}
}

// mockResponse 返回 Mock 响应
func (h *UserHandler) mockResponse(c *gin.Context, action string, data ...gin.H) {
	response := gin.H{
		"message": "Mock: " + action,
		"status":  "pending_implementation",
	}

	// 合并额外数据
	if len(data) > 0 {
		for k, v := range data[0] {
			response[k] = v
		}
	}

	h.handler.Success(c, response)
}

// withRealServiceOrMock 执行真实业务逻辑或返回 Mock 响应
func (h *UserHandler) withRealServiceOrMock(c *gin.Context, action string, fn func() error) {
	if h.isMock {
		h.mockResponse(c, action)
		return
	}

	if err := fn(); err != nil {
		h.handler.Error(c, err)
	}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	DisplayName *string                `json:"display_name,omitempty"`
	FirstName   *string                `json:"first_name,omitempty"`
	LastName    *string                `json:"last_name,omitempty"`
	Gender      *domainUser.UserGender `json:"gender,omitempty"`
	PhoneNumber *string                `json:"phone_number,omitempty"`
	AvatarURL   string                 `json:"avatar_url,omitempty"`
}

// CreateUser 创建用户
// POST /api/v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "User created", func() error {
		var req CreateUserRequest
		if !h.handler.BindJSON(c, &req) {
			return nil // BindJSON 已经处理了错误
		}

		cmd := &user.RegisterUserCommand{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		newUser, err := h.userService.RegisterUser(c.Request.Context(), cmd)
		if err != nil {
			return err
		}

		h.handler.Created(c, h.toUserResponse(newUser))
		return nil
	})
}

// GetUserByID 根据 ID 获取用户
// GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	h.withRealServiceOrMock(c, "Get user by ID", func() error {
		userID := util.ToInt64(c.Param("id"))
		if userID == 0 {
			h.handler.BadRequest(c, "invalid user id")
			return nil
		}

		u, err := h.userService.GetUserByID(c.Request.Context(), domainUser.NewUserID(userID))
		if err != nil {
			return err
		}

		h.handler.Success(c, h.toUserResponse(u))
		return nil
	})
}

// UpdateUser 更新用户
// PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "Update user", func() error {
		userID := util.ToInt64(c.Param("id"))
		if userID == 0 {
			h.handler.BadRequest(c, "invalid user id")
			return nil
		}

		var req UpdateUserRequest
		if !h.handler.BindJSON(c, &req) {
			return nil
		}

		cmd := &user.UpdateUserProfileCommand{
			UserID:      domainUser.NewUserID(userID),
			DisplayName: req.DisplayName,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Gender:      req.Gender,
			PhoneNumber: req.PhoneNumber,
		}

		if err := h.userService.UpdateUserProfile(c.Request.Context(), cmd); err != nil {
			return err
		}

		updatedUser, err := h.userService.GetUserByID(c.Request.Context(), domainUser.NewUserID(userID))
		if err != nil {
			return err
		}

		h.handler.Success(c, h.toUserResponse(updatedUser))
		return nil
	})
}

// DeleteUser 删除用户
// DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "Delete user", func() error {
		// TODO: 实现 DeleteUser 命令
		h.handler.NoContent(c)
		return nil
	})
}

// ActivateUser 激活用户
// POST /api/v1/users/:id/activate
func (h *UserHandler) ActivateUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "Activate user", func() error {
		userID := util.ToInt64(c.Param("id"))
		if userID == 0 {
			h.handler.BadRequest(c, "invalid user id")
			return nil
		}

		cmd := &user.ActivateUserCommand{
			UserID: domainUser.NewUserID(userID),
		}

		if err := h.userService.ActivateUser(c.Request.Context(), cmd); err != nil {
			return err
		}

		h.handler.Accepted(c, gin.H{"status": "activated", "user_id": userID})
		return nil
	})
}

// DeactivateUser 禁用用户
// POST /api/v1/users/:id/deactivate
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "Deactivate user", func() error {
		userID := util.ToInt64(c.Param("id"))
		if userID == 0 {
			h.handler.BadRequest(c, "invalid user id")
			return nil
		}

		var req struct {
			Reason string `json:"reason,omitempty"`
		}
		if !h.handler.BindJSON(c, &req) {
			return nil
		}

		cmd := &user.DeactivateUserCommand{
			UserID: domainUser.NewUserID(userID),
			Reason: req.Reason,
		}

		if err := h.userService.DeactivateUser(c.Request.Context(), cmd); err != nil {
			return err
		}

		h.handler.Accepted(c, gin.H{"status": "deactivated", "user_id": userID})
		return nil
	})
}

// AuthenticateUser 用户认证（登录）
// POST /api/v1/users/authenticate
func (h *UserHandler) AuthenticateUser(c *gin.Context) {
	h.withRealServiceOrMock(c, "Authenticate user", func() error {
		var req struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}
		if !h.handler.BindJSON(c, &req) {
			return nil
		}

		cmd := &user.AuthenticateUserCommand{
			Username:  req.Username,
			Password:  req.Password,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		result, err := h.userService.AuthenticateUser(c.Request.Context(), cmd)
		if err != nil {
			return err
		}

		h.handler.Success(c, h.toAuthResponse(result))
		return nil
	})
}

// toUserResponse 将领域用户转换为响应对象
func (h *UserHandler) toUserResponse(u *domainUser.User) gin.H {
	userID := u.ID().(domainUser.UserID)
	response := gin.H{
		"id":           userID.Int64(),
		"username":     u.Username().Value(),
		"email":        u.Email().Value(),
		"status":       u.Status().String(),
		"display_name": u.DisplayName(),
		"first_name":   u.FirstName(),
		"last_name":    u.LastName(),
		"gender":       u.Gender().String(),
		"phone_number": u.PhoneNumber(),
		"avatar_url":   u.AvatarURL(),
		"login_count":  u.LoginCount(),
		"created_at":   u.CreatedAt(),
		"updated_at":   u.UpdatedAt(),
	}

	if u.LastLoginAt() != nil {
		response["last_login_at"] = u.LastLoginAt()
	}

	return response
}

// toAuthResponse 将认证结果转换为响应对象
func (h *UserHandler) toAuthResponse(result *user.AuthenticationResult) gin.H {
	return gin.H{
		"user_id":       result.UserID.Int64(),
		"username":      result.Username,
		"email":         result.Email,
		"token":         result.Token,
		"refresh_token": result.RefreshToken,
		"expires_at":    result.ExpiresAt,
	}
}
