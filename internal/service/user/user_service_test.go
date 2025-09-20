package user_service

import (
	"errors"
	"testing"

	"taskflow/internal/domain/user"
	"taskflow/internal/dto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		req       *dto.CreateUserRequest
		mockSetup func(m *MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req:  &dto.CreateUserRequest{Email: "test@example.com", Password: "secret123"},
			mockSetup: func(m *MockUserRepository) {
				m.On("Create", mock.Anything).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req:  &dto.CreateUserRequest{Email: "dup@example.com", Password: "secret123"},
			mockSetup: func(m *MockUserRepository) {
				m.On("Create", mock.Anything).Return(errors.New("duplicate key")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
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
		mockSetup func(m *MockUserRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			req:  &dto.AuthRequest{Email: "user@example.com", Password: "mypassword"},
			mockSetup: func(m *MockUserRepository) {
				u := &user.User{ID: 1, Email: "user@example.com", Password: hashedPass}
				m.On("GetByEmail", "user@example.com").Return(u, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "invalid password",
			req:  &dto.AuthRequest{Email: "user@example.com", Password: "wrongpass"},
			mockSetup: func(m *MockUserRepository) {
				u := &user.User{ID: 1, Email: "user@example.com", Password: hashedPass}
				m.On("GetByEmail", "user@example.com").Return(u, nil).Once()
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name: "non-existent user",
			req:  &dto.AuthRequest{Email: "ghost@example.com", Password: "pass"},
			mockSetup: func(m *MockUserRepository) {
				m.On("GetByEmail", "ghost@example.com").Return((*user.User)(nil), nil).Once()
			},
			wantErr: true,
			errMsg:  "invalid credentials",
		},
		{
			name:      "missing email",
			req:       &dto.AuthRequest{Email: "", Password: "pass"},
			mockSetup: nil,
			wantErr:   true,
			errMsg:    "email is required",
		},
		{
			name:      "missing password",
			req:       &dto.AuthRequest{Email: "user@example.com", Password: ""},
			mockSetup: nil,
			wantErr:   true,
			errMsg:    "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.AuthenticateUser(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
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
		mockSetup func(m *MockUserRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success",
			req: &dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "oldpass",
				NewPassword: "newpass123",
			},
			mockSetup: func(m *MockUserRepository) {
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
			mockSetup: func(m *MockUserRepository) {
				u := &user.User{ID: 2, Email: "y@example.com", Password: oldHash}
				m.On("GetByID", 2).Return(u, nil).Once()
			},
			wantErr: true,
			errMsg:  "invalid old password",
		},
		{
			name: "non-existent user",
			req: &dto.UpdatePasswordRequest{
				ID:          999,
				OldPassword: "any",
				NewPassword: "newpass123",
			},
			mockSetup: func(m *MockUserRepository) {
				// return typed nil to avoid panic
				m.On("GetByID", 999).Return((*user.User)(nil), nil).Once()
			},
			wantErr: true,
			errMsg:  "user not found",
		},
		{
			name: "missing old password",
			req: &dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "",
				NewPassword: "newpass123",
			},
			wantErr: true,
			errMsg:  "old password is required",
		},
		{
			name: "missing new password",
			req: &dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "oldpass",
				NewPassword: "",
			},
			wantErr: true,
			errMsg:  "new password is required",
		},
		{
			name: "new password too short",
			req: &dto.UpdatePasswordRequest{
				ID:          1,
				OldPassword: "oldpass",
				NewPassword: "123",
			},
			wantErr: true,
			errMsg:  "new password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			svc := NewUserService(mockRepo)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			resp, err := svc.UpdatePassword(tt.req)
			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					require.Contains(t, err.Error(), tt.errMsg)
				}
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
		mockSetup func(m *MockUserRepository)
		wantErr   bool
	}{
		{
			name: "success",
			req:  &dto.DeleteUserRequest{ID: 1},
			mockSetup: func(m *MockUserRepository) {
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
			mockRepo := new(MockUserRepository)
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
