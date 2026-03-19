package loginlog

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/kernel"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
)

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

// EventHandler 登录日志领域事件处理器
type EventHandler struct {
	repo        LoginLogRepository
	idGenerator IDGenerator
	uaParser    UserAgentParser
}

// IDGenerator ID生成器接口
type IDGenerator interface {
	Generate() (int64, error)
}

// NewEventHandler 创建登录日志事件处理器
func NewEventHandler(repo LoginLogRepository, idGenerator IDGenerator, uaParser UserAgentParser) *EventHandler {
	return &EventHandler{
		repo:        repo,
		idGenerator: idGenerator,
		uaParser:    uaParser,
	}
}

// Handle 处理领域事件
func (h *EventHandler) Handle(ctx context.Context, event kernel.DomainEvent) error {
	switch e := event.(type) {
	case *user.UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	default:
		return nil // 忽略不关心的事件
	}
}

func (h *EventHandler) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	// 解析 User-Agent
	deviceInfo := DeviceInfo{DeviceType: "desktop", OS: "Unknown", Browser: "Unknown"}
	if h.uaParser != nil {
		deviceInfo = h.uaParser.Parse(event.UserAgent)
	}

	log := &LoginLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		LoginType:    "password",
		LoginStatus:  LoginStatusSuccess,
		IPAddress:    event.IPAddress,
		UserAgent:    event.UserAgent,
		DeviceType:   deviceInfo.DeviceType,
		OSInfo:       deviceInfo.OS,
		BrowserInfo:  deviceInfo.Browser,
		IsSuspicious: false, // TODO: 实现风控检测
		RiskScore:    0,
		OccurredAt:   event.LoginAt,
	}

	return h.repo.Save(ctx, log)
}

func (h *EventHandler) generateID() int64 {
	if h.idGenerator != nil {
		id, _ := h.idGenerator.Generate()
		return id
	}
	return 0
}
