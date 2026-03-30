package user

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handlers"
)

// UpdateProfileHandler 更新用户资料处理器
type UpdateProfileHandler struct {
	uc          *usecase.UpdateProfileUseCase
	respHandler *handlers.Handler
}

// NewUpdateProfileHandler 创建更新用户资料处理器
func NewUpdateProfileHandler(
	uc *usecase.UpdateProfileUseCase,
	respHandler *handlers.Handler,
) *UpdateProfileHandler {
	return &UpdateProfileHandler{
		uc:          uc,
		respHandler: respHandler,
	}
}

// ServeHTTP 更新用户资料
// @Summary 更新当前用户资料
// @Description 更新当前登录用户的详细信息（从 Token 中获取用户身份）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body usecase.UpdateProfileCommand true "用户更新信息"
// @Success 200 {object} handler.APIResponse "更新成功"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 401 {object} handler.APIResponse "未认证"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/profile [put]
func (h *UpdateProfileHandler) Handle(c *gin.Context) {
	// 从 JWT Token 中获取用户 ID（由认证中间件注入）
	userID, exists := c.Get("user_id")
	if !exists {
		h.respHandler.Unauthorized(c, "user not authenticated")
		return
	}

	var cmd usecase.UpdateProfileCommand
	cmd.UserID = vo.NewUserID(userID.(int64))

	result, err := h.uc.Execute(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	if !result.Success {
		h.respHandler.Error(c, fmt.Errorf("update profile failed"))
		return
	}

	h.respHandler.Success(c, nil)
}
