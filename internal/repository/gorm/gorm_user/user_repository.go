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
	return r.db.Create(u).Error
}

func (r *UserRepository) Update(u *user.User) error {
	if u.ID == 0 {
		return fmt.Errorf("missing user ID")
	}

	// Check if user exists
	var existing user.User
	if err := r.db.First(&existing, u.ID).Error; err != nil {
		return err // will be gorm.ErrRecordNotFound if not found
	}

	// Update only provided fields
	return r.db.Model(&existing).Updates(u).Error
}

func (r *UserRepository) Delete(id int) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("database not initialized")
	}

	if id == 0 {
		return fmt.Errorf("missing user ID")
	}

	return r.db.Delete(&user.User{}, id).Error
}

func (r *UserRepository) GetByID(id int) (*user.User, error) {
	if id == 0 {
		return nil, fmt.Errorf("missing user ID")
	}

	var u user.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	if email == "" {
		return nil, fmt.Errorf("missing email")
	}

	var u user.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}
