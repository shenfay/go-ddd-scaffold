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

	// 基于表名生成模型和DAO接口
	fmt.Println("正在生成模型和DAO接口...")
	// 按模块分组生成模型
	// 先生成基础模型
	g.ApplyBasic(
		// Knowledge Graph 相关模型
		g.GenerateModel("kg_domains"),
		g.GenerateModel("kg_trunks"),
		g.GenerateModel("kg_nodes"),
		g.GenerateModel("kg_node_relationships"),
		g.GenerateModel("kg_competency_levels"),
		g.GenerateModel("kg_academic_concepts"),

		// Learning System 相关模型
		g.GenerateModel("learn_assessment_items"),
		g.GenerateModel("learn_side_quests"),
		g.GenerateModel("learn_student_progress"),

		// User Management 相关模型
		g.GenerateModel("users"),
		g.GenerateModel("tenants"),
		g.GenerateModel("tenant_members"),
		g.GenerateModel("tenant_invitations"),

		// NPC System 相关模型
		g.GenerateModel("npc_memories"),
		g.GenerateModel("npc_profiles"),
		g.GenerateModel("npc_relationships"),

		// Tagging System 相关模型
		g.GenerateModel("taggables"),
		g.GenerateModel("tags"),
	)

	// 定义模型别名
	kgDomainModel := g.GenerateModelAs("kg_domains", "KgDomain")
	kgTrunkModel := g.GenerateModelAs("kg_trunks", "KgTrunk")
	kgNodeModel := g.GenerateModelAs("kg_nodes", "KgNode")
	kgNodeRelModel := g.GenerateModelAs("kg_node_relationships", "KgNodeRelationship")

	userModel := g.GenerateModelAs("users", "User")
	tenantModel := g.GenerateModelAs("tenants", "Tenant")
	memberModel := g.GenerateModelAs("tenant_members", "TenantMember")
	invitationModel := g.GenerateModelAs("tenant_invitations", "TenantInvitation")

	// 配置关联关系 - 分批处理
	// 首先配置Knowledge Graph关联关系
	g.ApplyBasic(
		g.GenerateModelAs("kg_domains", "KgDomain",
			gen.FieldRelate(field.HasMany, "Trunks", kgTrunkModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"DomainID"}}}),
		),
		g.GenerateModelAs("kg_trunks", "KgTrunk",
			gen.FieldRelate(field.BelongsTo, "Domain", kgDomainModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"DomainID"}}}),
			gen.FieldRelate(field.HasMany, "Nodes", kgNodeModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TrunkID"}}}),
		),
		g.GenerateModelAs("kg_nodes", "KgNode",
			gen.FieldRelate(field.BelongsTo, "Trunk", kgTrunkModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TrunkID"}}}),
			gen.FieldRelate(field.HasMany, "NodeRelationshipsAsSource", kgNodeRelModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"FromNodeID"}}}),
			gen.FieldRelate(field.HasMany, "NodeRelationshipsAsTarget", kgNodeRelModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"ToNodeID"}}}),
		),
		g.GenerateModelAs("kg_node_relationships", "KgNodeRelationship",
			gen.FieldRelate(field.BelongsTo, "SourceNode", kgNodeModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"FromNodeID"}}}),
			gen.FieldRelate(field.BelongsTo, "TargetNode", kgNodeModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"ToNodeID"}}}),
		),
		g.GenerateModelAs("kg_competency_levels", "KgCompetencyLevel"),
		g.GenerateModelAs("kg_academic_concepts", "KgAcademicConcept"),
	)

	// 然后配置用户管理关联关系
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

	// 接着配置学习系统关联关系
	g.ApplyBasic(
		g.GenerateModelAs("learn_assessment_items", "LearnAssessmentItem"),
		g.GenerateModelAs("learn_side_quests", "LearnSideQuest",
			gen.FieldRelate(field.BelongsTo, "TargetNode", kgNodeModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TargetNodeID"}}}),
		),
		g.GenerateModelAs("learn_student_progress", "LearnStudentProgress",
			gen.FieldRelate(field.BelongsTo, "Node", kgNodeModel, &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"NodeID"}}}),
		),
	)

	// NPC系统关联关系
	g.ApplyBasic(
		g.GenerateModelAs("npc_memories", "NpcMemory"),
		g.GenerateModelAs("npc_profiles", "NpcProfile"),
		g.GenerateModelAs("npc_relationships", "NpcRelationship"),
	)

	// 标签系统关联关系
	g.ApplyBasic(
		g.GenerateModelAs("tags", "Tag",
			gen.FieldRelate(field.HasMany, "Taggables", g.GenerateModelAs("taggables", "Taggable"), &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TagID"}}}),
		),
		g.GenerateModelAs("taggables", "Taggable",
			gen.FieldRelate(field.BelongsTo, "Tag", g.GenerateModelAs("tags", "Tag"), &field.RelateConfig{GORMTag: field.GormTag{"foreignKey": []string{"TagID"}}}),
		),
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
		user = "mathfun"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "math111"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "mathfun"
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
