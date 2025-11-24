package user_service

import (
	"errors"
	"fmt"
	"strings"

	"taskflow/internal/domain/user"
	"taskflow/internal/dto"
	"taskflow/internal/repository/gorm/gorm_user"
	"taskflow/pkg"
	"taskflow/pkg/jwt"
	"taskflow/pkg/validator"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo gorm_user.UserRepositoryInterface
}

func NewUserService(repo gorm_user.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

var _ UserServiceInterface = (*UserService)(nil)

func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	validator := validator.NewPasswordValidator()
	if err := validator.Validate(req.Password); err != nil {
		return nil, errors.New("password validation failed, choose a stronger password")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u := &user.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(u); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, errors.New("email already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.CreateUserResponse{
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

func (s *UserService) AuthenticateUser(req *dto.AuthRequest) (*dto.AuthResponse, error) {
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	u, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	secretKey := []byte(pkg.GetEnv("JWT_SECRET", "secret-key"))
	token, err := jwt.CreateToken(u.ID, u.Email, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		ID:    u.ID,
		Email: u.Email,
	}, nil
}

func (s *UserService) UpdatePassword(req *dto.UpdatePasswordRequest) (*dto.UpdatePasswordResponse, error) {
	if req.OldPassword == "" {
		return nil, errors.New("old password is required")
	}
	if req.NewPassword == "" {
		return nil, errors.New("new password is required")
	}

	validator := validator.NewPasswordValidator()
	if err := validator.Validate(req.NewPassword); err != nil {
		return nil, errors.New("password validation failed, choose a stronger password")
	}

	u, err := s.repo.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)); err != nil {
		return nil, errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = string(hashedPassword)
	if err := s.repo.Update(u); err != nil {
		return nil, fmt.Errorf("failed to update password: %w", err)
	}

	return &dto.UpdatePasswordResponse{
		Message: "Password updated successfully",
	}, nil
}

func (s *UserService) DeleteUser(req *dto.DeleteUserRequest) (*dto.DeleteUserResponse, error) {
	if req.ID == 0 {
		return nil, errors.New("user ID is required")
	}

	if err := s.repo.Delete(req.ID); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	return &dto.DeleteUserResponse{
		Message: "User account deleted successfully",
	}, nil
}

func HashPassword(p string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b), err
}
