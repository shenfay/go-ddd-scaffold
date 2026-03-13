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

// @Summary 创建用户
// @Description 创建一个新的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "用户创建信息"
// @Success 201 {object} CreateUserResponse "创建成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 409 {object} httpShared.APIResponse "用户已存在"
// @Router /users [post]
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

// @Summary 获取用户详情
// @Description 根据用户 ID 获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 200 {object} UserResponse "用户详情"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [get]
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

// @Summary 列出用户
// @Description 分页获取用户列表，支持关键词搜索和状态筛选
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param keyword query string false "搜索关键词"
// @Param status query string false "用户状态 (active/inactive/pending/locked)"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} UserListResponse "分页用户列表"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Router /users [get]
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

// @Summary 更新用户信息
// @Description 更新指定用户的详细信息（部分更新）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body UpdateUserRequest true "用户更新信息"
// @Success 200 {object} UserResponse "更新成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id} [put]
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

// @Summary 激活用户
// @Description 激活已禁用的用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Success 204 {object} httpShared.APIResponse "激活成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id}/activate [post]
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

// @Summary 禁用用户
// @Description 禁用指定用户账户（可填写原因）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body DeactivateUserRequest false "禁用原因"
// @Success 204 {object} httpShared.APIResponse "禁用成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id}/deactivate [post]
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

// @Summary 修改用户密码
// @Description 修改指定用户的登录密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户 ID"
// @Param request body ChangePasswordRequest true "原密码和新密码"
// @Success 204 {object} httpShared.APIResponse "修改成功"
// @Failure 400 {object} httpShared.APIResponse "请求参数错误"
// @Failure 404 {object} httpShared.APIResponse "用户不存在"
// @Router /users/{id}/password [put]
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
