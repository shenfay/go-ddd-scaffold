package authentication

import (
	"testing"
	"time"

	"github.com/shenfay/go-ddd-scaffold/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestToUserResponse(t *testing.T) {
	t.Run("should convert user entity to response", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUser := f.CreateUser()

		resp := ToUserResponse(mockUser)

		assert.NotNil(t, resp)
		assert.Equal(t, mockUser.ID, resp.ID)
		assert.Equal(t, mockUser.Email, resp.Email)
		assert.Equal(t, mockUser.EmailVerified, resp.EmailVerified)
		assert.Equal(t, mockUser.LastLoginAt, resp.LastLoginAt)
		assert.Equal(t, mockUser.CreatedAt, resp.CreatedAt)
	})

	t.Run("should return nil for nil user", func(t *testing.T) {
		resp := ToUserResponse(nil)
		assert.Nil(t, resp)
	})

	t.Run("should handle verified user", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUser := f.CreateUser()
		mockUser.VerifyEmail()

		resp := ToUserResponse(mockUser)

		assert.NotNil(t, resp)
		assert.True(t, resp.EmailVerified)
	})

	t.Run("should handle unverified user", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUser := f.CreateUser(factory.WithUnverified())

		resp := ToUserResponse(mockUser)

		assert.NotNil(t, resp)
		assert.False(t, resp.EmailVerified)
	})
}

func TestToAuthResponse(t *testing.T) {
	t.Run("should convert service response to auth response", func(t *testing.T) {
		f := factory.NewUserFactory()
		mockUser := f.CreateUser()
		tokenPair := f.CreateTokenPair(mockUser.ID, mockUser.Email)

		serviceResp := &ServiceAuthResponse{
			User:         mockUser,
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresIn:    3600 * time.Second,
		}

		resp := ToAuthResponse(serviceResp)

		assert.NotNil(t, resp)
		assert.Equal(t, tokenPair.AccessToken, resp.AccessToken)
		assert.Equal(t, tokenPair.RefreshToken, resp.RefreshToken)
		assert.Equal(t, int64(3600), resp.ExpiresIn)
		assert.NotNil(t, resp.User)
		assert.Equal(t, mockUser.Email, resp.User.Email)
	})

	t.Run("should return nil for nil response", func(t *testing.T) {
		resp := ToAuthResponse(nil)
		assert.Nil(t, resp)
	})

	t.Run("should handle nil user in response", func(t *testing.T) {
		serviceResp := &ServiceAuthResponse{
			User:         nil,
			AccessToken:  "access-token",
			RefreshToken: "refresh-token",
			ExpiresIn:    3600 * time.Second,
		}

		resp := ToAuthResponse(serviceResp)

		assert.NotNil(t, resp)
		assert.Equal(t, "access-token", resp.AccessToken)
		assert.Nil(t, resp.User)
	})

	t.Run("should convert duration to seconds", func(t *testing.T) {
		serviceResp := &ServiceAuthResponse{
			User:         nil,
			AccessToken:  "token",
			RefreshToken: "token",
			ExpiresIn:    7200 * time.Second,
		}

		resp := ToAuthResponse(serviceResp)

		assert.Equal(t, int64(7200), resp.ExpiresIn)
	})
}
