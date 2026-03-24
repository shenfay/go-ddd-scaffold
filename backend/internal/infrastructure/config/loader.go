package config

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// ConfigLoader 配置加载器
type ConfigLoader struct {
	viper    *viper.Viper
	validate *validator.Validate
	logger   *zap.Logger
}

// NewConfigLoader 创建新的配置加载器
func NewConfigLoader(logger *zap.Logger) *ConfigLoader {
	v := viper.New()

	// 基础配置
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 环境变量配置
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 如果 logger 为 nil，创建一个简单的 logger
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}

	return &ConfigLoader{
		viper:    v,
		validate: validator.New(),
		logger:   logger,
	}
}

// Load 加载配置
func (c *ConfigLoader) Load(env string) (*AppConfig, error) {
	// 如果指定了环境，尝试加载环境特定配置
	if env != "" {
		envConfigName := fmt.Sprintf("config_%s", env)
		c.viper.SetConfigName(envConfigName)

		if err := c.viper.MergeInConfig(); err != nil {
			// 环境配置文件不存在时不报错，继续使用主配置
			c.logger.Warn("failed to load environment config",
				zap.String("env", env),
				zap.String("config_name", envConfigName),
				zap.Error(err))
		} else {
			c.logger.Info("loaded environment config",
				zap.String("env", env),
				zap.String("config_name", envConfigName))
		}

		// 重新设置为主配置文件名
		c.viper.SetConfigName("config")
	}

	// 加载主配置文件
	if err := c.viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	c.logger.Info("loaded main config file",
		zap.String("config_file", c.viper.ConfigFileUsed()))

	// 反序列化到结构体
	var config AppConfig
	if err := c.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := c.validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	c.logger.Info("config loaded and validated successfully")
	return &config, nil
}

// WatchConfig 监听配置文件变化
func (c *ConfigLoader) WatchConfig(callback func(*AppConfig)) {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		c.logger.Info("config file changed", zap.String("file", e.Name))

		var config AppConfig
		if err := c.viper.Unmarshal(&config); err != nil {
			c.logger.Error("failed to reload config", zap.Error(err))
			return
		}

		if err := c.validate.Struct(config); err != nil {
			c.logger.Error("reloaded config validation failed", zap.Error(err))
			return
		}

		c.logger.Info("config reloaded successfully")
		callback(&config)
	})
}

// GetString 获取字符串配置值
func (c *ConfigLoader) GetString(key string) string {
	return c.viper.GetString(key)
}

// GetInt 获取整数配置值
func (c *ConfigLoader) GetInt(key string) int {
	return c.viper.GetInt(key)
}

// GetBool 获取布尔配置值
func (c *ConfigLoader) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

// GetDuration 获取时间间隔配置值
func (c *ConfigLoader) GetDuration(key string) string {
	return c.viper.GetString(key)
}

// GetAllSettings 获取所有配置设置
func (c *ConfigLoader) GetAllSettings() map[string]interface{} {
	return c.viper.AllSettings()
}
