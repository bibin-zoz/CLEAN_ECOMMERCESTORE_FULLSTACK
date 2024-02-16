// user_repository_impl.go
package repository

import (
	"cleancode/pkg/entity"
	"errors"
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

func (r *UserRepositoryImpl) FetchUser(Newmail string) (entity.Compare, error) {
	var compare entity.Compare
	fmt.Println(Newmail)
	if err := r.DB.Raw("SELECT ID, password, username,email, role, status FROM users WHERE email=$1", Newmail).Scan(&compare).Error; err != nil {
		fmt.Println("Error querying the database:", err)
		return compare, err
	}
	fmt.Println("compare", compare)
	return compare, nil
}

func (r *UserRepositoryImpl) FindUserByEmail(user entity.LoginDetail) (entity.User, error) {
	var userDetails entity.User
	err := r.DB.Raw("SELECT * FROM users WHERE email=?", user.Email).Scan(&userDetails).Error
	if err != nil {
		return entity.User{}, errors.New("error checking user details")
	}
	return userDetails, nil
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
	query := `SELECT count(1) FROM "users" WHERE email = ?`
	if err := r.DB.Raw(query, email).Scan(&count).Error; err != nil {
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
	fmt.Println("user", user)
	if err := r.DB.Create(user).Error; err != nil {
		return fmt.Errorf("wrror saving user: %v", err)
	}
	return nil
}
func (r *UserRepositoryImpl) AddAddress(userID int, address entity.UserAddress) error {
	query := "INSERT INTO user_addresses (user_id, street, city, state, postal_code, country, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result := r.DB.Exec(query, userID, address.Street, address.City, address.State, address.PostalCode, address.Country, address.PhoneNumber)

	if result.Error != nil {
		fmt.Println("Error creating address:", result.Error)
		return errors.New("could not add address")
	}

	return nil
}

func (r *UserRepositoryImpl) UserDetails(userID int) (entity.UserDetail, error) {
	var userDetails entity.UserDetail
	err := r.DB.Raw("SELECT u.username,u.email,u.phonenumber FROM users u WHERE u.id = ?", userID).Row().Scan(&userDetails.UserName, &userDetails.Email, &userDetails.PhoneNumber)
	if err != nil {
		return entity.UserDetail{}, errors.New("could not get user details")
	}
	return userDetails, nil
}
func (r *UserRepositoryImpl) GetAllAddress(userId int) ([]entity.AddressInfoResponse, error) {
	var addressInfoResponse []entity.AddressInfoResponse
	if err := r.DB.Raw("SELECT * FROM user_addresses WHERE user_id = ?", userId).Scan(&addressInfoResponse).Error; err != nil {
		return []entity.AddressInfoResponse{}, err
	}
	return addressInfoResponse, nil
}
