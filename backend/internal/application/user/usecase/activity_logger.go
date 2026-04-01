package usecase

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/event"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/model"
	vo "github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	idgen "github.com/shenfay/go-ddd-scaffold/internal/infrastructure/platform/idgen"
)

// activityLogger 活动日志记录器 (辅助工具)
// 用于简化 UseCase 中的 ActivityLog 创建和保存逻辑
type activityLogger struct {
	repo model.ActivityLogRepository
}

// newActivityLogger 创建活动日志记录器
func newActivityLogger(repo model.ActivityLogRepository) *activityLogger {
	return &activityLogger{repo: repo}
}

// LogUserAction 记录用户操作日志
func (l *activityLogger) LogUserAction(ctx context.Context, userID vo.UserID, action model.ActivityType, metadata map[string]interface{}) error {
	log := &model.ActivityLog{
		ID:         idgen.Generate(),
		UserID:     userID.Int64(),
		Action:     action,
		Status:     event.ActivityStatusSuccess,
		Metadata:   metadata,
		OccurredAt: time.Now(),
		CreatedAt:  time.Now(),
	}
	return l.repo.Save(ctx, log)
}

// LogUserActionWithResult 记录用户操作结果
func (l *activityLogger) LogUserActionWithResult(ctx context.Context, userID vo.UserID, action model.ActivityType, status model.ActivityStatus, metadata map[string]interface{}) error {
	log := &model.ActivityLog{
		ID:         idgen.Generate(),
		UserID:     userID.Int64(),
		Action:     action,
		Status:     status,
		Metadata:   metadata,
		OccurredAt: time.Now(),
		CreatedAt:  time.Now(),
	}
	return l.repo.Save(ctx, log)
}
