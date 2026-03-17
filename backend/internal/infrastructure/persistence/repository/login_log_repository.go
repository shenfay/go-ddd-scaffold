package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
)

type LoginLogRepositoryImpl struct {
	query *dao.Query
}

func NewLoginLogRepository(db *dao.Query) loginlog.LoginLogRepository {
	return &LoginLogRepositoryImpl{query: db}
}

func (r *LoginLogRepositoryImpl) Save(ctx context.Context, log *loginlog.LoginLog) error {
	daoModel := r.fromDomain(log)
	return r.query.LoginLog.WithContext(ctx).Create(daoModel)
}

func (r *LoginLogRepositoryImpl) FindByUserID(ctx context.Context, userID int64, limit int) ([]*loginlog.LoginLog, error) {
	daoModels, err := r.query.LoginLog.WithContext(ctx).
		Where(r.query.LoginLog.UserID.Eq(userID)).
		Order(r.query.LoginLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*loginlog.LoginLog, 0, len(daoModels))
	for _, m := range daoModels {
		logs = append(logs, r.toDomain(m))
	}
	return logs, nil
}

func (r *LoginLogRepositoryImpl) FindSuspiciousLogins(ctx context.Context, limit int) ([]*loginlog.LoginLog, error) {
	daoModels, err := r.query.LoginLog.WithContext(ctx).
		Where(r.query.LoginLog.IsSuspicious.Is(true)).
		Order(r.query.LoginLog.OccurredAt.Desc()).
		Limit(limit).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*loginlog.LoginLog, 0, len(daoModels))
	for _, m := range daoModels {
		logs = append(logs, r.toDomain(m))
	}
	return logs, nil
}

// fromDomain 将领域模型转换为 DAO 模型
func (r *LoginLogRepositoryImpl) fromDomain(log *loginlog.LoginLog) *model.LoginLog {
	return &model.LoginLog{
		ID:            log.ID,
		UserID:        log.UserID,
		TenantID:      log.TenantID,
		LoginType:     util.String(log.LoginType),
		LoginStatus:   log.LoginStatus,
		IPAddress:     log.IPAddress,
		UserAgent:     util.String(log.UserAgent),
		DeviceType:    util.String(log.DeviceType),
		OsInfo:        util.String(log.OSInfo),
		BrowserInfo:   util.String(log.BrowserInfo),
		Country:       util.String(log.Country),
		City:          util.String(log.City),
		FailureReason: util.String(log.FailureReason),
		IsSuspicious:  util.Bool(log.IsSuspicious),
		RiskScore:     util.Int32(int32(log.RiskScore)),
		SessionID:     util.String(log.SessionID),
		AccessTokenID: util.String(log.AccessTokenID),
		OccurredAt:    log.OccurredAt,
		CreatedAt:     util.Time(time.Now()),
	}
}

// toDomain 将 DAO 模型转换为领域模型
func (r *LoginLogRepositoryImpl) toDomain(m *model.LoginLog) *loginlog.LoginLog {
	return &loginlog.LoginLog{
		ID:            m.ID,
		UserID:        m.UserID,
		TenantID:      m.TenantID,
		LoginType:     util.StringValue(m.LoginType),
		LoginStatus:   m.LoginStatus,
		IPAddress:     m.IPAddress,
		UserAgent:     util.StringValue(m.UserAgent),
		DeviceType:    util.StringValue(m.DeviceType),
		OSInfo:        util.StringValue(m.OsInfo),
		BrowserInfo:   util.StringValue(m.BrowserInfo),
		Country:       util.StringValue(m.Country),
		City:          util.StringValue(m.City),
		FailureReason: util.StringValue(m.FailureReason),
		IsSuspicious:  util.BoolValue(m.IsSuspicious),
		RiskScore:     int(util.Int32Value(m.RiskScore)),
		SessionID:     util.StringValue(m.SessionID),
		AccessTokenID: util.StringValue(m.AccessTokenID),
		OccurredAt:    m.OccurredAt,
		CreatedAt:     util.TimeValue(m.CreatedAt),
	}
}
