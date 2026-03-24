package generators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// DAOGenerator DAO from database generator
type DAOGenerator struct {
	opts DAOOptions
}

// DAOOptions DAO from database generation options
type DAOOptions struct {
	OutputPath    string
	DSN           string
	ConfigFile    string
	WithUnitTest  bool
	FieldNullable bool
	Tables        []string
}

// NewDAOGenerator creates a new DAO from database generator
func NewDAOGenerator(opts DAOOptions) *DAOGenerator {
	if opts.OutputPath == "" {
		opts.OutputPath = "internal/infra/persistence/dao"
	}
	return &DAOGenerator{opts: opts}
}

// Generate generates DAO layer from database
func (g *DAOGenerator) Generate() error {
	fmt.Println("开始从数据库生成 DAO 层...")

	// 获取数据库配置
	fmt.Println("正在加载配置...")
	dbConfig, err := g.loadDatabaseConfig()
	if err != nil {
		log.Fatalf("加载数据库配置失败：%v", err)
	}

	// 构建 DSN
	dsn := g.opts.DSN
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Port,
			dbConfig.SSLMode)
	}

	// 确定输出路径
	outputPath := g.opts.OutputPath
	if !filepath.IsAbs(outputPath) {
		baseDir, _ := os.Getwd()
		outputPath = filepath.Join(baseDir, outputPath)
	}

	// 创建生成器
	genConfig := gen.Config{
		OutPath: outputPath,
		Mode:    gen.WithDefaultQuery | gen.WithoutContext,

		// 字段配置
		FieldNullable:     g.opts.FieldNullable,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,

		// 生成配置
		WithUnitTest: g.opts.WithUnitTest,
	}

	generator := gen.NewGenerator(genConfig)

	// 连接数据库
	fmt.Println("正在连接数据库...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败：%v\n", err)
		fmt.Println("将使用表名直接生成模型...")
	} else {
		fmt.Println("数据库连接成功")
		generator.UseDB(db)
	}

	// 生成模型和 DAO 接口
	fmt.Println("正在生成模型和 DAO 接口...")

	// 如果指定了特定表
	if len(g.opts.Tables) > 0 {
		return g.generateSpecificTables(generator, g.opts.Tables)
	}

	return g.generateAllTables(generator)
}

func (g *DAOGenerator) generateSpecificTables(generator *gen.Generator, tables []string) error {
	models := make([]interface{}, 0, len(tables))
	for _, table := range tables {
		models = append(models, generator.GenerateModel(table))
	}

	generator.ApplyBasic(models...)
	generator.Execute()

	fmt.Printf("已生成 %d 个表的 DAO 代码\n", len(tables))
	fmt.Println("生成位置:", generator.Config.OutPath)
	return nil
}

func (g *DAOGenerator) generateAllTables(generator *gen.Generator) error {
	// 按模块分组生成模型
	fmt.Println("正在生成 go-ddd-scaffold 核心业务模型...")

	// ============================================
	// 核心业务模型（基于实际 migrations）
	// ============================================
	coreModels := []interface{}{
		// 用户与认证
		generator.GenerateModel("users"),

		// 租户管理
		generator.GenerateModel("tenants"),
		generator.GenerateModel("tenant_members"),
		generator.GenerateModel("tenant_configs"),

		// RBAC 权限系统
		generator.GenerateModel("roles"),
		generator.GenerateModel("permissions"),
		generator.GenerateModel("role_permissions"),

		// 审计与日志
		generator.GenerateModel("audit_logs"),
		generator.GenerateModel("login_logs"),

		// DDD 基础设施
		generator.GenerateModel("domain_events"),
	}

	generator.ApplyBasic(coreModels...)

	// ============================================
	// 定义核心模型别名和关联关系
	// ============================================
	fmt.Println("正在配置模型关联关系...")

	// 基础模型
	userModel := generator.GenerateModelAs("users", "User")
	tenantModel := generator.GenerateModelAs("tenants", "Tenant")
	memberModel := generator.GenerateModelAs("tenant_members", "TenantMember")
	roleModel := generator.GenerateModelAs("roles", "Role")
	permissionModel := generator.GenerateModelAs("permissions", "Permission")
	rolePermissionModel := generator.GenerateModelAs("role_permissions", "RolePermission")

	// 配置核心关联关系
	relationshipModels := []interface{}{
		// 用户与租户关系
		generator.GenerateModelAs("users", "User",
			gen.FieldRelate(field.HasMany, "TenantMemberships", memberModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"UserID"}},
			}),
		),
		generator.GenerateModelAs("tenants", "Tenant",
			gen.FieldRelate(field.BelongsTo, "Owner", userModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"OwnerID"}},
			}),
			gen.FieldRelate(field.HasMany, "Members", memberModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}},
			}),
		),
		generator.GenerateModelAs("tenant_members", "TenantMember",
			gen.FieldRelate(field.BelongsTo, "Tenant", tenantModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}},
			}),
			gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"UserID"}},
			}),
			gen.FieldRelate(field.BelongsTo, "Role", roleModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"RoleID"}},
			}),
		),

		// RBAC 角色与权限关系
		generator.GenerateModelAs("roles", "Role",
			gen.FieldRelate(field.HasMany, "RolePermissions", rolePermissionModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"RoleID"}},
			}),
		),
		generator.GenerateModelAs("permissions", "Permission",
			gen.FieldRelate(field.HasMany, "RolePermissions", rolePermissionModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"PermissionID"}},
			}),
		),
		generator.GenerateModelAs("role_permissions", "RolePermission",
			gen.FieldRelate(field.BelongsTo, "Role", roleModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"RoleID"}},
			}),
			gen.FieldRelate(field.BelongsTo, "Permission", permissionModel, &field.RelateConfig{
				GORMTag: field.GormTag{"foreignKey": []string{"PermissionID"}},
			}),
		),

		// 其他表（无复杂关联）
		generator.GenerateModelAs("tenant_configs", "TenantConfig"),
		generator.GenerateModelAs("audit_logs", "AuditLog"),
		generator.GenerateModelAs("login_logs", "LoginLog"),
		generator.GenerateModelAs("domain_events", "DomainEvent"),
	}

	generator.ApplyBasic(relationshipModels...)

	// 执行生成
	fmt.Println("正在执行生成...")
	generator.Execute()
	fmt.Println("DAO 层生成完成!")
	fmt.Printf("生成位置：%s/{dao,model}\n", generator.Config.OutPath)

	return nil
}

func (g *DAOGenerator) loadDatabaseConfig() (*config.DatabaseConfig, error) {
	// 优先使用环境变量
	host := os.Getenv("APP_DATABASE_HOST")
	if host == "" {
		host = os.Getenv("DATABASE_HOST")
		if host == "" {
			host = "localhost"
		}
	}

	portStr := os.Getenv("APP_DATABASE_PORT")
	if portStr == "" {
		portStr = os.Getenv("DATABASE_PORT")
		if portStr == "" {
			portStr = "5432"
		}
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432
	}

	user := os.Getenv("APP_DATABASE_USER")
	if user == "" {
		user = os.Getenv("DATABASE_USER")
		if user == "" {
			user = "shenfay"
		}
	}

	password := os.Getenv("APP_DATABASE_PASSWORD")
	if password == "" {
		password = os.Getenv("DATABASE_PASSWORD")
		if password == "" {
			password = "postgres"
		}
	}

	dbname := os.Getenv("APP_DATABASE_NAME")
	if dbname == "" {
		dbname = os.Getenv("DATABASE_NAME")
		if dbname == "" {
			dbname = "go_ddd_scaffold"
		}
	}

	sslmode := os.Getenv("APP_DATABASE_SSL_MODE")
	if sslmode == "" {
		sslmode = os.Getenv("DATABASE_SSL_MODE")
		if sslmode == "" {
			sslmode = "disable"
		}
	}

	fmt.Printf("使用数据库配置：host=%s, port=%d, dbname=%s, user=%s\n",
		host, port, dbname, user)

	return &config.DatabaseConfig{
		Host:     host,
		Port:     port,
		Name:     dbname,
		User:     user,
		Password: password,
		SSLMode:  sslmode,
	}, nil
}
