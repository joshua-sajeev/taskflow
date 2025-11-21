package dto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"strongpassword123"`
}

type CreateUserResponse struct {
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"john@example.com"`
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"strongpassword123"`
}

type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ID    int    `json:"id" example:"1"`
	Email string `json:"email" example:"john@example.com"`
}

type DeleteUserRequest struct {
	ID int `json:"id" binding:"required" example:"1"`
}

type DeleteUserResponse struct {
	Message string `json:"message" example:"User deleted successfully"`
}

type UpdatePasswordRequest struct {
	ID          int    `json:"id" binding:"required" example:"1"`
	OldPassword string `json:"old_password" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"newsecurepassword456"`
}

type UpdatePasswordResponse struct {
	Message string `json:"message" example:"Password updated successfully"`
}
