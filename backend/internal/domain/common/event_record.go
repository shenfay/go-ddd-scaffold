package common

// EventRecord 事件记录（从数据库读取）
type EventRecord struct {
	ID            int64
	AggregateID   string
	AggregateType string
	EventType     string
	EventVersion  int32
	EventData     string
	OccurredOn    string
	Metadata      *string
}
