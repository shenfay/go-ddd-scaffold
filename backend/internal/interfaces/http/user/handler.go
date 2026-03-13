package user

import (
	"github.com/gin-gonic/gin"

	"github.com/shenfay/go-ddd-scaffold/internal/application/user/commands"
	"github.com/shenfay/go-ddd-scaffold/internal/application/user/queries"
	httpShared "github.com/shenfay/go-ddd-scaffold/internal/interfaces/http"
)

// Handler HTTP处理器
type Handler struct {
	createUserHandler     *commands.CreateUserHandler
	updateUserHandler     *commands.UpdateUserHandler
	getUserHandler        *queries.GetUserHandler
	listUsersHandler      *queries.ListUsersHandler
	deactivateUserHandler *commands.DeactivateUserHandler
	activateUserHandler   *commands.ActivateUserHandler
	changePasswordHandler *commands.ChangePasswordHandler
	mapper                *Mapper
	respHandler           *httpShared.Handler
}

// HandlerOption Handler配置选项
type HandlerOption func(*Handler)

// WithCreateUserHandler 设置创建用户处理器
func WithCreateUserHandler(h *commands.CreateUserHandler) HandlerOption {
	return func(handler *Handler) {
		handler.createUserHandler = h
	}
}

// WithUpdateUserHandler 设置更新用户处理器
func WithUpdateUserHandler(h *commands.UpdateUserHandler) HandlerOption {
	return func(handler *Handler) {
		handler.updateUserHandler = h
	}
}

// WithGetUserHandler 设置获取用户处理器
func WithGetUserHandler(h *queries.GetUserHandler) HandlerOption {
	return func(handler *Handler) {
		handler.getUserHandler = h
	}
}

// WithListUsersHandler 设置列出用户处理器
func WithListUsersHandler(h *queries.ListUsersHandler) HandlerOption {
	return func(handler *Handler) {
		handler.listUsersHandler = h
	}
}

// WithDeactivateUserHandler 设置禁用用户处理器
func WithDeactivateUserHandler(h *commands.DeactivateUserHandler) HandlerOption {
	return func(handler *Handler) {
		handler.deactivateUserHandler = h
	}
}

// WithActivateUserHandler 设置激活用户处理器
func WithActivateUserHandler(h *commands.ActivateUserHandler) HandlerOption {
	return func(handler *Handler) {
		handler.activateUserHandler = h
	}
}

// WithChangePasswordHandler 设置修改密码处理器
func WithChangePasswordHandler(h *commands.ChangePasswordHandler) HandlerOption {
	return func(handler *Handler) {
		handler.changePasswordHandler = h
	}
}

// WithResponseHandler 设置响应处理器
func WithResponseHandler(h *httpShared.Handler) HandlerOption {
	return func(handler *Handler) {
		handler.respHandler = h
	}
}

// NewHandler 创建处理器
func NewHandler(opts ...HandlerOption) *Handler {
	h := &Handler{
		mapper: NewMapper(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// CreateUser 创建用户
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if !h.respHandler.BindJSON(c, &req) {
		return
	}

	cmd := h.mapper.ToCreateUserCommand(&req)
	result, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Created(c, result)
}

// GetUser 获取用户
func (h *Handler) GetUser(c *gin.Context) {
	var req GetUserRequest
	if !h.respHandler.BindUri(c, &req) {
		return
	}

	query, err := h.mapper.ToGetUserQuery(req.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	result, err := h.getUserHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}

// ListUsers 列出用户
func (h *Handler) ListUsers(c *gin.Context) {
	var req ListUsersRequest
	if !h.respHandler.BindQuery(c, &req) {
		return
	}

	query := h.mapper.ToListUsersQuery(&req)
	result, err := h.listUsersHandler.Handle(c.Request.Context(), query)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Page(c, result.Items, result.TotalCount, result.Page, result.PageSize)
}

// UpdateUser 更新用户
func (h *Handler) UpdateUser(c *gin.Context) {
	var uriReq GetUserRequest
	if !h.respHandler.BindUri(c, &uriReq) {
		return
	}

	var bodyReq UpdateUserRequest
	if !h.respHandler.BindJSON(c, &bodyReq) {
		return
	}

	cmd, err := h.mapper.ToUpdateUserCommand(&bodyReq, uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	result, err := h.updateUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.Success(c, result)
}

// ActivateUser 激活用户
func (h *Handler) ActivateUser(c *gin.Context) {
	var req ActivateUserRequest
	if !h.respHandler.BindUri(c, &req) {
		return
	}

	cmd, err := h.mapper.ToActivateUserCommand(req.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	if err := h.activateUserHandler.Handle(c.Request.Context(), cmd); err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// DeactivateUser 禁用用户
func (h *Handler) DeactivateUser(c *gin.Context) {
	var uriReq GetUserRequest
	if !h.respHandler.BindUri(c, &uriReq) {
		return
	}

	var bodyReq DeactivateUserRequest
	if !h.respHandler.BindJSON(c, &bodyReq) {
		return
	}

	cmd, err := h.mapper.ToDeactivateUserCommand(&bodyReq, uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	if err := h.deactivateUserHandler.Handle(c.Request.Context(), cmd); err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}

// ChangePassword 修改密码
func (h *Handler) ChangePassword(c *gin.Context) {
	var uriReq GetUserRequest
	if !h.respHandler.BindUri(c, &uriReq) {
		return
	}

	var bodyReq ChangePasswordRequest
	if !h.respHandler.BindJSON(c, &bodyReq) {
		return
	}

	cmd, err := h.mapper.ToChangePasswordCommand(&bodyReq, uriReq.UserID)
	if err != nil {
		h.respHandler.BadRequest(c, "invalid user id")
		return
	}

	if err := h.changePasswordHandler.Handle(c.Request.Context(), cmd); err != nil {
		h.respHandler.Error(c, err)
		return
	}

	h.respHandler.NoContent(c)
}
