package user

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/usecase"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user/valueobject"
	"github.com/shenfay/go-ddd-scaffold/internal/interfaces/http/handler"
)

// GetUserHandler 获取用户处理器
type GetUserHandler struct {
	uc          *usecase.GetUserUseCase
	respHandler *handler.Handler
}

// NewGetUserHandler 创建获取用户处理器
func NewGetUserHandler(
	uc *usecase.GetUserUseCase,
	respHandler *handler.Handler,
) *GetUserHandler {
	return &GetUserHandler{
		uc:          uc,
		respHandler: respHandler,
	}
}

// ServeHTTP 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} interface{} "用户详情"
// @Failure 400 {object} handler.APIResponse "请求参数错误"
// @Failure 404 {object} handler.APIResponse "用户不存在"
// @Router /users/{id} [get]
func (h *GetUserHandler) ServeHTTP(c *gin.Context) {
	userIDStr := c.Param("id")
	userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}
	userID := vo.NewUserID(userIDInt)

	result, err := h.uc.Execute(c.Request.Context(), userID)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}
