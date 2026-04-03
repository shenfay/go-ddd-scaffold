package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	// TypeActivityLogRecord 活动日志记录任务类型
	TypeActivityLogRecord = "activity:record"
)

// ActivityLogPayload 活动日志任务载荷
// 使用 string 避免循环依赖
type ActivityLogPayload struct {
	UserID      string                 `json:"user_id"`
	Email       string                 `json:"email"`
	Action      string                 `json:"action"`
	Status      string                 `json:"status"`
	IP          string                 `json:"ip"`
	UserAgent   string                 `json:"user_agent"`
	Description string                 `json:"description"`
	Device      string                 `json:"device"`
	Browser     string                 `json:"browser"`
	OS          string                 `json:"os"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewActivityLogRecordTask 创建活动日志记录任务
func NewActivityLogRecordTask(payload ActivityLogPayload) (*asynq.Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeActivityLogRecord, b), nil
}
