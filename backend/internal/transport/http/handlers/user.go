package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/app/user"
)

// UserHandler 用户HTTP处理器
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler 创建用户处理器
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

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	cmd := user.CreateUserCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	dto, err := h.userService.CreateUser(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, dto)
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	dto, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	cmd := user.UpdateProfileCommand{
		UserID: userID,
		Email:  req.Email,
	}

	dto, err := h.userService.UpdateProfile(c.Request.Context(), cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, dto)
}
