package auth

import (
	"path/filepath"

	"gorm.io/gorm"

	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

const (
	// AuthConfigPath 权限模型配置文件路径
	AuthConfigPath = "config/auth/rbac_with_domains.conf"
	// CasbinTableName Casbin 策略存储表名
	CasbinTableName = "casbin_rule"
)

// NewCasbinAdapter 创建 Casbin Gorm 适配器
// 该适配器会自动创建 casbin_rule 表（如果不存在）
func NewCasbinAdapter(db *gorm.DB) (*gormadapter.Adapter, error) {
	// 创建适配器，使用指定的表名
	// 参数: db, prefix(空), tableName
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", CasbinTableName)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

// NewCasbinEnforcer 创建 Casbin 权限执行器
func NewCasbinEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	// 创建适配器
	adapter, err := NewCasbinAdapter(db)
	if err != nil {
		return nil, err
	}

	// 获取配置文件的绝对路径
	modelPath := filepath.Join("config", "auth", "rbac_with_domains.conf")

	// 创建 Enforcer
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}

	// 从数据库加载策略
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}

	return e, nil
}
