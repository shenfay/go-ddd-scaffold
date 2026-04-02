package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shenfay/go-ddd-scaffold/internal/auth"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cli migrate <up|down>")
		os.Exit(1)
	}

	command := os.Args[1]
	action := os.Args[2]

	if command != "migrate" {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	// 加载配置
	cfg, err := config.Load("development")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch action {
	case "up":
		if err := runMigrationsUp(db); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		fmt.Println("✓ Database migrations completed successfully")
	case "down":
		if err := runMigrationsDown(db); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		fmt.Println("✓ Database rollback completed successfully")
	default:
		fmt.Printf("Unknown action: %s\n", action)
		os.Exit(1)
	}
}

// runMigrationsUp 执行正向迁移
func runMigrationsUp(db *gorm.DB) error {
	log.Println("Running migrations up...")

	// 自动迁移表结构
	if err := db.AutoMigrate(&auth.UserPO{}); err != nil {
		return fmt.Errorf("failed to migrate user table: %w", err)
	}

	log.Println("All migrations applied successfully")
	return nil
}

// runMigrationsDown 回滚迁移
func runMigrationsDown(db *gorm.DB) error {
	log.Println("Rolling back migrations...")

	// 注意：GORM 不支持自动回滚
	// 这里只是示例，实际项目中应该手动编写回滚 SQL
	fmt.Println("Warning: GORM does not support automatic rollback.")
	fmt.Println("You need to manually drop tables if needed:")
	fmt.Println("  DROP TABLE IF EXISTS user_pos CASCADE;")

	return nil
}
