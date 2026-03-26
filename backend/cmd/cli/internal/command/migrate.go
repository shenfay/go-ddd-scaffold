package command

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
		Long:  `Manage database migrations using golang-migrate`,
	}

	cmd.AddCommand(migrateUpCmd())
	cmd.AddCommand(migrateDownCmd())
	cmd.AddCommand(migrateCreateCmd())
	cmd.AddCommand(migrateStatusCmd())
	cmd.AddCommand(migrateVersionCmd())

	// 全局标志
	cmd.PersistentFlags().String("path", "./migrations", "迁移文件目录")
	cmd.PersistentFlags().String("dsn", "", "数据库连接 DSN（优先于配置文件）")

	return cmd
}

// migrateUpCmd 应用所有待处理的迁移
func migrateUpCmd() *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Long:  `Apply all pending database migrations. Use --steps to limit the number of migrations to apply.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			dsn, _ := cmd.Flags().GetString("dsn")

			if dsn == "" {
				var err error
				dsn, err = getDatabaseDSN()
				if err != nil {
					return fmt.Errorf("failed to get database DSN: %w", err)
				}
			}

			return runMigrateUp(dsn, path, steps)
		},
	}

	cmd.Flags().IntVarP(&steps, "steps", "s", 0, "Number of migrations to apply")
	return cmd
}

// migrateDownCmd 回滚迁移
func migrateDownCmd() *cobra.Command {
	var steps int

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Rollback database migrations",
		Long:  `Rollback database migrations. Default is 1 step, use --steps to specify.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			dsn, _ := cmd.Flags().GetString("dsn")

			if dsn == "" {
				var err error
				dsn, err = getDatabaseDSN()
				if err != nil {
					return fmt.Errorf("failed to get database DSN: %w", err)
				}
			}

			return runMigrateDown(dsn, path, steps)
		},
	}

	cmd.Flags().IntVarP(&steps, "steps", "s", 1, "Number of migrations to rollback")
	return cmd
}

// migrateCreateCmd 创建新的迁移文件
func migrateCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration",
		Long:  `Create a new migration file with timestamp`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			name := args[0]

			return runMigrateCreate(path, name)
		},
	}

	return cmd
}

// migrateStatusCmd 查看迁移状态
func migrateStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Long:  `Show current database version and pending migrations`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			dsn, _ := cmd.Flags().GetString("dsn")

			if dsn == "" {
				var err error
				dsn, err = getDatabaseDSN()
				if err != nil {
					return fmt.Errorf("failed to get database DSN: %w", err)
				}
			}

			return showMigrationStatus(dsn, path)
		},
	}

	return cmd
}

// migrateVersionCmd 查看当前版本
func migrateVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show current version",
		Long:  `Show current database migration version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			dsn, _ := cmd.Flags().GetString("dsn")

			if dsn == "" {
				var err error
				dsn, err = getDatabaseDSN()
				if err != nil {
					return fmt.Errorf("failed to get database DSN: %w", err)
				}
			}

			return showCurrentVersion(dsn, path)
		},
	}

	return cmd
}

// runMigrateUp 执行向上迁移
func runMigrateUp(dsn, path string, steps int) error {
	fmt.Println("正在应用数据库迁移...")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✓ 数据库连接成功")

	// 确保 schema_migrations 表存在
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	// 获取当前版本
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Printf("当前版本：%d\n", currentVersion)

	// 获取待应用的迁移文件
	pendingMigrations, err := getPendingMigrations(path, currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pendingMigrations) == 0 {
		fmt.Println("✓ 没有待处理的迁移")
		return nil
	}

	// 限制步骤数
	if steps > 0 && steps < len(pendingMigrations) {
		pendingMigrations = pendingMigrations[:steps]
	}

	fmt.Printf("发现 %d 个待应用的迁移\n", len(pendingMigrations))

	// 应用每个迁移
	for _, migrationFile := range pendingMigrations {
		version := extractVersion(migrationFile)
		direction := "up"

		fmt.Printf("\n应用迁移：%s\n", migrationFile)

		// 读取 SQL 文件
		sqlContent, err := readMigrationFile(path, migrationFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		// 执行迁移
		if err := executeMigration(db, sqlContent, fmt.Sprintf("%d", version), direction); err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}

		fmt.Printf("✓ 迁移 %s 应用成功\n", migrationFile)
	}

	fmt.Println("\n✓ 所有迁移应用成功！")
	return nil
}

// runMigrateDown 执行向下迁移（回滚）
func runMigrateDown(dsn, path string, steps int) error {
	fmt.Println("正在回滚数据库迁移...")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✓ 数据库连接成功")

	// 获取当前版本
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if currentVersion == 0 {
		fmt.Println("✓ 已经是最低版本，无需回滚")
		return nil
	}

	fmt.Printf("当前版本：%d\n", currentVersion)

	// 获取要回滚的迁移
	migrationsToRollback, err := getMigrationsToRollback(path, currentVersion, steps)
	if err != nil {
		return fmt.Errorf("failed to get migrations to rollback: %w", err)
	}

	if len(migrationsToRollback) == 0 {
		fmt.Println("✓ 没有可回滚的迁移")
		return nil
	}

	fmt.Printf("准备回滚 %d 个迁移\n", len(migrationsToRollback))

	// 回滚每个迁移
	for _, migrationFile := range migrationsToRollback {
		version := extractVersion(migrationFile)
		direction := "down"

		// 查找对应的 down 文件
		downFile := strings.Replace(migrationFile, ".up.sql", ".down.sql", 1)

		fmt.Printf("\n回滚迁移：%s\n", downFile)

		// 读取 SQL 文件
		sqlContent, err := readMigrationFile(path, downFile)
		if err != nil {
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		// 执行回滚
		if err := executeMigration(db, sqlContent, fmt.Sprintf("%d", version), direction); err != nil {
			return fmt.Errorf("failed to execute rollback: %w", err)
		}

		fmt.Printf("✓ 迁移 %s 回滚成功\n", downFile)
	}

	fmt.Println("\n✓ 回滚完成！")
	return nil
}

// runMigrateCreate 创建新的迁移文件
func runMigrateCreate(path, name string) error {
	fmt.Printf("创建新的迁移文件：%s\n", name)

	// 确保迁移目录存在
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// 生成时间戳版本号
	timestamp := time.Now().Unix()
	versionStr := fmt.Sprintf("%d", timestamp)

	// 简化名称（移除特殊字符）
	sanitizedName := sanitizeFilename(name)

	// 生成文件名
	upFilename := fmt.Sprintf("%s_%s.up.sql", versionStr, sanitizedName)
	downFilename := fmt.Sprintf("%s_%s.down.sql", versionStr, sanitizedName)

	upPath := filepath.Join(path, upFilename)
	downPath := filepath.Join(path, downFilename)

	// 创建 up 文件
	upContent := fmt.Sprintf(`-- +migrate Up
-- TODO: Add your migration SQL here

`)
	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	// 创建 down 文件
	downContent := fmt.Sprintf(`-- +migrate Down
-- TODO: Add your rollback SQL here

`)
	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	fmt.Printf("✓ 创建成功:\n")
	fmt.Printf("  Up:   %s\n", upPath)
	fmt.Printf("  Down: %s\n", downPath)
	fmt.Printf("\n提示：编辑这两个文件添加你的迁移逻辑\n")

	return nil
}

// showMigrationStatus 显示迁移状态
func showMigrationStatus(dsn, path string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 确保 migrations 表存在
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	// 获取当前版本
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Println("数据库迁移状态")
	fmt.Println("================")
	fmt.Printf("当前版本：%d\n", currentVersion)

	// 获取所有迁移文件
	allMigrations, err := getAllMigrationFiles(path)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	fmt.Printf("\n迁移文件列表 (共 %d 个):\n", len(allMigrations))
	fmt.Println("------------------------")

	for _, file := range allMigrations {
		version := extractVersion(file)
		status := "⬜ 待应用"
		if version <= currentVersion {
			status = "✅ 已应用"
		}
		fmt.Printf("[%s] %s (版本：%d)\n", status, file, version)
	}

	return nil
}

// showCurrentVersion 显示当前版本
func showCurrentVersion(dsn, path string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 确保 migrations 表存在
	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	fmt.Printf("当前数据库版本：%d\n", currentVersion)
	return nil
}

// Helper functions

func getDatabaseDSN() (string, error) {
	// 从环境变量读取
	host := os.Getenv("APP_DATABASE_HOST")
	if host == "" {
		host = os.Getenv("DB_HOST")
	}
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("APP_DATABASE_PORT")
	if port == "" {
		port = os.Getenv("DB_PORT")
	}
	if port == "" {
		port = "5432"
	}

	dbname := os.Getenv("APP_DATABASE_NAME")
	if dbname == "" {
		dbname = os.Getenv("DB_NAME")
	}
	if dbname == "" {
		dbname = "go_ddd_scaffold"
	}

	user := os.Getenv("APP_DATABASE_USER")
	if user == "" {
		user = os.Getenv("DB_USER")
	}
	if user == "" {
		user = "shenfay"
	}

	password := os.Getenv("APP_DATABASE_PASSWORD")
	if password == "" {
		password = os.Getenv("DB_PASSWORD")
	}
	if password == "" {
		password = "postgres"
	}

	sslmode := os.Getenv("APP_DATABASE_SSL_MODE")
	if sslmode == "" {
		sslmode = os.Getenv("DB_SSLMODE")
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbname, user, password, sslmode)

	return dsn, nil
}

func ensureMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version BIGINT PRIMARY KEY,
		dirty BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

func getCurrentVersion(db *sql.DB) (int64, error) {
	var version int64
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations WHERE dirty = false").Scan(&version)
	return version, err
}

func getPendingMigrations(path string, currentVersion int64) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(path, "*.up.sql"))
	if err != nil {
		return nil, err
	}

	var pending []string
	for _, file := range files {
		filename := filepath.Base(file)
		version := extractVersion(filename)
		if version > currentVersion {
			pending = append(pending, filename)
		}
	}

	// 排序
	sortStrings(pending)
	return pending, nil
}

func getMigrationsToRollback(path string, currentVersion int64, steps int) ([]string, error) {
	if steps <= 0 {
		steps = 1
	}

	files, err := filepath.Glob(filepath.Join(path, "*.up.sql"))
	if err != nil {
		return nil, err
	}

	var applied []string
	for _, file := range files {
		filename := filepath.Base(file)
		version := extractVersion(filename)
		if version <= currentVersion {
			applied = append(applied, filename)
		}
	}

	// 降序排序
	sortStringsDesc(applied)

	// 返回指定数量的迁移
	if steps > len(applied) {
		steps = len(applied)
	}

	return applied[:steps], nil
}

func getAllMigrationFiles(path string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(path, "*.up.sql"))
	if err != nil {
		return nil, err
	}

	var result []string
	for _, file := range files {
		result = append(result, filepath.Base(file))
	}

	sortStrings(result)
	return result, nil
}

func extractVersion(filename string) int64 {
	parts := strings.Split(filename, "_")
	if len(parts) > 0 {
		var version int64
		fmt.Sscanf(parts[0], "%d", &version)
		return version
	}
	return 0
}

func readMigrationFile(path, filename string) (string, error) {
	fullPath := filepath.Join(path, filename)
	content, err := os.ReadFile(fullPath)
	return string(content), err
}

func executeMigration(db *sql.DB, sqlContent, versionStr, direction string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 执行 SQL
	if _, err := tx.Exec(sqlContent); err != nil {
		return err
	}

	// 更新版本号
	if direction == "up" {
		// 插入或更新版本记录
		query := `
		INSERT INTO schema_migrations (version, dirty, created_at) 
		VALUES ($1, false, NOW())
		ON CONFLICT (version) 
		DO UPDATE SET dirty = false, created_at = NOW()`
		if _, err := tx.Exec(query, versionStr); err != nil {
			return err
		}
	} else {
		// 删除版本记录
		if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = $1", versionStr); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func sortStrings(slice []string) {
	// 简单冒泡排序
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			v1 := extractVersion(slice[i])
			v2 := extractVersion(slice[j])
			if v1 > v2 {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

func sortStringsDesc(slice []string) {
	for i := 0; i < len(slice)-1; i++ {
		for j := i + 1; j < len(slice); j++ {
			v1 := extractVersion(slice[i])
			v2 := extractVersion(slice[j])
			if v1 < v2 {
				slice[i], slice[j] = slice[j], slice[i]
			}
		}
	}
}

func sanitizeFilename(filename string) string {
	// 简单的文件名清理，移除特殊字符
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filename)
	return result
}
