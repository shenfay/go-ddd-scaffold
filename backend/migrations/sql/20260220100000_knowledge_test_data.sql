-- +goose Up
-- +goose StatementBegin

-- ============================================
-- 知识图谱测试数据
-- ============================================

-- 扩展知识领域数据
INSERT INTO kg_domains (id, name_key, description_key, world_view_type, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000000303', 'domains.measurement.name', 'domains.measurement.description', 'FOUNDATION_BUILDING', '基础数学/计量', true),
('00000000-0000-0000-0000-000000000304', 'domains.statistics.name', 'domains.statistics.description', 'PUBLIC_SQUARE', '基础数学/统计与概率', true)
ON CONFLICT (id) DO NOTHING;

-- 扩展学术概念数据
INSERT INTO kg_academic_concepts (id, name, level, field) VALUES
('00000000-0000-0000-0000-000000000205', '时间计量', '基础', '计量'),
('00000000-0000-0000-0000-000000000206', '长度计量', '基础', '计量'),
('00000000-0000-0000-0000-000000000207', '数据统计', '基础', '统计'),
('00000000-0000-0000-0000-000000000208', '概率基础', '中级', '概率')
ON CONFLICT (id) DO NOTHING;

-- 为主干添加中文注释（如果需要）
COMMENT ON COLUMN kg_trunks.name_key IS '名称的国际化键';
COMMENT ON COLUMN kg_trunks.description_key IS '描述的国际化键';

-- 插入主干数据 - 数与代数领域
INSERT INTO kg_trunks (id, domain_id, name_key, description_key, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000001001', '00000000-0000-0000-0000-000000000301', 'trunk.addition.name', 'trunk.addition.description', '基础数学/数与代数/加法', true),
('00000000-0000-0000-0000-000000001002', '00000000-0000-0000-0000-000000000301', 'trunk.subtraction.name', 'trunk.subtraction.description', '基础数学/数与代数/减法', true),
('00000000-0000-0000-0000-000000001003', '00000000-0000-0000-0000-000000000301', 'trunk.multiplication.name', 'trunk.multiplication.description', '基础数学/数与代数/乘法', true),
('00000000-0000-0000-0000-000000001004', '00000000-0000-0000-0000-000000000301', 'trunk.division.name', 'trunk.division.description', '基础数学/数与代数/除法', true),
('00000000-0000-0000-0000-000000001005', '00000000-0000-0000-0000-000000000301', 'trunk.fraction.name', 'trunk.fraction.description', '基础数学/数与代数/分数', true),
('00000000-0000-0000-0000-000000001006', '00000000-0000-0000-0000-000000000301', 'trunk.decimal.name', 'trunk.decimal.description', '基础数学/数与代数/小数', true)
ON CONFLICT (id) DO NOTHING;

-- 插入主干数据 - 图形与几何领域
INSERT INTO kg_trunks (id, domain_id, name_key, description_key, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000001007', '00000000-0000-0000-0000-000000000302', 'trunk.basic_shapes.name', 'trunk.basic_shapes.description', '基础数学/图形与几何/基本图形', true),
('00000000-0000-0000-0000-000000001008', '00000000-0000-0000-0000-000000000302', 'trunk.measurement.name', 'trunk.measurement.description', '基础数学/图形与几何/测量', true),
('00000000-0000-0000-0000-000000001009', '00000000-0000-0000-0000-000000000302', 'trunk.spatial.name', 'trunk.spatial.description', '基础数学/图形与几何/空间与方位', true)
ON CONFLICT (id) DO NOTHING;

-- 插入主干数据 - 计量领域
INSERT INTO kg_trunks (id, domain_id, name_key, description_key, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000001010', '00000000-0000-0000-0000-000000000303', 'trunk.time.name', 'trunk.time.description', '基础数学/计量/时间', true),
('00000000-0000-0000-0000-000000001011', '00000000-0000-0000-0000-000000000303', 'trunk.money.name', 'trunk.money.description', '基础数学/计量/货币', true),
('00000000-0000-0000-0000-000000001012', '00000000-0000-0000-0000-000000000303', 'trunk.length.name', 'trunk.length.description', '基础数学/计量/长度', true)
ON CONFLICT (id) DO NOTHING;

-- 插入主干数据 - 统计与概率领域
INSERT INTO kg_trunks (id, domain_id, name_key, description_key, academic_source, is_active) VALUES
('00000000-0000-0000-0000-000000001013', '00000000-0000-0000-0000-000000000304', 'trunk.data_collection.name', 'trunk.data_collection.description', '基础数学/统计与概率/数据收集', true),
('00000000-0000-0000-0000-000000001014', '00000000-0000-0000-0000-000000000304', 'trunk.probability.name', 'trunk.probability.description', '基础数学/统计与概率/概率', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 加法主干 (C类型-概念节点)
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002001', '00000000-0000-0000-0000-000000001001', 'C', '00000000-0000-0000-0000-000000000101', 'nodes.addition_concept_lv1.name_child', 'nodes.addition_concept_lv1.name_parent', 'nodes.addition_concept_lv1.description', '00000000-0000-0000-0000-000000000201', true),
('00000000-0000-0000-0000-000000002002', '00000000-0000-0000-0000-000000001001', 'C', '00000000-0000-0000-0000-000000000102', 'nodes.addition_concept_lv2.name_child', 'nodes.addition_concept_lv2.name_parent', 'nodes.addition_concept_lv2.description', '00000000-0000-0000-0000-000000000201', true),
('00000000-0000-0000-0000-000000002003', '00000000-0000-0000-0000-000000001001', 'C', '00000000-0000-0000-0000-000000000103', 'nodes.addition_concept_lv3.name_child', 'nodes.addition_concept_lv3.name_parent', 'nodes.addition_concept_lv3.description', '00000000-0000-0000-0000-000000000201', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 加法主干 (S类型-支撑技能节点)
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002004', '00000000-0000-0000-0000-000000001001', 'S', '00000000-0000-0000-0000-000000000101', 'nodes.addition_skill_lv1.name_child', 'nodes.addition_skill_lv1.name_parent', 'nodes.addition_skill_lv1.description', '00000000-0000-0000-0000-000000000201', true),
('00000000-0000-0000-0000-000000002005', '00000000-0000-0000-0000-000000001001', 'S', '00000000-0000-0000-0000-000000000102', 'nodes.addition_skill_lv2.name_child', 'nodes.addition_skill_lv2.name_parent', 'nodes.addition_skill_lv2.description', '00000000-0000-0000-0000-000000000201', true),
('00000000-0000-0000-0000-000000002006', '00000000-0000-0000-0000-000000001001', 'S', '00000000-0000-0000-0000-000000000103', 'nodes.addition_skill_lv3.name_child', 'nodes.addition_skill_lv3.name_parent', 'nodes.addition_skill_lv3.description', '00000000-0000-0000-0000-000000000201', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 分数主干
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002007', '00000000-0000-0000-0000-000000001005', 'C', '00000000-0000-0000-0000-000000000102', 'nodes.fraction_concept_lv2.name_child', 'nodes.fraction_concept_lv2.name_parent', 'nodes.fraction_concept_lv2.description', '00000000-0000-0000-0000-000000000202', true),
('00000000-0000-0000-0000-000000002008', '00000000-0000-0000-0000-000000001005', 'C', '00000000-0000-0000-0000-000000000103', 'nodes.fraction_concept_lv3.name_child', 'nodes.fraction_concept_lv3.name_parent', 'nodes.fraction_concept_lv3.description', '00000000-0000-0000-0000-000000000202', true),
('00000000-0000-0000-0000-000000002009', '00000000-0000-0000-0000-000000001005', 'S', '00000000-0000-0000-0000-000000000102', 'nodes.fraction_skill_lv2.name_child', 'nodes.fraction_skill_lv2.name_parent', 'nodes.fraction_skill_lv2.description', '00000000-0000-0000-0000-000000000202', true),
('00000000-0000-0000-0000-000000002010', '00000000-0000-0000-0000-000000001005', 'S', '00000000-0000-0000-0000-000000000103', 'nodes.fraction_skill_lv3.name_child', 'nodes.fraction_skill_lv3.name_parent', 'nodes.fraction_skill_lv3.description', '00000000-0000-0000-0000-000000000202', true),
('00000000-0000-0000-0000-000000002011', '00000000-0000-0000-0000-000000001005', 'P', '00000000-0000-0000-0000-000000000104', 'nodes.fraction_problem_lv4.name_child', 'nodes.fraction_problem_lv4.name_parent', 'nodes.fraction_problem_lv4.description', '00000000-0000-0000-0000-000000000202', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 基本图形主干
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002012', '00000000-0000-0000-0000-000000001007', 'C', '00000000-0000-0000-0000-000000000101', 'nodes.circle_concept_lv1.name_child', 'nodes.circle_concept_lv1.name_parent', 'nodes.circle_concept_lv1.description', '00000000-0000-0000-0000-000000000204', true),
('00000000-0000-0000-0000-000000002013', '00000000-0000-0000-0000-000000001007', 'C', '00000000-0000-0000-0000-000000000102', 'nodes.triangle_concept_lv2.name_child', 'nodes.triangle_concept_lv2.name_parent', 'nodes.triangle_concept_lv2.description', '00000000-0000-0000-0000-000000000204', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 时间计量主干
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002014', '00000000-0000-0000-0000-000000001010', 'C', '00000000-0000-0000-0000-000000000101', 'nodes.time_concept_lv1.name_child', 'nodes.time_concept_lv1.name_parent', 'nodes.time_concept_lv1.description', '00000000-0000-0000-0000-000000000205', true),
('00000000-0000-0000-0000-000000002015', '00000000-0000-0000-0000-000000001010', 'S', '00000000-0000-0000-0000-000000000101', 'nodes.time_skill_lv1.name_child', 'nodes.time_skill_lv1.name_parent', 'nodes.time_skill_lv1.description', '00000000-0000-0000-0000-000000000205', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点数据 - 概率主干
INSERT INTO kg_nodes (id, trunk_id, type, competency_level_id, name_child_key, name_parent_key, description_key, academic_concept_id, is_active) VALUES
('00000000-0000-0000-0000-000000002016', '00000000-0000-0000-0000-000000001014', 'C', '00000000-0000-0000-0000-000000000103', 'nodes.probability_concept_lv3.name_child', 'nodes.probability_concept_lv3.name_parent', 'nodes.probability_concept_lv3.description', '00000000-0000-0000-0000-000000000208', true),
('00000000-0000-0000-0000-000000002017', '00000000-0000-0000-0000-000000001014', 'S', '00000000-0000-0000-0000-000000000104', 'nodes.probability_skill_lv4.name_child', 'nodes.probability_skill_lv4.name_parent', 'nodes.probability_skill_lv4.description', '00000000-0000-0000-0000-000000000208', true)
ON CONFLICT (id) DO NOTHING;

-- 插入节点关系数据 - 加法概念依赖关系
INSERT INTO kg_node_relationships (id, from_node_id, to_node_id, relationship_type, weight) VALUES
('00000000-0000-0000-0000-000000003001', '00000000-0000-0000-0000-000000002004', '00000000-0000-0000-0000-000000002001', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003002', '00000000-0000-0000-0000-000000002005', '00000000-0000-0000-0000-000000002002', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003003', '00000000-0000-0000-0000-000000002006', '00000000-0000-0000-0000-000000002003', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003004', '00000000-0000-0000-0000-000000002005', '00000000-0000-0000-0000-000000002004', 'PREREQ', 0.8),
('00000000-0000-0000-0000-000000003005', '00000000-0000-0000-0000-000000002006', '00000000-0000-0000-0000-000000002005', 'PREREQ', 0.8)
ON CONFLICT (id) DO NOTHING;

-- 插入节点关系数据 - 分数依赖关系
INSERT INTO kg_node_relationships (id, from_node_id, to_node_id, relationship_type, weight) VALUES
('00000000-0000-0000-0000-000000003006', '00000000-0000-0000-0000-000000002009', '00000000-0000-0000-0000-000000002007', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003007', '00000000-0000-0000-0000-000000002010', '00000000-0000-0000-0000-000000002008', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003008', '00000000-0000-0000-0000-000000002011', '00000000-0000-0000-0000-000000002010', 'PREREQ', 1.0),
('00000000-0000-0000-0000-000000003009', '00000000-0000-0000-0000-000000002011', '00000000-0000-0000-0000-000000002009', 'SUP_SKILL', 0.7)
ON CONFLICT (id) DO NOTHING;

-- 插入节点关系数据 - 支撑技能关系（加法 -> 减法）
INSERT INTO kg_node_relationships (id, from_node_id, to_node_id, relationship_type, weight) VALUES
('00000000-0000-0000-0000-000000003010', '00000000-0000-0000-0000-000000002004', '00000000-0000-0000-0000-000000002001', 'SUP_SKILL', 0.9)
ON CONFLICT (id) DO NOTHING;

-- 扩展标签数据
INSERT INTO tags (id, name_key, description_key, category, is_active) VALUES
('00000000-0000-0000-0000-000000000404', 'tags.fraction.name', 'tags.fraction.description', '知识主题', true),
('00000000-0000-0000-0000-000000000405', 'tags.addition.name', 'tags.addition.description', '知识主题', true),
('00000000-0000-0000-0000-000000000406', 'tags.time.name', 'tags.time.description', '知识主题', true),
('00000000-0000-0000-0000-000000000407', 'tags.probability.name', 'tags.probability.description', '知识主题', true)
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- 删除测试数据（按依赖顺序逆序删除）

-- 删除节点关系
DELETE FROM kg_node_relationships WHERE id IN (
    '00000000-0000-0000-0000-000000003001',
    '00000000-0000-0000-0000-000000003002',
    '00000000-0000-0000-0000-000000003003',
    '00000000-0000-0000-0000-000000003004',
    '00000000-0000-0000-0000-000000003005',
    '00000000-0000-0000-0000-000000003006',
    '00000000-0000-0000-0000-000000003007',
    '00000000-0000-0000-0000-000000003008',
    '00000000-0000-0000-0000-000000003009',
    '00000000-0000-0000-0000-000000003010'
);

-- 删除节点
DELETE FROM kg_nodes WHERE id IN (
    '00000000-0000-0000-0000-000000002001',
    '00000000-0000-0000-0000-000000002002',
    '00000000-0000-0000-0000-000000002003',
    '00000000-0000-0000-0000-000000002004',
    '00000000-0000-0000-0000-000000002005',
    '00000000-0000-0000-0000-000000002006',
    '00000000-0000-0000-0000-000000002007',
    '00000000-0000-0000-0000-000000002008',
    '00000000-0000-0000-0000-000000002009',
    '00000000-0000-0000-0000-000000002010',
    '00000000-0000-0000-0000-000000002011',
    '00000000-0000-0000-0000-000000002012',
    '00000000-0000-0000-0000-000000002013',
    '00000000-0000-0000-0000-000000002014',
    '00000000-0000-0000-0000-000000002015',
    '00000000-0000-0000-0000-000000002016',
    '00000000-0000-0000-0000-000000002017'
);

-- 删除主干
DELETE FROM kg_trunks WHERE id IN (
    '00000000-0000-0000-0000-000000001001',
    '00000000-0000-0000-0000-000000001002',
    '00000000-0000-0000-0000-000000001003',
    '00000000-0000-0000-0000-000000001004',
    '00000000-0000-0000-0000-000000001005',
    '00000000-0000-0000-0000-000000001006',
    '00000000-0000-0000-0000-000000001007',
    '00000000-0000-0000-0000-000000001008',
    '00000000-0000-0000-0000-000000001009',
    '00000000-0000-0000-0000-000000001010',
    '00000000-0000-0000-0000-000000001011',
    '00000000-0000-0000-0000-000000001012',
    '00000000-0000-0000-0000-000000001013',
    '00000000-0000-0000-0000-000000001014'
);

-- 删除新增标签
DELETE FROM tags WHERE id IN (
    '00000000-0000-0000-0000-000000000404',
    '00000000-0000-0000-0000-000000000405',
    '00000000-0000-0000-0000-000000000406',
    '00000000-0000-0000-0000-000000000407'
);

-- 删除新增领域
DELETE FROM kg_domains WHERE id IN (
    '00000000-0000-0000-0000-000000000303',
    '00000000-0000-0000-0000-000000000304'
);

-- 删除新增学术概念
DELETE FROM kg_academic_concepts WHERE id IN (
    '00000000-0000-0000-0000-000000000205',
    '00000000-0000-0000-0000-000000000206',
    '00000000-0000-0000-0000-000000000207',
    '00000000-0000-0000-0000-000000000208'
);

-- +goose StatementEnd
