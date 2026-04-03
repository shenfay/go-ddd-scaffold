package activitylog

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/shenfay/go-ddd-scaffold/pkg/logger"
)

// Service 活动日志服务
type Service struct {
	repo   ActivityLogRepository
	client *asynq.Client // Asynq Client
}

// NewService 创建活动日志服务
func NewService(repo ActivityLogRepository, client *asynq.Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

// LogParams 记录日志的参数
type LogParams struct {
	UserID      string
	Email       string
	Action      ActivityType
	Status      ActivityStatus
	IP          string
	UserAgent   string
	Description string
	Metadata    map[string]interface{}
}

// Record 记录活动日志（同步版本）
func (s *Service) Record(ctx context.Context, params LogParams) error {
	// 解析 User-Agent
	device, browser, os := parseUserAgent(params.UserAgent)

	// 构建元数据 JSON
	var metadataJSON string
	if len(params.Metadata) > 0 {
		data, _ := json.Marshal(params.Metadata)
		metadataJSON = string(data)
	} else {
		metadataJSON = "{}" // 空对象而不是空字符串
	}

	log := &ActivityLog{
		UserID:      params.UserID,
		Email:       params.Email,
		Action:      params.Action,
		Status:      params.Status,
		IP:          params.IP,
		UserAgent:   params.UserAgent,
		Device:      device,
		Browser:     browser,
		OS:          os,
		Description: params.Description,
		Metadata:    metadataJSON,
		CreatedAt:   time.Now(),
	}

	return s.repo.Create(ctx, log)
}

// RecordAsync 异步记录活动日志（使用 Asynq 队列）
func (s *Service) RecordAsync(params LogParams) {
	// 解析 User-Agent
	device, browser, os := parseUserAgent(params.UserAgent)

	// 构建 payload
	payload := map[string]interface{}{
		"user_id":     params.UserID,
		"email":       params.Email,
		"action":      string(params.Action),
		"status":      string(params.Status),
		"ip":          params.IP,
		"user_agent":  params.UserAgent,
		"description": params.Description,
		"device":      device,
		"browser":     browser,
		"os":          os,
		"metadata":    params.Metadata,
	}

	taskPayload, _ := json.Marshal(payload)
	task := asynq.NewTask("activity:record", taskPayload)

	info, err := s.client.Enqueue(task)
	if err != nil {
		logger.Error("Failed to enqueue activity log: ", err)
		return
	}

	logger.Debug("✓ Enqueued activity log task: ID=", info.ID)
}

// Close 关闭服务（Asynq Client 在 main 中管理）
func (s *Service) Close() {
	// 无需操作
}

// GetUserLogs 获取用户的活动日志
func (s *Service) GetUserLogs(ctx context.Context, userID string, limit, offset int) ([]*ActivityLog, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.repo.FindByUserID(ctx, userID, limit, offset)
}

// GetRecentLogs 获取用户最近的活动日志
func (s *Service) GetRecentLogs(ctx context.Context, userID string, limit int) ([]*ActivityLog, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	return s.repo.FindRecent(ctx, userID, limit)
}

// parseUserAgent 解析 User-Agent 字符串
func parseUserAgent(ua string) (device, browser, os string) {
	ua = strings.ToLower(ua)

	// 判断设备类型
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") && strings.Contains(ua, "wv") {
		device = "mobile"
	} else if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		device = "tablet"
	} else {
		device = "desktop"
	}

	// 简单判断浏览器和操作系统
	switch {
	case strings.Contains(ua, "chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "safari"):
		browser = "Safari"
	case strings.Contains(ua, "edge"):
		browser = "Edge"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident"):
		browser = "IE"
	default:
		browser = "Other"
	}

	switch {
	case strings.Contains(ua, "windows"):
		os = "Windows"
	case strings.Contains(ua, "mac os") || strings.Contains(ua, "macos"):
		os = "macOS"
	case strings.Contains(ua, "linux"):
		os = "Linux"
	case strings.Contains(ua, "android"):
		os = "Android"
	case strings.Contains(ua, "ios"):
		os = "iOS"
	default:
		os = "Other"
	}

	return device, browser, os
}

// Middleware 活动日志中间件（用于自动记录）
func Middleware(service *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在请求结束后记录日志
		c.Next()

		// 只记录特定路径
		shouldLog := shouldRecordPath(c.Request.URL.Path)
		if !shouldLog {
			return
		}

		// 获取用户信息（如果已认证）
		userID, _ := c.Get("user_id")
		email, _ := c.Get("user_email")

		if userIDStr, ok := userID.(string); ok && userIDStr != "" {
			action := getActionFromPath(c.Request.URL.Path, c.Request.Method)
			status := getStatusFromCode(c.Writer.Status())

			// 使用异步方式记录，不阻塞请求
			service.RecordAsync(LogParams{
				UserID:      userIDStr,
				Email:       email.(string),
				Action:      action,
				Status:      status,
				IP:          c.ClientIP(),
				UserAgent:   c.GetHeader("User-Agent"),
				Description: strings.ToUpper(c.Request.Method) + " " + c.Request.URL.Path,
				Metadata: map[string]interface{}{
					"status_code": c.Writer.Status(),
					"method":      c.Request.Method,
					"path":        c.Request.URL.Path,
				},
			})

		}
	}
}

// shouldRecordPath 判断是否需要记录日志的路径
func shouldRecordPath(path string) bool {
	// 不记录健康检查和静态资源
	if strings.HasPrefix(path, "/health") ||
		strings.HasPrefix(path, "/swagger") ||
		strings.HasPrefix(path, "/static") {
		return false
	}
	return true
}

// getActionFromPath 根据路径和方法推断活动类型
func getActionFromPath(path, method string) ActivityType {
	switch {
	case strings.Contains(path, "/login"):
		return ActivityLogin
	case strings.Contains(path, "/logout"):
		return ActivityLogout
	case strings.Contains(path, "/register"):
		return ActivityRegister
	case strings.Contains(path, "/refresh"):
		return ActivityRefreshToken
	default:
		if method == http.MethodGet {
			return ActivityProfileUpdate // 查询操作
		}
		return ActivityProfileUpdate // 默认为资料更新
	}
}

// getStatusFromCode 根据 HTTP 状态码判断状态
func getStatusFromCode(code int) ActivityStatus {
	if code >= 200 && code < 400 {
		return ActivitySuccess
	}
	return ActivityFailed
}
