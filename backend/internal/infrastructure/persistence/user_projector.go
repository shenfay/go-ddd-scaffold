package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/pkg/util"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserProjectorImpl 用户投影器实现
type UserProjectorImpl struct {
	db *sql.DB
}

// NewUserProjector 创建用户投影器
func NewUserProjector(db *sql.DB) *UserProjectorImpl {
	return &UserProjectorImpl{db: db}
}

// Project 投影领域事件到读模型
func (p *UserProjectorImpl) Project(event ddd.DomainEvent) error {
	ctx := context.Background()

	switch e := event.(type) {
	case *user.UserRegisteredEvent:
		return p.handleUserRegistered(ctx, e)
	case *user.UserActivatedEvent:
		return p.handleUserActivated(ctx, e)
	case *user.UserDeactivatedEvent:
		return p.handleUserDeactivated(ctx, e)
	case *user.UserLoggedInEvent:
		return p.handleUserLoggedIn(ctx, e)
	case *user.UserPasswordChangedEvent:
		return p.handleUserPasswordChanged(ctx, e)
	case *user.UserEmailChangedEvent:
		return p.handleUserEmailChanged(ctx, e)
	case *user.UserLockedEvent:
		return p.handleUserLocked(ctx, e)
	case *user.UserUnlockedEvent:
		return p.handleUserUnlocked(ctx, e)
	case *user.UserProfileUpdatedEvent:
		return p.handleUserProfileUpdated(ctx, e)
	default:
		// 忽略未知事件类型
		return nil
	}
}

// handleUserRegistered 处理用户注册事件
func (p *UserProjectorImpl) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	query := `
		INSERT INTO user_read_model (
			user_id, username, email, status, created_at, updated_at, login_count
		) VALUES ($1, $2, $3, $4, NOW(), NOW(), 0)
	`

	_, err := p.db.ExecContext(ctx, query,
		event.UserID.Int64(),
		event.Username,
		event.Email,
		int(user.UserStatusPending),
	)

	return err
}

// handleUserActivated 处理用户激活事件
func (p *UserProjectorImpl) handleUserActivated(ctx context.Context, event *user.UserActivatedEvent) error {
	query := `UPDATE user_read_model SET status = $1, updated_at = NOW() WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, int(user.UserStatusActive), event.UserID.Int64())
	return err
}

// handleUserDeactivated 处理用户禁用事件
func (p *UserProjectorImpl) handleUserDeactivated(ctx context.Context, event *user.UserDeactivatedEvent) error {
	query := `UPDATE user_read_model SET status = $1, updated_at = NOW() WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, int(user.UserStatusInactive), event.UserID.Int64())
	return err
}

// handleUserLoggedIn 处理用户登录事件
func (p *UserProjectorImpl) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	query := `
		UPDATE user_read_model 
		SET last_login_at = $1, login_count = login_count + 1, updated_at = NOW()
		WHERE user_id = $2
	`

	_, err := p.db.ExecContext(ctx, query, event.LoginAt, event.UserID.Int64())
	return err
}

// handleUserPasswordChanged 处理密码修改事件
func (p *UserProjectorImpl) handleUserPasswordChanged(ctx context.Context, event *user.UserPasswordChangedEvent) error {
	query := `UPDATE user_read_model SET updated_at = NOW() WHERE user_id = $1`

	_, err := p.db.ExecContext(ctx, query, event.UserID.Int64())
	return err
}

// handleUserEmailChanged 处理邮箱修改事件
func (p *UserProjectorImpl) handleUserEmailChanged(ctx context.Context, event *user.UserEmailChangedEvent) error {
	query := `UPDATE user_read_model SET email = $1, updated_at = NOW() WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, event.NewEmail, event.UserID.Int64())
	return err
}

// handleUserLocked 处理用户锁定事件
func (p *UserProjectorImpl) handleUserLocked(ctx context.Context, event *user.UserLockedEvent) error {
	query := `UPDATE user_read_model SET status = $1, updated_at = NOW() WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, int(user.UserStatusLocked), event.UserID.Int64())
	return err
}

// handleUserUnlocked 处理用户解锁事件
func (p *UserProjectorImpl) handleUserUnlocked(ctx context.Context, event *user.UserUnlockedEvent) error {
	query := `UPDATE user_read_model SET status = $1, updated_at = NOW() WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, int(user.UserStatusActive), event.UserID.Int64())
	return err
}

// handleUserProfileUpdated 处理用户资料更新事件
func (p *UserProjectorImpl) handleUserProfileUpdated(ctx context.Context, event *user.UserProfileUpdatedEvent) error {
	// 简化实现，实际应该根据 UpdatedFields 更新对应字段
	query := `UPDATE user_read_model SET updated_at = NOW() WHERE user_id = $1`

	_, err := p.db.ExecContext(ctx, query, event.UserID.Int64())
	return err
}

// GetUserProfile 获取用户资料
func (p *UserProjectorImpl) GetUserProfile(ctx context.Context, userID user.UserID) (*user.UserProfileDTO, error) {
	query := `
		SELECT user_id, username, email, display_name, first_name, last_name,
		       gender, phone_number, avatar_url, status, created_at, updated_at,
		       last_login_at, login_count
		FROM user_read_model
		WHERE user_id = $1
	`

	var dto user.UserProfileDTO
	err := p.db.QueryRowContext(ctx, query, userID.Int64()).Scan(
		&dto.UserID, &dto.Username, &dto.Email, &dto.DisplayName,
		&dto.FirstName, &dto.LastName, &dto.Gender, &dto.PhoneNumber,
		&dto.AvatarURL, &dto.Status, &dto.CreatedAt, &dto.UpdatedAt,
		&dto.LastLoginAt, &dto.LoginCount,
	)

	if err == sql.ErrNoRows {
		return nil, ddd.ErrAggregateNotFound
	}
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

// ListUsers 列出用户
func (p *UserProjectorImpl) ListUsers(ctx context.Context, criteria user.UserListCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.UserListItemDTO], error) {
	countQuery := `SELECT COUNT(*) FROM user_read_model WHERE 1=1`
	query := `
		SELECT user_id, username, email, display_name, status, created_at, last_login_at
		FROM user_read_model WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if criteria.Status != nil {
		query += " AND status = $" + util.ToString(argPos)
		countQuery += " AND status = $" + util.ToString(argPos)
		args = append(args, int(*criteria.Status))
		argPos++
	}

	if criteria.Keyword != "" {
		query += " AND (username LIKE $" + util.ToString(argPos) + " OR email LIKE $" + util.ToString(argPos) + " OR display_name LIKE $" + util.ToString(argPos) + ")"
		countQuery += " AND (username LIKE $" + util.ToString(argPos) + " OR email LIKE $" + util.ToString(argPos) + " OR display_name LIKE $" + util.ToString(argPos) + ")"
		keyword := "%" + criteria.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
		argPos += 3
	}

	// 排序
	if criteria.SortBy != "" {
		query += " ORDER BY " + criteria.SortBy
		if criteria.SortDesc {
			query += " DESC"
		} else {
			query += " ASC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	// 分页
	query += " LIMIT $" + util.ToString(argPos) + " OFFSET $" + util.ToString(argPos+1)
	args = append(args, pagination.PageSize, pagination.Offset())

	// 查询总数
	var total int64
	err := p.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// 查询数据
	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*user.UserListItemDTO
	for rows.Next() {
		var item user.UserListItemDTO
		err := rows.Scan(
			&item.UserID, &item.Username, &item.Email, &item.DisplayName,
			&item.Status, &item.CreatedAt, &item.LastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return &ddd.PaginatedResult[*user.UserListItemDTO]{
		Items:      items,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.CalculateTotalPages(total),
	}, nil
}

// SearchUsers 搜索用户
func (p *UserProjectorImpl) SearchUsers(ctx context.Context, keyword string, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.UserListItemDTO], error) {
	criteria := user.UserListCriteria{
		Keyword: keyword,
		SortBy:  "created_at",
	}
	return p.ListUsers(ctx, criteria, pagination)
}

// GetUserStatistics 获取用户统计
func (p *UserProjectorImpl) GetUserStatistics(ctx context.Context) (*user.UserStatisticsDTO, error) {
	dto := &user.UserStatisticsDTO{}

	// 总用户数
	err := p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model").Scan(&dto.TotalUsers)
	if err != nil {
		return nil, err
	}

	// 各状态用户数
	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = $1", int(user.UserStatusActive)).Scan(&dto.ActiveUsers)
	if err != nil {
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = $1", int(user.UserStatusInactive)).Scan(&dto.InactiveUsers)
	if err != nil {
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = $1", int(user.UserStatusPending)).Scan(&dto.PendingUsers)
	if err != nil {
		return nil, err
	}

	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = $1", int(user.UserStatusLocked)).Scan(&dto.LockedUsers)
	if err != nil {
		return nil, err
	}

	// 今日登录数
	today := time.Now().Format("2006-01-02")
	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) = $1", today).Scan(&dto.TodayLogins)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 本周登录数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) >= $1", weekStart).Scan(&dto.ThisWeekLogins)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 本月登录数
	monthStart := time.Now().Format("2006-01") + "-01"
	err = p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) >= $1", monthStart).Scan(&dto.ThisMonthLogins)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return dto, nil
}

// ProcessUnprocessedEvents 处理未处理的事件
func (p *UserProjectorImpl) ProcessUnprocessedEvents(ctx context.Context, limit int) error {
	query := `
		SELECT id, aggregate_id, aggregate_type, event_type, event_version,
		       event_data, occurred_on, processed, created_at
		FROM domain_events
		WHERE processed = false
		ORDER BY id
		LIMIT $1
	`

	rows, err := p.db.QueryContext(ctx, query, limit)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var eventID int64
		var aggregateID int64
		var aggregateType, eventType, eventData string
		var eventVersion int
		var occurredOn time.Time
		var processed bool
		var createdAt time.Time

		err := rows.Scan(
			&eventID, &aggregateID, &aggregateType, &eventType, &eventVersion,
			&eventData, &occurredOn, &processed, &createdAt,
		)
		if err != nil {
			return err
		}

		// 反序列化事件
		event, err := p.deserializeEvent(eventType, eventData)
		if err != nil {
			return err
		}

		// 投影事件
		if err := p.Project(event); err != nil {
			return err
		}

		// 标记为已处理
		_, err = p.db.ExecContext(ctx, "UPDATE domain_events SET processed = true WHERE id = $1", eventID)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

// deserializeEvent 反序列化事件
func (p *UserProjectorImpl) deserializeEvent(eventType, eventData string) (ddd.DomainEvent, error) {
	switch eventType {
	case "UserRegistered":
		var event user.UserRegisteredEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserActivated":
		var event user.UserActivatedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserDeactivated":
		var event user.UserDeactivatedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserLoggedIn":
		var event user.UserLoggedInEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserPasswordChanged":
		var event user.UserPasswordChangedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserEmailChanged":
		var event user.UserEmailChangedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserLocked":
		var event user.UserLockedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserUnlocked":
		var event user.UserUnlockedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	case "UserProfileUpdated":
		var event user.UserProfileUpdatedEvent
		err := json.Unmarshal([]byte(eventData), &event)
		return &event, err
	default:
		return nil, ddd.NewBusinessError("UNKNOWN_EVENT", "unknown event type: "+eventType)
	}
}
