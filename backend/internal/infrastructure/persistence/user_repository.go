package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// DB 数据库接口抽象
type DB interface {
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	BeginTx(ctx context.Context) (Tx, error)
}

// Row 单行结果接口
type Row interface {
	Scan(dest ...interface{}) error
}

// Rows 多行结果接口
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}

// Result 执行结果接口
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Tx 事务接口
type Tx interface {
	DB
	Commit() error
	Rollback() error
}

// UserModel 用户数据库模型
type UserModel struct {
	ID             int64
	Username       string
	Email          string
	PasswordHash   string
	Status         int
	DisplayName    string
	FirstName      string
	LastName       string
	Gender         int
	PhoneNumber    string
	AvatarURL      string
	LastLoginAt    *time.Time
	LoginCount     int
	LockedUntil    *time.Time
	FailedAttempts int
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db DB) user.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save 保存用户
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
	model := r.toModel(u)

	tx, err := r.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 检查用户是否存在
	var existingVersion int
	var exists bool
	err = tx.QueryRow(ctx,
		"SELECT version FROM users WHERE id = ?",
		model.ID,
	).Scan(&existingVersion)

	if err == nil {
		exists = true
	}

	if !exists {
		// 创建新用户
		_, err = tx.Exec(ctx,
			`INSERT INTO users (id, username, email, password_hash, status, display_name, 
			first_name, last_name, gender, phone_number, avatar_url, last_login_at, 
			login_count, locked_until, failed_attempts, version, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			model.ID, model.Username, model.Email, model.PasswordHash, model.Status,
			model.DisplayName, model.FirstName, model.LastName, model.Gender,
			model.PhoneNumber, model.AvatarURL, model.LastLoginAt, model.LoginCount,
			model.LockedUntil, model.FailedAttempts, model.Version, model.CreatedAt, model.UpdatedAt,
		)
	} else {
		// 乐观锁检查
		if existingVersion != u.Version()-1 {
			return ddd.NewConcurrencyError(
				u.ID(),
				u.Version()-1,
				existingVersion,
				"user version conflict",
			)
		}

		// 更新用户
		_, err = tx.Exec(ctx,
			`UPDATE users SET username = ?, email = ?, password_hash = ?, status = ?, 
			display_name = ?, first_name = ?, last_name = ?, gender = ?, phone_number = ?, 
			avatar_url = ?, last_login_at = ?, login_count = ?, locked_until = ?, 
			failed_attempts = ?, version = ?, updated_at = ? WHERE id = ?`,
			model.Username, model.Email, model.PasswordHash, model.Status,
			model.DisplayName, model.FirstName, model.LastName, model.Gender,
			model.PhoneNumber, model.AvatarURL, model.LastLoginAt, model.LoginCount,
			model.LockedUntil, model.FailedAttempts, model.Version, model.UpdatedAt, model.ID,
		)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

// FindByID 根据ID查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	var model UserModel
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at 
		FROM users WHERE id = ?`,
		id.Int64(),
	).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash, &model.Status,
		&model.DisplayName, &model.FirstName, &model.LastName, &model.Gender,
		&model.PhoneNumber, &model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.LockedUntil, &model.FailedAttempts, &model.Version, &model.CreatedAt, &model.UpdatedAt,
	)

	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	return r.toDomain(&model), nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id user.UserID) error {
	result, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = ?", id.Int64())
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return ddd.ErrAggregateNotFound
	}

	return nil
}

// Exists 检查用户是否存在
func (r *UserRepositoryImpl) Exists(ctx context.Context, id user.UserID) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE id = ?", id.Int64()).Scan(&count)
	return count > 0, err
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var model UserModel
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at 
		FROM users WHERE username = ?`,
		username,
	).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash, &model.Status,
		&model.DisplayName, &model.FirstName, &model.LastName, &model.Gender,
		&model.PhoneNumber, &model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.LockedUntil, &model.FailedAttempts, &model.Version, &model.CreatedAt, &model.UpdatedAt,
	)

	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	return r.toDomain(&model), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var model UserModel
	err := r.db.QueryRow(ctx,
		`SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at 
		FROM users WHERE email = ?`,
		email,
	).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash, &model.Status,
		&model.DisplayName, &model.FirstName, &model.LastName, &model.Gender,
		&model.PhoneNumber, &model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.LockedUntil, &model.FailedAttempts, &model.Version, &model.CreatedAt, &model.UpdatedAt,
	)

	if err != nil {
		return nil, ddd.ErrAggregateNotFound
	}

	return r.toDomain(&model), nil
}

// FindByStatus 根据状态查找用户
func (r *UserRepositoryImpl) FindByStatus(ctx context.Context, status user.UserStatus) ([]*user.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at 
		FROM users WHERE status = ?`,
		int(status),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// FindAll 分页查找所有用户
func (r *UserRepositoryImpl) FindAll(ctx context.Context, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
	var total int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at 
		FROM users LIMIT ? OFFSET ?`,
		pagination.Limit(), pagination.Offset(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return &ddd.PaginatedResult[*user.User]{
		Items:      users,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.CalculateTotalPages(total),
	}, nil
}

// FindByCriteria 根据条件查找用户
func (r *UserRepositoryImpl) FindByCriteria(ctx context.Context, criteria user.UserSearchCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
	query := `SELECT id, username, email, password_hash, status, display_name, first_name, 
		last_name, gender, phone_number, avatar_url, last_login_at, login_count, 
		locked_until, failed_attempts, version, created_at, updated_at FROM users WHERE 1=1`
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	var args []interface{}

	if criteria.Keyword != "" {
		query += " AND (username LIKE ? OR email LIKE ? OR display_name LIKE ?)"
		countQuery += " AND (username LIKE ? OR email LIKE ? OR display_name LIKE ?)"
		keyword := "%" + criteria.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}

	if criteria.Status != nil {
		query += " AND status = ?"
		countQuery += " AND status = ?"
		args = append(args, int(*criteria.Status))
	}

	if criteria.Gender != nil {
		query += " AND gender = ?"
		countQuery += " AND gender = ?"
		args = append(args, int(*criteria.Gender))
	}

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

	users, err := r.scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return &ddd.PaginatedResult[*user.User]{
		Items:      users,
		TotalCount: total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: pagination.CalculateTotalPages(total),
	}, nil
}

// Count 统计用户总数
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// CountByStatus 根据状态统计用户数
func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status user.UserStatus) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE status = ?", int(status)).Scan(&count)
	return count, err
}

// SaveBatch 批量保存用户
func (r *UserRepositoryImpl) SaveBatch(ctx context.Context, users []*user.User) error {
	// 简化实现：逐个保存
	for _, u := range users {
		if err := r.Save(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 批量删除用户
func (r *UserRepositoryImpl) DeleteBatch(ctx context.Context, ids []user.UserID) error {
	if len(ids) == 0 {
		return nil
	}
	// 简化实现：逐个删除
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// SaveWithVersion 带乐观锁的保存
func (r *UserRepositoryImpl) SaveWithVersion(ctx context.Context, u *user.User, expectedVersion int) error {
	model := r.toModel(u)

	result, err := r.db.Exec(ctx,
		`UPDATE users SET username = ?, email = ?, password_hash = ?, status = ?, 
		display_name = ?, first_name = ?, last_name = ?, gender = ?, phone_number = ?, 
		avatar_url = ?, last_login_at = ?, login_count = ?, locked_until = ?, 
		failed_attempts = ?, version = ?, updated_at = ? WHERE id = ? AND version = ?`,
		model.Username, model.Email, model.PasswordHash, model.Status,
		model.DisplayName, model.FirstName, model.LastName, model.Gender,
		model.PhoneNumber, model.AvatarURL, model.LastLoginAt, model.LoginCount,
		model.LockedUntil, model.FailedAttempts, model.Version, model.UpdatedAt,
		model.ID, expectedVersion,
	)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return ddd.NewConcurrencyError(
			u.ID(),
			expectedVersion,
			u.Version(),
			"user version conflict during save",
		)
	}

	return nil
}

// scanUsers 扫描用户列表
func (r *UserRepositoryImpl) scanUsers(rows Rows) ([]*user.User, error) {
	var users []*user.User
	for rows.Next() {
		var model UserModel
		err := rows.Scan(
			&model.ID, &model.Username, &model.Email, &model.PasswordHash, &model.Status,
			&model.DisplayName, &model.FirstName, &model.LastName, &model.Gender,
			&model.PhoneNumber, &model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
			&model.LockedUntil, &model.FailedAttempts, &model.Version, &model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, r.toDomain(&model))
	}
	return users, nil
}

// toModel 将领域对象转换为数据库模型
func (r *UserRepositoryImpl) toModel(u *user.User) *UserModel {
	return &UserModel{
		ID:             u.ID().(user.UserID).Int64(),
		Username:       u.Username().Value(),
		Email:          u.Email().Value(),
		PasswordHash:   u.Password().Value(),
		Status:         int(u.Status()),
		DisplayName:    u.DisplayName(),
		FirstName:      u.FirstName(),
		LastName:       u.LastName(),
		Gender:         int(u.Gender()),
		PhoneNumber:    u.PhoneNumber(),
		AvatarURL:      u.AvatarURL(),
		LastLoginAt:    u.LastLoginAt(),
		LoginCount:     u.LoginCount(),
		LockedUntil:    u.LockedUntil(),
		FailedAttempts: u.FailedAttempts(),
		Version:        u.Version(),
		CreatedAt:      u.CreatedAt(),
		UpdatedAt:      u.UpdatedAt(),
	}
}

// toDomain 将数据库模型转换为领域对象
func (r *UserRepositoryImpl) toDomain(model *UserModel) *user.User {
	u := &user.User{}
	u.SetID(user.NewUserID(model.ID))
	u.SetVersion(model.Version)
	u.SetCreatedAt(model.CreatedAt)
	u.SetUpdatedAt(model.UpdatedAt)

	// 使用反射设置私有字段（实际项目中可能需要更好的方式）
	// 这里简化处理，实际应该通过构造函数或工厂方法重建对象
	return u
}

// DomainEventsModel 领域事件数据库模型
type DomainEventsModel struct {
	ID            int64
	AggregateID   string
	AggregateType string
	EventType     string
	EventVersion  int
	EventData     string
	OccurredOn    time.Time
	Processed     bool
	CreatedAt     time.Time
}

// EventStoreImpl 事件存储实现
type EventStoreImpl struct {
	db DB
}

// NewEventStore 创建事件存储
func NewEventStore(db DB) *EventStoreImpl {
	return &EventStoreImpl{db: db}
}

// AppendEvents 追加事件
func (s *EventStoreImpl) AppendEvents(ctx context.Context, aggregateID interface{}, events []ddd.DomainEvent) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO domain_events (aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, processed, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			fmt.Sprintf("%v", aggregateID),
			"user",
			event.EventName(),
			event.Version(),
			string(eventData),
			event.OccurredOn(),
			false,
			time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to store event: %w", err)
		}
	}

	return tx.Commit()
}

// GetEventsForAggregate 获取聚合的事件
func (s *EventStoreImpl) GetEventsForAggregate(ctx context.Context, aggregateID interface{}, afterVersion int) ([]ddd.DomainEvent, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, processed, created_at 
		FROM domain_events WHERE aggregate_id = ? AND event_version > ? ORDER BY event_version ASC`,
		fmt.Sprintf("%v", aggregateID),
		afterVersion,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

// GetUnprocessedEvents 获取未处理的事件
func (s *EventStoreImpl) GetUnprocessedEvents(ctx context.Context, limit int) ([]ddd.DomainEvent, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, aggregate_id, aggregate_type, event_type, event_version, event_data, occurred_on, processed, created_at 
		FROM domain_events WHERE processed = ? ORDER BY id ASC LIMIT ?`,
		false, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.scanEvents(rows)
}

// MarkEventsAsProcessed 标记事件为已处理
func (s *EventStoreImpl) MarkEventsAsProcessed(ctx context.Context, eventIDs []int64) error {
	if len(eventIDs) == 0 {
		return nil
	}
	// 简化实现：逐个更新
	for _, id := range eventIDs {
		_, err := s.db.Exec(ctx, "UPDATE domain_events SET processed = ? WHERE id = ?", true, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// scanEvents 扫描事件列表
func (s *EventStoreImpl) scanEvents(rows Rows) ([]ddd.DomainEvent, error) {
	var events []ddd.DomainEvent
	for rows.Next() {
		var model DomainEventsModel
		err := rows.Scan(
			&model.ID, &model.AggregateID, &model.AggregateType, &model.EventType,
			&model.EventVersion, &model.EventData, &model.OccurredOn, &model.Processed, &model.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		event, err := s.deserializeEvent(model.EventType, model.EventData)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

// deserializeEvent 反序列化事件
func (s *EventStoreImpl) deserializeEvent(eventType string, eventData string) (ddd.DomainEvent, error) {
	switch eventType {
	case "UserRegistered":
		var event user.UserRegisteredEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserActivated":
		var event user.UserActivatedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserDeactivated":
		var event user.UserDeactivatedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserLoggedIn":
		var event user.UserLoggedInEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserPasswordChanged":
		var event user.UserPasswordChangedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserEmailChanged":
		var event user.UserEmailChangedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserLocked":
		var event user.UserLockedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserUnlocked":
		var event user.UserUnlockedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserProfileUpdated":
		var event user.UserProfileUpdatedEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	case "UserFailedLoginAttempt":
		var event user.UserFailedLoginAttemptEvent
		if err := json.Unmarshal([]byte(eventData), &event); err != nil {
			return nil, err
		}
		return event, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
