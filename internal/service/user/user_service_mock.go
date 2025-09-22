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

	var resp *dto.CreateUserResponse
	if r := args.Get(0); r != nil {
		resp = r.(*dto.CreateUserResponse)
	}

	return resp, args.Error(1)
}

func (m *UserServiceMock) AuthenticateUser(req *dto.AuthRequest) (*dto.AuthResponse, error) {
	args := m.Called(req)
	var resp *dto.AuthResponse
	if r := args.Get(0); r != nil {
		resp = r.(*dto.AuthResponse)
	}
	return resp, args.Error(1)
}

func (m *UserServiceMock) UpdatePassword(req *dto.UpdatePasswordRequest) (*dto.UpdatePasswordResponse, error) {
	args := m.Called(req)
	var resp *dto.UpdatePasswordResponse
	if r := args.Get(0); r != nil {
		resp = r.(*dto.UpdatePasswordResponse)
	}
	return resp, args.Error(1)
}

func (m *UserServiceMock) DeleteUser(req *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
	args := m.Called(req)
	var resp *dto.DeleteUserResponse
	if r := args.Get(0); r != nil {
		resp = r.(*dto.DeleteUserResponse)
	}
	return resp, args.Error(1)
}
