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
	Status      string                 `json:"status"`
	UserID      string                 `json:"user_id"`
	Email       string                 `json:"email"`
	IP          string                 `json:"ip"`
	UserAgent   string                 `json:"user_agent"`
	Device      string                 `json:"device"`
	Browser     string                 `json:"browser"`
	OS          string                 `json:"os"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GetPayload 获取任务数据
func (t *ActivityLogTask) GetPayload() interface{} {
	return map[string]interface{}{
		"user_id":     t.UserID,
		"email":       t.Email,
		"action":      t.Action,
		"status":      t.Status,
		"ip":          t.IP,
		"user_agent":  t.UserAgent,
		"device":      t.Device,
		"browser":     t.Browser,
		"os":          t.OS,
		"description": t.Description,
		"metadata":    t.Metadata,
	}
}

// GetType 获取任务类型
func (t *ActivityLogTask) GetType() string {
	return t.Type
}
