package mocks

import (
	"avito-tech-go/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(ID uint) (*domain.User, error) {
	args := m.Called(ID)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) GetUserByName(username string) (*domain.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) DeleteUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ChangeCoins(userID int, delta int) error {
	args := m.Called(userID, delta)
	return args.Error(0)
}
