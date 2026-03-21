package valueobject

// TenantConfig 租户配置值对象
type TenantConfig struct {
	MaxStorageGB      int               `json:"max_storage_gb"`
	MaxProjects       int               `json:"max_projects"`
	AllowedFeatures   []string          `json:"allowed_features"`
	CustomSettings    map[string]string `json:"custom_settings"`
	RequireMFA        bool              `json:"require_mfa"`
	SessionTimeoutMin int               `json:"session_timeout_min"`
}

// NewDefaultTenantConfig 创建默认租户配置
func NewDefaultTenantConfig() *TenantConfig {
	return &TenantConfig{
		MaxStorageGB:      10,
		MaxProjects:       10,
		AllowedFeatures:   []string{"basic", "api_access"},
		CustomSettings:    make(map[string]string),
		RequireMFA:        false,
		SessionTimeoutMin: 30,
	}
}
