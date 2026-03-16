package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/audit"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"github.com/spf13/cast"
)

type AuditLogRepositoryImpl struct {
}

func NewAuditLogRepository(db interface{}) audit.AuditLogRepository {
	return &AuditLogRepositoryImpl{}
}

func (r *AuditLogRepositoryImpl) Save(ctx context.Context, log *audit.AuditLog) error {
	daoModel := r.fromDomain(log)
	return dao.AuditLog.WithContext(ctx).Create(daoModel)
}

func (r *AuditLogRepositoryImpl) FindByUserID(ctx context.Context, userID int64, limit int) ([]*audit.AuditLog, error) {
	daoModels, err := dao.AuditLog.WithContext(ctx).
		Where(dao.AuditLog.UserID.Eq(userID)).
		Order(dao.AuditLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*audit.AuditLog, 0, len(daoModels))
	for _, m := range daoModels {
		logs = append(logs, r.toDomain(m))
	}
	return logs, nil
}

func (r *AuditLogRepositoryImpl) FindByAction(ctx context.Context, action string, limit int) ([]*audit.AuditLog, error) {
	daoModels, err := dao.AuditLog.WithContext(ctx).
		Where(dao.AuditLog.Action.Eq(action)).
		Order(dao.AuditLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*audit.AuditLog, 0, len(daoModels))
	for _, m := range daoModels {
		logs = append(logs, r.toDomain(m))
	}
	return logs, nil
}

// fromDomain 将领域模型转换为 DAO 模型
func (r *AuditLogRepositoryImpl) fromDomain(log *audit.AuditLog) *model.AuditLog {
	return &model.AuditLog{
		ID:           log.ID,
		TenantID:     log.TenantID,
		UserID:       util.Int64PtrNilIfZero(log.UserID),
		Action:       log.Action,
		ResourceType: util.StringPtrNilIfEmpty(log.ResourceType),
		ResourceID:   log.ResourceID,
		RequestID:    util.StringPtrNilIfEmpty(log.RequestID),
		IPAddress:    util.StringPtrNilIfEmpty(log.IPAddress),
		UserAgent:    util.StringPtrNilIfEmpty(log.UserAgent),
		Metadata:     r.metadataToString(log.Metadata),
		Status:       util.Int16PtrNilIfZero(cast.ToInt16(log.Status)),
		ErrorMessage: util.StringPtrNilIfEmpty(log.ErrorMessage),
		OccurredAt:   log.OccurredAt,
		CreatedAt:    util.Time(time.Now()),
	}
}

// toDomain 将 DAO 模型转换为领域模型
func (r *AuditLogRepositoryImpl) toDomain(m *model.AuditLog) *audit.AuditLog {
	return &audit.AuditLog{
		ID:       m.ID,
		TenantID: m.TenantID,
		UserID: func() int64 {
			if m.UserID != nil {
				return *m.UserID
			}
			return 0
		}(),
		Action: m.Action,
		ResourceType: func() string {
			if m.ResourceType != nil {
				return *m.ResourceType
			}
			return ""
		}(),
		ResourceID: m.ResourceID,
		RequestID: func() string {
			if m.RequestID != nil {
				return *m.RequestID
			}
			return ""
		}(),
		IPAddress: func() string {
			if m.IPAddress != nil {
				return *m.IPAddress
			}
			return ""
		}(),
		UserAgent: func() string {
			if m.UserAgent != nil {
				return *m.UserAgent
			}
			return ""
		}(),
		Metadata: r.stringToMetadata(m.Metadata),
		Status: func() int16 {
			if m.Status != nil {
				return *m.Status
			}
			return 0
		}(),
		ErrorMessage: func() string {
			if m.ErrorMessage != nil {
				return *m.ErrorMessage
			}
			return ""
		}(),
		OccurredAt: m.OccurredAt,
		CreatedAt: func() time.Time {
			if m.CreatedAt != nil {
				return *m.CreatedAt
			}
			return time.Time{}
		}(),
	}
}

func (r *AuditLogRepositoryImpl) metadataToString(m map[string]interface{}) *string {
	if len(m) == 0 {
		return nil
	}
	data, _ := json.Marshal(m)
	s := string(data)
	return &s
}

func (r *AuditLogRepositoryImpl) stringToMetadata(s *string) map[string]interface{} {
	if s == nil || *s == "" {
		return make(map[string]interface{})
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(*s), &m)
	return m
}
