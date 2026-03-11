package persistence

import (
	"context"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserReadModelImpl 用户读模型实现
type UserReadModelImpl struct {
	db DB
}

// NewUserReadModel 创建用户读模型
func NewUserReadModel(db DB) user.UserReadModel {
	return &UserReadModelImpl{db: db}
}

// GetUserProfile 获取用户资料
func (r *UserReadModelImpl) GetUserProfile(ctx context.Context, userID user.UserID) (*user.UserProfileDTO, error) {
	var dto user.UserProfileDTO
	err := r.db.QueryRow(ctx,
		`SELECT user_id, username, email, display_name, first_name, last_name, 
		gender, phone_number, avatar_url, status, created_at, updated_at, last_login_at, login_count 
		FROM user_read_model WHERE user_id = ?`,
		userID.Int64(),
	).Scan(
		&dto.UserID, &dto.Username, &dto.Email, &dto.DisplayName, &dto.FirstName,
		&dto.LastName, &dto.Gender, &dto.PhoneNumber, &dto.AvatarURL, &dto.Status,
		&dto.CreatedAt, &dto.UpdatedAt, &dto.LastLoginAt, &dto.LoginCount,
	)

	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	return &dto, nil
}

// ListUsers 列出用户
func (r *UserReadModelImpl) ListUsers(ctx context.Context, criteria user.UserListCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.UserListItemDTO], error) {
	query := `SELECT user_id, username, email, display_name, status, created_at, last_login_at 
		FROM user_read_model WHERE 1=1`
	countQuery := "SELECT COUNT(*) FROM user_read_model WHERE 1=1"
	var args []interface{}

	if criteria.Status != nil {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, int(*criteria.Status))
	}

	if criteria.Keyword != "" {
		query += " AND (username LIKE ? OR email LIKE ? OR display_name LIKE ?)"
		countQuery += " AND (username LIKE ? OR email LIKE ? OR display_name LIKE ?)"
		keyword := "%" + criteria.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	// 排序
	sortField := "created_at"
	if criteria.SortBy != "" {
		sortField = criteria.SortBy
	}
	sortOrder := "ASC"
	if criteria.SortDesc {
		sortOrder = "DESC"
	}
	query += " ORDER BY " + sortField + " " + sortOrder

	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, pagination.Limit(), pagination.Offset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*user.UserListItemDTO
	for rows.Next() {
		var dto user.UserListItemDTO
		err := rows.Scan(
			&dto.UserID, &dto.Username, &dto.Email, &dto.DisplayName,
			&dto.Status, &dto.CreatedAt, &dto.LastLoginAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &dto)
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
func (r *UserReadModelImpl) SearchUsers(ctx context.Context, keyword string, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.UserListItemDTO], error) {
	criteria := user.UserListCriteria{
		Keyword: keyword,
	}
	return r.ListUsers(ctx, criteria, pagination)
}

// GetUserStatistics 获取用户统计
func (r *UserReadModelImpl) GetUserStatistics(ctx context.Context) (*user.UserStatisticsDTO, error) {
	var dto user.UserStatisticsDTO

	// 总用户数
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_read_model").Scan(&dto.TotalUsers)
	if err != nil {
		return nil, err
	}

	// 各状态用户数
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = ?", int(user.UserStatusActive)).Scan(&dto.ActiveUsers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = ?", int(user.UserStatusInactive)).Scan(&dto.InactiveUsers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = ?", int(user.UserStatusPending)).Scan(&dto.PendingUsers)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM user_read_model WHERE status = ?", int(user.UserStatusLocked)).Scan(&dto.LockedUsers)
	if err != nil {
		return nil, err
	}

	// 今日登录数
	today := time.Now().Format("2006-01-02")
	err = r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) = ?",
		today,
	).Scan(&dto.TodayLogins)
	if err != nil {
		return nil, err
	}

	// 本周登录数
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Format("2006-01-02")
	err = r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) >= ?",
		weekStart,
	).Scan(&dto.ThisWeekLogins)
	if err != nil {
		return nil, err
	}

	// 本月登录数
	monthStart := time.Now().Format("2006-01") + "-01"
	err = r.db.QueryRow(ctx,
		"SELECT COUNT(*) FROM user_read_model WHERE DATE(last_login_at) >= ?",
		monthStart,
	).Scan(&dto.ThisMonthLogins)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

// UserProjector 用户投影器
type UserProjector struct {
	db DB
}

// NewUserProjector 创建用户投影器
func NewUserProjector(db DB) *UserProjector {
	return &UserProjector{db: db}
}

// Project 投影领域事件到读模型
func (p *UserProjector) Project(ctx context.Context, event ddd.DomainEvent) error {
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
		return nil
	}
}

// handleUserRegistered 处理用户注册事件
func (p *UserProjector) handleUserRegistered(ctx context.Context, event *user.UserRegisteredEvent) error {
	_, err := p.db.Exec(ctx,
		`INSERT INTO user_read_model (user_id, username, email, status, created_at, updated_at, login_count) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.UserID.Int64(),
		event.Username,
		event.Email,
		int(user.UserStatusPending),
		event.RegisteredAt,
		event.RegisteredAt,
		0,
	)
	return err
}

// handleUserActivated 处理用户激活事件
func (p *UserProjector) handleUserActivated(ctx context.Context, event *user.UserActivatedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET status = ?, updated_at = ? WHERE user_id = ?",
		int(user.UserStatusActive),
		event.ActivatedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserDeactivated 处理用户禁用事件
func (p *UserProjector) handleUserDeactivated(ctx context.Context, event *user.UserDeactivatedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET status = ?, updated_at = ? WHERE user_id = ?",
		int(user.UserStatusInactive),
		event.DeactivatedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserLoggedIn 处理用户登录事件
func (p *UserProjector) handleUserLoggedIn(ctx context.Context, event *user.UserLoggedInEvent) error {
	_, err := p.db.Exec(ctx,
		`UPDATE user_read_model SET last_login_at = ?, login_count = login_count + 1, updated_at = ? WHERE user_id = ?`,
		event.LoginAt,
		event.LoginAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserPasswordChanged 处理密码修改事件
func (p *UserProjector) handleUserPasswordChanged(ctx context.Context, event *user.UserPasswordChangedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET updated_at = ? WHERE user_id = ?",
		event.ChangedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserEmailChanged 处理邮箱修改事件
func (p *UserProjector) handleUserEmailChanged(ctx context.Context, event *user.UserEmailChangedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET email = ?, updated_at = ? WHERE user_id = ?",
		event.NewEmail,
		event.ChangedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserLocked 处理用户锁定事件
func (p *UserProjector) handleUserLocked(ctx context.Context, event *user.UserLockedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET status = ?, updated_at = ? WHERE user_id = ?",
		int(user.UserStatusLocked),
		event.LockedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserUnlocked 处理用户解锁事件
func (p *UserProjector) handleUserUnlocked(ctx context.Context, event *user.UserUnlockedEvent) error {
	_, err := p.db.Exec(ctx,
		"UPDATE user_read_model SET status = ?, updated_at = ? WHERE user_id = ?",
		int(user.UserStatusActive),
		event.UnlockedAt,
		event.UserID.Int64(),
	)
	return err
}

// handleUserProfileUpdated 处理用户资料更新事件
func (p *UserProjector) handleUserProfileUpdated(ctx context.Context, event *user.UserProfileUpdatedEvent) error {
	// 根据更新的字段动态构建更新语句
	for _, field := range event.UpdatedFields {
		switch field {
		case "display_name":
			// 这里需要从事件中获取新值，简化处理
			_, err := p.db.Exec(ctx,
				"UPDATE user_read_model SET updated_at = ? WHERE user_id = ?",
				event.UpdatedAt,
				event.UserID.Int64(),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ProjectionHandler 投影处理器（用于事件总线订阅）
type ProjectionHandler struct {
	projector *UserProjector
}

// NewProjectionHandler 创建投影处理器
func NewProjectionHandler(projector *UserProjector) *ProjectionHandler {
	return &ProjectionHandler{projector: projector}
}

// Handle 处理事件
func (h *ProjectionHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	return h.projector.Project(ctx, event)
}

// ProjectorEventHandler 适配 EventHandler 接口
type ProjectorEventHandler struct {
	projector *UserProjector
}

// NewProjectorEventHandler 创建投影事件处理器
func NewProjectorEventHandler(projector *UserProjector) *ProjectorEventHandler {
	return &ProjectorEventHandler{projector: projector}
}

// Handle 实现 EventHandler 接口
func (h *ProjectorEventHandler) Handle(ctx context.Context, event ddd.DomainEvent) error {
	return h.projector.Project(ctx, event)
}
