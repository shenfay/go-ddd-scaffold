package auth

import (
	"fmt"

	"github.com/casbin/casbin/v3"
)

// InitDefaultPolicies 初始化默认权限策略
// 在服务启动时调用，将预置的租户级权限加载到系统中
func InitDefaultPolicies(e *casbin.Enforcer) error {
	// 租户级默认角色权限模板
	// 格式: 角色, 租户ID(使用占位符), 资源, 操作
	defaultPolicies := [][]string{
		// 成员角色（原 parent 角色）
		{"member", ":tenant", "users", "read"},
		{"member", ":tenant", "users", "write"},
		{"member", ":tenant", "users", "delete"},
		{"member", ":tenant", "self", "read"},
		{"member", ":tenant", "self", "write"},
		{"member", ":tenant", "invitation", "manage"},
		{"member", ":tenant", "settings", "read"},
		{"member", ":tenant", "settings", "write"},

		// 访客角色（原 child 角色）
		{"guest", ":tenant", "self", "read"},
		{"guest", ":tenant", "self", "write"},

		// 爷爷角色
		{"grandpa", ":tenant", "children", "read"},
		{"grandpa", ":tenant", "progress", "read"},
		{"grandpa", ":tenant", "self", "read"},
		{"grandpa", ":tenant", "self", "write"},

		// 奶奶角色
		{"grandma", ":tenant", "children", "read"},
		{"grandma", ":tenant", "progress", "read"},
		{"grandma", ":tenant", "self", "read"},
		{"grandma", ":tenant", "self", "write"},
	}

	// 获取当前已存在的策略
	existingPolicies, _ := e.GetPolicy()

	// 添加默认策略
	for _, policy := range defaultPolicies {
		// 检查是否已存在
		exists := false
		for _, p := range existingPolicies {
			if len(p) >= 4 && p[0] == policy[0] && p[1] == policy[1] && p[2] == policy[2] && p[3] == policy[3] {
				exists = true
				break
			}
		}
		if !exists {
			// 转换为 []interface{} 类型
			params := make([]interface{}, len(policy))
			for i, v := range policy {
				params[i] = v
			}
			_, err := e.AddPolicy(params...)
			if err != nil {
				return fmt.Errorf("添加默认策略失败: %v", err)
			}
		}
	}

	// 保存策略到数据库
	if err := e.SavePolicy(); err != nil {
		return fmt.Errorf("保存默认策略失败: %v", err)
	}

	return nil
}

// SyncUserRoles 同步用户角色到 Casbin
// 在用户登录成功后调用，将用户的租户角色同步到 Casbin
func SyncUserRoles(e *casbin.Enforcer, userID, tenantID, role string) error {
	// 先获取该用户在租户内的所有角色，然后逐一删除
	roles := e.GetRolesForUserInDomain(userID, tenantID)
	for _, r := range roles {
		_, _ = e.DeleteRoleForUserInDomain(userID, r, tenantID)
	}

	// 添加新角色
	_, err := e.AddRoleForUserInDomain(userID, role, tenantID)
	if err != nil {
		return fmt.Errorf("同步用户角色失败: %v", err)
	}

	// 保存到数据库
	if err := e.SavePolicy(); err != nil {
		return fmt.Errorf("保存用户角色失败: %v", err)
	}

	return nil
}

// FamilyRoleToSystemRole 家庭角色转换为系统角色标识
// 用于 Casbin 中的角色标识
// 注意：此函数保留以支持向后兼容，新系统应直接使用 RoleMember/RoleGuest
func FamilyRoleToSystemRole(familyRole string) string {
	roleMap := map[string]string{
		"father":  "member",
		"mother":  "member",
		"child":   "guest",
		"grandpa": "member",
		"grandma": "member",
		"parent":  "member", // 通用父母角色
	}

	if role, ok := roleMap[familyRole]; ok {
		return role
	}
	return familyRole
}
