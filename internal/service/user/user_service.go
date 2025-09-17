package user_service

import "taskflow/internal/repository/gorm/gorm_user"

type UserService struct {
	repo gorm_user.UserRepositoryInterface
}

func NewUserService(repo gorm_user.UserRepositoryInterface) *UserService {

	return &UserService{repo: repo}
}
