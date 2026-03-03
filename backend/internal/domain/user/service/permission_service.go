package service

import (
	"errors"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/valueobject"
)

// PermissionService 权限领域服务
// 职责：处理复杂的权限判断逻辑，跨多个实体的权限规则
type PermissionService interface {
	// CanManageChild 家长是否可以管理子女
	CanManageChild(parent, child *entity.User) bool

	// CanAccessResource 用户是否可以访问资源
	CanAccessResource(user *entity.User, resource string, action string) bool

	// CheckTenantAccess 检查用户的租户访问权限
	CheckTenantAccess(user *entity.User, targetTenantID string) error

	// ValidateRoleTransition 验证角色转换是否合法
	ValidateRoleTransition(user *entity.User, newRole entity.UserRole) error
}

// permissionService 实现
type permissionService struct{}

// NewPermissionService 创建权限领域服务
func NewPermissionService() PermissionService {
	return &permissionService{}
}

// CanManageChild 家长是否可以管理子女
func (s *permissionService) CanManageChild(parent, child *entity.User) bool {
	// 业务规则1：必须是家长角色
	if parent.Role != entity.RoleParent {
		return false
	}

	// 业务规则2：子女必须是孩子角色
	if child.Role != entity.RoleChild {
		return false
	}

	// 业务规则3：必须属于同一租户
	if parent.TenantID == nil || child.TenantID == nil {
		return false
	}

	if *parent.TenantID != *child.TenantID {
		return false
	}

	// 业务规则4：子女的ParentID必须指向该家长
	if child.ParentID == nil || *child.ParentID != parent.ID {
		return false
	}

	// 业务规则5：双方都必须是激活状态
	if !parent.IsActive() || !child.IsActive() {
		return false
	}

	return true
}

// CanAccessResource 用户是否可以访问资源
func (s *permissionService) CanAccessResource(user *entity.User, resource string, action string) bool {
	// 业务规则1：超级管理员拥有所有权限
	if user.Role == entity.RoleSuperAdmin {
		return true
	}

	// 业务规则2：内容管理员可以管理内容
	if user.Role == entity.RoleContentAdmin {
		contentResources := map[string]bool{
			"kg_domains":       true,
			"kg_trunks":        true,
			"kg_nodes":         true,
			"kg_relationships": true,
			"learn_content":    true,
		}

		if contentResources[resource] {
			// 内容管理员可以读写内容，但不能删除
			if action == "read" || action == "create" || action == "update" {
				return true
			}
		}
	}

	// 业务规则3：运营管理员可以管理用户（除了超级管理员）
	if user.Role == entity.RoleOpsAdmin {
		if resource == "users" {
			if action == "read" || action == "update" {
				return true
			}
		}
	}

	// 业务规则4：家长可以管理自己租户下的孩子
	if user.Role == entity.RoleParent {
		if resource == "children" {
			return true
		}

		if resource == "learning_progress" && user.TenantID != nil {
			// 可以查看本租户内孩子的学习进度
			return action == "read"
		}
	}

	// 业务规则5：孩子只能查看和更新自己的数据
	if user.Role == entity.RoleChild {
		if resource == "self_profile" || resource == "learning_progress" {
			return action == "read" || action == "update"
		}
	}

	return false
}

// CheckTenantAccess 检查用户的租户访问权限
func (s *permissionService) CheckTenantAccess(user *entity.User, targetTenantID string) error {
	// 业务规则1：超级管理员可以访问所有租户
	if user.Role == entity.RoleSuperAdmin {
		return nil
	}

	// 业务规则2：用户必须属于目标租户
	if user.TenantID == nil {
		return errors.New("用户未分配租户")
	}

	if user.TenantID.String() != targetTenantID {
		return errors.New("无权访问该租户的数据")
	}

	return nil
}

// ValidateRoleTransition 验证角色转换是否合法
func (s *permissionService) ValidateRoleTransition(user *entity.User, newRole entity.UserRole) error {
	// 业务规则1：不能从管理员角色转为普通用户角色
	if user.Role == entity.RoleSuperAdmin || user.Role == entity.RoleContentAdmin || user.Role == entity.RoleOpsAdmin {
		if newRole == entity.RoleParent || newRole == entity.RoleChild {
			return errors.New("管理员角色不能转换为普通用户角色")
		}
	}

	// 业务规则2：孩子不能转换为家长
	if user.Role == entity.RoleChild && newRole == entity.RoleParent {
		return errors.New("孩子角色不能转换为家长角色")
	}

	// 业务规则3：家长不能转换为孩子
	if user.Role == entity.RoleParent && newRole == entity.RoleChild {
		return errors.New("家长角色不能转换为孩子角色")
	}

	// 业务规则4：新角色必须合法
	validRoles := []entity.UserRole{
		entity.RoleParent,
		entity.RoleChild,
		entity.RoleSuperAdmin,
		entity.RoleContentAdmin,
		entity.RoleOpsAdmin,
	}

	roleValid := false
	for _, role := range validRoles {
		if newRole == role {
			roleValid = true
			break
		}
	}

	if !roleValid {
		return errors.New("无效的角色类型")
	}

	return nil
}

// UserValidationService 用户验证领域服务
// 职责：处理用户相关的验证逻辑
type UserValidationService interface {
	// ValidatePasswordStrength 验证密码强度
	ValidatePasswordStrength(password string) error

	// ValidateEmailDomain 验证邮箱域名是否允许
	ValidateEmailDomain(email string) error

	// ValidateNicknameUniqueness 验证昵称唯一性（需要配合仓储）
	// 注意：这个方法签名只是示例，实际实现在应用层
	ValidateNicknameUniqueness(nickname string, excludeUserID string) error
}

// userValidationService 实现
type userValidationService struct{}

// NewUserValidationService 创建用户验证领域服务
func NewUserValidationService() UserValidationService {
	return &userValidationService{}
}

// ValidatePasswordStrength 验证密码强度
func (s *userValidationService) ValidatePasswordStrength(password string) error {
	pwd, err := valueobject.NewPassword(password)
	if err != nil {
		return err
	}

	strength := pwd.Strength()
	if strength == "弱" {
		return errors.New("密码强度过弱，建议使用大小写字母、数字和特殊字符的组合")
	}

	return nil
}

// ValidateEmailDomain 验证邮箱域名是否允许
func (s *userValidationService) ValidateEmailDomain(email string) error {
	e, err := valueobject.NewEmail(email)
	if err != nil {
		return err
	}

	domain := e.Domain()

	// 业务规则：黑名单域名
	blacklist := map[string]bool{
		"example.com":     true,
		"test.com":        true,
		"tempmail.com":    true,
		"throwaway.email": true,
	}

	if blacklist[domain] {
		return errors.New("不支持该邮箱域名")
	}

	return nil
}

// ValidateNicknameUniqueness 验证昵称唯一性
// 注意：这个方法需要访问仓储，实际应该在应用层实现
func (s *userValidationService) ValidateNicknameUniqueness(nickname string, excludeUserID string) error {
	// 这里只做格式验证
	_, err := valueobject.NewNickname(nickname)
	if err != nil {
		return err
	}

	// 唯一性检查需要在应用层通过仓储实现
	return nil
}
