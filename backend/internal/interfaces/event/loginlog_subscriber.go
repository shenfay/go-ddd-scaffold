package event

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	userEvent "github.com/shenfay/go-ddd-scaffold/internal/domain/user/event"
	logger "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/logging"
)

// LoginLogSubscriber 登录日志事件订阅者
// 负责监听用户登录事件并记录登录日志
type LoginLogSubscriber struct {
	repo        logger.LoginLogRepository
	idGenerator IDGenerator
	uaParser    UserAgentParser
}

// UserAgentParser User-Agent解析器接口
type UserAgentParser interface {
	Parse(ua string) DeviceInfo
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceType string
	OS         string
	Browser    string
}

// NewLoginLogSubscriber 创建登录日志事件订阅者
func NewLoginLogSubscriber(repo logger.LoginLogRepository, idGenerator IDGenerator, uaParser UserAgentParser) *LoginLogSubscriber {
	return &LoginLogSubscriber{
		repo:        repo,
		idGenerator: idGenerator,
		uaParser:    uaParser,
	}
}

// Handle 处理领域事件
func (s *LoginLogSubscriber) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *userEvent.UserLoggedInEvent:
		return s.handleUserLoggedIn(ctx, e)
	default:
		return nil // 忽略不关心的事件
	}
}

func (s *LoginLogSubscriber) handleUserLoggedIn(ctx context.Context, event *userEvent.UserLoggedInEvent) error {
	// 解析 User-Agent
	deviceInfo := DeviceInfo{DeviceType: "desktop", OS: "Unknown", Browser: "Unknown"}
	if s.uaParser != nil {
		deviceInfo = s.uaParser.Parse(event.UserAgent)
	}

	log := &logger.LoginLog{
		ID:           s.generateID(),
		UserID:       event.UserID.Int64(),
		LoginType:    "password",
		LoginStatus:  logger.LoginStatusSuccess,
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		DeviceType:   deviceInfo.DeviceType,
		OSInfo:       deviceInfo.OS,
		BrowserInfo:  deviceInfo.Browser,
		IsSuspicious: false, // TODO: 实现风控检测
		RiskScore:    0,
		OccurredAt:   event.LoginAt,
	}

	return s.repo.Save(ctx, log)
}

func (s *LoginLogSubscriber) generateID() int64 {
	if s.idGenerator != nil {
		id, _ := s.idGenerator.Generate()
		return id
	}
	return 0
}
