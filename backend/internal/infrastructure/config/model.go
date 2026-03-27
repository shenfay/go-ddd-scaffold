package config

import (
	"fmt"
	"time"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string        `mapstructure:"port" validate:"required"`
	Mode         string        `mapstructure:"mode" validate:"oneof=debug release test"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host" validate:"required"`
	Port            int           `mapstructure:"port" validate:"required,min=1,max=65535"`
	Name            string        `mapstructure:"name" validate:"required"`
	User            string        `mapstructure:"user" validate:"required"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr" validate:"required"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret" validate:"required"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `mapstructure:"level" validate:"oneof=debug info warn error"`
	Format     string `mapstructure:"format" validate:"oneof=console json"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`    // 单个日志文件最大大小 (MB)
	MaxBackups int    `mapstructure:"max_backups"` // 保留的旧日志文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 日志文件保留天数
}

// SnowflakeConfig Snowflake ID 生成器配置
type SnowflakeConfig struct {
	NodeID int64 `mapstructure:"node_id" validate:"min=0,max=1023"`
}

// PasswordPolicyConfig 密码策略配置
type PasswordPolicyConfig struct {
	MinLength           int    `mapstructure:"min_length" validate:"min=6,max=128"`
	MaxLength           int    `mapstructure:"max_length" validate:"min=6,max=128"`
	RequireUppercase    bool   `mapstructure:"require_uppercase"`
	RequireLowercase    bool   `mapstructure:"require_lowercase"`
	RequireDigits       bool   `mapstructure:"require_digits"`
	RequireSpecialChars bool   `mapstructure:"require_special_chars"`
	SpecialChars        string `mapstructure:"special_chars"`
	DisallowCommon      bool   `mapstructure:"disallow_common"`
}

// PasswordHasherConfig 密码哈希配置
type PasswordHasherConfig struct {
	Cost int `mapstructure:"cost" validate:"min=4,max=31"`
}

// SecurityConfig 安全配置（包含密码相关）
type SecurityConfig struct {
	PasswordPolicy PasswordPolicyConfig `mapstructure:"password_policy"`
	PasswordHasher PasswordHasherConfig `mapstructure:"password_hasher"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTPHost  string `mapstructure:"smtp_host"`
	SMTPPort  int    `mapstructure:"smtp_port"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"` // 授权码
	From      string `mapstructure:"from"`
	FromName  string `mapstructure:"from_name"`
	EnableTLS bool   `mapstructure:"enable_tls"`
}

// EventsConfig 事件系统配置
type EventsConfig struct {
	Asynq         AsynqConfig         `mapstructure:"asynq"`
	Deduplication DeduplicationConfig `mapstructure:"deduplication"`
	Monitoring    MonitoringConfig    `mapstructure:"monitoring"`
	Alerts        AlertsConfig        `mapstructure:"alerts"`
}

// AsynqConfig Asynq 配置
type AsynqConfig struct {
	Concurrency int            `mapstructure:"concurrency"`
	Queues      map[string]int `mapstructure:"queues"`
	Retry       RetryConfig    `mapstructure:"retry"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts int           `mapstructure:"max_attempts"`
	MinBackoff  time.Duration `mapstructure:"min_backoff"`
	MaxBackoff  time.Duration `mapstructure:"max_backoff"`
}

// DeduplicationConfig 去重配置
type DeduplicationConfig struct {
	TTL time.Duration `mapstructure:"ttl"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	EnableMetrics  bool          `mapstructure:"enable_metrics"`
	ScrapeInterval time.Duration `mapstructure:"scrape_interval"`
	EnableAlerts   bool          `mapstructure:"enable_alerts"`
}

// AlertsConfig 告警配置
type AlertsConfig struct {
	Channels []AlertChannel `mapstructure:"channels"`
	Rules    []AlertRule    `mapstructure:"rules"`
}

// AlertChannel 告警渠道
type AlertChannel struct {
	Type    string `mapstructure:"type"`
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url,omitempty"`
}

// AlertRule 告警规则
type AlertRule struct {
	Name      string  `mapstructure:"name"`
	Threshold float64 `mapstructure:"threshold"`
	Severity  string  `mapstructure:"severity"`
}

// AppConfig 应用完整配置
type AppConfig struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
	Security  SecurityConfig  `mapstructure:"security"`
	Email     EmailConfig     `mapstructure:"email"`
	Events    EventsConfig    `mapstructure:"events"`
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	sslMode := "disable"
	if c.SSLMode != "" {
		sslMode = c.SSLMode
	}

	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" +
		fmt.Sprintf("%d", c.Port) + "/" + c.Name + "?sslmode=" + sslMode
}

// GetSnowflakeNodeID 获取 Snowflake 节点 ID
func (c *AppConfig) GetSnowflakeNodeID() int64 {
	if c.Snowflake.NodeID >= 0 && c.Snowflake.NodeID <= 1023 {
		return c.Snowflake.NodeID
	}
	return 0 // 默认返回 0
}
