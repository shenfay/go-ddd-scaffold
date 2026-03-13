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

// AppConfig 应用完整配置
type AppConfig struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	Snowflake SnowflakeConfig `mapstructure:"snowflake"`
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
