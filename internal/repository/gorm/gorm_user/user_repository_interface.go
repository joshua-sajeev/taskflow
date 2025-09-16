package gorm_user

import "taskflow/internal/domain/user"

type UserRepositoryInterface interface {
	Create(user *user.User) error
	Update(user *user.User) error
	Delete(user *user.User) error
}
