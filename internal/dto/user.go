package dto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type CreateUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type DeleteUserRequest struct {
	ID int `json:"id" binding:"required"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type UpdatePasswordRequest struct {
	ID          int    `json:"id" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type UpdatePasswordResponse struct {
	Message string `json:"message"`
}
