package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConfig_DSN(t *testing.T) {
	t.Run("should generate correct DSN string", func(t *testing.T) {
		cfg := &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "secret",
			Name:     "testdb",
			SSLMode:  "disable",
		}

		dsn := cfg.DSN()
		assert.Contains(t, dsn, "host=localhost")
		assert.Contains(t, dsn, "port=5432")
		assert.Contains(t, dsn, "user=postgres")
		assert.Contains(t, dsn, "password=secret")
		assert.Contains(t, dsn, "dbname=testdb")
		assert.Contains(t, dsn, "sslmode=disable")
	})
}

func TestLoad_Success(t *testing.T) {
	// 保存当前工作目录
	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()

	t.Run("should load config from configs directory", func(t *testing.T) {
		// 切换到backend目录
		backendDir := filepath.Join(originalDir, "..")
		if err := os.Chdir(backendDir); err != nil {
			t.Skip("Cannot change to backend directory")
		}

		cfg, err := Load("test")
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		// 验证基本配置加载
		assert.NotEmpty(t, cfg.Server.Port)
		assert.NotEmpty(t, cfg.Database.Host)
	})
}

func TestLoad_WithEnvironment(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()

	t.Run("should respect environment variables", func(t *testing.T) {
		backendDir := filepath.Join(originalDir, "..")
		if err := os.Chdir(backendDir); err != nil {
			t.Skip("Cannot change to backend directory")
		}

		// 设置环境变量
		if err := os.Setenv("APP_SERVER_PORT", "9999"); err != nil {
			t.Skip("Cannot set environment variable")
		}
		defer func() { _ = os.Unsetenv("APP_SERVER_PORT") }()

		cfg, err := Load("test")
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}

func TestFindConfigDir(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalDir) }()

	t.Run("should find configs directory", func(t *testing.T) {
		backendDir := filepath.Join(originalDir, "..")
		if err := os.Chdir(backendDir); err != nil {
			t.Skip("Cannot change to backend directory")
		}

		configDir := findConfigDir()
		assert.NotEmpty(t, configDir)
		assert.Contains(t, configDir, "configs")
	})
}

func TestConfigStructs(t *testing.T) {
	t.Run("should create ServerConfig", func(t *testing.T) {
		cfg := ServerConfig{
			Port:         8080,
			Mode:         "debug",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
		}
		assert.Equal(t, 8080, cfg.Port)
		assert.Equal(t, "debug", cfg.Mode)
	})

	t.Run("should create RedisConfig", func(t *testing.T) {
		cfg := RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			PoolSize: 10,
		}
		assert.Equal(t, "localhost:6379", cfg.Addr)
		assert.Equal(t, 10, cfg.PoolSize)
	})

	t.Run("should create JWTConfig", func(t *testing.T) {
		cfg := JWTConfig{
			Secret:        "test-secret",
			AccessExpire:  3600,
			RefreshExpire: 86400,
			Issuer:        "test-issuer",
		}
		assert.Equal(t, "test-secret", cfg.Secret)
		assert.Equal(t, "test-issuer", cfg.Issuer)
	})

	t.Run("should create CORSConfig", func(t *testing.T) {
		cfg := CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST"},
			AllowedHeaders:   []string{"Authorization"},
			AllowCredentials: true,
			MaxAge:           3600,
		}
		assert.True(t, cfg.AllowCredentials)
		assert.Len(t, cfg.AllowedOrigins, 1)
	})

	t.Run("should create RateLimitConfig", func(t *testing.T) {
		cfg := RateLimitConfig{
			Enabled: true,
			General: RateLimitRule{
				Rate:  10,
				Burst: 20,
			},
		}
		assert.True(t, cfg.Enabled)
		assert.Equal(t, 10.0, cfg.General.Rate)
	})

	t.Run("should create DeviceConfig", func(t *testing.T) {
		cfg := DeviceConfig{
			MaxDevicesPerUser: 5,
			AutoRevokeOldest:  true,
		}
		assert.Equal(t, 5, cfg.MaxDevicesPerUser)
		assert.True(t, cfg.AutoRevokeOldest)
	})

	t.Run("should create MetricsConfig", func(t *testing.T) {
		cfg := MetricsConfig{
			Enabled:  true,
			HTTP:     MetricsHTTPConfig{Enabled: true},
			Database: MetricsDatabaseConfig{Enabled: true},
			Redis:    MetricsRedisConfig{Enabled: true},
		}
		assert.True(t, cfg.Enabled)
		assert.True(t, cfg.HTTP.Enabled)
	})
}
