package activitylog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Service 活动日志服务
type Service struct {
	repo       ActivityLogRepository
	queue      chan LogParams // 异步队列
	wg         sync.WaitGroup // 优雅关闭
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewService 创建活动日志服务
func NewService(repo ActivityLogRepository) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Service{
		repo:   repo,
		queue:  make(chan LogParams, 100), // 缓冲队列长度 100
		ctx:    ctx,
		cancel: cancel,
	}
	
	// 启动后台协程处理异步写入
	go s.processQueue()
	
	return s
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
	
	// 调试日志
	fmt.Printf("[DEBUG] Recording activity log: UserID=%s, Action=%s, Status=%s\n", 
		log.UserID, log.Action, log.Status)
	
	return s.repo.Create(ctx, log)
}

// RecordAsync 异步记录活动日志（非阻塞）
func (s *Service) RecordAsync(params LogParams) {
	select {
	case s.queue <- params:
		// 成功加入队列
		fmt.Printf("[DEBUG] Activity log queued: UserID=%s, Action=%s\n", 
			params.UserID, params.Action)
	default:
		// 队列已满，降级为同步写入或直接丢弃（根据业务需求）
		// 这里选择同步写入作为降级策略
		fmt.Println("[DEBUG] Activity log queue full, falling back to sync write")
		_ = s.Record(context.Background(), params)
	}
}

// processQueue 后台处理队列中的日志
func (s *Service) processQueue() {
	const batchSize = 10
	const flushInterval = 2 * time.Second
	
	timer := time.NewTicker(flushInterval)
	defer timer.Stop()
	
	batch := make([]LogParams, 0, batchSize)
	
	for {
		select {
		case <-s.ctx.Done():
			// 上下文取消，处理剩余日志后退出
			if len(batch) > 0 {
				s.flushBatch(batch)
			}
			return
			
		case params := <-s.queue:
			batch = append(batch, params)
			if len(batch) >= batchSize {
				s.flushBatch(batch)
				batch = batch[:0]
			}
			
		case <-timer.C:
			if len(batch) > 0 {
				s.flushBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// flushBatch 批量写入日志
func (s *Service) flushBatch(batch []LogParams) {
	if len(batch) == 0 {
		return
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	fmt.Printf("[DEBUG] Flushing batch of %d activity logs\n", len(batch))
	
	for _, params := range batch {
		// 忽略单个错误，继续处理其他日志
		if err := s.Record(ctx, params); err != nil {
			fmt.Printf("[ERROR] Failed to record activity log: %v\n", err)
		}
	}
}

// Close 优雅关闭服务
func (s *Service) Close() {
	s.cancel()
	s.wg.Wait()
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
			
			fmt.Printf("[DEBUG] Middleware capturing activity: UserID=%s, Path=%s, Status=%d\n",
				userIDStr, c.Request.URL.Path, c.Writer.Status())
			
			// 使用异步方式记录，不阻塞请求
			service.RecordAsync(LogParams{
				UserID:      userIDStr,
				Email:       email.(string),
				Action:      action,
				Status:      status,
				IP:          c.ClientIP(),
				UserAgent:   c.GetHeader("User-Agent"),
				Description: fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
				Metadata: map[string]interface{}{
					"status_code": c.Writer.Status(),
					"method":      c.Request.Method,
					"path":        c.Request.URL.Path,
				},
			})
		} else {
			fmt.Printf("[DEBUG] No user found in context, skipping activity log\n")
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
