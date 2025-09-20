package user_service

import (
	"taskflow/internal/domain/user"

	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements gorm_user.UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(email string) (*user.User, error) {
	args := m.Called(email)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id int) (*user.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
