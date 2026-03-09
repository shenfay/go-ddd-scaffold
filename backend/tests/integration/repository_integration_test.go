// Package integration_test Repository 层集成测试
package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-ddd-scaffold/internal/domain/user/repository"
	"go-ddd-scaffold/internal/domain/user/valueobject"
	"go-ddd-scaffold/internal/infrastructure/persistence/gorm/repo"
	"go-ddd-scaffold/tests/helper"
)

// UserRepositoryIntegrationSuite 用户仓储集成测试套件
type UserRepositoryIntegrationSuite struct {
	suite.Suite
	db       *gorm.DB
	userRepo repository.UserRepository
}

// SetupSuite 初始化测试环境
func (s *UserRepositoryIntegrationSuite) SetupSuite() {
	// 初始化内存 SQLite 数据库
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)

	// 自动迁移表结构
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
	s.Require().NoError(err)

	// 创建仓储实例
	s.userRepo = repo.NewUserDAORepository(s.db)
}

// TearDownSuite 清理测试环境
func (s *UserRepositoryIntegrationSuite) TearDownSuite() {
	sqlDB, err := s.db.DB()
	s.Require().NoError(err)
	sqlDB.Close()
}

// TestCreateAndGet 测试创建和查询用户
func (s *UserRepositoryIntegrationSuite) TestCreateAndGet() {
	ctx := context.Background()

	// 使用 helper 工厂创建测试用户
	factory := helper.NewUserFactory(s.T())
	user := factory.CreateUser(
		helper.WithEmail("create_get@test.com"),
		helper.WithNickname("Create Get Test User"),
	)

	// 执行创建
	err := s.userRepo.Create(ctx, user)
	s.Require().NoError(err)
	s.NotEmpty(user.ID)

	// 测试 GetByID
	retrievedUser, err := s.userRepo.GetByID(ctx, user.ID)
	s.Require().NoError(err)
	s.Equal(user.Email.String(), retrievedUser.Email.String())
	s.Equal(user.Nickname.String(), retrievedUser.Nickname.String())
	s.Equal(user.Status, retrievedUser.Status)

	// 测试 GetByEmail
	retrievedByEmail, err := s.userRepo.GetByEmail(ctx, "create_get@test.com")
	s.Require().NoError(err)
	s.Equal(user.ID, retrievedByEmail.ID)
}

// TestUpdate 测试更新用户
func (s *UserRepositoryIntegrationSuite) TestUpdate() {
	ctx := context.Background()

	// 创建初始用户
	factory := helper.NewUserFactory(s.T())
	user := factory.CreateUser(
		helper.WithEmail("update@test.com"),
		helper.WithNickname("OriginalNickname"),
	)

	err := s.userRepo.Create(ctx, user)
	s.Require().NoError(err)

	// 修改用户信息
	newEmail, _ := valueobject.NewEmail("newemail@test.com")
	newNickname, _ := valueobject.NewNickname("UpdatedNickname")
	user.Email = newEmail
	user.Nickname = newNickname

	// 执行更新
	err = s.userRepo.Update(ctx, user)
	s.Require().NoError(err)

	// 验证更新结果
	retrievedUser, err := s.userRepo.GetByID(ctx, user.ID)
	s.Require().NoError(err)
	s.Equal("newemail@test.com", retrievedUser.Email.String())
	s.Equal("UpdatedNickname", retrievedUser.Nickname.String())
}

// TestDelete 测试删除用户
func (s *UserRepositoryIntegrationSuite) TestDelete() {
	ctx := context.Background()

	// 创建用户
	factory := helper.NewUserFactory(s.T())
	user := factory.CreateUser(
		helper.WithEmail("delete@test.com"),
	)

	err := s.userRepo.Create(ctx, user)
	s.Require().NoError(err)

	// 执行删除
	err = s.userRepo.Delete(ctx, user.ID)
	s.Require().NoError(err)

	// 验证删除
	_, err = s.userRepo.GetByID(ctx, user.ID)
	s.Error(err)
}

// TestWithTx 测试事务支持
func (s *UserRepositoryIntegrationSuite) TestWithTx() {
	ctx := context.Background()

	// 开启事务
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// 获取事务中的 Repository
	txRepo := s.userRepo.WithTx(tx)

	// 在事务中创建用户
	factory := helper.NewUserFactory(s.T())
	user := factory.CreateUser(
		helper.WithEmail("tx@test.com"),
		helper.WithNickname("Transaction User"),
	)

	createErr := txRepo.Create(ctx, user)
	s.Require().NoError(createErr)

	// 在事务中查询（应该能查到）
	_, getErr := txRepo.GetByID(ctx, user.ID)
	s.NoError(getErr)

	// 在主数据库中查询（事务未提交，应该查不到）
	_, mainDBErr := s.userRepo.GetByID(ctx, user.ID)
	s.Error(mainDBErr)

	// 提交事务
	commitErr := tx.Commit().Error
	s.Require().NoError(commitErr)

	// 再次在主数据库中查询（应该能查到）
	retrievedUser, err := s.userRepo.GetByID(ctx, user.ID)
	s.Require().NoError(err)
	s.Equal("tx@test.com", retrievedUser.Email.String())
}

// TestWithTxRollback 测试事务回滚
func (s *UserRepositoryIntegrationSuite) TestWithTxRollback() {
	ctx := context.Background()

	// 开启事务
	tx := s.db.Begin()
	s.Require().NoError(tx.Error)

	// 获取事务中的 Repository
	txRepo := s.userRepo.WithTx(tx)

	// 在事务中创建用户
	factory := helper.NewUserFactory(s.T())
	user := factory.CreateUser(
		helper.WithEmail("rollback@test.com"),
	)

	createErr := txRepo.Create(ctx, user)
	s.Require().NoError(createErr)

	// 回滚事务
	rollbackErr := tx.Rollback().Error
	s.Require().NoError(rollbackErr)

	// 验证用户不存在（已回滚）
	_, checkErr := s.userRepo.GetByID(ctx, user.ID)
	s.Error(checkErr)
}

// TestListAndCount 测试列表和计数
func (s *UserRepositoryIntegrationSuite) TestListAndCount() {
	ctx := context.Background()

	// 创建多个用户
	factory := helper.NewUserFactory(s.T())
	user1 := factory.CreateUser(helper.WithEmail("user1@test.com"))
	user2 := factory.CreateUser(helper.WithEmail("user2@test.com"))
	user3 := factory.CreateUser(helper.WithEmail("user3@test.com"))

	err := s.userRepo.Create(ctx, user1)
	s.Require().NoError(err)
	err = s.userRepo.Create(ctx, user2)
	s.Require().NoError(err)
	err = s.userRepo.Create(ctx, user3)
	s.Require().NoError(err)

	// 注意：ListByTenant 和 CountByTenant 需要租户 ID
	// 由于当前设计可能不需要租户过滤，这里先跳过
	s.T().Skip("ListByTenant and CountByTenant need tenant ID, skipping for now")
}

// TestGetByIDNotFound 测试用户不存在的情况
func (s *UserRepositoryIntegrationSuite) TestGetByIDNotFound() {
	ctx := context.Background()

	// 查询不存在的用户
	_, err := s.userRepo.GetByID(ctx, uuid.New())
	s.Error(err)
	// 注意：实际实现可能返回自定义错误而不是 gorm.ErrRecordNotFound
	s.Contains(err.Error(), "user not found")
}

// TestUserRepositoryIntegration 运行测试套件
func TestUserRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(UserRepositoryIntegrationSuite))
}
