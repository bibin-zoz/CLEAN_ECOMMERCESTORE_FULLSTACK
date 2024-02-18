// yourproject/usecase/interfaces/user_usecase.go
package interfaces

import (
	"cleancode/pkg/entity"
)

type UserUseCase interface {
	RegisterUser(user *entity.User) error

	LoginUser(email, password string) (entity.Compare, entity.Invalid, error)
	AddAddress(userID int, address entity.UserAddress) error
	GetAllAddress(userId int) ([]entity.AddressInfoResponse, error)
	GetUserDetails(userId int) (entity.UserDetail, error)
}
