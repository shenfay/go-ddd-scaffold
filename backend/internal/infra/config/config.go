package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Logger    LoggerConfig    `mapstructure:"logger"`
	CORS      CORSConfig      `mapstructure:"cors"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
	Device    DeviceConfig    `mapstructure:"device"`
	Token     TokenConfig     `mapstructure:"token"`
	Email     EmailConfig     `mapstructure:"email"`
	Asynq     AsynqConfig     `mapstructure:"asynq"`
	Metrics   MetricsConfig   `mapstructure:"metrics"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Name            string        `mapstructure:"name"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// RedisConfig Redis 连接配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string        `mapstructure:"secret"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	OutputPath string `mapstructure:"output_path"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

// RateLimitConfig 速率限制配置
type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	General  RateLimitRule `mapstructure:"general"`
	Login    RateLimitRule `mapstructure:"login"`
	Register RateLimitRule `mapstructure:"register"`
}

// RateLimitRule 速率限制规则
type RateLimitRule struct {
	Rate  float64 `mapstructure:"rate"`
	Burst int     `mapstructure:"burst"`
}

// DeviceConfig 设备管理配置
type DeviceConfig struct {
	MaxDevicesPerUser int  `mapstructure:"max_devices_per_user"` // 每个用户最多登录设备数
	AutoRevokeOldest  bool `mapstructure:"auto_revoke_oldest"`   // 是否自动踢出最旧设备
}

// TokenConfig Token 过期配置
type TokenConfig struct {
	EmailVerificationExpire time.Duration `mapstructure:"email_verification_expire"` // 邮箱验证 token 过期时间
	PasswordResetExpire     time.Duration `mapstructure:"password_reset_expire"`     // 密码重置 token 过期时间
}

// EmailConfig 邮件配置
type EmailConfig struct {
	From                     string `mapstructure:"from"`
	VerificationURLTemplate  string `mapstructure:"verification_url_template"`
	PasswordResetURLTemplate string `mapstructure:"password_reset_url_template"`
}

// AsynqConfig Asynq 消息队列配置
type AsynqConfig struct {
	Addr        string         `mapstructure:"addr"`
	Concurrency int            `mapstructure:"concurrency"`
	Queues      map[string]int `mapstructure:"queues"`
}

// MetricsConfig 监控指标配置
type MetricsConfig struct {
	Enabled  bool                  `mapstructure:"enabled"`
	HTTP     MetricsHTTPConfig     `mapstructure:"http"`
	Database MetricsDatabaseConfig `mapstructure:"database"`
	Redis    MetricsRedisConfig    `mapstructure:"redis"`
}

// MetricsHTTPConfig HTTP 指标配置
type MetricsHTTPConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsDatabaseConfig 数据库指标配置
type MetricsDatabaseConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// MetricsRedisConfig Redis 指标配置
type MetricsRedisConfig struct {
	Enabled bool `mapstructure:"enabled"`
}

// Load 加载配置
func Load(env string) (*Config, error) {
	// 清理 viper 状态
	viper.Reset()

	// 确定配置目录
	configDir := findConfigDir()
	if configDir == "" {
		return nil, fmt.Errorf("config directory not found")
	}

	viper.AddConfigPath(configDir)

	// 1. 绑定环境变量（优先级最高）
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 2. 加载所有 YAML 文件（使用独立 viper 实例）
	yamlFiles := []string{
		"server",
		"database",
		"redis",
		"auth",
		"rate_limit",
		"app",
	}

	mergedConfig := make(map[string]interface{})

	for _, filename := range yamlFiles {
		v := viper.New()
		v.AddConfigPath(configDir)
		v.SetConfigName(filename)
		v.SetConfigType("yaml")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				continue // 文件不存在，跳过
			}
			return nil, fmt.Errorf("failed to read %s.yaml: %w", filename, err)
		}

		// 读取配置并合并
		var fileData map[string]interface{}
		if err := v.Unmarshal(&fileData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s.yaml: %w", filename, err)
		}

		// 合并到总配置
		for key, value := range fileData {
			mergedConfig[key] = value
		}

		// 调试输出
		if filename == "database" {
		}
	}

	// 3. 加载 .env 文件（覆盖 YAML）
	envViper := viper.New()
	envViper.AddConfigPath(configDir)
	envViper.SetConfigName(".env")
	envViper.SetConfigType("env")

	if err := envViper.ReadInConfig(); err == nil {
		var envData map[string]interface{}
		if err := envViper.Unmarshal(&envData); err == nil {
			for key, value := range envData {
				mergedConfig[key] = value
			}
		}
	}

	// 5. 将所有值设置到 viper（覆盖默认值）
	for key, value := range mergedConfig {
		viper.Set(key, value)
	}

	// 6. 反序列化到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// findConfigDir 查找配置目录
func findConfigDir() string {
	paths := []string{
		"configs",
		"../configs",
		"../../configs",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			// 返回相对路径即可，viper 会正确处理
			return path
		}
	}

	return "configs"
}

// setDefaults 设置默认值（不会覆盖已存在的值）
func setDefaults() {
	// 只在配置不存在时设置默认值
	if !viper.IsSet("server.port") {
		viper.Set("server.port", 8080)
	}
	if !viper.IsSet("server.mode") {
		viper.Set("server.mode", "debug")
	}
	if !viper.IsSet("database.host") {
		viper.Set("database.host", "localhost")
	}
	if !viper.IsSet("database.name") {
		viper.Set("database.name", "ddd_scaffold")
	}
	if !viper.IsSet("database.user") {
		viper.Set("database.user", "postgres")
	}

}
