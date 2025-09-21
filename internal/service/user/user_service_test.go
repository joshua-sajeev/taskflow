package user_service

import (
	"errors"
	"testing"

	"taskflow/internal/domain/user"
	"taskflow/internal/dto"
	"taskflow/internal/repository/gorm/gorm_user"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		req       *dto.CreateUserRequest
		mockSetup func(m *gorm_user.MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req:  &dto.CreateUserRequest{Email: "test@example.com", Password: "secret123"},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				m.On("Create", mock.Anything).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req:  &dto.CreateUserRequest{Email: "dup@example.com", Password: "secret123"},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				m.On("Create", mock.Anything).Return(errors.New("duplicate key")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(gorm_user.MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.CreateUser(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.req.Email, resp.Email)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthenticateUser(t *testing.T) {
	hashedPass, _ := HashPassword("mypassword")

	tests := []struct {
		name      string
		req       *dto.AuthRequest
		mockSetup func(m *gorm_user.MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req:  &dto.AuthRequest{Email: "user@example.com", Password: "mypassword"},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				u := &user.User{ID: 1, Email: "user@example.com", Password: hashedPass}
				m.On("GetByEmail", "user@example.com").Return(u, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "invalid password",
			req:  &dto.AuthRequest{Email: "user2@example.com", Password: "wrongpass"},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				u := &user.User{ID: 2, Email: "user2@example.com", Password: hashedPass}
				m.On("GetByEmail", "user2@example.com").Return(u, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(gorm_user.MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.AuthenticateUser(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.req.Email, resp.Email)
				require.NotEmpty(t, resp.Token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	oldHash, _ := HashPassword("oldpass")

	tests := []struct {
		name      string
		req       *dto.UpdatePasswordRequest
		mockSetup func(m *gorm_user.MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req: &dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "oldpass",
				NewPassword: "newpass123",
			},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				u := &user.User{ID: 1, Email: "x@example.com", Password: oldHash}
				m.On("GetByID", 1).Return(u, nil).Once()
				m.On("Update", mock.Anything).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "invalid old password",
			req: &dto.UpdatePasswordRequest{
				ID:          2,
				OldPassword: "wrongpass",
				NewPassword: "newpass",
			},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				u := &user.User{ID: 2, Email: "y@example.com", Password: oldHash}
				m.On("GetByID", 2).Return(u, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(gorm_user.MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.UpdatePassword(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, "Password updated successfully", resp.Message)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name      string
		req       *dto.DeleteUserRequest
		mockSetup func(m *gorm_user.MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req:  &dto.DeleteUserRequest{ID: 1},
			mockSetup: func(m *gorm_user.MockUserRepository) {
				m.On("Delete", 1).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:    "missing ID",
			req:     &dto.DeleteUserRequest{ID: 0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(gorm_user.MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.DeleteUser(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Equal(t, "User account deleted successfully", resp.Message)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
