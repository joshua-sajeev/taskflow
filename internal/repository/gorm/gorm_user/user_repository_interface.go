package gorm_user

import "taskflow/internal/domain/user"

type UserRepositoryInterface interface {
	Create(user *user.User) error
	GetByID(id int) (*user.User, error)
	GetByEmail(email string) (*user.User, error)
	Update(user *user.User) error
	Delete(id int) error
}
