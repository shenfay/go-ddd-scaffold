package auth

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	apperrors "github.com/shenfay/go-ddd-scaffold/pkg/errors"
)

// UserPO 用户持久化对象
type UserPO struct {
	ID             string    `gorm:"primaryKey;type:varchar(50)" json:"id"`
	Email          string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Password       string    `gorm:"type:varchar(255);not null" json:"-"`
	EmailVerified  bool      `gorm:"default:false" json:"email_verified"`
	Locked         bool      `gorm:"default:false" json:"locked"`
	FailedAttempts int       `gorm:"default:0" json:"failed_attempts"`
	LastLoginAt    *TimeNull `json:"last_login_at"`
	CreatedAt      TimeNull  `json:"created_at"`
	UpdatedAt      TimeNull  `json:"updated_at"`
}

// TimeNull 可空的时间类型
type TimeNull struct {
	Time  time.Time
	Valid bool
}

// Value 实现 driver.Valuer 接口，用于 GORM 数据库操作
func (t TimeNull) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan 实现 sql.Scanner 接口，用于 GORM 数据库操作
func (t *TimeNull) Scan(value interface{}) error {
	if value == nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time = v
		t.Valid = true
		return nil
	default:
		return fmt.Errorf("failed to scan TimeNull: %v", value)
	}
}

// MarshalJSON 实现 JSON 序列化
func (t TimeNull) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.Time.Format(time.RFC3339))
	}
	return json.Marshal(nil)
}

// UnmarshalJSON 实现 JSON 反序列化
func (t *TimeNull) UnmarshalJSON(data []byte) error {
	var s interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == nil {
		t.Valid = false
		return nil
	}

	str, ok := s.(string)
	if !ok {
		return fmt.Errorf("invalid time value")
	}

	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	t.Time = parsed
	t.Valid = true
	return nil
}

// ToDomain 转换为领域模型
func (po *UserPO) ToDomain() *User {
	if po == nil {
		return nil
	}

	// 安全处理时间字段，避免空指针
	createdAt := time.Time{}
	updatedAt := time.Time{}

	if po.CreatedAt.Valid {
		createdAt = po.CreatedAt.Time
	}
	if po.UpdatedAt.Valid {
		updatedAt = po.UpdatedAt.Time
	}

	user := &User{
		ID:             po.ID,
		Email:          po.Email,
		Password:       po.Password,
		EmailVerified:  po.EmailVerified,
		Locked:         po.Locked,
		FailedAttempts: po.FailedAttempts,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}

	// 安全处理可空的时间字段
	if lastLogin := po.LastLoginAt; lastLogin != nil && lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	return user
}

// ToPO 从领域模型转换
func ToPO(user *User) *UserPO {
	po := &UserPO{
		ID:             user.ID,
		Email:          user.Email,
		Password:       user.Password,
		EmailVerified:  user.EmailVerified,
		Locked:         user.Locked,
		FailedAttempts: user.FailedAttempts,
		CreatedAt:      TimeNull{Time: user.CreatedAt, Valid: true},
		UpdatedAt:      TimeNull{Time: user.UpdatedAt, Valid: true},
	}

	if user.LastLoginAt != nil {
		po.LastLoginAt = &TimeNull{Time: *user.LastLoginAt, Valid: true}
	}

	return po
}

// userRepository GORM 实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	po := ToPO(user)
	return r.db.WithContext(ctx).Create(po).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*User, error) {
	var po UserPO
	err := r.db.WithContext(ctx).First(&po, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var po UserPO
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}
	return po.ToDomain(), nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) bool {
	var count int64
	r.db.Model(&UserPO{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	po := ToPO(user)
	return r.db.WithContext(ctx).Save(po).Error
}
