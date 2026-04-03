package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shenfay/go-ddd-scaffold/internal/activitylog"
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

	// 1. 自动迁移表结构
	if err := db.AutoMigrate(&auth.UserPO{}); err != nil {
		return fmt.Errorf("failed to migrate user table: %w", err)
	}

	// 2. 迁移活动日志表
	if err := db.AutoMigrate(&activitylog.ActivityLog{}); err != nil {
		return fmt.Errorf("failed to migrate activity_logs table: %w", err)
	}

	// 3. 执行 SQL 迁移文件
	if err := executeSQLMigrations(db, "up"); err != nil {
		return fmt.Errorf("failed to execute SQL migrations: %w", err)
	}

	log.Println("All migrations applied successfully")
	return nil
}

// runMigrationsDown 回滚迁移
func runMigrationsDown(db *gorm.DB) error {
	log.Println("Rolling back migrations...")

	// 执行 SQL 回滚文件（按倒序）
	if err := executeSQLMigrations(db, "down"); err != nil {
		return fmt.Errorf("failed to execute SQL rollback: %w", err)
	}

	// 注意：GORM 不支持自动回滚
	// 这里只是示例，实际项目中应该手动编写回滚 SQL
	fmt.Println("Warning: GORM does not support automatic rollback.")
	fmt.Println("You need to manually drop tables if needed:")
	fmt.Println("  DROP TABLE IF EXISTS activity_logs CASCADE;")
	fmt.Println("  DROP TABLE IF EXISTS users CASCADE;")

	return nil
}

// executeSQLMigrations 执行 SQL 迁移文件
func executeSQLMigrations(db *gorm.DB, direction string) error {
	migrationsDir := "migrations"
	pattern := filepath.Join(migrationsDir, fmt.Sprintf("*.%s.sql", direction))

	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	if len(files) == 0 {
		log.Printf("No SQL migration files found for direction: %s", direction)
		return nil
	}

	// 如果是正向迁移，按顺序执行；如果是回滚，按倒序执行
	if direction == "down" {
		sort.Sort(sort.Reverse(sort.StringSlice(files)))
	}

	log.Printf("Found %d SQL migration files to execute", len(files))

	for _, file := range files {
		filename := filepath.Base(file)
		log.Printf("Executing migration: %s", filename)

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		sql := strings.TrimSpace(string(sqlBytes))
		if sql == "" {
			log.Printf("Skipping empty migration file: %s", filename)
			continue
		}

		// 执行 SQL
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		log.Printf("✓ Executed: %s", filename)
	}

	return nil
}
