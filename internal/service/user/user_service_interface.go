package user_service

import "taskflow/internal/dto"

type UserServiceInterface interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error)
	AuthenticateUser(req *dto.AuthRequest) (*dto.AuthResponse, error)
	UpdatePassword(req *dto.UpdatePasswordRequest) (*dto.UpdatePasswordResponse, error)
	DeleteUser(req *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error)
}
