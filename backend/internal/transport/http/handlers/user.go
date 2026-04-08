package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/app/user"
	_ "github.com/shenfay/go-ddd-scaffold/internal/transport/http/middleware" // for swagger types
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/response"
	validationErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/validation"
)

// UserHandler 用户管理 HTTP 处理器
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterRoutes 注册用户路由
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.CreateUser)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
	}
}

// CreateUser 创建用户（管理员接口）
//
// 使用邮箱和密码创建新用户账户。
// 邮箱必须唯一，密码必须符合安全要求。
//
// @Summary 创建用户
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object true "用户创建数据"
// @Success 201 {object} middleware.SuccessResponse{data=user.UserDTO} "创建成功"
// @Failure 400 {object} middleware.ErrorResponse "Validation error"
// @Failure 409 {object} middleware.ErrorResponse "邮箱已存在"
// @Failure 500 {object} middleware.ErrorResponse "服务器内部错误"
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := user.CreateUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	dto, err := h.userService.CreateUser(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, dto)
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	dto, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, validationErr.FromGinError(err))
		return
	}

	cmd := user.UpdateProfileCommand{
		UserID: userID,
		Email:  req.Email,
	}

	dto, err := h.userService.UpdateProfile(c.Request.Context(), cmd)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, dto)
}
