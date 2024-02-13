// user_repository_impl.go
package repository

import (
	"cleancode/entity"
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

func (r *UserRepositoryImpl) CreateUser(user *entity.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserRepositoryImpl) CheckExistingUsername(username string) (bool, error) {
	var count int64
	if err := r.DB.Table("users").Where("username = ?", username).Count(&count).Error; err != nil {
		return false, fmt.Errorf("Error checking existing username: %v", err)
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) CheckExistingEmail(email string) (bool, error) {
	var count int64
	if err := r.DB.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return false, fmt.Errorf("Error checking existing email: %v", err)
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) CheckExistingNumber(number string) (bool, error) {
	var count int64
	if err := r.DB.Table("users").Where("number = ?", number).Count(&count).Error; err != nil {
		return false, fmt.Errorf("Error checking existing number: %v", err)
	}
	return count > 0, nil
}

func (r *UserRepositoryImpl) SaveUser(user *entity.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return fmt.Errorf("Error saving user: %v", err)
	}
	return nil
}
