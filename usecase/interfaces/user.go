// yourproject/usecase/interfaces/user_usecase.go
package interfaces

import (
	"cleancode/entity"
)

type UserUseCase interface {
	RegisterUser(user *entity.User) error

	LoginUser(email, password string) (entity.Compare, entity.Invalid, error)
}
