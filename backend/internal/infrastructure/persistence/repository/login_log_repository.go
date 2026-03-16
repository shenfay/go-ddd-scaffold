package repository

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/dao"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/persistence/model"
)

type LoginLogRepositoryImpl struct {
}

func NewLoginLogRepository(db interface{}) loginlog.LoginLogRepository {
	return &LoginLogRepositoryImpl{}
}

func (r *LoginLogRepositoryImpl) Save(ctx context.Context, log *loginlog.LoginLog) error {
	daoModel := r.fromDomain(log)
	return dao.LoginLog.WithContext(ctx).Create(daoModel)
}

func (r *LoginLogRepositoryImpl) FindByUserID(ctx context.Context, userID int64, limit int) ([]*loginlog.LoginLog, error) {
	daoModels, err := dao.LoginLog.WithContext(ctx).
		Where(dao.LoginLog.UserID.Eq(userID)).
		Order(dao.LoginLog.OccurredAt.Desc()).
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
	daoModels, err := dao.LoginLog.WithContext(ctx).
		Where(dao.LoginLog.IsSuspicious.Is(true)).
		Order(dao.LoginLog.OccurredAt.Desc()).
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
	loginType := log.LoginType
	userAgent := log.UserAgent
	deviceType := log.DeviceType
	osInfo := log.OSInfo
	browserInfo := log.BrowserInfo
	country := log.Country
	city := log.City
	failureReason := log.FailureReason
	sessionID := log.SessionID
	accessTokenID := log.AccessTokenID

	return &model.LoginLog{
		ID:       log.ID,
		UserID:   log.UserID,
		TenantID: log.TenantID,
		LoginType: func() *string {
			if loginType != "" {
				return &loginType
			}
			return nil
		}(),
		LoginStatus: log.LoginStatus,
		IPAddress:   log.IPAddress,
		UserAgent: func() *string {
			if userAgent != "" {
				return &userAgent
			}
			return nil
		}(),
		DeviceType: func() *string {
			if deviceType != "" {
				return &deviceType
			}
			return nil
		}(),
		OsInfo: func() *string {
			if osInfo != "" {
				return &osInfo
			}
			return nil
		}(),
		BrowserInfo: func() *string {
			if browserInfo != "" {
				return &browserInfo
			}
			return nil
		}(),
		Country: func() *string {
			if country != "" {
				return &country
			}
			return nil
		}(),
		City: func() *string {
			if city != "" {
				return &city
			}
			return nil
		}(),
		FailureReason: func() *string {
			if failureReason != "" {
				return &failureReason
			}
			return nil
		}(),
		IsSuspicious: func() *bool { v := log.IsSuspicious; return &v }(),
		RiskScore:    func() *int32 { v := int32(log.RiskScore); return &v }(),
		SessionID: func() *string {
			if sessionID != "" {
				return &sessionID
			}
			return nil
		}(),
		AccessTokenID: func() *string {
			if accessTokenID != "" {
				return &accessTokenID
			}
			return nil
		}(),
		OccurredAt: log.OccurredAt,
		CreatedAt:  func() *time.Time { t := time.Now(); return &t }(),
	}
}

// toDomain 将 DAO 模型转换为领域模型
func (r *LoginLogRepositoryImpl) toDomain(m *model.LoginLog) *loginlog.LoginLog {
	return &loginlog.LoginLog{
		ID:       m.ID,
		UserID:   m.UserID,
		TenantID: m.TenantID,
		LoginType: func() string {
			if m.LoginType != nil {
				return *m.LoginType
			}
			return ""
		}(),
		LoginStatus: m.LoginStatus,
		IPAddress:   m.IPAddress,
		UserAgent: func() string {
			if m.UserAgent != nil {
				return *m.UserAgent
			}
			return ""
		}(),
		DeviceType: func() string {
			if m.DeviceType != nil {
				return *m.DeviceType
			}
			return ""
		}(),
		OSInfo: func() string {
			if m.OsInfo != nil {
				return *m.OsInfo
			}
			return ""
		}(),
		BrowserInfo: func() string {
			if m.BrowserInfo != nil {
				return *m.BrowserInfo
			}
			return ""
		}(),
		Country: func() string {
			if m.Country != nil {
				return *m.Country
			}
			return ""
		}(),
		City: func() string {
			if m.City != nil {
				return *m.City
			}
			return ""
		}(),
		FailureReason: func() string {
			if m.FailureReason != nil {
				return *m.FailureReason
			}
			return ""
		}(),
		IsSuspicious: func() bool {
			if m.IsSuspicious != nil {
				return *m.IsSuspicious
			}
			return false
		}(),
		RiskScore: func() int {
			if m.RiskScore != nil {
				return int(*m.RiskScore)
			}
			return 0
		}(),
		SessionID: func() string {
			if m.SessionID != nil {
				return *m.SessionID
			}
			return ""
		}(),
		AccessTokenID: func() string {
			if m.AccessTokenID != nil {
				return *m.AccessTokenID
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
