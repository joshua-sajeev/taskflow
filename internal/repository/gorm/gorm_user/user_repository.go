package gorm_user

import (
	"fmt"
	"taskflow/internal/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func (r *UserRepository) Create(u *user.User) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}
	return r.db.Create(u).Error
}

func (r *UserRepository) Update(u *user.User) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}

	if u.ID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	var existing user.User
	if err := r.db.First(&existing, u.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return err
	}

	return r.db.Model(&existing).Updates(u).Error
}

func (r *UserRepository) Delete(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	return r.db.Delete(&user.User{}, id).Error
}

func (r *UserRepository) GetByID(id int) (*user.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	var u user.User
	err := r.db.First(&u, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	var u user.User
	err := r.db.Where("email = ?", email).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Exists(id int) (bool, error) {
	if id <= 0 {
		return false, fmt.Errorf("invalid user ID")
	}

	var u user.User
	err := r.db.First(&u, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
