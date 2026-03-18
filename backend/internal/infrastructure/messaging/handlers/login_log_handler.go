package handlers

import (
	"context"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/loginlog"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/internal/infrastructure/snowflake"
	"github.com/shenfay/go-ddd-scaffold/shared/kernel"
)

// LoginLogHandler 登录日志事件处理器
type LoginLogHandler struct {
	repo      loginlog.LoginLogRepository
	snowflake *snowflake.Node
}

func NewLoginLogHandler(repo loginlog.LoginLogRepository, snowflake *snowflake.Node) *LoginLogHandler {
	return &LoginLogHandler{repo: repo, snowflake: snowflake}
}

// Handle 处理领域事件
func (h *LoginLogHandler) Handle(ctx context.Context, evt kernel.DomainEvent) error {
	switch e := evt.(type) {
	case *user.UserLoggedInEvent:
		return h.handleUserLoggedIn(ctx, e)
	default:
		return nil
	}
}

func (h *LoginLogHandler) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	// 解析 User-Agent（简化版本）
	deviceInfo := parseUserAgent(event.UserAgent)

	log := &loginlog.LoginLog{
		ID:           h.generateID(),
		UserID:       event.UserID.Int64(),
		LoginType:    "password",
		LoginStatus:  loginlog.LoginStatusSuccess,
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

// DeviceInfo 设备信息
type DeviceInfo struct {
	DeviceType string
	OS         string
	Browser    string
}

// parseUserAgent 解析 User-Agent（简化实现）
func parseUserAgent(ua string) DeviceInfo {
	info := DeviceInfo{
		DeviceType: "desktop",
		OS:         "Unknown",
		Browser:    "Unknown",
	}

	// 简单的关键词匹配
	if contains(ua, "Mobile") {
		info.DeviceType = "mobile"
	} else if contains(ua, "Tablet") {
		info.DeviceType = "tablet"
	}

	if contains(ua, "Windows") {
		info.OS = "Windows"
	} else if contains(ua, "Macintosh") {
		info.OS = "macOS"
	} else if contains(ua, "Android") {
		info.OS = "Android"
	} else if contains(ua, "iOS") || contains(ua, "iPhone") {
		info.OS = "iOS"
	}

	if contains(ua, "Chrome") {
		info.Browser = "Chrome"
	} else if contains(ua, "Safari") {
		info.Browser = "Safari"
	} else if contains(ua, "Firefox") {
		info.Browser = "Firefox"
	}

	return info
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != "" && substr != "" &&
		(len(s) >= len(substr)) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func (h *LoginLogHandler) generateID() int64 {
	if h.snowflake != nil {
		id, _ := h.snowflake.Generate()
		return id
	}
	return 0
}
