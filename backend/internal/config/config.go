package config

import (
	"time"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	LLM      LLMConfig      `mapstructure:"llm"`
	Log      LogConfig      `mapstructure:"log"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime string `mapstructure:"conn_max_idle_time"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host           string              `mapstructure:"host"`
	Port           int                 `mapstructure:"port"`
	Password       string              `mapstructure:"password"`
	DB             int                 `mapstructure:"db"`
	PoolSize       int                 `mapstructure:"pool_size"`
	MinIdleConns   int                 `mapstructure:"min_idle_conns"`
	EventBusConfig EventBusConfig      `mapstructure:"event_bus"`
}

// EventBusConfig 事件总线配置
type EventBusConfig struct {
	StreamKey      string        `mapstructure:"stream_key"`
	MaxRetries     int           `mapstructure:"max_retries"`
	RetryBaseDelay time.Duration `mapstructure:"retry_base_delay"`
	PollInterval   time.Duration `mapstructure:"poll_interval"`
	BatchSize      int           `mapstructure:"batch_size"`
}

type JWTConfig struct {
	SecretKey string        `mapstructure:"secret_key"`
	ExpireIn  time.Duration `mapstructure:"expire_in"`
}

type LLMConfig struct {
	Provider string `mapstructure:"provider"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"`
	Model    string `mapstructure:"model"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}
