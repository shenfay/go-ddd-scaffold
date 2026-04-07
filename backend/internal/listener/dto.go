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

// ActivityLogTask 活动日志任务
type ActivityLogTask struct {
	Type        string                 `json:"type"`
	Action      string                 `json:"action"`
	UserID      string                 `json:"user_id"`
	Email       string                 `json:"email"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GetPayload 获取任务数据
func (t *ActivityLogTask) GetPayload() interface{} {
	return map[string]interface{}{
		"user_id":     t.UserID,
		"email":       t.Email,
		"action":      t.Action,
		"description": t.Description,
		"metadata":    t.Metadata,
	}
}

// GetType 获取任务类型
func (t *ActivityLogTask) GetType() string {
	return t.Type
}
