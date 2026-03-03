// Package converter 提供实体与模型之间的转换功能
package converter

import (
	"encoding/json"
	"go-ddd-scaffold/internal/domain/knowledge/entity"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"

	"github.com/google/uuid"
)

// ===== 知识领域转换函数 =====

// DomainToModel 将领域实体转换为数据库模型
func DomainToModel(e *entity.Domain) *model.KgDomain {
	if e == nil {
		return nil
	}

	// 处理可选字段
	var id *string
	if e.ID != uuid.Nil {
		uuidStr := e.ID.String()
		id = &uuidStr
	}

	// 转换Metadata
	var metadata *string
	if e.Metadata != nil {
		if metaBytes, err := json.Marshal(e.Metadata); err == nil {
			metaStr := string(metaBytes)
			metadata = &metaStr
		}
	}

	// 转换WorldViewConfig
	var worldviewConfig *string
	if e.WorldViewConfig != nil {
		if configBytes, err := json.Marshal(e.WorldViewConfig); err == nil {
			configStr := string(configBytes)
			worldviewConfig = &configStr
		}
	}

	isActive := e.IsActive

	return &model.KgDomain{
		ID:              id,
		NameKey:         &e.NameKey,
		DescriptionKey:  &e.DescriptionKey,
		WorldViewType:   &e.WorldViewType,
		WorldViewConfig: worldviewConfig,
		AcademicSource:  &e.AcademicSource,
		IsActive:        &isActive,
		Metadata:        metadata,
	}
}

// ModelToDomain 将数据库模型转换为领域实体
func ModelToDomain(m *model.KgDomain) *entity.Domain {
	if m == nil {
		return nil
	}

	// 处理ID
	var id uuid.UUID
	if m.ID != nil {
		id, _ = uuid.Parse(*m.ID)
	}

	// 处理基础字段
	nameKey := ""
	if m.NameKey != nil {
		nameKey = *m.NameKey
	}

	descriptionKey := ""
	if m.DescriptionKey != nil {
		descriptionKey = *m.DescriptionKey
	}

	worldViewType := ""
	if m.WorldViewType != nil {
		worldViewType = *m.WorldViewType
	}

	academicSource := ""
	if m.AcademicSource != nil {
		academicSource = *m.AcademicSource
	}

	// 转换Metadata
	var metadata map[string]interface{}
	if m.Metadata != nil {
		json.Unmarshal([]byte(*m.Metadata), &metadata)
	}

	// 转换WorldViewConfig
	var worldviewConfig map[string]interface{}
	if m.WorldViewConfig != nil {
		json.Unmarshal([]byte(*m.WorldViewConfig), &worldviewConfig)
	}

	isActive := true
	if m.IsActive != nil {
		isActive = *m.IsActive
	}

	// 注意：这里假设entity.Domain有相应的构造函数
	// 实际使用时需要根据具体的实体构造方式调整
	return &entity.Domain{
		ID:              id,
		NameKey:         nameKey,
		DescriptionKey:  descriptionKey,
		WorldViewType:   worldViewType,
		WorldViewConfig: worldviewConfig,
		AcademicSource:  academicSource,
		IsActive:        isActive,
		Metadata:        metadata,
	}
}

// ModelsToDomains 批量转换
func ModelsToDomains(models []*model.KgDomain) []*entity.Domain {
	entities := make([]*entity.Domain, 0, len(models))
	for _, m := range models {
		if entity := ModelToDomain(m); entity != nil {
			entities = append(entities, entity)
		}
	}
	return entities
}

// DomainsToModels 批量转换
func DomainsToModels(entities []*entity.Domain) []*model.KgDomain {
	models := make([]*model.KgDomain, 0, len(entities))
	for _, e := range entities {
		if model := DomainToModel(e); model != nil {
			models = append(models, model)
		}
	}
	return models
}

// ===== 主干转换函数 =====

// ModelToTrunk 将数据库模型转换为主干领域实体
func ModelToTrunk(m *model.KgTrunk) *entity.Trunk {
	if m == nil {
		return nil
	}

	// 处理ID
	id := uuid.Nil
	if m.ID != nil {
		id, _ = uuid.Parse(*m.ID)
	}

	// 处理DomainID
	domainID := uuid.Nil
	if m.DomainID != nil {
		domainID, _ = uuid.Parse(*m.DomainID)
	}

	// 处理基础字段
	nameKey := ""
	if m.NameKey != nil {
		nameKey = *m.NameKey
	}

	descriptionKey := ""
	if m.DescriptionKey != nil {
		descriptionKey = *m.DescriptionKey
	}

	academicSource := ""
	if m.AcademicSource != nil {
		academicSource = *m.AcademicSource
	}

	// 转换Metadata
	var metadata map[string]interface{}
	if m.Metadata != nil {
		json.Unmarshal([]byte(*m.Metadata), &metadata)
	}

	isActive := true
	if m.IsActive != nil {
		isActive = *m.IsActive
	}

	return &entity.Trunk{
		ID:             id,
		DomainID:       domainID,
		NameKey:        nameKey,
		DescriptionKey: descriptionKey,
		AcademicSource: academicSource,
		IsActive:       isActive,
		Metadata:       metadata,
	}
}

// ModelsToTrunks 批量转换
func ModelsToTrunks(models []*model.KgTrunk) []*entity.Trunk {
	entities := make([]*entity.Trunk, 0, len(models))
	for _, m := range models {
		if e := ModelToTrunk(m); e != nil {
			entities = append(entities, e)
		}
	}
	return entities
}

// ===== 辅助函数 =====

// GetStringPtr 获取字符串指针
func GetStringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GetValue 从指针获取值，如果为nil则返回零值
func GetValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// GetBoolValue 从指针获取布尔值
func GetBoolValue(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}

// ===== 节点转换函数 =====

// ModelToNode 将数据库模型转换为领域实体
func ModelToNode(m *model.KgNode) *entity.Node {
	if m == nil {
		return nil
	}

	// 处理ID
	id := uuid.Nil
	if m.ID != nil {
		id, _ = uuid.Parse(*m.ID)
	}

	// 处理TrunkID
	trunkID := uuid.Nil
	if m.TrunkID != nil {
		trunkID, _ = uuid.Parse(*m.TrunkID)
	}

	// 处理Type
	nodeType := entity.NodeType(m.Type)

	// 处理基础字段
	nameChildKey := ""
	if m.NameChildKey != nil {
		nameChildKey = *m.NameChildKey
	}

	nameParentKey := ""
	if m.NameParentKey != nil {
		nameParentKey = *m.NameParentKey
	}

	descriptionKey := ""
	if m.DescriptionKey != nil {
		descriptionKey = *m.DescriptionKey
	}

	competencyLevelID := ""
	if m.CompetencyLevelID != nil {
		competencyLevelID = *m.CompetencyLevelID
	}

	academicConceptID := ""
	if m.AcademicConceptID != nil {
		academicConceptID = *m.AcademicConceptID
	}

	// 转换Resources
	var resources map[string]interface{}
	if m.Resources != nil {
		json.Unmarshal([]byte(*m.Resources), &resources)
	}

	// 转换Metadata
	var metadata map[string]interface{}
	if m.Metadata != nil {
		json.Unmarshal([]byte(*m.Metadata), &metadata)
	}

	isActive := true
	if m.IsActive != nil {
		isActive = *m.IsActive
	}

	return &entity.Node{
		ID:                id,
		TrunkID:           trunkID,
		Type:              nodeType,
		CompetencyLevelID: competencyLevelID,
		NameChildKey:      nameChildKey,
		NameParentKey:     nameParentKey,
		DescriptionKey:    descriptionKey,
		AcademicConceptID: academicConceptID,
		Resources:         resources,
		IsActive:          isActive,
		Metadata:          metadata,
	}
}

// ModelsToNodes 批量转换
func ModelsToNodes(models []*model.KgNode) []*entity.Node {
	entities := make([]*entity.Node, 0, len(models))
	for _, m := range models {
		if e := ModelToNode(m); e != nil {
			entities = append(entities, e)
		}
	}
	return entities
}
