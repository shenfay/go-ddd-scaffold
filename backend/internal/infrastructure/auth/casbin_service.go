package auth

import (
	"errors"

	"github.com/casbin/casbin/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CasbinService 权限服务接口
type CasbinService interface {
	// Enforce 权限检查
	// sub: 用户ID
	// dom: 租户ID
	// obj: 资源
	// act: 操作
	Enforce(sub, dom, obj, act string) (bool, error)

	// AddRoleForUser 在租户内添加角色
	AddRoleForUser(userID, tenantID uuid.UUID, role string) error

	// RemoveRoleForUser 移除用户在租户内的角色
	RemoveRoleForUser(userID, tenantID uuid.UUID, role string) error

	// GetRolesForUser 获取用户在租户内的所有角色
	GetRolesForUser(userID, tenantID uuid.UUID) []string

	// AddPermission 添加权限策略
	AddPermission(role, tenantID, resource, action string) error

	// RemovePermission 移除权限策略
	RemovePermission(role, tenantID, resource, action string) error

	// GetPermissionsForRole 获取角色的所有权限
	GetPermissionsForRole(role, tenantID string) [][]string

	// DeleteAllRolesForUser 删除用户在租户内的所有角色
	DeleteAllRolesForUser(userID, tenantID uuid.UUID) error
}

// casbinService Casbin 权限服务实现
type casbinService struct {
	enforcer *casbin.Enforcer
}

// NewCasbinService 创建 Casbin 权限服务
func NewCasbinService(enforcer *casbin.Enforcer) CasbinService {
	return &casbinService{
		enforcer: enforcer,
	}
}

// Enforce 权限检查
func (s *casbinService) Enforce(sub, dom, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, dom, obj, act)
}

// AddRoleForUser 在租户内添加角色
func (s *casbinService) AddRoleForUser(userID, tenantID uuid.UUID, role string) error {
	if userID == uuid.Nil {
		return errors.New("userID 不能为空")
	}
	if tenantID == uuid.Nil {
		return errors.New("tenantID 不能为空")
	}
	if role == "" {
		return errors.New("role 不能为空")
	}

	_, err := s.enforcer.AddRoleForUserInDomain(userID.String(), role, tenantID.String())
	return err
}

// RemoveRoleForUser 移除用户在租户内的角色
func (s *casbinService) RemoveRoleForUser(userID, tenantID uuid.UUID, role string) error {
	if userID == uuid.Nil {
		return errors.New("userID 不能为空")
	}
	if tenantID == uuid.Nil {
		return errors.New("tenantID 不能为空")
	}
	if role == "" {
		return errors.New("role 不能为空")
	}

	_, err := s.enforcer.DeleteRoleForUserInDomain(userID.String(), role, tenantID.String())
	return err
}

// GetRolesForUser 获取用户在租户内的所有角色
func (s *casbinService) GetRolesForUser(userID, tenantID uuid.UUID) []string {
	if userID == uuid.Nil {
		return nil
	}

	// 如果没有提供租户ID，获取用户在所有域的角色
	if tenantID == uuid.Nil {
		roles, _ := s.enforcer.GetRolesForUser(userID.String())
		return roles
	}

	// 获取在特定租户内的角色
	roles := s.enforcer.GetRolesForUserInDomain(userID.String(), tenantID.String())
	return roles
}

// AddPermission 添加权限策略
func (s *casbinService) AddPermission(role, tenantID, resource, action string) error {
	if role == "" {
		return errors.New("role 不能为空")
	}
	if tenantID == "" {
		return errors.New("tenantID 不能为空")
	}
	if resource == "" {
		return errors.New("resource 不能为空")
	}
	if action == "" {
		return errors.New("action 不能为空")
	}

	_, err := s.enforcer.AddPolicy(role, tenantID, resource, action)
	return err
}

// RemovePermission 移除权限策略
func (s *casbinService) RemovePermission(role, tenantID, resource, action string) error {
	if role == "" {
		return errors.New("role 不能为空")
	}
	if tenantID == "" {
		return errors.New("tenantID 不能为空")
	}
	if resource == "" {
		return errors.New("resource 不能为空")
	}
	if action == "" {
		return errors.New("action 不能为空")
	}

	_, err := s.enforcer.RemovePolicy(role, tenantID, resource, action)
	return err
}

// GetPermissionsForRole 获取角色的所有权限
func (s *casbinService) GetPermissionsForRole(role, tenantID string) [][]string {
	if role == "" {
		return nil
	}

	if tenantID == "" {
		perms, _ := s.enforcer.GetPermissionsForUser(role)
		return perms
	}

	perms := s.enforcer.GetPermissionsForUserInDomain(role, tenantID)
	return perms
}

// DeleteAllRolesForUser 删除用户在租户内的所有角色
func (s *casbinService) DeleteAllRolesForUser(userID, tenantID uuid.UUID) error {
	if userID == uuid.Nil {
		return errors.New("userID 不能为空")
	}
	if tenantID == uuid.Nil {
		return errors.New("tenantID 不能为空")
	}

	// 先获取用户在租户内的所有角色，然后逐一删除
	roles := s.GetRolesForUser(userID, tenantID)
	for _, role := range roles {
		_, err := s.enforcer.DeleteRoleForUserInDomain(userID.String(), role, tenantID.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCasbinServiceForTest 为测试创建 Casbin 服务（使用默认配置）
func NewCasbinServiceForTest(db *gorm.DB) (CasbinService, error) {
	enforcer, err := NewCasbinEnforcer(db)
	if err != nil {
		return nil, err
	}
	return NewCasbinService(enforcer), nil
}
