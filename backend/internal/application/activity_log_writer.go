package application

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/shared/aggregate"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
	"go.uber.org/zap"
)

// ActivityLogWriter 活动日志写入器
// 职责：为应用层提供便捷的活动日志写入功能
// 使用场景：UseCase 在事务内直接调用，同步写入活动日志
type ActivityLogWriter struct {
	repo   aggregate.ActivityLogRepository
	logger *zap.Logger
}

// NewActivityLogWriter 创建活动日志写入器
func NewActivityLogWriter(repo aggregate.ActivityLogRepository, logger *zap.Logger) *ActivityLogWriter {
	if logger == nil {
		logger = zap.L().Named("activity_log_writer")
	}
	return &ActivityLogWriter{
		repo:   repo,
		logger: logger,
	}
}

// Write 写入活动日志
// 参数说明：
// - userID: 用户 ID
// - action: 活动类型（如 aggregate.ActivityUserRegistered）
// - status: 活动状态（成功/失败）
// - metadata: 扩展元数据（可选）
// - opts: 可选配置（IP、User-Agent 等）
func (w *ActivityLogWriter) Write(
	ctx context.Context,
	userID int64,
	action aggregate.ActivityType,
	status aggregate.ActivityStatus,
	metadata map[string]interface{},
	opts ...ActivityLogOption,
) error {
	// 应用默认配置
	options := DefaultActivityLogOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}

	// 处理元数据
	var finalMetadata map[string]interface{}
	if metadata != nil {
		finalMetadata = metadata
	} else {
		finalMetadata = make(map[string]interface{})
	}

	// 添加额外信息到元数据
	if options.IPAddress != "" {
		finalMetadata["ip_address"] = options.IPAddress
	}
	if options.UserAgent != "" {
		finalMetadata["user_agent"] = options.UserAgent
	}
	if options.TenantID != nil {
		finalMetadata["tenant_id"] = *options.TenantID
	}

	// 创建活动日志
	log := &aggregate.ActivityLog{
		ID:         idgen.Generate(),
		UserID:     userID,
		Action:     action,
		Status:     status,
		Metadata:   finalMetadata,
		OccurredAt: time.Now(),
		CreatedAt:  time.Now(),
	}

	// 记录日志（调试级别）
	w.logger.Debug("Writing activity log",
		zap.String("action", string(action)),
		zap.Int64("user_id", userID),
		zap.Any("metadata", finalMetadata),
	)

	// 保存到数据库
	return w.repo.Save(ctx, log)
}

// WriteSuccess 写入成功的活动日志（快捷方法）
func (w *ActivityLogWriter) WriteSuccess(
	ctx context.Context,
	userID int64,
	action aggregate.ActivityType,
	metadata map[string]interface{},
	opts ...ActivityLogOption,
) error {
	return w.Write(ctx, userID, action, aggregate.ActivityStatusSuccess, metadata, opts...)
}

// WriteFailure 写入失败的活动日志（快捷方法）
func (w *ActivityLogWriter) WriteFailure(
	ctx context.Context,
	userID int64,
	action aggregate.ActivityType,
	errorMsg string,
	opts ...ActivityLogOption,
) error {
	metadata := map[string]interface{}{
		"error": errorMsg,
	}
	return w.Write(ctx, userID, action, aggregate.ActivityStatusFailed, metadata, opts...)
}

// ActivityLogOption 活动日志配置选项
type ActivityLogOption interface {
	Apply(*activityLogOptions)
}

type activityLogOptions struct {
	IPAddress string
	UserAgent string
	TenantID  *int64
}

// DefaultActivityLogOptions 返回默认配置
func DefaultActivityLogOptions() *activityLogOptions {
	return &activityLogOptions{}
}

// WithIPAddress 设置 IP 地址
func WithIPAddress(ip string) ActivityLogOption {
	return &ipAddressOption{ip: ip}
}

type ipAddressOption struct {
	ip string
}

func (o *ipAddressOption) Apply(opts *activityLogOptions) {
	opts.IPAddress = o.ip
}

// WithUserAgent 设置 User-Agent
func WithUserAgent(ua string) ActivityLogOption {
	return &userAgentOption{ua: ua}
}

type userAgentOption struct {
	ua string
}

func (o *userAgentOption) Apply(opts *activityLogOptions) {
	opts.UserAgent = o.ua
}

// WithTenantID 设置租户 ID
func WithTenantID(tenantID int64) ActivityLogOption {
	return &tenantIDOption{tenantID: tenantID}
}

type tenantIDOption struct {
	tenantID int64
}

func (o *tenantIDOption) Apply(opts *activityLogOptions) {
	opts.TenantID = &o.tenantID
}
