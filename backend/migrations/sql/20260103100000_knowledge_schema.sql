-- +goose Up
-- +goose StatementBegin

-- Create kg_competency_levels table
-- 能力等级表：定义学习能力的等级体系（如Bloom认知等级、SOLO分类等）
CREATE TABLE IF NOT EXISTS kg_competency_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 等级唯一标识符
    model_name VARCHAR(255) NOT NULL, -- 等级模型名称 (e.g., "Bloom_SOLO_Fusion")
    level_value VARCHAR(100) NOT NULL, -- 级别值 (e.g., "Lv1", "Advanced")
    description_key TEXT, -- 该级别描述的国际化键 (e.g., "levels.bloom_solo_fusion.lv1.description")
    min_mastery_score DECIMAL(3,2), -- 触发此级别的最低掌握度阈值
    node_id UUID, -- 关联的知识图谱节点ID，用于将能力等级与知识节点关联
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- 添加能力等级表注释
COMMENT ON TABLE kg_competency_levels IS '能力等级表：定义学习能力的等级体系（如Bloom认知等级、SOLO分类等）';
COMMENT ON COLUMN kg_competency_levels.id IS '等级唯一标识符';
COMMENT ON COLUMN kg_competency_levels.model_name IS '等级模型名称 (e.g., "Bloom_SOLO_Fusion")';
COMMENT ON COLUMN kg_competency_levels.level_value IS '级别值 (e.g., "Lv1", "Advanced")';
COMMENT ON COLUMN kg_competency_levels.description_key IS '该级别描述的国际化键 (e.g., "levels.bloom_solo_fusion.lv1.description")';
COMMENT ON COLUMN kg_competency_levels.min_mastery_score IS '触发此级别的最低掌握度阈值';
COMMENT ON COLUMN kg_competency_levels.node_id IS '关联的知识图谱节点ID，用于将能力等级与知识节点关联';
COMMENT ON COLUMN kg_competency_levels.metadata IS '额外元数据';
COMMENT ON COLUMN kg_competency_levels.created_at IS '创建时间';
COMMENT ON COLUMN kg_competency_levels.updated_at IS '更新时间';

-- Create kg_academic_concepts table
-- 学术概念表：存储数学学科的核心概念定义
CREATE TABLE IF NOT EXISTS kg_academic_concepts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 学术概念唯一标识符
    name VARCHAR(255) NOT NULL, -- 概念名称 (e.g., "方程")
    level VARCHAR(100), -- 学术层级 (e.g., "基础", "高级")
    field VARCHAR(255), -- 所属学术领域 (e.g., "代数")
    node_id UUID, -- 关联的知识图谱节点ID，用于将学术概念与知识节点关联
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- 添加学术概念表注释
COMMENT ON TABLE kg_academic_concepts IS '学术概念表：存储数学学科的核心概念定义';
COMMENT ON COLUMN kg_academic_concepts.id IS '学术概念唯一标识符';
COMMENT ON COLUMN kg_academic_concepts.name IS '概念名称 (e.g., "方程")';
COMMENT ON COLUMN kg_academic_concepts.level IS '学术层级 (e.g., "基础", "高级")';
COMMENT ON COLUMN kg_academic_concepts.field IS '所属学术领域 (e.g., "代数")';
COMMENT ON COLUMN kg_academic_concepts.node_id IS '关联的知识图谱节点ID，用于将学术概念与知识节点关联';
COMMENT ON COLUMN kg_academic_concepts.metadata IS '额外元数据';
COMMENT ON COLUMN kg_academic_concepts.created_at IS '创建时间';
COMMENT ON COLUMN kg_academic_concepts.updated_at IS '更新时间';

-- Create kg_domains table
CREATE TABLE IF NOT EXISTS kg_domains (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 领域唯一标识符
    name_key VARCHAR(500), -- 名称的国际化键 (e.g., "domains.numeracy_algebra.name")
    description_key TEXT, -- 描述的国际化键 (e.g., "domains.numeracy_algebra.description")
    world_view_type VARCHAR(100), -- 领域的世界观类型 (e.g., "FOUNDATION_BUILDING", "PUBLIC_SQUARE")
    world_view_config JSONB DEFAULT '{}', -- 世界观相关的配置信息 (e.g., `{ "building_style": "fortress", "floors": { "integer_ops": 1 } }`)
    academic_source TEXT, -- 学术源头 (e.g., "基础数学/代数")
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- 添加领域表注释
COMMENT ON TABLE kg_domains IS '知识图谱领域表：定义数学知识的不同领域（如数与代数、图形与几何等）';
COMMENT ON COLUMN kg_domains.id IS '领域唯一标识符';
COMMENT ON COLUMN kg_domains.name_key IS '名称的国际化键 (e.g., "domains.numeracy_algebra.name")';
COMMENT ON COLUMN kg_domains.description_key IS '描述的国际化键 (e.g., "domains.numeracy_algebra.description")';
COMMENT ON COLUMN kg_domains.world_view_type IS '领域的世界观类型 (e.g., "FOUNDATION_BUILDING", "PUBLIC_SQUARE")';
COMMENT ON COLUMN kg_domains.world_view_config IS '世界观相关的配置信息';
COMMENT ON COLUMN kg_domains.academic_source IS '学术源头 (e.g., "基础数学/代数")';
COMMENT ON COLUMN kg_domains.is_active IS '是否激活，默认 True';
COMMENT ON COLUMN kg_domains.metadata IS '额外元数据';
COMMENT ON COLUMN kg_domains.created_at IS '创建时间';
COMMENT ON COLUMN kg_domains.updated_at IS '更新时间';

-- Create kg_trunks table
CREATE TABLE IF NOT EXISTS kg_trunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 主线唯一标识符
    domain_id UUID REFERENCES kg_domains(id) ON DELETE CASCADE, -- 关联的领域 ID (kg_domains.id)
    name_key VARCHAR(500), -- 名称的国际化键 (e.g., "trunks.integer_ops_trunk.name")
    description_key TEXT, -- 描述的国际化键 (e.g., "trunks.integer_ops_trunk.description")
    academic_source TEXT, -- 学术源头
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create kg_nodes table
CREATE TABLE IF NOT EXISTS kg_nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 节点唯一标识符
    trunk_id UUID REFERENCES kg_trunks(id) ON DELETE CASCADE, -- 关联的主线 ID (kg_trunks.id)
    type VARCHAR(10) CHECK (type IN ('C', 'S', 'T', 'P')) NOT NULL, -- 节点类型: C (概念), S (支撑技能), T (思维模式), P (问题模型)
    competency_level_id UUID REFERENCES kg_competency_levels(id), -- 关联的能力等级 ID (kg_competency_levels.id)
    name_child_key VARCHAR(500), -- 孩子视角名称的国际化键 (e.g., "nodes.fraction_lv3_node.name_child")
    name_parent_key VARCHAR(500), -- 家长/教研视角名称的国际化键 (e.g., "nodes.fraction_lv3_node.name_parent")
    description_key TEXT, -- 节点描述的国际化键 (e.g., "nodes.fraction_lv3_node.description")
    academic_concept_id UUID REFERENCES kg_academic_concepts(id), -- 关联的学术概念 ID (kg_academic_concepts.id)
    resources JSONB DEFAULT '{}', -- 关联的学习资源 (e.g., `{ "video": "v001", "game": "g102" }`)
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create kg_node_relationships table
CREATE TABLE IF NOT EXISTS kg_node_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 关系唯一标识符
    from_node_id UUID REFERENCES kg_nodes(id) ON DELETE CASCADE, -- 依赖的源节点 ID
    to_node_id UUID REFERENCES kg_nodes(id) ON DELETE CASCADE, -- 被依赖的目标节点 ID
    relationship_type VARCHAR(50) CHECK (relationship_type IN ('PREREQ', 'SUP_SKILL', 'THINK_PAT')), -- 关系类型: PREREQUISITE (前置), SUPPORTING_SKILL (支撑技能), THINKING_PATTERN (思维模式)
    weight DECIMAL(3,2) DEFAULT 1.0, -- 关系权重 (可选)
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 标签唯一标识符
    name_key VARCHAR(500), -- 标签名的国际化键 (e.g., "tags.fraction.name")
    description_key TEXT, -- 标签描述的国际化键 (e.g., "tags.fraction.description")
    category VARCHAR(100), -- 标签类别 (e.g., "知识主题", "能力等级", "应用场景")
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create taggables table
CREATE TABLE IF NOT EXISTS taggables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 关联唯一标识符
    taggable_type VARCHAR(100) NOT NULL, -- 被标记实体的类型 (e.g., "kg_nodes", "learn_assessment_items")
    taggable_id UUID NOT NULL, -- 被标记实体的 ID
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE, -- 关联的标签 ID (tags.id)
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    UNIQUE(taggable_type, taggable_id, tag_id) -- 确保一个实体的同一个标签不重复
);

-- Create learn_student_progress table
CREATE TABLE IF NOT EXISTS learn_student_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 进度记录唯一标识符
    student_id VARCHAR(255) NOT NULL, -- 学员 ID
    node_id UUID REFERENCES kg_nodes(id) ON DELETE CASCADE, -- 关联的节点 ID (kg_nodes.id)
    mastery_score DECIMAL(5,4) DEFAULT 0.0000, -- 掌握度 (0.0 - 1.0)
    last_interacted_at TIMESTAMPTZ, -- 最后互动时间
    attempts_count INTEGER DEFAULT 0, -- 尝试次数
    metadata JSONB DEFAULT '{}', -- 额外元数据 (如答题详情)
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create learn_assessment_items table
CREATE TABLE IF NOT EXISTS learn_assessment_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 评测项唯一标识符
    question_text_key TEXT, -- 题目内容的国际化键 (e.g., "questions.fraction_addition_lv3.question_text")
    correct_answer TEXT, -- 正确答案
    difficulty_level VARCHAR(50), -- 难度等级 (e.g., "Easy", "Lv1")
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create learn_side_quests table
CREATE TABLE IF NOT EXISTS learn_side_quests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 任务唯一标识符
    name_key VARCHAR(500), -- 任务名称的国际化键 (e.g., "quests.gear_workshop.name")
    description_key TEXT, -- 任务描述的国际化键 (e.g., "quests.gear_workshop.description")
    target_node_id UUID REFERENCES kg_nodes(id) ON DELETE CASCADE, -- 强化的目标节点 ID (kg_nodes.id)
    category VARCHAR(100), -- 任务分类 (e.g., "drill_and_practice", "story_based", "game_like", "application")
    difficulty_level VARCHAR(50), -- 任务难度 (e.g., "Easy", "Medium", "Hard", "Lv1", "Lv2")
    weight DECIMAL(5,2) DEFAULT 1.00, -- 任务触发权重 (e.g., 1.0)
    trigger_condition_json JSONB DEFAULT '{}', -- 触发条件的详细描述 (可选)
    content JSONB DEFAULT '{}', -- 任务具体内容
    reward JSONB DEFAULT '{}', -- 奖励描述
    is_active BOOLEAN DEFAULT TRUE, -- 是否激活，默认 True
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create npc_profiles table (for LLM-driven NPCs)
-- NPC档案表：存储LLM驱动的NPC角色配置信息
CREATE TABLE IF NOT EXISTS npc_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- NPC唯一标识符
    name_key VARCHAR(500) UNIQUE NOT NULL, -- 名称国际化键
    base_personality TEXT, -- 基础性格描述（Prompt Base）
    math_expertise_level INTEGER DEFAULT 1, -- 数学专业度等级
    visual_config JSONB DEFAULT '{}', -- 形象配置（对应前端模型组件参数）
    node_id UUID, -- 关联的知识图谱节点ID，用于将NPC与知识节点关联
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW(), -- 更新时间
    deleted_at TIMESTAMPTZ -- 软删除时间戳
);

-- Create npc_memories table (for LLM-driven NPCs)
-- NPC记忆表：存储NPC与用户的交互记忆，支持RAG检索
CREATE TABLE IF NOT EXISTS npc_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 记忆唯一标识符
    npc_id UUID REFERENCES npc_profiles(id) ON DELETE CASCADE, -- 关联 NPC
    user_id VARCHAR(255) NOT NULL, -- 关联用户
    memory_type VARCHAR(20) CHECK (memory_type IN ('short_term', 'long_term')) NOT NULL, -- 长期记忆/短期记忆
    content TEXT NOT NULL, -- 事实性记忆内容（如"用户昨天没学会分数加法"）
    domain_tag VARCHAR(50), -- 记忆所属的知识领域标签 (e.g., 'arithmetic', 'geometry')
    vector_ref VARCHAR(255), -- 向量数据库索引（可选，用于 RAG 检索）
    node_id UUID, -- 关联的知识图谱节点ID，用于将记忆与知识节点关联
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create npc_relationships table (for LLM-driven NPCs)
CREATE TABLE IF NOT EXISTS npc_relationships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- 关系唯一标识符
    from_npc_id UUID REFERENCES npc_profiles(id) ON DELETE CASCADE, -- 关联发起关系的 NPC ID
    to_id VARCHAR(255) NOT NULL, -- 关联被指向的 NPC ID 或 User ID
    relationship_type VARCHAR(50) NOT NULL, -- 关系类型 (e.g., 'friend', 'colleague', 'student')
    affinity_level DECIMAL(3,2) DEFAULT 0.00, -- 亲密度等级 (FLOAT)
    metadata JSONB DEFAULT '{}', -- 额外元数据
    created_at TIMESTAMPTZ DEFAULT NOW(), -- 创建时间
    updated_at TIMESTAMPTZ DEFAULT NOW() -- 更新时间
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_kg_nodes_trunk_id ON kg_nodes(trunk_id);
CREATE INDEX IF NOT EXISTS idx_kg_nodes_type ON kg_nodes(type);
CREATE INDEX IF NOT EXISTS idx_kg_nodes_competency_level_id ON kg_nodes(competency_level_id);
CREATE INDEX IF NOT EXISTS idx_kg_nodes_academic_concept_id ON kg_nodes(academic_concept_id);
CREATE INDEX IF NOT EXISTS idx_kg_nodes_is_active ON kg_nodes(is_active);

CREATE INDEX IF NOT EXISTS idx_kg_node_relationships_from_node_id ON kg_node_relationships(from_node_id);
CREATE INDEX IF NOT EXISTS idx_kg_node_relationships_to_node_id ON kg_node_relationships(to_node_id);
CREATE INDEX IF NOT EXISTS idx_kg_node_relationships_relationship_type ON kg_node_relationships(relationship_type);

CREATE INDEX IF NOT EXISTS idx_learn_student_progress_student_id ON learn_student_progress(student_id);
CREATE INDEX IF NOT EXISTS idx_learn_student_progress_node_id ON learn_student_progress(node_id);

CREATE INDEX IF NOT EXISTS idx_kg_competency_levels_node_id ON kg_competency_levels(node_id);
CREATE INDEX IF NOT EXISTS idx_kg_academic_concepts_node_id ON kg_academic_concepts(node_id);

CREATE INDEX IF NOT EXISTS idx_npc_profiles_node_id ON npc_profiles(node_id);

CREATE INDEX IF NOT EXISTS idx_npc_memories_npc_id ON npc_memories(npc_id);
CREATE INDEX IF NOT EXISTS idx_npc_memories_user_id ON npc_memories(user_id);
CREATE INDEX IF NOT EXISTS idx_npc_memories_domain_tag ON npc_memories(domain_tag);
CREATE INDEX IF NOT EXISTS idx_npc_memories_memory_type ON npc_memories(memory_type);
CREATE INDEX IF NOT EXISTS idx_npc_memories_node_id ON npc_memories(node_id);

CREATE INDEX IF NOT EXISTS idx_npc_relationships_from_npc_id ON npc_relationships(from_npc_id);
CREATE INDEX IF NOT EXISTS idx_npc_relationships_to_id ON npc_relationships(to_id);

-- Create index for taggables table
CREATE INDEX IF NOT EXISTS idx_taggables_taggable_type_id ON taggables(taggable_type, taggable_id);
CREATE INDEX IF NOT EXISTS idx_taggables_tag_id ON taggables(tag_id);

-- Add foreign key constraints after all tables are created
-- 在所有表创建完成后添加外键约束，避免循环依赖问题
ALTER TABLE kg_competency_levels ADD CONSTRAINT fk_competency_levels_node_id 
    FOREIGN KEY (node_id) REFERENCES kg_nodes(id) ON DELETE SET NULL;
ALTER TABLE kg_academic_concepts ADD CONSTRAINT fk_academic_concepts_node_id 
    FOREIGN KEY (node_id) REFERENCES kg_nodes(id) ON DELETE SET NULL;
ALTER TABLE npc_profiles ADD CONSTRAINT fk_npc_profiles_node_id 
    FOREIGN KEY (node_id) REFERENCES kg_nodes(id) ON DELETE SET NULL;
ALTER TABLE npc_memories ADD CONSTRAINT fk_npc_memories_node_id 
    FOREIGN KEY (node_id) REFERENCES kg_nodes(id) ON DELETE SET NULL;

-- ============================================
-- 插入初始化数据
-- ============================================

-- 插入默认能力等级模型
INSERT INTO kg_competency_levels (id, model_name, level_value, description_key, min_mastery_score) VALUES
('00000000-0000-0000-0000-000000000101', 'Bloom_SOLO_Fusion', 'Lv1', 'levels.bloom_solo_fusion.lv1.description', 0.60),
('00000000-0000-0000-0000-000000000102', 'Bloom_SOLO_Fusion', 'Lv2', 'levels.bloom_solo_fusion.lv2.description', 0.70),
('00000000-0000-0000-0000-000000000103', 'Bloom_SOLO_Fusion', 'Lv3', 'levels.bloom_solo_fusion.lv3.description', 0.80),
('00000000-0000-0000-0000-000000000104', 'Bloom_SOLO_Fusion', 'Lv4', 'levels.bloom_solo_fusion.lv4.description', 0.85),
('00000000-0000-0000-0000-000000000105', 'Bloom_SOLO_Fusion', 'Lv5', 'levels.bloom_solo_fusion.lv5.description', 0.90)
ON CONFLICT (id) DO NOTHING;

-- 插入基础学术概念
INSERT INTO kg_academic_concepts (id, name, level, field) VALUES
('00000000-0000-0000-0000-000000000201', '整数运算', '基础', '算术'),
('00000000-0000-0000-0000-000000000202', '分数概念', '基础', '算术'),
('00000000-0000-0000-0000-000000000203', '代数表达式', '中级', '代数'),
('00000000-0000-0000-0000-000000000204', '平面几何', '基础', '几何')
ON CONFLICT (id) DO NOTHING;

-- 插入默认知识领域
INSERT INTO kg_domains (id, name_key, description_key, world_view_type, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000000301', 'domains.numeracy_algebra.name', 'domains.numeracy_algebra.description', 'FOUNDATION_BUILDING', '基础数学/数与代数', true),
('00000000-0000-0000-0000-000000000302', 'domains.geometry.name', 'domains.geometry.description', 'PUBLIC_SQUARE', '基础数学/图形与几何', true)
ON CONFLICT (id) DO NOTHING;

-- 插入默认标签
INSERT INTO tags (id, name_key, description_key, category, is_active) VALUES
('00000000-0000-0000-0000-000000000401', 'tags.arithmetic.name', 'tags.arithmetic.description', '知识主题', true),
('00000000-0000-0000-0000-000000000402', 'tags.algebra.name', 'tags.algebra.description', '知识主题', true),
('00000000-0000-0000-0000-000000000403', 'tags.geometry.name', 'tags.geometry.description', '知识主题', true)
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign key constraints first
ALTER TABLE IF EXISTS kg_competency_levels DROP CONSTRAINT IF EXISTS fk_competency_levels_node_id;
ALTER TABLE IF EXISTS kg_academic_concepts DROP CONSTRAINT IF EXISTS fk_academic_concepts_node_id;
ALTER TABLE IF EXISTS npc_profiles DROP CONSTRAINT IF EXISTS fk_npc_profiles_node_id;
ALTER TABLE IF EXISTS npc_memories DROP CONSTRAINT IF EXISTS fk_npc_memories_node_id;

-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS npc_relationships;
DROP TABLE IF EXISTS npc_memories;
DROP TABLE IF EXISTS npc_profiles;
DROP TABLE IF EXISTS learn_side_quests;
DROP TABLE IF EXISTS learn_assessment_items;
DROP TABLE IF EXISTS learn_student_progress;
DROP TABLE IF EXISTS taggables;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS kg_node_relationships;
DROP TABLE IF EXISTS kg_nodes;
DROP TABLE IF EXISTS kg_trunks;
DROP TABLE IF EXISTS kg_domains;
DROP TABLE IF EXISTS kg_academic_concepts;
DROP TABLE IF EXISTS kg_competency_levels;

-- +goose StatementEnd