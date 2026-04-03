package listener

// AuditLogTask 审计日志任务
type AuditLogTask struct {
	Type   string                 `json:"type"`
	Action string                 `json:"action"`
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// GetPayload 获取任务数据（实现 event.Event 接口）
func (t *AuditLogTask) GetPayload() interface{} {
	return t.Data
}

// GetType 获取任务类型
func (t *AuditLogTask) GetType() string {
	return t.Type
}
