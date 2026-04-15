package authentication

import (
	"context"

	domainAuth "github.com/shenfay/go-ddd-scaffold/internal/domain/authentication"
	"github.com/shenfay/go-ddd-scaffold/internal/domain/user"
	authErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/auth"
	userErr "github.com/shenfay/go-ddd-scaffold/pkg/errors/user"
)

// SendVerificationEmail 发送邮箱验证邮件
//
// 业务流程:
// 1. 查找用户
// 2. 检查邮箱是否已验证
// 3. 生成验证令牌
// 4. 保存令牌
// 5. 发布事件(触发邮件发送)
//
// 参数:
//   - ctx: 请求上下文
//   - userID: 用户ID
//
// 返回:
//   - error: 发送失败时返回错误
func (s *Service) SendVerificationEmail(ctx context.Context, userID string) error {
	// 1. 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return userErr.ErrNotFound
	}

	// 2. 检查邮箱是否已验证
	if u.EmailVerified {
		return nil // 已验证,无需重复发送
	}

	// 3. 生成验证令牌
	verifyToken, err := domainAuth.NewEmailVerificationToken(u.ID)
	if err != nil {
		return err
	}

	// 4. 保存令牌
	if err := s.emailTokenRepo.Create(ctx, verifyToken); err != nil {
		return err
	}

	// 5. 发布事件(异步触发邮件发送)
	s.publisher.Publish(ctx, &user.EmailVerificationRequested{
		UserID: u.ID,
		Email:  u.Email,
		Token:  verifyToken.Token,
	})

	return nil
}

// VerifyEmail 验证邮箱
//
// 业务流程:
// 1. 查找验证令牌
// 2. 检查令牌是否过期
// 3. 标记令牌已使用
// 4. 更新用户邮箱验证状态
// 5. 发布事件
//
// 参数:
//   - ctx: 请求上下文
//   - token: 验证令牌
//
// 返回:
//   - error: 验证失败时返回错误
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	// 1. 查找验证令牌
	verifyToken, err := s.emailTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return authErr.ErrInvalidToken
	}

	// 2. 检查令牌是否过期
	if verifyToken.IsExpired() {
		return authErr.ErrTokenExpired
	}

	// 3. 查找用户
	u, err := s.userRepo.FindByID(ctx, verifyToken.UserID)
	if err != nil {
		return userErr.ErrNotFound
	}

	// 4. 标记令牌已使用
	verifyToken.MarkAsUsed()
	if err := s.emailTokenRepo.MarkAsUsed(ctx, verifyToken.ID); err != nil {
		return err
	}

	// 5. 更新用户邮箱验证状态
	u.VerifyEmail()
	if err := s.userRepo.Update(ctx, u); err != nil {
		return err
	}

	// 6. 发布事件
	s.publisher.Publish(ctx, &user.EmailVerified{
		UserID: u.ID,
		Email:  u.Email,
	})

	return nil
}
