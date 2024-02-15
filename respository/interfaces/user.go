// repository/user_repository.go
package repository

import "cleancode/entity"

type UserRepository interface {
	CheckExistingUsername(username string) (bool, error)
	CheckExistingEmail(email string) (bool, error)
	CheckExistingNumber(number string) (bool, error)
	SaveUser(user *entity.User) error
	FetchUser(Newmail string) (entity.Compare, error)
}
