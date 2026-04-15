package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/response"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	validationErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/validation"
)

// VerifyEmailRequest 邮箱验证请求
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"` // 邮箱验证令牌
}

// VerifyEmail 处理邮箱验证
//
// @Summary 验证邮箱
// @Description 通过验证令牌确认邮箱所有权
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body VerifyEmailRequest true "邮箱验证请求"
// @Success 200 {object} response.Response "验证成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "令牌无效或已过期"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Email verified successfully",
	})
}

// ResendVerificationEmail 请求重新发送验证邮件
//
// @Summary 重新发送验证邮件
// @Description 向用户邮箱发送新的验证链接
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "邮件发送成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 404 {object} response.ErrorResponse "用户不存在"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Security BearerAuth
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	// 从上下文获取用户ID(需要已登录)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, authErr.ErrInvalidCredentials)
		return
	}

	if err := h.service.SendVerificationEmail(c.Request.Context(), userID.(string)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "Verification email sent",
	})
}
