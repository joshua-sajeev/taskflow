package user_service

import (
	"taskflow/internal/dto"

	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

var _ UserServiceInterface = (*UserServiceMock)(nil)

func (m *UserServiceMock) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.CreateUserResponse), nil
}

func (m *UserServiceMock) AuthenticateUser(req *dto.AuthRequest) (*dto.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.AuthResponse), nil
}

func (m *UserServiceMock) UpdatePassword(req *dto.UpdatePasswordRequest) (*dto.UpdatePasswordResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.UpdatePasswordResponse), nil
}

func (m *UserServiceMock) DeleteUser(req *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*dto.DeleteUserResponse), nil
}
