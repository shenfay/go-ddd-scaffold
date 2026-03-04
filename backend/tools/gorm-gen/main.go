package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("开始生成DAO层...")

	// 从环境变量或配置文件读取数据库连接信息
	dsn := getDSN()

	// 按照官方文档配置生成器
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../internal/infrastructure/persistence/gorm/dao",
		Mode:    gen.WithDefaultQuery | gen.WithoutContext,

		// 字段配置
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,

		// 生成配置
		WithUnitTest: false, // 是否生成单元测试
	})

	// 尝试连接数据库
	fmt.Println("正在连接数据库...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		fmt.Println("将使用表名直接生成模型...")
	} else {
		fmt.Println("数据库连接成功")
		g.UseDB(db)
	}

	// 基于表名生成模型和 DAO 接口
	fmt.Println("正在生成模型和 DAO 接口...")
	// 按模块分组生成模型
	// 先生成基础模型
	g.ApplyBasic(
		// User Management 相关模型
		g.GenerateModel("users"),
		g.GenerateModel("tenants"),
		g.GenerateModel("tenant_members"),
		g.GenerateModel("tenant_invitations"),

		// Casbin RBAC 相关模型
		g.GenerateModel("casbin_rule"),
	)

	// 定义模型别名
	userModel := g.GenerateModelAs("users", "User")
	tenantModel := g.GenerateModelAs("tenants", "Tenant")
	memberModel := g.GenerateModelAs("tenant_members", "TenantMember")
	invitationModel := g.GenerateModelAs("tenant_invitations", "TenantInvitation")

	// 配置关联关系 - 用户管理模块
	g.ApplyBasic(
		g.GenerateModelAs("users", "User",
			gen.FieldRelate(field.HasMany, "TenantMemberships", memberModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"UserID"}}}),
			gen.FieldRelate(field.HasMany, "CreatedInvitations", invitationModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"CreatorID"}}}),
			gen.FieldRelate(field.HasMany, "ReceivedInvitations", memberModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"InvitedBy"}}}),
		),
		g.GenerateModelAs("tenants", "Tenant",
			gen.FieldRelate(field.HasMany, "Members", memberModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}}}),
			gen.FieldRelate(field.HasMany, "Invitations", invitationModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}}}),
		),
		g.GenerateModelAs("tenant_members", "TenantMember",
			gen.FieldRelate(field.BelongsTo, "Tenant", tenantModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}}}),
			gen.FieldRelate(field.BelongsTo, "User", userModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"UserID"}}}),
			gen.FieldRelate(field.BelongsTo, "Inviter", userModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"InvitedBy"}}}),
		),
		g.GenerateModelAs("tenant_invitations", "TenantInvitation",
			gen.FieldRelate(field.BelongsTo, "Tenant", tenantModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TenantID"}}}),
			gen.FieldRelate(field.BelongsTo, "Creator", userModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"CreatorID"}}}),
		),
	)

	// Casbin RBAC 相关模型（无需关联关系）
	g.ApplyBasic(
		g.GenerateModelAs("casbin_rule", "CasbinRule"),
	)

	// 执行生成
	fmt.Println("正在执行生成...")
	g.Execute()
	fmt.Println("DAO层生成完成!")
	fmt.Println("生成位置: ../../internal/infrastructure/persistence/gorm/{dao,model}")
}

// getDSN 从环境变量获取数据库连接字符串
func getDSN() string {
	// 优先从环境变量读取
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "go_ddd_scaffold"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "math111"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "go_ddd_scaffold"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		host, user, password, dbname, port, sslmode)
}
