package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	"github.com/shenfay/go-ddd-scaffold/shared/ddd"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *sql.DB) user.UserRepository {
	return &UserRepositoryImpl{db: db}
}

// Save 保存用户（支持创建和更新，带乐观锁）
func (r *UserRepositoryImpl) Save(ctx context.Context, u *user.User) error {
	if u.Version() == 0 {
		return r.insert(ctx, u)
	}
	return r.update(ctx, u)
}

// insert 插入新用户
func (r *UserRepositoryImpl) insert(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, status, display_name,
			gender, phone_number, avatar_url, last_login_at,
			login_count, failed_attempts, locked_until,
			version, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.ExecContext(ctx, query,
		u.ID().(user.UserID).Int64(),
		u.Username().Value(),
		u.Email().Value(),
		u.Password().Value(),
		int(u.Status()),
		u.DisplayName(),
		int(u.Gender()),
		u.PhoneNumber(),
		u.AvatarURL(),
		u.LastLoginAt(),
		u.LoginCount(),
		u.FailedAttempts(),
		u.LockedUntil(),
		u.Version(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)

	if err != nil {
		return err
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// update 更新现有用户（带乐观锁检查）
func (r *UserRepositoryImpl) update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users SET
			username = $1, email = $2, password_hash = $3, status = $4,
			display_name = $5, gender = $6, phone_number = $7, avatar_url = $8,
			last_login_at = $9, login_count = $10, failed_attempts = $11,
			locked_until = $12, version = $13, updated_at = $14
		WHERE id = $15 AND version = $16
	`

	result, err := r.db.ExecContext(ctx, query,
		u.Username().Value(),
		u.Email().Value(),
		u.Password().Value(),
		int(u.Status()),
		u.DisplayName(),
		int(u.Gender()),
		u.PhoneNumber(),
		u.AvatarURL(),
		u.LastLoginAt(),
		u.LoginCount(),
		u.FailedAttempts(),
		u.LockedUntil(),
		u.Version(),
		u.UpdatedAt(),
		u.ID().(user.UserID).Int64(),
		u.Version()-1, // 期望的旧版本号
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ddd.NewConcurrencyError(
			u.ID(),
			u.Version()-1,
			u.Version(),
			"user was updated by another transaction",
		)
	}

	// 保存领域事件
	return r.saveEvents(ctx, u)
}

// saveEvents 保存领域事件到事件存储
func (r *UserRepositoryImpl) saveEvents(ctx context.Context, u *user.User) error {
	events := u.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	query := `
		INSERT INTO domain_events (
			aggregate_id, aggregate_type, event_type, event_version,
			event_data, occurred_on, processed, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, false, NOW())
	`

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return err
		}

		_, err = r.db.ExecContext(ctx, query,
			u.ID().(user.UserID).Int64(),
			"user",
			event.EventName(),
			event.Version(),
			string(eventData),
			event.OccurredOn(),
		)
		if err != nil {
			return err
		}
	}

	// 清除已保存的事件
	u.ClearUncommittedEvents()
	return nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, status, display_name,
		       gender, phone_number, avatar_url, last_login_at,
		       login_count, failed_attempts, locked_until,
		       version, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var model userModel
	err := r.db.QueryRowContext(ctx, query, id.Int64()).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash,
		&model.Status, &model.DisplayName, &model.Gender, &model.PhoneNumber,
		&model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.FailedAttempts, &model.LockedUntil, &model.Version,
		&model.CreatedAt, &model.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ddd.ErrAggregateNotFound
	}
	if err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, status, display_name,
		       gender, phone_number, avatar_url, last_login_at,
		       login_count, failed_attempts, locked_until,
		       version, created_at, updated_at
		FROM users
		WHERE username = $1 AND deleted_at IS NULL
	`

	var model userModel
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash,
		&model.Status, &model.DisplayName, &model.Gender, &model.PhoneNumber,
		&model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.FailedAttempts, &model.LockedUntil, &model.Version,
		&model.CreatedAt, &model.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ddd.ErrAggregateNotFound
	}
	if err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

// FindByEmail 根据邮箱查找用户
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, status, display_name,
		       gender, phone_number, avatar_url, last_login_at,
		       login_count, failed_attempts, locked_until,
		       version, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	var model userModel
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&model.ID, &model.Username, &model.Email, &model.PasswordHash,
		&model.Status, &model.DisplayName, &model.Gender, &model.PhoneNumber,
		&model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
		&model.FailedAttempts, &model.LockedUntil, &model.Version,
		&model.CreatedAt, &model.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ddd.ErrAggregateNotFound
	}
	if err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

// Delete 软删除用户
func (r *UserRepositoryImpl) Delete(ctx context.Context, id user.UserID) error {
	query := `UPDATE users SET deleted_at = NOW(), version = version + 1 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id.Int64())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ddd.ErrAggregateNotFound
	}

	return nil
}

// Exists 检查用户是否存在
func (r *UserRepositoryImpl) Exists(ctx context.Context, id user.UserID) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE id = $1 AND deleted_at IS NULL`

	var count int
	err := r.db.QueryRowContext(ctx, query, id.Int64()).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Count 统计用户总数
func (r *UserRepositoryImpl) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// CountByStatus 按状态统计用户数
func (r *UserRepositoryImpl) CountByStatus(ctx context.Context, status user.UserStatus) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE status = $1 AND deleted_at IS NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query, int(status)).Scan(&count)
	return count, err
}

// FindByStatus 根据状态查找用户
func (r *UserRepositoryImpl) FindByStatus(ctx context.Context, status user.UserStatus) ([]*user.User, error) {
	query := `
		SELECT id, username, email, password_hash, status, display_name,
		       gender, phone_number, avatar_url, last_login_at,
		       login_count, failed_attempts, locked_until,
		       version, created_at, updated_at
		FROM users
		WHERE status = $1 AND deleted_at IS NULL
	`

	rows, err := r.db.QueryContext(ctx, query, int(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanUsers(rows)
}

// FindAll 分页查询所有用户
func (r *UserRepositoryImpl) FindAll(ctx context.Context, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, username, email, password_hash, status, display_name,
		       gender, phone_number, avatar_url, last_login_at,
		       login_count, failed_attempts, locked_until,
		       version, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pagination.PageSize, pagination.Offset())
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

// FindByCriteria 根据条件查询用户
func (r *UserRepositoryImpl) FindByCriteria(ctx context.Context, criteria user.UserSearchCriteria, pagination ddd.Pagination) (*ddd.PaginatedResult[*user.User], error) {
	// 简化实现，后续可扩展
	return r.FindAll(ctx, pagination)
}

// SaveBatch 批量保存用户
func (r *UserRepositoryImpl) SaveBatch(ctx context.Context, users []*user.User) error {
	for _, u := range users {
		if err := r.Save(ctx, u); err != nil {
			return err
		}
	}
	return nil
}

// DeleteBatch 批量删除用户
func (r *UserRepositoryImpl) DeleteBatch(ctx context.Context, ids []user.UserID) error {
	for _, id := range ids {
		if err := r.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

// SaveWithVersion 带乐观锁的保存（已实现）
func (r *UserRepositoryImpl) SaveWithVersion(ctx context.Context, u *user.User, expectedVersion int) error {
	// 当前 Save 方法已经实现了乐观锁逻辑
	// 这里为了接口兼容性直接调用 Save
	return r.Save(ctx, u)
}

// scanUsers 扫描用户列表
func (r *UserRepositoryImpl) scanUsers(rows *sql.Rows) ([]*user.User, error) {
	var users []*user.User
	for rows.Next() {
		var model userModel
		err := rows.Scan(
			&model.ID, &model.Username, &model.Email, &model.PasswordHash,
			&model.Status, &model.DisplayName, &model.Gender, &model.PhoneNumber,
			&model.AvatarURL, &model.LastLoginAt, &model.LoginCount,
			&model.FailedAttempts, &model.LockedUntil, &model.Version,
			&model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, r.toDomain(&model))
	}
	return users, nil
}

// toDomain 将数据库模型转换为领域对象
func (r *UserRepositoryImpl) toDomain(model *userModel) *user.User {
	// 使用 Builder 模式优雅地重建用户对象
	return user.NewUserBuilder().
		WithID(model.ID).
		WithUsername(model.Username).
		WithEmail(model.Email).
		WithPasswordHash(model.PasswordHash).
		WithStatus(user.UserStatus(model.Status)).
		WithGender(user.UserGender(model.Gender)).
		WithDisplayName(model.DisplayName).
		WithPhoneNumber(model.PhoneNumber).
		WithAvatarURL(model.AvatarURL).
		WithLastLoginAt(model.LastLoginAt).
		WithLoginCount(model.LoginCount).
		WithFailedAttempts(model.FailedAttempts).
		WithLockedUntil(model.LockedUntil).
		WithVersion(model.Version).
		WithTimestamps(model.CreatedAt, model.UpdatedAt).
		Build()
}

// userModel 用户数据库模型
type userModel struct {
	ID             int64
	Username       string
	Email          string
	PasswordHash   string
	Status         int
	DisplayName    string
	Gender         int
	PhoneNumber    string
	AvatarURL      string
	LastLoginAt    *time.Time
	LoginCount     int
	FailedAttempts int
	LockedUntil    *time.Time
	Version        int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
