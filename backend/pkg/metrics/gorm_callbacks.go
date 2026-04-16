package metrics

import (
	"time"

	"gorm.io/gorm"
)

// SetupGORMCallbacks 为 GORM 注册指标收集 Callback
// 如果 m 为 nil 或 metricsEnabled 为 false，则不注册
func SetupGORMCallbacks(db *gorm.DB, m *Metrics, metricsEnabled bool) {
	if m == nil || !metricsEnabled {
		return
	}

	// Query 回调
	_ = db.Callback().Query().After("gorm:query").Register("metrics:query", func(db *gorm.DB) {
		if startTime, ok := db.InstanceGet("query_start_time"); ok {
			if start, ok := startTime.(time.Time); ok {
				duration := time.Since(start).Seconds()
				table := db.Statement.Table
				if table == "" {
					table = "unknown"
				}
				m.IncDBQuery("SELECT", table)
				m.ObserveDBQueryDuration("SELECT", table, duration)
			}
		}
	})

	// Create 回调
	_ = db.Callback().Create().After("gorm:create").Register("metrics:create", func(db *gorm.DB) {
		if startTime, ok := db.InstanceGet("query_start_time"); ok {
			if start, ok := startTime.(time.Time); ok {
				duration := time.Since(start).Seconds()
				table := db.Statement.Table
				if table == "" {
					table = "unknown"
				}
				m.IncDBQuery("CREATE", table)
				m.ObserveDBQueryDuration("CREATE", table, duration)
			}
		}
	})

	// Update 回调
	_ = db.Callback().Update().After("gorm:update").Register("metrics:update", func(db *gorm.DB) {
		if startTime, ok := db.InstanceGet("query_start_time"); ok {
			if start, ok := startTime.(time.Time); ok {
				duration := time.Since(start).Seconds()
				table := db.Statement.Table
				if table == "" {
					table = "unknown"
				}
				m.IncDBQuery("UPDATE", table)
				m.ObserveDBQueryDuration("UPDATE", table, duration)
			}
		}
	})

	// Delete 回调
	_ = db.Callback().Delete().After("gorm:delete").Register("metrics:delete", func(db *gorm.DB) {
		if startTime, ok := db.InstanceGet("query_start_time"); ok {
			if start, ok := startTime.(time.Time); ok {
				duration := time.Since(start).Seconds()
				table := db.Statement.Table
				if table == "" {
					table = "unknown"
				}
				m.IncDBQuery("DELETE", table)
				m.ObserveDBQueryDuration("DELETE", table, duration)
			}
		}
	})

	// Row 回调（COUNT 等）
	_ = db.Callback().Row().After("gorm:row").Register("metrics:row", func(db *gorm.DB) {
		if startTime, ok := db.InstanceGet("query_start_time"); ok {
			if start, ok := startTime.(time.Time); ok {
				duration := time.Since(start).Seconds()
				table := db.Statement.Table
				if table == "" {
					table = "unknown"
				}
				m.IncDBQuery("COUNT", table)
				m.ObserveDBQueryDuration("COUNT", table, duration)
			}
		}
	})

	// Before 回调：记录开始时间
	_ = db.Callback().Query().Before("gorm:query").Register("metrics:before:query", func(db *gorm.DB) {
		db.InstanceSet("query_start_time", time.Now())
	})
	_ = db.Callback().Create().Before("gorm:create").Register("metrics:before:create", func(db *gorm.DB) {
		db.InstanceSet("query_start_time", time.Now())
	})
	_ = db.Callback().Update().Before("gorm:update").Register("metrics:before:update", func(db *gorm.DB) {
		db.InstanceSet("query_start_time", time.Now())
	})
	_ = db.Callback().Delete().Before("gorm:delete").Register("metrics:before:delete", func(db *gorm.DB) {
		db.InstanceSet("query_start_time", time.Now())
	})
	_ = db.Callback().Row().Before("gorm:row").Register("metrics:before:row", func(db *gorm.DB) {
		db.InstanceSet("query_start_time", time.Now())
	})
}
