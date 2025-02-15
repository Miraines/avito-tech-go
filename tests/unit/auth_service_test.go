package unit

import (
	"avito-tech-go/internal/domain"
	"avito-tech-go/internal/services"
	"avito-tech-go/tests/unit/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	authSvc := services.NewAuthService(mockUserRepo, "test_secret")

	t.Run("user already exists", func(t *testing.T) {
		mockUserRepo.On("ExistsByUsername", "alex").
			Return(true, nil).Once()

		token, err := authSvc.Register("alex", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("db error on exists check", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.On("ExistsByUsername", "alex").
			Return(false, errors.New("db error")).Once()

		token, err := authSvc.Register("alex", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success user creation", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.On("ExistsByUsername", "alex").
			Return(false, nil).Once()
		mockUserRepo.On("CreateUser", mock.Anything).
			Return(nil).Once()

		token, err := authSvc.Register("alex", "12345")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	authSvc := services.NewAuthService(mockUserRepo, "test_secret")

	t.Run("user not found", func(t *testing.T) {
		mockUserRepo.On("GetUserByName", "alex").
			Return((*domain.User)(nil), nil).Once()

		token, err := authSvc.Login("alex", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("db error on get user", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		mockUserRepo.On("GetUserByName", "alex").
			Return((*domain.User)(nil), errors.New("db error")).Once()

		token, err := authSvc.Login("alex", "12345")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		user := &domain.User{ID: 1, Username: "alex", PasswordHash: "wrong-hash"}
		mockUserRepo.On("GetUserByName", "alex").
			Return(user, nil).Once()

		token, err := authSvc.Login("alex", "bad-pass")
		assert.Error(t, err)
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("success login", func(t *testing.T) {
		mockUserRepo.ExpectedCalls = nil
		hashed, _ := services.HashPassword("12345")
		user := &domain.User{ID: 1, Username: "alex", PasswordHash: hashed}

		mockUserRepo.On("GetUserByName", "alex").
			Return(user, nil).Once()

		token, err := authSvc.Login("alex", "12345")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		mockUserRepo.AssertExpectations(t)
	})
}
