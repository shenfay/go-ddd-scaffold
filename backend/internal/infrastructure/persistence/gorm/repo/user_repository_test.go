package repo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/domain/user/entity"
	"go-ddd-scaffold/internal/domain/user/service"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/model"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
)

// hashPassword 密码哈希辅助函数
func hashPassword(plain string) (string, error) {
	hasher := service.NewDefaultBcryptPasswordHasher()
	return hasher.Hash(plain)
}

// UserRepositoryTestSuite 用户仓储测试套件
type UserRepositoryTestSuite struct {
	suite.Suite
	db     *gorm.DB
	repo   repo.UserDAORepository
	userID uuid.UUID
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// 初始化内存 SQLite 数据库
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("failed to connect database: %v", err)
	}

	// 自动迁移表结构（简化版，避免 PostgreSQL 特定语法）
	err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			nickname TEXT NOT NULL,
			avatar TEXT,
			phone TEXT,
			bio TEXT,
			status TEXT NOT NULL DEFAULT 'active',
			role TEXT NOT NULL DEFAULT 'member',
			tenant_id TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
	if err != nil {
		s.T().Fatalf("failed to migrate database: %v", err)
	}

	// 创建仓储实例
	s.repo = *repo.NewUserDAORepository(s.db).(*repo.UserDAORepository)
	
	// 生成测试用的 UUID
	s.userID = uuid.New()
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	// 关闭数据库连接
	sqlDB, err := s.db.DB()
	if err != nil {
		s.T().Logf("failed to get sql.DB: %v", err)
	} else {
		sqlDB.Close()
	}
}

func (s *UserRepositoryTestSuite) TestCreate() {
	ctx := context.Background()

	// 准备测试数据
	email, err := valueobject.NewEmail("test_create@example.com")
	s.Require().NoError(err)

	nickname, err := valueobject.NewNickname("测试用户")
	s.Require().NoError(err)

	hashedPassword, err := hashPassword("Password123!")
	s.Require().NoError(err)

	user := &entity.User{
		ID:       s.userID,
		Email:    email,
		Password: entity.HashedPassword(hashedPassword),
		Nickname: nickname,
		Status:   entity.StatusActive,
	}

	// 执行创建
	err = s.repo.Create(ctx, user)
	s.NoError(err)

	// 验证数据库中存在
	var count int64
	err = s.db.Model(&model.User{}).Where("id = ?", s.userID.String()).Count(&count).Error
	s.NoError(err)
	s.Equal(int64(1), count)
}

func (s *UserRepositoryTestSuite) TestGetByID() {
	ctx := context.Background()

	// 先创建测试数据
	testUserID := uuid.New()
	email, _ := valueobject.NewEmail("test_getbyid@example.com")
	nickname, _ := valueobject.NewNickname("测试用户 2")
	hashedPassword, _ := hashPassword("Password123!")

	user := &entity.User{
		ID:       testUserID,
		Email:    email,
		Password: entity.HashedPassword(hashedPassword),
		Nickname: nickname,
		Status:   entity.StatusActive,
	}
	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)

	// 测试获取成功场景
	foundUser, err := s.repo.GetByID(ctx, testUserID)
	s.NoError(err)
	s.NotNil(foundUser)
	s.Equal(testUserID, foundUser.ID)
	s.Equal("test_getbyid@example.com", foundUser.Email.String())
	s.Equal("测试用户 2", foundUser.Nickname.String())

	// 测试用户不存在场景
	notFoundUser, err := s.repo.GetByID(ctx, uuid.New())
	s.Error(err)
	s.Nil(notFoundUser)
	s.Contains(err.Error(), "user not found")
}

func (s *UserRepositoryTestSuite) TestGetByEmail() {
	ctx := context.Background()

	// 先创建测试数据
	testEmail := "test_getbyemail@example.com"
	email, _ := valueobject.NewEmail(testEmail)
	nickname, _ := valueobject.NewNickname("测试用户 3")
	hashedPassword, _ := hashPassword("Password123!")

	user := &entity.User{
		ID:       uuid.New(),
		Email:    email,
		Password: entity.HashedPassword(hashedPassword),
		Nickname: nickname,
		Status:   entity.StatusActive,
	}
	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)

	// 测试获取成功场景
	foundUser, err := s.repo.GetByEmail(ctx, testEmail)
	s.NoError(err)
	s.NotNil(foundUser)
	s.Equal(testEmail, foundUser.Email.String())

	// 测试邮箱不存在场景
	notFoundUser, err := s.repo.GetByEmail(ctx, "nonexistent@example.com")
	s.Error(err)
	s.Nil(notFoundUser)
	s.Contains(err.Error(), "user not found")
}

func (s *UserRepositoryTestSuite) TestUpdate() {
	ctx := context.Background()

	// 先创建测试数据
	testUserID := uuid.New()
	email, _ := valueobject.NewEmail("test_update@example.com")
	nickname, _ := valueobject.NewNickname("原昵称")
	hashedPassword, _ := hashPassword("Password123!")

	user := &entity.User{
		ID:       testUserID,
		Email:    email,
		Password: entity.HashedPassword(hashedPassword),
		Nickname: nickname,
		Status:   entity.StatusActive,
	}
	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)

	// 更新用户资料
	newNickname, _ := valueobject.NewNickname("新昵称")
	user.Nickname = newNickname
	
	err = s.repo.Update(ctx, user)
	s.NoError(err)

	// 验证更新成功
	updatedUser, err := s.repo.GetByID(ctx, testUserID)
	s.NoError(err)
	s.NotNil(updatedUser)
	s.Equal("新昵称", updatedUser.Nickname.String())
}

func (s *UserRepositoryTestSuite) TestDelete() {
	ctx := context.Background()

	// 先创建测试数据
	testUserID := uuid.New()
	email, _ := valueobject.NewEmail("test_delete@example.com")
	nickname, _ := valueobject.NewNickname("删除测试")
	hashedPassword, _ := hashPassword("Password123!")

	user := &entity.User{
		ID:       testUserID,
		Email:    email,
		Password: entity.HashedPassword(hashedPassword),
		Nickname: nickname,
		Status:   entity.StatusActive,
	}
	err := s.repo.Create(ctx, user)
	s.Require().NoError(err)

	// 执行删除
	err = s.repo.Delete(ctx, testUserID)
	s.NoError(err)

	// 验证已删除
	deletedUser, err := s.repo.GetByID(ctx, testUserID)
	s.Error(err)
	s.Nil(deletedUser)
}

func (s *UserRepositoryTestSuite) TestListByTenant() {
	// 这个测试需要租户成员关系，暂时跳过
	// TODO: 实现完整的租户相关测试
	s.T().Skip("需要完整的租户成员关系支持")
}

func (s *UserRepositoryTestSuite) TestCountByTenant() {
	// 这个测试需要租户成员关系，暂时跳过
	// TODO: 实现完整的租户相关测试
	s.T().Skip("需要完整的租户成员关系支持")
}

// TestUserRepositorySuite 运行测试套件
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
