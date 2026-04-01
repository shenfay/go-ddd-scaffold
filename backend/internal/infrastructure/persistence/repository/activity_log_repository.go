package repository

import (
	"context"
	"encoding/json"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	daoModel "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

// ActivityLogRepositoryImpl 活动日志仓储实现
type ActivityLogRepositoryImpl struct {
	query *dao.Query
}

// NewActivityLogRepository 创建活动日志仓储
func NewActivityLogRepository(db *dao.Query) model.ActivityLogRepository {
	return &ActivityLogRepositoryImpl{query: db}
}

// Save 保存活动日志
func (r *ActivityLogRepositoryImpl) Save(ctx context.Context, log *model.ActivityLog) error {
	daoModel := r.fromDomain(log)
	return r.query.ActivityLog.WithContext(ctx).Create(daoModel)
}

// FindByUserID 按用户 ID 查询
func (r *ActivityLogRepositoryImpl) FindByUserID(ctx context.Context, userID int64, limit int) ([]*model.ActivityLog, error) {
	daoModels, err := r.query.ActivityLog.WithContext(ctx).
		Where(r.query.ActivityLog.UserID.Eq(userID)).
		Order(r.query.ActivityLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*model.ActivityLog, len(daoModels))
	for i, m := range daoModels {
		logs[i] = r.toDomain(m)
	}
	return logs, nil
}

// FindByAction 按操作类型查询
func (r *ActivityLogRepositoryImpl) FindByAction(ctx context.Context, action model.ActivityType, limit int) ([]*model.ActivityLog, error) {
	daoModels, err := r.query.ActivityLog.WithContext(ctx).
		Where(r.query.ActivityLog.Action.Eq(string(action))).
		Order(r.query.ActivityLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*model.ActivityLog, len(daoModels))
	for i, m := range daoModels {
		logs[i] = r.toDomain(m)
	}
	return logs, nil
}

// FindFailed 查询失败的活动
func (r *ActivityLogRepositoryImpl) FindFailed(ctx context.Context, limit int) ([]*model.ActivityLog, error) {
	daoModels, err := r.query.ActivityLog.WithContext(ctx).
		Where(r.query.ActivityLog.Status.Eq(int16(model.ActivityStatusFailed))).
		Order(r.query.ActivityLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*model.ActivityLog, len(daoModels))
	for i, m := range daoModels {
		logs[i] = r.toDomain(m)
	}
	return logs, nil
}

// fromDomain 将领域模型转换为 DAO 模型
func (r *ActivityLogRepositoryImpl) fromDomain(log *model.ActivityLog) *daoModel.ActivityLog {
	metadataJSON, _ := json.Marshal(log.Metadata)
	status := int16(log.Status)
	metadataStr := string(metadataJSON)

	return &daoModel.ActivityLog{
		ID:         log.ID,
		TenantID:   log.TenantID, // 指针直接复制
		UserID:     log.UserID,
		Action:     string(log.Action),
		Status:     &status,
		IPAddress:  util.StringPtrNilIfEmpty(log.IPAddress), // 空字符串转为 nil
		UserAgent:  util.StringPtrNilIfEmpty(log.UserAgent), // 空字符串转为 nil
		Metadata:   &metadataStr,                            // 取地址
		OccurredAt: log.OccurredAt,
		CreatedAt:  &log.CreatedAt,
	}
}

// toDomain 将 DAO 模型转换为领域模型
func (r *ActivityLogRepositoryImpl) toDomain(m *daoModel.ActivityLog) *model.ActivityLog {
	var metadata map[string]any
	if m.Metadata != nil && *m.Metadata != "" {
		json.Unmarshal([]byte(*m.Metadata), &metadata)
	}

	var tenantID *int64
	if m.TenantID != nil {
		tenantID = m.TenantID
	}

	var status model.ActivityStatus
	if m.Status != nil {
		status = model.ActivityStatus(*m.Status)
	}

	return &model.ActivityLog{
		ID:         m.ID,
		TenantID:   tenantID,
		UserID:     m.UserID,
		Action:     model.ActivityType(m.Action),
		Status:     status,
		IPAddress:  getString(m.IPAddress),
		UserAgent:  getString(m.UserAgent),
		Metadata:   metadata,
		OccurredAt: m.OccurredAt,
		CreatedAt:  *m.CreatedAt,
	}
}

// getString 辅助函数：安全获取 *string 的值
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
