package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenfay/go-ddd-scaffold/internal/app/user"
	"github.com/shenfay/go-ddd-scaffold/internal/transport/http/response"
	validationErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/validation"
)

// UserHandler handles user management HTTP requests.
type UserHandler struct {
	userService *user.Service
}

// NewUserHandler creates a new user handler instance.
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

// CreateUser handles user registration via admin API.
//
// Creates a new user account with email and password.
// The email must be unique and password must meet security requirements.
//
// @Summary Create a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body object true "User creation data"
// @Success 201 {object} response.SuccessResponse{data=user.UserDTO}
// @Failure 400 {object} response.ErrorResponse "Validation error"
// @Failure 409 {object} response.ErrorResponse "Email already exists"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
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
